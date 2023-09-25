import { IPersonaSharedProps, Persona, PersonaSize, PersonaPresence, PersonaInitialsColor } from '@fluentui/react/lib/Persona';
import { useState } from 'react';

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
    <header className="navbar navbar-dark sticky-top bg-blue flex-md-nowrap p-0 shadow">
      <a className="navbar-brand col-md-3 col-lg-2 me-0 px-3" href="#">
        Azure
      </a>
      <button className="navbar-toggler position-absolute d-md-none collapsed" type="button" data-bs-toggle="collapse" 
      data-bs-target="#sidebarMenu" aria-controls="sidebarMenu" aria-expanded="false" aria-label="Toggle navigation">
        <span className="navbar-toggler-icon"></span>
      </button>


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