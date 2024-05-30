#!/bin/bash

# To be used as user-data when starting an ec2 instance. This will install
# the utility scripts we need to manage the server.
# https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/user-data.html
POSTGRESQL_VERSION=14

# Install utilities
apt-get update
apt-get -y install gridsite-clients unzip

# Install postgresql
apt-get -y install postgresql-${POSTGRESQL_VERSION}

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
echo ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQCp+THN+mGaR2KpZB/zqlu82SbXGM9N5tVyk9Zq2UWtlgxp+Bf+l72zToOxq8jNCa9lm6toDCTK2x3m7/Gd4IdXmlbTpLBiuwDKwvqQexWvo880GJJ5qnbKWvQIOkAHNFF9huyD+Hk9uYJ3Uif+FphyBJZqMyFpTGzIxzSUUlw+ecMSBoWEwqIwgciE/Ag3lZaxf3QGZioibgFMKlr3sYsE6AKnktuIaYrGMtTni8za9olB8CPh31PO0rKjwOAYI5Pqqy5kSg7U4O95+A1FfX9FHE3Epg/L5JV1+MOQUWAz+u2Dm9Pv8MjAqj6PdfGRZbBC0vgV0ajUzdKccRskRVPGUwxGY9Ryy5/RMZ4Xnaf4IKIKwbBsVgcSwaPPALARHFLVEb7GLwxUIU8CSdkjWSxenwWkdKtDpcoylAtRhgA/ljru8Q1a8pKSu6jcodpQyKtoF08OTIQg0JapSpmya/NBxOb0T5AMF8fuhz2TURjsSjzqtPbA3ocriWVp0e0JcSxZlx5I+6JJmtg7cl5dq4qjN2G6wty7hcb73NArwwj4OsRujAzUzg+NYusPhr5X3dm8XpIsI4MDJ+FGb42GPfhXkCYHr2Cgpj7K2T0pqhv3DiLO64QgI2Sp3XBwm3d0HXOEsetMt7K9Cl8zd2B+uSaNHoGPlWgtwRLBUuzIK3yx5w== ci-ec2 \
  >> /home/ubuntu/.ssh/authorized_keys
