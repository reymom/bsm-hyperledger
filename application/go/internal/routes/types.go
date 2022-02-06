package routes

import "github.com/reymom/bsm-hyperledger/application/go/internal/connection"

type ViewData struct {
	Context string
	Name    string
}

type ChannelAuctions struct {
	Channel connection.Channel
	Auction Auction
}

type Auction struct {
	ID             string             `json:"auctionID"`
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
	opened statusTypes = iota
	closed
	finished
)

func (s statusTypes) ToString() string {
	switch s {
	case opened:
		return "Opened"
	case closed:
		return "Closed"
	case finished:
		return "Finished"
	default:
		return "Unknown"
	}
}
