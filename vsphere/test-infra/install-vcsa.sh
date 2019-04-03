#!/bin/bash
set -e

if [ "$OVFTOOL_URL" == "" ]; then
  echo "Please provide OVFTOOL_URL"
  exit 1
fi
if [ "$VCSA_ISO_URL" == "" ]; then
  echo "Please provide VCSA_ISO_URL"
  exit 1
fi
if [ "$VCSA_TPL_PATH" == "" ]; then
  echo "Please provide VCSA_TPL_PATH"
  exit 1
fi

# Install ovftool
echo "Downloading ovftool ..."
curl -f ${OVFTOOL_URL} -o ./vmware-ovftool.bundle
chmod a+x ./vmware-ovftool.bundle
echo "Installing ovftool ..."
TERM=dumb sudo ./vmware-ovftool.bundle --eulas-agreed

# Install vCenter Server Appliance
MOUNT_LOCATION=/mnt/vcenter
echo "Downloading vCenter Server Appliance ..."
curl -f ${VCSA_ISO_URL} -o ./vmware-vcenter.iso
sudo mkdir $MOUNT_LOCATION
echo "Mounting downloaded VCSA ISO to $MOUNT_LOCATION ..."
sudo mount -o loop ./vmware-vcenter.iso $MOUNT_LOCATION
echo "Installing VCSA ..."
sudo ${MOUNT_LOCATION}/vcsa-cli-installer/lin64/vcsa-deploy install --accept-eula ${VCSA_TPL_PATH}
echo "VCSA installed."
sudo umount $MOUNT_LOCATION
echo "$MOUNT_LOCATION unmounted."
