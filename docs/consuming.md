# Azure Marketplace Installer Documentation

## Introduction

Today, the Azure Marketplace primarily supports deployments via ARM templates. This poses a challenge for many publishers who have crafted solutions using different tools, such as Terraform. These solutions align well with the transactional nature of the Azure Marketplace, but the deployment mechanism creates a barrier.

Enter the Marketplace Offer Deployment Module (MODM). MODM is designed to bridge this gap, allowing publishers to bring their existing solutions, even those based on Terraform, to the Azure Marketplace with ease. By using MODM, publishers can seamlessly package their solutions and ensure their compatibility with the Azure Marketplace's deployment mechanisms.

When packaging a solution, publishers will include a `manifest.json` within their solution in a file called `content.zip`. This manifest informs the MODM about the solution type and how to process it. A sample `manifest.json` is provided below:

```json
{
    "deploymentType": "terraform",
    "mainTemplate": "main.tf",
    "offer": {
        "name": "VMOffer",
        "description": "VMOffer just for you!"
    }
}
```

A detailed specification of the `manifest.json` file can be found at [https://github.com/microsoft/commercial-marketplace-offer-deploy/blob/main/docs/manifest-json.md](https://github.com/microsoft/commercial-marketplace-offer-deploy/blob/main/docs/manifest-json.md).

## Packaging Your Solution

1. **Prepare your Solution**: Ensure your solution is in a state ready for deployment.
2. **Packaging**: Compress your solution along with `manifest.json` the  into a file named `content.zip`.
3. **Artifacts Preparation**: Alongside your `content.zip`, ensure you have the `createUiDefinition` and `mainTemplate.json` ready. These are essential for Azure Marketplace deployments. An example of these files can be found at [mainTemplate.json](https://github.com/microsoft/commercial-marketplace-offer-deploy/blob/main/build/managedapp/terraform/complex/mainTemplate.json) and [createUiDefinition.json]([mainTemplate.json](https://github.com/microsoft/commercial-marketplace-offer-deploy/blob/main/build/managedapp/terraform/complex/createUiDefinition.json)).
4. **Final Packaging**: All the above files (`content.zip`, `createUiDefinition`, and `mainTemplate.json`) should be packaged into a single `app.zip` for deployment.

## ARM Template Details

The ARM template provided is the backbone of the MODM solution. It provisions resources in Azure, most notably a virtual machine that will host the deployed solution.

Key aspects of the ARM template:

- **Parameters**: Essential details required for deployment, like `location`, `adminUsername`, `adminPassword`, and `imageReference`.  The value of these parameters are sent from the output of the `createUiDefinition.json` file.  
- **Variables**: Pre-defined values and dynamically constructed strings used throughout the template, like VM size, network configurations, and more.
- **Resources**: The Azure resources that will be provisioned. This includes:
  - Virtual Network and Subnet for networking.
  - Network Security Groups with rules for HTTP and HTTPS.
  - Public IP Address.
  - Virtual Machine with the specified `plan` and `imageReference` which represents the MODM offer.
- **userDataObject**: The content that is sent to the MODM installer virtual machine.
  - **artifactsUri**: The uri of where the `content.zip` can be downloaded.
  - **artifactsHash**: The SHA256 hash of the `content.zip`. This hash is used by the MODM virtual machine to verify the contents of the `content.zip` file has not been tampered with after being uploaded.
  - **parameters**: These are the parameters that will be sent as input to the solution packaged in the `content.zip` file. In the case of a terraform solution, the parameters will be converted into a variables file before exeuting the terraform template.
- **Outputs**: Any value that needs to be returned post-deployment. In this template, it's the `adminUsername`.

## Deployment Process

1. **Deploy your `app.zip`**: You will deploy your `app.zip` either through the Azure Service Catalog or the Partnercenter Marketplace.  Detailed instructions on deploying your `app.zip` can be found at [https://github.com/microsoft/commercial-marketplace-offer-deploy/blob/main/docs/deploy-app-zip.md](https://github.com/microsoft/commercial-marketplace-offer-deploy/blob/main/docs/deploy-app-zip.md).

2. **Deploy using ARM Template**: Use the provided ARM template to deploy. Ensure the `_artifactsLocation` parameter points to the URI where `app.zip` resides.

3. **MODM Execution**: Once deployed, the MODM VM will boot up, and perform the following:

    - Retrieve the `content.zip` from the specified `artifactsUri`  
    - Pull the parameters from the userDataObject via the verify [Azure Instance Metadata Service](https://learn.microsoft.com/en-us/azure/virtual-machines/instance-metadata-service?tabs=linux)  
    - Verify the hash to ensure `content.zip` has not been tampered with
    - Kick off the installation process.

4. **Accessing Your Solution**: Once MODM completes the installation, you should be able to access and interact with your solution as defined in your packaged `content.zip`.

## Conclusion

MODM abstracts away the intricacies of Azure Marketplace deployments, allowing publishers to focus on their solutions. By following the packaging and deployment steps, publishers can have their offerings on the Azure Marketplace with significantly reduced overhead.
