package fly

import (
	"fmt"
	"github.com/umahmood/haversine"
)

type Region struct {
	Code             string `json:"code"`
	Name             string `json:"name"`
	GatewayAvailable bool   `json:"gatewayAvailable"`
	Location         haversine.Coord
	UserDistance     float64
}

// Distance calculates the distance between region and a given point
func (r *Region) Distance(p haversine.Coord) float64 {
	mi, km := haversine.Distance(p, r.Location)
	fmt.Printf("Miles: %f, Kilometers: %f", mi, km)
	return mi
}

var Regions = []Region{
	{
		Code: "ams",
		Name: "Amsterdam, Netherlands",
		Location: haversine.Coord{
			Lat: 52.3676,
			Lon: 4.9041,
		},
	},

	{
		Code: "cdg",
		Name: "Paris, France",
		Location: haversine.Coord{
			Lat: 48.8566,
			Lon: 2.3522,
		},
	},

	{
		Code: "dfw",
		Name: "Dallas, Texas (US)",
		Location: haversine.Coord{
			Lat: 32.7767,
			Lon: 96.7970,
		},
	},

	{
		Code: "ewr",
		Name: "Secaucus, NJ (US)",
		Location: haversine.Coord{
			Lat: 40.7895,
			Lon: 74.0565,
		},
	},

	{
		Code: "fra",
		Name: "Frankfurt, Germany",
		Location: haversine.Coord{
			Lat: 50.1109,
			Lon: 8.6821,
		},
	},

	{
		Code: "gru",
		Name: "São Paulo",
		Location: haversine.Coord{
			Lat: 23.5558,
			Lon: 46.6396,
		},
	},

	{
		Code: "hkg",
		Name: "Hong Kong, Hong Kong",
		Location: haversine.Coord{
			Lat: 22.3193,
			Lon: 114.1694,
		},
	},

	{
		Code: "iad",
		Name: "Ashburn, Virginia (US)",
		Location: haversine.Coord{
			Lat: 39.0438,
			Lon: 77.4874,
		},
	},

	{
		Code: "lax",
		Name: "Los Angeles, California (US)",
		Location: haversine.Coord{
			Lat: 34.0522,
			Lon: 118.2437,
		},
	},

	{
		Code: "lhr",
		Name: "London, United Kingdom",
		Location: haversine.Coord{
			Lat: 51.5072,
			Lon: 0.1276,
		},
	},

	{
		Code: "maa",
		Name: "Chennai (Madras), India",
		Location: haversine.Coord{
			Lat: 13.0827,
			Lon: 80.2707,
		},
	},

	{
		Code: "mad",
		Name: "Madrid, Spain",
		Location: haversine.Coord{
			Lat: 40.4168,
			Lon: 3.7038,
		},
	},

	{
		Code: "mia",
		Name: "Miami, Florida (US)",
		Location: haversine.Coord{
			Lat: 25.7617,
			Lon: 80.1918,
		},
	},

	{
		Code: "nrt",
		Name: "Tokyo, Japan",
		Location: haversine.Coord{
			Lat: 35.6762,
			Lon: 139.6503,
		},
	},

	{
		Code: "ord",
		Name: "Chicago, Illinois (US)",
		Location: haversine.Coord{
			Lat: 41.8781,
			Lon: 87.6298,
		},
	},

	{
		Code: "scl",
		Name: "Santiago, Chile",
		Location: haversine.Coord{
			Lat: 33.4489,
			Lon: 70.6693,
		},
	},

	{
		Code: "sea",
		Name: "Seattle, Washington (US)",
		Location: haversine.Coord{
			Lat: 47.6062,
			Lon: 122.3321,
		},
	},

	{
		Code: "sin",
		Name: "Singapore",
		Location: haversine.Coord{
			Lat: 1.3521,
			Lon: 103.8198,
		},
	},

	{
		Code: "sjc",
		Name: "Sunnyvale, California (US)",
		Location: haversine.Coord{
			Lat: 37.3688,
			Lon: 122.0363,
		},
	},

	{
		Code: "syd",
		Name: "Sydney, Australia",
		Location: haversine.Coord{
			Lat: 33.8688,
			Lon: 151.2093,
		},
	},

	{
		Code: "yul",
		Name: "Montreal, Canada",
		Location: haversine.Coord{
			Lat: 45.5017,
			Lon: 73.5673,
		},
	},

	{
		Code: "yyz",
		Name: "Toronto, Canada",
		Location: haversine.Coord{
			Lat: 43.6532,
			Lon: 79.3832,
		},
	},
}

func ClosestRegion(coord haversine.Coord) *Region {
	for i, r := range Regions {
		Regions[i].UserDistance = r.Distance(coord)
	}

	shortestDistance := Regions[0].UserDistance
	shortestIndex := 0
	for idx, r := range Regions {
		if shortestDistance > r.UserDistance {
			shortestDistance = r.UserDistance
			shortestIndex = idx
		}
	}

	return &Regions[shortestIndex]
}
