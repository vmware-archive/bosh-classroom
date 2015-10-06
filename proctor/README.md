# proctor
A tool for running BOSH 101 classrooms.

[Pivotal Tracker project](https://www.pivotaltracker.com/n/projects/1434846)

### Installation
Download a [binary release](https://github.com/pivotal-cf-experimental/bosh-classroom/releases), and `chmod +x` it. 


### Basic usage
0. Load credentials for your AWS environment
    ```
    export AWS_DEFAULT_REGION=us-east-1
    export AWS_ACCESS_KEY_ID=YOUR-ACCESS-KEY
    export AWS_SECRET_ACCESS_KEY=YOUR-SECRET-KEY
    ```
    
0. Create a new classroom
    ```
    proctor create -name my-classroom -number 3
    ```
    This will spin up 3 EC2 instances in your AWS account.
    
0. Watch your classroom get created
    ```
    proctor describe -name my-classroom
    ```
    The SSH key was generated at `create` time and is world-readable.

0. Run a command on all VMs
    ```
    proctor run -name my-classroom -c 'bosh status'
    ```

0. Destroy your classroom
    ```
    proctor destroy -name my-classroom
    ```
    

### Building from source
You need to have a Go development environment set up.  Make sure that `$GOPATH/bin` is in your `$PATH`
```
go get github.com/pivotal-cf-experimental/bosh-classroom/proctor/...
go install github.com/pivotal-cf-experimental/bosh-classroom
```    
