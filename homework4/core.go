package homework4

type id int
type rank string

type Resource interface {
	GetId() id
}

type CreateResource interface {
	Create() (Resource, error)
}

type ReadResource interface {
	Read() (Resource, error)
}

type ReadResourceCollection interface {
	ReadCollection() ([]Resource, error)
}

type UpdateResource interface {
	Update() error
}

type DeleteResource interface {
	Delete() error
}

type UpdateResourcePosition interface {
	GetNearestRanks() ([2]rank, error)
	UpdatePosition(rank) error
}

func Create(r CreateResource) (Resource, error) {
	return r.Create()
}

func Read(r ReadResource) (Resource, error) {
	return r.Read()
}

func ReadCollection(r ReadResourceCollection) ([]Resource, error) {
	return r.ReadCollection()
}

func Update(r UpdateResource) error {
	return r.Update()
}

func Delete(r DeleteResource) error {
	return r.Delete()
}

func UpdatePosition(r UpdateResourcePosition) error {
	nearestRanks, err := r.GetNearestRanks()
	if err != nil {
		return err
	}
	newRank := calculateRank(nearestRanks[0], nearestRanks[1])
	return r.UpdatePosition(newRank)
}

// naive implementation of lexicographic ranking algorithm
func calculateRank(rankA rank, rankB rank) rank {
	var smaller, bigger rank
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
			return rank(res)
		}
	}
}
