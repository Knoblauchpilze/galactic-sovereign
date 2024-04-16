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

# Allow ci to access the instance
echo ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQCIBJdCN2m+OWxIAc3GhvJPNJ+zusKXIP37gEcYn3iqRrCFqkbtwn9ebwZiCGZFuMfsnDF/MrwaeJp8rBWSZptX9+iarJIrE9Qd4jnyADN1Jdolkoq3tyuvm3qJKlapTPYdM1noRO8oE14y4pnf9kCvQ7nOfLfErTYIvMmC6LojvrOVBfv4l4VGybFS4fCbDAlaUyaR9RMERLIA16CyDlsfLWijn7/uP6SAL1LP04QfSy2qjUX3NcqQzjzRxkZ7G2oHOg9SYLoFX7gIGHlP5+L1C1kQmH4wBiK+K3cjdWdP/uuFnWR3I9PrqugfhC4LvFafqRXh0CuCOfodKmr9fW6o4IQstOBnOCxc7bqVtrrS9Lhj8As6nOVXowBneI+HWEXyHkSCGVGVOy2vqUriNYohB6hcXiOCKN3ZO3b5ztpD+9FuvlppFJzRO5v/B1twc+h6AbXnD5DkmhtgQo71JiiHpuw2foezPxEQq9zYc8rrtOpobaJp48Um8RMhAQyxistHSZjyP/1yPgWFSmbuzsQtyygn5ioXhqnU/zmGK+IAfBegoF1u0ZcyeSGYWc1iGIcqFAnOYwMscc2k58IlRxtX9P81ltqAb+NiK1aTD1qUA2gjEDrcPBtNu+2wiPrSHGezwIdmONUDlspFGuXhHkTMEyDbuzQv4VmeoveRPPE1iw== ci-ec2 \
  >> /home/ubuntu/.ssh/authorized_keys
