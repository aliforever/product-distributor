package distributor

import (
	"github.com/aliforever/product-distributor/lib/distributor/models"
	"testing"
)

func Test_dist(t *testing.T) {
	packages := []models.Package{
		{
			ID:       "1",
			Quantity: 250,
		},
		{
			ID:       "2",
			Quantity: 500,
		},
		{
			ID:       "3",
			Quantity: 1000,
		},
		{
			ID:       "4",
			Quantity: 2000,
		},
		{
			ID:       "5",
			Quantity: 5000,
		},
	}

	type args struct {
		items    []models.Package
		quantity int
	}
	tests := []struct {
		name string
		args args
		want []models.Package
	}{
		{
			name: "test 1",
			args: args{
				items:    packages,
				quantity: 1,
			},
			want: []models.Package{
				{
					ID:       "1",
					Quantity: 250,
				},
			},
		},
		{
			name: "test 2",
			args: args{
				items:    packages,
				quantity: 250,
			},
			want: []models.Package{
				{
					ID:       "1",
					Quantity: 250,
				},
			},
		},
		{
			name: "test 3",
			args: args{
				items:    packages,
				quantity: 251,
			},
			want: []models.Package{
				{
					ID:       "2",
					Quantity: 500,
				},
			},
		},
		{
			name: "test 4",
			args: args{
				items:    packages,
				quantity: 501,
			},
			want: []models.Package{
				{
					ID:       "1",
					Quantity: 250,
				},
				{
					ID:       "2",
					Quantity: 500,
				},
			},
		},
		{
			name: "test 5",
			args: args{
				items:    packages,
				quantity: 12001,
			},
			want: []models.Package{
				{
					ID:       "1",
					Quantity: 250,
				},
				{
					ID:       "4",
					Quantity: 2000,
				},
				{
					ID:       "5",
					Quantity: 5000,
				},
				{
					ID:       "5",
					Quantity: 5000,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Distribute(tt.args.items, tt.args.quantity); !areSlicesEqualWithoutOrder(got, tt.want) {
				t.Errorf("dist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func areSlicesEqualWithoutOrder(slice1, slice2 []models.Package) bool {
	if len(slice1) != len(slice2) {
		return false
	}

	counts := make(map[string]int)

	// Count occurrences in slice1
	for _, value := range slice1 {
		counts[value.ID]++
	}

	// Subtract occurrences based on slice2
	for _, value := range slice2 {
		counts[value.ID]--
		if counts[value.ID] < 0 {
			return false
		}
	}

	return true
}
