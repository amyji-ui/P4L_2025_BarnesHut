package main

//BarnesHut is our highest level function.
//Input: initial Universe object, a number of generations, and a time interval.
//Output: collection of Universe objects corresponding to updating the system
//over indicated number of generations every given time interval.
func BarnesHut(initialUniverse *Universe, numGens int, time, theta float64) []*Universe {
	timePoints := make([]*Universe, numGens+1)

	// Your code goes here. Use subroutines! :)

	return timePoints
}

func GenerateQuadTree (currentUniverse Universe) QuadTree {

	rootQuadrant := Quadrant{0, 0, currentUniverse.width}
	rootNode := BuildNode(rootQuadrant, currentUniverse.stars)

	return QuadTree{root: rootNode}

}


func BuildNode (quadrant Quadrant, stars []*Star) *Node {
	starList := CountStarsInQuadrant(quadrant, stars)
	n := len(starList)

	// Two base cases: no star in the quadrant or ony 1 star in the quadrant
	if n == 0 {
		return nil
	}
	if n == 1 {
		return &Node{
			children: nil,
			star: starList[0],
			sector: quadrant,
		}
	}

	// If there's more than one star in this quadrant, we need to keep splitting.
	subQuads := SplitQuadrant(quadrant)
	// Partition all stars into four groups.
	// We can then assign them base on the locations: nw, ne, sw, se (sequentially!).
	buckets := make([][]*Star, 4)
	for _, star := range starList {
		i := childIndex(quadrant, star.position)
		// Recurse using each of the four subquadrant and subset of stars.
		buckets[i] = append(buckets[i], star)
	}
	
	children := make([]*Node, 4)
	for i := 0; i < 4; i++ {
		children[i] = BuildNode(subQuads[i], buckets[i])
	}
	
	quadMass := SumStarMasses(starList)
	quadCom := CenterOfMass(starList)
	dummy := &Star{position : quadCom, mass : quadMass}

	return &Node {
		children : children,
		star: dummy,
		sector: quadrant,
	}

}

// childIndex takes as input a parent Quadrant and an OrderedPair position,
// and returns the index (0 to 3) of the child quadrant that contains the position.
func childIndex(parent Quadrant, p OrderedPair) int {
	cx := parent.x + parent.width/2
	cy := parent.y + parent.width/2
	top := p.y >= cy
	right := p.x >= cx
	switch {
	case  top && !right: return 0 // nw
	case  top &&  right: return 1 // ne
	case !top && !right: return 2 // sw
	default:              return 3 // se
	}
}

// CountStarsInQuadrant takes as input a Quadrant and a slice of Star pointers,
// and returns a slice of Star pointers corresponding to the stars within that Quadrant.
func CountStarsInQuadrant (quadrant Quadrant, stars []*Star) []*Star {

	var starList []*Star
	
	for _, star := range stars{
		if (star.position.x >= quadrant.x && star.position.x < quadrant.x + quadrant.width) && (star.position.y >= quadrant.y && star.position.y < quadrant.y + quadrant.width){
			starList = append(starList, star)
		}
	}
	return starList
}

// SplitQuadrant takes as input a Quadrant and returns a slice of four Quadrants
// corresponding to the northwest, northeast, southwest, and southeast sub-quadrants.
func SplitQuadrant (quadrant Quadrant) []Quadrant {
	subQuadrants := make([]Quadrant, 4)
	mid := quadrant.width / 2.0

	subQuadrants[0] = Quadrant{quadrant.x, quadrant.y + mid, mid}        // nw
	subQuadrants[1] = Quadrant{quadrant.x + mid, quadrant.y + mid, mid}  // ne
	subQuadrants[2] = Quadrant{quadrant.x, quadrant.y, mid}              // sw
	subQuadrants[3] = Quadrant{quadrant.x + mid, quadrant.y, mid}        // se

	return subQuadrants
}

// SumStarMasses takes as input a slice of Star pointers and returns the sum of their masses.	
func SumStarMasses (stars []*Star) float64 {

	SumMass := 0.0
	for _,star := range stars {
		SumMass = SumMass + star.mass
	}

	return SumMass
}

// CenterOfMass takes as input a slice of Star pointers and returns an OrderedPair corresponding to the center of mass of those stars.
func CenterOfMass (stars []*Star) OrderedPair {
	var center OrderedPair

	SumMass := SumStarMasses(stars)

	if SumMass == 0 {
		return center
	}

	x := 0.0
	y := 0.0

	for _,star := range stars {
		x += star.position.x * star.mass
		y += star.position.y * star.mass
	}

	center.x = x/SumMass
	center.y = y/SumMass

	return center
}

