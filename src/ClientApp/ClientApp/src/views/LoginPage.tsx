import React, { useState, FormEvent } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../security/AuthContext';
import Logo from './azure.png';
import './LoginPage.css'; 

const LoginPage: React.FC = () => {
  const [username, setUsername] = useState<string>('');
  const [password, setPassword] = useState<string>('');
  const [errorMessage, setErrorMessage] = useState<string | null>(null); // New state for error message
  const { login } = useAuth();
  const navigate = useNavigate();

  const handleLogin = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    try {
      await login(username, password);
      navigate('/');
    } catch (error) {
      // Set the error message state if login fails
      setErrorMessage('Login failed. Please check your username and password.');
    }
  };

  return (
    <div className="login-container">
      <div className="login-logo">
        <img src={Logo} alt="Logo" className="logo" /> 
        <h1>Marketplace App Installer</h1>
      </div>
      <div className="login-form">
        <form onSubmit={handleLogin}>
          {errorMessage && <div className="login-error">{errorMessage}</div>} {/* Display error message */}
          <div className="form-group">
            <label htmlFor="username">Username</label>
            <input
              type="text"
              id="username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              placeholder="Enter your username"
            />
          </div>
          <div className="form-group">
            <label htmlFor="password">Password</label>
            <input
              type="password"
              id="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="Enter your password"
            />
          </div>
          <button type="submit" className="login-button">Login</button>
        </form>
        <div className="login-instructions">
        <p>To find your login credentials:</p>
        <ol>
            <li>Go to the <strong>Azure portal</strong> and navigate to your <strong>managed resource group</strong>.</li>
            <li>Find and select the <strong>deployment</strong>.</li>
            <li>In the deployment's left blade, click on the <strong>Outputs tab</strong> to view your credentials.</li>
        </ol>
        </div>
      </div>
    </div>
  );
};

export default LoginPage;
