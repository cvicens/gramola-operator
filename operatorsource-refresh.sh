#!/bin/sh
. ./settings.sh

ACME_OPSRC_POD=$(oc get pods -n openshift-marketplace | grep acme |  awk '{ print $1 }')
oc delete pod ${ACME_OPSRC_POD} -n openshift-marketplace
