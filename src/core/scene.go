package core

import (
	. "geometry"
	"accelerators"
	"math"
	mrand "math/rand"
	"image"
	c "image/color"
	"os"
	png "image/png"
//	"fmt"
)

//rgbLuminance := NewColor(0.2126, 0.7152, 0.0722)
//gammaEncode := 0.45

type Scene struct {
	opts *SceneOpts
	camera *Camera
	world *World
	prims []*Triangle
	lights []*Triangle
	tree accelerators.Tree
	enveloppe *BoundingBox
}

func NewScene(sceneOpts *SceneOpts, camera *Camera, world *World, prims []*Triangle, lights []*Triangle, tree accelerators.Tree, enveloppe *BoundingBox) *Scene {
	return &Scene{sceneOpts, camera, world, prims, lights, tree, enveloppe}
}

type SceneOpts struct {
	iterations, imWidth, imHeight int
}

func NewOpts(it int64, width int64, height int64) *SceneOpts {
	opts := &SceneOpts{int(it), int(width), int(height)}
	return opts
}

type Camera struct {
	position Point3
	direction Vector3
	viewAngle float64
	right Vector3
	up Vector3
	tanViewAngle float64
	
}

func NewCamera(pos Point3, dir Vector3, fov float64) *Camera {
	viewAngle := math.Min(math.Max(10.0, fov), fov) * math.Pi / 180.0
	tanViewAngle := math.Tan(viewAngle * 0.5)
	
	var viewDirection Vector3
	
	if IsNillVector(UnitizeV(dir)) {
		viewDirection = *NewVector(0.0, 0.0, 1.0)	
	} else {
		viewDirection = UnitizeV(dir)
	}
	
	var right Vector3
	var up Vector3
	
	right = UnitizeV(NewVector(0.0, 1.0, 0.0).CrossProduct(UnitizeV(dir)))
	
	if IsNillVector(right) {
		if viewDirection.Y() > 0 {
			up = *NewVector(0.0, 0.0, 1.0)
		} else {
			up = *NewVector(0.0, 0.0, -1.0)
		}
		right = UnitizeV(up.CrossProduct(viewDirection))
	} else {
		up = UnitizeV(viewDirection.CrossProduct(right))
	}
	
	camera := &Camera{pos, viewDirection, viewAngle, right, up, tanViewAngle}
	return camera
}

type World struct {
	skyEmission, groundReflexion Color
}

func NewWorld(skyEmission Color, groundReflexion Color) *World {
	world := &World{skyEmission, groundReflexion}
	return world
}

func (scene *Scene) Render(epoch int64) {
	
	img := image.NewRGBA(image.Rect(0,0,scene.opts.imWidth-1,scene.opts.imHeight-1))
	
	aspect := float64(scene.opts.imWidth) / float64(scene.opts.imHeight)
	
	colors := make([]Color, (scene.opts.imWidth * scene.opts.imHeight))
	
//	for x := 0; x < scene.opts.imWidth ; x++ {
//		for y := 0 ; y < scene.opts.imHeight ; y++ {
//			xCoeff := ((float64(x) + mrand.Float64()) * 2.0 / float64(scene.opts.imWidth)) - 1.0
//			yCoeff := ((float64(y) + mrand.Float64()) * 2.0 / float64(scene.opts.imHeight)) - 1.0
//			offset := MultV(scene.camera.right, xCoeff).AddV(MultV(scene.camera.up, aspect * yCoeff))
//			var sampleDirection *Vector3 = new(Vector3)
//			*sampleDirection = UnitizeV(scene.camera.direction.AddV(MultV(offset, scene.camera.tanViewAngle)))
//			color := scene.getRadiance(&scene.camera.position, &sampleDirection, nil)
//			colors[x+(scene.opts.imWidth*y)] = *AddColor(*color, colors[x+(scene.opts.imWidth*y)])
////			color := scene.getRadiance(&scene.camera.position, &sampleDirection, nil)
////			img.Set(x,y,c.RGBA{color.R(), color.G(), color.B(),255})
////			fmt.Printf("X=%d | Y=%d\n",x,y)
//		}
//	}
	
	sem := make(chan int, 4)  // Buffering optional but sensible.
	
	for i := 0; i < 4 ; i++ {
		go func(){
			for i := 0; i < scene.opts.iterations/4 ; i++ {
				for x := 0; x < scene.opts.imWidth ; x++ {
					for y := 0 ; y < scene.opts.imHeight ; y++ {
						xCoeff := ((float64(x) + mrand.Float64()) * 2.0 / float64(scene.opts.imWidth)) - 1.0
						yCoeff := ((float64(y) + mrand.Float64()) * 2.0 / float64(scene.opts.imHeight)) - 1.0
						offset := MultV(scene.camera.right, xCoeff).AddV(MultV(scene.camera.up, aspect * yCoeff))
						var sampleDirection *Vector3 = new(Vector3)
						*sampleDirection = UnitizeV(scene.camera.direction.AddV(MultV(offset, scene.camera.tanViewAngle)))
						colors[x+(scene.opts.imWidth*y)] = *AddColor(*scene.getRadiance(&scene.camera.position, &sampleDirection, nil),colors[x+(scene.opts.imWidth*y)])
			//			color := scene.getRadiance(&scene.camera.position, &sampleDirection, nil)
		//				fmt.Printf("X=%d | Y=%d\n",x,y)
					}
				}
			}
			sem <- 1
		}()
	}
	
	for i := 0; i < 4; i++ {
        <-sem    // wait for one task to complete
    }
	
	for xx := 0; xx < scene.opts.imWidth ; xx++ {
			for yy := 0 ; yy < scene.opts.imHeight ; yy++ {
				color := MultC(&colors[xx+(scene.opts.imWidth*yy)],1/float64(scene.opts.iterations))
				img.Set(xx,yy,c.RGBA{color.R(), color.G(), color.B(),255})
//				fmt.Printf("X=%d | Y=%d\n",xx,yy)
			}
	}
	f, _ := os.Create("result.png")
	png.Encode(f,img)
	f.Close()
}

