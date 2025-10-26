package main

import (
	"math"
	"testing"
	"strconv"
)

// minimal, stable 2-body system so something moves each step
func smallUniverse() *Universe {
	u := &Universe{width: 1e12}
	heavy := &Star{
		position:     OrderedPair{x: 5e11, y: 5e11},
		velocity:     OrderedPair{x: 0, y: 0},
		acceleration: OrderedPair{},
		mass:         solarMass,
		radius:       7e8,
		red:          255, green: 220, blue: 0,
	}
	light := &Star{
		position:     OrderedPair{x: 5e11 + 2e10, y: 5e11},
		velocity:     OrderedPair{x: 0, y: 800}, // small tangential speed
		acceleration: OrderedPair{},
		mass:         1e24,
		radius:       3e6,
		red:          200, green: 200, blue: 255,
	}
	u.stars = []*Star{heavy, light}
	return u
}

func finite(x float64) bool { return !math.IsNaN(x) && !math.IsInf(x, 0) }

func TestBarnesHut_MultipleNumGens(t *testing.T) {
	cases := []struct {
		numGens int
	}{
		{3}, {4}, {5}, {6},
	}

	const (
		dt    = 1e3   // modest step to keep it stable but moving
		theta = 0.6
	)

	for _, tc := range cases {
		t.Run(
			// subtest name
			func() string { return "numGens=" + strconv.Itoa(tc.numGens) }(),
			func(t *testing.T) {
				initial := smallUniverse()
				tps := BarnesHut(initial, tc.numGens, dt, theta)

				// length
				if got, want := len(tps), tc.numGens+1; got != want {
					t.Fatalf("len(timepoints)=%d, want %d", got, want)
				}

				// basic invariants per frame
				for i := 0; i < len(tps); i++ {
					if tps[i] == nil {
						t.Fatalf("timePoints[%d] is nil", i)
					}
					if tps[i].width != initial.width {
						t.Fatalf("width changed at %d: got %g want %g", i, tps[i].width, initial.width)
					}
					for j, s := range tps[i].stars {
						if s == nil {
							t.Fatalf("nil star at frame %d index %d", i, j)
						}
						if !finite(s.position.x) || !finite(s.position.y) ||
							!finite(s.velocity.x) || !finite(s.velocity.y) {
							t.Fatalf("non-finite state at frame %d index %d", i, j)
						}
					}
				}

				// no aliasing and some motion between consecutive frames
				for i := 1; i < len(tps); i++ {
					if tps[i] == tps[i-1] {
						t.Fatalf("timePoints[%d] aliases timePoints[%d]", i, i-1)
					}
					if len(tps[i].stars) != len(tps[i-1].stars) {
						t.Fatalf("star count changed between %d and %d", i-1, i)
					}

					moved := false
					for k := range tps[i].stars {
						a := tps[i-1].stars[k]
						b := tps[i].stars[k]
						dx := b.position.x - a.position.x
						dy := b.position.y - a.position.y
						if math.Hypot(dx, dy) > 0 { // any star moved
							moved = true
							break
						}
					}
					if !moved {
						t.Fatalf("no star moved between frames %d and %d", i-1, i)
					}
				}
			},
		)
	}
}
