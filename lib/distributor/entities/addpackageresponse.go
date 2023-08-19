package entities

import "github.com/aliforever/product-distributor/lib/distributor/models"

type AddPackageResponse struct {
	Package *models.Package `json:"package,omitempty"`
}
