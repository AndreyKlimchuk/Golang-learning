package homework4

import "errors"

type column struct {
	Id    id     `json:"id"`
	Tasks []task `json:"tasks,omitempty"`
	columnSettableFields
}

type columnSettableFields struct {
	Name string `json:"name" validate:"min=1,max=255"`
}

type createColumn struct {
	ProjectId id
	columnSettableFields
}

type readColumn struct {
	ProjectId id
	ColumnId  id
}

type readColumnCollection struct {
	ProjectId id
}

type updateColumn struct {
	ProjectId id
	ColumnId  id
	columnSettableFields
}

type deleteColumn struct {
	ProjectId id
	ColumnId  id
}

type updateColumnPosition struct {
	ProjectId     id
	ColumnId      id
	AfterColumnId id
}

func (c column) GetId() id {
	return c.Id
}

func (r createColumn) Create() (column, error) {
	maxRank, err := db_blocking_select_max_column_rank(r.ProjectId)
	if err != nil {
		return column{}, errors.New("500")
	}
	maxRank = calculateRank(maxRank, "")
	return db_create_column(r, maxRank)
}

func (r readColumn) Read() (column, error) {
	return db_read_column(r)
}

func (r readColumnCollection) ReadCollection() ([]column, error) {
	return db_read_column_collection(r)
}

func (r updateColumn) Update() error {
	return db_update_column(r)
}

func (r deleteColumn) Delete() error {
	rank, err := db_delete_column(r.ProjectId, r.ColumnId)
	if err != nil {
		return errors.New("500")
	} else if rank == "" {
		return errors.New("404")
	} else {
		successorColumnId, err := db_blocking_select_successor_column_id(r.ProjectId, rank)
		if err != nil {
			return errors.New("500")
		} else if successorColumnId == 0 {
			return errors.New("409")
		} else {
			tasksIds, err := db_blocking_select_tasks_ids(r.ColumnId)
			if err != nil {
				return errors.New("500")
			}
			maxRank, err := db_blocking_select_max_task_rank(successorColumnId)
			if err != nil {
				return errors.New("500")
			}
			for _, Id := range tasksIds {
				maxRank = calculateRank(maxRank, "")
				updateTaskPosition{
					ProjectId:   r.ProjectId,
					ColumnId:    r.ColumnId,
					TaskId:      Id,
					NewColumnId: successorColumnId,
				}.UpdatePosition(maxRank)
			}
			return nil
		}
	}
}

func (r updateColumnPosition) GetNearestRanks() ([2]rank, error) {
	return db_blocking_select_nearest_column_ranks(r)
}

func (r updateColumnPosition) UpdatePosition(rank rank) error {
	return db_update_column_rank(r, rank)
}

func db_blocking_select_max_column_rank(columnId id) (rank, error) {
	return "", nil
}

func db_create_column(r createColumn, maxRank rank) (column, error) {
	return column{}, nil
}

func db_read_column(r readColumn) (column, error) {
	return column{}, nil
}

func db_read_column_collection(r readColumnCollection) ([]column, error) {
	return []column{}, nil
}

func db_update_column(r updateColumn) error {
	return nil
}

func db_delete_column(projectId id, columnId id) (rank, error) {
	// sql := "DELETE FROM columns WHERE project_id = ? AND id = ? RETURNING rank"
	return "", nil
}

func db_blocking_select_successor_column_id(projectId id, r rank) (id, error) {
	// sql := "
	// 	WITH before AS (
	// 		SELECT id FROM columns
	// 		WHERE project_id = ? AND rank < ?
	// 		ORDER BY rank
	// 		LIMIT 1
	// 		FOR UPDATE
	// 	WITH after AS (
	// 		SELECT id FROM columns
	// 		WHERE NOT EXIST(SELECT * FROM before) AND project_id = ? AND rank > ?
	// 		ORDER BY rank
	// 		LIMIT 1
	// 		FOR UPDATE
	// 	)
	// 	SELECT * FROM before UNION SELECT * FROM after
	// "
	return 0, nil
}

func db_blocking_select_tasks_ids(columnId id) ([]id, error) {
	return []id{}, nil
}

func db_blocking_select_max_task_rank(columnId id) (rank, error) {
	return "", nil
}

func db_blocking_select_nearest_column_ranks(r updateColumnPosition) ([2]rank, error) {
	return [2]rank{"", ""}, nil
}

func db_update_column_rank(r updateColumnPosition, rank rank) error {
	return nil
}
