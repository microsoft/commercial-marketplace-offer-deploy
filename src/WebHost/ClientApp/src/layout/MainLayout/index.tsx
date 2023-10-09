import "bootstrap/dist/css/bootstrap.min.css";
import "bootstrap/dist/js/bootstrap.bundle.min";

import "./index.css"
import Sidebar from "./Sidebar";
import Header from "./Header"
import Container from "./Container";
import { initializeIcons, registerIcons } from '@fluentui/react';

//import { initializeIcons, registerIcons } from '@fluentui/react/lib/Icons';

// This will initialize and register the Fluent UI icons.
initializeIcons();

const icons = {
  'GreenCheckCircle': {
    // SVG path for a green circle with a white checkmark
    svg: `
      <svg width="16" height="16" viewBox="0 0 16 16" fill="none" xmlns="http://www.w3.org/2000/svg">
        <circle cx="8" cy="8" r="8" fill="green"/>
        <path d="M12 5L7 11L5 9" stroke="white" stroke-width="1.5"/>
      </svg>
    `,
  },
  'RedExclamationCircle': {
    // SVG path for a red circle with a white exclamation mark
    svg: `
      <svg width="16" height="16" viewBox="0 0 16 16" fill="none" xmlns="http://www.w3.org/2000/svg">
        <circle cx="8" cy="8" r="8" fill="red"/>
        <path d="M8 11V9" stroke="white" stroke-width="1.5"/>
        <circle cx="8" cy="5" r="0.5" fill="white"/>
      </svg>
    `,
  },
};

// Register custom icons
registerIcons({
  icons: {
    ...icons,
  },
});

const MainLayout = () => {
  return (
    <>
    <Header />
    <div className='d-flex'>
    <Sidebar></Sidebar>
    <Container></Container>
    </div>
    </>)
}

export default MainLayout;