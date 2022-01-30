package connection

import "github.com/hyperledger/fabric-sdk-go/pkg/gateway"

//gateway contracts
type GatewayContract struct {
	name       Contract
	gwContract *gateway.Contract
}

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

//contracts
type Contract string

const (
	AuctionContract1 Contract = "auction1"
	AuctionContract2 Contract = "auction2"

	LogisticsContract11 Contract = "logistics11"
	LogisticsContract12 Contract = "logistics12"
	LogisticsContract13 Contract = "logistics13"
	LogisticsContract21 Contract = "logistics21"
	LogisticsContract22 Contract = "logistics22"
	LogisticsContract23 Contract = "logistics23"
)

func (c Channel) GetContract() Contract {
	switch c {
	case Public1Channel:
		return AuctionContract1
	case Public2Channel:
		return AuctionContract2
	case Logistics11Channel:
		return LogisticsContract11
	case Logistics12Channel:
		return LogisticsContract12
	case Logistics13Channel:
		return LogisticsContract13
	case Logistics21Channel:
		return LogisticsContract21
	case Logistics22Channel:
		return LogisticsContract22
	case Logistics23Channel:
		return LogisticsContract23
	default:
		return ""
	}
}
