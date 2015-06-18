#!/bin/bash
set -e

INSTALL=""
OS=$(uname -s)
BOSH_LITES_DIR=$HOME/tmp/bosh-lites

abspath() {
    local DIR=$(dirname "$1")
    cd $DIR
    printf "%s/%s\n" "$(pwd)" "$(basename "$1")" | perl -pe 's{/{2,}}{/}g'
    cd "$OLDPWD"
}

$THIS_DIR/mac-check-prerequsites.bash

THIS_DIR=$(abspath $0)
seq 1 25 |
    parallel --tag --line-buffer -I % -j 5 "$HOME/workspace/deployments-aws/thansmann/scripts/run-a-vagrant-bosh-lite.bash $BOSH_LITES_DIR/%; sleep 1"

$THIS_DIR/fixup-bosh-lite-vms.bash
