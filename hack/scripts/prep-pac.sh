#!/usr/bin/env bash

CLOUD_INIT_DEV="/dev/sr0"
MNT_POINT="/mnt"
HOME_DIR="/root"
AUTH_KEYS_FILE="${HOME_DIR}/.ssh/authorized_keys"
USER_DATA="${MNT_POINT}/openstack/latest/user_data"

if [[ ! -b "$CLOUD_INIT_DEV" ]]; then
    echo "INFO: $CLOUD_INIT_DEV doesn't exist, exiting gracefully."
    exit 0
fi

if mount | grep ${CLOUD_INIT_DEV}; then
    echo "INFO: already found the mounted"
else
    echo "Mounting the ${CLOUD_INIT_DEV}...."
    if mount ${CLOUD_INIT_DEV} ${MNT_POINT}; then
       echo "INFO: successfully mounted ${CLOUD_INIT_DEV} to ${MNT_POINT}"
    else
       echo "ERROR: Failed to mount ${CLOUD_INIT_DEV} to ${MNT_POINT}, hence exiting with error."
        exit 1
    fi
fi

if [[ -f "$USER_DATA" ]]; then
    echo "INFO: $USER_DATA exist."
else
    echo "ERROR: $USER_DATA doesn't exist, hence exiting with error."
    exit 1
fi

mkdir -m 700 -p ${HOME_DIR}/.ssh

if [[ -f "$AUTH_KEYS_FILE" ]]; then
    echo "INFO: $AUTH_KEYS_FILE exists."
else
    install -m 600 /dev/null ${AUTH_KEYS_FILE}
    echo "Created a ${AUTH_KEYS_FILE}"
fi

if grep "$(cat $USER_DATA)" "$AUTH_KEYS_FILE"; then
    echo "INFO: $AUTH_KEYS_FILE already contains the $USER_DATA key, exiting gracefully"
    exit 0
fi

if cat $USER_DATA >> "$AUTH_KEYS_FILE"; then
    echo "INFO: Added the public ssh key to the $AUTH_KEYS_FILE"
else
    echo "ERROR: Fails to add $USER_DATA content to the $AUTH_KEYS_FILE, hence exiting with error"
    exit 1
fi
