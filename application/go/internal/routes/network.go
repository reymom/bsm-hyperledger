package routes

import (
	"net/http"
	"os"

	"github.com/reymom/bsm-hyperledger/application/go/internal/connection"
	"github.com/rs/zerolog/log"
	"github.com/urfave/negroni"
)

func generateNetworkRoutes(mux *http.ServeMux) error {
	mux.Handle("/networkInfo", negroni.New(negroni.Wrap(http.HandlerFunc(networkInfoHandler))))

	return nil
}

func networkInfoHandler(w http.ResponseWriter, r *http.Request) {
	var e error

	if !loggedIn {
		if loggedIn, e = sessionStore.CheckLoginFromSession(r, connectionConfig.UsersLoginMap); !loggedIn {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
	}

	template := templates.Lookup("networkInfo")
	if template == nil {
		log.Err(e).Msg("Error while looking up \"networkInfo\" template")
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}

	m := struct {
		NetworkContracts connection.NetworkContract
	}{
		NetworkContracts: sessionStore.NetworkContracts,
	}

	w.Header().Set("Content-Type", "text/html")
	e = template.Execute(w, m)
	if e != nil {
		log.Err(e).Msg("Error while executing template")
		http.Error(w, e.Error(), http.StatusInternalServerError)
		os.Exit(1)
	}
}
