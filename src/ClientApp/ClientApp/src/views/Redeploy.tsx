import React from 'react';
import { PrimaryButton, IButtonStyles } from '@fluentui/react/lib/Button';
import { useLocation, useNavigate } from 'react-router-dom';
import { TextField, ITextFieldStyles } from '@fluentui/react/lib/TextField';
import { Stack } from '@fluentui/react/lib/Stack';
import { AppConstants } from '../constants/app-constants';
import { useAuth } from '../security/AuthContext';

const buttonStyles: Partial<IButtonStyles> = {
    root: {
      maxWidth: '200px', 
      margin: '0 auto', 
    }
  };
  
const textFieldStyles: Partial<ITextFieldStyles> = {
    root: {
      width: '100%' 
    }
  };

export const Redeploy = () => {
    const location = useLocation();
    const navigate = useNavigate();
    const queryParams = new URLSearchParams(location.search);
    const deploymentId = queryParams.get('deploymentId') || '0';
    const [loading, setLoading] = React.useState<boolean>(true);
    const [deploymentParams, setDeploymentParams] = React.useState<any | null>(null);
    const { userToken } = useAuth();

    const getDeploymentName = (deploymentId: string) => {
        return `deployment-${deploymentId}`;
    }

    const getAuthHeader = (): HeadersInit | undefined => {
        if (userToken && userToken.token) {
          return {
            'Authorization': `Bearer ${userToken.token}`
          };
        }
    };
    
    const doGetDeploymentParams = async () => {
        try {
            const backendUrl = AppConstants.baseUrl;
            const headers = getAuthHeader();
            console.log(`backendUrl: ${backendUrl}`);
            const response = await fetch(`${backendUrl}/api/deployments/${deploymentId}/parameters`, {
                headers: {
                    Accept: 'application/json',
                    ...headers,
                },
            });
            console.log(`parameters response: ${JSON.stringify(response)}`);
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            const result = await response.json();
            console.log(JSON.stringify(result, null, 2));
            setDeploymentParams(result);
        } catch (error) {
            console.error(error);
            setLoading(false); 
        } finally {
            setLoading(false);
        }
    }
    
    React.useEffect(() => {
        doGetDeploymentParams();
    }, [deploymentId]);

    const handleTextFieldChange = (key: any, event: any) => {
        console.log(`inside handleTextFieldChange with key: ${key} and value: ${event.target.value}`);
        setDeploymentParams({ ...deploymentParams, [key]: event.target.value });
        console.log(`deploymentParams: ${JSON.stringify(deploymentParams)}`);
    };

    const renderDeploymentParams = () => {
        return Object.keys(deploymentParams).map(key => (
            <TextField
                key={key}
                label={key}
                value={deploymentParams[key]}
                onChange={(event) => handleTextFieldChange(key, event)}
            />
        ));
    };

    const getRedeploymentRequest = async() => {
        return {
            "deploymentId": deploymentId,
            "parameters": deploymentParams
        };
    };

    const handleRedeploy = async () => {
        const backendUrl = AppConstants.baseUrl;
        const redeployUrl = `${backendUrl}/api/deployments/redeploy`; 
        console.log(`redeployUrl: ${redeployUrl}`);
        const headers = getAuthHeader();
        try {
            const redeploymentRequest = await getRedeploymentRequest();
            console.log(`redeploymentRequest: ${JSON.stringify(redeploymentRequest)}`);
            const response = await fetch(redeployUrl, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Accept': 'application/json',
                    ...headers,
                },
                body: JSON.stringify(deploymentParams)
            });
    
            if (response.ok) {
                const result = await response.json();
                console.log('Redeployment successful', result);
                navigate('/');
            } else {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
        } catch (error) {
            console.error('Redeployment failed', error);
        }
    };

    return (
        <>
            <div>
                <h1>Redeploy - {getDeploymentName(deploymentId)}</h1>
                {loading ? (
                    <p>Loading deployment parameters...</p>
                ) : (
                    <Stack tokens={{ childrenGap: 15 }}>
                        {renderDeploymentParams()}
                        <PrimaryButton text="Redeploy" onClick={handleRedeploy} styles={buttonStyles} />
                    </Stack>
                )}
            </div>
        </>
    );
}

export default Redeploy;