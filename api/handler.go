package api

import (
	"log"
	"net/http"
	"os"

	"github.com/suconghou/videoproxy/route"
)

var (
	logger = log.New(os.Stdout, "", 0)
)

// Handler hand
func Handler(w http.ResponseWriter, r *http.Request) {
	routeMatch(w, r)
}

func routeMatch(w http.ResponseWriter, r *http.Request) {
	for _, p := range route.Route {
		if p.Reg.MatchString(r.URL.Path) {
			if err := p.Handler(w, r, p.Reg.FindStringSubmatch(r.URL.Path)); err != nil {
				logger.Print(err)
			}
			return
		}
	}
	fallback(w, r)
}

func fallback(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}
