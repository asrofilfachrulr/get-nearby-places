package models

type (
	Category struct {
		ID    uint8
		Name  string
		Count uint
	}
)

var MapCategories = map[string][]Category{
	"city": {
		{ID: 1, Name: "Rumah Sakit", Count: 3},
		{ID: 2, Name: "SMA", Count: 20},
		{ID: 3, Name: "Kantor Pemerintah", Count: 1},
	},
	"district": {
		{ID: 1, Name: "Puskesmas", Count: 5},
		{ID: 2, Name: "SMP", Count: 3},
		{ID: 3, Name: "Kantor Pemerintah Kecamatan", Count: 1},
	},
	"village": {
		{ID: 1, Name: "SD", Count: 5},
		{ID: 2, Name: "Tempat Ibadah", Count: 20},
		{ID: 3, Name: "Kantor Pemerintah", Count: 1},
	},
}
