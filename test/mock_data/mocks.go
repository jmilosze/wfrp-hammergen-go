package mock_data

import (
	"context"
	"github.com/jmilosze/wfrp-hammergen-go/internal/domain/user"
	"github.com/jmilosze/wfrp-hammergen-go/internal/domain/warhammer"
)

func InitUser(ctx context.Context, s user.UserService) {
	s.SeedUsers(ctx, NewMockUsers())

}

func InitWh(ctx context.Context, s warhammer.WhService) {
	s.SeedWh(ctx, warhammer.WhTypeMutation, NewMockMutations())
	s.SeedWh(ctx, warhammer.WhTypeSpell, NewMockSpells())
	s.SeedWh(ctx, warhammer.WhTypeProperty, NewMockProperties())
	s.SeedWh(ctx, warhammer.WhTypeItem, NewMockItems())
	s.SeedWh(ctx, warhammer.WhTypeTalent, NewMockTalents())
	s.SeedWh(ctx, warhammer.WhTypeSkill, NewMockSkills())
}
