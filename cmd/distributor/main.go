package main

import (
	"github.com/aliforever/product-distributor/lib/distributor"
	"github.com/aliforever/product-distributor/lib/distributor/dbs"
	"os"
)

func main() {
	logger := distributor.NewDefaultLogger(os.Stdout)

	repository := dbs.NewRepository()

	err := repository.SeedPackages()
	if err != nil {
		logger.Error(err)
		return
	}

	logger.Error(distributor.NewDistributor(repository, logger).StartAPI(":8080"))
}
