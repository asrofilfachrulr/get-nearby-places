package models

import (
	"math"
	"math/rand"
)

type DirectionV int
type DirectionH int

// simplified direction for location shiftting

// vertical shift toward DirectionV constants (south or north) means latitude increased or decreased but longitude remains the same
const (
	North DirectionV = iota
	South
)

// horizontal shift toward DirectionH constants (east or west) means longitude increased or decreased but latitude remains the same
const (
	East DirectionH = iota
	West
)

// grouping constants
var vdirections = []DirectionV{North, South}
var hdirections = []DirectionH{East, West}

// Assuming that the Earth is a sphere with a circumference of 40075 km. So,
// Length in meters of 1° of latitude = always 111.32 km
// Length in meters of 1° of longitude = 40075 km * cos( latitude ) / 360
const (
	LATITUDE_PER_METRE = 1 / (111.32 * 1000)
)

func LONGITUDE_PER_METRE(l float64) float64 {
	return 1 / (40075 * 1000 * RadToDeg(math.Cos(DegToRad(l))) / 360)
}

// converter utilities
func DegToRad(deg float64) float64 {
	return deg * 2 * (math.Phi / 360)
}

func RadToDeg(rad float64) float64 {
	return rad * (360 / 2 * math.Phi)
}

// pure function, returning new object location based given direction and distance
func ShiftLoc(l Location, dist float64, dv DirectionV, dh DirectionH) Location {
	if dv != 0 {
		if dv == North {
			// latitude is increased toward north
			return Location{
				Longitude: l.Longitude,
				Latitude:  l.Latitude + LATITUDE_PER_METRE*dist,
			}
		} else if dv == South {
			// latitude is decreased toward south
			return Location{
				Longitude: l.Longitude,
				Latitude:  l.Latitude - LATITUDE_PER_METRE*dist,
			}
		}

	} else if dh != 0 {
		if dh == East {
			// longitude is increased toward east
			newLong := l.Longitude + LONGITUDE_PER_METRE(l.Latitude)*dist
			return Location{
				Longitude: newLong,
				Latitude:  l.Latitude,
			}
		} else if dh == West {
			// longitude is decreased toward west
			newLong := l.Longitude - LONGITUDE_PER_METRE(l.Latitude)*dist
			return Location{
				Longitude: newLong,
				Latitude:  l.Latitude,
			}
		}

	}
	return Location{}
}

// Perform diagonal shift randomly from given location
func RandShiftLoc(l Location, dist float64) Location {
	if dist < 0 {
		return Location{}
	}

	// shift vertically by one a half of given distance
	dFirst := vdirections[rand.Intn(2)]
	loc := ShiftLoc(l, dist/1.5, dFirst, 0)

	// then shift horizonatallly by one a half of given distance
	dSecond := hdirections[rand.Intn(2)]
	loc = ShiftLoc(loc, dist/1.5, 0, dSecond)

	return loc
}
