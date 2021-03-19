package controller

import (
	"encoding/json"
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

func (c Controller) sort(companies []*Company) {

	for i := range companies {
		for j := range companies {
			if companies[i].Rank < companies[j].Rank {
				a := *companies[i]
				b := *companies[j]
				*companies[i] = b
				*companies[j] = a
			}
		}
	}
}

func (c Controller) buildRanks(companies []*Company) {
	allCompanies, err := c.r.ZRevRange(keyLeaderBoard, 0, -1)
	if err != nil {
		log.Println(err)
	}
	ranks := make(map[string]int, len(allCompanies))
	for i := range allCompanies {
		key := allCompanies[i].GetKey()
		ranks[key] = i + 1
	}
	for i := range companies {
		if rank, ok := ranks[companies[i].GetKey()]; ok {
			companies[i].Rank = rank
		}
	}
}

func (c Controller) buildCompany(company *Company) {
	key := company.GetKey()
	data, err := c.r.HGetAll(key)
	if err == redis.Nil {
		return
	} else if err != nil {
		log.Println(err)
		return
	}

	company.Symbol = data["symbol"]
	company.Company = data["company"]
	company.Country = data["country"]
}

func (c Controller) buildCompanies(companies []*Company) {
	for i := range companies {
		companies[i].Rank = i + 1
		c.buildCompany(companies[i])
	}
}

func (c Controller) All() ([]*Company, error) {
	companies, err := c.r.ZRevRange(keyLeaderBoard, 0, -1)
	if err != nil {
		return nil, err
	}
	c.buildCompanies(companies)
	return companies, nil
}

func (c Controller) Top10() ([]*Company, error) {
	companies, err := c.r.ZRevRange(keyLeaderBoard, 0, 9)
	if err != nil {
		return nil, err
	}
	c.buildCompanies(companies)
	c.buildRanks(companies)
	return companies, nil
}

func (c Controller) Bottom10() ([]*Company, error) {
	companies, err := c.r.ZRange(keyLeaderBoard, 0, 9)
	if err != nil {
		return nil, err
	}
	c.buildCompanies(companies)
	c.buildRanks(companies)
	return companies, nil
}

func (c Controller) InRank(start, end int64) ([]*Company, error) {
	companies, err := c.r.ZRevRange(keyLeaderBoard, start, end)
	if err != nil {
		return nil, err
	}
	c.buildCompanies(companies)
	c.buildRanks(companies)
	c.sort(companies)
	return companies, nil
}

func (c Controller) GetBySymbol(symbols []string) ([]*Company, error) {
	companies := make([]*Company, 0, len(symbols))

	for i := range symbols {
		score, err := c.r.ZScore(keyLeaderBoard, strings.ToLower(symbols[i]))
		if err != nil {
			return nil, err
		}

		company := &Company{
			Symbol:    symbols[i],
			MarketCap: score,
		}

		c.buildCompany(company)

		companies = append(companies, company)
	}
	c.buildRanks(companies)
	c.sort(companies)
	return companies, nil
}

func (c Controller) UpdateRank(symbol string, amount float64) error  {
	err := c.r.ZIncrBy(keyLeaderBoard, amount, symbol)
	if err != nil {
		return err
	}
	return nil
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