func (scene *Scene) getRadiance(pos *Point3, dir **Vector3, lastHit *Triangle) *Color {
	var radiance *Color = NewColor(0, 0, 0)
	var hitObject *Triangle
	var hitPosition *Point3
	
	rayBackDirection := NegativeV(**dir)
	
	scene.intersection(pos, *dir, lastHit, &hitObject, &hitPosition)
	
	if hitObject != nil {
	
		sfp := NewSurfacePoint(hitPosition,hitObject)
		
		localEmission := map[bool](*Color){true:NewColor(0,0,0), false:sfp.SurfacePointEmission(pos,&rayBackDirection,false)}[lastHit!=nil]
		
		emitterSample := scene.sampleEmitters(&rayBackDirection, sfp)
		
		/* recursed reflection */
		var recursedReflection *Color = NewColor(0, 0, 0)
		
		/* single hemisphere sample, ideal diffuse BRDF:
	               reflected = (inradiance * pi) * (cos(in) / pi * color) *
	                  reflectance
	            -- reflectance magnitude is 'scaled' by the russian roulette,
	            cos is importance sampled (both done by SurfacePoint),
	            and the pi and 1/pi cancel out -- leaving just:
	               inradiance * reflectance color */
		var nextDirection *Vector3
		var color *Color
		
		if sfp.SurfacePointNextDirection(&rayBackDirection, &nextDirection, &color) {
			recursed := scene.getRadiance(sfp.HitPosition(), &nextDirection, sfp.Object())
			recursedReflection = ColorMultC(recursed, color)
		}
		
		radiance = AddColor(*localEmission,*emitterSample)
		radiance = AddColor(*radiance, *recursedReflection)
	
	}
	return radiance
}

func (scene *Scene) intersection(pos *Point3, dir *Vector3, lastHit *Triangle, hitObject **Triangle, hitPosition **Point3) {
	ray := NewRay(*pos, *dir)
	sceneIntersection := RayIntersectsPrimitive(ray, scene.enveloppe)
	if IsHit(sceneIntersection) {
		rayBBox := NewBBoxFromIntersection(pos, dir, sceneIntersection)
		intersections := accelerators.Intersect(scene.tree, rayBBox, ray, lastHit)
//		fmt.Printf("Intersections size : %d\n",len(intersections))
		intersection := MinIntersections(intersections)
//		fmt.Printf("adresse : %p\n",dir)
		if intersection.Object() != nil {
			tempRay := MultV(*dir,intersection.Dist())
			*hitPosition = NewPointFromVector(pos,&tempRay)
			*hitObject = intersection.Object()
		} else {
			*hitObject = nil
		}
		
	} else {
		*hitObject = nil
	}
}

func (scene *Scene) getEmitter(emitterPosition **Point3, emitterObject **Triangle) {
	if len(scene.lights) > 0 {
		index := int(math.Floor(mrand.Float64() * float64(len(scene.lights))))
		index = map[bool]int{true:index, false:len(scene.lights)-1}[index < len(scene.lights)]
		*emitterObject = scene.lights[index]
		*emitterPosition = SamplePoint(*emitterObject)
		
	} else {
		*emitterPosition = NewPoint(0,0,0)
		*emitterObject = nil
	}
}

func (scene *Scene) sampleEmitters(rayBackDirection *Vector3, sfp *SurfacePoint) *Color {
	radiance := NewColor(0,0,0)
	var emitterPosition *Point3
	var emitterObject *Triangle
	var hitObject *Triangle
	var hitPosition *Point3
	
	scene.getEmitter(&emitterPosition,&emitterObject)
	
	if emitterObject != nil {
		emitDirection := UnitizeV(*NewVectorFromPoints(*sfp.HitPosition(),*emitterPosition))
		
		scene.intersection(sfp.HitPosition(),&emitDirection,sfp.Object(),&hitObject,&hitPosition)
		
		//if hitObject == nil || emitterObject == hitObject {
			sp := NewSurfacePoint(emitterPosition, emitterObject)
			backEmitDirection := NegativeV(emitDirection)
			emissionIn := sp.SurfacePointEmission(sfp.HitPosition(), &backEmitDirection, true)
			emissionAll := MultC(emissionIn,float64(len(scene.lights)))
			radiance = sfp.SurfacePointReflection(&emitDirection,emissionAll,rayBackDirection)
			return radiance
		//}
		
		//return radiance
	} else {
		return radiance
	}
	
}