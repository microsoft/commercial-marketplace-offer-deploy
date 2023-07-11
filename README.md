# Commercial Marketplace Offer Deployment Manager (MODM)

This marketplace offer deployment manager (MODM) simplifies the deployment of complex managed and packaged applications for the Azure Commercial Marketplace.


## How it works

- Create a `Controller` (some call it a driver--same thing) codebase that consumes the MODM client SDK
- Build your controller binary into a Docker image along with your managed/package app templates
- Create the createUIDefinition.json as usual, but now have the mainTemplate.json represent the deployment of your controller + MODM (See our starter bicep templates)
- Create your marketplace package


<img src="https://github.com/microsoft/commercial-marketplace-offer-deploy/blob/main/docs/img/modm-architecture.png?raw=true" />


## Feature Overview

- Simplified deployment semantics including deployment "Stages" for child Azure deployments
- Automatic retries
- Dry Run operation support
- Async operations built-in
- Web Hook registration to receive only relevant deployment events
- Client SDK (Go, C#, Python) 

A full description of each feature can be found in the [features](./docs/features.md) documentation.


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


### Development Setup

**Prerequisites**
- Go Version: 1.18+
- Docker Version: v4+
- An Azure subscription
- Azure CLI (latest)
- IDE that works with Go
- Ngrok ([Create a free account](https://ngrok.com/))
* Setup a .env file in /bin (see ./configs for the template)

**Developer activities**

- [Building the Docker image](./docs/docker-image.md)
- [Running locally](./docs/run-locally.md)
- [Client SDK usage (Go)](./docs/sdk-usage-go.md)

## Trademarks

This project may contain trademarks or logos for projects, products, or services. Authorized use of Microsoft 
trademarks or logos is subject to and must follow 
[Microsoft's Trademark & Brand Guidelines](https://www.microsoft.com/en-us/legal/intellectualproperty/trademarks/usage/general).
Use of Microsoft trademarks or logos in modified versions of this project must not cause confusion or imply Microsoft sponsorship.
Any use of third-party trademarks or logos are subject to those third-party's policies.

## Credits

- [Kevin Hillinger](https://github.com/kevinhillinger) (Author)
- [Bob Jacobs](https://github.com/bobjac) (Author)
- [Contributors](https://github.com/microsoft/commercial-marketplace-offer-deploy/graphs/contributors)
