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
		log.Info().Msg("appUser exists")
	}

	ccpPath := filepath.Join(
		"..",
		"..",
		"organizations",
		"peerOrganizations",
		"logistics.example.com",
		"connection-logistics.yaml",
	)

	gw, e := gateway.Connect(
		gateway.WithConfig(fabricConfig.FromFile(filepath.Clean(ccpPath))),
		gateway.WithIdentity(wallet, "appUser"),
	)
	if e != nil {
		log.Err(e).Msg("Failed to connect to gateway")
	}
	defer gw.Close()

	network, e := gw.GetNetwork("cursochannel")
	if e != nil {
		log.Err(e).Msg("Failed to get network")
	}

	auctionContract = network.GetContract("auctionContract")
	logisticsContract = network.GetContract("logisticsContract")

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
	template := templates.Lookup("root")
	println(template.Name())

	host := r.Host
	println("host = ", host)
}
