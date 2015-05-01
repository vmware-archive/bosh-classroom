#!/bin/bash
if [[ -z "$AWS_ACCESS_KEY_ID" ]] ; then
   . $HOME/workspace/deployments-aws/thansmann/bosh_environment
fi

aws ec2 describe-instances --output text | tr '   ' ' ' | egrep  'INSTANCES\b|STATE\b' | paste - - | grep running| awk '{print $7,$8,$9,$11,$12,$15,$16,$NF}'
