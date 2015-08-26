#!/bin/bash

aws ec2 describe-instances --output text | tr '   ' ' ' | egrep  'INSTANCES\b|STATE\b' | paste - - | grep running | egrep -o '\si-........\s' | tr -d ' \t' |
 parallel -j30 -rt "aws ec2 terminate-instances --instance-ids {}"

