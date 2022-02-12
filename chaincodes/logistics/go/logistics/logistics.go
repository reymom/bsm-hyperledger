package logistics

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func (s *SmartContract) CreateDelivery(ctx contractapi.TransactionContextInterface,
	auctionID, destinyOrg, deliveryOrgMSPID, country, city, street, number string) error {

	return nil
}

func (s *SmartContract) CreateAuctionDelivery(ctx contractapi.TransactionContextInterface,
	auctionID, destinyOrg, deliveryOrgMSPID, country, city, street, number string) error {

	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get client org identity %v", err)
	}

	supplierNum := string(clientOrgID[len(clientOrgID)-4])
	buyerNum := string(destinyOrg[len(destinyOrg)-1])

	if strings.Contains(clientOrgID, "buyer") {
		return fmt.Errorf("just suppliers or or the delivery company can create an auction delivery")
	}

	collection, err := getPrivateCollection(ctx, supplierNum, buyerNum)
	if err != nil {
		return fmt.Errorf("failed to get collection %v", err)
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

	return ctx.GetStub().PutPrivateData(collection, auctionID, deliveryJSON)
}

func (s *SmartContract) UpdateDeliveryStatus(ctx contractapi.TransactionContextInterface,
	supplierNum, buyerNum string, auctionID string, newStatusUint uint) error {

	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get client org identity %v", err)
	}

	delivery, err := s.QueryDelivery(ctx, supplierNum, buyerNum, auctionID)
	if err != nil {
		return fmt.Errorf("failed to get delivery from collection %v", err)
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

	collection, err := getPrivateCollection(ctx, supplierNum, buyerNum)
	if err != nil {
		return fmt.Errorf("failed to get collection %v", err)
	}
	return ctx.GetStub().PutPrivateData(collection, auctionID, deliveryJSON)
}
