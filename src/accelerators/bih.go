package accelerators

import (
	"geometry"
	"math"
//	"fmt"
)

type Tree interface {
	insert(t *geometry.Triangle) Tree
	intersect(rayBBox *geometry.BoundingBox, ray *geometry.Ray, lastHit *geometry.Triangle) []geometry.Intersectable
}

type EmptyTree struct {	
}

func (t *EmptyTree) insert(triangle *geometry.Triangle) Tree {
//	fmt.Println("inserting a triangle in an emptytree")
	result := &Element{triangle.Box(), triangle}
	return result
}

func (t *EmptyTree) intersect(bbox *geometry.BoundingBox, ray *geometry.Ray, lastHit *geometry.Triangle) []geometry.Intersectable {
//	fmt.Println(bbox)
	panic("Cannot intersect with an empty tree")
}

type Element struct {	
	bbox geometry.BoundingBox
	triangle *geometry.Triangle
}

func (e1 *Element) insert(triangle *geometry.Triangle) Tree {
//	fmt.Println("inserting a triangle in an element")
	e2 := &Element{triangle.Box(), triangle}
	axis := [3]int{0,1,2}
	var lefts [3]*Element
	var rights [3]*Element
	for _,v := range axis {
		if e1.bbox.GetLowerFromAxis(v) < e2.bbox.GetLowerFromAxis(v) {
			lefts[v] = e1
			rights[v] = e2
		} else {
			lefts[v] = e2
			rights[v] = e1
		}
	}
	
	var differences [3]float64
	var difference float64
	var index int
	for _,v := range axis {
		differences[v] = math.Max(lefts[v].bbox.GetUpperFromAxis(v), rights[v].bbox.GetUpperFromAxis(v)) - lefts[v].bbox.GetLowerFromAxis(v)
	}
	
	difference = differences[0]
	for _,v := range axis {
		difference = math.Max(difference, differences[v])
		if difference == differences[v] {
			index = v
		}
	}
	
	il := lefts[index].bbox.GetUpperFromAxis(index)
	ir := rights[index].bbox.GetLowerFromAxis(index)
	
	node := &Node{il, ir, index, lefts[index], rights[index]}
	
	return node
}

func (t *Element) intersect(bbox *geometry.BoundingBox, ray *geometry.Ray, lastHit *geometry.Triangle) []geometry.Intersectable {
//	fmt.Println(bbox)
	if t.triangle == lastHit {
		miss := make([]geometry.Intersectable,1)
		miss[0] = &geometry.Miss{}
		return miss
	} else {
//		fmt.Println("Possible hit on an object")
		return geometry.RayIntersectsPrimitive(ray, t.triangle)
	}
}

type Node struct {
		il float64
		ir float64
		axis int
		left Tree
		right Tree
}

func (n *Node) insert(triangle *geometry.Triangle) Tree {
//	fmt.Println("inserting a triangle in a node")
	var result *Node
	if chooseLeft(triangle, n) {
		result = &Node{math.Max(n.il,triangle.Box().GetUpperFromAxis(n.axis)), n.ir, n.axis, n.left.insert(triangle), n.right}
	} else {
		result = &Node{n.il, math.Min(n.ir,triangle.Box().GetLowerFromAxis(n.axis)), n.axis, n.left, n.right.insert(triangle)}
	}
	
	return result
}

func clipRayStart(ray *geometry.BoundingBox, axis int, splitPlane float64) *geometry.BoundingBox {
	tdist := (splitPlane - ray.GetLowerFromAxis(axis)) / (ray.GetUpperFromAxis(axis) -ray.GetLowerFromAxis(axis))
	interval := func(low float64, up float64, axe int) *geometry.Interval {
		var i *geometry.Interval
		if axe == axis {
			i = geometry.NewInterval(splitPlane, up)
			return i
		} else {
			i = geometry.NewInterval(low +(up-low)*tdist, up)
			return i
		}
	}
	i0 := interval(ray.GetLowerFromAxis(0), ray.GetUpperFromAxis(0), 0)
	i1 := interval(ray.GetLowerFromAxis(1), ray.GetUpperFromAxis(1), 1)
	i2 := interval(ray.GetLowerFromAxis(2), ray.GetUpperFromAxis(2), 2)
	result := geometry.NewBBox(*i0, *i1, *i2)
//	fmt.Println(ray)
//	fmt.Println(result)
	return result
}

