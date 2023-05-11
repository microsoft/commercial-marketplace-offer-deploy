# GO Client SDK Usage

## Creating a Deployment

MODM simplifies the deployment semantics of an Azure marketplace managed / package application by acting as an intermediary between Azure.

The client SDK exposes a `Create` and a `Start`, separately. `Create` is where you pass "mainTemplate.json" to it during the creation of a deployment

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

result, err := client.Start(ctx, request)
```

```go
res, err := client.Start(ctx, deploymentId, templateParameters)
```

