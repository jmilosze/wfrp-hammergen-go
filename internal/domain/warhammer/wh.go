package warhammer

import (
	"fmt"
	"strings"
)

type Wh struct {
	Id      string
	OwnerId string
	CanEdit bool
	Object  WhObject
}

func NewWh(t WhType) (Wh, error) {
	var wh Wh

	switch t {
	case WhTypeMutation:
		wh.Object = &WhMutation{}
	case WhTypeSpell:
		wh.Object = &WhSpell{}
	case WhTypeProperty:
		wh.Object = &WhProperty{}
	default:
		return wh, fmt.Errorf("invalid Wh type %s", t)
	}

	return wh, nil
}

func (w *Wh) Copy() *Wh {
	if w == nil {
		return nil
	}

	return &Wh{
		Id:      strings.Clone(w.Id),
		OwnerId: strings.Clone(w.OwnerId),
		Object:  w.Object.Copy(),
	}
}

func (w *Wh) IsShared() bool {
	return w.Object.IsShared()
}

type WhObject interface {
	Copy() WhObject
	IsShared() bool
}
