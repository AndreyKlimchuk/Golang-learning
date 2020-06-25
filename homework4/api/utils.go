package api

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/AndreyKlimchuk/golang-learning/homework4/logger"
	rsrc "github.com/AndreyKlimchuk/golang-learning/homework4/resources"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

var validate = validator.New()

func StartHttpServer() {
	err := http.ListenAndServe(":8080", NewRouter())
	logger.Zap.Fatal("http server termination", zap.Error(err))
}

func handleRequest(w http.ResponseWriter, httpReq *http.Request, req interface{}) {
	body, err := ioutil.ReadAll(httpReq.Body)
	if err != nil {
		logger.Zap.Error("error while reading request body", zap.Error(err))
		httpServerError(w)
		return
	}
	defer httpReq.Body.Close()
	if len(body) > 0 {
		if err := json.Unmarshal(body, req); err != nil {
			http.Error(w, "Invalid json", http.StatusBadRequest)
			return
		}
	}
	if err := validate.Struct(req); err != nil {
		http.Error(w, formatValidationErrors(err), http.StatusUnprocessableEntity)
		return
	}
	resp, err := req.(rsrc.Request).Handle()
	sendResponse(w, httpReq, resp, err)
}

func formatValidationErrors(err error) string {
	// TODO: customize errors
	return err.(validator.ValidationErrors).Error()
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
			logger.Zap.Error("internal error", zap.Error(err))
			httpServerError(w)
		}
	} else {
		logger.Zap.Error("unhandled internal error", zap.Error(err))
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
		logger.Zap.Error("error while marshaling response body", zap.Error(err))
		httpServerError(w)
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(statusCode)
	if _, err := w.Write(binBody); err != nil {
		logger.Zap.Error("error while writing http response", zap.Error(err))
	}
}

func httpServerError(w http.ResponseWriter) {
	http.Error(w, "server error", http.StatusInternalServerError)
}

func getId(r *http.Request, key string) rsrc.Id {
	id, _ := strconv.Atoi(chi.URLParam(r, key))
	return rsrc.Id(id)
}

func getProjectId(r *http.Request) rsrc.Id { return getId(r, "projectID") }

func getColumnId(r *http.Request) rsrc.Id { return getId(r, "columnID") }

func getTaskId(r *http.Request) rsrc.Id { return getId(r, "taskID") }

func getCommentId(r *http.Request) rsrc.Id { return getId(r, "commentID") }

func getExpanded(r *http.Request) bool {
	query := r.URL.Query()
	expanded, prs := query["expanded"]
	return prs && (expanded[0] == "" || expanded[0] == "true")
}
