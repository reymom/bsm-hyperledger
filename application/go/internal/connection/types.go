package connection

type ConnectionProfile struct {
	OrgName  organisation
	Channels []channel
}

func NewConnectionProfile(orgName organisation, channels ...channel) (*ConnectionProfile, error) {
	return &ConnectionProfile{
		OrgName:  orgName,
		Channels: channels,
	}, nil
}

type organisation string

const (
	Supplier1 organisation = "supplier1"
	Supplier2 organisation = "supplier2"
	Buyer1    organisation = "buyer1"
	Buyer2    organisation = "buyer2"
	Buyer3    organisation = "buyer3"
	Logistics organisation = "logistics"
)

type channel string

const (
	Public1Channel     channel = "public1channel"
	Public2Channel     channel = "public2channel"
	Private112Channel  channel = "private112channel"
	Private123Channel  channel = "private123channel"
	Private212Channel  channel = "private212channel"
	Private223Channel  channel = "private223channel"
	Logistics11Channel channel = "logistics11channel"
	Logistics12Channel channel = "logistics12channel"
	Logistics13Channel channel = "logistics13channel"
	Logistics21Channel channel = "logistics21channel"
	Logistics22Channel channel = "logistics22channel"
	Logistics23Channel channel = "logistics23channel"
)
