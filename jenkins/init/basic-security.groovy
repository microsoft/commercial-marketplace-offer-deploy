#!groovy

import jenkins.model.*
import hudson.security.*
import hudson.util.*;
import jenkins.install.*;
import hudson.security.csrf;

def instance = Jenkins.get()

println "--> creating local user 'admin'"

def securityRealm = new HudsonPrivateSecurityRealm(false)
securityRealm.createAccount('admin','admin')

instance.setSecurityRealm(securityRealm)

def strategy = new FullControlOnceLoggedInAuthorizationStrategy()
strategy.setAllowAnonymousRead(false)
instance.setAuthorizationStrategy(strategy)

instance.setCrumbIssuer(GlobalCrumbIssuerConfiguration.createDefaultCrumbIssuer());

instance.save()