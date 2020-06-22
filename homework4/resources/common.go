package resources

type Id int
type Rank string

type Project struct {
	Id      Id       `json:"id"`
	Columns []Column `json:"columns,omitempty"`
	ProjectSettableFields
}

type ProjectSettableFields struct {
	Name        string `json:"name" validate:"min=1,max=500"`
	Description string `json:"description" validate:"min=0,max=1000"`
}

type Column struct {
	Id    Id     `json:"id"`
	Tasks []Task `json:"tasks,omitempty"`
	ColumnSettableFields
}

type ColumnSettableFields struct {
	Name string `json:"name" validate:"min=1,max=255"`
}

type Task struct {
	ProjectId Id        `json:"project_id"`
	ColumnId  Id        `json:"column_id"`
	Id        Id        `json:"id"`
	TaskSettableFields
	Comments  []Comment `json:"comments,omitempty"`
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

//func GenericUpdatePosition(r UpdatePositionRequest) error {
//	tx, err := pg.Begin()
//	if err != nil {
//		return rsrc.NewInternalError("cannot begin transaction", err)
//	}
//	defer pg.Rollback(tx)
//	var prevRank rsrc.Rank = ""
//	if r.AfterTargetId() > 0 {
//		prevRank, err = r.GetAndBlockPrevRank(tx)
//		if pg.IsNoRowsError(err) {
//			return rsrc.NewConflictError("after target doesn't exist")
//		} else if err != nil {
//			return rsrc.NewInternalError("cannot get previous task rank", err)
//		}
//	}
//	nextRank, err := r.GetNextRank(tx, prevRank)
//	if pg.IsNoRowsError(err) {
//		nextRank = ""
//	} else if err != nil {
//		return rsrc.NewInternalError("cannot get next task rank", err)
//	}
//	newRank := rsrc.CalculateRank(prevRank, nextRank)
//	err = r.UpdatePositionFinal(tx, newRank)
//	if err != nil {
//		return rsrc.NewNotFoundOrInternalError("cannot update position", err)
//	}
//	if err := pg.Commit(tx); err != nil {
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
