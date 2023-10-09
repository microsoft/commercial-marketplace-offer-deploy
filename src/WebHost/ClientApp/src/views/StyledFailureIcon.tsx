import React from 'react';
import FailureIcon from './FailureIcon';

const StyledFailureIcon = (props: { style: React.CSSProperties | undefined; }) => (
    <span style={props.style}>
      <FailureIcon />
    </span>
  );

export default StyledFailureIcon;