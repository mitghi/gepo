package gepo

import (
	"github.com/mitghi/gepo/storage"
	"math"
	"sort"
)

type Map struct {
	tiles     *grid
	positions map[string]*Point
	distances distanceMap
}

func New(resolution float64) *Map {
	builder := func() storage.Storage[*Point] {
		return storage.New[*Point]()
	}
	return WithBuilder(resolution, builder)
}

func WithBuilder(resolution float64, builder storage.Builder[*Point]) *Map {
	return &Map{
		tiles:     newGrid(resolution, builder),
		positions: make(map[string]*Point),
		distances: make(distanceMap),
	}
}

func (gm *Map) Add(point *Point) {
	gm.Remove(point.ID)
	tile := gm.tiles.AddEntry(point)
	tile.Add(point.ID, point)
	gm.positions[point.ID] = point
}

func (gm *Map) AddPoints(points ...*Point) {
	for _, p := range points {
		gm.Remove(p.ID)
		tile := gm.tiles.AddEntry(p)
		tile.Add(p.ID, p)
		gm.positions[p.ID] = p
	}
}

func (gm *Map) Get(id string) (p *Point) {
	if p, ok := gm.positions[id]; ok {
		if r, ok := gm.tiles.GetEntry(p).Get(p.ID); ok {
			return r
		}
	}

	return nil
}

func (gm *Map) Remove(id string) {
	if v, ok := gm.positions[id]; ok {
		s := gm.tiles.GetEntry(v)
		s.Delete(v.ID)
		delete(gm.positions, v.ID)
	}
}

func (gm *Map) All() (all map[string]Point) {
	all = make(map[string]Point, 0)
	for k, v := range gm.positions {
		all[k] = *v
	}
	return all
}

func (gm *Map) Range(tl *Point, br *Point) []*Point {
	stores := gm.tiles.Range(tl, br)
	fn := func(p *Point) bool {
		return isBetween(p.Lat, br.Lat, tl.Lat) && isBetween(p.Lon, tl.Lon, br.Lon)
	}
	return getPoints(stores, fn)
}

func (gm *Map) Nearest(point *Point, threshold int, maxDistance float64, accept AcceptFn) []*Point {
	entry := gm.tiles.GetEntry(point)
	points := getPoints([]storage.Storage[*Point]{entry}, accept)
	output := append([]*Point{}, points...)
	index := point.cell(gm.tiles.resolution)
	coarseMaxDis := math.Max(float64(maxDistance)*2.0, float64(gm.tiles.resolution)*2.0+0.01)
	count := 0

	for d := 1; float64(d)*float64(gm.tiles.resolution) <= coarseMaxDis; d++ {
		currentCount := len(output)
		output = seekPoints(output, gm.tiles.get(index.x-d, index.x+d, index.y+d, index.y+d), accept)
		output = seekPoints(output, gm.tiles.get(index.x-d, index.x+d, index.y-d, index.y-d), accept)
		output = seekPoints(output, gm.tiles.get(index.x-d, index.x-d, index.y-d+1, index.y+d-1), accept)
		output = seekPoints(output, gm.tiles.get(index.x+d, index.x+d, index.y-d+1, index.y+d-1), accept)
		count += len(output) - currentCount
		if count > threshold {
			break
		}
	}

	psorted := pointAggregator{output, point, gm.distances}
	sort.Sort(&psorted)

	threshold = min(threshold, len(psorted.points))
	for i, nearby := range psorted.points {
		if point.Distance(nearby) > maxDistance || i == threshold {
			threshold = i
			break
		}
	}

	return psorted.points[0:threshold]
}
