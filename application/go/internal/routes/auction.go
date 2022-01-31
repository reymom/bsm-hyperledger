package routes

import (
	"encoding/json"
	"net/http"

	"github.com/reymom/bsm-hyperledger/application/go/internal/connection"
	"github.com/rs/zerolog/log"
	"github.com/urfave/negroni"
)

func generateAuctionRoutes(mux *http.ServeMux) error {
	mux.Handle("/auctions", negroni.New(negroni.Wrap(http.HandlerFunc(auctionsHandler))))

	return nil
}

func auctionsHandler(w http.ResponseWriter, r *http.Request) {
	var e error

	template := templates.Lookup("auctions")
	if template == nil {
		log.Err(e).Msg("Error while looking up \"auctions\" template")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	auctionsJSON, e := networkContracts[connection.Public1Channel].GwContract.EvaluateTransaction("GetAllAuctions", "", "", "")
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
