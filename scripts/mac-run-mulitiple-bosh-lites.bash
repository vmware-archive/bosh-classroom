#!/bin/bash

INSTALL=""
OS=$(uname -s)
BOSH_LITES_DIR=$HOME/tmp/bosh-lites
brew cask --help 2>&1 > /dev/null || INSTALL+="brew install caskroom/cask/brew-cask\n"
parallel --version 2>&1 > /dev/null || INSTALL+="brew install parallel\n"
aws --version 2>&1 > /dev/null || INSTALL+="brew install awscli\n "
vagrant --version 2>&1 > /dev/null || INSTALL+="brew cask install vagrant\n"

if (vagrant plugin list > /dev/null 2>&1) ; then
  (vagrant plugin list |grep -q vagrant-aws) || INSTALL+="vagrant plugin install vagrant-aws\n"
else
  INSTALL+="vagrant plugin install vagrant-aws\n"
fi

if [[ ! -z "$INSTALL" && $OS == "Darwin" ]] ; then
  echo -n "INFO: need to run [ $INSTALL ]"
  echo -e $INSTALL | while read i ;do
    bash -c "$i"
  done
fi

if [[ -d $HOME/workspace/deployments-aws ]] ; then
  cd $HOME/workspace/deployments-aws
  git pull
else
  git clone git@github.com:pivotal-cf/deployments-aws.git $HOME/workspace/deployments-aws
fi

source $HOME/workspace/deployments-aws/thansmann/bosh_environment
chmod 600 $HOME/workspace/deployments-aws/thansmann/config/id_rsa_bosh
ssh-add $HOME/workspace/deployments-aws/thansmann/config/id_rsa_bosh
mkdir -p $HOME/.ssh
ln -svf $HOME/workspace/deployments-aws/thansmann/config/id_rsa_bosh $HOME/.ssh/id_rsa_bosh

if (aws ec2 describe-instances | egrep -q Reservations) ; then
   echo "INFO: aws cli is working"
else
   echo "ERROR: aws cli is not working; please run 'aws configure' and try again"
   exit 3
fi

seq 1 30 |
    parallel --tag --line-buffer -I % -j 5 "$HOME/workspace/deployments-aws/thansmann/scripts/run-a-vagrant-bosh-lite.bash $BOSH_LITES_DIR/%; sleep 1"

fixup-bosh-lite-vms.bash
