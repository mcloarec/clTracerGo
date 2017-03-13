package util

import (
	"regexp"
	"strconv"
	"core"
	"geometry"
	"spacepartitionning"
//	"fmt"
)

func ParseFile(s string) *core.Scene {
	sceneOpts := ParseSceneOpts(s)
	camera := ParseCamera(s)
	world := ParseWorld(s)
	primitives := ParsePrimitives(s)
	lights := mapBool(geometry.IsLight,primitives)
	
	var tree spacepartitionning.Tree
	
	tree = &spacepartitionning.EmptyTree{}
	
	for _,v := range primitives {
		tree = spacepartitionning.Insert(tree, v)
	}
	
	
	var enveloppe geometry.Expandable
	enveloppe = &geometry.EmptyBBox{}
	
	for _,v := range primitives {
		bbox := v.Box()
		enveloppe = geometry.ExpandBBox(enveloppe,&bbox)
	}
	
	return core.NewScene(sceneOpts,camera,world,primitives,lights,tree,enveloppe.(*geometry.BoundingBox))
}

func ParseSceneOpts(s string) (opts *core.SceneOpts) {
	itRE := regexp.MustCompile(`(?m)^([0-9]+)$`)
	index := itRE.FindStringSubmatch(s)
	
	var iterations int64
	
	if len(index) > 1 {
		iterations,_ = strconv.ParseInt(index[1],10,0)
	}
	
	frSizeRE := regexp.MustCompile(`(?m)^([0-9]+) ([0-9]+)$`)
	index = frSizeRE.FindStringSubmatch(s)
	var width, height int64
	
	if len(index) > 1 {
		width,_ = strconv.ParseInt(index[1],10,0)
		height,_ = strconv.ParseInt(index[2],10,0)
	}
	
	opts = core.NewOpts(iterations, width, height)
	
	return
}

func ParseCamera(s string) (camera *core.Camera) {
	cameraRE := regexp.MustCompile(`(?m)^\((-?[0-9]+\.[0-9]+) (-?[0-9]+\.[0-9]+) (-?[0-9]+\.[0-9]+)\) \((-?[0-9]+\.?[0-9]*) (-?[0-9]+\.?[0-9]*) (-?[0-9]+\.?[0-9]*)\) ([0-9]+)$`)
	index := cameraRE.FindStringSubmatch(s)
	
	var xPos, yPos, zPos float64
	var xDir, yDir, zDir float64
	var fov float64
	var pos *geometry.Point3
	var dir *geometry.Vector3
	
	if len(index) > 1 {
		xPos,_ = strconv.ParseFloat(index[1],64)
		yPos,_ = strconv.ParseFloat(index[2],64)
		zPos,_ = strconv.ParseFloat(index[3],64)
		
		pos = geometry.NewPoint(xPos, yPos, zPos)
		
		xDir,_ = strconv.ParseFloat(index[4],64)
		yDir,_ = strconv.ParseFloat(index[5],64)
		zDir,_ = strconv.ParseFloat(index[6],64)
		
		dir = geometry.NewVector(xDir, yDir, zDir)
		
		fov,_ = strconv.ParseFloat(index[7],64)
	}
	
	camera = core.NewCamera(*pos, *dir, fov)
	
	return
}

func ParseWorld(s string) (world *core.World) {
	worldRE := regexp.MustCompile(`(?m)^\((-?[0-9]+\.[0-9]+) (-?[0-9]+\.[0-9]+) (-?[0-9]+\.[0-9]+)\) \((-?[0-9]+\.[0-9]+) (-?[0-9]+\.[0-9]+) (-?[0-9]+\.[0-9]+)\)$`)
	index := worldRE.FindStringSubmatch(s)
	
	var rSky, gSky, bSky float64
	var rGrd, gGrd, bGrd float64
	var skyEmission, groundReflexion *geometry.Color
	
	if len(index) > 1 {
		rSky,_ = strconv.ParseFloat(index[1],64)
		gSky,_ = strconv.ParseFloat(index[2],64)
		bSky,_ = strconv.ParseFloat(index[3],64)
		
		skyEmission = geometry.NewColor(rSky, gSky, bSky)
		
		rGrd,_ = strconv.ParseFloat(index[4],64)
		gGrd,_ = strconv.ParseFloat(index[5],64)
		bGrd,_ = strconv.ParseFloat(index[6],64)
		
		groundReflexion = geometry.NewColor(rGrd, gGrd, bGrd)
	}
	
	world = core.NewWorld(*skyEmission, *groundReflexion)
	
	return
}

