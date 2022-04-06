package models

import (
	"fmt"
	"math"
)

type Direction int

const (
	North Direction = iota
	East
	South
	West
)

// Given you're looking for a simple formula, this is probably the simplest way to do it, assuming that the Earth is a sphere with a circumference of 40075 km.

// Length in meters of 1° of latitude = always 111.32 km

// Length in meters of 1° of longitude = 40075 km * cos( latitude ) / 360

const (
	LATITUDE_PER_METRE = 1 / (111.32 * 1000)
)

func DegToRad(deg float64) float64 {
	return deg * (math.Phi / 360)
}

func RadToDeg(rad float64) float64 {
	return rad * (360 / 2 * math.Phi)
}

func ShiftLocation(l Location, dist float64, d Direction) (Location, error) {
	if dist < 0 {
		return Location{}, fmt.Errorf("distance can't be negative")
	}

	LONGITUDE_PER_METRE := func(l float64) float64 {
		return 1 / (40075 * 1000 * RadToDeg(math.Cos(DegToRad(l))) / 360)
	}

	if d == North {
		return Location{
			Longitude: l.Longitude,
			Latitude:  l.Latitude + LATITUDE_PER_METRE*dist,
		}, nil
	} else if d == South {
		return Location{
			Longitude: l.Longitude,
			Latitude:  l.Latitude - LATITUDE_PER_METRE*dist,
		}, nil
	} else if d == East {
		newLong := l.Longitude + LONGITUDE_PER_METRE(l.Latitude)*dist
		return Location{
			Longitude: newLong,
			Latitude:  l.Latitude,
		}, nil
	} else if d == West {
		newLong := l.Longitude - LONGITUDE_PER_METRE(l.Latitude)*dist
		return Location{
			Longitude: newLong,
			Latitude:  l.Latitude,
		}, nil
	} else {
		return Location{}, fmt.Errorf("direction unknown")
	}
}
