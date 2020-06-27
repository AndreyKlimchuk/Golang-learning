package projects

import (
	"context"
	"github.com/AndreyKlimchuk/golang-learning/homework4/db/common"
	rsrc "github.com/AndreyKlimchuk/golang-learning/homework4/resources"
	"github.com/jackc/pgx/v4"
)

type QueryerWrap common.QueryerWrap

func (w QueryerWrap) Create(name string, description string) (rsrc.Project, error) {
	project := rsrc.Project{ProjectSettableFields: rsrc.ProjectSettableFields{Name: name, Description: description}}
	const q = "INSERT INTO projects (name, description) VALUES ($1, $2) RETURNING id"
	err := w.Q.QueryRow(context.Background(), q, name, description).Scan(&project.Id)
	return project, err
}

func (w QueryerWrap) Get(projectId rsrc.Id) (rsrc.Project, error) {
	project := rsrc.Project{Id: projectId}
	const q = "SELECT name, description FROM projects WHERE id = $1"
	err := w.Q.QueryRow(context.Background(), q, projectId).Scan(&project.Name, &project.Description)
	return project, err
}

func (w QueryerWrap) GetExpanded(projectId rsrc.Id) (rsrc.ProjectExpanded, error) {
	const q = `
		SELECT p.name, p.description,
			   c.id, c.name,
			   COALESCE(t.id, 0), COALESCE(t.name, ''), COALESCE(t.description, '')
		FROM projects p
		JOIN columns c ON p.id = c.project_id
		LEFT JOIN tasks t ON c.id = t.column_id
		WHERE p.id = $1
		ORDER BY c.rank, t.rank ASC
	`
	rows, err := w.Q.Query(context.Background(), q, projectId)
	if err != nil {
		return rsrc.ProjectExpanded{}, err
	}
	defer rows.Close()
	return buildExpanded(projectId, rows)
}

func buildExpanded(projectId rsrc.Id, rows pgx.Rows) (rsrc.ProjectExpanded, error) {
	p := rsrc.Project{Id: projectId}
	c := rsrc.Column{}
	t := rsrc.Task{}
	columns := make([]rsrc.ColumnExpanded, 0, 1)
	i := -1
	for rows.Next() {
		err := rows.Scan(&p.Name, &p.Description, &c.Id, &c.Name, &t.Id, &t.Name, &t.Description)
		if err != nil {
			return rsrc.ProjectExpanded{}, err
		}
		if i == -1 || columns[i].Id != c.Id {
			columns = append(columns, rsrc.ColumnExpanded{Column: c, Tasks: make([]rsrc.Task, 0)})
			i++
		}
		if t.Id != 0 {
			columns[i].Tasks = append(columns[i].Tasks, t)
		}
	}
	if len(columns) > 0 {
		return rsrc.ProjectExpanded{Project: p, Columns: columns}, nil
	} else {
		return rsrc.ProjectExpanded{}, common.ErrNoRows
	}
}

func (w QueryerWrap) GetMultiple() ([]rsrc.Project, error) {
	const q = "SELECT id, name, description FROM projects ORDER BY name"
	rows, err := w.Q.Query(context.Background(), q)
	if err != nil {
		return []rsrc.Project{}, err
	}
	defer rows.Close()
	p := rsrc.Project{}
	projects := make([]rsrc.Project, 0)
	for rows.Next() {
		err := rows.Scan(p.Id, p.Name, p.Description)
		if err != nil {
			return []rsrc.Project{}, err
		}
		projects = append(projects, p)
	}
	return projects, nil
}

func (w QueryerWrap) Update(projectId rsrc.Id, name string, description string) error {
	const q = "UPDATE projects SET name = $2, description = $3 WHERE id = $1"
	return common.ErrorIfNoAffectedRows(w.Q.Exec(context.Background(), q, projectId, name, description))
}

func (w QueryerWrap) Delete(projectId rsrc.Id) error {
	const q = "DELETE projects WHERE id = $1"
	return common.ErrorIfNoAffectedRows(w.Q.Exec(context.Background(), q, projectId))
}
