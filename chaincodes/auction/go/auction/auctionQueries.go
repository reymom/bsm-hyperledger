package auction

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// QueryAuction
func (s *SmartContract) QueryAuction(ctx contractapi.TransactionContextInterface, auctionID string) (*Auction, error) {

	auctionJSON, err := ctx.GetStub().GetState(auctionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get auction object %v: %v", auctionID, err)
	}
	if auctionJSON == nil {
		return nil, fmt.Errorf("auction does not exist")
	}

	var auction *Auction
	err = json.Unmarshal(auctionJSON, &auction)
	if err != nil {
		return nil, err
	}

	return auction, nil
}

func (s *SmartContract) QueryPrivateAuction(ctx contractapi.TransactionContextInterface, auctionID string, collectionOrgNums string) (*Auction, error) {

	collectionName, err := getPrivateCollectionChannel(ctx, collectionOrgNums)
	if err != nil {
		return nil, fmt.Errorf("failed to get private collection name: %v", err)
	}

	auctionJSON, err := ctx.GetStub().GetPrivateData(collectionName, auctionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get private auction object %v: %v", auctionID, err)
	}
	if auctionJSON == nil {
		return nil, fmt.Errorf("auction does not exist")
	}

	var auction *Auction
	err = json.Unmarshal(auctionJSON, &auction)
	if err != nil {
		return nil, err
	}

	return auction, nil
}

// GetAllAuctions
func (s *SmartContract) GetAllAuctions(ctx contractapi.TransactionContextInterface) ([]*Auction, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var auctions []*Auction
	for resultsIterator.HasNext() {
		auctionJSON, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var auction Auction
		err = json.Unmarshal(auctionJSON.Value, &auction)
		if err != nil {
			return nil, err
		}

		auctions = append(auctions, &auction)
	}

	return auctions, nil
}

func (s *SmartContract) GetAllPrivateAuctions(ctx contractapi.TransactionContextInterface, collectionOrgNums string) ([]*Auction, error) {
	collectionName, err := getPrivateCollectionChannel(ctx, collectionOrgNums)
	if err != nil {
		return nil, fmt.Errorf("failed to get private collection name: %v", err)
	}

	resultsIterator, err := ctx.GetStub().GetPrivateDataByRange(collectionName, "", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var auctions []*Auction
	for resultsIterator.HasNext() {
		auctionJSON, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var auction Auction
		err = json.Unmarshal(auctionJSON.Value, &auction)
		if err != nil {
			return nil, err
		}

		auctions = append(auctions, &auction)
	}

	return auctions, nil
}
