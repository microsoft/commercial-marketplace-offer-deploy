
export type ApiSettings = {
  baseUrl: string
}

export type AppSettings = {
  environment: string | undefined,
  api: ApiSettings,
  auth: { redirectUri: string }
}

const appSettings: AppSettings = {
  environment: process.env.NODE_ENV,
  api: {
    baseUrl: process.env.REACT_APP_API_BASE_URL as string
  },
  auth: {
    redirectUri: process.env.REACT_APP_AUTH_REDIRECT_URI as string
  }
}

export default appSettings;