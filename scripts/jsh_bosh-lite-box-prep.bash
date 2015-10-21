#!/bin/bash
_USERID=$(id --name -u)
sudo chown -R $_USERID .
mkdir -p ~/workspace ~/tmp

if [[ ! -f ~/bosh-stemcell-389-warden-boshlite-ubuntu-trusty-go_agent.tgz ]]; then
    cd && nohup wget -q https://bosh-jenkins-artifacts.s3.amazonaws.com/bosh-stemcell/warden/bosh-stemcell-389-warden-boshlite-ubuntu-trusty-go_agent.tgz &
disown $!
else
    echo "INFO: Stemcell is already downloaded"
fi

DUMMY_RELEASE=~/workspace/dummy

export  UUID=$(bosh -n status |grep UUID| awk '{print $NF}')

[[ -f ~/.ssh/id_rsa ]] || {
    mkdir -p ~/.ssh
    chmod 700 ~/.ssh
    ssh-keygen -t rsa -P foobar -f ~/.ssh/id_rsa
}
if [[ ! -f ~/did-apt-get-update ]]; then
    sudo apt-get update
    date >> ~/did-apt-get-update
fi
function install_spiff()
{
    cd ~/tmp/ && wget -q https://github.com/cloudfoundry-incubator/spiff/releases/download/v1.0.6/spiff_linux_amd64.zip
    unzip spiff_linux_amd64.zip
    sudo cp -av spiff /usr/local/bin/
    type -a spiff
}

function install_packages()
{
    sudo apt-get install -y zip unzip git  || echo "Failed to install apt packages"
    install_spiff || echo "Failed to install spiff"
}

install_packages

if [[ ! -d ~/basic-env ]] ; then
    cd && git clone https://github.com/pivotal-cf-experimental/basic-env.git
    . basic-env/.profile
    new_env
    rm -rf ~/parrallel-*
fi

if [[ ! -d $DUMMY_RELEASE ]]; then
    set -x
    git clone https://github.com/pivotal-cf-experimental/dummy-boshrelease.git $DUMMY_RELEASE
    cd  $DUMMY_RELEASE
    mkdir -p $HOME/workspace/classroom
    echo -e "---\nname: first\ndirector_uuid: $UUID" >  $HOME/workspace/classroom/stub-first.yml
    echo -e "---\nname: second\ndirector_uuid: $UUID" > $HOME/workspace/classroom/stub-second.yml
    set +x
fi


if [[  ! -f  $HOME/workspace/classroom/first.yml ]]; then
    echo "INFO: Generating dummy deploy manifest"
    cd  $DUMMY_RELEASE && \
    ./generate_deployment_manifest warden ~/workspace/classroom/stub-first.yml > ~/workspace/classroom/first.yml
    else
      echo "INFO: dummy deploy manifest 'first.yml' already generated"
fi

if [[  ! -f  $HOME/workspace/classroom/second.yml ]]; then
    cd  $DUMMY_RELEASE && \
    spiff merge $DUMMY_RELEASE/templates/all-jobs-dummy-deployment.yml \
              $HOME/workspace/classroom/stub-second.yml > ~/workspace/classroom/second.yml
    ls -al  $HOME/workspace/classroom/
    cd -
    else
      echo "INFO: dummy deploy manifest 'second.yml' already generated"
fi

if ( ! (egrep -q sFnRXKn6gwnutEwDSvxwyl19pk4EKtQz ~/.ssh/authorized_keys) ); then
    echo "adding pub key file"
cat<<EOF >> ~/.ssh/authorized_keys
ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCsc8KIhDSPz5bl+KbGvqU9r3CbIJDICxfMrhW9fqNMu27HIL8oVvQOdPo8D+oTAOQUvObzDMw98dm0W7cQYaEh3L41IWQeV/ueaJnwvMsDvYZb3JPYWIYB72fEzf/Bg/J3RxloTU5W9sn9G2otMPEXjVC17Fsy8q4dlSnu3iNs9koiNWR+5mDenRSHFB/FllI+AB2HhlGzH0HewWuBZCNBbV4vT4sy9kkiisYmziOMwnUB32xGtOWH6wJFex/RZxX9FsqbP6GqxCfLzBNr81ZCHjGnG8NomNeqCEKX+qPT353ZBCHWFBamsFnRXKn6gwnutEwDSvxwyl19pk4EKtQz thansmann@bullfinch
EOF
else
    echo "INFO: pub key file has already been added"
fi
sudo updatedb

sudo perl -i.old -pe 's{(minimum_down_jobs:).*$}{$1 1\n}xms' /var/vcap/jobs/health_monitor/config/health_monitor.yml
