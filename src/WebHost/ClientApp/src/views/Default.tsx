import * as React from 'react';
import { TextField } from '@fluentui/react/lib/TextField';
import { Toggle } from '@fluentui/react/lib/Toggle';
import { Announced } from '@fluentui/react/lib/Announced';
import { DetailsList, DetailsListLayoutMode, Selection, SelectionMode, IColumn } from '@fluentui/react/lib/DetailsList';
import { MarqueeSelection } from '@fluentui/react/lib/MarqueeSelection';
import { mergeStyleSets } from '@fluentui/react/lib/Styling';
import { TooltipHost } from '@fluentui/react';
import { AppConstants } from '../constants/app-constants';
import { ProvisionState } from '../constants/deployment.constants';
import { DeploymentResource } from 'models/deployment-models';


export const Default = () => {
  const [deployedResources, setDeployedResources] = React.useState<DeploymentResource[]>([]);
  const columns: IColumn[] = [
    { key: 'name', name: 'Name', fieldName: 'name', minWidth: 100, maxWidth: 200, isResizable: true },
    { key: 'type', name: 'Type', fieldName: 'type', minWidth: 100, maxWidth: 200, isResizable: true },
    { key: 'state', name: 'State', fieldName: 'state', minWidth: 100, maxWidth: 200, isResizable: true },
    { key: 'timestamp', name: 'Timestamp', fieldName: 'timestamp', minWidth: 100, maxWidth: 200, isResizable: true },
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
        type: resource.type,
        state: resource.state,
        timestamp: resource.timestamp
    }));
    setDeployedResources(formattedResources);
    } else {
      console.log('result.deployment.resources is false');
      //setDeployedResources([{ name: "Resource1", type: "Storage Account", state: ProvisionState.SUCCEEDED, timestamp: "9/18/2023" }, { name: "Resource2", type: "Storage Account Container", state: ProvisionState.RUNNING, timestamp: "9/18/2023" }]);
    }

    //console.log(result);
  }


  return (
    <>
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