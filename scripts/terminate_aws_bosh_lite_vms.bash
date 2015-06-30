#!/bin/bash

aws ec2 describe-instances --output text | tr '   ' ' ' | egrep  'INSTANCES\b|STATE\b' | paste - - | grep running |  awk '{print $7}' |
 parallel -j30 -rt "aws ec2 terminate-instances --instance-ids {}"

