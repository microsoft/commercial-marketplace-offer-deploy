import React from 'react';
import { PrimaryButton, IButtonStyles } from '@fluentui/react/lib/Button';
import { useLocation } from 'react-router-dom';
import { TextField, ITextFieldStyles } from '@fluentui/react/lib/TextField';
import { Stack } from '@fluentui/react/lib/Stack';
import { AppConstants } from '../constants/app-constants';


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
    const queryParams = new URLSearchParams(location.search);
    const deploymentId = queryParams.get('deploymentId') || '0';
    const [loading, setLoading] = React.useState<boolean>(true);
    const [deploymentParams, setDeploymentParams] = React.useState<any | null>(null);

    const getDeploymentName = (deploymentId: string) => {
        return `deployment-${deploymentId}`;
    }
    
    const doGetDeploymentParams = async () => {
        try {
            const backendUrl = AppConstants.baseUrl;
            console.log(`backendUrl: ${backendUrl}`);
            const response = await fetch(`${backendUrl}/api/deployments/${deploymentId}/parameters`, {
                headers: {
                    Accept: 'application/json',
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

    const handleTextFieldChange = (key, event) => {
        setDeploymentParams({ ...deploymentParams, [key]: event.target.value });
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

    const handleRedeploy = async () => {
        const backendUrl = AppConstants.baseUrl;
        const redeployUrl = `${backendUrl}/api/deployments/${deploymentId}/redeploy`; 
        console.log(`redeployUrl: ${redeployUrl}`);

        try {
            const response = await fetch(redeployUrl, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Accept': 'application/json'
                },
                body: JSON.stringify(deploymentParams)
            });
    
            if (response.ok) {
                const result = await response.json();
                console.log('Redeployment successful', result);
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