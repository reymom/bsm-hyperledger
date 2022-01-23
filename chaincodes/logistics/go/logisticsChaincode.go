package main

import (
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/reymom/bsm-hyperledger/chaincodes/logistics/go/logistics"
)

func main() {
	logisticsSmartContract, err := contractapi.NewChaincode(&logistics.SmartContract{})
	if err != nil {
		log.Panicf("Error creating logistics chaincode: %v", err)
	}

	if err := logisticsSmartContract.Start(); err != nil {
		log.Panicf("Error starting logistics chaincode: %v", err)
	}
}
