import { constants } from "./constants";
import axios, { AxiosRequestConfig, AxiosPromise } from "axios";

export class ServiceHelper {
    public static handleServiceError = (err: any) => {
        if (err && err.response) {
            if (err.response.status === 401) {
                const message = (err.response.data && err.response.data.error && err.response.data.error.message) || "Please make sure the API key is correct.";
                throw new AppError(ErrorCode.HttpStatusUnauthorized, message, "Permission Denied");
            } else if (err.response.status === 404) {
                throw new AppError(
                    ErrorCode.HttpStatusNotFound,
                    "Please make sure the service endpoint is correct.",
                    "Endpoint not found");
            } else if (err.response.status === 429) {
                const response = err.response;
                let errorCode = ErrorCode.Unknown;
                let errorMessage = "";
                let errorTitle = "";
                if (response.data && response.data.error && response.data.error.code === "1014") {
                    errorCode = ErrorCode.ModelCountLimitExceeded;
                    errorMessage = "The number of models associated with the given API key has exceeded the maximum allowed value.";
                    errorTitle = "Too many models";
                } else {
                    errorCode = ErrorCode.HttpStatusTooManyRequests;
                    errorMessage = "We've got too many requests in a short period of time. Please try again later.";
                    errorTitle = "Too many requests";
                }
                throw new AppError(errorCode, errorMessage, errorTitle);
            } else if (err.response.data && err.response.data.error && err.response.data.error.code === "1001") {
                throw new AppError(
                    ErrorCode.ModelNotFound,
                    err.response.data.error.message);
            } else if (err.response.data && err.response.data.error && err.response.data.error.message) {
                throw new AppError(
                    ErrorCode.Unknown,
                    err.response.data.error.message);
            } else {
                throw new AppError(
                    ErrorCode.Unknown,
                    "An error occurred in the service. Please try again later.",
                    "Error");
            }
        } else {
            // Network Error
            throw new AppError(
                ErrorCode.HttpStatusNotFound,
                "Cannot resolve the host name. Please make sure the service endpoint is correct.",
                "Endpoint not found");
        }
    }

    public static postWithAutoRetry = <T>(
        url: string,
        data?: any,
        config?: AxiosRequestConfig,
        apiKey?: string,
    ): AxiosPromise<T> =>
        ServiceHelper.sendRequestWithAutoRetry(() => axios.post(url, data, ServiceHelper.applyApiKey(config, apiKey)))

    public static getWithAutoRetry = <T = any>(
        url: string,
        config?: AxiosRequestConfig,
        apiKey?: string,
    ): AxiosPromise<T> =>
        ServiceHelper.sendRequestWithAutoRetry(() => axios.get(url, ServiceHelper.applyApiKey(config, apiKey)))

    private static applyApiKey = (config?: AxiosRequestConfig, apiKey?: string) => ({
        ...config,
        headers: {
            ...config.headers,
            ...(apiKey ? { [constants.apiKeyHeader]: apiKey } : {}),
        },
    })

    private static sendRequestWithAutoRetry = async <T>(request: () => AxiosPromise<T>) => {
        let currentRetry = 0;
        while (true) {
            try {
                return await request();
            } catch (err) {
                currentRetry++;
                if (currentRetry > constants.maxRetry || !ServiceHelper.isTransient(err)) {
                    throw err;
                }

                await delay(constants.initialRetryInterval * Math.pow(2, currentRetry - 1));
            }
        }
    }

    private static isTransient = (err) => {
        if (err && err.response) {
            const response = err.response;
            if (response.status === 429
                && response.data
                && response.data.error
                && response.data.error.code === "1014") {
                return false;
            }
            return [408, 429, 444, 500, 503, 504].includes(err.response.status);
        }
        return false;
    }
}


export function delay(ms: number) {
    return new Promise<void>((resolve) => {
        setTimeout(() => {
            resolve();
        }, ms);
    });
}

export class AppError extends Error implements IAppError {
    public errorCode: ErrorCode;
    public message: string;
    public title?: string;

    constructor(errorCode: ErrorCode, message: string, title: string = null) {
        super(message);
        this.errorCode = errorCode;
        this.message = message;
        this.title = title;
    }
}

export interface IAppError {
    errorCode: ErrorCode;
    message: any;
    title?: string;
}

export enum ErrorCode {
    // Note that the value of the enum is in camelCase while
    // the enum key is in Pascal casing
    Unknown = "unknown",
    GenericRenderError = "genericRenderError",
    ProjectInvalidJson = "projectInvalidJson",
    ProjectInvalidSecurityToken = "projectInvalidSecurityToken",
    ProjectDuplicateName = "projectDuplicateName",
    SecurityTokenNotFound = "securityTokenNotFound",
    OverloadedKeyBinding = "overloadedKeyBinding",
    BlobContainerIONotFound = "blobContainerIONotFound",
    BlobContainerIOForbidden = "blobContainerIOForbidden",
    PredictWithoutTrainForbidden = "predictWithoutTrainForbidden",
    ModelNotFound = "modelNotFound",
    ModelCountLimitExceeded = "modelCountLimitExceeded",
    HttpStatusUnauthorized = "unauthorized",
    HttpStatusNotFound = "notFound",
    HttpStatusTooManyRequests = "tooManyRequests",
}



