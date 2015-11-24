package apis

type InstanceState string

const (
	CREATED   InstanceState = "created"
	STARTED   InstanceState = "started"
	STOPPED   InstanceState = "stopped"
	DESTROYED InstanceState = "stopped"
)

type Instance struct {
	Name       string
	State      InstanceState
	Connection Connection
	Tags       map[string]string
}

type ConnectionType string

const (
	SSH ConnectionType = "SSH"
)

type Connection struct {
	Type   ConnectionType
	Config map[string]interface{}
}
