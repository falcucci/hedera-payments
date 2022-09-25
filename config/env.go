package config

import (
	"os"
	"strconv"
)

var (
	// OperatorKEY OperatorKEY
	OperatorKEY string

	// NodeAddress NodeAddress
	NodeAddress string

	// NodeID NodeID
	NodeID int64

	// OperatorID OperatorID
	OperatorID int64
)

// LoadEnvVars load all envs
func LoadEnvVars() {
	// var err error

	OperatorKEY = os.Getenv("OPERATOR_KEY")

	NodeAddress = os.Getenv("NODE_ADDRESS")

	nodeID, err := strconv.Atoi(os.Getenv("NODE_ID"))
	if err != nil {
		NodeID = 3
	}

	NodeID = int64(nodeID)

	operatorID, err := strconv.Atoi(os.Getenv("OPERATOR_ID"))
	if err != nil {
		NodeID = 3
	}

	OperatorID = int64(operatorID)

}
