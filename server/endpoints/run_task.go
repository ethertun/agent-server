package endpoints

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/ethertun/agent-server/server/errors"
	"github.com/go-chi/render"
	"github.com/mitchellh/mapstructure"
)

type RunTaskRequest struct {
	Command   string         `json:"command"` // command to execute
	StartTime time.Time      `json:"at"`      // when to execute this command
	Options   map[string]any `json:"options"` // agent-specific options to be passed through
}

type RunTaskResponse struct {
	QueueTime time.Time `json:"queued"` // the time this task was queued to run
}

type RunTaskOutput struct {
	Stdout string
	Stderr string
}

type RunTaskCallback func(*RunTaskRequest) (*RunTaskOutput, error)

func (ctr *RunTaskRequest) Bind(r *http.Request) error {
	if ctr.Command == "" {
		return fmt.Errorf("a command is required")
	}

	return nil
}

func (ctr *RunTaskRequest) ParseOptions(v any) error {
    cfg := &mapstructure.DecoderConfig{
        TagName: "json",
        Result: v,
    }

    decoder, err := mapstructure.NewDecoder(cfg)
    if err != nil {
        return fmt.Errorf("unable to build decoder for request options: %w", err)
    }

    err = decoder.Decode(ctr.Options)
    if err != nil {
        return fmt.Errorf("unable to parse task request options: %w", err)
    }

    return nil
}

func RunTask(cb RunTaskCallback) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var request RunTaskRequest
		if err := render.Bind(r, &request); err != nil {
			// failed to decode response, return 400 bad request
			render.Render(w, r, errors.ErrInvalidRequest(err))
			return
		}

		// trim whitespace
		request.Command = strings.TrimSpace(request.Command)

		// spawn a goroutine to run this task
		go exec(cb, &request)

		resp := RunTaskResponse{
			QueueTime: time.Now(),
		}

		render.JSON(w, r, resp)
	}

	return http.HandlerFunc(fn)
}

func exec(cb RunTaskCallback, r *RunTaskRequest) {
	output, err := cb(r)
	if err != nil {
		slog.Error("unable to run task", "error", err)
	} else {
		slog.Debug("task output", "stdout", output.Stdout, "stderr", output.Stderr)
	}

	// TODO post response to server
}
