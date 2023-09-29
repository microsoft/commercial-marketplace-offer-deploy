import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faDashboard, faArchive } from '@fortawesome/free-solid-svg-icons'

import { useState } from "react";
import Offcanvas from "react-bootstrap/Offcanvas";

export interface SidebarProps {
  className?: string | undefined;
}

const Sidebar: React.FC<SidebarProps> = ({ className }) => {
  const [show, setShow] = useState(true);

  return (
    <Offcanvas id="offcanvas" placement="start" backdrop={false} show={show} onHide={() => setShow(false)}>
        <Offcanvas.Header>
          <Offcanvas.Title as="h5">Dashboard</Offcanvas.Title>
        </Offcanvas.Header>
        <Offcanvas.Body >

          <ul className="nav nav-pills flex-column mb-sm-auto mb-0 align-items-start" id="menu">
            <li className="nav-item">
              <a href={'/'} className="nav-link text-truncate">
                <FontAwesomeIcon icon={faDashboard} size='xl' />
                <span className="ms-1 d-none d-sm-inline text-dark">Deployment</span>
              </a>
            </li>
            <li>
              <a href={'/diagnostics'} className="nav-link text-truncate">
                <FontAwesomeIcon icon={faArchive} size='xl' />
                <span className="ms-1 d-none d-sm-inline text-dark ml-2">Diagnostics</span> 
                </a>
            </li>
          </ul>
          
        </Offcanvas.Body>
      </Offcanvas>
  )
}

export default Sidebar;