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