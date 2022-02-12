package routes

import (
	"encoding/json"
	"net/http"

	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"github.com/reymom/bsm-hyperledger/application/go/internal/connection"
	"github.com/rs/zerolog/log"
	"github.com/urfave/negroni"
)

func generateLogisticsRoutes(mux *http.ServeMux) error {

	mux.Handle("/delivery/list", negroni.New(negroni.Wrap(http.HandlerFunc(deliveryListHandler))))
	mux.Handle("/delivery/create", negroni.New(negroni.Wrap(http.HandlerFunc(deliveryCreateHandler))))
	mux.Handle("/delivery/updateStatus", negroni.New(negroni.Wrap(http.HandlerFunc(deliveryUpdateStatusHandler))))
	mux.Handle("/delivery/history", negroni.New(negroni.Wrap(http.HandlerFunc(deliveryHistoryHandler))))
	return nil

}

func deliveryListHandler(w http.ResponseWriter, r *http.Request) {
	var e error

	if !loggedIn {
		if loggedIn, e = sessionStore.CheckLoginFromSession(r, connectionConfig.UsersLoginMap); !loggedIn {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
	}

	template := templates.Lookup("delivery")
	if template == nil {
		log.Err(e).Msg("Error while looking up \"delivery\" template")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	deliveries := make([]*Delivery, 0)
	gw := sessionStore.NetworkContracts[connection.LogisticsChannel]
	orgNums := sessionStore.Login.Name.GetLogisticsCollectionsNums()
	for _, orgNum := range orgNums {
		tmpDeliveries := make([]*Delivery, 0)

		endorsingPeerOption := gateway.WithEndorsingPeers(connection.Logistics.GetEndorsingPeer())
		txn, e := gw.GwContract.CreateTransaction("GetAllDeliveries", endorsingPeerOption)
		if e != nil {
			log.Err(e).Msg("Error while creating transaction")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		deliveriesJSON, e := txn.Evaluate(orgNum[0], orgNum[1])
		if e != nil {
			log.Err(e).Msg("Error while submiting transaction")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if e != nil {
			log.Err(e).Msg("Error while getting deliveries from hyperledger state")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if deliveriesJSON != nil {
			e = json.Unmarshal(deliveriesJSON, &tmpDeliveries)
			if e != nil {
				log.Err(e).Msg("Error while unmarshaling deliveries")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			deliveries = append(deliveries, tmpDeliveries...)
		}
	}

	auctionID := r.FormValue("auctionID")
	m := struct {
		AuctionID  string
		Deliveries []*Delivery
	}{
		AuctionID:  auctionID,
		Deliveries: deliveries,
	}

	w.Header().Set("Content-Type", "text/html")
	e = template.Execute(w, m)
	if e != nil {
		log.Err(e).Msg("Error while executing delivery template")
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}

}

func deliveryCreateHandler(w http.ResponseWriter, r *http.Request) {
	var e error

	if !loggedIn {
		if loggedIn, _ = sessionStore.CheckLoginFromSession(r, connectionConfig.UsersLoginMap); !loggedIn {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
	}

	winner := r.FormValue("winner")
	auctionID := r.FormValue("auctionID")
	//create delivery
	country, city, street, number, e := connection.Organization(winner).GetAddress()
	if e != nil {
		log.Err(e).Msg("Error while getting address for delivery on the hyperledger state")
		http.Redirect(w, r, "/delivery/list", http.StatusSeeOther)
		return
	}

	_, e = sessionStore.NetworkContracts[connection.LogisticsChannel].GwContract.SubmitTransaction(
		"CreateAuctionDelivery", auctionID, winner, "LogisticsMSP", country, city, street, number)
	if e != nil {
		log.Err(e).Msg("Error while getting creating delivery on the hyperledger state")
		http.Redirect(w, r, "/delivery/list", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/delivery/list?auctionID="+auctionID, http.StatusSeeOther)
}

func deliveryUpdateStatusHandler(w http.ResponseWriter, r *http.Request) {

	if !loggedIn {
		if loggedIn, _ = sessionStore.CheckLoginFromSession(r, connectionConfig.UsersLoginMap); !loggedIn {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
	}

	redirectPath := "/delivery/list"

	supplierNum := string(r.FormValue("supplier")[len(r.FormValue("supplier"))-4])
	buyerNum := string(r.FormValue("buyer")[len(r.FormValue("buyer"))-1])

	endorsingPeerOption := gateway.WithEndorsingPeers(connection.Logistics.GetEndorsingPeer())
	txn, e := sessionStore.NetworkContracts[connection.LogisticsChannel].GwContract.CreateTransaction("UpdateDeliveryStatus", endorsingPeerOption)
	if e != nil {
		log.Err(e).Msg("Error while creating transaction")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, e = txn.Submit(supplierNum, buyerNum, r.FormValue("auctionID"), r.FormValue("toStatus"))
	if e != nil {
		log.Err(e).Msg("Error when submiting transaction to change status of delivery")
		http.Redirect(w, r, redirectPath, http.StatusSeeOther)
		return
	}

	// _, e := sessionStore.NetworkContracts[connection.LogisticsChannel].GwContract.SubmitTransaction(
	// 	"UpdateDeliveryStatus", supplierNum, buyerNum, r.FormValue("auctionID"), r.FormValue("toStatus"))
	// if e != nil {
	// 	log.Err(e).Msg("Error when updating status of delivery")
	// 	http.Redirect(w, r, redirectPath, http.StatusSeeOther)
	// 	return
	// }

	http.Redirect(w, r, redirectPath, http.StatusSeeOther)
}

func deliveryHistoryHandler(w http.ResponseWriter, r *http.Request) {
	var e error

	if !loggedIn {
		if loggedIn, _ = sessionStore.CheckLoginFromSession(r, connectionConfig.UsersLoginMap); !loggedIn {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
	}

	template := templates.Lookup("deliveryHistory")
	if template == nil {
		log.Err(e).Msg("Error while looking up \"deliveryHistory\" template")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	historicJSON, e := sessionStore.NetworkContracts[connection.LogisticsChannel].GwContract.EvaluateTransaction(
		"GetDeliveryHistory", r.FormValue("auctionID"))
	if e != nil {
		log.Err(e).Msg("Error when getting delivery history")
		http.Redirect(w, r, "/delivery/list", http.StatusSeeOther)
		return
	}

	historic := make([]*DeliveryHistory, 0)
	if historicJSON != nil {
		e = json.Unmarshal(historicJSON, &historic)
		if e != nil {
			log.Err(e).Msg("Error while unmarshaling history of delivery")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	m := struct {
		History    []*DeliveryHistory
		DeliveryID string
	}{
		History:    historic,
		DeliveryID: r.FormValue("auctionID"),
	}

	w.Header().Set("Content-Type", "text/html")
	e = template.Execute(w, m)
	if e != nil {
		log.Err(e).Msg("Error while executing delivery template")
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}
}
