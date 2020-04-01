#!/bin/sh
. ./settings.sh

oc apply -n openshift-marketplace -f ./deploy/operator-source.yaml