func clipRayEnd(ray *geometry.BoundingBox, axis int, splitPlane float64) *geometry.BoundingBox {
	tdist := (splitPlane - ray.GetLowerFromAxis(axis)) / (ray.GetUpperFromAxis(axis) -ray.GetLowerFromAxis(axis))
	interval := func(low float64, up float64, axe int) *geometry.Interval {
		var i *geometry.Interval
		if axe == axis {
			i = geometry.NewInterval(low, splitPlane)
			return i
		} else {
			i = geometry.NewInterval(low, low +(up-low)*tdist)
			return i
		}
	}
	i0 := interval(ray.GetLowerFromAxis(0), ray.GetUpperFromAxis(0), 0)
	i1 := interval(ray.GetLowerFromAxis(1), ray.GetUpperFromAxis(1), 1)
	i2 := interval(ray.GetLowerFromAxis(2), ray.GetUpperFromAxis(2), 2)
	result := geometry.NewBBox(*i0, *i1, *i2)
//	fmt.Println(ray)
//	fmt.Println(result)
	return result
	
}

func (t *Node) intersect(rayBBox *geometry.BoundingBox, ray *geometry.Ray, lastHit *geometry.Triangle) []geometry.Intersectable {
//	fmt.Println(rayBBox)
	rsAxisA := rayBBox.GetLowerFromAxis(t.axis)
	rsAxisB := rayBBox.GetUpperFromAxis(t.axis)
	
	var temp []geometry.Intersectable
	
	if rsAxisA < t.ir && rsAxisA > t.il {
		/*ray between places*/
		if rsAxisB <= t.il {
			return t.left.intersect(clipRayStart(rayBBox, t.axis, t.il), ray, lastHit)
		} else if rsAxisB >= t.ir {
			return t.right.intersect(clipRayStart(rayBBox, t.axis, t.ir), ray, lastHit)
		} else {
			temp = make([]geometry.Intersectable,1)
			temp[0] = &geometry.Miss{}
		}
	} else if rsAxisA <= t.il && rsAxisA >= t.ir {
		/*ray starts in both node*/
		newRs1 := func() *geometry.BoundingBox {
			if rsAxisB < t.ir {
				return clipRayEnd(rayBBox, t.axis, t.ir)
			} else {
				return rayBBox
			}
		}
		
		newRs2 := func() *geometry.BoundingBox {
			if rsAxisB > t.il {
				return clipRayEnd(rayBBox, t.axis, t.il)
			} else {
				return rayBBox
			}
		}
		temp = append(temp, t.right.intersect(newRs1(), ray, lastHit)...)
		temp = append(temp, t.left.intersect(newRs2(), ray, lastHit)...)
		
	} else if rsAxisA <= t.il {
		/*ray start in left node*/
		if rsAxisB < t.ir {
			if rsAxisB <= t.il {
				return t.left.intersect(rayBBox, ray, lastHit)
			} else {
				return t.left.intersect(clipRayEnd(rayBBox, t.axis, t.il), ray, lastHit)
			}
			
		} else {
			temp = append(temp, t.right.intersect(clipRayStart(rayBBox, t.axis, t.ir), ray, lastHit)...)
			temp = append(temp, t.left.intersect(clipRayEnd(rayBBox, t.axis, t.il), ray, lastHit)...)
		}
		
	} else if rsAxisA >= t.ir {
		/*ray start in right node*/
		if rsAxisB > t.il {
			if rsAxisB >= t.ir {
				return t.right.intersect(rayBBox, ray, lastHit)
			} else {
				return t.right.intersect(clipRayEnd(rayBBox, t.axis, t.ir), ray, lastHit)
			}
		} else {
			temp = append(temp, t.left.intersect(clipRayStart(rayBBox, t.axis, t.il), ray, lastHit)...)
			temp = append(temp, t.right.intersect(clipRayEnd(rayBBox, t.axis, t.ir), ray, lastHit)...)
		}
	}
	return temp
//	else {
//		miss := make([]geometry.Intersectable,1)
//		miss[0] = &geometry.Miss{}
//		return miss
//	}
}

func chooseLeft(triangle *geometry.Triangle, node *Node) bool {
	if triangle.Box().GetUpperFromAxis(node.axis) <= node.il {
		return true
	}
	
	if triangle.Box().GetLowerFromAxis(node.axis) >= node.ir {
		return false
	}
	
	delta_size_left := triangle.Box().GetUpperFromAxis(node.axis) - node.il
	delta_size_right := node.ir - triangle.Box().GetLowerFromAxis(node.axis)
	
	return delta_size_left <= delta_size_right
}

func Insert(t Tree, triangle *geometry.Triangle) Tree {
	return t.insert(triangle)
}

func Intersect(t Tree, rayBBox *geometry.BoundingBox, ray *geometry.Ray, lastHit *geometry.Triangle) []geometry.Intersectable {
	return t.intersect(rayBBox, ray, lastHit)
}