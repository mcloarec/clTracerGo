package geometry

import (
	"math"
)

type Point3 struct {
	x float64
	y float64
	z float64
}

func NewPoint(x float64, y float64, z float64) *Point3 {
	p := &Point3{x, y, z}
	return p
}

func NewPointFromVector(pos *Point3, v *Vector3) *Point3 {
	return &Point3{pos.x + v.x, pos.y + v.y, pos.z + v.z}
}

type Color struct {
	r float64
	g float64
	b float64
}

func NewColor(r float64, g float64, b float64) *Color {
	c := &Color{r, g, b}
	return c
}

func (c1 Color) AddColor(c2 Color) *Color {
	c := &Color{c1.r + c2.r, c1.g + c2.g, c1.b + c2.b}
	return c
}

func (c Color) R() uint8 {
	return uint8(math.Max(math.Min(255,c.r*255),0))
}

func (c Color) G() uint8 {
	return uint8(math.Max(math.Min(255,c.g*255),0))
}

func (c Color) B() uint8 {
	return uint8(math.Max(math.Min(255,c.b*255),0))
}

func AddColor(c1 Color, c2 Color) *Color {
	return c1.AddColor(c2)
}

func MultC(c0 *Color, f float64) *Color {
	return &Color{c0.r * f, c0.g * f, c0.b * f}
}

func ColorMultC(c0 *Color, c1 *Color) *Color {
	return &Color{c0.r*c1.r, c0.g*c1.g, c0.b*c1.b}
}

func (c Color) IsNotBlack() bool {
	if c.r > 0 || c.g > 0 || c.b > 0 {
		return true
	} else {
		return false
	}
}

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

type Ray struct {
	origin Point3
	direction Vector3
}

func NewRay(origin Point3, direction Vector3) *Ray {
	r := &Ray{origin, direction}
	return r
}

type Interval struct {
	lower float64
	upper float64
}

func NewInterval(lower float64, upper float64) *Interval {
	i := &Interval{lower, upper}
	return i
}

func (inter1 Interval) expand(inter2 Interval) Interval {
	result := &Interval{math.Min(inter1.lower, inter2.lower), math.Max(inter1.upper, inter2.upper)}
	return *result
}

func MinMaxFloat(x float64, y float64) (float64, float64) {
	if x > y {
		return y, x
	} else {
		return x, y
	}
}