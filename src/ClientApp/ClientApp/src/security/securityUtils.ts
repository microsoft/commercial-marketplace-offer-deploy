import { AppConstants } from '../constants/app-constants';

export interface AuthToken {
    expires: string; 
    token: string;
    id: string;
  }
  

  export const validateToken = (): boolean => {
    const tokenString = localStorage.getItem('jwtToken');

    if (!tokenString) {
      return false;
    }

    try {
      const token: AuthToken = JSON.parse(tokenString);

      // Convert the expiration date to a Date object
      const expirationDate = new Date(token.expires);
      const currentTimestamp = new Date();

      // Check if the current date/time is before the expiration date/time
      if (currentTimestamp < expirationDate) {
        return true;
      } else {
        console.log('token expired - currentTimestamp > expirationDate');
        return false;
      }
    } catch (error) {
      console.error('Error parsing the token from local storage:', error);
      return false;
    }
  };
  
  export const login = async (username: string, password: string): Promise<AuthToken> => {
    const backendUrl = AppConstants.baseUrl;
    const response = await fetch(`${backendUrl}/api/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
        },
        body: JSON.stringify({ username, password }), 
      });

      if (!response.ok) {
        throw new Error(`Server responded with ${response.status}: ${response.statusText}`);
      }

    const authToken: AuthToken = await response.json();
    return authToken;
  };

  export const refreshToken = async (sessionId: string): Promise<AuthToken> => {
    const backendUrl = AppConstants.baseUrl;
    const response = await fetch(`${backendUrl}/api/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json',
        },
        body: JSON.stringify({ sessionId }), 
      });

      if (!response.ok) {
        throw new Error(`Server responded with ${response.status}: ${response.statusText}`);
      }

    const authToken: AuthToken = await response.json();
    return authToken;
  };