package mock_data

import "github.com/jmilosze/wfrp-hammergen-go/internal/domain"

type UserSeed struct {
	Id             string
	Username       string
	Password       string
	Admin          bool
	SharedAccounts []string
}

func NewMockUsers() []*UserSeed {
	return []*UserSeed{
		{
			Id:             "000000000000000000000000",
			Username:       "user1@test.com",
			Password:       "123456",
			Admin:          true,
			SharedAccounts: []string{},
		},
		{
			Id:             "000000000000000000000001",
			Username:       "user2@test.com",
			Password:       "789123",
			Admin:          false,
			SharedAccounts: []string{"user1@test.com"},
		},
		{
			Id:             "000000000000000000000002",
			Username:       "user3@test.com",
			Password:       "111111",
			Admin:          false,
			SharedAccounts: []string{"user1@test.com", "user2@test.com"},
		},
	}
}

func NewMockMutations() []*domain.Wh {
	return []*domain.Wh{
		{
			Id:      "100000000000000000000000",
			OwnerId: "000000000000000000000000",
			Object: domain.WhMutation{
				Name:        "mutation 1",
				Description: "owned by user1",
				Type:        0,
				Modifiers: domain.WhModifiers{
					Size:     1,
					Movement: 1,
					Attributes: domain.WHAttributes{
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
				},
				Shared: true,
			},
		},
		{
			Id:      "100000000000000000000001",
			OwnerId: "000000000000000000000001",
			Object: domain.WhMutation{
				Name:        "mutation 2",
				Description: "owned by user2",
				Type:        1,
				Modifiers: domain.WhModifiers{
					Size:     0,
					Movement: 0,
					Attributes: domain.WHAttributes{
						WS:  0,
						BS:  0,
						S:   0,
						T:   0,
						I:   0,
						Ag:  0,
						Dex: 0,
						Int: 0,
						WP:  0,
						Fel: 0,
					},
				},
				Shared: false,
			},
		},
	}
}
