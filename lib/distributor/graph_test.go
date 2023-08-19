package distributor

import (
	"reflect"
	"testing"
)

func TestGraphPackCalculator_Calculate(t *testing.T) {
	packSizes := []int{250, 500, 1000, 2000, 5000}

	type fields struct {
		PackSizes []int
	}
	type args struct {
		quantity int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    RequiredPacks
		wantErr bool
	}{
		{
			name:    "Test 1",
			fields:  fields{PackSizes: packSizes},
			args:    args{quantity: 1},
			want:    RequiredPacks{250: 1},
			wantErr: false,
		},
		{
			name:    "Test 2",
			fields:  fields{PackSizes: packSizes},
			args:    args{quantity: 250},
			want:    RequiredPacks{250: 1},
			wantErr: false,
		},
		{
			name:    "Test 3",
			fields:  fields{PackSizes: packSizes},
			args:    args{quantity: 251},
			want:    RequiredPacks{500: 1},
			wantErr: false,
		},
		{
			name:    "Test 4",
			fields:  fields{PackSizes: packSizes},
			args:    args{quantity: 501},
			want:    RequiredPacks{500: 1, 250: 1},
			wantErr: false,
		},
		{
			name:   "Test 5",
			fields: fields{PackSizes: packSizes},
			args:   args{quantity: 12001},
			want:   RequiredPacks{5000: 2, 2000: 1, 250: 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := GraphPackCalculator{
				PackSizes: tt.fields.PackSizes,
			}
			got, err := c.Calculate(tt.args.quantity)
			if (err != nil) != tt.wantErr {
				t.Errorf("Calculate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Calculate() got = %v, want %v", got, tt.want)
			}
		})
	}
}
