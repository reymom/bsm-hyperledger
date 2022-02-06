package routes

import (
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/reymom/bsm-hyperledger/application/go/cmd/app/config"
	"github.com/reymom/bsm-hyperledger/application/go/internal/connection"
	"github.com/reymom/bsm-hyperledger/application/go/internal/sessionstore"
	"github.com/rs/zerolog/log"
)

var (
	templates        *template.Template
	connectionConfig *config.Config
	sessionStore     = new(sessionstore.SessionStore)
	loggedIn         bool
)

func GenerateRoutes(conf *config.Config) (http.Handler, error) {
	e := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	if e != nil {
		log.Err(e).Msg("Error setting env variable DISCOVERY_AS_LOCALHOST")
	}

	connectionConfig = conf

	sessionStore, e = sessionstore.NewSessionStore(conf.SessionKey)
	if e != nil {
		log.Err(e).Msg("Error creating session store")
		return nil, e
	}

	//template routing
	funcMap := template.FuncMap{
		"loggedIn": func() bool {
			return loggedIn
		},
		"getOrganizationName": func() string {
			return string(sessionStore.Login.Name)
		},
		"orgType": func() string {
			switch {
			case strings.Contains(string(sessionStore.Login.Name), "supplier"):
				return "supplier"
			case strings.Contains(string(sessionStore.Login.Name), "buyer"):
				return "buyer"
			case strings.Contains(string(sessionStore.Login.Name), "logistics"):
				return "logistics"
			default:
				return ""
			}
		},
		"stringToNetwork": func(str string) connection.Channel {
			return connection.Channel(str)
		},
	}
	t := template.New("steelPlatform.gohtml").Funcs(funcMap)
	templates, e = t.ParseGlob("www/templates/views/*/*.html")
	if e != nil {
		return nil, e
	}

	mux := http.NewServeMux()

	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./www/static/css"))))
	mux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./www/static/images"))))
	mux.Handle("/home", http.HandlerFunc(homeHandler))

	e = generateLoginRoutes(mux)
	if e != nil {
		return nil, e
	}
	e = generateNetworkRoutes(mux)
	if e != nil {
		return nil, e
	}
	e = generateAuctionRoutes(mux)
	if e != nil {
		return nil, e
	}
	e = generateBidRoutes(mux)
	if e != nil {
		return nil, e
	}
	e = generateLogisticsRoutes(mux)
	if e != nil {
		return nil, e
	}

	mux.Handle("/", http.HandlerFunc(homeHandler))

	return mux, nil
}
