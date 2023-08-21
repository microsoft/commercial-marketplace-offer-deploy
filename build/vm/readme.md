
## Setup

## Input variables

### Option 1 - environment variables

The following environment variables are required to build the modm vm. export each one prior to building

```
PKR_VAR_client_id
PKR_VAR_client_secret
PKR_VAR_subscription_id
PKR_VAR_tenant_id
```

Create an env file in `./bin` called `.env.packer` and set the values. When executing `./build/vm/deployNewManagedApp.sh` the values will be extracted and set as environment variables for Packer to pick up as input values for the variables defined in `build/vm/modm.pkr.hcl`.

### Option 2 - pkr vars file (todo)

- create a pkr.vars.hcl file and reference when building.