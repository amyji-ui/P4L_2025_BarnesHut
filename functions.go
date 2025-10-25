// Amy Ji
package main
// Also check the helperfunctions.go file for helper functions. It's too crowded to put everything here.

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

// GenerateQuadTree takes as input a Universe object and returns a QuadTree
func GenerateQuadTree (currentUniverse *Universe) QuadTree {

	rootQuadrant := Quadrant{0, 0, currentUniverse.width} // The root quadrant covers the entire universe
	// No for loop needed here since BuildNode is recursive!
	rootNode := BuildNode(rootQuadrant, currentUniverse.stars) // The first quad tree node corresponds to the root quadrant and all stars in the universe
	
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
	dummy := &Star{position : quadCom, mass : quadMass}

	// Return the internal node.
	return &Node {
		children : children,
		star: dummy,
		sector: quadrant,
	}

}

// node is the root of the tree.
func CalculateNetForce(node *Node, currStar *Star,theta float64) OrderedPair {

	var NetForce OrderedPair
	
	// empty
	if node == nil || node.star == nil || currStar == nil {
		return NetForce
	}
	// a single star
	if node.children == nil && node.star != nil {
		if node.star == currStar {
			return NetForce
		}
		return CalcForce(currStar, node.star,G)
	}
	// A cluster/galaxy
	s := node.sector.width
	// Because cluster is the aggregate of COM and total star mass, 
	// we can just call node.star.position
	d := CalcDistance(currStar.position, node.star.position)

	if d>0 && (s/d)<=theta {
		return CalcForce(currStar, node.star, G)
	}
	// if s/d > theta, we look inside this cluster.
	for _,child := range node.children{
		if child == nil {
			continue
		}
		f := CalculateNetForce(child, currStar, theta)
		NetForce.x += f.x
		NetForce.y += f.y
	}

    return NetForce
}


func UpdateUniverse(currentUniverse *Universe, time float64, theta float64) *Universe {
	
	newUniverse := CopyUniverse(currentUniverse)
	tree := GenerateQuadTree(currentUniverse)
	
	for i,s := range newUniverse.stars {
		oldAcceleration, oldVelocity := s.acceleration, s.velocity 

		newUniverse.stars[i].acceleration = UpdateAcceleration(tree.root, newUniverse.stars[i], theta )

		newUniverse.stars[i].velocity = UpdateVelocity(newUniverse.stars[i],oldAcceleration,time)
		
		newUniverse.stars[i].position = UpdatePosition(newUniverse.stars[i],oldAcceleration,oldVelocity,time)
	}
	return newUniverse 
}


