package models

import (
	"math"
	"math/rand"
)

type DirectionV int
type DirectionH int

const (
	North DirectionV = iota
	South
)
const (
	East DirectionH = iota
	West
)

var vdirections = []DirectionV{North, South}
var hdirections = []DirectionH{East, West}

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

func ShiftLoc(l Location, dist float64, dv DirectionV, dh DirectionH) Location {

	LONGITUDE_PER_METRE := func(l float64) float64 {
		return 1 / (40075 * 1000 * RadToDeg(math.Cos(DegToRad(l))) / 360)
	}

	if dv != 0 {
		if dv == North {
			return Location{
				Longitude: l.Longitude,
				Latitude:  l.Latitude + LATITUDE_PER_METRE*dist,
			}
		} else if dv == South {
			return Location{
				Longitude: l.Longitude,
				Latitude:  l.Latitude - LATITUDE_PER_METRE*dist,
			}
		}

	} else {
		if dh == East {
			newLong := l.Longitude + LONGITUDE_PER_METRE(l.Latitude)*dist
			return Location{
				Longitude: newLong,
				Latitude:  l.Latitude,
			}
		} else if dh == West {
			newLong := l.Longitude - LONGITUDE_PER_METRE(l.Latitude)*dist
			return Location{
				Longitude: newLong,
				Latitude:  l.Latitude,
			}
		}

	}
	return Location{}
}

func RandShiftLoc(l Location, dist float64) Location {
	if dist < 0 {
		return Location{}
	}

	dFirst := vdirections[rand.Intn(2)]
	loc := ShiftLoc(l, dist/2, dFirst, 0)
	dSecond := hdirections[rand.Intn(2)]
	loc = ShiftLoc(loc, dist/2, 0, dSecond)

	return loc
}
