package logistics

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func (s *SmartContract) CreateDelivery(ctx contractapi.TransactionContextInterface, auctionID, destinyOrg, deliveryOrgMSPID, country, city, street, number string) error {

	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get client org identity %v", err)
	}

	address := Address{
		Country: country,
		City:    city,
		Street:  street,
		Number:  number,
	}

	delivery := SteelDelivery{
		AuctionID:   auctionID,
		Creator:     clientOrgID,
		DestinyOrg:  destinyOrg,
		DeliveryOrg: deliveryOrgMSPID,
		Address:     &address,
		Updated:     time.Now(),
		Status:      processing,
	}

	deliveryJSON, err := json.Marshal(delivery)
	if err != nil {
		return fmt.Errorf("failed to put delivery in public data: %v", err)
	}

	err = ctx.GetStub().PutState(auctionID, deliveryJSON)
	if err != nil {
		return fmt.Errorf("failed to put delivery in public data: %v", err)
	}

	return nil
}

func (s *SmartContract) UpdateDeliveryStatus(ctx contractapi.TransactionContextInterface, auctionID string, newStatusUint uint) error {

	delivery, err := s.QueryDelivery(ctx, auctionID)
	if err != nil {
		return fmt.Errorf("failed to get delivery from public state %v", err)
	}

	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get client org identity %v", err)
	}

	if clientOrgID != delivery.DeliveryOrg {
		return fmt.Errorf("delivery state can only be updated by the delivery organisation")
	}

	status := delivery.Status
	newStatus := statusTypes(newStatusUint)
	switch newStatus {
	case processing:
		return fmt.Errorf("cannot update to processing")
	case onDelivery:
		if status != processing {
			return fmt.Errorf("on delivery comes only after processing status")
		}
	case delivered, failed:
		if status != onDelivery {
			return fmt.Errorf("delivered or failed comes only after being on delivery")
		}
	default:
		return fmt.Errorf("status %v unknown", newStatus)
	}

	delivery.Status = newStatus
	delivery.Updated = time.Now()

	deliveryJSON, err := json.Marshal(delivery)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(auctionID, deliveryJSON)
}
