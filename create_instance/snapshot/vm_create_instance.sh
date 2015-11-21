#!/bin/bash
#
# Creates the compute instance for taking pushable snapshot images.
#
set -x

# The name of instance where skia docs is running on.
INSTANCE_NAME=mathinate-systemd-snapshot-maker
MACHINE_TYPE=n1-standard-4
IMAGE_TYPE="ubuntu-15-04"
IP_ADDRESS=104.154.85.138
PROJECT_ID=heroic-muse-88515
ZONE=us-central1-f

#gcloud compute --project $PROJECT_ID instances create $INSTANCE_NAME \
#  --zone $ZONE \
#  --machine-type $MACHINE_TYPE \
#  --network "default" \
#  --maintenance-policy "MIGRATE" \
#  --tags "http-server,https-server" \
#  --image $IMAGE_TYPE \
#  --boot-disk-type "pd-standard" \
#  --boot-disk-device-name $INSTANCE_NAME \
#  --address=$IP_ADDRESS
#
## Wait until the instance is up.
#until nc -w 1 -z $IP_ADDRESS 22; do
#    echo "Waiting for VM to come up."
#    sleep 2
#done

gcloud compute --project $PROJECT_ID copy-files ./setup-script.sh default@${INSTANCE_NAME}:setup-script.sh \
  --zone $ZONE

gcloud compute --project $PROJECT_ID ssh default@${INSTANCE_NAME} \
  --zone $ZONE \
  --command "sudo bash setup-script.sh"
