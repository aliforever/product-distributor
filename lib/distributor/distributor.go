package distributor

import (
	"github.com/aliforever/product-distributor/lib/distributor/dbs"
	"net/http"
)

type Distributor struct {
	repository *dbs.Repository

	router *http.ServeMux

	logger Logger
}

func NewDistributor(repository *dbs.Repository, logger Logger) *Distributor {
	if logger == nil {
		logger = NewDefaultLogger(nil)
	}

	return &Distributor{
		repository: repository,
		router:     http.NewServeMux(),
		logger:     logger,
	}
}

// StartAPI starts distributor API
func (d *Distributor) StartAPI(address string) error {
	d.registerRoutes()

	return http.ListenAndServe(address, d.router)
}

func (d *Distributor) registerRoutes() {
	d.router.HandleFunc("/packages", d.handlePackages)
	d.router.HandleFunc("/packages/add", d.handleAddPackage)
	d.router.HandleFunc("/packages/remove", d.handleRemovePackage)

	d.router.HandleFunc("/orders/add", d.handleOrder)
}