func ParsePrimitives(s string) (prims []*geometry.Triangle) {
	trianglesRE := regexp.MustCompile(`(?m)^([A-Za-z]+) \((-?[0-9]+\.[0-9]+) (-?[0-9]+\.[0-9]+) (-?[0-9]+\.[0-9]+)\) \((-?[0-9]+\.[0-9]+) (-?[0-9]+\.[0-9]+) (-?[0-9]+\.[0-9]+)\) \((-?[0-9]+\.[0-9]+) (-?[0-9]+\.[0-9]+) (-?[0-9]+\.[0-9]+)\)  \((-?[0-9]+\.[0-9]+) (-?[0-9]+\.[0-9]+) (-?[0-9]+\.[0-9]+)\) \((-?[0-9]+\.?[0-9]*) (-?[0-9]+\.?[0-9]*) (-?[0-9]+\.?[0-9]*)\)$`)
	
	index2 := trianglesRE.FindAllStringSubmatch(s,-1)
	
	if len(index2) > 0 {
		
		prims = make([]*geometry.Triangle,1)
		
		for i:=0;i<len(index2);i++{
			var xP1,yP1,zP1 float64
			var xP2,yP2,zP2 float64
			var xP3,yP3,zP3 float64
			var p1, p2, p3 *geometry.Point3
			
			var rEmit, gEmit, bEmit float64
			var rDiffuse, gDiffuse, bDiffuse float64
			var emit, diffuse *geometry.Color
			var triangleID string
			
			if len(index2[i]) == 17 {
				triangleID = index2[i][1]
				xP1,_ = strconv.ParseFloat(index2[i][2],64)
				yP1,_ = strconv.ParseFloat(index2[i][3],64)
				zP1,_ = strconv.ParseFloat(index2[i][4],64)
				xP2,_ = strconv.ParseFloat(index2[i][5],64)
				yP2,_ = strconv.ParseFloat(index2[i][6],64)
				zP2,_ = strconv.ParseFloat(index2[i][7],64)
				xP3,_ = strconv.ParseFloat(index2[i][8],64)
				yP3,_ = strconv.ParseFloat(index2[i][9],64)
				zP3,_ = strconv.ParseFloat(index2[i][10],64)
				rDiffuse,_ = strconv.ParseFloat(index2[i][11],64)
				gDiffuse,_ = strconv.ParseFloat(index2[i][12],64)
				bDiffuse,_ = strconv.ParseFloat(index2[i][13],64)
				rEmit,_ = strconv.ParseFloat(index2[i][14],64)
				gEmit,_ = strconv.ParseFloat(index2[i][15],64)
				bEmit,_ = strconv.ParseFloat(index2[i][16],64)
			}
			
			p1 = geometry.NewPoint(xP1,yP1,zP1)
			p2 = geometry.NewPoint(xP2,yP2,zP2)
			p3 = geometry.NewPoint(xP3,yP3,zP3)
			emit = geometry.NewColor(rEmit, gEmit, bEmit)
			diffuse = geometry.NewColor(rDiffuse, gDiffuse, bDiffuse)
			
			p := geometry.NewTriangle(triangleID, p1, p2, p3, emit, diffuse)
			
			if i == 0 {
				prims[i] = p
			} else {
				temp := append(prims, p)
				prims = temp
			}
		}
	}
	
	return
}
