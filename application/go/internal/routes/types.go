package routes

import "github.com/reymom/bsm-hyperledger/application/go/internal/connection"

type ViewData struct {
	Context string
	Name    string
}

type Auction struct {
	ID       string
	Seller   string
	Type     string
	Form     string
	MinPrice uint
}

type ChannelAuctions struct {
	Channel connection.Channel
	Auction Auction
}
