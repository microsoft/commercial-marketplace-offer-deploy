import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faBars, faDashboard, faArchive } from '@fortawesome/free-solid-svg-icons'

const Sidebar = () => {
  return (
    <>
      <div className="offcanvas offcanvas-start offcanvas-sm-1 sidebar" tabIndex={-1} id="sidebar" 
        data-bs-keyboard="false"  aria-labelledby="sidebarLabel">
        <div className="offcanvas-header navbar-light" id="sidebarLabel">
        <button className="navbar-toggler" type="button" data-bs-toggle="collapse"
          data-bs-dismiss="offcanvas" 
          aria-expanded="false" aria-label="Toggle navigation">
          <FontAwesomeIcon icon={faBars} />
        </button>

        </div>
        <div className="offcanvas-body px-0">
          <ul className="nav nav-pills flex-column mb-sm-auto mb-0 align-items-start" id="menu">
            <li className="nav-item">
              <a href={'/'} className="nav-link text-truncate">
                <FontAwesomeIcon icon={faDashboard} />
                <span className="ms-1 d-none d-sm-inline">Dashboard</span>
              </a>
            </li>
            <li>
              <a href={'/diagnostics'} className="nav-link text-truncate">
                <FontAwesomeIcon icon={faArchive} />
                <span className="ms-1 d-none d-sm-inline">Diagnostics</span> 
                </a>
            </li>
          </ul>
        </div>
      </div>
    </>
  )
}

export default Sidebar;