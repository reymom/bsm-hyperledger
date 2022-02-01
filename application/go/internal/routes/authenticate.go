package routes

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/reymom/bsm-hyperledger/application/go/internal/connection"
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

	if !loggedIn {
		loggedIn, e = sessionStore.CheckLoginFromSession(r, connectionConfig.UsersLoginMap)
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

	w.Header().Set("Content-Type", "text/html")
	e = template.Execute(w, m)
	if e != nil {
		log.Err(e).Msg("Error while executing template")
		http.Error(w, e.Error(), http.StatusInternalServerError)
		os.Exit(1)
	}

}

func submitLogin(w http.ResponseWriter, r *http.Request) {
	var e error

	session, e := sessionStore.Store.Get(r, "auth-session")
	if e != nil {
		log.Err(e).Msg("Error while getting session storage in submit loggin")
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}

	user, password := connection.Organization(r.FormValue("user")), r.FormValue("password")
	if connection.IsRegistered(connectionConfig.UsersLoginMap, user, password) {
		sessionStore.Login = &connection.Login{
			Name:     user,
			Password: password,
		}
		loggedIn = true
		log.Info().Msgf("%s logged in", user)

		session.Values["login"] = &sessionStore.Login
		e = session.Save(r, w)
		if e != nil {
			log.Err(e).Msg("Error while saving session login")
			http.Error(w, e.Error(), http.StatusInternalServerError)
			return
		}

		_, sessionStore.NetworkContracts, e = connection.GetGatewayObjects(user)
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
	session, e := sessionStore.Store.Get(r, "auth-session")
	if e != nil {
		log.Err(e).Msg("Error while getting session storage on logout")
	}

	session.Values["login"] = connection.Login{}
	session.Options.MaxAge = -1
	e = sessionStore.Store.Save(r, w, session)
	if e != nil {
		http.Error(w, e.Error(), 400)
		return
	}

	cookie := http.Cookie{
		Name:       "auth-session",
		Value:      "",
		Path:       "",
		Domain:     "",
		Expires:    time.Time{},
		RawExpires: "",
		MaxAge:     -1,
		Secure:     false,
		HttpOnly:   false,
		SameSite:   0,
		Raw:        "",
		Unparsed:   nil,
	}
	http.SetCookie(w, &cookie)

	sessionStore.Login, loggedIn = new(connection.Login), false
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
