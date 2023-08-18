#!/bin/sh

set -e

##############
### CONFIG ###
##############

KEYS_DIRECTORY=${KEYS_DIRECTORY:-./keys}

if [ ! -d "${KEYS_DIRECTORY}" ]; then
  echo "Directory ${KEYS_DIRECTORY} does not exist."
  exit 1
fi

if [ ! -f "./config.json" ]; then
  echo "File ./config.json does not exist."
  exit 1
fi

if [ ! -f "${KEYS_DIRECTORY}/system-api-sa.pem" ]; then
  echo "File ${KEYS_DIRECTORY}/system-api-sa.pem does not exist."
  exit 1
fi

if [ ! -f "${KEYS_DIRECTORY}/org-level-admin-sa.json" ]; then
  echo "File ${KEYS_DIRECTORY}/org-level-admin-sa.json does not exist."
  echo "Did ZITADEL set up correctly?"
  exit 1
fi

SYSTEM_API_PEM_KEY=${KEYS_DIRECTORY}/system-api-sa.pem
echo "Using path ${SYSTEM_API_PEM_KEY} to read the system api service account pem key from ."

ORG_LEVEL_DOMAIN=$(cat ./config.json | jq --raw-output '.orgLevel.domain')
INSTANCE_LEVEL_DOMAIN=$(cat ./config.json | jq --raw-output '.instanceLevel.domain')

ORG_LEVEL_KEY=${KEYS_DIRECTORY}/org-level-admin-sa.json
echo "Using path ${ORG_LEVEL_KEY} to read the ${ORG_LEVEL_DOMAIN} instances admin service account key from."

INSTANCE_LEVEL_KEY=${KEYS_DIRECTORY}/instance-level-admin-sa.json
echo "Using path ${INSTANCE_LEVEL_KEY} to write the ${INSTANCE_LEVEL_DOMAIN} instances admin service account key to."

