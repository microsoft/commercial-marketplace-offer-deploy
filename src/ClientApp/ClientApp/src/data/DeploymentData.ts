import { AxiosRequestConfig } from "axios";
import { AppConstants } from "../constants/app-constants";
import { DeploymentConstants, ProvisionState } from "../constants/deployment.constants";
import { DeploymentResource } from "../models/deployment-models";

export const isSuccess = (deploymentResource: DeploymentResource): boolean => {
    return (deploymentResource.state == ProvisionState.SUCCEEDED);
}

export const isFailure = (deploymentResource: DeploymentResource): boolean => {
    return (deploymentResource.state == ProvisionState.FAILED);
}

export const isRunning = (deploymentResource: DeploymentResource): boolean => {
    return (deploymentResource.state == ProvisionState.RUNNING);
};


