package prometheus

import (
	"fmt"
	"github.com/mahmoudKheyrati/marketplace-backend/pkg"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

type Prometheus struct {
	port int
	mux  *http.ServeMux
}

func NewPrometheus(port int) *Prometheus {
	mux := http.NewServeMux()
	return &Prometheus{
		port: port,
		mux:  mux,
	}
}

func (p *Prometheus) RunHTTPServer() {
	p.mux.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(fmt.Sprintf(":%d", p.port), p.mux)
	if err != nil {
		pkg.Logger().Warn("Failed to listen prometheus server: ", err)
	}
}
