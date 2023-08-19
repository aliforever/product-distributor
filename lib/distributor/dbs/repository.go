package dbs

import (
	"github.com/aliforever/product-distributor/lib/distributor/models"
)

type Repository struct {
	Packages Packages
}

func NewRepository() *Repository {
	return &Repository{
		Packages: NewInMemoryPackages(),
	}
}

// SeedPackages seeds packages
func (r *Repository) SeedPackages() error {
	packages := []models.Package{
		{
			ID:       "P_1",
			Quantity: 250,
		},
		{
			ID:       "P_2",
			Quantity: 500,
		},
		{
			ID:       "P_3",
			Quantity: 1000,
		},
		{
			ID:       "P_4",
			Quantity: 2000,
		},
		{
			ID:       "P_5",
			Quantity: 5000,
		},
	}

	availablePackages, err := r.Packages.GetAll()
	if err != nil {
		return err
	}

	for _, p := range packages {
		found := false
		for _, ap := range availablePackages {
			if ap.Quantity == p.Quantity {
				found = true
				break
			}
		}

		if !found {
			err = r.Packages.InsertPackage(&p)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
