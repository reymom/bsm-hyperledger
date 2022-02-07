package routes

import (
	"time"

	"github.com/reymom/bsm-hyperledger/application/go/internal/connection"
)

type ViewData struct {
	Context string
	Name    string
}

type ChannelAuctions struct {
	Channel connection.Channel
	Auction Auction
}

type Auction struct {
	ID             string         `json:"auctionID"`
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

func (s statusTypes) ToString() string {
	switch s {
	case opened:
		return "Opened"
	case finished:
		return "Finished"
	default:
		return "Unknown"
	}
}
