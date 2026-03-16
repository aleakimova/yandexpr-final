package api

import "net/http"

func Init() {
	http.HandleFunc("/api/nextdate", logRequest(nextDayHandler))
	http.HandleFunc("/api/signin", logRequest(signInHandler))
	http.HandleFunc("/api/task", logRequest(auth(taskHandler)))
	http.HandleFunc("/api/tasks", logRequest(auth(tasksHandler)))
	http.HandleFunc("/api/task/done", logRequest(auth(doneTaskHandler)))
}
