package homework4

type project struct {
	Id      id       `json:"id"`
	Columns []column `json:"columns,omitempty"`
	projectSettableFields
}

type projectSettableFields struct {
	Name        string `json:"name" validate:"min=1,max=500"`
	Description string `json:"description" validate:"min=0,max=1000"`
}

type createProject struct {
	projectSettableFields
}

type readProject struct {
	ProjectId id
	Expanded  bool
}

type readProjectCollection struct {
}

type updateProject struct {
	ProjectId id
	projectSettableFields
}

type deleteProject struct {
	ProjectId id
}

func (p project) GetId() id {
	return p.Id
}

func (r createProject) Create() (project, error) {
	return project{}, nil
}

func (r readProject) Read() (project, error) {
	return project{}, nil
}

func (r readProjectCollection) ReadCollection() ([]project, error) {
	return []project{}, nil
}

func (r updateProject) Update() error {
	return nil
}

func (r deleteProject) Delete() error {
	return nil
}