AUDIENCE=${AUDIENCE:-http://$ORG_LEVEL_DOMAIN:8080}
echo "Using audience ${AUDIENCE} for which the key is used."

SERVICE=${SERVICE:-$AUDIENCE}
echo "Using the service ${SERVICE} to connect to ZITADEL. For example in docker compose this can differ from the audience."

######################################
### CREATE INSTANCE LEVEL INSTANCE ###
######################################

echo "Creating SA_JWT for system API user"
SYSTEM_API_TOKEN=$(zitadel-tools key2jwt --key ${SYSTEM_API_PEM_KEY} --audience ${AUDIENCE} --issuer "system-api-sa")

INSTANCE_LEVEL_INSTANCE_PAYLOAD=$(cat <<EOM
{
  "instanceName": "instance-level-tests",
  "customDomain": "${INSTANCE_LEVEL_DOMAIN}",
  "machine": {
    "userName": "instance-level-admin-sa",
    "name": "instance-level-admin-sa",
    "machineKey": {
      "type": "KEY_TYPE_JSON"
    }
  }
}
EOM
)
echo "Creating isolated instance"
echo "${INSTANCE_LEVEL_INSTANCE_PAYLOAD}" | jq

SECOND_INSTANCE_CREATED=$(curl -s --request POST \
  --url ${SERVICE}/system/v1/instances/_create \
  --header 'Content-Type: application/json' \
  --header "Authorization: Bearer ${SYSTEM_API_TOKEN}" \
  --data-raw "${INSTANCE_LEVEL_INSTANCE_PAYLOAD}")
echo "Got response from instance creation:"
echo "${SECOND_INSTANCE_CREATED}" | jq

if [ "$(echo "${SECOND_INSTANCE_CREATED}" | jq --raw-output '.code')" = "6" ]; then
  echo "second instance already exists"
else
  echo "second instance created"

  SECOND_INSTANCE_KEY_DATA=$(echo ${SECOND_INSTANCE_CREATED} | jq --raw-output '.machineKey' | base64 -d | sed 's/\n//g')
  echo "Extracted machine key ${SECOND_INSTANCE_KEY_DATA}"

  echo "Writing second instance key to ${INSTANCE_LEVEL_KEY}"
  echo "${SECOND_INSTANCE_KEY_DATA}" > ${INSTANCE_LEVEL_KEY}
fi

#########################
### HUMAN USERS #########
#########################

create_human_admin () {
  MGMT_AUDIENCE=$1
  MGMT_JSON_KEY=$2

  echo "Creating human admin for ${MGMT_AUDIENCE}"

  AUDIENCE_HOST="$(echo $MGMT_AUDIENCE | cut -d/ -f3)"
  echo "Deferred the Host header ${AUDIENCE_HOST} which will be sent in requests that ZITADEL then maps to a virtual instance"

  SA_JWT=$(zitadel-tools key2jwt --key ${MGMT_JSON_KEY} --audience ${MGMT_AUDIENCE})
  echo "Created JWT from Admin service account key ${SA_JWT}"

  TOKEN_RESPONSE=$(curl -s --request POST \
    --url ${SERVICE}/oauth/v2/token \
    --header 'Content-Type: application/x-www-form-urlencoded' \
    --header "Host: ${AUDIENCE_HOST}" \
    --data grant_type=urn:ietf:params:oauth:grant-type:jwt-bearer \
    --data scope='openid profile email urn:zitadel:iam:org:project:id:zitadel:aud' \
    --data assertion="${SA_JWT}")
  echo "Got response from token endpoint:"
  echo "${TOKEN_RESPONSE}" | jq

  TOKEN=$(echo ${TOKEN_RESPONSE} | jq --raw-output '.access_token')
  echo "Extracted access token ${TOKEN}"

  ORG_RESPONSE=$(curl -s --request GET \
    --url ${SERVICE}/admin/v1/orgs/default \
    --header 'Accept: application/json' \
    --header "Authorization: Bearer ${TOKEN}" \
    --header "Host: ${AUDIENCE_HOST}")
  echo "Got default org response:"
  echo "${ORG_RESPONSE}" | jq

  ORG_ID=$(echo ${ORG_RESPONSE} | jq --raw-output '.org.id')
  echo "Extracted default org id ${ORG_ID}"

  HUMAN_USER_USERNAME="zitadel-admin@zitadel.localhost"
  HUMAN_USER_PASSWORD="Password1!"

  HUMAN_USER_PAYLOAD=$(cat <<EOM
  {
    "userName": "${HUMAN_USER_USERNAME}",
    "profile": {
      "firstName": "ZITADEL",
      "lastName": "Admin",
      "displayName": "ZITADEL Admin",
      "preferredLanguage": "en"
    },
    "email": {
      "email": "zitadel-admin@zitadel.localhost",
      "isEmailVerified": true
    },
    "password": "${HUMAN_USER_PASSWORD}",
    "passwordChangeRequired": false
  }
EOM
  )
  echo "Creating human user"
  echo "${HUMAN_USER_PAYLOAD}" | jq

  HUMAN_USER_RESPONSE=$(curl -s --request POST \
    --url ${SERVICE}/management/v1/users/human/_import \
    --header 'Content-Type: application/json' \
    --header 'Accept: application/json' \
    --header "Authorization: Bearer ${TOKEN}" \
    --header "Host: ${AUDIENCE_HOST}" \
    --data-raw "${HUMAN_USER_PAYLOAD}")
  echo "Create human user response"
  echo "${HUMAN_USER_RESPONSE}" | jq

  if [ "$(echo "${HUMAN_USER_RESPONSE}" | jq --raw-output '.code')" = "6" ]; then
    echo "admin user already exists"
    return
  fi

  HUMAN_USER_ID=$(echo ${HUMAN_USER_RESPONSE} | jq --raw-output '.userId')
  echo "Extracted human user id ${HUMAN_USER_ID}"

  HUMAN_ADMIN_PAYLOAD=$(cat <<EOM
  {
    "userId": "${HUMAN_USER_ID}",
    "roles": [
      "IAM_OWNER"
    ]
  }
EOM
  )
  echo "Granting iam owner to human user"
  echo "${HUMAN_ADMIN_PAYLOAD}" | jq

  HUMAN_ADMIN_RESPONSE=$(curl -s --request POST \
    --url ${SERVICE}/admin/v1/members \
    --header 'Content-Type: application/json' \
    --header 'Accept: application/json' \
    --header "Authorization: Bearer ${TOKEN}" \
    --header "Host: ${AUDIENCE_HOST}" \
    --data-raw "${HUMAN_ADMIN_PAYLOAD}")

  echo "Grant iam owner to human user response"
  echo "${HUMAN_ADMIN_RESPONSE}" | jq
}

ORG_LEVEL_AUDIENCE="http://${ORG_LEVEL_DOMAIN}:8080"
INSTANCE_LEVEL_AUDIENCE="http://${INSTANCE_LEVEL_DOMAIN}:8080"

create_human_admin "${ORG_LEVEL_AUDIENCE}" "$ORG_LEVEL_KEY"
create_human_admin "${INSTANCE_LEVEL_AUDIENCE}" "$INSTANCE_LEVEL_KEY"

GREEN='\033[1;32m'
NC='\033[0m' # No Color

echo
printf "${GREEN}"
echo "Done setting up ZITADEL for acceptance tests"
echo "You can now log in at the following instances:"
printf "${NC}"
echo "${ORG_LEVEL_AUDIENCE}/ui/login"
echo "${INSTANCE_LEVEL_AUDIENCE}/ui/login"
echo
printf "${GREEN}"
echo "For either instance, you can log in with the following credentials:"
printf "${NC}"
echo "username: ${HUMAN_USER_USERNAME}"
echo "password: ${HUMAN_USER_PASSWORD}"
