package aws_test

import (
	"fmt"
	"math/rand"
	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cloudformation", func() {
	const minimalTemplate = `
{
  "AWSTemplateFormatVersion" : "2010-09-09",
  "Description" : "Minimal template for unit-testing proctor",

  "Parameters" : {
    "SomePortNumber" : {
      "Description" : "A port number",
      "Type": "Number"
		},
    "SomeCIDRBlock" : {
      "Description" : "A CIDR block",
      "Type": "String"
   }
  },

  "Resources" : {
    "InstanceSecurityGroup" : {
      "Type" : "AWS::EC2::SecurityGroup",
      "Properties" : {
        "GroupDescription" : "A security group for unit testing proctor",
        "SecurityGroupIngress" : [ {
          "IpProtocol" : "tcp",
					"FromPort" : "22",
          "ToPort" : { "Ref" : "SomePortNumber"},
          "CidrIp" : { "Ref" : "SomeCIDRBlock"}
        } ]
      }
    }
  }
}
`

	It("should create and delete stacks", func() {
		stackName := fmt.Sprintf("test-stack-%d", rand.Int63())
		template := minimalTemplate

		randomPort := rand.Intn(65500) + 22
		randomOctet := rand.Intn(250)
		parameters := map[string]string{
			"SomePortNumber": strconv.Itoa(randomPort),
			"SomeCIDRBlock":  fmt.Sprintf("10.%d.0.0/24", randomOctet),
		}

		stackOperationTimeout := "60s"

		stackID, err := awsClient.CreateStack(stackName, template, parameters)
		Expect(err).NotTo(HaveOccurred())

		stackStatus, returnedStackID, returnedParams, err := awsClient.DescribeStack(stackName)
		Expect(err).NotTo(HaveOccurred())
		Expect(stackStatus).To(Equal("CREATE_IN_PROGRESS"))
		Expect(returnedParams).To(Equal(parameters))
		Expect(returnedStackID).To(Equal(stackID))

		Eventually(func() (string, error) {
			status, _, _, err := awsClient.DescribeStack(stackID)
			return status, err
		}, stackOperationTimeout).
			Should(Equal("CREATE_COMPLETE"))

		hosts, err := awsClient.GetHostsFromStackID(stackID)
		Expect(err).NotTo(HaveOccurred())
		Expect(hosts).To(HaveLen(0))

		Expect(awsClient.DeleteStack(stackName)).To(Succeed())

		Eventually(func() (string, error) {
			status, _, _, err := awsClient.DescribeStack(stackID)
			return status, err
		}, stackOperationTimeout).
			Should(Equal("DELETE_COMPLETE"))
	})
})
