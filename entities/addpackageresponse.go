package entities

import (
	"github.com/aliforever/product-distributor/models"
)

type AddPackageResponse struct {
	Package *models.Package `json:"package,omitempty"`
}
