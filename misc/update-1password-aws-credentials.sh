#!/bin/bash

# This script is save aws-sso-go results to 1Password.
#
# Usage: aws-sso-go | update-1password-aws-credentials.sh <1password item name>
#   e.g. aws-sso-go | update-1password-aws-credentials.sh aws-credentials 
#
# You can handle AWS credentials by running `op run --env-file .env` with a `.env` file like:
#   AWS_ACCESS_KEY_ID="op://Private/aws-credentials/access key id"
#   AWS_SECRET_ACCESS_KEY="op://Private/aws-credentials/secret access key"
#   AWS_SESSION_TOKEN="op://Private/aws-credentials/session token"

while read -r line; do
    AWS_SECRET_ACCESS_KEY=$(echo $line | jq -r '.SecretAccessKey')
    AWS_ACCESS_KEY_ID=$(echo $line | jq -r '.AccessKeyId')
    AWS_SESSION_TOKEN=$(echo $line | jq -r '.SessionToken')

    op item edit "$1" "secret access key=$AWS_SECRET_ACCESS_KEY" "access key id=$AWS_ACCESS_KEY_ID" "session token=$AWS_SESSION_TOKEN"
done