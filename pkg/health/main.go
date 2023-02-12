package health

import (
	"fmt"
	"net/http"

	"github.com/michaelraskansky/ntokin/pkg/dts"
)

func healthcheck(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "{\"status\":\"healthy\"}")
}

func Start(ctx *dts.Ctx) {
	host := fmt.Sprintf("0.0.0.0:%v", ctx.HealthcheckPort)
	http.HandleFunc("/healthcheck", healthcheck)
	go func() { http.ListenAndServe(host, nil) }()
}
