import { IPersonaSharedProps, Persona, PersonaSize, PersonaPresence, PersonaInitialsColor, IPersonaProps } from '@fluentui/react/lib/Persona';
import { useState } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faBars, faDashboard, faArchive } from '@fortawesome/free-solid-svg-icons'
import { IRenderFunction } from '@fluentui/react';
import { useAuth } from '../../security/AuthContext';

const Header = () => {
  const [renderDetails, updateRenderDetails] = useState(true);
  const { isAuthenticated } = useAuth();
  const onChange = (ev: unknown, checked: boolean | undefined) => {
    updateRenderDetails(!!checked);
  };

  const adminPersona: IPersonaSharedProps = {
    imageInitials: 'MA',
    text: 'MODM Admin',
  };

  const renderPrimaryTextHandler: IRenderFunction<IPersonaProps> = (props) => {
    return <span className='position-absolute text-white' style={{ left: 'auto', right: 50, top: 0 }}>{props?.text}</span>;
  };


  return (
    <header className="navbar navbar-dark sticky-top bg-blue flex-md-nowrap p-0">
       <a className="position-absolute btn-link" role="button"  
       style={{ left: '20px'}}
        data-bs-toggle="offcanvas" data-bs-target="#sidebar" aria-controls="sidebar">
          <FontAwesomeIcon icon={faBars} size="sm" inverse />
        </a>

      <a className="navbar-brand col-md-3 col-lg-2 me-0 px-5 pt-2 pb-2 font-weight-bold" style={{ fontSize: 15, marginLeft: '2em' }} href={'/'}>
        Marketplace Application Installer 
      </a>
    

      <div className="collapse navbar-collapse">
      Marketplace Installer
      </div>
      <div className="navbar-nav d-flex">
        <div className="nav-item text-nowrap">
        {isAuthenticated && (
            <Persona
            {...adminPersona}
            text="MODM Admin"
            size={PersonaSize.size24}
            presence={PersonaPresence.none}
            hidePersonaDetails={!renderDetails}
            initialsColor={PersonaInitialsColor.gray}
            onRenderPrimaryText={renderPrimaryTextHandler}
            color='white'
          />
        )}
        
        </div>
        
      </div>
    </header>
  )
}

export default Header;