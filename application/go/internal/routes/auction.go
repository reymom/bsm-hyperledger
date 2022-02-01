package routes

import (
	"encoding/json"
	"net/http"

	"github.com/reymom/bsm-hyperledger/application/go/internal/connection"
	"github.com/rs/zerolog/log"
	"github.com/urfave/negroni"
)

func generateAuctionRoutes(mux *http.ServeMux) error {
	mux.Handle("/auctions/supplier/list", negroni.New(negroni.Wrap(http.HandlerFunc(supplierAuctionsHandler))))
	mux.Handle("/auctions/buyer/list", negroni.New(negroni.Wrap(http.HandlerFunc(buyerAuctionsHandler))))
	mux.Handle("/auctions/create", negroni.New(negroni.Wrap(http.HandlerFunc(auctionCreateHandler))))
	mux.Handle("/auctions/create/submit", negroni.New(negroni.Wrap(http.HandlerFunc(auctionSubmitHandler))))
	return nil
}

func supplierAuctionsHandler(w http.ResponseWriter, r *http.Request) {
	var e error

	if !loggedIn {
		if loggedIn, e = sessionStore.CheckLoginFromSession(r, connectionConfig.UsersLoginMap); !loggedIn {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
	}

	template := templates.Lookup("auctions")
	if template == nil {
		log.Err(e).Msg("Error while looking up \"auctions\" template")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	auctionsJSON, e := sessionStore.NetworkContracts[connection.Public1Channel].GwContract.EvaluateTransaction("GetAllAuctions", "", "", "")
	if e != nil {
		log.Err(e).Msg("Error while getting \"auctions\" from hyperledger state")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var auctions []*Auction
	e = json.Unmarshal(auctionsJSON, &auctions)
	if e != nil {
		log.Err(e).Msg("Error while unmarshaling \"auctions\"")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	m := struct {
		Auctions []*Auction
	}{
		Auctions: auctions,
	}

	w.Header().Set("Content-Type", "text/html")
	e = template.Execute(w, m)
	if e != nil {
		log.Err(e).Msg("Error while executing \"auctions\" template")
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}
}

func buyerAuctionsHandler(w http.ResponseWriter, r *http.Request) {
}

func auctionCreateHandler(w http.ResponseWriter, r *http.Request) {
}

func auctionSubmitHandler(w http.ResponseWriter, r *http.Request) {
}
