package bot

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/yuhanfang/riot/apiclient"
	"github.com/yuhanfang/riot/constants/region"
	"github.com/yuhanfang/riot/ratelimit"
)

func (bot *Bot) getLeagueElo(name string) (string, error) {
	var apikey string
	apikeyFile, found := os.LookupEnv("RIOT_APIKEY_FILE")
	if found {
		buf, err := ioutil.ReadFile(apikeyFile)
		if err != nil {
			return "", err
		}
		apikey = string(buf)
	} else {
		apikey, found = os.LookupEnv("RIOT_APIKEY")
		if !found {
			return "", errors.New("RIOT_APIKEY env variable not set")
		}
		log.Warnln("Using unencrypted Riot API Key from ENV, consider switching to RIOT_APIKEY_FILE")
	}

	client := apiclient.New(apikey, http.DefaultClient, ratelimit.NewLimiter())

	summ, err := client.GetBySummonerName(context.TODO(), region.EUW1, name)
	if err != nil {
		return "", err
	}

	list, err := client.GetAllLeaguePositionsForSummoner(context.TODO(), region.EUW1, summ.ID)
	if err != nil {
		return "", err
	}

	for _, league := range list {
		if league.QueueType == "RANKED_SOLO_5x5" {
			return fmt.Sprintf("%s %s %s %d LPs\n", name, league.Tier, league.Rank, league.LeaguePoints), nil
		}
	}
	return "", nil
}
