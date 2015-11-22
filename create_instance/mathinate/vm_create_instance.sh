#!/bin/bash
#
# Creates the compute instance for mathinate.
#
set -x

PROJECT_ID=heroic-muse-88515
INSTANCE_NAME=mathinate
IP_ADDRESS=104.154.85.138
MACHINE_TYPE=f1-micro
SOURCE_SNAPSHOT=mathinate-systemd-snapshot
SCOPES='https://www.googleapis.com/auth/devstorage.full_control https://www.googleapis.com/auth/compute.readonly'
ZONE=us-central1-f

# Create a boot disk from the pushable base snapshot.
gcloud compute --project $PROJECT_ID disks create $INSTANCE_NAME \
  --zone $ZONE \
  --source-snapshot $SOURCE_SNAPSHOT \
  --type "pd-standard"

gcloud compute --project $PROJECT_ID instances create $INSTANCE_NAME \
  --zone $ZONE \
  --machine-type $MACHINE_TYPE \
  --network "default" \
  --maintenance-policy "MIGRATE" \
  --scopes $SCOPES \
  --tags "http-server" "https-server" \
  --disk "name=mathinate" "device-name=mathinate" "mode=rw" "boot=yes" "auto-delete=yes" \
  --address $IP_ADDRESS
