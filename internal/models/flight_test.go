package models

import (
	"testing"
)

func TestSearchRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     SearchRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: SearchRequest{
				Origin:        "CGK",
				Destination:   "DPS",
				DepartureDate: "2025-12-15",
				Passengers:    1,
				CabinClass:    "economy",
			},
			wantErr: false,
		},
		{
			name: "missing origin",
			req: SearchRequest{
				Destination:   "DPS",
				DepartureDate: "2025-12-15",
				Passengers:    1,
				CabinClass:    "economy",
			},
			wantErr: true,
		},
		{
			name: "missing destination",
			req: SearchRequest{
				Origin:        "CGK",
				DepartureDate: "2025-12-15",
				Passengers:    1,
				CabinClass:    "economy",
			},
			wantErr: true,
		},
		{
			name: "invalid passengers",
			req: SearchRequest{
				Origin:        "CGK",
				Destination:   "DPS",
				DepartureDate: "2025-12-15",
				Passengers:    0,
				CabinClass:    "economy",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}