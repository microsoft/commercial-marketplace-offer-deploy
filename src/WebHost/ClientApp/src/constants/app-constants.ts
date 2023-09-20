// set the AppConstants.  If running in production, it will send http requests back to the hostname that the app is running on.  
// If running in development, it will send http requests to localhost:7258.
// If you are running in a local docker compose instance, change localhost:7258 to whatever is in the docker compose file (localhost:5000)
export const AppConstants = {
    "baseUrl":location.host.indexOf("localhost")==-1? `https://${location.host}`:"https://localhost:7258",
    "apiTimeOut": 650000,
    dateFormat : "MM/DD/YYYY HH:mm:ss"
}

export enum ValidationStatus {
    Validated = 'Validated',
    PendingValidation = 'Pending Validation',
}