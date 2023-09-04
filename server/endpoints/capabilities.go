package endpoints

import (
	"net/http"
	"time"

	"github.com/go-chi/render"
)

type CapabilitiesResponse struct {
	StartTime time.Time `json:"startTime"`
	UpTime    float64   `json:"upTime"`
}

func Capabilities(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, NewHealtzResponse())
}
