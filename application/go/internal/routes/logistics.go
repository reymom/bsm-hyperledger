package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

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

	reqChannel := r.FormValue("channel")

	deliveries := make([]*Delivery, 0)
	if reqChannel != "" {
		ch := connection.Channel(reqChannel)
		gw := sessionStore.NetworkContracts[ch]
		deliveriesJSON, e := gw.GwContract.EvaluateTransaction("GetAllDeliveries")
		if e != nil {
			log.Err(e).Msg("Error while getting deliveries from hyperledger state")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if deliveriesJSON != nil {
			e = json.Unmarshal(deliveriesJSON, &deliveries)
			if e != nil {
				log.Err(e).Msg("Error while unmarshaling deliveries")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	} else {
		for channel, contract := range sessionStore.NetworkContracts {
			deliveriesTmp := make([]*Delivery, 0)
			if strings.Contains(string(channel), "logistics") {
				deliveriesJSON, e := contract.GwContract.EvaluateTransaction("GetAllDeliveries")
				if e != nil {
					log.Err(e).Msg("Error while getting deliveries from hyperledger state")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				if deliveriesJSON != nil {
					e = json.Unmarshal(deliveriesJSON, &deliveriesTmp)
					if e != nil {
						log.Err(e).Msg("Error while unmarshaling deliveries")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					deliveries = append(deliveries, deliveriesTmp...)
				}
			}
		}
	}

	auctionID := r.FormValue("auctionID")
	m := struct {
		AuctionID  string
		Deliveries []*Delivery
		Channel    string
	}{
		AuctionID:  auctionID,
		Deliveries: deliveries,
		Channel:    reqChannel,
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
	fmt.Printf("winner = %v, %T", winner, winner)
	fmt.Println("auctionID = ", auctionID)
	fmt.Println("winner = ", winner)
	fmt.Println("winner = ", winner)
	//create delivery
	country, city, street, number, e := connection.Organization(winner).GetAddress()
	if e != nil {
		log.Err(e).Msg("Error while getting address for delivery on the hyperledger state")
		http.Redirect(w, r, "/delivery/list", http.StatusSeeOther)
		return
	}
	fmt.Println(country, city, street, number)
	fmt.Println(sessionStore.Login.Name.GetLogisticsChannel(connection.Organization(winner)))
	_, e = sessionStore.NetworkContracts[sessionStore.Login.Name.GetLogisticsChannel(connection.Organization(winner))].GwContract.EvaluateTransaction(
		"CreateDelivery", auctionID, winner, "LogisticsMSP", country, city, street, number)
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
	if reqChannel := r.FormValue("channel"); reqChannel != "" {
		redirectPath += "?channel=" + reqChannel
	}

	_, e := sessionStore.NetworkContracts[connection.Channel(
		sessionStore.Login.Name.GetLogisticsChannel(
			connection.Organization(r.FormValue("destinyOrg"))))].GwContract.SubmitTransaction(
		"UpdateDeliveryStatus", r.FormValue("auctionID"), r.FormValue("toStatus"))
	if e != nil {
		log.Err(e).Msg("Error when updating status of delivery")
		http.Redirect(w, r, redirectPath, http.StatusSeeOther)
		return
	}

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

	historicJSON, e := sessionStore.NetworkContracts[connection.Channel(
		sessionStore.Login.Name.GetLogisticsChannel(
			connection.Organization(r.FormValue("destinyOrg"))))].GwContract.SubmitTransaction(
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
