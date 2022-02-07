package auction

import (
	"encoding/base64"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func (s *SmartContract) GetSubmittingClientIdentity(ctx contractapi.TransactionContextInterface) (string, error) {

	b64ID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "", fmt.Errorf("failed to read clientID: %v", err)
	}
	decodeID, err := base64.StdEncoding.DecodeString(b64ID)
	if err != nil {
		return "", fmt.Errorf("failed to base64 decode clientID: %v", err)
	}
	return string(decodeID), nil
}

func (s *SmartContract) setWinnerOfAuction(ctx contractapi.TransactionContextInterface, auction *Auction) (*Auction, error) {

	maxPrice := uint(0)
	bids := auction.Bids
	var winner string
	for _, bidder := range auction.Bidders {
		if bids[bidder].Price > maxPrice {
			maxPrice = bids[bidder].Price
			winner = bids[bidder].Buyer
		}
	}

	auction.Winner = winner
	auction.Price = maxPrice

	return auction, nil
}

// getPrivateCollenction is an internal helper function to get a private auction channel.
func getPrivateCollectionChannel(ctx contractapi.TransactionContextInterface, channel string) (string, error) {

	_, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return "", fmt.Errorf("failed to get verified MSPID: %v", err)
	}

	// Create the collection name
	orgCollection := "privateCollection" + channel

	return orgCollection, nil
}
