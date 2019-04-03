#!/bin/bash
set -e

TMP_FILE=$(mktemp)
echo "${SSH_PRIV_KEY}" > $TMP_FILE
echo "private key stored temporarily at "$TMP_FILE

echo "scp'ing file from ${FROM} to ${TO} ..."
scp -o StrictHostKeyChecking=no -i $TMP_FILE $FROM $TO

rm -f $TMP_FILE
echo "private key removed from "$TMP_FILE
