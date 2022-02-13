package routes

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
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

	o := sessionStore.Login.Name
	channel := connection.Channel(r.FormValue("channel"))

	type channelAuction struct {
		Auctions []*Auction
		Channel  connection.Channel
	}
	channelAuctions := make([]channelAuction, 0)
	if channel != "" {
		auctions := make([]*Auction, 0)
		privateAuctions := make([]*Auction, 0)
		gw := sessionStore.NetworkContracts[channel]
		auctionsTmp := make([]*Auction, 0)
		auctionsJSON, e = gw.GwContract.EvaluateTransaction("GetAllAuctions")
		if e != nil {
			log.Err(e).Msg("Error while getting auctions from hyperledger state")
			http.Redirect(w, r, "/home", http.StatusSeeOther)
			return
		}
		if auctionsJSON != nil {
			e = json.Unmarshal(auctionsJSON, &auctions)
			if e != nil {
				log.Err(e).Msg("Error while unmarshaling auctions")
				http.Redirect(w, r, "/home", http.StatusSeeOther)
				return
			}
			auctionsTmp = append(auctionsTmp, auctions...)
		}

		if o.GetAuctionCollections(channel) != "" {
			endorsingPeerOption := gateway.WithEndorsingPeers(channel.GetEndorsingPeer())
			txn, e := gw.GwContract.CreateTransaction("GetAllPrivateAuctions", endorsingPeerOption)
			if e != nil {
				log.Err(e).Msg("Error while creating transaction")
				http.Redirect(w, r, "/home", http.StatusSeeOther)
				return
			}
			privateAuctionsJSON, e = txn.Evaluate(o.GetAuctionCollections(channel))
			if e != nil {
				log.Err(e).Msg("Error when evaluating get private auctions transaction")
				http.Redirect(w, r, "/home", http.StatusSeeOther)
				return
			}
			if privateAuctionsJSON != nil {
				e = json.Unmarshal(privateAuctionsJSON, &privateAuctions)
				if e != nil {
					log.Err(e).Msg("Error while unmarshaling private auctions")
					http.Redirect(w, r, "/home", http.StatusSeeOther)
					return
				}
				auctionsTmp = append(auctionsTmp, privateAuctions...)
			}
		}

		channelAuctions = append(channelAuctions, channelAuction{
			Auctions: auctionsTmp,
			Channel:  channel,
		})
	} else {
		for channel, contract := range sessionStore.NetworkContracts {
			auctionsTmp := make([]*Auction, 0)
			auctions := make([]*Auction, 0)
			privateAuctions := make([]*Auction, 0)
			if !strings.Contains(string(channel), "logistics") {
				auctionsJSON, e = contract.GwContract.EvaluateTransaction("GetAllAuctions")
				if e != nil {
					log.Err(e).Msg("Error while getting auctions from hyperledger state")
					http.Redirect(w, r, "/home", http.StatusSeeOther)
					return
				}
				if auctionsJSON != nil {
					e = json.Unmarshal(auctionsJSON, &auctions)
					if e != nil {
						log.Err(e).Msg("Error while unmarshaling auctions")
						http.Redirect(w, r, "/home", http.StatusSeeOther)
						return
					}
					auctionsTmp = append(auctionsTmp, auctions...)
				}
				collection := o.GetAuctionCollections(channel)
				if collection != "" {
					endorsingPeerOption := gateway.WithEndorsingPeers(channel.GetEndorsingPeer())
					txn, e := contract.GwContract.CreateTransaction("GetAllPrivateAuctions", endorsingPeerOption)
					if e != nil {
						log.Err(e).Msg("Error while creating transaction")
						http.Redirect(w, r, "/home", http.StatusSeeOther)
						return
					}
					privateAuctionsJSON, e = txn.Evaluate(collection)
					if e != nil {
						log.Err(e).Msg("Error when evaluating get private auctions transaction")
						http.Redirect(w, r, "/home", http.StatusSeeOther)
						return
					}
					if privateAuctionsJSON != nil {
						e = json.Unmarshal(privateAuctionsJSON, &privateAuctions)
						if e != nil {
							log.Err(e).Msg("Error while unmarshaling private auctions")
							http.Redirect(w, r, "/home", http.StatusSeeOther)
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
		http.Redirect(w, r, "/auctions/list", http.StatusSeeOther)
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

	if r.FormValue("private") == "true" {
		endorsingPeerOption := gateway.WithEndorsingPeers(sessionStore.Login.Name.GetPublicNetwork().GetEndorsingPeer())
		txn, e := sessionStore.NetworkContracts[sessionStore.Login.Name.GetPublicNetwork()].GwContract.CreateTransaction(
			"CreateAuction", endorsingPeerOption)
		if e != nil {
			log.Err(e).Msg("Error while creating transaction")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, e = txn.Submit("true", r.FormValue("collNum"), r.FormValue("material"),
			r.FormValue("form"), r.FormValue("weight"), r.FormValue("minPrice"), r.FormValue("hours"))
		if e != nil {
			log.Err(e).Msg("Error when submiting transaction to create auction")
			http.Redirect(w, r, "/auctions/list", http.StatusSeeOther)
			return
		}
	} else {
		_, e := sessionStore.NetworkContracts[sessionStore.Login.Name.GetPublicNetwork()].GwContract.SubmitTransaction(
			"CreateAuction", "false", "", r.FormValue("material"),
			r.FormValue("form"), r.FormValue("weight"), r.FormValue("minPrice"), r.FormValue("hours"))
		if e != nil {
			log.Err(e).Msg("Error while submiting auction creation to the hyperledger state")
			http.Redirect(w, r, "/auctions/list", http.StatusSeeOther)
			return
		}
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

	var (
		e          error
		winnerJSON []byte
	)
	if r.FormValue("private") == "true" {
		endorsingPeerOption := gateway.WithEndorsingPeers(sessionStore.Login.Name.GetPublicNetwork().GetEndorsingPeer())
		txn, e := sessionStore.NetworkContracts[sessionStore.Login.Name.GetPublicNetwork()].GwContract.CreateTransaction(
			"FinishAuction", endorsingPeerOption)
		if e != nil {
			log.Err(e).Msg("Error while creating transaction")
			http.Redirect(w, r, redirectPath, http.StatusSeeOther)
			return
		}
		winnerJSON, e = txn.Submit("true", r.FormValue("auctionID"), r.FormValue("colNums"))
		if e != nil {
			log.Err(e).Msg("Error when submiting transaction to finish auction")
			http.Redirect(w, r, redirectPath, http.StatusSeeOther)
			return
		}
	} else {
		winnerJSON, e = sessionStore.NetworkContracts[sessionStore.Login.Name.GetPublicNetwork()].GwContract.SubmitTransaction(
			"FinishAuction", "false", r.FormValue("auctionID"), "")
		if e != nil {
			log.Err(e).Msg("Error while finishing auction in the hyperledger state")
			http.Redirect(w, r, redirectPath, http.StatusSeeOther)
			return
		}
	}

	winner := strings.ToLower(strings.ReplaceAll(string(winnerJSON), "MSP", ""))
	if winner != "" {
		http.Redirect(w, r, "/delivery/create?auctionID="+r.FormValue("auctionID")+"&winner="+winner, http.StatusSeeOther)
	} else {
		http.Redirect(w, r, redirectPath, http.StatusSeeOther)
	}
}
