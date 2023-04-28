# Commercial Marketplace Offer Deployment Manager (MODM)

This marketplace offer deployment manager (MODM) simplifies the deployment of complex managed and packaged applications for the Azure Commercial Marketplace.


## How it works
The 

## Feature Overview

- Deployment Dry Run
- Simplified interaction
- Client SDK (Go, C#, Python)

## Contributing

## Developer Setup

Running locally
```
docker compose -f ./deployments/docker-compose.yml up  
```

Minimum Requirements
* Go Version: 1.20.2+
* Docker Version: v4+

Tools Needed
* Ngrok
  * Developer will need to install Ngrok locally
  * [Officlal Getting Started Documentation](https://ngrok.com/docs/using-ngrok-with/go/)
* Setup a .env file in /bin (see ./configs for the template)

## Getting Started

- [Building the Docker image](./docs/docker-image.md)
- [Running locally](./docs/run-locally.md)
- [Client SDK usage (Go)](./docs/sdk-usage-go.md)

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

## Trademarks

This project may contain trademarks or logos for projects, products, or services. Authorized use of Microsoft 
trademarks or logos is subject to and must follow 
[Microsoft's Trademark & Brand Guidelines](https://www.microsoft.com/en-us/legal/intellectualproperty/trademarks/usage/general).
Use of Microsoft trademarks or logos in modified versions of this project must not cause confusion or imply Microsoft sponsorship.
Any use of third-party trademarks or logos are subject to those third-party's policies.
