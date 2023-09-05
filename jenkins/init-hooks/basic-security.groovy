#!groovy

import jenkins.model.*
import hudson.security.*
import hudson.util.*;
import jenkins.install.*;
import jenkins.model.Jenkins
import java.io.File
import hudson.model.FreeStyleProject
import com.cloudbees.hudson.plugins.folder.*

def jenkinsHome = System.getenv('JENKINS_HOME')
def jobConfigXmlPath = '/solutions/job-definitions/terraform.xml'
def jobName = 'modmserviceprincipal'
def cliJarPath = "${jenkinsHome}/jenkins-cli.jar" // Define the path to jenkins-cli.jar
def folderName = "${jenkinsHome}/jobs/${jobName}"

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

def folder = instance.getItem(jobName)
// Create the folder if it doesn't exist
if (folder == null) {
  folder = instance.createProject(Folder.class, jobName)
}

def xmlContent = new File(jobConfigXmlPath).text
def xmlStream = new StringBufferInputStream(xmlContent)
// Check if the job already exists
def job = folder.getItem(jobName)
// Remove it if it already exists
if (job != null) {
  folder.remove(job)
}
// Create job in the folder
job = folder.createProjectFromXML(jobName, xmlStream)

instance.save()