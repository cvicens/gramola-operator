#!/bin/sh
. ./settings.sh

operator-sdk generate k8s
operator-sdk generate crds
./bin/openapi-gen --logtostderr=true -o "" -i ./pkg/apis/gramola/v1alpha1 -O zz_generated.openapi -p ./pkg/apis/gramola/v1alpha1 -h ./hack/boilerplate.go.txt -r "-"

operator-sdk generate csv --csv-version ${OPERATOR_VERSION} --update-crds

go mod vendor

operator-sdk build quay.io/${USERNAME}/${OPERATOR_IMAGE}:${OPERATOR_VERSION}
docker push quay.io/${USERNAME}/${OPERATOR_IMAGE}:${OPERATOR_VERSION}

echo -n "Password for ${USERNAME}: "
read -s PASSWORD 
echo

AUTH_TOKEN=$(curl -sH "Content-Type: application/json" \
-XPOST https://quay.io/cnr/api/v1/users/login \
-d '{"user": {"username": "'"${USERNAME}"'", "password": "'"${PASSWORD}"'"}}' | jq -r '.token')

source ./venv/bin/activate

operator-courier push ./deploy/olm-catalog/gramola-operator cvicensa gramola-operator ${OPERATOR_VERSION} "${AUTH_TOKEN}" 

deactivate