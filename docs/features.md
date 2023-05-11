# MODM Features

## Simplified Deployment semantics

MODM simplifies the deployment semantics of an Azure marketplace managed / package application by acting as an intermediary between Azure.
With the deployment created in MODM, all you have to do is start it with your template parameters:


- Simplified "mainTemplate.json" with a larget set of resources
- Deployment Stages - as nested Azure deployment template equivalent is a Bicep module)
- Automatical retries of a Deployment and/or Stage
- Deployment Dry Run - automatically check for Azure Policies, constraints & quota limits that will prevent the Managed App from successfully being deployed
- Simplified interaction
- Client SDK (Go, C#, Python)