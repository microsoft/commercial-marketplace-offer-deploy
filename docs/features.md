# MODM Features

## Simplified Deployment semantics

MODM simplifies the deployment semantics of an Azure marketplace managed / package application by acting as an intermediary between Azure.

The client SDK exposes a `CreateDeployment` and a `StartDeployment` where you pass "mainTemplate.json" to it during the creation of a deployment

```go

// summarized reading of json. ARM template is unmarshaled json to map
template, err := ReadJson(path)

deploymentName := "myDeployment"
request := api.CreateDeployment{
    Name:           &deploymentName,
    Template:       template,
    Location:       &location,
    ResourceGroup:  &resourceGroup,
    SubscriptionID: &subscription,
}

result, err := client.CreateDeployment(ctx, request)
```

With the deployment created in MODM, all you have to do is start it with your template parameters:


```go
res, err := client.StartDeployment(ctx, deploymentId, templateParameters)
```

- Simplified "mainTemplate.json" with a larget set of resources
- Deployment Stages - as nested Azure deployment template equivalent is a Bicep module)
- Automatical retries of a Deployment and/or Stage
- Deployment Dry Run - automatically check for Azure Policies, constraints & quota limits that will prevent the Managed App from successfully being deployed
- Simplified interaction
- Client SDK (Go, C#, Python)