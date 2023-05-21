package warhammer

type WhModifiers struct {
	Size       int `json:"size" validate:"min=-3,max=3"`
	Movement   int `json:"movement" validate:"min=-3,max=3"`
	Attributes WHAttributes
}

func (m WhModifiers) Copy() WhModifiers {
	return WhModifiers{
		Size:       m.Size,
		Movement:   m.Movement,
		Attributes: m.Attributes.Copy(),
	}
}

type WHAttributes struct {
	WS  int `json:"WS" validate:"min=-99,max=99"`
	BS  int `json:"BS" validate:"min=-99,max=99"`
	S   int `json:"S" validate:"min=-99,max=99"`
	T   int `json:"T" validate:"min=-99,max=99"`
	I   int `json:"I" validate:"min=-99,max=99"`
	Ag  int `json:"Ag" validate:"min=-99,max=99"`
	Dex int `json:"Dex" validate:"min=-99,max=99"`
	Int int `json:"Int" validate:"min=-99,max=99"`
	WP  int `json:"WP" validate:"min=-99,max=99"`
	Fel int `json:"Fel" validate:"min=-99,max=99"`
}

func (a WHAttributes) Copy() WHAttributes {
	return WHAttributes{
		WS:  a.WS,
		BS:  a.BS,
		S:   a.S,
		T:   a.T,
		I:   a.I,
		Ag:  a.Ag,
		Dex: a.Dex,
		Int: a.Int,
		WP:  a.WP,
		Fel: a.Fel,
	}
}
