import React from 'react';
import "../../assets/styles/common.scss";
import Table from 'react-bootstrap/Table';
import { Container } from 'react-bootstrap';
import { DeploymentResource } from 'models/deployment-models';
import { DeploymentCardGroup } from '../DeploymentCardGroup/index';
import { ProvisionState } from '../../constants/deployment.constants';

export const Status = () => {
  const [deployedResources, setDeployedResources] = React.useState<DeploymentResource[]>([]);

  React.useEffect(() => {
    doGetDeployedResources();
  }, []);

  const doGetDeployedResources = async () => {
    //const response = await fetch('/api/status');
    setDeployedResources([{ resourceName: "Resource1", resourceStatus: ProvisionState.SUCCEEDED }, { resourceName: "Resource2", resourceStatus: ProvisionState.RUNNING }]);
  };

  const onCardClick = (val: string) => {
    // TODO: Implement
  }

  return (
    <>
        <DeploymentCardGroup 
                data={deployedResources} 
                onClick={onCardClick}/>
                
        <Container fluid className="overview">
            <h2>Deployment Status</h2>
            <div className='container'>
                <Table striped bordered hover>
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
                </Table>
            </div>
        </Container>
    </>
  );
}
