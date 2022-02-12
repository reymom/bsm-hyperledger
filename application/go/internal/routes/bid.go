package routes

import (
	"encoding/json"
	"net/http"

	"github.com/reymom/bsm-hyperledger/application/go/internal/connection"
	"github.com/rs/zerolog/log"
	"github.com/urfave/negroni"
)

func generateBidRoutes(mux *http.ServeMux) error {
	mux.Handle("/bid", negroni.New(negroni.Wrap(http.HandlerFunc(bidCreateHandler))))
	mux.Handle("/bid/submit", negroni.New(negroni.Wrap(http.HandlerFunc(bidSubmitHandler))))
	return nil
}

func bidCreateHandler(w http.ResponseWriter, r *http.Request) {
	var e error

	if !loggedIn {
		if loggedIn, e = sessionStore.CheckLoginFromSession(r, connectionConfig.UsersLoginMap); !loggedIn {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
	}

	template := templates.Lookup("createBid")
	if template == nil {
		log.Err(e).Msg("Error while looking up \"createBid\" template")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	channel := connection.Channel(r.FormValue("channel"))
	colNums := r.FormValue("colNums")
	var auctionJSON []byte
	var auction *Auction
	if colNums != "" {
		auctionJSON, e = sessionStore.NetworkContracts[connection.Channel(channel)].GwContract.EvaluateTransaction(
			"QueryPrivateAuction", r.FormValue("auctionID"), colNums)
	} else {
		auctionJSON, e = sessionStore.NetworkContracts[connection.Channel(channel)].GwContract.EvaluateTransaction(
			"QueryAuction", r.FormValue("auctionID"))
	}

	if e != nil {
		log.Err(e).Msg("Error while getting auction from hyperledger state")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if auctionJSON != nil {
		e = json.Unmarshal(auctionJSON, &auction)
		if e != nil {
			log.Err(e).Msg("Error while unmarshaling auctions")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	m := struct {
		Channel connection.Channel
		Auction *Auction
	}{
		Channel: channel,
		Auction: auction,
	}

	w.Header().Set("Content-Type", "text/html")
	e = template.Execute(w, m)
	if e != nil {
		log.Err(e).Msg("Error while executing auctions template")
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}
}

func bidSubmitHandler(w http.ResponseWriter, r *http.Request) {

	if !loggedIn {
		if loggedIn, _ = sessionStore.CheckLoginFromSession(r, connectionConfig.UsersLoginMap); !loggedIn {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
	}

	_, e := sessionStore.NetworkContracts[connection.Channel(r.FormValue("channel"))].GwContract.SubmitTransaction(
		"Bid", r.FormValue("isPrivate"), r.FormValue("auctionID"), r.FormValue("colNums"), r.FormValue("price"))
	if e != nil {
		log.Err(e).Msg("Error while submiting bid to the hyperledger state")
		w.WriteHeader(http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/auctions/list?channel="+r.FormValue("channel"), http.StatusSeeOther)
}
