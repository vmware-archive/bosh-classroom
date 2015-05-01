#!/bin/bash

aws ec2 describe-instances --output text | tr '   ' ' ' | egrep  'INSTANCES\b|STATE\b' | paste - - | grep running|  awk '{print $8}' |
  parallel -j30 -rt "aws ec2 terminate-instances --instance-ids {}"

