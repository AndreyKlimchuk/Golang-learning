package homework4

type task struct {
	Id       id        `json:"id"`
	Comments []comment `json:"comments,omitempty"`
	taskSettableFields
}

type taskSettableFields struct {
	Name        string `json:"name" validate:"min=1,max=500"`
	Description string `json:"description" validate:"min=0,max=5000"`
}

type createTask struct {
	ProjectId id
	ColumnId  id
	taskSettableFields
}

type readTask struct {
	ProjectId id
	ColumnId  id
	TaskId    id
	expanded  bool
}

type readTaskCollection struct {
	ProjectId id
	ColumnId  id
}

type updateTask struct {
	ProjectId id
	ColumnId  id
	TaskId    id
	taskSettableFields
}

type deleteTask struct {
	ProjectId id
	ColumnId  id
	TaskId    id
}

type updateTaskPosition struct {
	ProjectId   id
	ColumnId    id
	TaskId      id
	NewColumnId id
	AfterTaskId id
}

func (t task) GetId() id {
	return t.Id
}

func (r createTask) Create() (task, error) {
	return task{}, nil
}

func (r readTask) Read() (task, error) {
	return task{}, nil
}

func (r readTaskCollection) ReadCollection() ([]task, error) {
	return []task{}, nil
}

func (r updateTask) Update() error {
	return nil
}

func (r deleteTask) Delete() error {
	return nil
}

func (u updateTaskPosition) GetNearestRanks() ([2]rank, error) {
	return [2]rank{"", ""}, nil
}

func (u updateTaskPosition) UpdatePosition(r rank) error {
	return nil
}
