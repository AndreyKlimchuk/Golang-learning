package api

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/AndreyKlimchuk/golang-learning/homework4/docs"
	"github.com/AndreyKlimchuk/golang-learning/homework4/resources/columns"
	"github.com/AndreyKlimchuk/golang-learning/homework4/resources/comments"
	_ "github.com/AndreyKlimchuk/golang-learning/homework4/resources/common"
	"github.com/AndreyKlimchuk/golang-learning/homework4/resources/projects"
	"github.com/AndreyKlimchuk/golang-learning/homework4/resources/tasks"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// @title API
// @version 1.0
// @description This is Trello-like task management application

// @contact.name Andrew
// @contact.email ua.challenger@gmail.com

// @host localhost:8080
// @BasePath /api/v1

const basePath = "/api/v1"

func NewRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)

	r.Route(basePath, func(r chi.Router) {
		r.Route("/projects", func(r chi.Router) {
			r.Post("/", createProject)
			r.Get("/", getProjects)

			r.Route("/{projectID:[\\d]+}", func(r chi.Router) {
				r.Get("/", getProject)
				r.Put("/", updateProject)
				r.Delete("/", deleteProject)

				r.Route("/columns", func(r chi.Router) {
					r.Post("/", createColumn)
					r.Get("/", getColumns)

					r.Route("/{columnID:[\\d]+}", func(r chi.Router) {
						r.Get("/", getColumn)
						r.Put("/", updateColumn)
						r.Delete("/", deleteColumn)

						r.Put("/position", updateColumnPosition)
						r.Post("/tasks", createTask)
					})
				})
			})
		})

		r.Route("/tasks", func(r chi.Router) {
			r.Route("/{taskID:[\\d]+}", func(r chi.Router) {
				r.Get("/", getTask)
				r.Put("/", updateTask)
				r.Delete("/", deleteTask)

				r.Put("/position", updateTaskPosition)

				r.Route("/comments", func(r chi.Router) {
					r.Post("/", createComment)
					r.Get("/", getComments)

					r.Route("/{commentID:[\\d]+}", func(r chi.Router) {
						r.Get("/", getComment)
						r.Put("/", updateComment)
						r.Delete("/", deleteComment)
					})
				})
			})
		})
	})

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))
	return r
}

// createProject godoc
// @Summary Create project
// @Description Create new project with single "default" column
// @Tags projects
// @Accept  json
// @Produce  json
// @Param body body common.ProjectSettableFields true "request body"
// @Success 201 {object} common.ProjectExpanded
// @Header 201 {string} Location "/project/1"
// @Router /projects [post]
func createProject(w http.ResponseWriter, httpReq *http.Request) {
	var req = projects.CreateRequest{}
	handleRequest(w, httpReq, &req)
}

// getProjects godoc
// @Summary Get projects
// @Description Get all projects
// @Tags projects
// @Produce  json
// @Success 200 {array} common.Project{}
// @Router /projects [get]
func getProjects(w http.ResponseWriter, httpReq *http.Request) {
	var req = projects.ReadCollectionRequest{}
	handleRequest(w, httpReq, &req)
}

// getProject godoc
// @Summary Get project
// @Description Get project
// @Tags projects
// @Accept  json
// @Produce  json
// @Param project_id path int true "Project ID"
// @Param expanded query bool false "expand by sub-resources" default(false)
// @Success 200 {object} common.ProjectExpanded
// @Router /projects/{project_id} [get]
func getProject(w http.ResponseWriter, httpReq *http.Request) {
	var req = projects.ReadRequest{
		ProjectId: getProjectId(httpReq),
		Expanded:  getExpanded(httpReq),
	}
	handleRequest(w, httpReq, &req)
}

// updateProject godoc
// @Summary Update project
// @Description Update project
// @Tags projects
// @Accept  json
// @Param project_id path int true "Project ID"
// @Param body body common.ProjectSettableFields true "request body"
// @Success 204
// @Router /projects/{project_id} [put]
func updateProject(w http.ResponseWriter, httpReq *http.Request) {
	var req = projects.UpdateRequest{
		ProjectId: getProjectId(httpReq),
	}
	handleRequest(w, httpReq, &req)
}

// deleteProject godoc
// @Summary Delete project
// @Description Delete project and all sub-resources
// @Tags projects
// @Param project_id path int true "Project ID"
// @Success 204
// @Router /projects/{project_id} [delete]
func deleteProject(w http.ResponseWriter, httpReq *http.Request) {
	var req = projects.DeleteRequest{
		ProjectId: getProjectId(httpReq),
	}
	handleRequest(w, httpReq, &req)
}

