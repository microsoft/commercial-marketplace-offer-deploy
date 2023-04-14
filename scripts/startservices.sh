#!/bin/bash


# Verify that required env variables are set
if [ -z ${ACME_ACCOUNT_EMAIL} ]; then
  echo "Environment variable ACME_ACCOUNT_EMAIL missing."
  exit 1
fi

if [ -z ${INSTALLER_DOMAIN_NAME} ]; then
  echo "Environment variable INSTALLER_DOMAIN_NAME missing."
  exit 1
fi

function log {
  echo "[$(date -u)] $1"
}

waitForNginxToStop() {
  tail --pid="$NGINX_PID" -f /dev/null
}

stop() {
  log "Caught signal, stopping server and nginx"
  if [[ -n "$SERVER_PID" ]]; then
    kill -TERM "$SERVER_PID" 2> /dev/null
    wait "$SERVER_PID"
  fi
  nginx -s quit
  waitForNginxToStop
}

trap stop EXIT

log "Creating required directories and links..."
# Create folders for nginx and certs files. These might exist already in case of container restart
mkdir -p /installerstore/nginx
mkdir -p /installerstore/.acme.sh
mkdir -p /etc/letsencrypt/live/aapinstaller
# Create a soft link to the directory in persistent volume if it does not exist yet
test ! -d ~/.acme.sh && ln -s /installerstore/.acme.sh

# update the server name and redirect URLs in nginx configuration files (if they exist)
log "Updating host name in nginx configuration..."
test -f /etc/nginx/sites-enabled/aapinstaller.http_only.conf && \
  sed -i "s|aapinstaller;|${INSTALLER_DOMAIN_NAME};|g" /etc/nginx/sites-enabled/aapinstaller.http_only.conf
test -f /etc/nginx/sites-enabled/aapinstaller.http_and_https.conf.disabled && \
  sed -i "s|aapinstaller;|${INSTALLER_DOMAIN_NAME};|g" /etc/nginx/sites-enabled/aapinstaller.http_and_https.conf.disabled && \
  sed -i "s|https://aapinstaller|https://${INSTALLER_DOMAIN_NAME}|g" /etc/nginx/sites-enabled/aapinstaller.http_and_https.conf.disabled

# at this point it is ok to start nginx, it will run in background
log "Starting nginx..."
nginx

# start generating dhparam file if not generated yet and let it run in background
DHPARAM_PID=""
if [ ! -f /installerstore/nginx/dhparam.pem ]; then
  log "Starting to generate dhparam file..."
  # following command might take time to run so it will be running in background, its output file is only needed by nginx for HTTPS
  $(openssl dhparam -out /installerstore/nginx/dhparam.pem 2048 2> /dev/null ; cp /installerstore/nginx/dhparam.pem /etc/nginx ) &
  DHPARAM_PID=$!
elif [ ! -f /etc/nginx/dhparam.pem ]; then
  log "File dhparam exists, copying to nginx..."
  cp /installerstore/nginx/dhparam.pem /etc/nginx
else
  log "File dhparam already in nginx."
fi

# install acme.sh if not installed yet
if [ ! -f ~/.acme.sh/.acme.sh.installed. ]; then
  # Run the acme.sh installation step
  log "Installing acme.sh..."
  ./acme.sh --install --email aoc-automation@redhat.com --force --no-color
  RC=$?
  if [ ${RC} -ne 0 ]; then
    log "Failed to install acme.sh script. See previous output for more information."
    exit ${RC}
  fi
  touch ~/.acme.sh/.acme.sh.installed.
else
  log "ACME script already installed."
fi

# request certificate if not obtained yet
if [ ! -f ~/.acme.sh/.acme.sh.certissued. ]; then
  # verify that the requsts to the external host name are being routed to nginx running here
  log "Making requests to http://${INSTALLER_DOMAIN_NAME}/status to verify traffic is reaching nginx..."
  curl -s -o /dev/null \
    --retry-all-errors \
    --retry-connrefused \
    --retry-delay 10 \
    --retry-max-time 600 \
    --retry 60 \
    http://${INSTALLER_DOMAIN_NAME}/status
  RC=$?
  if [ ${RC} -ne 0 ]; then
    log "Failed to get response from http://${INSTALLER_DOMAIN_NAME}/status, aborting."
    exit ${RC}
  fi

  log "Requesting certificate..."
  cd .acme.sh
  ./acme.sh --no-color --set-default-ca --server letsencrypt
  ./acme.sh --no-color --issue --nginx -d ${INSTALLER_DOMAIN_NAME} --log
  RC=$?
  if [ ${RC} -ne 0 ]; then
    log "Failed to issue certificate. See previous output for more information."
    exit ${RC}
  else
    touch ~/.acme.sh/.acme.sh.certissued.
  fi
  cd ~
else
  log "Certificate already issued:"
  ./.acme.sh/acme.sh list
fi

# install certificate if not already installed
if [ ! -f /etc/letsencrypt/live/aapinstaller/.acme.sh.certinstalled. ]; then
  log "Installing certificate..."
  cd .acme.sh
  ./acme.sh --no-color --install-cert \
    -d ${INSTALLER_DOMAIN_NAME} \
    --key-file /etc/letsencrypt/live/aapinstaller/privkey.pem  \
    --fullchain-file /etc/letsencrypt/live/aapinstaller/fullchain.pem \
    --ca-file /etc/letsencrypt/live/aapinstaller/chain.pem
  RC=$?
  if [ ${RC} -ne 0 ]; then
    log "Failed to install certificate. See previous output for more information."
    exit ${RC}
  else
    touch /etc/letsencrypt/live/aapinstaller/.acme.sh.certinstalled.
  fi
  cd ~
else
  log "Certificate already installed."
fi

if [ ! -f /etc/nginx/sites-enabled/.aapinstaller.https.enabled. ]; then
  # Rename configuration files to enable HTTPS
  log "Enabling HTTP and HTTPS nginx configuration files..."
  mv /etc/nginx/sites-enabled/aapinstaller.http_and_https.conf.disabled /etc/nginx/sites-enabled/aapinstaller.http_and_https.conf
  mv /etc/nginx/sites-enabled/aapinstaller.http_only.conf /etc/nginx/sites-enabled/aapinstaller.http_only.conf.disabled
  # wait for dhparam to finish if it was being generated before starting nginx in HTTPS mode
  if [ -n "$DHPARAM_PID" ]; then
    log "Waiting for process generating dhparam file to finish..."
    wait $DHPARAM_PID
  fi
  log "Calling nginx to reload configuration ..."
  nginx -s reload
  touch /etc/nginx/sites-enabled/.aapinstaller.https.enabled.
else
  log "Nginx already configured for HTTPS."
fi


# Start the api server
/apiserver & 

# Start the operator server
/operator &

# Wait for any process to exit
#wait -n

# Exit with status of process that exited first
#exit $?

SERVER_PID=$!
NGINX_PID=$(cat /var/run/nginx.pid)

log "Start-up done, will wait here while nginx is running."
waitForNginxToStop