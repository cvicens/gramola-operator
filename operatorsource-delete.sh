#!/bin/sh
. ./settings.sh

oc delete -n openshift-marketplace -f ./deploy/operator-source.yaml