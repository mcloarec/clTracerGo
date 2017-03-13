package geometry

import (

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