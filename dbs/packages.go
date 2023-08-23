package dbs

import (
	"github.com/aliforever/product-distributor/models"
	"sync"
)

type Packages interface {
	// GetAll returns all packages
	GetAll() ([]models.Package, error)
	// InsertPackage inserts a package
	InsertPackage(packageModel *models.Package) error
	// RemovePackageByID removes a package by ID
	RemovePackageByID(id string) error
}

// inMemoryPackages implements Packages interface
type inMemoryPackages struct {
	sync.Mutex

	packages []models.Package
}

// NewInMemoryPackages returns a new inMemoryPackages
func NewInMemoryPackages() Packages {
	return &inMemoryPackages{}
}

// GetAll returns all packages
func (imp *inMemoryPackages) GetAll() ([]models.Package, error) {
	imp.Lock()
	defer imp.Unlock()

	return imp.packages, nil
}

// InsertPackage inserts a package
func (imp *inMemoryPackages) InsertPackage(packageModel *models.Package) error {
	imp.Lock()
	defer imp.Unlock()

	imp.packages = append(imp.packages, *packageModel)
	return nil
}

// RemovePackageByID removes a package by ID
func (imp *inMemoryPackages) RemovePackageByID(id string) error {
	imp.Lock()
	defer imp.Unlock()

	for i, p := range imp.packages {
		if p.ID == id {
			imp.packages = append(imp.packages[:i], imp.packages[i+1:]...)
			return nil
		}
	}

	return nil
}
