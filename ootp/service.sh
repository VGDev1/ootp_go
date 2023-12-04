#!/bin/bash

# Exit if any command fails
set -e


go build -o "./target/release/" . 

#sudo systemctl stop ootp.service

# Define the source paths for the files
CONFIG_JSON_SRC_PATH="./config.json"          # Replace with actual source path
OOTP_BIN_SRC_PATH="./target/release/oo"          # Replace with actual source path
OOTP_SERVICE_SRC_PATH="./ubuntu/ootp.service" # Replace with actual source path

# Define the destination paths
CONFIG_JSON_DEST_PATH="/home/victor/config.json"
OOTP_BIN_DEST_PATH="/usr/bin/ootp"
OOTP_SERVICE_DEST_PATH="/etc/systemd/system/ootp.service"

# Copy config.json to the user's home directory
echo "Copying config.json to $CONFIG_JSON_DEST_PATH"
cp "$CONFIG_JSON_SRC_PATH" "$CONFIG_JSON_DEST_PATH"

# Copy the ootp binary to /usr/bin
echo "Copying ootp binary to $OOTP_BIN_DEST_PATH"
sudo cp "$OOTP_BIN_SRC_PATH" "$OOTP_BIN_DEST_PATH"
sudo chmod +x "$OOTP_BIN_DEST_PATH"

echo "Copying ootp.service to $OOTP_SERVICE_DEST_PATH"
sudo cp "$OOTP_SERVICE_SRC_PATH" "$OOTP_SERVICE_DEST_PATH"

sudo systemctl start ootp.service

echo "Setup complete."