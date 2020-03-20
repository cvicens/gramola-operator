#!/bin/sh
# Doc: https://github.com/operator-framework/operator-sdk/blob/master/doc/test-framework/scorecard.md
# See: .osdk-scorecard.yaml

. ./settings.sh

operator-sdk scorecard
