package geometry

import (
	"math"
)

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

func zipWith(f func(Color, Color) *Color, l1 []Color, l2 []Color) []Color {
	if len(l1) != len(l2) {
		panic("impossible d'additionner les deux listes : elles devraient être égales")
	}
	
	temp := make([]Color, len(l1))
	
	for i,v := range l1 {
		temp[i] = *f(v, l2[i])
	}
	
	return temp
}