package auction

import (
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type Auction struct {
	ID             string         `json:"auctionID"`
	ClientID       string         `json:"clientID"`
	StartDate      time.Time      `json:"startDate"`
	EndDate        time.Time      `json:"endDate"`
	IsPrivate      bool           `json:"isPrivate"`
	CollectionName string         `json:"collectionName"`
	Type           string         `json:"type"`
	Form           string         `json:"form"`
	Weight         uint           `json:"weight"`
	Seller         string         `json:"seller"`
	Bidders        []string       `json:"bidders"`
	Bids           map[string]Bid `json:"bids"`
	Winner         string         `json:"winner"`
	MinPrice       uint           `json:"minimumPrice"`
	Price          uint           `json:"price"`
	Status         statusTypes    `json:"status"`
}

type Bid struct {
	ID    string `json:"bidID"`
	Buyer string `json:"buyer"`
	Price uint   `json:"price"`
}

type statusTypes uint

const (
	opened statusTypes = iota
	finished
)
