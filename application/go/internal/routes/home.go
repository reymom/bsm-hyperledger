package routes

import (
	"net/http"
	"os"

	"github.com/rs/zerolog/log"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	var e error

	template := templates.Lookup("home")
	if template == nil {
		log.Err(e).Msg("Error while looking up \"home\" template")
		http.Error(w, e.Error(), http.StatusInternalServerError)
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
	e = template.Execute(w, vd)
	if e != nil {
		log.Err(e).Msg("Error while executing template")
		http.Error(w, e.Error(), http.StatusInternalServerError)
		os.Exit(1)
	}
}
