package league

import (
	"context"
	"errors"
	"fmt"
	"os"

	httpc "github.com/milgradesec/go-libs/http"
	"github.com/rs/zerolog/log"
	"github.com/yuhanfang/riot/apiclient"
	"github.com/yuhanfang/riot/constants/region"
	"github.com/yuhanfang/riot/ratelimit"
)

var (
	Client apiclient.Client
)

func NewClient() (apiclient.Client, error) {
	var apikey string
	apikeyFile, found := os.LookupEnv("RIOT_APIKEY_FILE")
	if found {
		buf, err := os.ReadFile(apikeyFile)
		if err != nil {
			return nil, err
		}
		apikey = string(buf)
	} else {
		apikey, found = os.LookupEnv("RIOT_APIKEY")
		if !found {
			return nil, errors.New("RIOT_APIKEY env variable not set")
		}
		log.Warn().Msg("Using unencrypted Riot API Key from env, consider switching to RIOT_APIKEY_FILE")
	}

	return apiclient.New(apikey, httpc.NewHTTPClient(), ratelimit.NewLimiter()), nil
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

/*func (bot *Bot) eloHandler(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	var name string

	if len(args) == 1 {
		name = defaultSummonerName
	} else {
		name = strings.Join(args[1:], " ")
	}

	msg, err := bot.getLeagueElo(name)
	if err != nil {
		log.Error().Err(err).Msgf("failed to get league data for '%s'", name)
		return
	}
	s.ChannelMessageSend(m.ChannelID, msg)
}

func (bot *Bot) getLeagueElo(name string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	summ, err := bot.riotapi.GetBySummonerName(ctx, region.EUW1, name)
	if err != nil {
		return "", err
	}

	list, err := bot.riotapi.GetAllLeaguePositionsForSummoner(ctx, region.EUW1, summ.ID)
	if err != nil {
		return "", err
	}

	for _, league := range list {
		if league.QueueType == "RANKED_SOLO_5x5" {
			wr := float64(league.Wins) / float64(league.Wins+league.Losses) * 100
			return fmt.Sprintf("%s %s %s %d LPs -- %dW/%dL %.2f%% WR\n", name, league.Tier, league.Rank,
				league.LeaguePoints, league.Wins, league.Losses, wr), nil
		}
	}
	return "", nil
}*/

const (
	defaultRegion = region.EUW1
)
