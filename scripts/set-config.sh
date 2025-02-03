#!/bin/bash

# Usage: ./set-config.sh customer1 dev|prod key value
# Example: ./set-config.sh customer1 dev dbPassword "mypassword"

CUSTOMER=$1
ENV=$2
KEY=$3
VALUE=$4

if [ -z "$CUSTOMER" ] || [ -z "$ENV" ] || [ -z "$KEY" ] || [ -z "$VALUE" ]; then
    echo "Usage: ./set-config.sh customer env key value"
    exit 1
fi

cd "customers/$CUSTOMER/$ENV"

if [ "$KEY" == "dbPassword" ] || [ "$KEY" == "apiKey" ]; then
    # Set as secure config
    pulumi config set --secret "$CUSTOMER:$KEY" "$VALUE"
else
    # Set as regular config
    pulumi config set "$CUSTOMER:$KEY" "$VALUE"
fi