import * as React from 'react';
import { DetailsList, DetailsListLayoutMode, SelectionMode, IColumn, ConstrainMode } from '@fluentui/react/lib/DetailsList';
import { FocusTrapZone } from '@fluentui/react/lib/FocusTrapZone';
import { Checkbox } from '@fluentui/react/lib/Checkbox';
import { Icon } from '@fluentui/react/lib/Icon';
import { CommandBar, ICommandBarItemProps } from '@fluentui/react/lib/CommandBar';
import { AppConstants } from '../constants/app-constants';
import { DeploymentResource } from '@/models/deployment-models';
import { DeploymentProgressBar } from '@/components/DeploymentProgressBar';
import StyledSuccessIcon from './StyledSuccessIcon';
import StyledFailureIcon from './StyledFailureIcon';
import { toLocalDateTime } from '../utils/DateUtils';
import { Separator } from '@fluentui/react';

export const Default = () => {

  const [filter, setFilter] = React.useState<'All' | 'Succeeded' | 'Failed'>('All');
  const [backendUrl, setBackendUrl] = React.useState<string | null>(null);
  const [offerName, setOfferName] = React.useState<string | null>(null);
  const [deploymentId, setDeploymentId] = React.useState<string | null>(null);
  const [subscriptionId, setSubscriptionId] = React.useState<string | null>(null);
  const [deploymentResourceGroup, setDeploymentResourceGroup] = React.useState<string | null>(null);
  const [deployedResources, setDeployedResources] = React.useState<DeploymentResource[]>([]);
  const [loading, setLoading] = React.useState<boolean>(true);
  const [isHealthy, setIsHealthy] = React.useState(false);
  const [enableFocusTrap, setEnableFocusTrap] = React.useState(false);

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

  const doGetBackendUrl = async () => {
    const response = await fetch(`${AppConstants.baseUrl}/api/settings?key=backendUrl`, {
        headers: {
          Accept: 'application/json',
        },
      });
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const result = await response.json();
      console.log(JSON.stringify(result));

      if (result.backendUrl ) {
        setBackendUrl(result.backendUrl);
      }
  }

  const doGetDeployedResources = async () => {
    try {
     // const backendUrl = AppConstants.baseUrl;
      const response = await fetch(`${backendUrl}/api/Deployments`, {
        headers: {
          Accept: 'application/json',
        },
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const result = await response.json();
      console.log(JSON.stringify(result, null, 2));
      
      if (result.deployment && result.deployment.id !== undefined && result.deployment.id !== null) {
        setDeploymentId(result.deployment.id);
      }

      if (result.deployment && result.deployment.resourceGroup !== undefined && result.deployment.resourceGroup !== null) {
        setDeploymentResourceGroup(result.deployment.resourceGroup);
      }

      if (result.deployment && result.deployment.offerName !== undefined && result.deployment.offerName !== null) {
        setOfferName(result.deployment.offerName);
      }

      if (result.deployment && result.deployment.subscriptionId !== undefined && result.deployment.subscriptionId !== null) {
        setSubscriptionId(result.deployment.subscriptionId);
      }

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
      setLoading(false); // Set loading to false once data is fetched or an error occurred
    }
  }

  React.useEffect(() => {
    const checkEngineHealth = async () => {
        try {

          const response = await fetch(`${backendUrl}/api/status`, {
            headers: {
              Accept: 'application/json',
            },
          });
  
          if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
          }
  
          const statusData = await response.json();
          console.log(JSON.stringify(statusData));
          setIsHealthy(statusData.isHealthy);
          
          if (statusData.isHealthy) {
            console.log('Engine is healthy');
            doGetDeployedResources();
          }
        } catch (error) {
          console.error(error);
        }
      };

    // Start the interval
    const intervalId = setInterval(() => {
        checkEngineHealth();
    }, 5000); // 5 seconds

    // Clear the interval when the component unmounts
    return () => clearInterval(intervalId);
  }, []);

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

  const _items: ICommandBarItemProps[] = [
    {
      key: 'redeploy',
      text: 'Redeploy',
      iconProps: { iconName: 'Upload' },
      onClick: () => console.log('Redeploy clicked'),
    },
    {
        key: 'delete',
        text: 'Delete',
        iconProps: { iconName: 'Delete' }, // Using 'Delete' as the iconName
        onClick: async () => {
          // Here, you can make your API call or any other logic for the delete action
          try {
            const deleteResponse = await fetch(`${backendUrl}/api/resources/${deploymentResourceGroup}/deletemodmresources`, {
              method: 'POST',
            });
            if (!deleteResponse.ok) {
              throw new Error(`HTTP error! status: ${deleteResponse.status}`);
            }
            const deleteResult = await deleteResponse.json();
            console.log(deleteResult);
            // You can also update your component's state or trigger other side effects here if necessary
          } catch (error) {
            console.error("Error deleting:", error);
          }
        },
      }
  ];

  const earliestTimestamp = deployedResources.length > 0
  ? new Date(Math.min(...deployedResources.map(resource => new Date(resource.timestamp).getTime())))
  : null;

  if (!isHealthy) {
    return <h4>Engine is loading...</h4>;
  }

  return (
    <>
      
      <Separator />

      <div style={{ display: 'flex', alignItems: 'center' }}>
    
            {(() => {
              if (loading) return <h4>Installer is Loading...</h4>; 

              const failedCount = deployedResources.filter(r => r.state === "Failed").length;
              const successCount = deployedResources.filter(r => r.state === "Succeeded").length;

              if (failedCount > 0 && (failedCount + successCount) === deployedResources.length) {
                return <h4>{offerName} failed</h4>;
              }
              if (successCount === deployedResources.length) {
                return (
                  <h4 style={{ display: 'inline-flex', alignItems: 'center' }}>
                    <StyledSuccessIcon style={{ marginRight: '4px' }} />
                    <span>{offerName} succeeded</span>
                  </h4>
                );
              }
              return <h4>... {offerName} is in progress</h4>;
            })()}

          <FocusTrapZone style={{ marginTop: '-5px' }}  disabled={!enableFocusTrap}>
                <CommandBar
                  items={_items}
                  ariaLabel="Inbox actions"
                  primaryGroupAriaLabel="Email actions"
                  farItemsGroupAriaLabel="More actions"
                />
          </FocusTrapZone>  
      </div>

      <div className="row mt-3"> {/* Added margin-top for some spacing */}
        <div className="col-md-6">
          <strong>Deployment Id: </strong> {deploymentId}
        </div>
        <div className="col-md-6">
          <strong>Start time: </strong> {earliestTimestamp ? toLocalDateTime(earliestTimestamp.toISOString()) : 'N/A'}
        </div>
      </div>

      <div className="row mt-3">
        <div className="col-md-6">
          <strong>Subscription: </strong> {subscriptionId}
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

      <Separator />
      
      <DetailsList
        items={filteredDeployedResources}
        columns={columns}
        selectionMode={SelectionMode.none}
        layoutMode={DetailsListLayoutMode.fixedColumns}
        constrainMode={ConstrainMode.unconstrained}
      />

    </>
  );
}

export default Default