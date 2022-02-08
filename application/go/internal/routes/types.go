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

func (s statusTypes) ToString() string {
	switch s {
	case opened:
		return "Opened"
	case finished:
		return "Finished"
	default:
		return ""
	}
}

const (
	opened statusTypes = iota
	finished
)

//deliveries
type Delivery struct {
	AuctionID   string      `json:"auctionID"`
	DestinyOrg  string      `json:"destinyOrg"`
	Creator     string      `json:"creator"`
	DeliveryOrg string      `json:"deliveryOrg"`
	Address     *Address    `json:"address"`
	Updated     time.Time   `json:"timestamp"`
	Status      statusTypes `json:"status"`
}

type HistoryQueryResult struct {
	Record    *Delivery `json:"record"`
	TxId      string    `json:"txId"`
	Timestamp time.Time `json:"timestamp"`
	Deleted   bool      `json:"delete"`
}

type Address struct {
	Country string `json:"country"`
	City    string `json:"city"`
	Street  string `json:"street"`
	Number  string `json:"number"`
}

type deliveryStatusType uint

const (
	processing deliveryStatusType = iota
	onDelivery
	delivered
	failed
)

func (s deliveryStatusType) ToString() string {
	switch s {
	case processing:
		return "Processing"
	case onDelivery:
		return "On Delivery"
	case delivered:
		return "Delivered"
	case failed:
		return "Failed"
	default:
		return ""
	}
}
