package routes

import (
	"net/http"
	"os"

	"github.com/reymom/bsm-hyperledger/application/go/internal/connection"
	"github.com/reymom/bsm-hyperledger/application/go/internal/sessionstore"
	"github.com/rs/zerolog/log"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	var e error

	session, e := sessionstore.Store.Get(r, "coockie-name")
	if e != nil {
		log.Err(e).Msg("Error while getting session storage")
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}
	authUser = sessionstore.GetLoginFromSession(session)
	if (connection.Login{}) != *authUser {
		loggedIn = true
	}

	if !loggedIn {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	template := templates.Lookup("home")
	if template == nil {
		log.Err(e).Msg("Error while looking up \"home\" template")
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	e = template.Execute(w, new(interface{}))
	if e != nil {
		log.Err(e).Msg("Error while executing template")
		http.Error(w, e.Error(), http.StatusInternalServerError)
		os.Exit(1)
	}
}
