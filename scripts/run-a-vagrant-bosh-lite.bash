#!/bin/bash

BOSH_LITES_DIR=$HOME/tmp/bosh-lites
aws --version 2>&1 > /dev/null || INSTALL+=awscli
vagrant --version 2>&1 > /dev/null || INSTALL+=vagrant

if [[ ! -z "$INSTALL" && $OS == "Darwin" ]] ; then
  for i in $INSTALL ; do
    brew install $i
  done
fi

mkdir -p $1
 cd $1
 rm -rfv Vagrant .vagrant
 cd $1 && {
    vagrant init cloudfoundry/bosh-lite
    vagrant up --provider=aws
 }
