package geometry

import (
	"math"
	mrand "math/rand"
)

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

func MapBool(f func(*Triangle) bool, l []*Triangle) []*Triangle {
	temp := make([]*Triangle, 1)
	
	var isFirst bool = true
	var j int = 0
	
	for _,v := range l {
		if f(v) {
			if isFirst {
				temp[j] = v
				isFirst = false
				j = j+1
			} else {
				t := append(temp,v)
				temp = t
				j = j+1
			}
		}
	}
	return temp
}