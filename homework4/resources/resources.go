package resources

import "github.com/AndreyKlimchuk/golang-learning/homework4/resources/common"

type Request interface {
	// If error is not nil, first value should be ignored
	Handle() (interface{}, error)
}

type Resource interface {
	GetId() common.Id
}
