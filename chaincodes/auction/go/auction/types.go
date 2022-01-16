package auction

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type Auction struct {
	Type         string             `json:"type"`
	Form         string             `json:"form"`
	Weight       uint               `json:"weight"`
	Seller       string             `json:"seller"`
	Orgs         []string           `json:"organizations"`
	PrivateBids  map[string]BidHash `json:"privateBids"`
	RevealedBids map[string]FullBid `json:"revealedBids"`
	Winner       string             `json:"winner"`
	MinPrice     uint               `json:"minimumPrice"`
	Price        uint               `json:"price"`
	Status       statusTypes        `json:"status"`
}

type FullBid struct {
	Type   string `json:"type"`
	Price  uint   `json:"price"`
	Form   string `json:"form"`
	Weight string `json:"weight"`
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
