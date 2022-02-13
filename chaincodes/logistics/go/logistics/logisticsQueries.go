package logistics

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func (s *SmartContract) QueryDelivery(ctx contractapi.TransactionContextInterface, supplierNum, buyerNum, auctionID string) (*SteelDelivery, error) {

	collection, err := getPrivateCollection(ctx, supplierNum, buyerNum)
	if err != nil {
		return nil, fmt.Errorf("failed to get collection %v", err)
	}

	deliveryJSON, err := ctx.GetStub().GetPrivateData(collection, auctionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get delivery object %v: %v", auctionID, err)
	}
	if deliveryJSON == nil {
		return nil, fmt.Errorf("delivery %v does not exist", auctionID)
	}

	var delivery *SteelDelivery
	err = json.Unmarshal(deliveryJSON, &delivery)
	if err != nil {
		return nil, err
	}

	return delivery, nil
}

func (s *SmartContract) DeliveryExists(ctx contractapi.TransactionContextInterface,
	supplierNum, buyerNum, auctionID string) (bool, error) {

	collection, err := getPrivateCollection(ctx, supplierNum, buyerNum)
	if err != nil {
		return false, fmt.Errorf("failed to get collection %v", err)
	}

	deliveryJSON, err := ctx.GetStub().GetPrivateData(collection, auctionID)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return deliveryJSON != nil, nil
}

func (s *SmartContract) GetAllDeliveries(ctx contractapi.TransactionContextInterface, supplierNum, buyerNum string) ([]*SteelDelivery, error) {

	collection, err := getPrivateCollection(ctx, supplierNum, buyerNum)
	if err != nil {
		return nil, fmt.Errorf("failed to get collection %v", err)
	}

	resultsIterator, err := ctx.GetStub().GetPrivateDataByRange(collection, "", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var deliveries []*SteelDelivery
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var delivery SteelDelivery
		err = json.Unmarshal(queryResponse.Value, &delivery)
		if err != nil {
			return nil, err
		}
		deliveries = append(deliveries, &delivery)
	}

	return deliveries, nil
}

func (s *SmartContract) GetDeliveryHistory(ctx contractapi.TransactionContextInterface, auctionID string) ([]HistoryQueryResult, error) {

	resultsIterator, err := ctx.GetStub().GetHistoryForKey(auctionID)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var records []HistoryQueryResult
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var hashTracker HistoryPublicHashTracker
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &hashTracker)
			if err != nil {
				return nil, err
			}
		} else {
			hashTracker = HistoryPublicHashTracker{
				HashID: auctionID,
			}
		}

		record := HistoryQueryResult{
			Record:    &hashTracker,
			TxId:      response.TxId,
			Timestamp: response.Timestamp.AsTime(),
			Deleted:   response.IsDelete,
		}
		records = append(records, record)
	}

	return records, nil
}
