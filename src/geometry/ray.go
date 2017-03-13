package geometry

import (

)

type Ray struct {
	origin Point3
	direction Vector3
}

func NewRay(origin Point3, direction Vector3) *Ray {
	r := &Ray{origin, direction}
	return r
}