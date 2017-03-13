package util

import (
	"geometry"
)

func mapBool(f func(*geometry.Triangle) bool, l []*geometry.Triangle) []*geometry.Triangle {
	temp := make([]*geometry.Triangle, 1)
	
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

func zipWith(f func(geometry.Color, geometry.Color) *geometry.Color, l1 []geometry.Color, l2 []geometry.Color) []geometry.Color {
	if len(l1) != len(l2) {
		panic("impossible d'additionner les deux listes : elles devraient être égales")
	}
	
	temp := make([]geometry.Color, len(l1))
	
	for i,v := range l1 {
		temp[i] = *f(v, l2[i])
	}
	
	return temp
}

type Tuple [2]float64

type Sample [3]Tuple