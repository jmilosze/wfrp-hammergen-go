package mock_data

import (
	"fmt"
	"github.com/jmilosze/wfrp-hammergen-go/internal/domain/warhammer"
)

var itemMelee = warhammer.Wh{
	Id:      "400000000000000000000000",
	OwnerId: user1.Id,
	Object: warhammer.WhItem{
		Name:        "melee item",
		Description: fmt.Sprintf("owned by %s", user1.Username),
		Price:       2.31,
		Enc:         1.5,
		Properties:  []string{property0.Id, property1.Id},
		Type:        0,
		Shared:      true,
		Source: map[warhammer.WhSource]string{
			warhammer.WhSourceArchivesOfTheEmpireVolI: "g",
			warhammer.WhSourceUpInArms:                "f",
		},

		Melee: warhammer.WhItemMelee{
			Hands:     1,
			Dmg:       5,
			DmgSbMult: 1.0,
			Reach:     3,
			Group:     0,
		},
	},
}

func NewMockItems() []*warhammer.Wh {
	return []*warhammer.Wh{&itemMelee}
}
