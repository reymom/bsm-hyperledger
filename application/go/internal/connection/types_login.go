package connection

import "strings"

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
		return []Channel{Public1Channel, Public1Channel, Logistics11Channel, Logistics21Channel}
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
