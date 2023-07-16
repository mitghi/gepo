package gepo

import "github.com/mitghi/gepo/storage"

type grid struct {
	resolution float64
	index      map[cell]storage.Storage[*Point]
	builder    storage.Builder[*Point]
}

func newGrid(resolution float64, builder storage.Builder[*Point]) *grid {
	return &grid{
		resolution: resolution,
		index:      make(map[cell]storage.Storage[*Point]),
		builder:    builder,
	}
}

func (g *grid) Clone() (output *grid) {
	output = newGrid(g.resolution, g.builder)
	for k, v := range g.index {
		output.index[k] = v.Clone()
	}
	return output
}

func (g *grid) AddEntry(point *Point) storage.Storage[*Point] {
	index := point.cell(g.resolution)
	v, ok := g.index[index]
	if !ok {
		tile := g.builder()
		g.index[index] = tile
		return tile
	}
	return v
}

func (g *grid) GetEntry(point *Point) (tile storage.Storage[*Point]) {
	index := point.cell(g.resolution)
	tile, ok := g.index[index]
	if !ok {
		return g.builder()
	}
	return tile
}

func (g *grid) Range(topLeft *Point, bottomRight *Point) []storage.Storage[*Point] {
	tlIndex := topLeft.cell(g.resolution)
	brIndex := bottomRight.cell(g.resolution)
	return g.get(brIndex.x, tlIndex.x, tlIndex.y, brIndex.y)
}

func (g *grid) Cells(minx int, maxx int, miny int, maxy int) (cells []cell) {
	cells = make([]cell, 0)
	for x := minx; x <= maxx; x++ {
		for y := miny; y <= maxy; y++ {
			cells = append(cells, cell{x, y})
		}
	}
	return cells
}

func (g *grid) get(minx int, maxx int, miny int, maxy int) (entries []storage.Storage[*Point]) {
	entries = make([]storage.Storage[*Point], 0)
	for x := minx; x <= maxx; x++ {
		for y := miny; y <= maxy; y++ {
			if index, ok := g.index[cell{x, y}]; ok {
				entries = append(entries, index)
			}
		}
	}
	return entries
}
