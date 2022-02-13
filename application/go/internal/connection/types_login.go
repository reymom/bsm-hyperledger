package connection

import (
	"fmt"
	"strings"
)

//organisations
type Organization string

const (
	Supplier1 Organization = "supplier1"
	Supplier2 Organization = "supplier2"
	Buyer1    Organization = "buyer1"
	Buyer2    Organization = "buyer2"
	Buyer3    Organization = "buyer3"
	Logistics Organization = "logistics"
)

type Login struct {
	Name     Organization `json:"name"`
	Password string       `json:"password"`
}

func (o Organization) getMSP() string {
	return strings.Title(string(o)) + "MSP"
}

type UsersLoginMap map[Organization]string

func IsRegistered(u UsersLoginMap, org Organization, pw string) bool {
	if uPw, ok := u[org]; ok {
		return uPw == pw
	}

	return false
}

func (o Organization) getNetworks() []Channel {
	switch o {
	case Supplier1:
		return []Channel{Public1Channel, LogisticsChannel}
	case Supplier2:
		return []Channel{Public2Channel, LogisticsChannel}
	case Buyer1:
		return []Channel{Public1Channel, Public2Channel, LogisticsChannel}
	case Buyer2:
		return []Channel{Public1Channel, Public2Channel, LogisticsChannel}
	case Buyer3:
		return []Channel{Public1Channel, Public2Channel, LogisticsChannel}
	case Logistics:
		return []Channel{LogisticsChannel}
	default:
		return []Channel{}
	}
}

func (o Organization) GetPublicNetwork() Channel {
	if o == Supplier1 {
		return Public1Channel
	} else if o == Supplier2 {
		return Public2Channel
	}
	return Channel("")
}

func (o Organization) GetAuctionCollections(ch Channel) string {
	switch o {
	case Buyer1:
		if ch == Public1Channel {
			return "12"
		} else if ch == Public2Channel {
			return ""
		} else {
			return ""
		}
	case Buyer2:
		if ch == Public1Channel {
			return "12"
		} else if ch == Public2Channel {
			return "23"
		} else {
			return ""
		}
	case Buyer3:
		if ch == Public1Channel {
			return ""
		} else if ch == Public2Channel {
			return "23"
		} else {
			return ""
		}
	case Supplier1:
		if ch == Public1Channel {
			return "12"
		} else if ch == Public2Channel {
			return ""
		} else {
			return ""
		}
	case Supplier2:
		if ch == Public1Channel {
			return ""
		} else if ch == Public2Channel {
			return "23"
		} else {
			return ""
		}
	default:
		return ""
	}
}

func (o Organization) GetLogisticsCollectionsNums() [][2]string {
	switch o {
	case Buyer1:
		return [][2]string{{"1", "1"}, {"2", "1"}}
	case Buyer2:
		return [][2]string{{"1", "2"}, {"2", "2"}}
	case Buyer3:
		return [][2]string{{"1", "3"}, {"2", "3"}}
	case Supplier1:
		return [][2]string{{"1", "1"}, {"1", "2"}, {"1", "3"}}
	case Supplier2:
		return [][2]string{{"2", "1"}, {"2", "2"}, {"2", "3"}}
	case Logistics:
		return [][2]string{
			{"1", "1"}, {"1", "2"}, {"1", "3"}, {"2", "1"}, {"2", "2"}, {"2", "3"},
		}
	default:
		return [][2]string{}
	}
}

func (o Organization) GetAddress() (string, string, string, string, error) {
	var e error
	switch o {
	case Buyer1:
		return "Spain", "Barcelona", "c/ Muntaner", "43", e
	case Buyer2:
		return "Spain", "Zaragoza", "c/ Ramon y Cajal", "14", e
	case Buyer3:
		return "France", "Toulouse", "rue Alsace-Lorraine", "3", e
	default:
		return "", "", "", "", fmt.Errorf("organization without address")
	}
}
