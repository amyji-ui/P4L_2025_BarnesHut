// Amy Ji
package main

import (
		"math")


//BarnesHut is our highest level function.
//Input: initial Universe object, a number of generations, and a time interval.
//Output: collection of Universe objects corresponding to updating the system
//over indicated number of generations every given time interval.
func BarnesHut(initialUniverse *Universe, numGens int, time, theta float64) []*Universe {
	timePoints := make([]*Universe, numGens+1)
	timePoints[0] = initialUniverse
	for i :=1; i < numGens+1; i++{
		u := UpdateUniverse(timePoints[i-1],time, theta)
		timePoints[i] = u
	}
	return timePoints

}


func UpdateUniverse(currentUniverse *Universe, time float64, theta float64) *Universe {
	
	newUniverse := CopyUniverse(currentUniverse)
	tree := GenerateQuadTree(currentUniverse)
	
	// First pass: calculate all new accelerations
    for i := range newUniverse.stars {
        newUniverse.stars[i].acceleration = UpdateAcceleration(tree.root, newUniverse.stars[i], theta)
    }
    
    // Second pass: update velocities and positions using OLD values from currentUniverse
    for i := range newUniverse.stars {
        oldStar := currentUniverse.stars[i]
        newUniverse.stars[i].velocity = UpdateVelocity(newUniverse.stars[i], oldStar.acceleration, time)
        newUniverse.stars[i].position = UpdatePosition(newUniverse.stars[i], oldStar.acceleration, oldStar.velocity, time)
    }
    return newUniverse 
}

// GenerateQuadTree takes as input a Universe object and returns a QuadTree
func GenerateQuadTree(currentUniverse *Universe) QuadTree {
    if len(currentUniverse.stars) == 0 {
        panic("No stars in universe for QuadTree construction")
    }
    
    rootQuadrant := Quadrant{0, 0, currentUniverse.width*2}
    rootNode := BuildNode(rootQuadrant, currentUniverse.stars)
    
    if rootNode == nil {
        panic("QuadTree root is nil - no stars were placed in tree")
    }
    
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
		// determine which sub-quadrant each star in starList belongs to.
		i := ChildIndex(quadrant, star.position)
		// append the star to the corresponding bucket.
		buckets[i] = append(buckets[i], star)
	}
	
	// Recursively build child nodes for each sub-quadrant.
	children := make([]*Node, 4)
	for i := 0; i < 4; i++ {
		// for one of the quadrants, call BuildNode on that quadrant and the corresponding bucket of stars.
		children[i] = BuildNode(subQuads[i], buckets[i])
	}

	// Create dummy node for this quadrant.
	quadMass := SumStarMasses(starList)
	quadCom := CenterOfMass(starList)
	dummy := &Star{
    	position: quadCom, 
    	mass: quadMass,
    	radius: 0, // Dummy stars shouldn't be drawn
    	red: 0, green: 0, blue: 0, // Make them invisible
	}

	// Return the internal node.
	return &Node {
		children : children,
		star: dummy,
		sector: quadrant,
	}

}

// node is the root of the tree.
func CalculateNetForce(node *Node, currStar *Star, theta float64) OrderedPair {
    var NetForce OrderedPair
    
    if node == nil || node.star == nil || currStar == nil {
        return NetForce
    }
    
    // a single star
    if node.children == nil && node.star != nil {
        if node.star == currStar {
            return NetForce
        }
        force := CalcForce(currStar, node.star, G)
        return force
    }
    
    // A cluster/galaxy
	s := node.sector.width
	d := CalcDistance(currStar.position, node.star.position)

	if d > 0 && (s/d) <= theta {
    	// Use the cluster approximation
    	force := CalcForce(currStar, node.star, G)
    	return force
	} else {
    	// Otherwise, look inside this cluster
    	for _, child := range node.children {
     		if child == nil {
            	continue
        	}
        	f := CalculateNetForce(child, currStar, theta)
        	NetForce.x += f.x
        	NetForce.y += f.y
    	}
	}
	return NetForce
}

func CalcForce (s1, s2 *Star, G float64) OrderedPair {
	var Force OrderedPair
	d := CalcDistance(s1.position, s2.position)

	if d == 0.0 {
		return Force
	}

	F := G * s1.mass * s2.mass / (d*d)
	deltaX := s2.position.x -s1.position.x
	deltaY := s2.position.y - s1.position.y 

	Force.x = F * deltaX/d 
	Force.y = F * deltaY/d 

	return Force

}

func CalcDistance(p1, p2 OrderedPair) float64 {
	// this is the distance formula from days of precalculus long ago ...
	deltaX := p1.x - p2.x
	deltaY := p1.y - p2.y
	return math.Sqrt(deltaX*deltaX + deltaY*deltaY)
}


// childIndex takes as input a parent Quadrant and a star position,
// and returns the index (0 to 3) of the child quadrant that contains the position.
func ChildIndex(parent Quadrant, p OrderedPair) int {
    cx := parent.x + parent.width/2
    cy := parent.y + parent.width/2

    if p.y >= cy && p.x <  cx { return 0 } // NW
    if p.y >= cy && p.x >= cx { return 1 } // NE
    if p.y <  cy && p.x <  cx { return 2 } // SW
    return 3                               // SE
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


func UpdateAcceleration(root *Node, s *Star, theta float64) OrderedPair {
	var accel OrderedPair 

	force := CalculateNetForce(root, s , theta)

	accel.x = force.x/s.mass
	accel.y = force.y/s.mass 

	return accel
}

func UpdateVelocity(s *Star, oldAcceleration OrderedPair, time float64) OrderedPair {
	var currentVelocity OrderedPair 

	currentVelocity.x = s.velocity.x + 0.5*(s.acceleration.x + oldAcceleration.x)*time

	currentVelocity.y = s.velocity.y + 0.5*(s.acceleration.y + oldAcceleration.y)*time

	return currentVelocity
}

func UpdatePosition(s *Star, oldAcceleration OrderedPair, oldVelocity OrderedPair, time float64) OrderedPair {
	var pos OrderedPair 

	pos.x = s.position.x + oldVelocity.x*time + 0.5*oldAcceleration.x*time*time 

	pos.y = s.position.y + oldVelocity.y*time + 0.5*oldAcceleration.y*time*time 

	return pos
}

func CopyUniverse(currentUniverse *Universe) *Universe {
	var newUniverse Universe

	newUniverse.width = currentUniverse.width

	numStars := len(currentUniverse.stars)

	newUniverse.stars = make([]*Star, numStars)

	for i := range newUniverse.stars {
		newUniverse.stars[i] = CopyStar(currentUniverse.stars[i])
	}

	return &newUniverse
}

func CopyStar(s *Star) *Star {
	var s2 Star

	s2.position.x = s.position.x
	s2.position.y = s.position.y

	s2.velocity.x = s.velocity.x
	s2.velocity.y = s.velocity.y  

	s2.acceleration.x = s.acceleration.x
	s2.acceleration.y = s.acceleration.y

	s2.mass = s.mass
	s2.radius = s.radius

	return &s2
}

