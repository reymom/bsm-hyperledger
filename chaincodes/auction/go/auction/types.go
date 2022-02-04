package auction

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type Auction struct {
	ID             string             `json:"auctionID"`
	ClientID       string             `json:"clientID"`
	IsPrivate      bool               `json:"isPrivate"`
	CollectionName string             `json:"collectionName"`
	Type           string             `json:"type"`
	Form           string             `json:"form"`
	Weight         uint               `json:"weight"`
	Seller         string             `json:"seller"`
	Orgs           []string           `json:"organizations"`
	PrivateBids    map[string]BidHash `json:"privateBids"`
	RevealedBids   map[string]FullBid `json:"revealedBids"`
	Winner         string             `json:"winner"`
	MinPrice       uint               `json:"minimumPrice"`
	Price          uint               `json:"price"`
	Status         statusTypes        `json:"status"`
}

type FullBid struct {
	Price  uint   `json:"price"`
	Org    string `json:"org"`
	Bidder string `json:"bidder"`
}

type BidHash struct {
	Org  string `json:"org"`
	Hash string `json:"hash"`
}

type statusTypes uint

const (
	created statusTypes = iota
	opened
	closed
	finished
)
