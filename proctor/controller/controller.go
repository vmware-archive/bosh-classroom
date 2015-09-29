package controller

type atlasClient interface {
	GetLatestAMIs(string) (map[string]string, error)
}

type awsClient interface {
	CreateKey(name string) (string, error)
	DeleteKey(name string) error
	ListKeys(prefix string) ([]string, error)
	StoreObject(name string, bytes []byte, downloadFileName, contentType string) error
	DeleteObject(name string) error
	URLForObject(name string) string
	CreateStack(name string, template string, parameters map[string]string) (string, error)
	DeleteStack(name string) error
	DescribeStack(name string) (string, string, map[string]string, error)
	GetHostsFromStackID(stackID string) (map[string]string, error)
}

type cliLogger interface {
	Println(indentation int, format string, args ...interface{})
	Green(format string, args ...interface{}) string
}

type Controller struct {
	AtlasClient atlasClient
	AWSClient   awsClient
	Log         cliLogger

	VagrantBoxName string
	Region         string
	Template       string
}

func prefix(classroomName string) string {
	return "classroom-" + classroomName
}
