#!/bin/bash

DUMMY_RELEASE=~/workspace/dummy

pushd $DUMMY_RELEASE
if ( ! (bosh releases |egrep -q dummy) ); then
    bosh -n create release --force
    bosh -n upload release
fi

if ( ! (bosh stemcells | egrep -q bosh-warden-boshlite) ) ; then
    if [[ -f $HOME/bosh-stemcell-389-warden-boshlite-ubuntu-trusty-go_agent.tgz ]] ; then
       bosh upload stemcell $HOME/bosh-stemcell-389-warden-boshlite-ubuntu-trusty-go_agent.tgz
    else
    bosh upload stemcell  https://bosh-jenkins-artifacts.s3.amazonaws.com/bosh-stemcell/warden/bosh-stemcell-389-warden-boshlite-ubuntu-trusty-go_agent.tgz
    fi
fi

for i in $DUMMY_RELEASE/classroom/{first,second}.yml ; do
    bosh -d $i -n deploy
done
