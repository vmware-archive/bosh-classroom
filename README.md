# bosh-classroom
Tools to start and setup BOSH-lite VMs with vagrant aws-plugin for classroom training

# How to setup the bosh-classroom with the current scripts (Mac OSX)

### Before the classroom

- Install vagrant 1.6.3
- `cp bin/jsh /usr/local/bin/jsh`
- `cd scripts`
- `shell_helpers.bash`
- `run-multiple-bosh-lites.bash <number of classroom vms>` #default 20
 

### Running commands during the classroom:

# Get list of IPs
- `our_boshlites`

Example: 

- `run_on_classroom 'ls'`
- `run_on_classroom 'bosh status'`

### Ending the classroom

- `terminate_aws_bosh_lite_vms.bash`
