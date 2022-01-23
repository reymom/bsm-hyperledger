package logistics

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func (s *SmartContract) QueryDelivery(ctx contractapi.TransactionContextInterface, deliveryID string) (*SteelDelivery, error) {

	deliveryJSON, err := ctx.GetStub().GetState(deliveryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get delivery object %v: %v", deliveryID, err)
	}
	if deliveryJSON == nil {
		return nil, fmt.Errorf("delivery %v does not exist", deliveryID)
	}

	var delivery *SteelDelivery
	err = json.Unmarshal(deliveryJSON, &delivery)
	if err != nil {
		return nil, err
	}

	return delivery, nil
}

func (s *SmartContract) DeliveryExists(ctx contractapi.TransactionContextInterface, deliveryID string) (bool, error) {

	deliveryJSON, err := ctx.GetStub().GetState(deliveryID)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return deliveryJSON != nil, nil
}

func (s *SmartContract) GetAllDeliveries(ctx contractapi.TransactionContextInterface) ([]*SteelDelivery, error) {

	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
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

func (t *SmartContract) GetDeliveryHistory(ctx contractapi.TransactionContextInterface, deliveryID string) ([]HistoryQueryResult, error) {

	resultsIterator, err := ctx.GetStub().GetHistoryForKey(deliveryID)
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

		var delivery SteelDelivery
		if len(response.Value) > 0 {
			err = json.Unmarshal(response.Value, &delivery)
			if err != nil {
				return nil, err
			}
		} else {
			delivery = SteelDelivery{
				ID: deliveryID,
			}
		}

		record := HistoryQueryResult{
			TxId:      response.TxId,
			Timestamp: response.Timestamp.AsTime(),
			Record:    &delivery,
			Deleted:   response.IsDelete,
		}
		records = append(records, record)
	}

	return records, nil
}
