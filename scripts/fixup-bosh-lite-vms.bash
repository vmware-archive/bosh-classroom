#!/bin/bash

# hard code to make sure we don't delete anything we might really want
BOSH_LITE_AMI=ami-905e65f8
export JSH_LOG_ALL=1

function abspath (){
    local DIR=$(dirname "$1");
    cd $DIR;
    printf "%s/%s\n" "$(pwd)" "$(basename "$1")" | perl -pe 's{/{2,}}{/}g';
    cd "$OLDPWD"
}


THIS_DIR=$(dirname $(abspath $0))

. $THIS_DIR/shell_helpers.bash

OUR_BOSH_LITES=$($THIS_DIR/aws-running-vms.bash| grep $BOSH_LITE_AMI | pcut -f -2 | nl2comma)
echo "our bosh-lite VMs are $OUR_BOSH_LITES"

# add host keys
our_boshlites 1 ; parallel -j 50 "ssh -o StrictHostKeyChecking=no ubuntu@{} id" ::: $OUR_BOSHLITES


chmod 755 $THIS_DIR/jsh_bosh-lite-box-prep.bash
jsh -e -w $OUR_BOSH_LITES -l ubuntu -s $THIS_DIR/jsh_bosh-lite-box-prep.bash
