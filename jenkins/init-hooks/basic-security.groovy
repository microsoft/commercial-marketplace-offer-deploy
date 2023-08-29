#!groovy

import jenkins.model.*
import hudson.security.*
import hudson.util.*;
import jenkins.install.*;
import jenkins.model.Jenkins

println "--> creating local user 'admin'"

def instance = Jenkins.get()

def securityRealm = new HudsonPrivateSecurityRealm(false)

def password = System.getenv('DEFAULT_ADMIN_PASSWORD')
securityRealm.createAccount('admin', password)

instance.setSecurityRealm(securityRealm)

def strategy = new FullControlOnceLoggedInAuthorizationStrategy()
strategy.setAllowAnonymousRead(false)
instance.setAuthorizationStrategy(strategy)

instance.setCrumbIssuer(hudson.security.csrf.GlobalCrumbIssuerConfiguration.createDefaultCrumbIssuer());

instance.save()