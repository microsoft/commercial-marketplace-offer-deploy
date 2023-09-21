import React from 'react';
import { AppConstants } from '../constants/app-constants';
import { useEffect } from 'react';

interface DiagnosticsRespoonse {
  deploymentEngine: string;
}

export const Diagnostics = () => {

  const [diagnostics, setDiagnostics] = React.useState<DiagnosticsRespoonse>({ deploymentEngine: "" });

  useEffect(() => {
    (async () => {
      const backendUrl = AppConstants.baseUrl;
      const response = await fetch(`${backendUrl}/api/diagnostics`, {
        headers: {
          Accept: 'application/json', 
        },
      });

      const result = await response.json();
      setDiagnostics(result);
    })();
  });

  return (<>
  <h5>Deployment Engine Output</h5>
  <div style={{whiteSpace: "pre-wrap", fontSize: 11, borderBottom: '1px solid #ccc', paddingBottom: 5}}>
    {diagnostics.deploymentEngine}
  </div>
  </>)
}

export default Diagnostics;