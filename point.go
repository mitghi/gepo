package gepo

import (
	"fmt"
	"github.com/mitghi/gepo/storage"
	"math"
)

type Point struct {
	ID   string  `json:"id"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
	Meta any
}

type cell struct{ x, y int }

type pointAggregator struct {
	points    []*Point
	point     *Point
	distances distanceMap
}

func NewPoint(id string, lat float64, lon float64) *Point {
	return &Point{id, lat, lon, nil}
}

func (p *Point) String() string {
	return fmt.Sprintf("Point(ID=(%s), Lat=(%f), Lon=(%f))", p.ID, p.Lat, p.Lon)
}

func (p *Point) cell(resolution float64) cell {
	x, y := coordinate(p.Lat, p.Lon, resolution)
	return cell{x, y}
}
func (p *Point) Direction(target *Point) int   { return direction(p.Bearing(target)) }
func (p *Point) Bearing(target *Point) float64 { return bearing(p.Lat, p.Lon, target.Lat, target.Lon) }
func (p *Point) Distance(target *Point) float64 {
	return distance(p.Lat, p.Lon, target.Lat, target.Lon)
}

func (p *pointAggregator) approximateSquareDistance(p1, p2 *Point) float64 {
	avgLat := (p1.Lat + p2.Lat) / 2.0
	latLen := math.Abs(p1.Lat-p2.Lat) * float64(latDegreeLength)
	lonLen := math.Abs(p1.Lon-p2.Lon) * float64(p.distances.get(avgLat))
	return latLen*latLen + lonLen*lonLen
}

func (p *pointAggregator) Len() int      { return len(p.points) }
func (p *pointAggregator) Swap(i, j int) { p.points[i], p.points[j] = p.points[j], p.points[i] }
func (p *pointAggregator) Less(i, j int) bool {
	return p.approximateSquareDistance(p.points[i], p.point) < p.approximateSquareDistance(p.points[j], p.point)
}

func getPoints(stores []storage.Storage[*Point], accept AcceptFn) []*Point {
	return seekPoints(make([]*Point, 0), stores, accept)
}

func seekPoints(points []*Point, stores []storage.Storage[*Point], accept AcceptFn) []*Point {
	for _, s := range stores {
		for _, v := range s.Values() {
			if accept(v) {
				points = append(points, v)
			}
		}
	}
	return points
}
