package projects

import (
	"context"
	"github.com/AndreyKlimchuk/golang-learning/homework4/db/common"
	rcommon "github.com/AndreyKlimchuk/golang-learning/homework4/resources/common"
	"github.com/jackc/pgx/v4"
)

type QueryerWrap common.QueryerWrap

func (w QueryerWrap) Create(name string, description string) (rcommon.Project, error) {
	project := rcommon.Project{ProjectSettableFields: rcommon.ProjectSettableFields{Name: name, Description: description}}
	const q = "INSERT INTO projects (name, description) VALUES ($1, $2) RETURNING id"
	err := w.Q.QueryRow(context.Background(), q, name, description).Scan(&project.Id)
	return project, err
}

func (w QueryerWrap) Get(projectId rcommon.Id) (rcommon.Project, error) {
	project := rcommon.Project{Id: projectId}
	const q = "SELECT name, description FROM projects WHERE id = $1"
	err := w.Q.QueryRow(context.Background(), q, projectId).Scan(&project.Name, &project.Description)
	return project, err
}

func (w QueryerWrap) GetExpanded(projectId rcommon.Id) (rcommon.ProjectExpanded, error) {
	const q = `
		SELECT p.Id, p.name, p.description,
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
		return rcommon.ProjectExpanded{}, err
	}
	defer rows.Close()
	return buildExpanded(rows)
}

func buildExpanded(rows pgx.Rows) (rcommon.ProjectExpanded, error) {
	p := rcommon.Project{}
	c := rcommon.Column{}
	t := rcommon.Task{}
	columns := make([]rcommon.ColumnExpanded, 0, 1)
	i := -1
	for rows.Next() {
		err := rows.Scan(&p.Id, &p.Name, &p.Description, &c.Id, &c.Name, &t.Id, &t.Name, &t.Description)
		if err != nil {
			return rcommon.ProjectExpanded{}, err
		}
		if i == -1 || columns[i].Id != c.Id {
			columns = append(columns, rcommon.ColumnExpanded{Column: c, Tasks: make([]rcommon.Task, 0)})
			i++
		}
		if t.Id != 0 {
			t.ProjectId = p.Id
			t.ColumnId = c.Id
			columns[i].Tasks = append(columns[i].Tasks, t)
		}
	}
	if p.Id != 0 {
		return rcommon.ProjectExpanded{Project: p, Columns: columns}, nil
	} else {
		return rcommon.ProjectExpanded{}, common.ErrNoRows
	}
}

func (w QueryerWrap) GetMultiple() ([]rcommon.Project, error) {
	projects := []rcommon.Project{}
	const q = "SELECT id, name, description FROM projects ORDER BY name"
	rows, err := w.Q.Query(context.Background(), q)
	if err != nil {
		return projects, err
	}
	defer rows.Close()
	p := rcommon.Project{}
	for rows.Next() {
		err := rows.Scan(&p.Id, &p.Name, &p.Description)
		if err != nil {
			return projects, err
		}
		projects = append(projects, p)
	}
	return projects, nil
}

func (w QueryerWrap) Update(projectId rcommon.Id, name string, description string) error {
	const q = "UPDATE projects SET name = $2, description = $3 WHERE id = $1"
	return common.ErrorIfNoAffectedRows(w.Q.Exec(context.Background(), q, projectId, name, description))
}

func (w QueryerWrap) Delete(projectId rcommon.Id) error {
	const q = "DELETE FROM projects WHERE id = $1"
	return common.ErrorIfNoAffectedRows(w.Q.Exec(context.Background(), q, projectId))
}
