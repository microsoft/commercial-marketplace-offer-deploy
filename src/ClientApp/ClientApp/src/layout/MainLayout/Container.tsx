import { Outlet, useLocation } from "react-router-dom";
import Sidebar from "./Sidebar";

const Container = () => {

  return (
    <div className="container-fluid">
  
      <div className="row">
        <main className="pt-3 px-3 pr-5" style={{ position: 'relative', left: '225px', width: '84%' }}>
          <Outlet />
        </main>
      </div>
    </div>
  )
}

export default Container;