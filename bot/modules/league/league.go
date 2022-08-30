package league

import (
	"context"
	"errors"
	"fmt"

	httpc "github.com/milgradesec/go-libs/http"
	"github.com/spf13/viper"
	"github.com/yuhanfang/riot/apiclient"
	"github.com/yuhanfang/riot/constants/region"
	"github.com/yuhanfang/riot/ratelimit"
)

const (
	defaultRegion = region.EUW1
)

var (
	Client apiclient.Client
)

func NewClient() (apiclient.Client, error) {
	if !viper.IsSet("RIOT_API_KEY") {
		return nil, errors.New("RIOT_API_KEY not set")
	}

	return apiclient.New(
		viper.GetString("RIOT_API_KEY"),
		httpc.NewHTTPClient(),
		ratelimit.NewLimiter()), nil
}

func GetRankedSummary(ctx context.Context, name string) (string, error) {
	summ, err := Client.GetBySummonerName(ctx, defaultRegion, name)
	if err != nil {
		return "", err
	}

	leagues, err := Client.GetAllLeaguePositionsForSummoner(ctx, defaultRegion, summ.ID)
	if err != nil {
		return "", err
	}

	for _, league := range leagues {
		if league.QueueType == "RANKED_SOLO_5x5" {
			wr := float64(league.Wins) / float64(league.Wins+league.Losses) * 100
			return fmt.Sprintf("%s %s %s %d LPs -- %dW/%dL %.2f%% WR\n", name, league.Tier, league.Rank,
				league.LeaguePoints, league.Wins, league.Losses, wr), nil
		}
	}
	return "", nil
}
