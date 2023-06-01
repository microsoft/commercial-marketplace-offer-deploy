# MarketplaceOfferDeploymentManager

> see https://aka.ms/autorest

This is the AutoRest configuration file for MarketplaceOfferDeploymentManager.

---

## Getting Started

To build the SDK for MarketplaceOfferDeploymentManager, simply [Install AutoRest](https://aka.ms/autorest/install) and in this folder, run:

> `autorest`

To see additional help and options, run:

> `autorest --help`

---

## Configuration

### Basic Information

These are the global settings for the MarketplaceOfferDeploymentManager API.

``` yaml
openapi-type: "data-plane"
data-plane: true
override-client-name: deploymentManagementClient
security: "AADToken"
tag: preview-2023-03-01
```

### Tag: preview-2023-03-01

These settings apply only when `--tag=preview-2023-03-01` is specified on the command line.

```yaml $(tag) == 'preview-2023-03-01'
input-file:
  - 2023-03-01-preview/api.json
```
---

## Go

See configuration in [readme.go.md](./readme.go.md)

## Suppression

``` yaml
# example suppression
# directive:
#   - suppress: R4009
#     from: apimapis.json
#     reason: Warning raised to error while PR was being reviewed. SystemData will implement in next preview version.
```