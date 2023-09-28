# Azure Marketplace Installer Documentation

## Introduction

The Marketplace Offer Deployment Module (MODM) is designed to streamline and automate the deployment process for the Azure Marketplace. Publishers often need to create various artifacts like `mainTemplate.json`, `createUiDefinition`, and package them into `app.zip` to have their offerings on the marketplace.

MODM offers an elegant solution to this challenge. It allows publishers to package existing solutions without getting entangled in the intricacies of marketplace artifact preparation. The essence of MODM is a VM offer that can deploy any solution compatible with Jenkins. However, the underlying Jenkins implementation is abstracted away for simplicity and future flexibility.

## Packaging Your Solution

1. **Prepare your Solution**: Ensure your solution is in a state ready for deployment.
2. **Packaging**: Compress your solution into a file named `content.zip`.
3. **Artifacts Preparation**: Alongside your `content.zip`, ensure you have the `createUiDefinition` and `mainTemplate.json` ready. These are essential for Azure Marketplace deployments.
4. **Final Packaging**: All the above files (`content.zip`, `createUiDefinition`, and `mainTemplate.json`) should be packaged into a single `app.zip` for deployment.

## ARM Template Details

The ARM template provided is the backbone of the MODM solution. It provisions resources in Azure, most notably a virtual machine that will host the deployed solution.

Key aspects of the ARM template:

- **Parameters**: Essential details required for deployment, like `location`, `adminUsername`, `adminPassword`, and `imageReference`.
- **Variables**: Pre-defined values and dynamically constructed strings used throughout the template, like VM size, network configurations, and more.
- **Resources**: The Azure resources that will be provisioned. This includes:
  - Virtual Network and Subnet for networking.
  - Network Security Groups with rules for HTTP and HTTPS.
  - Public IP Address.
  - Virtual Machine with the specified `plan` and `imageReference` which represents the MODM offer.
- **Outputs**: Any value that needs to be returned post-deployment. In this template, it's the `adminUsername`.

## Deployment Process

1. **Upload your `app.zip`**: Ensure your `app.zip` (containing `content.zip`, `createUiDefinition`, and `mainTemplate.json`) is uploaded to a location accessible by Azure.
2. **Deploy using ARM Template**: Use the provided ARM template to deploy. Ensure the `_artifactsLocation` parameter points to the URI where `app.zip` resides.
3. **MODM Execution**: Once deployed, the MODM VM will boot up, retrieve the `content.zip` from the specified `artifactsUri`, and kick off the installation process.
4. **Accessing Your Solution**: Once MODM completes the installation, you should be able to access and interact with your solution as defined in your packaged `content.zip`.

## Conclusion

MODM abstracts away the intricacies of Azure Marketplace deployments, allowing publishers to focus on their solutions. By following the packaging and deployment steps, publishers can have their offerings on the Azure Marketplace with significantly reduced overhead.
