package homework4

type comment struct {
	Id id `json:"id"`
	commentSettableFields
}

type commentSettableFields struct {
	Text string `json:"text" validate:"min=1,max=5000"`
}

type createComment struct {
	ProjectId id
	ColumnId  id
	TaskId    id
	commentSettableFields
}

type readComment struct {
	ProjectId id
	ColumnId  id
	TaskId    id
	CommentId id
}

type readCommentCollection struct {
	ProjectId id
	ColumnId  id
	TaskId    id
}

type updateComment struct {
	ProjectId id
	ColumnId  id
	TaskId    id
	CommentId id
	commentSettableFields
}

type deleteComment struct {
	ProjectId id
	ColumnId  id
	TaskId    id
	CommentId id
}

func (c comment) GetId() id {
	return c.Id
}

func (r createComment) Create() (comment, error) {
	return comment{}, nil
}

func (r readComment) Read() (comment, error) {
	return comment{}, nil
}

func (r readCommentCollection) ReadCollection() ([]comment, error) {
	return []comment{}, nil
}

func (r updateComment) Update() error {
	return nil
}

func (r deleteComment) Delete() error {
	return nil
}
