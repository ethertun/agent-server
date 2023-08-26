package endpoints

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/ethertun/agent-server/server/errors"
	"github.com/go-chi/render"
)

type ScheduleTaskCallback func(*ScheduleTaskRequest) error

type ScheduleTaskRequest struct {
	Command   string            `json:"command"` // command to execute
	StartTime time.Time         `json:"at"`      // when to execute this command
	Options   map[string]string `json:"options"` // agent-specific options to be passed through
}

func (ctr *ScheduleTaskRequest) Bind(r *http.Request) error {
	if ctr.Command == "" {
		return fmt.Errorf("command is required")
	}

	return nil
}

func ScheduleTask(cb ScheduleTaskCallback) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var request ScheduleTaskRequest
		if err := render.Bind(r, &request); err != nil {
			// failed to decode response, return 400 bad request
			render.Render(w, r, errors.ErrInvalidRequest(err))
			return
		}

		if err := cb(&request); err != nil {
			slog.Error("unable to schedule task", "error", err)
		}

		render.Render(w, r, errors.ErrInvalidRequest(fmt.Errorf("not implemented")))
	}

	return http.HandlerFunc(fn)
}
