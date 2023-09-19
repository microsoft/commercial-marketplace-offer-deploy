export const AppConstants = {
    "baseUrl":location.host.indexOf("localhost")==-1? `https://${location.host}`:"http://localhost:5000",
    "apiTimeOut": 650000,
    dateFormat : "MM/DD/YYYY HH:mm:ss"
}

export enum ValidationStatus {
    Validated = 'Validated',
    PendingValidation = 'Pending Validation',
}