package common

type Id int
type Rank string

const DefaultColumnName string = "default"

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
	Tasks []Task `json:"tasks"`
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

func CalculateRankHigher(rank Rank) Rank {
	return CalculateRankBetween(rank, "{{{{{{{{{{{{{{{{")
}

func CalculateRankInitial() Rank {
	return CalculateRankBetween("", "")
}

// naive implementation of lexicographic ranking algorithm
func CalculateRankBetween(rankA Rank, rankB Rank) Rank {
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
			upperBound = '{' // 'z' + 1
		}
		if upperBound > res[i]+1 {
			res[i] += (upperBound - res[i]) / 2
			return Rank(res)
		}
	}
}
