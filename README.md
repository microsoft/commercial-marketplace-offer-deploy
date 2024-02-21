# Commercial Marketplace Offer Deployment Manager (MODM)

[![Build](https://github.com/microsoft/commercial-marketplace-offer-deploy/actions/workflows/ci.yml/badge.svg)](https://github.com/microsoft/commercial-marketplace-offer-deploy/actions/workflows/ci.yml)

The marketplace offer deployment manager (MODM) is an Application Installer for the 
Azure Marketplace that supports deployment of Terraform and Bicep templates.

Supported deployment templates:

- Terraform
- Bicep


## How it works

- Build your Application Package through the Partner Center CLI
- Publish to the Partner Center (using the CLI or through the portal)


<img src="https://github.com/microsoft/commercial-marketplace-offer-deploy/blob/main/docs/img/modm-architecture.png?raw=true" />


## Contributing

This project welcomes contributions and suggestions.  Most contributions require you to agree to a
Contributor License Agreement (CLA) declaring that you have the right to, and actually do, grant us
the rights to use your contribution. For details, visit https://cla.opensource.microsoft.com.

When you submit a pull request, a CLA bot will automatically determine whether you need to provide
a CLA and decorate the PR appropriately (e.g., status check, comment). Simply follow the instructions
provided by the bot. You will only need to do this once across all repos using our CLA.

This project has adopted the [Microsoft Open Source Code of Conduct](https://opensource.microsoft.com/codeofconduct/).
For more information see the [Code of Conduct FAQ](https://opensource.microsoft.com/codeofconduct/faq/) or
contact [opencode@microsoft.com](mailto:opencode@microsoft.com) with any additional questions or comments.

## Walkthrough

There is a Jupyter notebook that contains a walkthrough of the packaging process for MODM located [here](./docs/notebooks/package-terraform-solution.ipynb). 
To execute the the walkthrough, you will need to install the following pre-requisites:

- [Python version 3.10 or higher](https://www.python.org/downloads)
- [Jupyter Notebooks](https://jupyter.org/install)  

Once you have installed the required pre-requisited, you can launch the Jupyter Notebook. More information on launching the notebook can be found at [https://docs.jupyter.org/en/latest](https://docs.jupyter.org/en/latest).

## Video Tutorials

- [Packaging Terraform Solutions](https://youtu.be/j-8jTDjU7S4?si=zhr_-XcbGHxPupah)
- [Installing Published Solutions](https://youtu.be/uA-8PpxexbY?si=7dO80qgTqKQPwxv7)

## Trademarks

This project may contain trademarks or logos for projects, products, or services. Authorized use of Microsoft 
trademarks or logos is subject to and must follow 
[Microsoft's Trademark & Brand Guidelines](https://www.microsoft.com/en-us/legal/intellectualproperty/trademarks/usage/general).
Use of Microsoft trademarks or logos in modified versions of this project must not cause confusion or imply Microsoft sponsorship.
Any use of third-party trademarks or logos are subject to those third-party's policies.

## Credits

- [Kevin Hillinger](https://github.com/kevinhillinger) (Author)
- [Bob Jacobs](https://github.com/bobjac) (Author)
