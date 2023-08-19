package distributor

import (
	"encoding/json"
	"github.com/aliforever/product-distributor/lib/distributor/entities"
	"net/http"
)

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

	result, err := GraphPackCalculator{PackSizes: packSizes}.Calculate(orderRequest.Quantity)
	if err != nil {
		d.logger.Errorf("error calculating packs: %s", err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	j, err := json.Marshal(entities.NewOkResponse(result))
	if err != nil {
		d.logger.Errorf("error marshaling packages: %s", err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, _ = writer.Write(j)
}
