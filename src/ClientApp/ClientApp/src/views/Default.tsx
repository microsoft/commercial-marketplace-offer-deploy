import React, { useState, useEffect, useRef } from 'react';
import { DetailsList, DetailsListLayoutMode, SelectionMode, IColumn, ConstrainMode } from '@fluentui/react/lib/DetailsList';
import { Dialog, DialogType, DialogFooter, PrimaryButton, DefaultButton } from '@fluentui/react';
import { FocusTrapZone } from '@fluentui/react/lib/FocusTrapZone';
import { CommandBar, ICommandBarItemProps } from '@fluentui/react/lib/CommandBar';
import { AppConstants } from '../constants/app-constants';
import { DeploymentResource } from '@/models/deployment-models';
import StyledSuccessIcon from './StyledSuccessIcon';
import StyledFailureIcon from './StyledFailureIcon';
import { toLocalDateTime } from '../utils/DateUtils';
import { Separator } from '@fluentui/react';
import { useAuth } from '../security/AuthContext';

export const Default = () => {
  const [filter, setFilter] = React.useState<'All' | 'Succeeded' | 'Failed'>('All');
  const [isConfirmDialogVisible, setIsConfirmDialogVisible] = useState(false);
  const [offerName, setOfferName] = React.useState<string | null>(null);
  const [deploymentId, setDeploymentId] = React.useState<string | null>(null);
  const [deploymentType, setDeploymentType] = React.useState<string | null>(null);
  const [deploymentStatus, setDeploymentStatus] = React.useState<string | null>(null);
  const [subscriptionId, setSubscriptionId] = React.useState<string | null>(null);
  const [deploymentResourceGroup, setDeploymentResourceGroup] = React.useState<string | null>(null);
  const [deployedResources, setDeployedResources] = React.useState<DeploymentResource[]>([]);
  const [loading, setLoading] = React.useState<boolean>(true);
  const [isFinal, setIsFinal] = React.useState<boolean>(false);
  const [isHealthy, setIsHealthy] = React.useState(false);
  const [enableFocusTrap, setEnableFocusTrap] = React.useState(false);
  const { userToken } = useAuth();

  const checkEngineIntervalRef = useRef<number | null>(null);
  const updateResourcesIntervalRef = useRef<number | null>(null);

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
        setIsHealthy(statusData.isHealthy);
        setLoading(false);
    } catch (error) {
        console.error(error);
    }
  };

  const isValidDeploymentType = (deployment) => {
    return deployment?.definition?.deploymentType;
  };

  const isValidDeploymentResourceGroup = (deployment) => {
    return deployment?.resourceGroup ?? null;
  };

  const isValidOfferName = (deployment) => {
    return deployment?.offerName ?? null;
  }

  const isDeploymentFinal = (deployment) => {
    return deployment?.status === "success" || deployment?.status === "failed";
  }

  const isValidSubscriptionId = (deployment) => {
    return deployment?.subscriptionId ?? null;
  }

  const isValidDeploymentStatus = (deployment) => {
    return deployment?.status ?? null;
  }

  const getDeployedResources = async () => {
    try {
        const backendUrl = AppConstants.baseUrl;
        const headers = getAuthHeader();
        const response = await fetch(`${backendUrl}/api/Deployments`, {
          headers: {
            Accept: 'application/json',
            ...headers,
          },
        });
  
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
  
        const result = await response.json();
        
        const deploymentType = isValidDeploymentType(result.deployment);
        if (deploymentType) {
            setDeploymentType(deploymentType);
        }

        const deploymentResourceGroup = isValidDeploymentResourceGroup(result.deployment);
        if (deploymentResourceGroup) {
            setDeploymentResourceGroup(deploymentResourceGroup);
        }
  
        const offerName = isValidOfferName(result.deployment);
        if (offerName) {
            setOfferName(offerName);
        }
        
        const subscriptionId = isValidSubscriptionId(result.deployment);
        if (subscriptionId) {
            setSubscriptionId(subscriptionId);
        }

        const deploymentStatus = isValidDeploymentStatus(result.deployment);
        if (deploymentStatus) {
            setDeploymentStatus(deploymentStatus);
        }

        const isFinal = isDeploymentFinal(result.deployment);
        setIsFinal(isFinal);
          
        if (result.deployment.resources) {
          const formattedResources = result.deployment.resources.map((resource: any) => ({
            name: resource.name,
            state: resource.state,
            type: resource.type,
            timestamp: resource.timestamp
          }));
          setDeployedResources(formattedResources);
        }
      } catch (error) {
        console.error(error);
      } finally {
        setLoading(false); 
      }
  };

  // Function to start checking engine health
  const startEngineHealthCheck = () => {
    checkEngineIntervalRef.current = setInterval(checkEngineHealth, 1000) as unknown as number; // Check every second
  };
  
  // Function to start getting resources
  const startGettingResources = () => {
    getDeployedResources(); // Call immediately
    updateResourcesIntervalRef.current = setInterval(getDeployedResources, 5000) as unknown as number; // Then every 5 seconds
  };
  

  const columns: IColumn[] = [
    { key: 'name', name: 'Name', fieldName: 'name', minWidth: 100, maxWidth: 300, isResizable: true },
    {
      key: 'state',
      name: 'Status',
      minWidth: 50,
      maxWidth: 150,
      isResizable: true,
      onRender: (item: DeploymentResource) => {
        if (item.state === "Succeeded") {
          return <>
            <span>
              <StyledSuccessIcon style={{ marginRight: '4px' }}/>
              {item.state}
            </span>
          </>;
        } else if (item.state === "Failed") {
          return <>
            <span>
              <StyledFailureIcon style={{ marginRight: '4px' }}/>
              {item.state}
            </span>
          </>;
        } else {
          return item.state;
        }
      }
    },
    { key: 'type', name: 'Type', fieldName: 'type', minWidth: 100, maxWidth: 300, isResizable: true },
    {
      key: 'timestamp',
      name: 'Timestamp',
      fieldName: 'timestamp',
      minWidth: 100,
      maxWidth: 300,
      isResizable: true,
      onRender: (item: DeploymentResource) => toLocalDateTime(item.timestamp)
    },
  ];

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
            getDeployedResources();  
            startGettingResources();

            if (checkEngineIntervalRef.current) {
                clearInterval(checkEngineIntervalRef.current);
                checkEngineIntervalRef.current = null;
            }
        }
    
        // Cleanup function to clear the resources interval when the component unmounts
        return () => {
          if (updateResourcesIntervalRef.current) {
            clearInterval(updateResourcesIntervalRef.current);
          }
        };
    }, [isHealthy]);

  const filteredDeployedResources = React.useMemo(() => {
    if (filter === 'All') {
        return deployedResources;
    }
    return deployedResources.filter(r => r.state === filter);
  }, [deployedResources, filter]);

  const onChangeEnableFocusTrap = React.useCallback(
    (ev?: React.FormEvent<HTMLElement | HTMLInputElement> | undefined, checked?: boolean | undefined) => {
      setEnableFocusTrap(!!checked);
    },
    [],
  );

  const handleConfirmDelete = async () => {
    setIsConfirmDialogVisible(false); // Close the dialog
    try {
      const backendUrl = AppConstants.baseUrl;
      const headers = getAuthHeader();
      const deleteResponse = await fetch(`${backendUrl}/api/resources/${deploymentResourceGroup}/deletemodmresources`, {
        method: 'POST',
        headers: {
          Accept: 'application/json',
          ...headers,
        },
      });
      if (!deleteResponse.ok) {
        throw new Error(`HTTP error! status: ${deleteResponse.status}`);
      }
      const deleteResult = await deleteResponse.json();
      // Add any additional logic needed after successful deletion
    } catch (error) {
      console.error("Error deleting:", error);
    }
  };

