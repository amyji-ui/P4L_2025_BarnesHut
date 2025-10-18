package main

import (
	"testing"
	"math"
)

func almostEqual(a, b, eps float64) bool {
	return math.Abs(a-b) <= eps
}

func TestCenterOfMass (t *testing.T) {
	eps := 1e-9

	tests := []struct {
		name   string
		stars  []*Star
		wantX  float64
		wantY  float64
	}{
		{
			name:  "single star",
			stars: []*Star{{position: OrderedPair{x: 3, y: -2}, mass: 5}},
			wantX: 3,
			wantY: -2,
		},
		{
			name: "two stars equal mass",
			stars: []*Star{
				{position: OrderedPair{x: 0, y: 0}, mass: 1},
				{position: OrderedPair{x: 10, y: 10}, mass: 1},
			},
			wantX: 5,
			wantY: 5,
		},
		{
			name: "two stars unequal mass",
			stars: []*Star{
				{position: OrderedPair{x: 0, y: 0}, mass: 1},
				{position: OrderedPair{x: 10, y: 0}, mass: 3}, // heavier star pulls COM toward x=10
			},
			// x = (0*1 + 10*3) / (1+3) = 30/4 = 7.5; y = 0
			wantX: 7.5,
			wantY: 0,
		},
		{
			name: "mixed positions and masses",
			stars: []*Star{
				{position: OrderedPair{x: -2, y: 4}, mass: 2},
				{position: OrderedPair{x: 3, y: -1}, mass: 1},
				{position: OrderedPair{x: 1, y: 5}, mass: 3},
			},
			// SumMass = 6
			// x = (-2*2 + 3*1 + 1*3)/6 = (-4 + 3 + 3)/6 = 2/6 = 0.333333...
			// y = (4*2 + (-1)*1 + 5*3)/6 = (8 - 1 + 15)/6 = 22/6 = 3.666666...
			wantX: 1.0 / 3.0,
			wantY: 22.0 / 6.0,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := CenterOfMass(tc.stars)
			if !almostEqual(got.x, tc.wantX, eps) || !almostEqual(got.y, tc.wantY, eps) {
				t.Fatalf("CenterOfMass() = (%v, %v), want (%v, %v)", got.x, got.y, tc.wantX, tc.wantY)
			}
		})
	}
}