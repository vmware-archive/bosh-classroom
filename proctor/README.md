# proctor
re-writing some (all?) of the bosh classroom scripts as a Go binary

### to do
- cloudformation template for a "classroom"
  - ~~parameters: AMI, keypair, num_instances~~
  - resources:
    - ~~ec2 instances~~
    - hosted zone
    - CNAMEs in hosted zone --> ec2 instances
- cli actions on a classroom instance named `<name>`:
  - `proctor create -name <name> -number <number>`
      - ~~creates a new CloudFormation stack with the given name and number of instances~~
      - ~~generates a new SSH key pair~~
      - ~~determines the latest bosh-lite AMI~~
      - ~~uploads SSH private key to `https://bosh101.s3.amazonaws.com/keys/<name>`~~
  - `proctor destroy -name <name>`
    - ~~destroys the CloudFormation stack~~
    - ~~deletes the SSH keypair from EC2~~
    - ~~removes the key object from S3~~
  - `proctor stop -name <name>`
    - stops (but does not destroy) all the ec2 instances
  - `proctor start -name <name>`
    - starts any/all stopped ec2 instances
    - updates CNAMEs to point to new IP addresses
  - `proctor run -name <name> -c <command>`
    - ssh's to every ec2 instance in parallel and runs `<command>`
  - `proctor run -name <name> -f <path/to/script.sh>`
    - ssh's to every ec2 instance in parallel and runs script

