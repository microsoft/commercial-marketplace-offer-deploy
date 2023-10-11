import React from 'react';
import SuccessIcon from './SuccessIcon';

const StyledSuccessIcon = (props: { style: React.CSSProperties | undefined; }) => (
    <span style={props.style}>
      <SuccessIcon />
    </span>
  );

  export default StyledSuccessIcon;