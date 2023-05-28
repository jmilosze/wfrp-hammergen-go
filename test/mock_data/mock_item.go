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
		Type:        warhammer.WhItemTypeMelee,
		Shared:      true,
		Source: map[warhammer.WhSource]string{
			warhammer.WhSourceArchivesOfTheEmpireVolI: "g",
			warhammer.WhSourceUpInArms:                "f",
		},

		Melee: warhammer.WhItemMelee{
			Hands:     warhammer.WhItemHandsOne,
			Dmg:       5,
			DmgSbMult: 1.0,
			Reach:     warhammer.WhItemMeleeReachAverage,
			Group:     warhammer.WhItemMeleeGroupBasic,
		},
	},
}

var itemRanged = warhammer.Wh{
	Id:      "400000000000000000000001",
	OwnerId: user1.Id,
	Object: warhammer.WhItem{
		Name:        "ranged item",
		Description: fmt.Sprintf("owned by %s", user1.Username),
		Price:       5,
		Enc:         8,
		Type:        warhammer.WhItemTypeRanged,
		Shared:      true,

		Ranged: warhammer.WhItemRanged{
			Hands:     warhammer.WhItemHandsOne,
			Dmg:       2,
			DmgSbMult: 2.0,
			Rng:       100,
			RngSbMult: 0,
			Group:     warhammer.WhItemRangedGroupCrossbow,
		},
	},
}

func NewMockItems() []*warhammer.Wh {
	return []*warhammer.Wh{&itemMelee, &itemRanged}
}