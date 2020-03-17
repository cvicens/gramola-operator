#!/bin/sh

export GO111MODULE=on

export OPERATOR_NAME="gramola-operator"
export OPERATOR_IMAGE="gramola-operator-image"
export API_VERSION="gramola.redhat.com/v1alpha1"

export PROJECT_NAME=${OPERATOR_NAME}-project

export USERNAME=<USERNAME>

export OPERATOR_VERSION=0.0.1