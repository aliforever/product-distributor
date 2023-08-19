package distributor

import (
	"encoding/json"
	"github.com/aliforever/product-distributor/lib/distributor/entities"
	"github.com/aliforever/product-distributor/lib/distributor/models"
	"net/http"
)

func (d *Distributor) handlePackages(writer http.ResponseWriter, request *http.Request) {
	data, err := d.repository.Packages.GetAll()
	if err != nil {
		d.logger.Errorf("error getting packages: %s", err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	j, err := json.Marshal(entities.NewOkResponse(entities.PackagesResponse{Packages: data}))
	if err != nil {
		d.logger.Errorf("error marshaling packages: %s", err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, _ = writer.Write(j)
}

func (d *Distributor) handleAddPackage(writer http.ResponseWriter, request *http.Request) {
	var addPackageRequest entities.AddPackageRequest

	err := json.NewDecoder(request.Body).Decode(&addPackageRequest)
	if err != nil {
		d.logger.Errorf("error decoding add package request: %s", err.Error())
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	packages, err := d.repository.Packages.GetAll()
	if err != nil {
		d.logger.Errorf("error getting packages: %s", err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, p := range packages {
		if p.Quantity == addPackageRequest.Quantity {
			d.logger.Errorf("package with quantity %d already exists", addPackageRequest.Quantity)
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	p := &models.Package{
		ID:       addPackageRequest.ID,
		Quantity: addPackageRequest.Quantity,
	}

	err = d.repository.Packages.InsertPackage(p)
	if err != nil {
		d.logger.Errorf("error inserting package: %s", err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	j, err := json.Marshal(entities.NewOkResponse(entities.AddPackageResponse{Package: p}))
	if err != nil {
		d.logger.Errorf("error marshaling add package response: %s", err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, _ = writer.Write(j)
}

func (d *Distributor) handleRemovePackage(writer http.ResponseWriter, request *http.Request) {
	var removePackageRequest entities.RemovePackageRequest

	err := json.NewDecoder(request.Body).Decode(&removePackageRequest)
	if err != nil {
		d.logger.Errorf("error decoding remove package request: %s", err.Error())
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err = d.repository.Packages.RemovePackageByID(removePackageRequest.ID)
	if err != nil {
		d.logger.Errorf("error removing package: %s", err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	j, err := json.Marshal(entities.NewOkResponse(nil))
	if err != nil {
		d.logger.Errorf("error marshaling remove package response: %s", err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, _ = writer.Write(j)
}
