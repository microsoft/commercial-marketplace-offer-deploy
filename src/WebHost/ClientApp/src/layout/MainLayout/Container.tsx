import { Outlet, useLocation } from "react-router-dom";

const Container = () => {

  return (
    <div className="container-fluid">
      <div className="row">
        <main className="col-md-9 ms-sm-auto col-lg-10 px-md-4">
          <Outlet />
        </main>
      </div>
    </div>
  )
}

export default Container;