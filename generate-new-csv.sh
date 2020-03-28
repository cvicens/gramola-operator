#!/bin/sh
. ./settings.sh

operator-sdk generate k8s
operator-sdk generate crds
./bin/openapi-gen --logtostderr=true -o "" -i ./pkg/apis/gramola/v1alpha1 -O zz_generated.openapi -p ./pkg/apis/gramola/v1alpha1 -h ./hack/boilerplate.go.txt -r "-"

operator-sdk generate csv --csv-version ${OPERATOR_VERSION} --from-version "0.0.1" --update-crds