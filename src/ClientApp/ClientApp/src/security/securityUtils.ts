
export interface AuthToken {
    expires: string; 
    token: string;
    id: string;
  }
  

  export const validateToken = (): boolean => {
    console.log('validateToken');
    const tokenString = localStorage.getItem('jwtToken');
    console.log('tokenString', tokenString);
    if (!tokenString) {
        console.log('!tokenString');
      return false;
    }
    try {
      const token: AuthToken = JSON.parse(tokenString);
      console.log('token', token);
  
      // Convert the expiration date to a Date object
      const expirationDate = new Date(token.expires);
      console.log('expirationDate', expirationDate);
      const currentTimestamp = new Date();
      console.log('currentTimestamp', currentTimestamp);

      // Check if the current date/time is before the expiration date/time
      if (currentTimestamp < expirationDate) {
        console.log('currentTimestamp < expirationDate');
        return true;
      } else {
        console.log('token expired - currentTimestamp > expirationDate');
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
        expires: "2023-11-07T16:46:35.57768+00:00", 
        token: 'eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjVhNGJmMWQ0LTM1NzctNDhjMi1iNWEyLTRmZmQ0NDYwZmI3ZiIsInN1YiI6IkFkbWluaXN0cmF0b3IiLCJqdGkiOiJiYmYzMzM1MS1jMGI5LTQ0NjQtOTViNy0xZjEyNzc1NGEzNmUiLCJuYmYiOjE2OTkzNzM3OTUsImV4cCI6MTY5OTM5MzU5NSwiaWF0IjoxNjk5MzczNzk1LCJpc3MiOiJsb2NhbGhvc3QuYXp1cmV3ZWJzaXRlcy5uZXQvIiwiYXVkIjoibG9jYWxob3N0LmF6dXJld2Vic2l0ZXMubmV0LyJ9.1qi7dBHRS488AW1PhZgChJymazN5I_IeQqgkVX5ZGm5H65j5rLPzBFWkEIra8fcAtUKp_AfuQv-vg9Uf3TIByA', 
        id: '5a4bf1d4-3577-48c2-b5a2-4ffd4460fb7f', 
      };
      console.log('login', fakeToken);
      return fakeToken;
  };

  export const refreshToken = async (userId: string): Promise<AuthToken> => {
    const fakeToken: AuthToken = {
        expires: "2023-11-08T16:46:35.57768+00:00", 
        token: 'eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjVhNGJmMWQ0LTM1NzctNDhjMi1iNWEyLTRmZmQ0NDYwZmI3ZiIsInN1YiI6IkFkbWluaXN0cmF0b3IiLCJqdGkiOiJiYmYzMzM1MS1jMGI5LTQ0NjQtOTViNy0xZjEyNzc1NGEzNmUiLCJuYmYiOjE2OTkzNzM3OTUsImV4cCI6MTY5OTM5MzU5NSwiaWF0IjoxNjk5MzczNzk1LCJpc3MiOiJsb2NhbGhvc3QuYXp1cmV3ZWJzaXRlcy5uZXQvIiwiYXVkIjoibG9jYWxob3N0LmF6dXJld2Vic2l0ZXMubmV0LyJ9.1qi7dBHRS488AW1PhZgChJymazN5I_IeQqgkVX5ZGm5H65j5rLPzBFWkEIra8fcAtUKp_AfuQv-vg9Uf3TIByA', 
        id: '5a4bf1d4-3577-48c2-b5a2-4ffd4460fb7f', 
      };
      console.log('refreshToken', fakeToken);
      return fakeToken;
  };