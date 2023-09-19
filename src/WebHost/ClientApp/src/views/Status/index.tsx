import React from 'react';
// import "../../assets/styles/common.scss";
// import Table from 'react-bootstrap/Table';
import { Container } from 'react-bootstrap';

import { DeploymentResource } from 'models/deployment-models';
import { DeploymentCardGroup } from '../DeploymentCardGroup/index';
import { ProvisionState } from '../../constants/deployment.constants';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Paper from '@mui/material/Paper';
import { AppConstants } from '../../constants/app-contants';


export const Status = () => {
  const [deployedResources, setDeployedResources] = React.useState<DeploymentResource[]>([]);

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
    
  };

  const onCardClick = (val: string) => {
    // TODO: Implement
  }

  return (
    <>
        <DeploymentCardGroup 
                data={deployedResources} 
                onClick={onCardClick}/>
                
        {/* <Container fluid className="overview">
            <h2>Deployment Status</h2>
            <div className='container'> */}
            <TableContainer component={Paper}>
      <Table sx={{ minWidth: 650 }} aria-label="simple table">
        <TableHead>
          <TableRow>
            <TableCell align="right">Resource Name</TableCell>
            <TableCell align="right">Resource Type</TableCell>
            <TableCell align="right">Resource Status</TableCell>
            <TableCell align="right">Timestamp</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {deployedResources.map((row) => (
            <TableRow
              key={row.name}
              sx={{ '&:last-child td, &:last-child th': { border: 0 } }}
            >
              <TableCell align="right">{row.name}</TableCell>
              <TableCell align="right">{row.type}</TableCell>
              <TableCell align="right">{row.state}</TableCell>
              <TableCell align="right">{row.timestamp}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </TableContainer>
                {/* <Table striped bordered hover>
                    <tbody>
                    <tr>
                        <th>Resource</th>
                        <th>Status</th>
                    </tr>
                    {deployedResources?.map((deployment, index) => {
                        return (
                            <tr key={index}>
                                <td>{deployment.resourceName}</td>
                                <td>{deployment.resourceStatus}</td>
                            </tr>
                        )
                    })}          
                    </tbody>
                </Table> */}
            {/* </div>
        </Container> */}
    </>
  );
}
