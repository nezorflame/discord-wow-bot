#!/bin/bash
echo $0: creating public and private key files

# Create the .ssh directory
mkdir -p ${HOME}/.ssh
chmod 700 ${HOME}/.ssh

# Create the public and private key files from the environment variables.
echo "${HEROKU_PUBLIC_KEY}" > ${HOME}/.ssh/heroku_id_rsa.pub
chmod 644 ${HOME}/.ssh/heroku_id_rsa.pub

# Note use of double quotes, required to preserve newlines
echo "${HEROKU_PRIVATE_KEY}" > ${HOME}/.ssh/heroku_id_rsa
chmod 600 ${HOME}/.ssh/heroku_id_rsa

# Start the SSH agent and add host
eval `ssh-agent -s`
ssh-add ~/.ssh/heroku_id_rsa
touch ssh-add ~/.ssh/known_hosts
ssh-keygen -R 104.155.2.110

# Copy config from the server
scp ${SSH_USER}@${SSH_ADDRESS}:~/config.toml .
