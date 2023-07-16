package gepo

import "math"

type distanceMap map[int]float64

func (dm distanceMap) euclid(p1, p2 Point) float64 {
	latLen := math.Abs(p1.Lat-p2.Lat) * latDegreeLength
	lonLen := math.Abs(p1.Lon-p2.Lon) * dm.averageLat(p1.Lat, p2.Lat)
	return (latLen*latLen + lonLen*lonLen)
}

func (dm distanceMap) get(lat float64) float64 {
	latIndex := int(lat * 10)
	if v, ok := dm[latIndex]; ok {
		return v
	}

	latRnd := float64(latIndex) / 10
	dist := distance(latRnd, 0.0, latRnd, 1.0)
	dm[latIndex] = dist

	return dist
}

func (dm distanceMap) averageLat(lat1, lat2 float64) float64 { return dm.get((lat1 + lat2) / 2.0) }
