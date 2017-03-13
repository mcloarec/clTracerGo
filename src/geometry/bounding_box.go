package geometry

import (
	"math"
)

type BoundingBox struct {
	xInterval Interval
	yInterval Interval
	zInterval Interval
}

func NewBBox(itX Interval, itY Interval, itZ Interval) *BoundingBox {
	return &BoundingBox{itX, itY, itZ}
}

func NewBBoxFromTriangle(t Triangle) *BoundingBox {
	lowerX := math.Min(t.p0.x, math.Min(t.p1.x,t.p2.x))
	lowerY := math.Min(t.p0.y, math.Min(t.p1.y,t.p2.y))
	lowerZ := math.Min(t.p0.z, math.Min(t.p1.z,t.p2.z))
	upperX := math.Max(t.p0.x, math.Max(t.p1.x,t.p2.x))
	upperY := math.Max(t.p0.y, math.Max(t.p1.y,t.p2.y))
	upperZ := math.Max(t.p0.z, math.Max(t.p1.z,t.p2.z))
	realUpperX := upperX + ((math.Abs(upperX) + 1)*tolerance)
	realUpperY := upperY + ((math.Abs(upperY) + 1)*tolerance)
	realUpperZ := upperZ + ((math.Abs(upperZ) + 1)*tolerance)
	realLowerX := lowerX - ((math.Abs(lowerX) + 1)*tolerance)
	realLowerY := lowerY - ((math.Abs(lowerY) + 1)*tolerance)
	realLowerZ := lowerZ - ((math.Abs(lowerZ) + 1)*tolerance)
	itX := NewInterval(realLowerX, realUpperX)
	itY := NewInterval(realLowerY, realUpperY)
	itZ := NewInterval(realLowerZ, realUpperZ)
	b := &BoundingBox{*itX, *itY, *itZ}
	return b
}

func NewBBoxFromIntersection(pos *Point3, dir *Vector3, l []Intersectable) *BoundingBox {
	near := l[0].(*IntersectBBox)
	far := l[1].(*IntersectBBox)
	itX := NewInterval(pos.x + (dir.x * near.dist), pos.x + (dir.x * far.dist))
	itY := NewInterval(pos.y + (dir.y * near.dist), pos.x + (dir.y * far.dist))
	itZ := NewInterval(pos.z + (dir.z * near.dist), pos.x + (dir.z * far.dist))
	return &BoundingBox{*itX, *itY, *itZ}
}

func (bbox BoundingBox) GetLowerFromAxis(axis int) float64 {
	var result float64
	switch axis {
		case 0 : result = bbox.xInterval.lower
		case 1 : result = bbox.yInterval.lower
		case 2 : result = bbox.zInterval.lower
	}
	return result
}

func (bbox BoundingBox) GetUpperFromAxis(axis int) float64 {
	var result float64
	switch axis {
		case 0 : result = bbox.xInterval.upper
		case 1 : result = bbox.yInterval.upper
		case 2 : result = bbox.zInterval.upper
	}
	return result
}

func (bbox1 *BoundingBox) expand(bbox2 *BoundingBox) Expandable {
	result := &BoundingBox{bbox1.xInterval.expand(bbox2.xInterval),
		bbox1.yInterval.expand(bbox2.yInterval),
		bbox1.zInterval.expand(bbox2.zInterval)}
	return result
}

func (bbox *BoundingBox) intersectedByRay(ray *Ray) []Intersectable {
	xMin, xMax := MinMaxFloat((bbox.xInterval.lower - ray.origin.x) / ray.direction.x, (bbox.xInterval.upper - ray.origin.x) / ray.direction.x)
	yMin, yMax := MinMaxFloat((bbox.yInterval.lower - ray.origin.y) / ray.direction.y, (bbox.yInterval.upper - ray.origin.y) / ray.direction.y)
	zMin, zMax := MinMaxFloat((bbox.zInterval.lower - ray.origin.z) / ray.direction.z, (bbox.zInterval.upper - ray.origin.z) / ray.direction.z)
	
	nearDist := math.Max(xMin, math.Max(yMin, zMin))
	farDist := math.Min(xMax, math.Min(yMax, zMax))
	
	isParallel := !((!(ray.direction.x == 0.0) || (ray.origin.x > bbox.xInterval.lower && ray.origin.x < bbox.xInterval.upper)) &&
		(!(ray.direction.y == 0.0) || (ray.origin.y > bbox.yInterval.lower && ray.origin.y < bbox.yInterval.upper)) &&
		(!(ray.direction.z == 0.0) || (ray.origin.z > bbox.zInterval.lower && ray.origin.z < bbox.zInterval.upper)))
	
	if isParallel || nearDist > farDist {
//		fmt.Println("Missed !!")
		result := make([]Intersectable,1)
		miss := new(Miss)
		result[0] = miss
		return  result
	} else {
		result := make([]Intersectable,2)
		near := &IntersectBBox{nearDist}
		far := &IntersectBBox{farDist}
//		fmt.Printf("near : %f || far %f\n",nearDist, farDist)
		result[0] = near
		result[1] = far
		return result
	}
}

type EmptyBBox struct {
}

func (bbox1 *EmptyBBox) expand(bbox2 *BoundingBox) Expandable {
	return bbox2
}

func ExpandBBox(bbox1 Expandable, bbox2 *BoundingBox) Expandable {
	return bbox1.expand(bbox2)
}

func MinMaxFloat(x float64, y float64) (float64, float64) {
	if x > y {
		return y, x
	} else {
		return x, y
	}
}