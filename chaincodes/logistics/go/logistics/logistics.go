package logistics

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {

// 	delivery := new(SteelDelivery)
// 	for i := 0; i < 2; i++ {
// 		delivery = SteelDelivery{
// 			ID: uuid.NewString(),
// 			User: uuid.NewString(),
// 			Address: &Address{
// 				Country: "ES",
// 				City: "Barcelona",
// 				ZipCode: uuid.NewString{},
// 				Street: "SomeStreet",
// 				Number: "56",
// 			},
// 			Status: processing,
// 		}

// 		deliveryJSON, err := json.Marshal(delivery)
// 		if err != nil {
// 			return err
// 		}

// 		err = ctx.GetStub().PutState(delivery.ID, deliveryJSON)
// 		if err != nil {
// 			return fmt.Errorf("failed to put to world state. %v", err)
// 		}
// 	}

// 	return nil
// }

func (s *SmartContract) CreateDelivery(ctx contractapi.TransactionContextInterface, deliveryID, country, city, street, number string) error {

	clientID, err := s.GetSubmittingClientIdentity(ctx)
	if err != nil {
		return fmt.Errorf("failed to get client identity %v", err)
	}

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
		ID:          deliveryID,
		DeliveryOrg: clientID,
		Address:     &address,
		Updated:     time.Now(),
		Status:      processing,
	}

	deliveryJSON, err := json.Marshal(delivery)
	if err != nil {
		return fmt.Errorf("failed to put delivery in public data: %v", err)
	}

	err = ctx.GetStub().PutState(deliveryID, deliveryJSON)
	if err != nil {
		return fmt.Errorf("failed to put delivery in public data: %v", err)
	}

	err = setAssetStateBasedEndorsement(ctx, deliveryID, clientOrgID)
	if err != nil {
		return fmt.Errorf("failed setting state based endorsement for new organization: %v", err)
	}

	return nil
}

func (s *SmartContract) UpdateDeliveryStatus(ctx contractapi.TransactionContextInterface, deliveryID string, newStatus statusTypes) error {

	delivery, err := s.QueryDelivery(ctx, deliveryID)
	if err != nil {
		return fmt.Errorf("failed to get delivery from public state %v", err)
	}

	clientID, err := s.GetSubmittingClientIdentity(ctx)
	if err != nil {
		return fmt.Errorf("failed to get client identity %v", err)
	}

	if clientID != delivery.DeliveryOrg {
		return fmt.Errorf("delivery state can only be updated by the delivery organisation")
	}

	status := delivery.Status

	switch newStatus {
	case processing:
		return fmt.Errorf("cannot update to processing")
	case onDelivery:
		if status != processing {
			return fmt.Errorf("delivery comes only after processing")
		}
	case delivered, failed:
		if status != onDelivery {
			return fmt.Errorf("delivered or failed comes only after being on delivery")
		}
	default:
		return fmt.Errorf("status %v unknown\n", newStatus)
	}

	delivery.Status = newStatus

	deliveryJSON, err := json.Marshal(delivery)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(deliveryID, deliveryJSON)
}

func (s *SmartContract) DeleteDelivery(ctx contractapi.TransactionContextInterface, deliveryID string) error {

	exists, err := s.DeliveryExists(ctx, deliveryID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the delivery %s does not exist", deliveryID)
	}

	return ctx.GetStub().DelState(deliveryID)
}
