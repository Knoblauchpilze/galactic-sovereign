#!/bin/bash

echo "Configuring project..."

echo "Creating build folders..."
if [ ! -d "bin" ]; then
  mkdir bin
fi

echo "Copying scripts..."
cp scripts/run.sh bin
