package geometry

import (
	"math"
//	"fmt"
)

type Intersectable interface {
	Dist() float64
	isHit() bool
	Object() *Triangle
}

type Miss struct {
	
}

func (o *Miss) Dist() float64 {
	return -math.MaxFloat64
}

func (o *Miss) isHit() bool {
	return false
}

func (o *Miss) Object() *Triangle {
	return nil
}

type IntersectBBox struct {
	dist float64
}

func (o *IntersectBBox) Dist() float64 {
	return o.dist
}

func (o *IntersectBBox) isHit() bool {
	return true
}

func (o *IntersectBBox) Object() *Triangle {
	return nil
}

type IntersectDist struct {
	dist float64
	triangle *Triangle
}

func (o *IntersectDist) Dist() float64 {
	return o.dist
}

func (o *IntersectDist) isHit() bool {
	return true
}

func (o *IntersectDist) Object() *Triangle {
	return o.triangle
}

func IsHit(l [] Intersectable) bool {
	result := false
	for _,v := range l {
		result = result || v.isHit()
	}
	return result
}

func MinIntersections(l []Intersectable) Intersectable {
	var result Intersectable
	
	min := func (i1 Intersectable, i2 Intersectable) Intersectable {
		var result Intersectable
		switch i1.(type) {
			case *Miss :
				switch i2.(type) {
					case *Miss :
						result = i1
					case *IntersectDist : 
						result = i2	
				}
			case *IntersectDist :
				switch i2.(type) {
					case *Miss : 
						result = i1
					case *IntersectDist :
						if i1.Dist() < i2.Dist() {
							result = i1
						} else {
							result = i2
						}
				}
		}
		return result
	}
	
	result = fold(min,l)
	if result.Dist() != -math.MaxFloat64 {
//		fmt.Printf("Found a minimum : %f\n",result.Dist())
	}
	return result
}

func fold(f func(i1 Intersectable, i2 Intersectable) Intersectable, l []Intersectable) Intersectable {
	intermediaryResult := f(l[0], l[1])
	temp := make([]Intersectable,1)
	if len(l) > 2 {
		temp[0] = intermediaryResult
		temp = append(temp,l[2:]...)
		return fold(f,temp) 
	} else {
		return intermediaryResult
	}
	
}