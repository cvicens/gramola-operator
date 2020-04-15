#!/bin/sh
. ./settings.sh

echo -n "Password for ${USERNAME}: "
read -s PASSWORD 
echo

AUTH_TOKEN=$(curl -sH "Content-Type: application/json" \
-XPOST https://quay.io/cnr/api/v1/users/login \
-d '{"user": {"username": "'"${USERNAME}"'", "password": "'"${PASSWORD}"'"}}' | jq -r '.token')


PACKAGE_MANIFEST_PATH=./deploy/olm-catalog/gramola-operator/gramola-operator.package.yaml
yq w -i -s update_package_manifest_0.0.1.yaml ${PACKAGE_MANIFEST_PATH}

source ./venv/bin/activate

operator-courier push ./deploy/olm-catalog/gramola-operator cvicensa gramola-operator "0.0.1" "${AUTH_TOKEN}" 

deactivate