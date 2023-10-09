import * as React from 'react';
import { DetailsList, DetailsListLayoutMode, SelectionMode, IColumn } from '@fluentui/react/lib/DetailsList';
import { Icon } from '@fluentui/react/lib/Icon';
import { AppConstants } from '../constants/app-constants';
import { DeploymentResource } from '@/models/deployment-models';
import { DeploymentProgressBar } from '@/components/DeploymentProgressBar';
import SuccessIcon from './SuccessIcon';
import FailureIcon from './FailureIcon';
import { toLocalDateTime } from '../utils/DateUtils';

export const Default = () => {
  const [deployedResources, setDeployedResources] = React.useState<DeploymentResource[]>([]);
  const columns: IColumn[] = [
    { key: 'name', name: 'Name', fieldName: 'name', minWidth: 100, maxWidth: 200, isResizable: true },
    {
      key: 'state',
      name: 'State',
      minWidth: 100,
      maxWidth: 200,
      isResizable: true,
      onRender: (item: DeploymentResource) => {
        if (item.state === "Succeeded") {
          return <>
            <SuccessIcon />
            {item.state}
          </>;
        } else if (item.state === "Failed") {
          return <>
            <FailureIcon />
            {item.state}
          </>;
        } else {
          return item.state;
        }
      }
    },
    { key: 'type', name: 'Type', fieldName: 'type', minWidth: 100, maxWidth: 200, isResizable: true },
    {
      key: 'timestamp',
      name: 'Timestamp',
      fieldName: 'timestamp',
      minWidth: 100,
      maxWidth: 200,
      isResizable: true,
      onRender: (item: DeploymentResource) => toLocalDateTime(item.timestamp)
    },
  ];

  React.useEffect(() => {
    doGetDeployedResources();
  }, []);

  const doGetDeployedResources = async () => {
    const backendUrl = AppConstants.baseUrl;
    const response = await fetch(`${backendUrl}/api/Deployments`, {
      headers: {
        Accept: 'application/json',
      },
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const result = await response.json();

    if (result.deployment.resources) {
      console.log('result.deployment.resources is true');
      const formattedResources = result.deployment.resources.map((resource: any) => ({
        name: resource.name,
        state: resource.state,
        type: resource.type,
        timestamp: resource.timestamp
      }));
      console.log(`Got formatted resources ${formattedResources}`);
      setDeployedResources(formattedResources);
    } else {
      console.log('result.deployment.resources is false');
      //setDeployedResources([{ name: "Resource1", type: "Storage Account", state: ProvisionState.SUCCEEDED, timestamp: "9/18/2023" }, { name: "Resource2", type: "Storage Account Container", state: ProvisionState.RUNNING, timestamp: "9/18/2023" }]);
    }
  }


  return (
    <>
      <div className="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3">
        <h1 className="h2">Deployment Details</h1>
        <div className="btn-toolbar mb-2 mb-md-0">
          <div className="btn-group me-2">
            <div className="alert alert-primary mx-2 p-1 px-2">
              Success: <strong>{deployedResources.filter(r => r.state == "Succeeded").length}</strong>
            </div>
            <div className="alert alert-danger mx-2 p-1 px-2">
              Failures: <strong>{deployedResources.filter(r => r.state == "Failed").length}</strong>
            </div>
          </div>
        </div>
      </div>
      <div className='border-bottom pb-5'>
      <DeploymentProgressBar />
      </div>
      <DetailsList
        items={deployedResources}
        columns={columns}
        selectionMode={SelectionMode.none}
        layoutMode={DetailsListLayoutMode.justified}
      />
    </>
  );
}

export default Default