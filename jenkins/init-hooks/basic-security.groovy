#!groovy

import jenkins.model.*
import hudson.security.*
import hudson.util.*;
import jenkins.install.*;
import jenkins.model.Jenkins
import java.io.File
import hudson.model.FreeStyleProject
import com.cloudbees.hudson.plugins.folder.*

println "--> creating local user 'admin'"

def instance = Jenkins.get()

def securityRealm = new HudsonPrivateSecurityRealm(false)

def password = System.getenv('DEFAULT_ADMIN_PASSWORD')
println "--> password is set to ${password}"

securityRealm.createAccount('admin', password)

instance.setSecurityRealm(securityRealm)

def strategy = new FullControlOnceLoggedInAuthorizationStrategy()
strategy.setAllowAnonymousRead(false)
instance.setAuthorizationStrategy(strategy)

instance.setCrumbIssuer(hudson.security.csrf.GlobalCrumbIssuerConfiguration.createDefaultCrumbIssuer());

instance.save()