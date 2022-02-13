package logistics

import (
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type SteelDelivery struct {
	AuctionID   string      `json:"auctionID"`
	DestinyOrg  string      `json:"destinyOrg"`
	Creator     string      `json:"creator"`
	DeliveryOrg string      `json:"deliveryOrg"`
	Address     *Address    `json:"address"`
	Updated     time.Time   `json:"timestamp"`
	Status      statusTypes `json:"status"`
}

type HistoryPublicHashTracker struct {
	HashID string      `json:"hashID"`
	Status statusTypes `json:"status"`
}

type HistoryQueryResult struct {
	Record    *HistoryPublicHashTracker `json:"record"`
	TxId      string                    `json:"txId"`
	Timestamp time.Time                 `json:"timestamp"`
	Deleted   bool                      `json:"delete"`
}

type Address struct {
	Country string `json:"country"`
	City    string `json:"city"`
	Street  string `json:"street"`
	Number  string `json:"number"`
}

type statusTypes uint

const (
	processing statusTypes = iota
	onDelivery
	delivered
	failed
)
