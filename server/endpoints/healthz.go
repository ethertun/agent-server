package endpoints

import (
	"net/http"
	"time"

	"github.com/go-chi/render"
)

var (
	// Timestamp of when server was started
    startTime time.Time
)

func init() {
    startTime = time.Now()
}

type HealthzResponse struct {
	StartTime time.Time `json:"startTime"`
	UpTime    float64   `json:"upTime"`
}

func NewHealtzResponse() *HealthzResponse {
	return &HealthzResponse{
		StartTime: startTime,
		UpTime:    time.Now().Sub(startTime).Minutes(),
	}
}

func Healthz(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, NewHealtzResponse())
}
