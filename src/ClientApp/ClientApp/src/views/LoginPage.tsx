// LoginPage.tsx

import React, { useState, FormEvent } from 'react';
import { login } from '../security/securityutils'; 

const LoginPage: React.FC = () => {
  const [username, setUsername] = useState<string>('');
  const [password, setPassword] = useState<string>('');

  const handleLogin = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    try {
      
      const authData = await login(username, password);
      console.log('Login successful:', authData);
      // TODO: Redirect to the dashboard or wherever you need
    } catch (error) {
      // Handle login error, show user feedback
      console.error('Login failed:', error);
    }
  };

  return (
    <form onSubmit={handleLogin}>
      <input
        type="text"
        value={username}
        onChange={(e) => setUsername(e.target.value)}
        placeholder="Username"
      />
      <input
        type="password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        placeholder="Password"
      />
      <button type="submit">Login</button>
    </form>
  );
};

export default LoginPage;
