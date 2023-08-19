package entities

import "github.com/aliforever/product-distributor/lib/distributor/models"

type PackagesResponse struct {
	Packages []models.Package `json:"packages,omitempty"`
}
