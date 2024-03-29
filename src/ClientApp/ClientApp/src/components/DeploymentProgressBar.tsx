import * as React from 'react';
import { ProgressIndicator } from '@fluentui/react/lib/ProgressIndicator';

const intervalDelay = 10000;
const intervalIncrement = 0.01;

export const DeploymentProgressBar: React.FunctionComponent = () => {
  const [percentComplete, setPercentComplete] = React.useState(0);

  React.useEffect(() => {
    const id = setInterval(() => {
      setPercentComplete((intervalIncrement + percentComplete) % 1);
    }, intervalDelay);
    return () => {
      clearInterval(id);
    };
  });

  return (
    <ProgressIndicator barHeight={10} description="Progress" percentComplete={percentComplete} />
  );
};
