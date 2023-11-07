
export interface AuthToken {
    expires: number; 
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
  
      // Check if the token has expired
      const currentTimestamp = Math.floor(Date.now() / 1000); // Current time in seconds
      if (token.expires && token.expires > currentTimestamp) {
        return true;
      }
  
      // Token is expired
      return false;
    } catch (error) {
      console.error('Error parsing the token from local storage:', error);
      return false;
    }
  };
  
  export const login = async (username: string, password: string): Promise<AuthToken> => {
    const fakeToken: AuthToken = {
        expires: Date.now() + 1000 * 60 * 60, 
        token: 'fake-jwt-token', 
        id: 'fake-id', 
      };

      return fakeToken;
  };