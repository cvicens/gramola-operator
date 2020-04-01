#!/bin/sh
. ./settings.sh

ACME_OPSRC_POD=$(oc get pods -n openshift-marketplace | grep acme |  awk '{ print $1 }')
oc delete pod ${ACME_OPSRC_POD} -n openshift-marketplace


#sleep 60
#
## This is just to solve a (temporary?) problem with OLM... SA, Role, etc. are deleted when upgrading...
#oc apply -f deploy/service_account.yaml -n gramola-olm-test
#oc apply -f deploy/role.yaml -n gramola-olm-test
#oc apply -f deploy/role_binding.yaml -n gramola-olm-test