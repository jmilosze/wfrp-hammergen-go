package mock_data

import (
	"fmt"
	"github.com/jmilosze/wfrp-hammergen-go/internal/domain/warhammer"
)

var character0 = warhammer.Wh{
	Id:      "800000000000000000000000",
	OwnerId: user1.Id,
	Object: warhammer.WhCharacter{
		Name:        "character 0",
		Description: fmt.Sprintf("owned by %s", user1.Username),
		Notes:       "some notes",
		EquippedItems: warhammer.WhIdNumberMap{
			itemArmour.Id: 1,
			itemMelee.Id:  2,
		},
		CarriedItems: warhammer.WhIdNumberMap{
			itemRanged.Id:     1,
			itemAmmunition.Id: 200,
		},
		StoredItems: warhammer.WhIdNumberMap{
			itemGrimoire.Id: 1,
			itemOther.Id:    100,
		},
		Skills: warhammer.WhIdNumberMap{
			skill0.Id: 1,
			skill1.Id: 10,
		},
		Talents: warhammer.WhIdNumberMap{
			talent0.Id: 1,
			talent1.Id: 5,
		},
		Species: warhammer.WhCharacterSpeciesHalflingBrandysnap,
		BaseAttributes: warhammer.WhAttributes{
			WS:  1,
			BS:  2,
			S:   3,
			T:   4,
			I:   5,
			Ag:  6,
			Dex: 7,
			Int: 8,
			WP:  9,
			Fel: 10,
		},
		AttributeAdvances: warhammer.WhAttributes{
			WS:  10,
			BS:  9,
			S:   8,
			T:   7,
			I:   6,
			Ag:  5,
			Dex: 4,
			Int: 3,
			WP:  2,
			Fel: 1,
		},
		CareerPath: []string{career0.Id},
		Career:     career1.Id,
		Fate:       2,
		Fortune:    1,
		Resilience: 3,
		Resolve:    1,
		CurrentExp: 100,
		SpentExp:   1500,
		Status:     warhammer.WhStatusSilver,
		Standing:   warhammer.WhStandingOne,
		Brass:      100,
		Silver:     15,
		Gold:       1,
		Spells:     []string{spell0.Id, spell1.Id},
		Sin:        0,
		Corruption: 0,
		Mutations:  []string{},
		Shared:     false,
	},
}

var character1 = warhammer.Wh{
	Id:      "800000000000000000000001",
	OwnerId: user1.Id,
	Object: warhammer.WhCharacter{
		Name:        "character 1",
		Description: fmt.Sprintf("owned by %s", user1.Username),
		Species:     warhammer.WhCharacterSpeciesDwarfAltdorf,
		Career:      career0.Id,
	},
}

func NewMockCharacter() []*warhammer.Wh {
	return []*warhammer.Wh{&character0, &character1}
}
