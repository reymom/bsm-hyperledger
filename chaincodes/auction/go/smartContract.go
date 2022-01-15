package main

import (
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	// "github.com/reymom/bsm-hyperledger/chaincodes/auction/go/smart-contract"
)

func main() {
	auctionSmartContract, err := contractapi.NewChaincode(&auction.SmartContract{})
	if err != nil {
		log.Panicf("Error creating auction chaincode: %v", err)
	}

	if err := auctionSmartContract.Start(); err != nil {
		log.Panicf("Error starting auction chaincode: %v", err)
	}
}
