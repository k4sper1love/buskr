package profile

import (
	"context"

	"github.com/k4sper1love/buskr/internal/domain/user"
	"github.com/k4sper1love/buskr/internal/usecase/keys"
	"github.com/k4sper1love/buskr/internal/usecase/response"
)

func (uc *Usecase) Profile(ctx context.Context, u *user.User) (response.Reply, error) {
	stats, err := uc.users.GetStats(ctx, u.ID)
	if err != nil {
		stats = &user.UserStats{}
	}

	name := u.Username
	if u.Name != "" {
		name = u.Name
	}

	karmaText := keys.TextProfileMainLblKarmaGood
	if uc.users.IsLowKarma(u) {
		karmaText = keys.TextProfileMainLblKarmaBad
	}

	badgeKey := ""
	if u.Role == user.RoleAdmin {
		badgeKey = "profile.main.admin_badge"
	}

	// Map domain values to i18n keys for human-readable display
	noiseKey := "common.lbl.noise_" + string(u.NoiseLevel) // e.g. "common.lbl.noise_light" → "🟢 Light"

	return response.Reply{
		Kind: response.KindEdit,
		Text: response.Text{Key: keys.TextProfileMainTitle, Args: map[string]any{
			"name":                name,
			"noise":               noiseKey,
			"badge":               badgeKey,
			"total_bookings":      stats.TotalBookings,
			"successful_checkins": stats.SuccessfulCheckins,
			"no_shows":            stats.NoShows,
			"karma":               karmaText,
		},
			SubKeyArgs: []string{"karma", "noise", "badge"},
		},
		Keyboard: response.Keyboard{
			InlineRows: [][]response.Button{
				{
					{
						Text: response.Text{Key: keys.TextProfileMainBtnEditName},
						Data: response.CallbackData{Unique: keys.BtnProfileEditName},
					},
				},
				{
					{
						Text: response.Text{Key: keys.TextProfileMainBtnUpgNoise},
						Data: response.CallbackData{Unique: keys.BtnProfileNoiseUp},
					},
				},
				{
					{
						Text: response.Text{Key: keys.TextCommonBtnMenu},
						Data: response.CallbackData{Unique: keys.BtnCommonMenu},
					},
				},
			},
		},
	}, nil
}
