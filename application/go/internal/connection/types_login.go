package connection

import (
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
		return []Channel{Public1Channel, Logistics11Channel, Logistics12Channel, Logistics13Channel}
	case Supplier2:
		return []Channel{Public2Channel, Logistics21Channel, Logistics22Channel, Logistics23Channel}
	case Buyer1:
		return []Channel{Public1Channel, Public2Channel, Logistics11Channel, Logistics21Channel}
	case Buyer2:
		return []Channel{Public1Channel, Public2Channel, Logistics12Channel, Logistics22Channel}
	case Buyer3:
		return []Channel{Public1Channel, Public2Channel, Logistics13Channel, Logistics23Channel}
	case Logistics:
		return []Channel{Logistics11Channel, Logistics12Channel, Logistics13Channel, Logistics21Channel, Logistics22Channel, Logistics23Channel}
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

func (o Organization) GetLogisticsChannel(orgDestiny Organization) Channel {
	if orgDestiny == Buyer1 || orgDestiny == Buyer2 || orgDestiny == Buyer3 {
		return Channel("logistics" + string(string(o)[len(o)]) + string(string(orgDestiny)[len(orgDestiny)]) + "channel")
	}
	return Channel("")
}

func (o Organization) GetCollections(ch Channel) string {
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

func (o Organization) GetAddress() (country, city, street, number string) {
	switch o {
	case Buyer1:
		return "Spain", "Barcelona", "c/ Muntaner", "43"
	case Buyer2:
		return "Spain", "Zaragoza", "c/ Ramon y Cajal", "14"
	case Buyer3:
		return "France", "Toulouse", "rue Alsace-Lorraine", "3"
	default:
		return "", "", "", ""
	}
}
