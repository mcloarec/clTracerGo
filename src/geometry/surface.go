package geometry

import (
	mrand "math/rand"
	"math"
)

type SurfacePoint struct {
	pTriangle *Triangle
	pHitPosition *Point3
}

func NewSurfacePoint(pPos *Point3, pT *Triangle) *SurfacePoint {
	return &SurfacePoint{pT, pPos}
}

func (pSp *SurfacePoint) Object() *Triangle {
	return pSp.pTriangle
}

func (pSp *SurfacePoint) HitPosition() *Point3 {
	return pSp.pHitPosition
}

func (pSp *SurfacePoint) SurfacePointEmission(pPos *Point3, pDir *Vector3, isSolidAngle bool) *Color {
	var solidAngle float64
	ray := NewVectorFromPoints(*pSp.pHitPosition,*pPos)
	distance2 := ray.DotProduct(*ray)
	distance2 = math.Max(distance2, 1e-6)
	normal := pSp.pTriangle.normal
	cosout := pDir.DotProduct(normal)
	cosArea := cosout * pSp.pTriangle.Area()
	/* Emit from front face of surface only with infinity clamped out*/
	solidAngle = map[bool]float64{true:(cosArea/distance2),false:1.0}[isSolidAngle]
	result := map[bool](*Color){true:MultC(&pSp.pTriangle.emit,solidAngle),false:NewColor(0,0,0)}[cosArea > 0.0]
	return result
}

func (pSp *SurfacePoint) SurfacePointNextDirection(pInDirection *Vector3, pOutDirection **Vector3, pColor **Color) bool {
	
	reflectivityMean := (pSp.pTriangle.diffuse.r + pSp.pTriangle.diffuse.g + pSp.pTriangle.diffuse.b) / 3.0
	
	/* russian-roulette for reflectance 'magnitude' */
	isAlive := mrand.Float64() < reflectivityMean
	
	if isAlive {
		/* cosine-weighted importance sample hemisphere */
		twopr1 := math.Pi * 2.0 * mrand.Float64()
		sr2 := math.Sqrt(mrand.Float64())
		
		/* make coord frame coefficients (z in normal direction) */
		x := math.Cos(twopr1) * sr2
		y := math.Sin(twopr1) * sr2
		z := math.Sqrt(1.0 - (sr2 * sr2))
		
		/* make coord frame */
		t := pSp.pTriangle.tangent
		n := pSp.pTriangle.normal
		
		/* put normal on inward ray side of surface (preventing transmission) */
		if n.DotProduct(*pInDirection) < 0.0 {
			n = NegativeV(n)
		}
		
		c := n.CrossProduct(t)
		
		/* scale frame by coefficients */
		tx := MultV(t, x)
		cy := MultV(c, y)
		nz := MultV(n, z)
		
		/* make direction from sum of scaled components */
		sum := AddV(&tx, &cy)
		*pOutDirection = AddV(sum, &nz)
		
		/* make color by dividing-out mean from reflectivity */
		*pColor = MultC(&pSp.pTriangle.diffuse, 1.0/ reflectivityMean) 
	}
	
	/* discluding degenerate result direction */
	return isAlive && ((*pOutDirection).length() > 0.0)
}

func (pSp *SurfacePoint) SurfacePointReflection(pInDirection *Vector3, pInRadiance *Color, pOutDirection *Vector3) *Color {
	inDot := pInDirection.DotProduct(pSp.pTriangle.normal)
	outDot := pOutDirection.DotProduct(pSp.pTriangle.normal)
	
	/* directions must be on same side of surface (no transmission) */
	isSameSide := !(((inDot < 0.0) || (outDot < 0.0)) && (!((inDot < 0.0) && (outDot < 0.0))))
	
	/* ideal diffuse BRDF:
      radiance scaled by reflectivity, cosine, and 1/pi  */
	r := ColorMultC(pInRadiance, &pSp.pTriangle.diffuse)
	return MultC(r, (math.Abs(inDot)/math.Pi)*(map[bool](float64){true:1.0,false:0.0}[isSameSide]))
}