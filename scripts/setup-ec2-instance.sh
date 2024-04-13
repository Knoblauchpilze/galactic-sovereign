#!/bin/bash

# To be used as user-data when starting an ec2 instance. This will install
# the utility scripts we need to manage the server.
# https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/user-data.html

# Install utilities
apt-get update
apt-get -y install gridsite-clients unzip

# Install postgresql
apt-get -y install postgresql-14

# Install docker
# https://docs.docker.com/engine/install/ubuntu/#install-using-the-repository
apt-get -y install ca-certificates curl
install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
chmod a+r /etc/apt/keyrings/docker.asc

echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
  tee /etc/apt/sources.list.d/docker.list > /dev/null
apt-get update

apt-get -y install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

# Install aws cli
# https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html#getting-started-install-instructions
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
unzip awscliv2.zip
./aws/install
