package models

type Unmarshallable interface {
	DataVillageReq | DataDistrictReq | DataCityReq
}

type Regionable interface {
	Village | District | City
}
