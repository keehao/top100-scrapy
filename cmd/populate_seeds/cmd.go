package main

import (
	"os"
	"fmt"
	"github.com/romanyx/polluter"
	"github.com/LiamYabou/top100-scrapy/variable"
	"github.com/LiamYabou/top100-pkg/logger"
	"github.com/LiamYabou/top100-scrapy/pkg/automation"
)

func main() {
	// Populate the data into the tables of the database.
	seedPath := fmt.Sprintf("%s/model/data.yml", variable.FixturesURI)
	seed, err := os.Open(seedPath)
	if err != nil {
		logger.Error("Failed to opent the seed.", err)
	}
	defer seed.Close()
	pqConn, err := automation.InitPQconn()
	if err != nil {
		logger.Error("Failed to establish a connection of PQ.", err)
	}
	defer pqConn.Close()
	poluter := polluter.New(polluter.PostgresEngine(pqConn))
	if err := poluter.Pollute(seed); err != nil {
		logger.Error("Failed to pollute the seed.", err)
	}
}
