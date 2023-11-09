import React, { useState, useEffect, useRef } from 'react';
import { AppConstants } from '../constants/app-constants';
import { useAuth } from '../security/AuthContext';

interface DiagnosticsRespoonse {
  deploymentEngine: string;
}

export const Diagnostics = () => {
  const { userToken } = useAuth();
  const [diagnostics, setDiagnostics] = React.useState<DiagnosticsRespoonse>({ deploymentEngine: "" });
  const [isHealthy, setIsHealthy] = React.useState(false);
  
  const checkEngineIntervalRef = useRef<number | null>(null);
  const fetchDiagnosticsIntervalRef = useRef<number | null>(null);
  const fetchCount = useRef(0);

  const getAuthHeader = (): HeadersInit | undefined => {
    if (userToken && userToken.token) {
      return {
        'Authorization': `Bearer ${userToken.token}`
      };
    }
  };

  const checkEngineHealth = async () => {
    try {
        const backendUrl = AppConstants.baseUrl;
        const headers = getAuthHeader();
        console.log(`inside checkEngineHealth with a backendUrl of ${backendUrl}}`);
        const response = await fetch(`${backendUrl}/api/status`, {
          headers: {
            Accept: 'application/json',
            ...headers,
          },
        });

        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }

        const statusData = await response.json();
        console.log(JSON.stringify(statusData));
        setIsHealthy(statusData.isHealthy);
    } catch (error) {
        console.error(error);
    }
  };

  // Function to start checking engine health
  const startEngineHealthCheck = () => {
    checkEngineIntervalRef.current = setInterval(checkEngineHealth, 1000) as unknown as number; // Check every second
  };

//   const fetchDiagnosticsTest = async () => {
//     // Simulate an API call delay
//     await new Promise(resolve => setTimeout(resolve, 1000));

//     // Simulate the API response by appending new text
//     const newLine = `Line ${fetchCount.current + 1}: New diagnostics information received.\n`;
//     setDiagnostics(prevState => ({
//       deploymentEngine: prevState.deploymentEngine + newLine
//     }));

//     fetchCount.current += 1;
//   };

  // Function to start getting resources
  const startGettingDiagnostics = () => {
    fetchDiagnostics(); // Call immediately
    fetchDiagnosticsIntervalRef.current = setInterval(fetchDiagnostics, 5000) as unknown as number; // Then every 5 seconds
  };

  const fetchDiagnostics = async () => {
    const backendUrl = AppConstants.baseUrl;
    const headers = getAuthHeader();
    try {
      const response = await fetch(`${backendUrl}/api/diagnostics`, {
        headers: {
          Accept: 'application/json',
          ...headers,
        },
      });

      if (!response.ok) {
        throw new Error(`Error: ${response.statusText}`);
      }

      const result = await response.json();
      setDiagnostics(result);
    } catch (error) {
      console.error('Failed to fetch diagnostics:', error);
    }
  };

  // Start engine health check interval
  useEffect(() => {
    startEngineHealthCheck(); 

    return () => {
      if (checkEngineIntervalRef.current) {
        clearInterval(checkEngineIntervalRef.current);
      }
    };
  }, []);

  
  useEffect(() => {
    if (isHealthy) {
      fetchDiagnostics(); // Fetch immediately once healthy
      startGettingDiagnostics(); // Start interval

      if (checkEngineIntervalRef.current) {
        clearInterval(checkEngineIntervalRef.current);
        checkEngineIntervalRef.current = null;
      }
    }

    return () => {
      if (fetchDiagnosticsIntervalRef.current) {
        clearInterval(fetchDiagnosticsIntervalRef.current);
      }
    };
  }, [isHealthy]);

  return (<>
    <div className="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <h1 className="h2">Diagnostics</h1>
      <div className="btn-toolbar mb-2 mb-md-0">
        <div className="btn-group me-2">
          <button type="button" className="btn btn-sm btn-outline-secondary">Share</button>
          <button type="button" className="btn btn-sm btn-outline-secondary">Export</button>
        </div>
      </div>
    </div>
    <div style={{ whiteSpace: "pre-wrap", fontSize: 11, paddingBottom: 5 }}>
      {diagnostics.deploymentEngine}
    </div>
  </>)
}

export default Diagnostics;