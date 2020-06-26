package resources

type Id int
type Rank string

type Project struct {
	Id Id `json:"id"`
	ProjectSettableFields
}

type ProjectExpanded struct {
	Project
	Columns []ColumnExpanded `json:"columns"`
}

type ProjectSettableFields struct {
	Name        string `json:"name" validate:"min=1,max=500"`
	Description string `json:"description" validate:"min=0,max=1000"`
}

type Column struct {
	Id Id `json:"id"`
	ColumnSettableFields
}

type ColumnExpanded struct {
	Column
	Tasks []TaskExpanded `json:"tasks"`
}

type ColumnSettableFields struct {
	Name string `json:"name" validate:"min=1,max=255"`
}

type Task struct {
	ProjectId Id `json:"project_id"`
	ColumnId  Id `json:"column_id"`
	Id        Id `json:"id"`
	TaskSettableFields
}

type TaskExpanded struct {
	Task
	Comments []Comment `json:"comments"`
}

type TaskSettableFields struct {
	Name        string `json:"name" validate:"min=1,max=500"`
	Description string `json:"description" validate:"min=0,max=5000"`
}

type Comment struct {
	Id Id `json:"id"`
	CommentSettableFields
}

type CommentSettableFields struct {
	Text string `json:"text" validate:"min=1,max=5000"`
}

type Request interface {
	Handle() (interface{}, error)
}

type Resource interface {
	GetId() Id
}

func (resource Project) GetId() Id {
	return resource.Id
}

func (resource ProjectExpanded) GetId() Id {
	return resource.Id
}

func (resource Column) GetId() Id {
	return resource.Id
}

func (resource ColumnExpanded) GetId() Id {
	return resource.Id
}

func (resource Task) GetId() Id {
	return resource.Id
}

func (resource TaskExpanded) GetId() Id {
	return resource.Id
}

func (resource Comment) GetId() Id {
	return resource.Id
}

//func GenericUpdatePosition(r UpdatePositionRequest) error {
//	tx, err := db.Begin()
//	if err != nil {
//		return rsrc.NewInternalError("cannot begin transaction", err)
//	}
//	defer db.Rollback(tx)
//	var prevRank rsrc.Rank = ""
//	if r.AfterTargetId() > 0 {
//		prevRank, err = r.GetAndBlockPrevRank(tx)
//		if db.IsNoRowsError(err) {
//			return rsrc.NewConflictError("after target doesn't exist")
//		} else if err != nil {
//			return rsrc.NewInternalError("cannot get previous task rank", err)
//		}
//	}
//	nextRank, err := r.GetNextRank(tx, prevRank)
//	if db.IsNoRowsError(err) {
//		nextRank = ""
//	} else if err != nil {
//		return rsrc.NewInternalError("cannot get next task rank", err)
//	}
//	newRank := rsrc.CalculateRank(prevRank, nextRank)
//	err = r.UpdatePositionFinal(tx, newRank)
//	if err != nil {
//		return rsrc.NewNotFoundOrInternalError("cannot update position", err)
//	}
//	if err := db.Commit(tx); err != nil {
//		return rsrc.NewInternalError("cannot commit transaction", err)
//	}
//}

// naive implementation of lexicographic ranking algorithm
func CalculateRank(rankA Rank, rankB Rank) Rank {
	var smaller, bigger Rank
	var upperBound byte
	if rankA < rankB {
		smaller, bigger = rankA, rankB
	} else {
		smaller, bigger = rankB, rankA
	}
	res := make([]byte, 0)
	for i := 0; ; i++ {
		if i < len(smaller) {
			res = append(res, smaller[i])
		} else {
			res = append(res, 'a')
		}
		if i < len(bigger) {
			upperBound = bigger[i]
		} else {
			upperBound = 'z' + 1
		}
		if upperBound > res[i]+1 {
			res[i] += (upperBound - res[i]) / 2
			return Rank(res)
		}
	}
}
