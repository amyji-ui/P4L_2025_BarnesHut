package main

// parser_universe.go
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// ParseUniverseFromFile reads your width/G/blocks format:
//   line 1: <width>
//   line 2: <G>         (ignored; code uses const G)
//   then for each body:
//     >Name
//     R,G,B
//     mass
//     radius
//     x,y
//     vx,vy
func ParseUniverseFromFile(path string) (*Universe, error) {
	f, err := os.Open(path)
	if err != nil { return nil, err }
	defer f.Close()
	return parseUniverse(bufio.NewScanner(f))
}

func parseUniverse(sc *bufio.Scanner) (*Universe, error) {
	next := func() (string, error) {
		for sc.Scan() {
			s := strings.TrimSpace(sc.Text())
			if s == "" || strings.HasPrefix(s, "#") { // allow blank/comments
				continue
			}
			return s, nil
		}
		if err := sc.Err(); err != nil { return "", err }
		return "", io.EOF
	}

	// width
	wstr, err := next()
	if err != nil { return nil, fmt.Errorf("width: %w", err) }
	width, err := strconv.ParseFloat(wstr, 64)
	if err != nil { return nil, fmt.Errorf("parse width %q: %w", wstr, err) }

	// G (ignored)
	if _, err := next(); err != nil {
		return nil, fmt.Errorf("G line: %w", err)
	}

	var stars []*Star
	for {
		h, err := next()
		if err == io.EOF { break }
		if err != nil { return nil, err }
		if !strings.HasPrefix(h, ">") {
			return nil, fmt.Errorf("expected '>' header, got %q", h)
		}

		rgb, err := next()
		if err != nil { return nil, err }
		R, Gc, B, err := parseRGB(rgb)
		if err != nil { return nil, fmt.Errorf("rgb %q: %w", rgb, err) }

		massLine, err := next()
		if err != nil { return nil, err }
		mass, err := strconv.ParseFloat(massLine, 64)
		if err != nil { return nil, fmt.Errorf("mass %q: %w", massLine, err) }

		radLine, err := next()
		if err != nil { return nil, err }
		radius, err := strconv.ParseFloat(radLine, 64)
		if err != nil { return nil, fmt.Errorf("radius %q: %w", radLine, err) }

		posLine, err := next()
		if err != nil { return nil, err }
		px, py, err := parsePair(posLine)
		if err != nil { return nil, fmt.Errorf("pos %q: %w", posLine, err) }

		velLine, err := next()
		if err != nil { return nil, err }
		vx, vy, err := parsePair(velLine)
		if err != nil { return nil, fmt.Errorf("vel %q: %w", velLine, err) }

		stars = append(stars, &Star{
			position:     OrderedPair{px, py},
			velocity:     OrderedPair{vx, vy},
			acceleration: OrderedPair{},
			mass:         mass,
			radius:       radius,
			// file: R,G,B  ; struct fields: red, blue, green
			red:   uint8(R),
			blue:  uint8(B),
			green: uint8(Gc),
		})
	}

	return &Universe{stars: stars, width: width}, nil
}

func parseRGB(s string) (int, int, int, error) {
	parts := splitCSV(s)
	if len(parts) != 3 { return 0,0,0, fmt.Errorf("want 3, got %d", len(parts)) }
	r, err := strconv.Atoi(parts[0]); if err != nil { return 0,0,0, err }
	g, err := strconv.Atoi(parts[1]); if err != nil { return 0,0,0, err }
	b, err := strconv.Atoi(parts[2]); if err != nil { return 0,0,0, err }
	return r, g, b, nil
}

func parsePair(s string) (float64, float64, error) {
	parts := splitCSV(s)
	if len(parts) != 2 { return 0,0, fmt.Errorf("want 2, got %d", len(parts)) }
	x, err := strconv.ParseFloat(parts[0], 64); if err != nil { return 0,0, err }
	y, err := strconv.ParseFloat(parts[1], 64); if err != nil { return 0,0, err }
	return x, y, nil
}

