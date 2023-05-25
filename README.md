# ZITADEL Terraform Provider

## A Better Identity and Access Management Solution

ZITADEL combines the best of Auth0 and Keycloak.
It is built for the serverless era.

Learn more about ZITADEL by checking out the [source repository on GitHub](https://github.com/zitadel/zitadel)

## About this Terraform Provider

This is the official ZITADEL Terraform provider.
It lets you declare ZITADEL resources declaratively and apply the resources to any ZITADEL instance.
Importing existing resources to the Terraform state is supported.
Therefore, as long as you have all resources declared, it is even easy to migrate state between instances.

The provider is currently in [Beta state](https://zitadel.com/docs/support/general) and the support level is Enterprise Support.

For general availability the following issues need to be resolved

- [ ] https://github.com/zitadel/terraform-provider-zitadel/issues/85
- [ ] https://github.com/zitadel/terraform-provider-zitadel/issues/92
- [ ] https://github.com/zitadel/terraform-provider-zitadel/issues/91

## Usage

[Follow the Guide in our Docs](https://zitadel.com/docs/guides/manage/terraform/basics).
Note that you need to create an authorized service user to access the ZITADEL APIs through the provider, as noted in the prerequisites.

We don't guarantee that all resources are available in the provider.
In case you miss something you are welcome to [contribute](#contributing).

## Contributing

If you found a bug or want to request a new feature, please open an [issue](https://github.com/zitadel/terraform-provider-zitadel/issues).
Contributions to the provider are very welcome, please follow the general guidance in the [Contribution Guide](https://github.com/zitadel/terraform-provider-zitadel/blob/main/CONTRIBUTING.md).

## Contributors

<a href="https://github.com/zitadel/terraform-provider-zitadel/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=zitadel/terraform-provider-zitadel" />
</a>
