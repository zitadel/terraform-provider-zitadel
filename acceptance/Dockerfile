FROM golang:1.19-alpine
RUN apk add curl jq
RUN go install github.com/zitadel/zitadel-tools@v0.4.0
COPY setup.sh config.json /
RUN chmod +x /setup.sh
WORKDIR /
ENTRYPOINT [ "/setup.sh" ]