// createColumn godoc
// @Description Create new column
// @Summary Create column
// @Tags columns
// @Accept  json
// @Produce  json
// @Param project_id path int true "Project ID"
// @Param body body common.ColumnSettableFields true "request body"
// @Success 201 {object} common.Column
// @Header 201 {string} Location "/project/1/columns/1"
// @Router /projects/{project_id}/columns [post]
func createColumn(w http.ResponseWriter, httpReq *http.Request) {
	var req = columns.CreateRequest{
		ProjectId: getProjectId(httpReq),
	}
	handleRequest(w, httpReq, &req)
}

// getColumns godoc
// @Summary Get columns
// @Description Get all columns within project
// @Tags columns
// @Produce  json
// @Param project_id path int true "Project ID"
// @Success 200 {array} common.Column{}
// @Router /projects/{project_id}/columns [get]
func getColumns(w http.ResponseWriter, httpReq *http.Request) {
	var req = columns.ReadCollectionRequest{
		ProjectId: getProjectId(httpReq),
	}
	handleRequest(w, httpReq, &req)
}

// getColumn godoc
// @Summary Get column
// @Description Get column
// @Tags columns
// @Produce  json
// @Param project_id path int true "Project ID"
// @Param column_id path int true "Column ID"
// @Success 200 {object} common.Column
// @Router /projects/{project_id}/columns/{column_id} [get]
func getColumn(w http.ResponseWriter, httpReq *http.Request) {
	var req = columns.ReadRequest{
		ProjectId: getProjectId(httpReq),
		ColumnId:  getColumnId(httpReq),
	}
	handleRequest(w, httpReq, &req)
}

// updateColumn godoc
// @Summary Update column
// @Description Update column
// @Tags columns
// @Accept  json
// @Param project_id path int true "Project ID"
// @Param column_id path int true "Column ID"
// @Param body body common.ColumnSettableFields true "request body"
// @Success 204
// @Router /projects/{project_id}/columns/{column_id} [put]
func updateColumn(w http.ResponseWriter, httpReq *http.Request) {
	var req = columns.UpdateRequest{
		ProjectId: getProjectId(httpReq),
		ColumnId:  getColumnId(httpReq),
	}
	handleRequest(w, httpReq, &req)
}

// updateColumnPosition godoc
// @Summary Update column's position
// @Description Place column after column specified by after_column_id
// @Description if it is grater than 0, otherwise at the beginning
// @Tags columns
// @Accept  json
// @Param project_id path int true "Project ID"
// @Param column_id path int true "Column ID"
// @Param body body columns.UpdatePositionRequestBody true "request body"
// @Success 204
// @Router /projects/{project_id}/columns/{column_id}/position [put]
func updateColumnPosition(w http.ResponseWriter, httpReq *http.Request) {
	var req = columns.UpdatePositionRequest{
		ProjectId: getProjectId(httpReq),
		ColumnId:  getColumnId(httpReq),
	}
	handleRequest(w, httpReq, &req)
}

// deleteColumn godoc
// @Summary Delete column
// @Description Delete column and move all tasks to the neighbor
// @Tags columns
// @Param project_id path int true "Project ID"
// @Param column_id path int true "Column ID"
// @Success 204
// @Router /projects/{project_id}/columns/{column_id} [delete]
func deleteColumn(w http.ResponseWriter, httpReq *http.Request) {
	var req = columns.DeleteRequest{
		ProjectId: getProjectId(httpReq),
		ColumnId:  getColumnId(httpReq),
	}
	handleRequest(w, httpReq, &req)
}

// createTask godoc
// @Summary Create task
// @Description Create new task
// @Tags tasks
// @Accept  json
// @Produce  json
// @Param project_id path int true "Project ID"
// @Param column_id path int true "Column ID"
// @Param body body common.TaskSettableFields true "request body"
// @Success 201 {object} common.Task
// @Header 201 {string} Location "/tasks/1"
// @Router /projects/{project_id}/columns/{column_id}/tasks [post]
func createTask(w http.ResponseWriter, httpReq *http.Request) {
	var req = tasks.CreateRequest{
		ProjectId: getProjectId(httpReq),
		ColumnId:  getColumnId(httpReq),
	}
	handleRequest(w, httpReq, &req)
}

// getTask godoc
// @Summary Get task
// @Description Get task
// @Tags tasks
// @Produce  json
// @Param task_id path int true "Task ID"
// @Param expanded query bool false "expand by sub-resources" default(false)
// @Success 200 {object} common.TaskExpanded
// @Router /tasks/{task_id} [get]
func getTask(w http.ResponseWriter, httpReq *http.Request) {
	var req = tasks.ReadRequest{
		TaskId:   getTaskId(httpReq),
		Expanded: getExpanded(httpReq),
	}
	handleRequest(w, httpReq, &req)
}