//   const _items: ICommandBarItemProps[] = [
//     // {
//     //   key: 'redeploy',
//     //   text: 'Redeploy',
//     //   iconProps: { iconName: 'Upload' },
//     //   onClick: () => console.log('Redeploy clicked'),
//     // },
//     {
//         key: 'delete',
//         text: 'Delete Installer',
//         iconProps: { iconName: 'Delete' }, // Using 'Delete' as the iconName
//         onClick: () => setIsConfirmDialogVisible(true),
//     }
//   ];

  const _items = React.useMemo(() => {
    const items = [
        // Include other command items here if needed
    ];

    if (isFinal) {
        items.push({
            key: 'delete',
            text: 'Delete Installer',
            iconProps: { iconName: 'Delete' },
            onClick: () => setIsConfirmDialogVisible(true),
        });
    }

    return items;
  }, [isFinal]); // Re-calculate _items when isFinal changes

  const earliestTimestamp = deployedResources.length > 0
  ? new Date(Math.min(...deployedResources.map(resource => new Date(resource.timestamp).getTime())))
  : null;

  return (
    <>

        <div className='row'>
        <FocusTrapZone disabled={!enableFocusTrap}>
                <CommandBar
                  items={_items}
                  ariaLabel="Inbox actions"
                  primaryGroupAriaLabel="Email actions"
                  farItemsGroupAriaLabel="More actions"
                />
          </FocusTrapZone>  
        </div>
      
      <Separator />

      <div style={{ display: 'flex', alignItems: 'center' }}>
    
            {(() => {
              if (!isHealthy) return <h4>Deployment pending...</h4>; 

              const failedCount = deployedResources.filter(r => r.state === "Failed").length;
              const successCount = deployedResources.filter(r => r.state === "Succeeded").length;

              const failedCountIndicatesFailure = failedCount > 0 && (failedCount + successCount) === deployedResources.length;
              if (failedCountIndicatesFailure || deploymentStatus === "failed") {
                return <h4>{offerName} failed</h4>;
              }
              if (deploymentStatus === "success") {
                return (
                  <h5 style={{ display: 'inline-flex', alignItems: 'center' }}>
                    <StyledSuccessIcon style={{ marginRight: '4px' }} />
                    <span>{offerName} deployment succeeded</span>
                  </h5>
                );
              }
              return <h4>... {offerName} is in progress</h4>;
            })()}
      </div>

      <div className="row mt-3"> {/* Added margin-top for some spacing */}
        <div className="col-md-6">
          <strong>Deployment Type: </strong> {deploymentType}
        </div>
        <div className="col-md-6">
        <strong>Subscription: </strong> {subscriptionId}
        </div>
      </div>

      <div className="row mt-3">
        <div className="col-md-6">
          <strong>Start time: </strong> {earliestTimestamp ? toLocalDateTime(earliestTimestamp.toISOString()) : 'N/A'}
        </div>
        <div className="col-md-6">
          <strong>Resource Group: </strong> {deploymentResourceGroup}
        </div>
      </div>

      <div className="btn-toolbar mb-2 mb-md-0 mt-4">
          <div className="btn-group me-2">
            <div 
              className="alert alert-primary mx-2 p-1 px-2"
              onClick={() => {
                if (filter === 'Succeeded') {
                    setFilter('All');
                } else {
                    setFilter('Succeeded');
                }
            }} 
              style={{cursor: 'pointer'}}
            >
              Success: <strong>{deployedResources.filter(r => r.state == "Succeeded").length}</strong>
            </div>
            <div 
              className="alert alert-danger mx-2 p-1 px-2"
              onClick={() => {
                if (filter === 'Failed') {
                    setFilter('All');
                } else {
                    setFilter('Failed');
                }
            }}
              style={{cursor: 'pointer'}}
            >
              Failures: <strong>{deployedResources.filter(r => r.state == "Failed").length}</strong>
            </div>
          </div>
        </div>

      <div className='row'>
      <DetailsList
        items={filteredDeployedResources}
        columns={columns}
        selectionMode={SelectionMode.none}
        layoutMode={DetailsListLayoutMode.justified}
      />
      </div>

      <Dialog
        hidden={!isConfirmDialogVisible}
        onDismiss={() => setIsConfirmDialogVisible(false)}
        dialogContentProps={{
            type: DialogType.normal,
            title: 'Confirm Deletion',
            subText: 'Are you sure you want to delete the installer?'
        }}>
        <DialogFooter>
            <PrimaryButton onClick={handleConfirmDelete} text="Yes" />
            <DefaultButton onClick={() => setIsConfirmDialogVisible(false)} text="No" />
        </DialogFooter>
      </Dialog>   

    </>
  );
}

export default Default