package geometry

import (
	"math"
	mrand "math/rand"
//	"fmt"
)

const tolerance float64 = 1/1024.0

type Expandable interface {
	expand(ex *BoundingBox) Expandable
}

type Primitive interface {
	intersectedByRay(ray *Ray) []Intersectable
}

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



type Triangle struct {
	id string
	p0, p1, p2 Point3
	emit, diffuse Color
	bbox BoundingBox
	edge0, edge1, edge2 Vector3
	tangent, normal Vector3
}

func NewTriangle(id string, p0 *Point3, p1 *Point3, p2 *Point3, emit *Color, diffuse *Color) *Triangle {
	t := &Triangle{}
	t.p0 = *p0
	t.p1 = *p1
	t.p2 = *p2
	t.emit = *emit
	t.diffuse = *diffuse
	t.id = id
	t.init()
	return t
}

func (t *Triangle) init() {
	t.bbox = *NewBBoxFromTriangle(*t)
	t.edge0 = *NewVectorFromPoints(t.p0, t.p1)
	t.edge1 = *NewVectorFromPoints(t.p1, t.p2)
	t.edge2 = *NewVectorFromPoints(t.p0, t.p2)
	t.tangent = UnitizeV(t.edge0)
	t.normal = t.edge0.CrossProduct(t.edge1)
}

func (t *Triangle) Area() float64 {
	pa2 := t.normal
	return math.Sqrt(pa2.DotProduct(pa2)) * 0.5
}
func IsLight(t *Triangle) bool {
	return t.emit.IsNotBlack()
}

func (t *Triangle) Box() BoundingBox {
	return t.bbox
}

func (t *Triangle) intersectedByRay(ray *Ray) []Intersectable {
	temp := make([]Intersectable, 1)
	p := ray.direction.CrossProduct(t.edge2)
	det := p.DotProduct(t.edge0)
	
//	fmt.Println("Det : ", det)
	if math.Abs(det) < 0.000001 {
		temp[0] = &Miss{}
		return temp
	}
	tPrim := NewVectorFromPoints(t.p0, ray.origin)
	u := p.DotProduct(*tPrim) / det
	
//	fmt.Println("u : ", u)
	if u < 0.0 || u > 1.0 {
		temp[0] = &Miss{}
		return temp
	}
	q := tPrim.CrossProduct(t.edge0)
	v := q.DotProduct(ray.direction) / det
	
//	fmt.Println("v : ", u)
	if v < 0.0 || (u + v) > 1.0 {
		temp[0] = &Miss{}
		return temp
	}
	
	tSecond := q.DotProduct(t.edge2) / det
	
//	fmt.Println("dist : ", tSecond)
	
	if tSecond < 1e-8 {
		temp[0] = &Miss{}
		return temp
	}
	
//	fmt.Println("Distance : ", tSecond)
	intersection := &IntersectDist{tSecond, t}
	temp[0] = intersection
	return temp
}

func (triangle *Triangle) samplePoint() *Point3 {
	sqr1 := math.Sqrt(mrand.Float64())
	r2 := mrand.Float64()
	c0 := 1.0 - sqr1
	c1 := (1.0 - r2) * sqr1
	
	ac0 := triangle.edge0.multV(c0)
	ac2 := triangle.edge2.multV(c1)
	
	sum := ac0.AddV(ac2)
	
	return NewPointFromVector(&triangle.p0, &sum)
}

func SamplePoint(triangle *Triangle) *Point3 {
	return triangle.samplePoint()
}

func RayIntersectsPrimitive(ray *Ray, p Primitive) []Intersectable {
	return p.intersectedByRay(ray)
}