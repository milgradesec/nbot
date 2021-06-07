package bot

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/lus/dgc"
	httpc "github.com/milgradesec/go-libs/http"
	log "github.com/sirupsen/logrus"
	"github.com/yuhanfang/riot/apiclient"
	"github.com/yuhanfang/riot/constants/region"
	"github.com/yuhanfang/riot/ratelimit"
)

func newRiotAPIClient() (apiclient.Client, error) {
	var apikey string
	apikeyFile, found := os.LookupEnv("RIOT_APIKEY_FILE")
	if found {
		buf, err := ioutil.ReadFile(apikeyFile)
		if err != nil {
			return nil, err
		}
		apikey = string(buf)
	} else {
		apikey, found = os.LookupEnv("RIOT_APIKEY")
		if !found {
			return nil, errors.New("RIOT_APIKEY env variable not set")
		}
		log.Warnln("Using unencrypted Riot API Key from env, consider switching to RIOT_APIKEY_FILE")
	}

	return apiclient.New(apikey, httpc.NewHTTPClient(), ratelimit.NewLimiter()), nil
}

func (bot *Bot) eloHandler(ctx *dgc.Ctx) {
	args := ctx.Arguments

	if args.Amount() == 0 {
		msg, err := bot.getLeagueElo("PEIN PACKER")
		if err != nil {
			log.Errorf("error: failed to get league data: %v", err)
			return
		}
		ctx.RespondText(msg) //nolint
	} else {
		name := args.Raw()
		msg, err := bot.getLeagueElo(name)
		if err != nil {
			log.Errorf("error: failed to get league data for '%s': %v", name, err)
			return
		}
		ctx.RespondText(msg) //nolint
	}
}

func (bot *Bot) getLeagueElo(name string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
}
