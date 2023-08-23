package entities

import (
	"github.com/aliforever/product-distributor/models"
)

type PackagesResponse struct {
	Packages []models.Package `json:"packages,omitempty"`
}
