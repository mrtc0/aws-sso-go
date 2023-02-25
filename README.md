# aws-sso-go

[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/mrtc0/aws-sso-go/badge)](https://api.securityscorecards.dev/projects/github.com/mrtc0/aws-sso-go)
[![CodeQL](https://github.com/mrtc0/aws-sso-go/actions/workflows/codeql.yml/badge.svg)](https://github.com/mrtc0/aws-sso-go/actions/workflows/codeql.yml)

# Motivation

1Password's AWS Shell Plugin is very useful for managing AWS credentials. However, it does not support AWS SSO (`aws sso login`). This project aims to provide a solution to this problem.  
aws-sso-go is output credentials as STDOUT instead of storing them in `~/.aws/sso/cache`. By saving the output to 1Password with a tool like `misc/update-1password-aws-credentials.sh`, you can use `op run --env-file .env` to handle AWS credentials.

# Install

```shell
$ go install github.com/mrtc0/aws-sso-go@latest

$ cat <<EOF > /usr/local/bin/update-1password-aws-credentials.sh
#!/bin/bash

# This script is save aws-sso-go results to 1Password.
#
# Usage: aws-sso-go | update-1password-aws-credentials.sh <1Password item name>
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
EOF

$ chmod +x /usr/local/bin/update-1password-aws-credentials.sh
```

# Usage

```shell
$ aws-sso-go --profile <profile> | update-1password-aws-credentials.sh <1Password item name>
$ cat .env
AWS_ACCESS_KEY_ID="op://Private/<1Password item name>/access key id"
AWS_SECRET_ACCESS_KEY="op://Private/<1Password item name>/secret access key"
AWS_SESSION_TOKEN="op://Private/<1Password item name>/session token"

$ op run --env-file .env -- aws s3 ls
```
