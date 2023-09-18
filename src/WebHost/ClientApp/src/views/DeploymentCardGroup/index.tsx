import React from 'react';
import Row from 'react-bootstrap/Row';
import Col from 'react-bootstrap/Col';
import Container from 'react-bootstrap/Container';
import './index.scss';
import './deploymentcardgroup.scss';
import overViewIcon from './overview-icon.svg';
import { DeploymentResource } from '../../models/deployment-models';
import { isFailure, isRunning, isSuccess } from '../../data/DeploymentData';
import { OverViewCard } from '../DeploymentOverviewCard/index';
import { ProvisionState } from '../../constants/deployment.constants';

interface Props {
    data: DeploymentResource[] | null;
    onClick(filter: string): any;
}

export const DeploymentCardGroup = ({ data, onClick }: Props) => {
    const sucessfulDeployments = data?.filter(deploymentResource => isSuccess(deploymentResource));
    const failedDeployments = data?.filter(deploymentOperation => isFailure(deploymentOperation));
    const runningDeployments = data?.filter(deploymentOperation => isRunning(deploymentOperation));

    const onCardClick = (val: string) => {
        onClick(val);
    };

    return (
        <Container fluid className="overview">
            <Row className="m-b-20">
                <h2 className="page-header">Deployment Overview</h2>
            </Row>
            
            <div className="overview-container row">   
                <Col>
                    <OverViewCard 
                            selectedFilter="success"
                            value={ProvisionState.SUCCEEDED}
                            onClick={val => onCardClick(val)}
                            title={sucessfulDeployments?.length.toString()}
                            subTitle="Successful Deployments"
                            image={overViewIcon}
                            className=" c-pointer"
                            color="green"
                    ></OverViewCard>
                    
                </Col>
                <Col>
                    <OverViewCard 
                        selectedFilter="failure"
                        value={ProvisionState.FAILED}
                        onClick={val => onCardClick(val)}
                        title={failedDeployments?.length.toString()}
                        subTitle="Failed Deployments"
                        image={overViewIcon}
                        className="c-pointer"
                        color="red"
                    ></OverViewCard>
                </Col>
                
                <Col>
                    <OverViewCard 
                        selectedFilter="running"
                        value={ProvisionState.RUNNING}
                        onClick={val => onCardClick(val)}
                        title={runningDeployments?.length.toString()}
                        subTitle="Running Deployments"
                        image={overViewIcon}
                        className=" c-pointer"
                        color="grey"
                    ></OverViewCard>
                </Col>
            </div>
        </Container>
    );
};