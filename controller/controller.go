package controller

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"os"
	"strings"
)

const (
	keyPrefixCompany = "leaderboard:"
	keyLeaderBoard   = "REDIS_LEADERBOARD"
)

type Controller struct {
	r Redis
}

func (c Controller) All() ([]*Company, error) {
	list, err := c.r.ZRevRange(keyLeaderBoard, 0, -1)
	if err != nil {
		return nil, err
	}
	log.Println(list)
	companies := make([]*Company, 0)
	rank := 0
	for symbol := range list {
		key := fmt.Sprintf("%s%s", keyPrefixCompany, symbol)
		data, err := c.r.HGetAll(key)
		if err == redis.Nil {
			continue
		} else if err != nil {
			return nil, err
		}
		rank++
		companies = append(companies, &Company{
			Symbol:    data["symbol"],
			Company:   data["company"],
			Country:   data["country"],
			MarketCap: list[symbol],
			Rank:      rank,
		})

	}

	return companies, nil
}

func (c Controller) Top10() ([]*Company, error) {
	return c.All()
}

func (c Controller) Bottom10() ([]*Company, error) {
	return c.All()
}

func (c Controller) InRank(start, end int) ([]*Company, error) {
	return c.All()
}

func (c Controller) GetBySymbol(symbols []string) ([]*Company, error) {
	return c.All()
}

var controller = &Controller{}

func Instance() *Controller {
	return controller
}

func ImportCompanies(filePath string, r Redis) error {
	fp, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		return err
	}
	defer fp.Close()

	companies := make([]Company, 0)
	dec := json.NewDecoder(fp)
	err = dec.Decode(&companies)
	if err != nil {
		return err
	}

	for i := range companies {

		if err := r.ZAdd(keyLeaderBoard, strings.ToLower(companies[i].Symbol), companies[i].MarketCap); err != nil {
			return err
		}

		if err := r.HSet(companies[i].GetKey(), "symbol", strings.ToLower(companies[i].Symbol)); err != nil {
			return err
		}
		if err := r.HSet(companies[i].GetKey(), "company", companies[i].Company); err != nil {
			return err
		}
		if err := r.HSet(companies[i].GetKey(), "country", companies[i].Country); err != nil {
			return err
		}
	}

	return err
}

func SetRedis(redis Redis) {
	controller.r = redis
}
