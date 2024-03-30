#!/bin/bash

echo "Configuring project..."

echo "Creating bin folders..."
if [ ! -d "bin" ]; then
  mkdir -p bin
fi
if [ ! -d "bin/configs" ]; then
  mkdir -p bin/configs
fi

echo "Generating configuration files..."
cp configs/server-template-dev.yml configs/server-dev.yml

cp -r configs/*.yml bin/configs

echo "Copying scripts..."
cp scripts/run.sh bin
