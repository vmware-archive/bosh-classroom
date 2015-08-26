source $HOME/workspace/deployments-aws/thansmann/bosh_environment
PATH+=":$HOME/workspace/bosh-classroom/scripts"
ssh-add $HOME/workspace/deployments-aws/thansmann/config/id_rsa_bosh


function grep_ip(){
grep -E "(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)" $*

}

function our_boshlites() {
    if [[ ! -z "$*" || -z "$OUR_BOSHLITES" ]]; then
        export OUR_BOSHLITES=$(aws-running-vms.bash | grep_ip -o)
      fi
    export OUR_BOSHLITES_JSH=$(echo $OUR_BOSHLITES | tr ' ' ',')
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

function vagrant_ssh_setup () {
  mkdir -p ~/.ssh
  curl -s https://raw.githubusercontent.com/mitchellh/vagrant/master/keys/vagrant > ~/.ssh/id_rsa_vagrant
  chmod 600 ~/.ssh/id_rsa_vagrant
  cat << EOF >> ~/.ssh/config

Host vagrant
	HostName 192.168.50.4
	User vagrant
	IdentityFile ~/.ssh/id_rsa_vagrant
EOF
}

function vagrant_ssh_setup () {
  mkdir -p ~/.ssh
  curl -s https://raw.githubusercontent.com/mitchellh/vagrant/master/keys/vagrant > ~/.ssh/id_rsa_vagrant
  chmod 600 ~/.ssh/id_rsa_vagrant

  if (egrep -q id_rsa_vagrant $HOME/.ssh/config) ; then
    echo "INFO: vagrant ssh config found in $HOME/.ssh/config"
  else
cat << EOF >> $HOME/.ssh/config
Host vagrant-default
        HostName 192.168.50.4
        User vagrant
        IdentityFile $ALT_HOME/.ssh/id_rsa_vagrant
EOF
  fi
}
