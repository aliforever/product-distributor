package main

import (
	"encoding/json"
	"fmt"
	"github.com/aliforever/product-distributor/dbs"
	entities2 "github.com/aliforever/product-distributor/entities"
	"github.com/aliforever/product-distributor/models"
	"math"
	"net/http"
	"sort"
)

var packages dbs.Packages

func main() {
	packages = dbs.NewInMemoryPackages()

	err := SeedDefaultPackages()
	if err != nil {
		fmt.Println(err)
		return
	}

	router := http.NewServeMux()

	router.HandleFunc("/packages", handlePackages)
	router.HandleFunc("/packages/add", handleAddPackage)
	router.HandleFunc("/packages/remove", handleRemovePackage)

	router.HandleFunc("/orders/add", handleOrder)

	router.Handle("/", http.FileServer(http.Dir("./public")))

	fmt.Println("Listening on port 8080")

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func handlePackages(writer http.ResponseWriter, request *http.Request) {
	data, err := packages.GetAll()
	if err != nil {
		_ = fmt.Errorf("error getting packages: %s", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	j, err := json.Marshal(entities2.NewOkResponse(entities2.PackagesResponse{Packages: data}))
	if err != nil {
		_ = fmt.Errorf("error marshaling packages: %s", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, _ = writer.Write(j)
}

func handleAddPackage(writer http.ResponseWriter, request *http.Request) {
	var addPackageRequest entities2.AddPackageRequest

	err := json.NewDecoder(request.Body).Decode(&addPackageRequest)
	if err != nil {
		_ = fmt.Errorf("error decoding add package request: %s", err.Error())
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	allPackages, err := packages.GetAll()
	if err != nil {
		_ = fmt.Errorf("error getting packages: %s", err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, p := range allPackages {
		if p.Quantity == addPackageRequest.Quantity {
			_ = fmt.Errorf("package with quantity %d already exists", addPackageRequest.Quantity)
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	p := &models.Package{
		ID:       addPackageRequest.ID,
		Quantity: addPackageRequest.Quantity,
	}

	err = packages.InsertPackage(p)
	if err != nil {
		_ = fmt.Errorf("error inserting package: %s", err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	j, err := json.Marshal(entities2.NewOkResponse(entities2.AddPackageResponse{Package: p}))
	if err != nil {
		_ = fmt.Errorf("error marshaling add package response: %s", err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, _ = writer.Write(j)
}

func handleRemovePackage(writer http.ResponseWriter, request *http.Request) {
	var removePackageRequest entities2.RemovePackageRequest

	err := json.NewDecoder(request.Body).Decode(&removePackageRequest)
	if err != nil {
		_ = fmt.Errorf("error decoding remove package request: %s", err.Error())
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	err = packages.RemovePackageByID(removePackageRequest.ID)
	if err != nil {
		_ = fmt.Errorf("error removing package: %s", err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	j, err := json.Marshal(entities2.NewOkResponse(nil))
	if err != nil {
		_ = fmt.Errorf("error marshaling remove package response: %s", err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, _ = writer.Write(j)
}

func handleOrder(writer http.ResponseWriter, request *http.Request) {
	var orderRequest entities2.SubmitOrderRequest

	err := json.NewDecoder(request.Body).Decode(&orderRequest)
	if err != nil {
		_ = fmt.Errorf("error decoding order request: %s", err.Error())
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	packages, err := packages.GetAll()
	if err != nil {
		_ = fmt.Errorf("error getting packages: %s", err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	var packSizes []int

	for _, p := range packages {
		packSizes = append(packSizes, p.Quantity)
	}

	result := Distribute(packages, orderRequest.Quantity)

	j, err := json.Marshal(entities2.NewOkResponse(result))
	if err != nil {
		_ = fmt.Errorf("error marshaling packages: %s", err.Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, _ = writer.Write(j)
}

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

func SeedDefaultPackages() error {
	defaultPackages := []models.Package{
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

	availablePackages, err := packages.GetAll()
	if err != nil {
		return err
	}

	for _, p := range defaultPackages {
		found := false
		for _, ap := range availablePackages {
			if ap.Quantity == p.Quantity {
				found = true
				break
			}
		}

		if !found {
			err = packages.InsertPackage(&p)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