func splitCSV(s string) []string {
	raw := strings.Split(s, ",")
	out := make([]string, 0, len(raw))
	for _, p := range raw {
		out = append(out, strings.TrimSpace(p))
	}
	return out
}


type pair = OrderedPair

func TestCalculateNetForce_IO(t *testing.T) {
	base := "Tests/CalculateNetForce"
	inDir := filepath.Join(base, "input")
	outDir := filepath.Join(base, "output")

	inputs, err := filepath.Glob(filepath.Join(inDir, "*.txt"))
	if err != nil {
		t.Fatalf("glob: %v", err)
	}
	if len(inputs) == 0 {
		t.Fatalf("no input files in %q", inDir)
	}

	for _, inPath := range inputs {
		name := filepath.Base(inPath)
		outPath := filepath.Join(outDir, name)

		t.Run(name, func(t *testing.T) {
			u, err := ParseUniverseFromFile(inPath)
			if err != nil {
				t.Fatalf("parse %s: %v", inPath, err)
			}

			theta, relTol, expected, err := readExpected(outPath)
			if err != nil {
				t.Fatalf("read expected %s: %v", outPath, err)
			}
			if len(expected) != len(u.stars) {
				t.Fatalf("expected %d force rows, got %d",
					len(u.stars), len(expected))
			}

			qt := GenerateQuadTree(u)
			for i := range u.stars {
				got := CalculateNetForce(qt.root, u.stars[i], theta)
				want := expected[i]
				if !forceClose(got, want, relTol) {
					t.Fatalf("star %d: got (%.6e, %.6e), want (%.6e, %.6e) [rel %.2g theta=%.3f]",
					 i, got.x, got.y, want.x, want.y, relTol, theta)
				}
			}
		})
	}
}

// readExpected supports headers and forces:
//   theta=0.6
//   reltol=0.15
//   fx, fy
//   # comments and blanks ignored
func readExpected(path string) (theta, relTol float64, forces []pair, err error) {
	theta = 0.6
	relTol = 0.15

	f, err := os.Open(path)
	if err != nil { return 0, 0, nil, err }
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") { continue }

		// headers
		lc := strings.ToLower(line)
		switch {
		case strings.HasPrefix(lc, "theta="):
			val := strings.TrimSpace(line[len("theta="):])
			if theta, err = strconv.ParseFloat(val, 64); err != nil {
				return 0, 0, nil, fmt.Errorf("theta %q: %w", val, err)
			}
			continue
		case strings.HasPrefix(lc, "reltol="):
			val := strings.TrimSpace(line[len("reltol="):])
			if relTol, err = strconv.ParseFloat(val, 64); err != nil {
				return 0, 0, nil, fmt.Errorf("reltol %q: %w", val, err)
			}
			continue
		}

		// allow "i: fx, fy" -> strip index if present
		if idx := strings.Index(line, ":"); idx >= 0 {
			line = strings.TrimSpace(line[idx+1:])
		}

		parts := splitCSV(line)
		if len(parts) != 2 {
			return 0, 0, nil, fmt.Errorf("want 2 values, got %d in %q", len(parts), line)
		}
		fx, err := strconv.ParseFloat(parts[0], 64); if err != nil {
			return 0, 0, nil, fmt.Errorf("fx %q: %w", parts[0], err)
		}
		fy, err := strconv.ParseFloat(parts[1], 64); if err != nil {
			return 0, 0, nil, fmt.Errorf("fy %q: %w", parts[1], err)
		}
		forces = append(forces, pair{fx, fy})
	}
	if err := sc.Err(); err != nil {
		return 0, 0, nil, err
	}
	return theta, relTol, forces, nil
}

func forceClose(a, b pair, rel float64) bool {
	diff := math.Hypot(a.x-b.x, a.y-b.y)
	base := math.Max(math.Hypot(b.x, b.y), 1e-30)
	return diff/base <= rel
}