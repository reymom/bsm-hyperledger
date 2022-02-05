package routes

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/reymom/bsm-hyperledger/application/go/internal/connection"
	"github.com/rs/zerolog/log"
	"github.com/urfave/negroni"
)

func generateAuctionRoutes(mux *http.ServeMux) error {
	mux.Handle("/auctions/list", negroni.New(negroni.Wrap(http.HandlerFunc(auctionsListHandler))))
	mux.Handle("/auctions/create", negroni.New(negroni.Wrap(http.HandlerFunc(auctionCreateHandler))))
	mux.Handle("/auctions/create/submit", negroni.New(negroni.Wrap(http.HandlerFunc(auctionSubmitHandler))))
	mux.Handle("/auctions/close", negroni.New(negroni.Wrap(http.HandlerFunc(auctionCloseHandler))))
	mux.Handle("/auctions/end", negroni.New(negroni.Wrap(http.HandlerFunc(auctionEndHandler))))
	return nil
}

func auctionsListHandler(w http.ResponseWriter, r *http.Request) {
	var (
		e            error
		auctionsJSON []byte
	)

	channel := r.FormValue("channel")

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

	auctions := make([]*Auction, 0)
	if channel != "" {
		auctionsJSON, e = sessionStore.NetworkContracts[connection.Channel(channel)].GwContract.EvaluateTransaction("GetAllAuctions", "", "", "")
		if e != nil {
			log.Err(e).Msg("Error while getting auctions from hyperledger state")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if auctionsJSON != nil {
			e = json.Unmarshal(auctionsJSON, &auctions)
			if e != nil {
				log.Err(e).Msg("Error while unmarshaling auctions")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	} else {
		auctionsTmp := make([]*Auction, 0)
		for channel, contract := range sessionStore.NetworkContracts {
			if !strings.Contains(string(channel), "logistics") {
				auctionsJSON, e = contract.GwContract.EvaluateTransaction("GetAllAuctions", "", "", "")
				if e != nil {
					log.Err(e).Msg("Error while getting auctions from hyperledger state")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				if auctionsJSON != nil {
					e = json.Unmarshal(auctionsJSON, &auctionsTmp)
					if e != nil {
						log.Err(e).Msg("Error while unmarshaling auctions")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					auctions = append(auctions, auctionsTmp...)
				}
			}
		}
	}

	m := struct {
		Auctions []*Auction
		Channel  string
	}{
		Auctions: auctions,
		Channel:  channel,
	}

	w.Header().Set("Content-Type", "text/html")
	e = template.Execute(w, m)
	if e != nil {
		log.Err(e).Msg("Error while executing auctions template")
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}
}

func auctionCreateHandler(w http.ResponseWriter, r *http.Request) {
	var e error

	if !loggedIn {
		if loggedIn, e = sessionStore.CheckLoginFromSession(r, connectionConfig.UsersLoginMap); !loggedIn {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
	}

	template := templates.Lookup("createAuction")
	if template == nil {
		log.Err(e).Msg("Error while looking up \"auctions\" template")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	channels := make([]connection.Channel, 0)
	for channel := range sessionStore.NetworkContracts {
		channels = append(channels, channel)
	}
	m := struct {
		Channels []connection.Channel
	}{
		Channels: channels,
	}

	w.Header().Set("Content-Type", "text/html")
	e = template.Execute(w, m)
	if e != nil {
		log.Err(e).Msg("Error while executing auctions template")
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}
}

func auctionSubmitHandler(w http.ResponseWriter, r *http.Request) {

	if !loggedIn {
		if loggedIn, _ = sessionStore.CheckLoginFromSession(r, connectionConfig.UsersLoginMap); !loggedIn {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
	}

	_, e := sessionStore.NetworkContracts[connection.Channel(
		sessionStore.Login.Name.GetPublicNetwork())].GwContract.SubmitTransaction(
		"CreateAuction", r.FormValue("isPrivate"), r.FormValue("collection"), r.FormValue("material"),
		r.FormValue("form"), r.FormValue("weight"), r.FormValue("minPrice"))
	if e != nil {
		log.Err(e).Msg("Error while submiting auction creation to the hyperledger state")
		w.WriteHeader(http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/auctions/list", http.StatusSeeOther)
}

func auctionCloseHandler(w http.ResponseWriter, r *http.Request) {

	if !loggedIn {
		if loggedIn, _ = sessionStore.CheckLoginFromSession(r, connectionConfig.UsersLoginMap); !loggedIn {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
	}

	_, e := sessionStore.NetworkContracts[connection.Channel(
		sessionStore.Login.Name.GetPublicNetwork())].GwContract.SubmitTransaction(
		"CloseAuction", r.FormValue("auctionID"))
	if e != nil {
		log.Err(e).Msg("Error while writing close auction in the hyperledger state")
		w.WriteHeader(http.StatusInternalServerError)
	}

	redirectPath := "/auctions/list"
	if r.FormValue("channel") != "" {
		redirectPath += "?channel=" + r.FormValue("channel")
	}
	http.Redirect(w, r, redirectPath, http.StatusSeeOther)
}

func auctionEndHandler(w http.ResponseWriter, r *http.Request) {

	if !loggedIn {
		if loggedIn, _ = sessionStore.CheckLoginFromSession(r, connectionConfig.UsersLoginMap); !loggedIn {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
	}

	_, e := sessionStore.NetworkContracts[connection.Channel(
		sessionStore.Login.Name.GetPublicNetwork())].GwContract.SubmitTransaction(
		"EndAuction", r.FormValue("auctionID"))
	if e != nil {
		log.Err(e).Msg("Error while writing end auction in the hyperledger state")
		w.WriteHeader(http.StatusInternalServerError)
	}

	redirectPath := "/auctions/list"
	if r.FormValue("channel") != "" {
		redirectPath += "?channel=" + r.FormValue("channel")
	}
	http.Redirect(w, r, redirectPath, http.StatusSeeOther)
}
