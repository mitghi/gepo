package gepo

import "math"

// see: https://manoa.hawaii.edu/exploringourfluidearth/physical/world-ocean/locating-points-globe
const (
	minLon          float64 = (-180.0 * 1000)
	minLat          float64 = (-90.0 * 1000)
	latDegreeLength float64 = (111.0 * 1000)
	lonDegreeLength float64 = (85.0 * 1000)
	cEarthRadius    float64 = (6371.0 * 1000)
)

const (
	NorthEast int = iota
	East
	SouthEast
	South
	SouthWest
	West
	NorthWest
	North
)

// AcceptFn gets called when a point
// matches on the plane and this function
// can further decide if the match should
// be included in result
type AcceptFn func(*Point) bool

func Km(km float64) float64                                  { return km * 1000 }
func toDegrees(x float64) float64                            { return x * 180.0 / math.Pi }
func toRadians(x float64) float64                            { return x * math.Pi / 180.0 }
func isBetween(value float64, min float64, max float64) bool { return (value >= min) && (value <= max) }
func min(a int, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func coordinate(lat float64, lon float64, resolution float64) (x int, y int) {
	x = int((-minLat + lat) * latDegreeLength / resolution)
	y = int((-minLon + lon) * lonDegreeLength / resolution)
	return x, y
}

func direction(bearing float64) int {
	index := bearing - 22.5
	if index < 0 {
		index += 360
	}
	return int(index / 45.0)
}

func bearing(lat1, lon1, lat2, lon2 float64) float64 {
	disLon := toRadians(lon2 - lon1)
	rlat1 := toRadians(lat1)
	rlat2 := toRadians(lat2)
	x := math.Cos(rlat1)*math.Sin(rlat2) - math.Sin(rlat1)*math.Cos(rlat2)*math.Cos(disLon)
	y := math.Sin(disLon) * math.Cos(rlat2)
	return toDegrees(math.Atan2(y, x))
}

func distance(lat1, lon1, lat2, lon2 float64) float64 {
	sLat := math.Sin(toRadians(lat2-lat1) / 2)
	sLon := math.Sin(toRadians(lon2-lon1) / 2)
	a := math.Pow(sLat, 2) + math.Pow(sLon, 2)*math.Cos(toRadians(lat1)*math.Cos(toRadians(lat2)))
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return cEarthRadius * c
}
