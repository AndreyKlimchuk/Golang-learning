package api

import (
	"bytes"
	"encoding/json"
	"github.com/AndreyKlimchuk/golang-learning/homework4/db"
	"github.com/AndreyKlimchuk/golang-learning/homework4/logger"
	"github.com/AndreyKlimchuk/golang-learning/homework4/resources/columns"
	"github.com/AndreyKlimchuk/golang-learning/homework4/resources/common"
	"github.com/AndreyKlimchuk/golang-learning/homework4/resources/tasks"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

var client *http.Client
var URL string

const nonExistentId = 9999

func TestMain(m *testing.M) {
	if err := logger.InitZap(); err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	if err := db.ApplyMigrationsDown(); err != nil {
		log.Fatalf("can't apply down migrations: %v", err)
	}
	if err := db.Init(); err != nil {
		log.Fatalf("can't initialize db: %v", err)
	}
	srv := httptest.NewServer(NewRouter())
	defer srv.Close()
	client = srv.Client()
	URL = srv.URL
	os.Exit(m.Run())
}

var project1 = common.ProjectExpanded{
	Project: common.Project{Id: 1, ProjectSettableFields: common.ProjectSettableFields{Name: "c", Description: "desc"}},
	Columns: []common.ColumnExpanded{},
}

var project2 = common.ProjectExpanded{
	Project: common.Project{Id: 2, ProjectSettableFields: common.ProjectSettableFields{Name: "a", Description: "desc"}},
	Columns: []common.ColumnExpanded{},
}

var project3 = common.ProjectExpanded{
	Project: common.Project{Id: 3, ProjectSettableFields: common.ProjectSettableFields{Name: "b", Description: "desc"}},
	Columns: []common.ColumnExpanded{},
}

var column1P1Def = common.ColumnExpanded{
	Column: common.Column{Id: 1, ColumnSettableFields: common.ColumnSettableFields{Name: common.DefaultColumnName}},
	Tasks:  []common.Task{},
}

var column2P2Def = common.ColumnExpanded{
	Column: common.Column{Id: 2, ColumnSettableFields: common.ColumnSettableFields{Name: common.DefaultColumnName}},
	Tasks:  []common.Task{},
}

var column3P3Def = common.ColumnExpanded{
	Column: common.Column{Id: 3, ColumnSettableFields: common.ColumnSettableFields{Name: common.DefaultColumnName}},
	Tasks:  []common.Task{},
}

var column4P3 = common.ColumnExpanded{
	Column: common.Column{Id: 4, ColumnSettableFields: common.ColumnSettableFields{Name: "a"}},
	Tasks:  []common.Task{},
}

var column5P3 = common.ColumnExpanded{
	Column: common.Column{Id: 5, ColumnSettableFields: common.ColumnSettableFields{Name: "b"}},
	Tasks:  []common.Task{},
}

var task1 = common.TaskExpanded{
	Task: common.Task{
		ProjectId:          project3.Id,
		ColumnId:           column4P3.Id,
		Id:                 1,
		TaskSettableFields: common.TaskSettableFields{Name: "a", Description: "desc"},
	},
	Comments: []common.Comment{},
}

var task2 = common.TaskExpanded{
	Task: common.Task{
		ProjectId:          project3.Id,
		ColumnId:           column5P3.Id,
		Id:                 2,
		TaskSettableFields: common.TaskSettableFields{Name: "b", Description: "desc"},
	},
	Comments: []common.Comment{},
}

var task3 = common.TaskExpanded{
	Task: common.Task{
		ProjectId:          project3.Id,
		ColumnId:           column5P3.Id,
		Id:                 3,
		TaskSettableFields: common.TaskSettableFields{Name: "c", Description: "desc"},
	},
	Comments: []common.Comment{},
}

var comment1T3 = common.Comment{Id: 1, CommentSettableFields: common.CommentSettableFields{Text: "text"}}

var comment2T3 = common.Comment{Id: 2, CommentSettableFields: common.CommentSettableFields{Text: "text"}}

func Test_Complex(t *testing.T) {
	assertPost201(t, projectsPath(), project1.ProjectSettableFields, project1.Project)
	assertPost201(t, projectsPath(), project2.ProjectSettableFields, project2.Project)
	assertPost201(t, projectsPath(), project3.ProjectSettableFields, project3.Project)

	assertPost201(t, columnsPath(project3.Id), column4P3.ColumnSettableFields, column4P3.Column)
	assertPost201(t, columnsPath(project3.Id), column5P3.ColumnSettableFields, column5P3.Column)

	assertPost201(t, tasksPath(project3.Id, column4P3.Id), task1.TaskSettableFields, task1.Task)
	assertPost201(t, tasksPath(project3.Id, column5P3.Id), task2.TaskSettableFields, task2.Task)
	assertPost201(t, tasksPath(project3.Id, column5P3.Id), task3.TaskSettableFields, task3.Task)

	assertPost201(t, commentsPath(task3.Id), comment1T3.CommentSettableFields, comment1T3)
	assertPost201(t, commentsPath(task3.Id), comment2T3.CommentSettableFields, comment2T3)

	runSubtestsCreate(t)
	runSubtestsGet(t)
	runSubtestsUpdate(t)
	runSubtestsUpdateColumnPosition(t)
	runSubtestsUpdateTaskPosition(t)
	runSubtestsDelete(t)
}

func runSubtestsCreate(t *testing.T) {
	t.Run("cannot create column with duplicate name", func(t *testing.T) {
		assertPost409(t, columnsPath(project3.Id), column4P3.ColumnSettableFields)
	})
}

func runSubtestsGet(t *testing.T) {
	t.Run("get projects ordered by name", func(t *testing.T) {
		projects := []common.Project{project2.Project, project3.Project, project1.Project}
		assertGet200(t, projectsPath(), projects)
	})
	t.Run("get expanded project", func(t *testing.T) {
		c1 := column3P3Def
		c2 := column4P3
		c2.Tasks = []common.Task{task1.Task}
		c3 := column5P3
		c3.Tasks = []common.Task{task2.Task, task3.Task}
		p := project3
		p.Columns = []common.ColumnExpanded{c1, c2, c3}
		assertGet200(t, projectPath(p.Id)+"?expanded", p)
	})
	t.Run("get non existent expanded project", func(t *testing.T) {
		assertGet404(t, projectPath(nonExistentId)+"?expanded")
	})
	t.Run("get columns", func(t *testing.T) {
		columns := []common.Column{column3P3Def.Column, column4P3.Column, column5P3.Column}
		assertGet200(t, columnsPath(project3.Id), columns)
	})
	t.Run("get expanded task", func(t *testing.T) {
		task := task3
		task.Comments = []common.Comment{comment1T3, comment2T3}
		assertGet200(t, taskPath(task.Id)+"?expanded", task)
	})
	t.Run("get non existent expanded task", func(t *testing.T) {
		assertGet404(t, taskPath(nonExistentId)+"?expanded")
	})
	t.Run("get comments", func(t *testing.T) {
		comments := []common.Comment{comment1T3, comment2T3}
		assertGet200(t, commentsPath(task3.Id), comments)
	})
}

func runSubtestsUpdate(t *testing.T) {
	t.Run("update project", func(t *testing.T) {
		project3.Name = "c1"
		project3.Description = "desc1"
		assertPut204(t, projectPath(project3.Id), project3.ProjectSettableFields)
		assertGet200(t, projectPath(project3.Id), project3.Project)
	})
	t.Run("cannot set duplicate name to column", func(t *testing.T) {
		c := column4P3
		c.Name = "b"
		assertPut409(t, columnPath(project3.Id, c.Id), c.ColumnSettableFields)
	})
	t.Run("update column", func(t *testing.T) {
		column4P3.Name = "a1"
		assertPut204(t, columnPath(project3.Id, column4P3.Id), column4P3.ColumnSettableFields)
		assertGet200(t, columnPath(project3.Id, column4P3.Id), column4P3.Column)
	})
	t.Run("update task", func(t *testing.T) {
		task3.Name = "c1"
		task3.Description = "desc1"
		assertPut204(t, taskPath(task3.Id), task3.TaskSettableFields)
		assertGet200(t, taskPath(task3.Id), task3.Task)
	})
	t.Run("update comment", func(t *testing.T) {
		comment1T3.Text = "text1"
		assertPut204(t, commentPath(task3.Id, comment1T3.Id), comment1T3.CommentSettableFields)
		assertGet200(t, commentPath(task3.Id, comment1T3.Id), comment1T3)
	})
}

func runSubtestsUpdateColumnPosition(t *testing.T) {
	t.Run("cannot update position of non existent column", func(t *testing.T) {
		body := columns.UpdatePositionRequestBody{AfterColumnId: column3P3Def.Id}
		assertPut404(t, columnPositionPath(project3.Id, nonExistentId), body)
	})
	t.Run("cannot place column after non existent column", func(t *testing.T) {
		body := columns.UpdatePositionRequestBody{AfterColumnId: nonExistentId}
		assertPut409(t, columnPositionPath(project3.Id, column3P3Def.Id), body)
	})
	t.Run("cannot place column after itself", func(t *testing.T) {
		body := columns.UpdatePositionRequestBody{AfterColumnId: column3P3Def.Id}
		assertPut422(t, columnPositionPath(project3.Id, column3P3Def.Id), body)
	})
	t.Run("cannot place column after column from another project", func(t *testing.T) {
		body := columns.UpdatePositionRequestBody{AfterColumnId: column1P1Def.Id}
		assertPut409(t, columnPositionPath(project3.Id, column3P3Def.Id), body)
	})
	t.Run("update column position, place in the middle", func(t *testing.T) {
		body := columns.UpdatePositionRequestBody{AfterColumnId: column3P3Def.Id}
		assertPut204(t, columnPositionPath(project3.Id, column5P3.Id), body)
		cls := []common.Column{column3P3Def.Column, column5P3.Column, column4P3.Column}
		assertGet200(t, columnsPath(project3.Id), cls)
	})
	t.Run("update column position, place at the beginning", func(t *testing.T) {
		body := columns.UpdatePositionRequestBody{AfterColumnId: 0}
		assertPut204(t, columnPositionPath(project3.Id, column4P3.Id), body)
		cls := []common.Column{column4P3.Column, column3P3Def.Column, column5P3.Column}
		assertGet200(t, columnsPath(project3.Id), cls)
	})
}

func runSubtestsUpdateTaskPosition(t *testing.T) {
	t.Run("cannot update position of non existent task", func(t *testing.T) {
		body := tasks.UpdatePositionRequestBody{NewColumnId: column4P3.Id, AfterTaskId: task1.Id}
		assertPut404(t, taskPositionPath(nonExistentId), body)
	})
	t.Run("cannot place task after non existent task", func(t *testing.T) {
		body := tasks.UpdatePositionRequestBody{NewColumnId: column5P3.Id, AfterTaskId: nonExistentId}
		assertPut409(t, taskPositionPath(task1.Id), body)
	})
	t.Run("cannot place task in non existent column", func(t *testing.T) {
		body := tasks.UpdatePositionRequestBody{NewColumnId: nonExistentId, AfterTaskId: task2.Id}
		assertPut409(t, taskPositionPath(task1.Id), body)
	})
	t.Run("cannot place task in column from another project", func(t *testing.T) {
		body := tasks.UpdatePositionRequestBody{NewColumnId: column1P1Def.Id, AfterTaskId: 0}
		assertPut409(t, taskPositionPath(task1.Id), body)
	})
	t.Run("cannot place task after itself", func(t *testing.T) {
		body := tasks.UpdatePositionRequestBody{NewColumnId: column4P3.Id, AfterTaskId: task1.Id}
		assertPut422(t, taskPositionPath(task1.Id), body)
	})
	t.Run("update task position, place at the beginning of new column", func(t *testing.T) {
		task3.ColumnId = column3P3Def.Id
		column4P3.Tasks = []common.Task{task1.Task}
		column3P3Def.Tasks = []common.Task{task3.Task}
		column5P3.Tasks = []common.Task{task2.Task}
		project3.Columns = []common.ColumnExpanded{column4P3, column3P3Def, column5P3}
		body := tasks.UpdatePositionRequestBody{NewColumnId: column3P3Def.Id, AfterTaskId: 0}
		assertPut204(t, taskPositionPath(task3.Id), body)
		assertGet200(t, projectPath(project3.Id)+"?expanded", project3)
	})
	t.Run("update task position, place at the end of new column", func(t *testing.T) {
		task1.ColumnId = column5P3.Id
		column4P3.Tasks = []common.Task{}
		column3P3Def.Tasks = []common.Task{task3.Task}
		column5P3.Tasks = []common.Task{task2.Task, task1.Task}
		project3.Columns = []common.ColumnExpanded{column4P3, column3P3Def, column5P3}
		body := tasks.UpdatePositionRequestBody{NewColumnId: column5P3.Id, AfterTaskId: task2.Id}
		assertPut204(t, taskPositionPath(task1.Id), body)
		assertGet200(t, projectPath(project3.Id)+"?expanded", project3)
	})
}

func runSubtestsDelete(t *testing.T) {
	t.Run("delete comment", func(t *testing.T) {
		assertDelete204(t, commentPath(task3.Id, comment1T3.Id))
	})
	t.Run("cannot delete last column", func(t *testing.T) {
		path := columnPath(project1.Id, column1P1Def.Id)
		assertDelete409(t, path)
		assertGet200(t, path, column1P1Def.Column)
	})
	t.Run("delete second column", func(t *testing.T) {
		task3.ColumnId = column4P3.Id
		column4P3.Tasks = []common.Task{task3.Task}
		column5P3.Tasks = []common.Task{task2.Task, task1.Task}
		project3.Columns = []common.ColumnExpanded{column4P3, column5P3}
		path := columnPath(project3.Id, column3P3Def.Id)
		assertDelete204(t, path)
		assertGet200(t, projectPath(project3.Id)+"?expanded", project3)
	})
	t.Run("delete first column", func(t *testing.T) {
		task3.ColumnId = column5P3.Id
		column5P3.Tasks = []common.Task{task2.Task, task1.Task, task3.Task}
		project3.Columns = []common.ColumnExpanded{column5P3}
		path := columnPath(project3.Id, column4P3.Id)
		assertDelete204(t, path)
		assertGet200(t, projectPath(project3.Id)+"?expanded", project3)
	})
	t.Run("delete task", func(t *testing.T) {
		assertDelete204(t, taskPath(task3.Id))
		assertGet404(t, taskPath(task3.Id))
	})
	t.Run("delete project", func(t *testing.T) {
		assertDelete204(t, projectPath(project3.Id))
		assertGet200(t, projectsPath(), []common.Project{project2.Project, project1.Project})
	})
}

func assertGet200(t *testing.T, path string, wantBody interface{}) {
	resp := sendGetRequest(t, path)
	defer resp.Body.Close()
	assertEqualStatusCode(t, resp, http.StatusOK)
	assertEqualBody(t, resp, wantBody)
}

func assertGet404(t *testing.T, path string) {
	resp := sendGetRequest(t, path)
	defer resp.Body.Close()
	assertEqualStatusCode(t, resp, http.StatusNotFound)
}

func assertPost201(t *testing.T, path string, reqBody interface{}, wantResource interface{}) {
	resp := sendPostRequest(t, path, reqBody)
	defer resp.Body.Close()
	assertEqualStatusCode(t, resp, http.StatusCreated)
	locations, prs := resp.Header["Location"]
	if !prs {
		t.Fatalf("201 response doesn't contain 'Location' header")
	}
	location := strings.TrimPrefix(locations[0], basePath)
	assertGet200(t, location, wantResource)
}

func assertPost409(t *testing.T, path string, reqBody interface{}) {
	resp := sendPostRequest(t, path, reqBody)
	defer resp.Body.Close()
	assertEqualStatusCode(t, resp, http.StatusConflict)
}

func assertPut204(t *testing.T, path string, body interface{}) {
	resp := sendPutRequest(t, path, body)
	assertEqualStatusCode(t, resp, http.StatusNoContent)
}

func assertPut409(t *testing.T, path string, body interface{}) {
	resp := sendPutRequest(t, path, body)
	assertEqualStatusCode(t, resp, http.StatusConflict)
}

func assertPut422(t *testing.T, path string, body interface{}) {
	resp := sendPutRequest(t, path, body)
	assertEqualStatusCode(t, resp, http.StatusUnprocessableEntity)
}

func assertPut404(t *testing.T, path string, body interface{}) {
	resp := sendPutRequest(t, path, body)
	assertEqualStatusCode(t, resp, http.StatusNotFound)
}

func assertDelete204(t *testing.T, path string) {
	resp := sendDeleteRequest(t, path)
	assertEqualStatusCode(t, resp, http.StatusNoContent)
	assertGet404(t, path)
}

func assertDelete409(t *testing.T, path string) {
	resp := sendDeleteRequest(t, path)
	assertEqualStatusCode(t, resp, http.StatusConflict)
}

func sendGetRequest(t *testing.T, path string) *http.Response {
	return sendRequest(t, "GET", path, nil)
}

func sendPostRequest(t *testing.T, path string, body interface{}) *http.Response {
	return sendRequest(t, "POST", path, body)
}

func sendPutRequest(t *testing.T, path string, body interface{}) *http.Response {
	return sendRequest(t, "PUT", path, body)
}

func sendDeleteRequest(t *testing.T, path string) *http.Response {
	return sendRequest(t, "DELETE", path, nil)
}

func sendRequest(t *testing.T, method, path string, body interface{}) *http.Response {
	var reqBody io.Reader = nil
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("error while marshaling request body: %v", err)
		}
		reqBody = bytes.NewBuffer(bodyBytes)
	}
	req, err := http.NewRequest(method, URL+basePath+path, reqBody)
	if err != nil {
		t.Fatalf("new request failed: %v", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("error while sending request: %v", err)
	}
	return resp
}

