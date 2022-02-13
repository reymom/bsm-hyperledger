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

func (c Channel) IsAuctionChannel() bool {
	return strings.Contains(string(c), "public")
}

func (c Channel) GetEndorsingPeer() string {
	switch c {
	case LogisticsChannel:
		return "peer0.logistics.steelplatform.com:17051"
	case Public1Channel:
		return "peer0.supplier1.steelplatform.com:7051"
	case Public2Channel:
		return "peer0.supplier2.steelplatform.com:9051"
	default:
		return ""
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
	case LogisticsChannel:
		return LogisticsContract
	default:
		return Contract("")
	}
}
