package gepo

import "testing"

func TestNearest(t *testing.T) {
	include := func(input ...string) map[string]struct{} {
		output := make(map[string]struct{})
		for _, v := range input {
			output[v] = struct{}{}
		}
		return output
	}

	points := []*Point{
		NewPoint("a", 52.522973, 13.404598),
		NewPoint("b", 52.528383, 13.387744),
		NewPoint("luxstr", 52.543751, 13.353576),
		NewPoint("wd", 52.552046, 13.361131),
		NewPoint("gb", 52.558554, 13.371625),
	}

	tests := []struct {
		given         *Point
		radius        float64
		shouldInclude map[string]struct{}
	}{
		{
			given:         points[0],
			radius:        Km(1),
			shouldInclude: include("a"),
		},
		{
			given:         points[0],
			radius:        Km(2),
			shouldInclude: include("a", "b"),
		},
		{
			given:         points[2],
			radius:        Km(2),
			shouldInclude: include("luxstr", "wd"),
		},
		{
			given:         points[3],
			radius:        Km(2),
			shouldInclude: include("luxstr", "wd", "gb"),
		},
	}

	gm := New(Km(100.0))
	gm.AddPoints(points...)

	for _, tc := range tests {
		result := gm.Nearest(tc.given, 1000, tc.radius, func(_ *Point) bool { return true })
		lenResult := len(result)
		lenExpect := len(tc.shouldInclude)
		t.Log(result)
		if lenResult != lenExpect {
			t.Fatalf("invalid length, expected %d, got %d", lenExpect, lenResult)
		}
		for _, p := range result {
			_, ok := tc.shouldInclude[p.ID]
			if !ok {
				t.Fatalf("invalid point, %s not found", p.ID)
			}
		}
	}
}
