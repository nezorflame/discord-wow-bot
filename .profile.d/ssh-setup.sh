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

# Preload the known_hosts file
echo '|1|nGa/epzYJei5Vohi6/wWqNsiobU=|Tih5nQX6aftXmv7wtVcBHpzK2ZQ= ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBJhx3I5S1scVvmtB9hkfv+tTdT759fUI899fvdjyF7gGq1Bb2xZ2K72gay/iS+a6zUoFw2GYp1dsSqooDhrWTA4=' > ${HOME}/.ssh/known_hosts

# Start the SSH agent and add host
eval `ssh-agent -s`
ssh-add ~/.ssh/heroku_id_rsa

# Copy config from the server
scp ${SSH_USER}@${SSH_ADDRESS}:~/config.toml .
