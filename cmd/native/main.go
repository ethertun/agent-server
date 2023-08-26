package main

import (
	"log/slog"

	"github.com/ethertun/agent-server/server"
	"github.com/ethertun/agent-server/server/endpoints"
)

func schedule_task(r *endpoints.ScheduleTaskRequest) error {
	slog.Info("creating task", "request", r)

	return nil
}

func main() {
	callbacks := server.Callbacks{
		ScheduleTask: schedule_task,
	}

	server := server.NewServer("token", callbacks)
	server.Start(1337)
}
