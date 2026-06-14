package adminuser

import (
	"context"
	"time"

	"github.com/k4sper1love/buskr/internal/domain/user"
	"github.com/k4sper1love/buskr/internal/usecase/keys"
	"github.com/k4sper1love/buskr/internal/usecase/response"
)

func (uc *Usecase) Ban(ctx context.Context, actor *user.User, targetID string) (response.Reply, error) {
	if actor.Role != user.RoleAdmin {
		return response.Reply{}, nil
	}

	target, err := uc.users.GetByID(ctx, targetID)
	if err != nil {
		return response.Reply{
			Kind: response.KindEdit,
			Text: response.Text{Key: keys.TextAdminUsersMsgNotFound},
			Keyboard: response.Keyboard{
				InlineRows: [][]response.Button{
					{{Text: response.Text{Key: keys.TextCommonBtnBack}, Data: response.CallbackData{Unique: keys.BtnAdminUsers}}},
				},
			},
		}, nil
	}

	if target.Role == user.RoleAdmin {
		return response.Reply{
			Kind: response.KindEdit,
			Text: response.Text{Key: keys.TextAdminUsersMsgIsAdmin},
			Keyboard: response.Keyboard{
				InlineRows: [][]response.Button{
					{{Text: response.Text{Key: keys.TextCommonBtnBack}, Data: response.CallbackData{Unique: keys.BtnAdminUsers}}},
				},
			},
		}, nil
	}

	target.Status = user.StatusBanned
	target.UpdatedAt = time.Now()

	err = uc.users.Update(ctx, target)
	if err != nil {
		return response.Reply{
			Kind: response.KindEdit,
			Text: response.Text{Key: keys.TextCommonErrGeneral},
		}, nil
	}

	return uc.UserDetail(ctx, actor, targetID)
}

func (uc *Usecase) Unban(ctx context.Context, actor *user.User, targetID string) (response.Reply, error) {
	if actor.Role != user.RoleAdmin {
		return response.Reply{}, nil
	}

	target, err := uc.users.GetByID(ctx, targetID)
	if err != nil {
		return response.Reply{
			Kind: response.KindEdit,
			Text: response.Text{Key: keys.TextAdminUsersMsgNotFound},
			Keyboard: response.Keyboard{
				InlineRows: [][]response.Button{
					{{Text: response.Text{Key: keys.TextCommonBtnBack}, Data: response.CallbackData{Unique: keys.BtnAdminUsers}}},
				},
			},
		}, nil
	}

	if target.Role == user.RoleAdmin {
		return response.Reply{
			Kind: response.KindEdit,
			Text: response.Text{Key: keys.TextAdminUsersMsgIsAdmin},
			Keyboard: response.Keyboard{
				InlineRows: [][]response.Button{
					{{Text: response.Text{Key: keys.TextCommonBtnBack}, Data: response.CallbackData{Unique: keys.BtnAdminUsers}}},
				},
			},
		}, nil
	}

	target.Status = user.StatusActive
	target.UpdatedAt = time.Now()

	err = uc.users.Update(ctx, target)
	if err != nil {
		return response.Reply{
			Kind: response.KindEdit,
			Text: response.Text{Key: keys.TextCommonErrGeneral},
		}, nil
	}

	return uc.UserDetail(ctx, actor, targetID)
}

func (uc *Usecase) Promote(ctx context.Context, actor *user.User, targetID string) (response.Reply, error) {
	if actor.Role != user.RoleAdmin {
		return response.Reply{}, nil
	}

	target, err := uc.users.GetByID(ctx, targetID)
	if err != nil {
		return response.Reply{
			Kind: response.KindEdit,
			Text: response.Text{Key: keys.TextAdminUsersMsgNotFound},
			Keyboard: response.Keyboard{
				InlineRows: [][]response.Button{
					{{Text: response.Text{Key: keys.TextCommonBtnBack}, Data: response.CallbackData{Unique: keys.BtnAdminUsers}}},
				},
			},
		}, nil
	}

	target.Role = user.RoleAdmin
	target.Status = user.StatusActive
	target.UpdatedAt = time.Now()

	err = uc.users.Update(ctx, target)
	if err != nil {
		return response.Reply{
			Kind: response.KindEdit,
			Text: response.Text{Key: keys.TextCommonErrGeneral},
		}, nil
	}

	return uc.UserDetail(ctx, actor, targetID)
}

func (uc *Usecase) Demote(ctx context.Context, actor *user.User, targetID string) (response.Reply, error) {
	if actor.Role != user.RoleAdmin {
		return response.Reply{}, nil
	}

	if targetID == actor.ID {
		return response.Reply{}, nil
	}

	target, err := uc.users.GetByID(ctx, targetID)
	if err != nil {
		return response.Reply{
			Kind: response.KindEdit,
			Text: response.Text{Key: keys.TextAdminUsersMsgNotFound},
			Keyboard: response.Keyboard{
				InlineRows: [][]response.Button{
					{{Text: response.Text{Key: keys.TextCommonBtnBack}, Data: response.CallbackData{Unique: keys.BtnAdminUsers}}},
				},
			},
		}, nil
	}

	target.Role = user.RoleMusician
	target.UpdatedAt = time.Now()

	err = uc.users.Update(ctx, target)
	if err != nil {
		return response.Reply{
			Kind: response.KindEdit,
			Text: response.Text{Key: keys.TextCommonErrGeneral},
		}, nil
	}

	return uc.UserDetail(ctx, actor, targetID)
}
