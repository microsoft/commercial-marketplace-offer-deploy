import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faDashboard, faArchive } from '@fortawesome/free-solid-svg-icons'

import { useState } from "react";
import Offcanvas from "react-bootstrap/Offcanvas";
import { Link } from 'react-router-dom';

export interface SidebarProps {
  className?: string | undefined;
}

const Sidebar: React.FC<SidebarProps> = ({ className }) => {
  const [show, setShow] = useState(true);

  return (

    <div className="d-flex flex-column flex-shrink-0 bg-light px-2 fixed-top" style={{ zIndex: 1, width: '225px', height: '100%', top: 0, paddingTop: 65 }}>

      <ul className="nav flex-column mb-auto">
        <li className="nav-item p-2">
          <Link to={'/'} className='nav-link text-dark' >
            <FontAwesomeIcon icon={faDashboard} size='xl' className='text-primary mr-1' /> <span className='ml-2'>Dashboard</span>
          </Link>
        </li>
        <li className="nav-item p-2">
          <Link to={'/diagnostics'} className='nav-link text-dark' >
            <FontAwesomeIcon icon={faArchive} size='xl' className='text-primary' /> Diagnostics
          </Link>
        </li>
      </ul>
    </div>
  )
}

export default Sidebar;