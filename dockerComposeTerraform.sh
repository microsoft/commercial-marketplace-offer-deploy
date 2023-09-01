#!/bin/bash

java -jar ./jenkins-cli.jar -s http://localhost:8083 create-job modmserviceprincipal3 < ./jenkins/definitions/terraform.xml

java -jar ./jenkins-cli.jar -s http://localhost:8083 build modmserviceprincipal3