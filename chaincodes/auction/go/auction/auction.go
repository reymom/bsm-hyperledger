package auction

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const (
	bidKeyType         = "bid"
	supplierMSPPreffix = "Supplier"
	buyerMSPPreffix    = "Buyer"
)

func (s *SmartContract) CreateAuction(
	ctx contractapi.TransactionContextInterface, private bool,
	collectionOrgNums, steelType, form string, weight, minPrice uint, durationHours float64,
) error {

	clientID, err := s.GetSubmittingClientIdentity(ctx)
	if err != nil {
		return fmt.Errorf("failed to get client identity %v", err)
	}

	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get client identity %v", err)
	}
	if !strings.Contains(string(clientOrgID), supplierMSPPreffix) {
		return fmt.Errorf("just suppliers can create auctions")
	}

	var privateCollectionName string
	if private {
		privateCollectionName, err = getPrivateCollectionChannel(ctx, collectionOrgNums)
		if err != nil {
			return fmt.Errorf("failed to get private collection name: %v", err)
		}
	}

	bids := make(map[string]Bid)
	auction := Auction{
		ID:             uuid.NewString(),
		ClientID:       clientID,
		StartDate:      time.Now(),
		EndDate:        time.Now().Add(time.Duration(durationHours) * time.Hour),
		IsPrivate:      private,
		CollectionName: privateCollectionName,
		Type:           steelType,
		Form:           form,
		Weight:         weight,
		Seller:         clientOrgID,
		Bidders:        []string{},
		Bids:           bids,
		Winner:         "",
		MinPrice:       minPrice,
		Price:          0,
		Status:         opened,
	}

	auctionJSON, err := json.Marshal(auction)
	if err != nil {
		return err
	}

	// put auction into state
	if private {
		// transientMap, err := ctx.GetStub().GetTransient()
		// if err != nil {
		// 	return fmt.Errorf("error getting transient: %v", err)
		// }
		err = ctx.GetStub().PutPrivateData(privateCollectionName, auction.ID, auctionJSON)
		if err != nil {
			return fmt.Errorf("failed to put private data: %v", err)
		}
	} else {
		err = ctx.GetStub().PutState(auction.ID, auctionJSON)
		if err != nil {
			return fmt.Errorf("failed to put auction in public data: %v", err)
		}
	}

	// set the seller of the auction as an endorser
	// err = setAssetStateBasedEndorsement(ctx, auction.ID, clientOrgID)
	// if err != nil {
	// 	return fmt.Errorf("failed setting state based endorsement for new organization: %v", err)
	// }

	return nil
}

func (s *SmartContract) Bid(ctx contractapi.TransactionContextInterface, private bool, auctionID, collectionOrgNums string, price uint) error {
	var (
		err                   error
		privateCollectionName string
	)

	err = verifyClientOrgMatchesPeerOrg(ctx)
	if err != nil {
		return fmt.Errorf("cannot store bid on this peer, not a member of this org: Error %v", err)
	}

	clientOrgID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get client identity %v", err)
	}
	if !strings.Contains(string(clientOrgID), buyerMSPPreffix) {
		return fmt.Errorf("just buyers can submit bids")
	}

	var auction *Auction
	if private {
		privateCollectionName, err = getPrivateCollectionChannel(ctx, collectionOrgNums)
		if err != nil {
			return fmt.Errorf("failed to get private collection name: %v", err)
		}
		auction, err = s.QueryPrivateAuction(ctx, auctionID, collectionOrgNums)
		if err != nil {
			return fmt.Errorf("failed to query private auction with ID: %s", auctionID)
		}
	} else {
		auction, err = s.QueryAuction(ctx, auctionID)
		if err != nil {
			return fmt.Errorf("failed to query auction with ID: %s", auctionID)
		}
	}

	if auction.Status != opened {
		return fmt.Errorf("auction must be opened to bid")
	}

	if auction.MinPrice > price {
		return fmt.Errorf("price must be above minimum price %v", auction.MinPrice)
	}
	bid := Bid{
		ID:    uuid.NewString(),
		Buyer: clientOrgID,
		Price: price,
	}

	// err = addAssetStateBasedEndorsement(ctx, auctionID, clientOrgID)
	// if err != nil {
	// 	return fmt.Errorf("failed setting state based endorsement for new organization: %v", err)
	// }

	fmt.Printf("auction before = %+v\n", auction)
	auction.Bidders = append(auction.Bidders, clientOrgID)
	auction.Bids[clientOrgID] = bid
	fmt.Printf("auction after = %+v\n", auction)

	updatedAuctionJSON, err := json.Marshal(auction)
	if err != nil {
		return fmt.Errorf("error marshaling new auction: %v", err)
	}

	if private {
		err = ctx.GetStub().PutPrivateData(privateCollectionName, auctionID, updatedAuctionJSON)
		if err != nil {
			return fmt.Errorf("failed to input price into collection: %v", err)
		}
	} else {
		fmt.Println("putting new state")
		err = ctx.GetStub().PutState(auctionID, updatedAuctionJSON)
		if err != nil {
			return fmt.Errorf("failed to put auction in public data: %v", err)
		}
	}

	return nil
}

func (s *SmartContract) FinishAuction(ctx contractapi.TransactionContextInterface, private bool, auctionID, collectionOrgNums string) error {

	var (
		err                   error
		privateCollectionName string
		auction               *Auction
	)
	if private {
		privateCollectionName, err = getPrivateCollectionChannel(ctx, collectionOrgNums)
		if err != nil {
			return fmt.Errorf("failed to get private collection name: %v", err)
		}
		auction, err = s.QueryPrivateAuction(ctx, auctionID, collectionOrgNums)
		if err != nil {
			return fmt.Errorf("failed to query auction with ID: %s", auctionID)
		}
	} else {
		auction, err = s.QueryAuction(ctx, auctionID)
		if err != nil {
			return fmt.Errorf("failed to query auction with ID: %s", auctionID)
		}
	}

	clientID, err := s.GetSubmittingClientIdentity(ctx)
	if err != nil {
		return fmt.Errorf("failed to get client identity %v", err)
	}

	if auction.ClientID != clientID {
		return fmt.Errorf("auction can only be closed by owner: %v", err)
	}

	if auction.Status != opened {
		return fmt.Errorf("cannot finish auction that is not open")
	}

	if auction.EndDate.After(time.Now()) {
		return fmt.Errorf("auction cannot be closed before %v", auction.EndDate)
	}

	auction.Status = finished
	finishedAuction, _ := s.setWinnerOfAuction(ctx, auction)
	finishedAuctionJSON, _ := json.Marshal(finishedAuction)

	if private {
		err = ctx.GetStub().PutPrivateData(privateCollectionName, auctionID, finishedAuctionJSON)
		if err != nil {
			return fmt.Errorf("failed to finish auction: %v", err)
		}
	} else {
		err = ctx.GetStub().PutState(auctionID, finishedAuctionJSON)
		if err != nil {
			return fmt.Errorf("failed to finish auction: %v", err)
		}
	}

	return nil
}
