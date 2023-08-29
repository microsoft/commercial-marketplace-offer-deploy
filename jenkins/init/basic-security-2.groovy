#!groovy

import jenkins.model.*
import hudson.security.*
import hudson.util.*;
import jenkins.install.*;

def instance = Jenkins.getInstance()

def hudsonRealm = new HudsonPrivateSecurityRealm(false)
def users = hudsonRealm.getAllUsers()
users_s = users.collect { it.toString() }

// Create the admin user account if it doesn't already exist.
if ("admin" in users_s) {
    println "Admin user already exists - updating password"

    def user = hudson.model.User.get('admin');
    def password = hudson.security.HudsonPrivateSecurityRealm.Details.fromPlainPassword('abcdefg1234')
    user.addProperty(password)
    user.save()
}
else {
    println "--> creating local admin user"

    hudsonRealm.createAccount('admin', 'abcdefg1234')
    instance.setSecurityRealm(hudsonRealm)

    def strategy = new FullControlOnceLoggedInAuthorizationStrategy()
    instance.setAuthorizationStrategy(strategy)
    instance.save()
}

instance.setInstallState(InstallState.INITIAL_SETUP_COMPLETED)