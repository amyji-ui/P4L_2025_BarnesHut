// Amy Ji
package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

// CenterOfMassTest holds one case worth of stars and the expected COM.
type CenterOfMassTest struct {
	stars         []*Star
	expectedX     float64
	expectedY     float64
}

func ReadCenterOfMassTests(directory string) []CenterOfMassTest {
	inputDir := filepath.Join(directory, "Input")
	outputDir := filepath.Join(directory, "Output")

	inputFiles := mustReadDir(inputDir)
	outputFiles := mustReadDir(outputDir)
	if len(inputFiles) != len(outputFiles) {
		panic("Error: number of input and output files do not match!")
	}

	// index output files by name for 1:1 pairing
	outMap := map[string]string{}
	for _, f := range outputFiles {
		outMap[f.Name()] = filepath.Join(outputDir, f.Name())
	}

	var tests []CenterOfMassTest
	for _, in := range inputFiles {
		name := in.Name()
		inPath := filepath.Join(inputDir, name)
		outPath, ok := outMap[name]
		if !ok {
			panic("Error: missing output file for " + name)
		}

		stars := readStarsFromFile(inPath)
		ex, ey := readTwoFloatsFromFile(outPath)

		tests = append(tests, CenterOfMassTest{
			stars:     stars,
			expectedX: ex,
			expectedY: ey,
		})
	}

	return tests
}

func mustReadDir(dir string) []os.DirEntry {
	files, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	return files
}

func readStarsFromFile(file string) []*Star {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var stars []*Star
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 3 {
			panic(fmt.Errorf("%s: each non-comment line must be 'x y m'", file))
		}
		x, err1 := strconv.ParseFloat(fields[0], 64)
		y, err2 := strconv.ParseFloat(fields[1], 64)
		m, err3 := strconv.ParseFloat(fields[2], 64)
		if err1 != nil || err2 != nil || err3 != nil {
			panic(fmt.Errorf("%s: bad star numbers on line: %q", file, line))
		}
		stars = append(stars, &Star{position: OrderedPair{x: x, y: y}, mass: m})
	}
	if err := sc.Err(); err != nil {
		panic(err)
	}
	return stars
}

func readTwoFloatsFromFile(file string) (float64, float64) {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			panic(fmt.Errorf("%s: need two floats 'ex ey'", file))
		}
		ex, err1 := strconv.ParseFloat(fields[0], 64)
		ey, err2 := strconv.ParseFloat(fields[1], 64)
		if err1 != nil || err2 != nil {
			panic(fmt.Errorf("%s: cannot parse expect_com numbers", file))
		}
		return ex, ey
	}
	if err := sc.Err(); err != nil {
		panic(err)
	}
	panic(fmt.Errorf("%s: no data found", file))
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}


// ------- actual tests -------
// TestCenterOfMass reads paired txt files from Tests/CenterOfMass/input|output
// and checks CenterOfMass result.
func TestCenterOfMass(t *testing.T) {
	tests := ReadCenterOfMassTests("Tests/CenterOfMass/")
	for i, tc := range tests {
		t.Run(fmt.Sprintf("case_%02d", i+1), func(t *testing.T) {
			got := CenterOfMass(tc.stars)
			if roundFloat(got.x, 6) != roundFloat(tc.expectedX, 6) ||
				roundFloat(got.y, 6) != roundFloat(tc.expectedY, 6) {
				t.Fatalf("CenterOfMass mismatch: got=(%.6f, %.6f) want=(%.6f, %.6f)",
					got.x, got.y, tc.expectedX, tc.expectedY)
			}
		})
	}
}