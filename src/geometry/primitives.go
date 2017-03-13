package geometry

import (

)

const tolerance float64 = 1/1024.0

type Expandable interface {
	expand(ex *BoundingBox) Expandable
}

type Primitive interface {
	intersectedByRay(ray *Ray) []Intersectable
}

func RayIntersectsPrimitive(ray *Ray, p Primitive) []Intersectable {
	return p.intersectedByRay(ray)
}