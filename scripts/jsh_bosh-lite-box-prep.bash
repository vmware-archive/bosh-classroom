#!/bin/bash
sudo chown -R ubuntu .

if [[ ! -f ~/did-apt-get-update ]]; then
    sudo apt-get update
    date >> ~/did-apt-get-update
fi

DUMMY_RELEASE=~/workspace/dummy
export OUR_AWS_ID=$(curl -s http://169.254.169.254/latest/meta-data/instance-id)
export  UUID=$(bosh -n status |grep UUID| awk '{print $NF}')

[[ -f ~/.ssh/id_rsa ]] || {
    mkdir -p ~/.ssh
    chmod 700 ~/.ssh
    ssh-keygen -t rsa -P foobar -f ~/.ssh/id_rsa
}

type -a git || sudo apt-get install -y git

mkdir -p ~/workspace ~/tmp

if [[ ! -d ~/basic-env ]] ; then
    cd && git clone https://github.com/pivotal-cf-experimental/basic-env.git
    . basic-env/.profile
    new_env
    rm -rf ~/parrallel-*
    ln -svf .profile .bashrc
fi

if [[ ! -d $DUMMY_RELEASE ]]; then
    git clone https://github.com/pivotal-cf-experimental/dummy-boshrelease.git $DUMMY_RELEASE
    cd  $DUMMY_RELEASE
    mkdir -p $DUMMY_RELEASE/classroom
    echo -e "---\nname: first-one\ndirector_uuid: $UUID" > $DUMMY_RELEASE/classroom/first-one-stub.yml
    cat  $DUMMY_RELEASE/classroom/this-bosh-lite.yml
    echo -e "---\nname: second-one\ndirector_uuid: $UUID" > $DUMMY_RELEASE/classroom/second-one-stub.yml
fi

type -a spiff 2>&1 > /dev/null || {
    cd ~/tmp/ && wget -q https://github.com/cloudfoundry-incubator/spiff/releases/download/v1.0.6/spiff_linux_amd64.zip
    type -a unzip || sudo apt-get install -y zip unzip
    unzip spiff_linux_amd64.zip
    sudo cp -av spiff /usr/local/bin/
    type -a spiff
}

if [[  ! -f  $DUMMY_RELEASE/classroom/deploy-${OUR_AWS_ID}-manifest.yml ]]; then
    echo "INFO: Generating dummy deploy manifest"
    cd  $DUMMY_RELEASE && {
        bash -x ./generate_deployment_manifest warden  $DUMMY_RELEASE/classroom/first-one-stub.yml >  $DUMMY_RELEASE/classroom/deploy-${OUR_AWS_ID}-manifest.yml
        spiff merge $DUMMY_RELEASE/templates/all-jobs-dummy-deployment.yml \
              $DUMMY_RELEASE/classroom/this-all-dummy-bosh-lite.yml > $DUMMY_RELEASE/classroom/all-dummy-deploy-${OUR_AWS_ID}-manifest.yml
        ls -al  $DUMMY_RELEASE/classroom/
    }
    cd -
else
    echo "INFO: dummy deploy manifest already generated"
fi

if [[ ! -f ~/bosh-stemcell-389-warden-boshlite-ubuntu-trusty-go_agent.tgz ]]; then
    cd && wget -q https://bosh-jenkins-artifacts.s3.amazonaws.com/bosh-stemcell/warden/bosh-stemcell-389-warden-boshlite-ubuntu-trusty-go_agent.tgz
    ls -lart bosh-stemcell-389-warden-boshlite-ubuntu-trusty-go_agent.tgz
else
    echo "INFO: Stemcell is already downloaded"
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

echo "PATH+=:/var/vcap/bosh/bin" >> ~/.profile

sudo perl -i.old -pe 's{(minimum_down_jobs:).*$}{$1 1\n}xms' /var/vcap/jobs/health_monitor/config/health_monitor.yml
