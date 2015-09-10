# proctor
re-writing some (all?) of the bosh classroom scripts as a Go binary

### to do
- cloudformation template for a "classroom"
  - template parameters:
    - number of instances
    - AMI
    - SSH Keypair
  - resources:
    - hosted zone
    - ec2 instances
    - CNAMEs in hosted zone --> ec2 instances
- cli actions:
  - `proctor create <name> <number>
      - generates a new SSH key pair
      - determines the latest bosh-lite AMI
      - creates a new CloudFormation stack with the given name and number of instances
      - uploads SSH private key to `https://s3.amazonaws.com/bosh-classroom/classrooms/<name>/key`
      - up to instructor to maintain a shortened url pointing to that
  - `proctor destroy <name>`
    - destroys the CloudFormation stack
    - deletes the SSH keypair from EC2
    - removes the key object from S3
  - `proctor stop <name>`
    - stops (but does not destroy) all the ec2 instances
  - `proctor start <name>`
    - starts any/all stopped ec2 instances
    - updates CNAMEs to point to new IP addresses
  - `proctor run <name> -c <command>`
    - ssh's to every ec2 instance in parallel and runs `<command>`
  - `proctor run <name> -f <path/to/script.sh>`
    - ssh's to every ec2 instance in parallel and runs script

