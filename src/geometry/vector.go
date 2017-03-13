package geometry

import (
	"math"
)

type Vector3 struct {
	x float64
	y float64
	z float64
}

func NewVector(x float64, y float64, z float64) *Vector3 {
	v := &Vector3{x, y, z}
	return v
}

func (v Vector3) X() float64 {
	return v.x
}
func (v Vector3) Y() float64 {
	return v.y
}
func (v Vector3) Z() float64 {
	return v.z
}

func NewVectorFromPoints(p1 Point3, p2 Point3) *Vector3 {
	x := p2.x - p1.x
	y := p2.y - p1.y
	z := p2.z - p1.z
	v := &Vector3{x, y, z}
	return v
}

func UnitizeV(v Vector3) Vector3 {
	return v.multV(1/v.length())
}

func (v Vector3) AddV(v2 Vector3) Vector3 {
	result := &Vector3{v.x + v2.x, v.y + v2.y, v.z + v2.z}
	return *result
}

func AddV(v0 *Vector3, v1 *Vector3) *Vector3 {
	return &Vector3{v0.x + v1.x, v0.y + v1.y, v0.z + v1.z}
}

func (v Vector3) multV (f float64) Vector3 {
	result := &Vector3{v.x*f,v.y*f,v.z*f}
	return *result
}

func MultV(v Vector3, f float64) Vector3 {
	return v.multV(f)
}

func NegativeV(v Vector3) Vector3 {
	r := &Vector3{-v.x, -v.y, -v.z}
	return *r
}

func (v Vector3) length() float64 {
	return math.Sqrt((v.x*v.x) + (v.y*v.y) + (v.z*v.z))
}

func (v Vector3) CrossProduct(v0 Vector3) Vector3 {
	x := v.y*v0.z - v.z*v0.y
	y := v.z*v0.x - v.x*v0.z
	z := v.x*v0.y - v.y*v0.x
	result := &Vector3{x, y, z}
	return *result
}

func CrossProduct(v Vector3, v0 Vector3) Vector3 {
	x := v.y*v0.z - v.z*v0.y
	y := v.z*v0.x - v.x*v0.z
	z := v.x*v0.y - v.y*v0.x
	result := &Vector3{x, y, z}
	return *result
}

func (v Vector3) DotProduct(v0 Vector3) float64 {
	return v.x*v0.x + v.y*v0.y + v.z*v0.z
}

func DotProduct(v Vector3, v0 Vector3) float64 {
	return v.x*v0.x + v.y*v0.y + v.z*v0.z
}

func IsNillVector(v Vector3) bool {
	return v.x == 0 && v.y == 0 && v.z == 0
}