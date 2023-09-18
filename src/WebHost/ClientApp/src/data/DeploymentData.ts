import { AxiosRequestConfig } from "axios";
import { AppConstants } from "../constants/app-contants";
import { DeploymentConstants, ProvisionState } from "../constants/deployment.constants";
import { DeploymentResource } from "models/deployment-models";

export const isSuccess = (deploymentResource: DeploymentResource): boolean => {
    return (deploymentResource.resourceStatus == ProvisionState.SUCCEEDED);
}

export const isFailure = (deploymentResource: DeploymentResource): boolean => {
    return (deploymentResource.resourceStatus == ProvisionState.FAILED);
}

export const isRunning = (deploymentResource: DeploymentResource): boolean => {
    return (deploymentResource.resourceStatus == ProvisionState.RUNNING);
};