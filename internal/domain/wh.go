package domain

var WhTypes = []string{"mutation", "spell"}

const (
	MutationPhysical int = 0
	MutationMental       = 1
)

type MutationWrite struct {
	Name        string `json:"name" validate:"omitempty,min=0,max=200,excludesall=<>"`
	Description string `json:"description" validate:"omitempty,min=0,max=100000,excludesall=<>"`
	Shared      *bool  `json:"shared" validate:"omitempty"`
	Type        *int   `json:"type" validate:"omitempty,oneof 0 1"`
}

type Mutation struct {
	MutationWrite
	Id      string
	OwnerId bool
}

type SpellWrite struct {
	Name        string `json:"name" validate:"omitempty,min=0,max=200,excludesall=<>"`
	Cn          *int   `json:"cn" validate:"omitempty,min=-1,max=99"`
	Range       string `json:"range" validate:"omitempty,min=0,max=200,excludesall=<>"`
	Target      string `json:"target" validate:"omitempty,min=0,max=200,excludesall=<>"`
	Duration    string `json:"duration" validate:"omitempty,min=0,max=200,excludesall=<>"`
	Description string `json:"description" validate:"omitempty,min=0,max=100000,excludesall=<>"`
}

type Spell struct {
	SpellWrite
	Id      string
	OwnerId bool
}
