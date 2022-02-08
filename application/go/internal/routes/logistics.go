package routes

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/reymom/bsm-hyperledger/application/go/internal/connection"
	"github.com/rs/zerolog/log"
	"github.com/urfave/negroni"
)

func generateLogisticsRoutes(mux *http.ServeMux) error {

	mux.Handle("/delivery/list", negroni.New(negroni.Wrap(http.HandlerFunc(orderListHandler))))
	mux.Handle("/delivery/updateStatus", negroni.New(negroni.Wrap(http.HandlerFunc(orderUpdateStatusHandler))))
	return nil

}

func orderListHandler(w http.ResponseWriter, r *http.Request) {
	var (
		e              error
		deliveriesJSON []byte
	)

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
		deliveriesJSON, e = gw.GwContract.EvaluateTransaction("GetAllDeliveries")
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
				deliveriesJSON, e = contract.GwContract.EvaluateTransaction("GetAllDeliveries")
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

func orderUpdateStatusHandler(w http.ResponseWriter, r *http.Request) {

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
	}

	http.Redirect(w, r, redirectPath, http.StatusSeeOther)
}
