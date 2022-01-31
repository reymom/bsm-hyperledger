package routes

import (
	"net/http"
	"os"
	"strconv"

	"github.com/reymom/bsm-hyperledger/application/go/internal/connection"
	"github.com/reymom/bsm-hyperledger/application/go/internal/sessionstore"
	"github.com/rs/zerolog/log"
	"github.com/urfave/negroni"
)

func generateLoginRoutes(mux *http.ServeMux) error {
	mux.Handle("/login", negroni.New(negroni.Wrap(http.HandlerFunc(loginHandler))))
	mux.Handle("/login/submit", negroni.New(negroni.Wrap(http.HandlerFunc(submitLogin))))
	mux.Handle("/logout", negroni.New(negroni.Wrap(http.HandlerFunc(logoutHandler))))
	return nil
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
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

	if loggedIn {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	template := templates.Lookup("login")
	if template == nil {
		log.Err(e).Msg("Error while looking up \"login\" template")
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}

	var context string
	switch r.Host {
	case "reymom.steelplatform.com":
		context = "PROD"
	case "localhost:8080":
		context = "LOC"
	default:
		context = "NAH"
	}

	falseAttempt, _ := strconv.ParseBool(r.FormValue("falseAttempt"))
	m := struct {
		Context   string
		Attempted bool
	}{
		Context:   context,
		Attempted: falseAttempt,
	}

	if !loggedIn {
		w.Header().Set("Content-Type", "text/html")
		e = template.Execute(w, m)
		if e != nil {
			log.Err(e).Msg("Error while executing template")
			http.Error(w, e.Error(), http.StatusInternalServerError)
			os.Exit(1)
		}
	}

}

func submitLogin(w http.ResponseWriter, r *http.Request) {
	var e error

	session, e := sessionstore.Store.Get(r, "coockie-name")
	if e != nil {
		log.Err(e).Msg("Error while getting session storage")
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}

	user, password := connection.Organization(r.FormValue("user")), r.FormValue("password")
	if connection.IsRegistered(connectionConfig.UsersLoginMap, user, password) {
		authUser = &connection.Login{
			Name:     user,
			Password: password,
		}
		loggedIn = true
		log.Info().Msgf("%s logged in", user)

		session.Values["login"] = authUser
		e = session.Save(r, w)
		if e != nil {
			log.Err(e).Msg("Error while saving session login")
			http.Error(w, e.Error(), http.StatusInternalServerError)
			return
		}

		_, networkContracts, e = connection.GetGatewayObjects(user)
		if e != nil {
			log.Err(e).Msg("Error getting gateway objects")
		}
	} else {
		http.Redirect(w, r, "/login?falseAttempt=true", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, e := sessionstore.Store.Get(r, "coockie-name")
	if e != nil {
		log.Err(e).Msg("Error while getting session storage")
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}

	session.Values["login"] = connection.Login{}
	session.Options.MaxAge = -1

	e = session.Save(r, w)
	if e != nil {
		log.Err(e).Msg("Error while saving session logout")
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}

	authUser, loggedIn = new(connection.Login), false
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
