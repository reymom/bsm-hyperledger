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

	Logistics11Channel Channel = "logistics11channel"
	Logistics12Channel Channel = "logistics12channel"
	Logistics13Channel Channel = "logistics13channel"
	Logistics21Channel Channel = "logistics21channel"
	Logistics22Channel Channel = "logistics22channel"
	Logistics23Channel Channel = "logistics23channel"
)

func (c *Channel) IsAuctionChannel() bool {
	return strings.Contains(string(*c), "public")
}

func (c *Channel) GetCollections() []string {
	switch *c {
	case Public1Channel:
		//if we add more collections, we can add them in the lists here
		return []string{"12"}
	case Public2Channel:
		return []string{"23"}
	default:
		return []string{""}
	}
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
	case Logistics11Channel:
		return LogisticsContract
	case Logistics12Channel:
		return LogisticsContract
	case Logistics13Channel:
		return LogisticsContract
	case Logistics21Channel:
		return LogisticsContract
	case Logistics22Channel:
		return LogisticsContract
	case Logistics23Channel:
		return LogisticsContract
	default:
		return ""
	}
}
