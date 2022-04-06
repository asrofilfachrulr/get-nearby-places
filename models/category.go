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
		{ID: 2, Name: "Sekolah Menengah Atas", Count: 20},
		{ID: 3, Name: "Kantor Pemerintah Kota", Count: 1},
		{ID: 4, Name: "Kantor Pemerintah Kabupaten", Count: 1},
	},
	"district": {
		{ID: 1, Name: "Puskesmas", Count: 5},
		{ID: 2, Name: "Sekolah Menengah Pertama", Count: 3},
		{ID: 3, Name: "Kantor Pemerintah Kecamatan", Count: 1},
	},
	"village": {
		{ID: 1, Name: "Sekolah Dasar", Count: 5},
		{ID: 2, Name: "Tempat Ibadah", Count: 20},
		{ID: 3, Name: "Kantor Pemerintah Desa", Count: 1},
		{ID: 4, Name: "Kantor Pemerintah Kelurahan", Count: 1},
	},
}
