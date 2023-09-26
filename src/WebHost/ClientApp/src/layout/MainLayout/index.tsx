import "bootstrap/dist/css/bootstrap.min.css";
import "bootstrap/dist/js/bootstrap.bundle.min";

import "./index.css"
import Sidebar from "./Sidebar";
import Header from "./Header"
import Container from "./Container";

const MainLayout = () => {
  return (
    <>
    <Sidebar></Sidebar>
    <Header />
    <Container></Container>

    </>)
}

export default MainLayout;