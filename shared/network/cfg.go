package network

const netPath = "./net-data"

// nodePath This needs to match the name of
// the node as found in the network.json file
// ToDo: configure the path on network start-up.
const nodePath = netPath + "/primary"

func Path() string {
	return netPath
}

func NodePath() string {
	return nodePath
}
