source $HOME/workspace/deployments-aws/thansmann/bosh_environment
PATH+=":$HOME/workspace/bosh-classroom/scripts"
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

function bosh-lite-ssh-id-add (){
  our_boshlites 1
  parallel -j 50 "ssh -o StrictHostKeyChecking=no ubuntu@{} id" ::: $OUR_BOSHLITES
}

function aws-stop-bosh-lites (){
  RUNNING_BOSHLITES=$(aws-running-vms.bash | grep ami-905e65f8 |grep running | awk '{print $2}'| nl2space)
  parallel -j 2 "aws ec2 stop-instances --instance-ids {}" ::: $RUNNING_BOSHLITES
}

function aws-start-bosh-lites (){
  STOPPED_BOSHLITES=$(aws-stopped-vms.bash| awk '{print $2}')
  parallel -j 2 "aws ec2 start-instances --instance-ids {}" ::: $STOPPED_BOSHLITES
}

function bosh-lite-ips (){
  aws-running-vms.bash | grep ami-905e65f8 |grep running | grep_ip -o
}
