package api

import (
	"net/http"

	"github.com/aleakimova/yandexpr-final/pkg/middleware"
)

func Init(pass string) {
	http.HandleFunc("/api/nextdate", middleware.LogRequest(nextDayHandler))
	http.HandleFunc("/api/signin", middleware.LogRequest(signInHandler))
	http.HandleFunc("/api/task", middleware.LogRequest(middleware.Auth(pass, taskHandler)))
	http.HandleFunc("/api/tasks", middleware.LogRequest(middleware.Auth(pass, tasksHandler)))
	http.HandleFunc("/api/task/done", middleware.LogRequest(middleware.Auth(pass, doneTaskHandler)))
}
