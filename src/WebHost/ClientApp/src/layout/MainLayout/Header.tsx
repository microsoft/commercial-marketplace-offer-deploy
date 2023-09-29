import { IPersonaSharedProps, Persona, PersonaSize, PersonaPresence, PersonaInitialsColor } from '@fluentui/react/lib/Persona';
import { useState } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faBars, faDashboard, faArchive } from '@fortawesome/free-solid-svg-icons'

const Header = () => {
  const [renderDetails, updateRenderDetails] = useState(true);
  const onChange = (ev: unknown, checked: boolean | undefined) => {
    updateRenderDetails(!!checked);
  };

  const examplePersona: IPersonaSharedProps = {
    imageInitials: 'AL',
    text: 'Annie Lindqvist',
    secondaryText: 'Software Engineer',
    tertiaryText: 'In a meeting',
    optionalText: 'Available at 4:00pm'
  };

  return (
    <header className="navbar navbar-dark sticky-top bg-blue flex-md-nowrap p-0">
       <a className="position-absolute btn-link" role="button"  
       style={{ left: '20px'}}
        data-bs-toggle="offcanvas" data-bs-target="#sidebar" aria-controls="sidebar">
          <FontAwesomeIcon icon={faBars} size="sm" inverse />
        </a>

      <a className="navbar-brand col-md-3 col-lg-2 me-0 px-5 pt-2 pb-2" style={{ fontSize: 16 }} href="#">
        Azure Marketplace Application Installer 
      </a>
    

      <div className="collapse navbar-collapse">
      Marketplace Installer
      </div>
      <div className="navbar-nav d-flex">
        <div className="nav-item text-nowrap">
        <Persona
          {...examplePersona}
          text="Annie Lindqvist (Available)"
          size={PersonaSize.size24}
          presence={PersonaPresence.online}
          hidePersonaDetails={!renderDetails}
          initialsColor={PersonaInitialsColor.coolGray}
          color='white'
          imageAlt="Annie Lindqvist, status is online"
        />
        </div>
        
      </div>
    </header>
  )
}

export default Header;