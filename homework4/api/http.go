package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	rsrc "github.com/AndreyKlimchuk/golang-learning/homework4/resources"
	"github.com/AndreyKlimchuk/golang-learning/homework4/resources/tasks"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	validator "github.com/go-playground/validator/v10"
)

var validate = validator.New()

func StartHttpServer() {
	_ = http.ListenAndServe(":8080", NewHttpRouter())
}

func NewHttpRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)

	r.Route("/projects", func(r chi.Router) {
		r.Post("/", CreateProject)
		r.Get("/", GetProjects)

		r.Route("/{projectID:[\\d]+}", func(r chi.Router) {
			r.Get("/", GetProject)
			r.Put("/", UpdateProject)
			r.Delete("/", DeleteProject)

			r.Route("/columns", func(r chi.Router) {
				r.Post("/", CreateColumn)
				r.Get("/", GetColumns)

				r.Route("/{columnID:[\\d]+}", func(r chi.Router) {
					r.Get("/", GetColumn)
					r.Put("/", UpdateColumn)
					r.Delete("/", DeleteColumn)

					r.Put("/position", UpdateColumnPosition)
					r.Post("/tasks", CreateTask)
				})
			})
		})
	})

	r.Route("/tasks", func(r chi.Router) {
		r.Route("/{taskID:[\\d]+}", func(r chi.Router) {
			r.Get("/", GetTask)
			r.Put("/", UpdateTask)
			r.Delete("/", DeleteTask)

			r.Put("/position", UpdateTaskPosition)

			r.Route("/comments", func(r chi.Router) {
				r.Post("/", CreateComment)

				r.Route("/{commentId:[\\d]+}", func(r chi.Router) {
					r.Get("/", GetComment)
					r.Put("/", UpdateComment)
					r.Delete("/", DeleteComment)
				})
			})
		})
	})
	return r
}

func CreateTask(w http.ResponseWriter, httpReq *http.Request) {
	handleRequest(w, httpReq, tasks.CreateRequest{
		ProjectId: getId(httpReq, "projectID"),
		ColumnId:  getId(httpReq, "columnID"),
	})
}

func handleRequest(w http.ResponseWriter, httpReq *http.Request, req rsrc.Request) {
	body, err := ioutil.ReadAll(httpReq.Body)
	if err != nil {
		// log
		httpServerError(w)
		return
	}
	defer httpReq.Body.Close()
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if err := validate.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}
	resp, err := req.Handle()
	sendResponse(w, httpReq, resp, err)
}

func sendResponse(w http.ResponseWriter, httpReq *http.Request, body interface{}, err error) {
	if err != nil {
		sendError(w, err)
		return
	}
	switch httpReq.Method {
	case "GET":
		sendJSONResponse(w, http.StatusOK, body)
	case "POST":
		w.Header().Set("Location", getLocation(httpReq, body.(rsrc.Resource)))
		sendJSONResponse(w, http.StatusCreated, body)
	case "PUT", "DELETE":
		w.WriteHeader(http.StatusNoContent)
	}
}

func sendError(w http.ResponseWriter, err error) {
	var genError = rsrc.Error{}
	if yes := errors.As(err, &genError); yes {
		switch genError.Type {
		case rsrc.NotFound:
			http.Error(w, "not found", http.StatusNotFound)
		case rsrc.Conflict:
			http.Error(w, genError.Description, http.StatusConflict)
		case rsrc.InternalError:
			// log
			httpServerError(w)
		}
	} else {
		// log
		httpServerError(w)
	}
}

func getLocation(httpReq *http.Request, resource rsrc.Resource) string {
	var location string
	id := strconv.Itoa(int(resource.GetId()))
	switch resource.(type) {
	// task resource has different base path after creation
	case rsrc.Task:
		location = "/tasks/" + id
	default:
		location = httpReq.URL.Path + "/" + id
	}
	return location
}

func sendJSONResponse(w http.ResponseWriter, statusCode int, body interface{}) {
	binBody, err := json.Marshal(body)
	if err != nil {
		// log
		httpServerError(w)
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(statusCode)
	if _, err := w.Write(binBody); err != nil {
		// log
	}
}

func httpServerError(w http.ResponseWriter) {
	http.Error(w, "server error", http.StatusInternalServerError)
}

func getId(r *http.Request, key string) rsrc.Id {
	id, _ := strconv.Atoi(chi.URLParam(r, key))
	return rsrc.Id(id)
}

func getExpanded(r *http.Request) bool {
	query := r.URL.Query()
	expanded, prs := query["expanded"]
	return prs && (expanded[0] == "" || expanded[0] == "true")
}