// updateTask godoc
// @Summary Update task
// @Description Update task
// @Tags tasks
// @Accept  json
// @Param task_id path int true "Task ID"
// @Param body body common.TaskSettableFields true "request body"
// @Success 204
// @Router /tasks/{task_id} [put]
func updateTask(w http.ResponseWriter, httpReq *http.Request) {
	var req = tasks.UpdateRequest{
		TaskId: getTaskId(httpReq),
	}
	handleRequest(w, httpReq, &req)
}

// updateTaskPosition godoc
// @Summary Update task's position
// @Description Place task after task specified by after_task_id
// @Description if it is grater than 0, otherwise at the top of specified by new_column_id column
// @Tags tasks
// @Accept  json
// @Param task_id path int true "Task ID"
// @Param body body tasks.UpdatePositionRequestBody true "request body"
// @Success 204
// @Router /tasks/{task_id}/position [put]
func updateTaskPosition(w http.ResponseWriter, httpReq *http.Request) {
	var req = tasks.UpdatePositionRequest{
		TaskId: getTaskId(httpReq),
	}
	handleRequest(w, httpReq, &req)
}

// deleteTask godoc
// @Summary Delete task
// @Description Delete task with all sub-resources
// @Tags tasks
// @Param task_id path int true "Task ID"
// @Success 204
// @Router /tasks/{task_id} [delete]
func deleteTask(w http.ResponseWriter, httpReq *http.Request) {
	var req = tasks.DeleteRequest{
		TaskId: getTaskId(httpReq),
	}
	handleRequest(w, httpReq, &req)
}

// createComment godoc
// @Summary Create comment
// @Description Create new comment
// @Tags comments
// @Accept  json
// @Produce  json
// @Param task_id path int true "Task ID"
// @Param body body common.CommentSettableFields true "request body"
// @Success 201 {object} common.Comment
// @Header 201 {string} Location "/tasks/1/comments/1"
// @Router /tasks/{task_id}/comments [post]
func createComment(w http.ResponseWriter, httpReq *http.Request) {
	var req = comments.CreateRequest{
		TaskId: getTaskId(httpReq),
	}
	handleRequest(w, httpReq, &req)
}

// getComments godoc
// @Summary Get comments
// @Description Get all comments within task
// @Tags comments
// @Produce  json
// @Param task_id path int true "Task ID"
// @Success 200 {array} common.Comment{}
// @Router /tasks/{task_id}/comments [get]
func getComments(w http.ResponseWriter, httpReq *http.Request) {
	var req = comments.ReadCollectionRequest{
		TaskId: getTaskId(httpReq),
	}
	handleRequest(w, httpReq, &req)
}

// getComment godoc
// @Summary Get comment
// @Description Get comment
// @Tags comments
// @Produce  json
// @Param task_id path int true "Task ID"
// @Param comment_id path int true "Comment ID"
// @Success 200 {object} common.Comment
// @Router /tasks/{task_id}/comments/{comment_id} [get]
func getComment(w http.ResponseWriter, httpReq *http.Request) {
	var req = comments.ReadRequest{
		TaskId:    getTaskId(httpReq),
		CommentId: getCommentId(httpReq),
	}
	handleRequest(w, httpReq, &req)
}

// updateComment godoc
// @Summary Update comment
// @Description Update comment
// @Tags comments
// @Accept  json
// @Param task_id path int true "Task ID"
// @Param comment_id path int true "Comment ID"
// @Param body body common.CommentSettableFields true "request body"
// @Success 204
// @Router /tasks/{task_id}/comments/{comment_id} [put]
func updateComment(w http.ResponseWriter, httpReq *http.Request) {
	var req = comments.UpdateRequest{
		TaskId:    getTaskId(httpReq),
		CommentId: getCommentId(httpReq),
	}
	handleRequest(w, httpReq, &req)
}

// deleteComment godoc
// @Summary Delete comment
// @Description Delete comment
// @Tags comments
// @Param task_id path int true "Task ID"
// @Param comment_id path int true "Comment ID"
// @Success 204
// @Router /tasks/{task_id}/comments/{comment_id} [delete]
func deleteComment(w http.ResponseWriter, httpReq *http.Request) {
	var req = comments.DeleteRequest{
		TaskId:    getTaskId(httpReq),
		CommentId: getCommentId(httpReq),
	}
	handleRequest(w, httpReq, &req)
}
