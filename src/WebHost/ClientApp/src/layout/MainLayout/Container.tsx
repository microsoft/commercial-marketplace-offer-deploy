import { Outlet, useLocation } from "react-router-dom";
import Sidebar from "./Sidebar";

const Container = () => {

  return (
    <div className="container-fluid">
  
      <div className="row">
        <main className="col-md-8 ms-sm-auto col-lg-10 px-lg-5 px-md-5">
          <Outlet />
        </main>
      </div>
    </div>
  )
}

export default Container;