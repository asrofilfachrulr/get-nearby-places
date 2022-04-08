package models

type Category struct {
	ID    uint8
	Name  string
	Count uint
}

// Predefined categories, mapped by their region level
// ID 1-3 city level
// ID 4-6 district level
// ID 7-9 village level
var MapCategories = map[string][]Category{
	"city": {
		{ID: 1, Name: "Kantor Pemerintah", Count: 1},
		{ID: 2, Name: "Rumah Sakit", Count: 3},
		{ID: 3, Name: "SMA", Count: 20},
	},
	"district": {
		{ID: 4, Name: "Kantor Pemerintah Kecamatan", Count: 1},
		{ID: 5, Name: "Puskesmas", Count: 5},
		{ID: 6, Name: "SMP", Count: 3},
	},
	"village": {
		{ID: 7, Name: "Kantor Pemerintah", Count: 1},
		{ID: 8, Name: "Tempat Ibadah", Count: 20},
		{ID: 9, Name: "SD", Count: 5},
	},
}
