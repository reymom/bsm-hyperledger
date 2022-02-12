package connection

import (
	"strings"

	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
)

//gateway contracts
type GatewayContract struct {
	Name       Contract
	GwContract *gateway.Contract
}

type NetworkContract map[Channel]GatewayContract

//channels
type Channel string

const (
	Public1Channel Channel = "public1channel"
	Public2Channel Channel = "public2channel"

	LogisticsChannel Channel = "logisticschannel"
)

func (c *Channel) IsAuctionChannel() bool {
	return strings.Contains(string(*c), "public")
}

//contracts
type Contract string

const (
	AuctionContract   Contract = "auction"
	LogisticsContract Contract = "logistics"
)

func (c Channel) GetContract() Contract {
	switch c {
	case Public1Channel:
		return AuctionContract
	case Public2Channel:
		return AuctionContract
	case LogisticsChannel:
		return LogisticsContract
	default:
		return Contract("")
	}
}
