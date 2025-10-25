// Amy Ji
package main
import("math")

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


func CalcDistance(p1, p2 OrderedPair) float64 {
	// this is the distance formula from days of precalculus long ago ...
	deltaX := p1.x - p2.x
	deltaY := p1.y - p2.y
	return math.Sqrt(deltaX*deltaX + deltaY*deltaY)
}


func CalcForce (s1, s2 *Star, G float64) OrderedPair {
	var Force OrderedPair

	if s1 == nil || s2 == nil {
		return Force
	}
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

