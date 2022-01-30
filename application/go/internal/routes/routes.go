package routes

import (
	"html/template"
	"net/http"
	"os"

	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"github.com/reymom/bsm-hyperledger/application/go/cmd/app/config"
	"github.com/reymom/bsm-hyperledger/application/go/internal/connection"
	"github.com/rs/zerolog/log"
)

var (
	templates        *template.Template
	connectionConfig *config.Config

	gw          *gateway.Gateway
	gwContracts []connection.GatewayContract

	authUser = new(connection.Login)
	loggedIn bool
)

func GenerateRoutes(conf *config.Config) (http.Handler, error) {
	e := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	if e != nil {
		log.Err(e).Msg("Error setting env variable DISCOVERY_AS_LOCALHOST")
	}

	connectionConfig = conf

	//template routing
	funcMap := template.FuncMap{
		"loggedIn": func() bool {
			return loggedIn
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
	e = generateAuctionRoutes(mux)
	if e != nil {
		return nil, e
	}
	e = generateLogisticsRoutes(mux)
	if e != nil {
		return nil, e
	}
	e = generateSupplierRoutes(mux)
	if e != nil {
		return nil, e
	}
	e = generateBuyerRoutes(mux)
	if e != nil {
		return nil, e
	}

	mux.Handle("/", http.HandlerFunc(homeHandler))

	return mux, nil
}
