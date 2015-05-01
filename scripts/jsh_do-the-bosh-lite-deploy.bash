#!/bin/bash

DUMMY_RELEASE=~/workspace/dummy

pushd $DUMMY_RELEASE
if ( ! (bosh releases |egrep -q dummy) ); then
    bosh -n create release --force
    bosh -n upload release
fi

if ( ! (bosh stemcells | egrep -q bosh-warden-boshlite) ) ; then
    bosh upload stemcell  https://bosh-jenkins-artifacts.s3.amazonaws.com/bosh-stemcell/warden/bosh-stemcell-389-warden-boshlite-ubuntu-trusty-go_agent.tgz
fi

bosh deployment $DUMMY_RELEASE/meetup/deploy-i-*-manifest.yml
bosh -n deploy
bosh deployment $DUMMY_RELEASE/meetup/all-dummy-deploy-*-manifest.yml
bosh -n deploy
