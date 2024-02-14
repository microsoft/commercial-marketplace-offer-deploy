// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

export const DeploymentConstants = {
    "defaultDeploymentName": "bobjac_test1",
    "defaultResourceGroup": "DIP",
    "defaultRoute": "dashboard/DIP/bobjac_test1",
};

export enum ResourceType {
    DEPLOYMENT = "Microsoft.Resources/deployments"
}

export enum ProvisionState {
    FAILED = "Failed",
    SUCCEEDED = "Succeeded",
    RUNNING = "Running",
}