func projectsPath() string {
	return "/projects"
}

func projectPath(projectId common.Id) string {
	return "/projects/" + idToStr(projectId)
}

func columnsPath(projectId common.Id) string {
	return "/projects/" + idToStr(projectId) + "/columns"
}

func columnPath(projectId, columnId common.Id) string {
	return "/projects/" + idToStr(projectId) + "/columns/" + idToStr(columnId)
}

func columnPositionPath(projectId, columnId common.Id) string {
	return "/projects/" + idToStr(projectId) + "/columns/" + idToStr(columnId) + "/position"
}

func tasksPath(projectId, columnId common.Id) string {
	return "/projects/" + idToStr(projectId) + "/columns/" + idToStr(columnId) + "/tasks"
}

func taskPath(taskId common.Id) string {
	return "/tasks/" + idToStr(taskId)
}

func taskPositionPath(taskId common.Id) string {
	return "/tasks/" + idToStr(taskId) + "/position"
}

func commentsPath(taskId common.Id) string {
	return "/tasks/" + idToStr(taskId) + "/comments"
}

func commentPath(taskId, commentId common.Id) string {
	return "/tasks/" + idToStr(taskId) + "/comments/" + idToStr(commentId)
}

func idToStr(id common.Id) string {
	return strconv.Itoa(int(id))
}

func assertEqualBody(t *testing.T, r *http.Response, wantBody interface{}) {
	bodyJSON, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("error while read response body: %v", err)
	}
	gotBody := reflect.New(reflect.TypeOf(wantBody)).Interface()
	err = json.Unmarshal(bodyJSON, gotBody)
	if err != nil {
		t.Fatalf("error while unmarshal response body: %v", err)
	}
	gotBody = reflect.Indirect(reflect.ValueOf(gotBody)).Interface()
	if !assert.Equal(t, wantBody, gotBody, "body mismatch") {
		t.FailNow()
	}
}

func assertEqualStatusCode(t *testing.T, r *http.Response, WantCode int) {
	if !assert.Equal(t, WantCode, r.StatusCode, "status code mismatch") {
		t.FailNow()
	}
}
