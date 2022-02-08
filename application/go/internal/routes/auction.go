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

func generateAuctionRoutes(mux *http.ServeMux) error {
	mux.Handle("/auctions/list", negroni.New(negroni.Wrap(http.HandlerFunc(auctionsListHandler))))
	mux.Handle("/auctions/create", negroni.New(negroni.Wrap(http.HandlerFunc(auctionCreateHandler))))
	mux.Handle("/auctions/create/submit", negroni.New(negroni.Wrap(http.HandlerFunc(auctionSubmitHandler))))
	mux.Handle("/auctions/finish", negroni.New(negroni.Wrap(http.HandlerFunc(auctionFinishHandler))))
	return nil
}

func auctionsListHandler(w http.ResponseWriter, r *http.Request) {
	var (
		e                   error
		auctionsJSON        []byte
		privateAuctionsJSON []byte
	)

	o := sessionStore.Login.Name
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

	type channelAuction struct {
		Auctions []*Auction
		Channel  connection.Channel
	}
	auctions := make([]*Auction, 0)
	privateAuctions := make([]*Auction, 0)
	channelAuctions := make([]channelAuction, 0)
	if channel != "" {
		ch := connection.Channel(channel)
		gw := sessionStore.NetworkContracts[ch]
		auctionsTmp := make([]*Auction, 0)
		auctionsJSON, e = gw.GwContract.EvaluateTransaction("GetAllAuctions")
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
			auctionsTmp = append(auctionsTmp, auctions...)
		}

		if o.GetCollections(ch) != "" {
			privateAuctionsJSON, e = gw.GwContract.EvaluateTransaction("GetAllPrivateAuctions", o.GetCollections(ch))
			if e != nil {
				log.Err(e).Msg("Error while getting private auctions from hyperledger state")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			fmt.Println(" privateAuctionsJSON = ", privateAuctionsJSON)
			if privateAuctionsJSON != nil {
				e = json.Unmarshal(privateAuctionsJSON, &privateAuctions)
				if e != nil {
					log.Err(e).Msg("Error while unmarshaling private auctions")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				auctionsTmp = append(auctionsTmp, privateAuctions...)
			}
		}

		channelAuctions = append(channelAuctions, channelAuction{
			Auctions: auctionsTmp,
			Channel:  connection.Channel(channel),
		})
	} else {
		for channel, contract := range sessionStore.NetworkContracts {
			auctionsTmp := make([]*Auction, 0)
			if !strings.Contains(string(channel), "logistics") {
				auctionsJSON, e = contract.GwContract.EvaluateTransaction("GetAllAuctions")
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
					auctionsTmp = append(auctionsTmp, auctions...)
				}

				if o.GetCollections(channel) != "" {
					privateAuctionsJSON, e = contract.GwContract.EvaluateTransaction("GetAllPrivateAuctions", o.GetCollections(channel))
					if e != nil {
						log.Err(e).Msg("Error while getting private auctions from hyperledger state")
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					if privateAuctionsJSON != nil {
						e = json.Unmarshal(privateAuctionsJSON, &privateAuctions)
						if e != nil {
							log.Err(e).Msg("Error while unmarshaling private auctions")
							w.WriteHeader(http.StatusInternalServerError)
							return
						}
						auctionsTmp = append(auctionsTmp, privateAuctions...)
					}
				}

				channelAuctions = append(channelAuctions, channelAuction{
					Auctions: auctionsTmp,
					Channel:  channel,
				})
			}
		}
	}

	m := struct {
		ChannelAuctions []channelAuction
	}{
		ChannelAuctions: channelAuctions,
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
		"CreateAuction", r.FormValue("isPrivate"), r.FormValue("collNum"), r.FormValue("material"),
		r.FormValue("form"), r.FormValue("weight"), r.FormValue("minPrice"), r.FormValue("hours"))
	if e != nil {
		log.Err(e).Msg("Error while submiting auction creation to the hyperledger state")
		http.Redirect(w, r, "/auctions/list", http.StatusSeeOther)
	}

	http.Redirect(w, r, "/auctions/list", http.StatusSeeOther)
}

func auctionFinishHandler(w http.ResponseWriter, r *http.Request) {

	if !loggedIn {
		if loggedIn, _ = sessionStore.CheckLoginFromSession(r, connectionConfig.UsersLoginMap); !loggedIn {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
	}

	redirectPath := "/auctions/list"
	if r.FormValue("channel") != "" {
		redirectPath += "?channel=" + r.FormValue("channel")
	}

	var winner string
	winnerJSON, e := sessionStore.NetworkContracts[connection.Channel(
		sessionStore.Login.Name.GetPublicNetwork())].GwContract.SubmitTransaction(
		"FinishAuction", r.FormValue("private"), r.FormValue("auctionID"), r.FormValue("colNums"))
	if e != nil {
		log.Err(e).Msg("Error while finishing auction in the hyperledger state")
		http.Redirect(w, r, redirectPath, http.StatusSeeOther)
	}
	if winnerJSON != nil {
		e = json.Unmarshal(winnerJSON, &winner)
		if e != nil {
			log.Err(e).Msg("Error while unmarshaling winner")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	//create delivery
	country, city, street, number := connection.Organization(winner).GetAddress()
	_, e = sessionStore.NetworkContracts[connection.Channel(
		sessionStore.Login.Name.GetLogisticsChannel(connection.Organization(winner)))].GwContract.EvaluateTransaction(
		"CreateDelivery", r.FormValue("auctionID"), winner, "LogisticsMSP", country, city, street, number)
	if e != nil {
		log.Err(e).Msg("Error while getting auctions from hyperledger state")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, redirectPath, http.StatusSeeOther)
}
