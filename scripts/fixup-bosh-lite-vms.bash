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

# add host keys
our_boshlites 1 
echo "our bosh-lite VMs are $OUR_BOSHLITES"
parallel -j 50 "ssh -o StrictHostKeyChecking=no ubuntu@{} id" ::: $OUR_BOSHLITES
echo parallel -j 50 "ssh -o StrictHostKeyChecking=no ubuntu@{} id" ::: $OUR_BOSHLITES

chmod 755 $THIS_DIR/jsh_*

jcp -l ubuntu -w $OUR_BOSHLITES_JSH jcp_fixed_up_sshd_config /tmp/sshd_config
jsh -e -w $OUR_BOSHLITES_JSH -l ubuntu -s $THIS_DIR/jsh_enable_ubuntu_c1oudc0w_login

jsh -e -w $OUR_BOSHLITES_JSH -l ubuntu -s $THIS_DIR/jsh_bosh-lite-box-prep.bash
