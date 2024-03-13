# Getting Started

This documentation is intended to guide you through creating a MODM package and offer.

## Prerequisites

* Python (v3.10+)
  * [Python virtual environment (venv)](https://docs.python.org/3/library/venv.html)
* Azure CLI
* Azure Subscription tied to Partner Center account

## Getting Started - Create Offer

For a video walkthrough, see [Packaging Terraform Solutions](https://youtu.be/j-8jTDjU7S4?si=zhr_-XcbGHxPupah)

To get started, these steps will walk you thought:

* Login via `az cli` and install Partner Center Extension
* Clone repo
* Create (or update) Offer and Plan
* Create app package and upload to Plan

---

Prep the Azure CLI by logging into your Partner Center account and installing the [Partner Center Azure CLI Extension](https://github.com/Azure/partnercenter-cli-extension) 
```
az login # Might need: az login --use-device-code

az extension add --source "https://github.com/Azure/partnercenter-cli-extension/releases/download/v0.2.5-alpha/partnercenter-0.2.5-py3-none-any.whl"
```

Clone the Github repo
```
git clone git@github.com:microsoft/commercial-marketplace-offer-deploy.git
cd commercial-marketplace-offer-deploy
```

Create the Offer and Plan (Can be done via CLI or [Partner Center](https://partner.microsoft.com/en-us/dashboard/marketplace-offers/overview))

```
OFFER_ID=<OFFER ID FROM PARTNER CENTER>
OFFER_ALIAS=<OFFER ALIAS FROM PARTNER CENTER>
PLAN_ID=<PLAN ID FOR THE OFFER>
PLAN_NAME=<PLAN NAME FOR THE OFFER>
PLAN_DESCRIPTION="test MODM plan"

# Create the Offer
az partnercenter marketplace offer create --offer-id $OFFER_ID --alias $OFFER_ALIAS --type AzureApplication
OFFER_DESCRIPTION="test MODM offer"

az partnercenter marketplace offer listing update --offer-id $OFFER_ID --summary $OFFER_DESCRIPTION --short-description $OFFER_DESCRIPTION --description $OFFER_DESCRIPTION

az partnercenter marketplace offer list --output table

# Create a Offer Plan
az partnercenter marketplace offer plan create --plan-id $PLAN_ID --name $PLAN_NAME --offer-id $OFFER_ID --subtype managed-application
az partnercenter marketplace offer plan listing update --offer-id $OFFER_ID --plan-id $PLAN_ID --description $PLAN_DESCRIPTION --summary $PLAN_DESCRIPTION 
```

Create the app package and upload to Partner Center (done via CLI AND Partner Center)
```
cd docs
EXAMPLE_DIR=notebooks/complexterraform
az partnercenter marketplace offer package build --id $OFFER_ID --src $EXAMPLE_DIR  --package-type AppInstaller -o tsv --query 'file'

# View contents of the zip file
unzip -l app.zip

```

Now that the app package has been created, do the following in Partner Center:

* Add subscription to [Preview audience](https://learn.microsoft.com/en-us/partner-center/marketplace/azure-app-preview)
* [Upload Technical Config](https://learn.microsoft.com/en-us/partner-center/marketplace/azure-app-technical-configuration)
* [Test and Publish your offer](https://learn.microsoft.com/en-us/partner-center/marketplace/azure-app-test-publish)

## Getting Started - Deploy Offer

For a video walkthrough, see [Installing Published Solutions](https://youtu.be/uA-8PpxexbY?si=7dO80qgTqKQPwxv7).

In Partner Center, once the offer has hit "Publisher signoff" state, use the "Azure marketplace preview" link to deploy into your Azure Subscription.

After the deployment is complete, grab the `applicationInstallerUrl` from the Deployment Inputs and load in your browser.  It will take a few minutes to be available, and the URL will look something like: `https://modmquxaa353lmpr.azurewebsites.net/diagnostics`



## How to Contribute

* Don't be afraid to contribute
* Create Branch (main is protected), Open PR, get approved
* There are automated tests for catching issues, to help you create PR's with confidence
* PREREQ: .NET SDK (only for contribution)

### GH Actions

The project uses the following GH Actions to build and test MODM

TODO: Is this right?
* [Build Service Catalog Managed Applications](../.github/workflows/sccd.yml) - Builds MODM
* [Build](../.github/workflows/ci.yml) - Validates MODM tests
* [Build and Publish Virtual Machine Images](../.github/workflows/cd.yml) - Builds VM Images which MODM uses