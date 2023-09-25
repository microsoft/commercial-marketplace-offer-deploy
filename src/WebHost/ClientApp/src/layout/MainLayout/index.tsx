import "bootstrap/dist/css/bootstrap.min.css";
import bootstrap from "bootstrap/dist/js/bootstrap.bundle";

import "./index.css"
import { useEffect } from "react";
import Sidebar from "./Sidebar";
import Header from "./Header"
import Container from "./Container";

const MainLayout = () => {
  return (
    <>
      <Header />
      <Container></Container>
    </>)
}

export default MainLayout;