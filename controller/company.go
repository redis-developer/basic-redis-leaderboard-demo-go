package controller

import (
	"fmt"
	"strings"
)

type Company struct {
	Symbol    string  `json:"symbol"`
	Company   string  `json:"company"`
	Country   string  `json:"country"`
	MarketCap float64 `json:"marketCap"`
	Rank      int     `json:"rank"`
}

func (c *Company) GetKey() string {
	return fmt.Sprintf("leaderboard:%s", strings.ToLower(c.Symbol))
}
