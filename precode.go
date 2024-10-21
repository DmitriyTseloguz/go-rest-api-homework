package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

func getTasks(response http.ResponseWriter, request *http.Request) {
	fmt.Println("get tasks")
	var jsonTasks, marshalError = json.Marshal(tasks)

	if marshalError != nil {
		http.Error(response, marshalError.Error(), http.StatusInternalServerError)
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	response.Write(jsonTasks)
}

func createTasks(response http.ResponseWriter, request *http.Request) {
	var task Task
	var buffer bytes.Buffer

	var _, readError = buffer.ReadFrom(request.Body)

	if readError != nil {
		fmt.Println(readError)
		http.Error(response, readError.Error(), http.StatusBadRequest)
	}

	var unmarshalError = json.Unmarshal(buffer.Bytes(), &task)

	if unmarshalError != nil {
		fmt.Println(unmarshalError)
		http.Error(response, unmarshalError.Error(), http.StatusBadRequest)
	}

	tasks[task.ID] = task

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
}

func main() {
	router := chi.NewRouter()

	router.Get("/tasks", getTasks)

	router.Post("/tasks", createTasks)

	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
