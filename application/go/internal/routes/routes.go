package routes

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"

	fabricConfig "github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/gateway"
	"github.com/reymom/bsm-hyperledger/application/go/cmd/app/config"
	"github.com/rs/zerolog/log"
)

var templates *template.Template
var auctionContract *gateway.Contract
var logisticsContract *gateway.Contract

func GenerateRoutes(conf *config.Config) (http.Handler, error) {
	e := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	if e != nil {
		log.Err(e).Msg("Error setting env variable DISCOVERY_AS_LOCALHOST")
	}

	wallet, e := gateway.NewFileSystemWallet("wallet")
	if e != nil {
		log.Err(e).Msg("Failed to create wallet")
	}

	if !wallet.Exists("appUser") {
		e = populateWallet(wallet)
		if e != nil {
			log.Err(e).Msg("Failed to populate wallet contents")
		}
	}

	ccpPath := filepath.Join(
		"..",
		"..",
		"organizations",
		"peerOrganizations",
		"supplier1.steelplatform.com",
		"connection-supplier1.yaml",
	)

	gw, e := gateway.Connect(
		gateway.WithConfig(fabricConfig.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, "appUser"),
	)
	if e != nil {
		log.Err(e).Msg("Failed to connect to gateway")
	}
	defer gw.Close()

	network, e := gw.GetNetwork("public1channel")
	if e != nil {
		log.Err(e).Msg("Failed to get network")
	}

	auctionContract = network.GetContract("auction")
	log.Info().Msgf("Loaded Contract:%s", auctionContract.Name())
	// logisticsContract = network.GetContract("logistics")

	//template routing
	funcMap := template.FuncMap{
		"dummyFunc": func(str string) string {
			return str + " is dummy"
		},
	}
	t := template.New("appTemplate").Funcs(funcMap)
	templates, e = t.ParseGlob("www/views/pages/*/*.html")
	if e != nil {
		return nil, e
	}

	mux := http.NewServeMux()

	mux.Handle("/home", http.HandlerFunc(homeHandler))
	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./www/public/css"))))
	// mux.Handle("/callback", callBackHandler)
	// mux.Handle("/login", login.LoginHandler)
	// mux.Handle("/logout", logout.LogoutHandler)

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

	mux.Handle("/", http.HandlerFunc(rootHandler))

	return mux, nil
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	template := templates.Lookup("root")
	if template == nil {
		log.Err(err).Msg("Error while looking up \"root\" template")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var Context string
	switch r.Host {
	case "reymom.steelplatform.com":
		Context = "PROD"
	case "localhost:8080":
		Context = "LOC"
	default:
		Context = "NAH"
	}

	vd := ViewData{
		Context: Context,
		Name:    "",
	}
	w.Header().Set("Content-Type", "text/html")
	err = template.Execute(w, vd)
	if err != nil {
		log.Err(err).Msg("Error while executing template")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		os.Exit(1)
	}
}
