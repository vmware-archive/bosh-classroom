source $HOME/workspace/deployments-aws/thansmann/bosh_environment
PATH+=":$HOME/workspace/deployments-aws/thansmann/scripts"
ssh-add $HOME/workspace/deployments-aws/thansmann/config/id_rsa_bosh

function our_boshlites() {
    if [[ ! -z "$*" || -z "$OUR_BOSHLITES" ]]; then
        export OUR_BOSHLITES=$(aws-running-vms.bash | pcut -f -2)
    fi
    OUR_BOSHLITES_JSH=$(echo $OUR_BOSHLITES|space2comma)
    echo $OUR_BOSHLITES

}

function boshlites() {
    if [[ ! -z "$*" ]] ; then
        parallel -j 100 -rt --keep --tag "bosh -t {} $*" ::: $(our_boshlites)
    else
        parallel -j 100 -rt --keep --tag "bosh -t {} deployments" ::: $(our_boshlites)
    fi
}

