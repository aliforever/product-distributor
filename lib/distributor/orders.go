package distributor

import (
	"encoding/json"
	"github.com/aliforever/product-distributor/lib/distributor/entities"
	"github.com/aliforever/product-distributor/lib/distributor/models"
	"math"
	"net/http"
	"sort"
)

func Distribute(items []models.Package, quantity int) []models.Package {
	packages := make([]int, len(items))

	sort.Slice(items, func(i, j int) bool {
		return items[i].Quantity > items[j].Quantity
	})

	for i := range items {
		packages[i] = items[i].Quantity
	}

	result := []models.Package{}

	leastQuantity := packages[len(packages)-1]

	minSize := math.Ceil(float64(quantity)/float64(leastQuantity)) * float64(leastQuantity)

	for i, packSize := range packages {
		packs := math.Floor(minSize / float64(packSize))
		if packs > 0 {
			minSize -= packs * float64(packSize)
			result = append(result, items[i])
		}
	}

	return result
}

func (d *Distributor) handleOrder(writer http.ResponseWriter, request *http.Request) {
	var orderRequest entities.SubmitOrderRequest

	err := json.NewDecoder(request.Body).Decode(&orderRequest)
	if err != nil {
		d.logger.Errorf("error decoding order request: %s", err.Error())
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	packages, err := d.repository.Packages.GetAll()
	if err != nil {
		d.logger.Errorf("error getting packages: %s", err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	var packSizes []int

	for _, p := range packages {
		packSizes = append(packSizes, p.Quantity)
	}

	result := Distribute(packages, orderRequest.Quantity)

	j, err := json.Marshal(entities.NewOkResponse(result))
	if err != nil {
		d.logger.Errorf("error marshaling packages: %s", err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, _ = writer.Write(j)
}
