#!/bin/sh
. ./settings.sh

operator-sdk generate k8s
operator-sdk generate crds
./bin/openapi-gen --logtostderr=true -o "" -i ./pkg/apis/gramola/v1alpha1 -O zz_generated.openapi -p ./pkg/apis/gramola/v1alpha1 -h ./hack/boilerplate.go.txt -r "-"

operator-sdk build quay.io/${USERNAME}/${OPERATOR_IMAGE}:${OPERATOR_VERSION}
docker push quay.io/${USERNAME}/${OPERATOR_IMAGE}:${OPERATOR_VERSION}

echo -n "Password for ${USERNAME}: "
read -s PASSWORD 
echo

AUTH_TOKEN=$(curl -sH "Content-Type: application/json" \
-XPOST https://quay.io/cnr/api/v1/users/login \
-d '{"user": {"username": "'"${USERNAME}"'", "password": "'"${PASSWORD}"'"}}' | jq -r '.token')

source ./venv/bin/activate

operator-courier push ./deploy/olm-catalog/gramola-operator cvicensa gramola-operator 0.0.1 "${AUTH_TOKEN}" 

deactivate