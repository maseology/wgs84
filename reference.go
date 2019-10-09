package wgs84

import (
	"math"
)

type XYZ struct {
	GeodeticDatum GeodeticDatum
}

func (crs XYZ) To(to CoordinateReferenceSystem) func(a, b, c float64) (a2, b2, c2 float64) {
	return Transform(crs, to)
}

func (crs XYZ) MajorAxis() float64 {
	return spheroid(crs.GeodeticDatum).MajorAxis()
}

func (crs XYZ) InverseFlattening() float64 {
	return spheroid(crs.GeodeticDatum).InverseFlattening()
}

func (crs XYZ) ToWGS84(x, y, z float64) (x0, y0, z0 float64) {
	return toWGS84(crs.GeodeticDatum, x, y, z)
}

func (crs XYZ) FromWGS84(x0, y0, z0 float64) (x, y, z float64) {
	return fromWGS84(crs.GeodeticDatum, x0, y0, z0)
}

func (crs XYZ) ToXYZ(a, b, c float64, gs GeodeticSpheroid) (x, y, z float64) {
	return a, b, c
}

func (crs XYZ) FromXYZ(x, y, z float64, gs GeodeticSpheroid) (a, b, c float64) {
	return x, y, z
}

type LonLat struct {
	GeodeticDatum GeodeticDatum
}

func (crs LonLat) To(to CoordinateReferenceSystem) func(a, b, c float64) (a2, b2, c2 float64) {
	return Transform(crs, to)
}

func (crs LonLat) MajorAxis() float64 {
	return spheroid(crs.GeodeticDatum).MajorAxis()
}

func (crs LonLat) InverseFlattening() float64 {
	return spheroid(crs.GeodeticDatum).InverseFlattening()
}

func (crs LonLat) ToWGS84(x, y, z float64) (x0, y0, z0 float64) {
	return toWGS84(crs.GeodeticDatum, x, y, z)
}

func (crs LonLat) FromWGS84(x0, y0, z0 float64) (x, y, z float64) {
	return fromWGS84(crs.GeodeticDatum, x0, y0, z0)
}

func (crs LonLat) ToXYZ(a, b, c float64, gs GeodeticSpheroid) (x, y, z float64) {
	s := spheroid(gs, crs.GeodeticDatum)
	x = (crs._N(radian(b), s) + c) * math.Cos(radian(a)) * math.Cos(radian(b))
	y = (crs._N(radian(b), s) + c) * math.Cos(radian(b)) * math.Sin(radian(a))
	z = (crs._N(radian(b), s)*math.Pow(s.MajorAxis()*(1-s.F()), 2)/(s.A2()) + c) * math.Sin(radian(b))
	return x, y, z
}

func (crs LonLat) FromXYZ(x, y, z float64, gs GeodeticSpheroid) (a, b, c float64) {
	s := spheroid(gs, crs.GeodeticDatum)
	sd := math.Sqrt(x*x + y*y)
	T := math.Atan(z * s.MajorAxis() / (sd * s.B()))
	B := math.Atan((z + s.E2()*(s.A2())/s.B()*
		math.Pow(math.Sin(T), 3)) / (sd - s.E2()*s.MajorAxis()*math.Pow(math.Cos(T), 3)))
	c = sd/math.Cos(B) - crs._N(B, s)
	a = degree(math.Atan2(y, x))
	b = degree(B)
	return a, b, c
}

func (crs LonLat) _N(φ float64, s Spheroid) float64 {
	return s.MajorAxis() / math.Sqrt(1-s.E2()*math.Pow(math.Sin(φ), 2))
}

type Projection struct {
	GeodeticDatum        GeodeticDatum
	CoordinateProjection CoordinateProjection
}

func (crs Projection) To(to CoordinateReferenceSystem) func(a, b, c float64) (a2, b2, c2 float64) {
	return Transform(crs, to)
}

func (crs Projection) MajorAxis() float64 {
	return spheroid(crs.GeodeticDatum).MajorAxis()
}

func (crs Projection) InverseFlattening() float64 {
	return spheroid(crs.GeodeticDatum).InverseFlattening()
}

func (crs Projection) ToWGS84(x, y, z float64) (x0, y0, z0 float64) {
	return toWGS84(crs.GeodeticDatum, x, y, z)
}

func (crs Projection) FromWGS84(x0, y0, z0 float64) (x, y, z float64) {
	return fromWGS84(crs.GeodeticDatum, x0, y0, z0)
}

func (crs Projection) ToXYZ(a, b, c float64, gs GeodeticSpheroid) (x, y, z float64) {
	if crs.CoordinateProjection == nil {
		return WebMercator{
			GeodeticDatum: crs.GeodeticDatum,
		}.ToXYZ(a, b, c, gs)
	}
	s := spheroid(gs, crs.GeodeticDatum)
	a, b = crs.CoordinateProjection.ToLonLat(a, b, s)
	return LonLat{
		GeodeticDatum: crs.GeodeticDatum,
	}.ToXYZ(a, b, c, s)
}

func (crs Projection) FromXYZ(x, y, z float64, gs GeodeticSpheroid) (a, b, c float64) {
	if crs.CoordinateProjection == nil {
		return WebMercator{
			GeodeticDatum: crs.GeodeticDatum,
		}.FromXYZ(x, y, z, gs)
	}
	s := spheroid(gs, crs.GeodeticDatum)
	a, b, c = LonLat{
		GeodeticDatum: crs.GeodeticDatum,
	}.FromXYZ(x, y, z, s)
	a, b = crs.CoordinateProjection.FromLonLat(a, b, s)
	return a, b, c
}

type WebMercator struct {
	GeodeticDatum GeodeticDatum
}

func (crs WebMercator) To(to CoordinateReferenceSystem) func(a, b, c float64) (a2, b2, c2 float64) {
	return Transform(crs, to)
}

func (crs WebMercator) MajorAxis() float64 {
	return spheroid(crs.GeodeticDatum).MajorAxis()
}

func (crs WebMercator) InverseFlattening() float64 {
	return spheroid(crs.GeodeticDatum).InverseFlattening()
}

func (crs WebMercator) ToWGS84(x, y, z float64) (x0, y0, z0 float64) {
	return toWGS84(crs.GeodeticDatum, x, y, z)
}

func (crs WebMercator) FromWGS84(x0, y0, z0 float64) (x, y, z float64) {
	return fromWGS84(crs.GeodeticDatum, x0, y0, z0)
}

func (crs WebMercator) ToXYZ(a, b, c float64, gs GeodeticSpheroid) (x, y, z float64) {
	s := spheroid(gs, crs.GeodeticDatum)
	return Projection{
		GeodeticDatum:        crs.GeodeticDatum,
		CoordinateProjection: crs,
	}.ToXYZ(a, b, c, s)
}

func (crs WebMercator) FromXYZ(x, y, z float64, gs GeodeticSpheroid) (a, b, c float64) {
	s := spheroid(gs, crs.GeodeticDatum)
	return Projection{
		GeodeticDatum:        crs.GeodeticDatum,
		CoordinateProjection: crs,
	}.FromXYZ(x, y, z, s)
}

func (crs WebMercator) ToLonLat(east, north float64, gs GeodeticSpheroid) (lon, lat float64) {
	s := spheroid(gs, crs.GeodeticDatum)
	lon = degree(east / s.MajorAxis())
	lat = math.Atan(math.Exp(north/s.MajorAxis()))*degree(1)*2 - 90
	return lon, lat
}

func (crs WebMercator) FromLonLat(lon, lat float64, gs GeodeticSpheroid) (east, north float64) {
	s := spheroid(gs, crs.GeodeticDatum)
	east = radian(lon) * s.MajorAxis()
	north = math.Log(math.Tan(radian((90+lat)/2))) * s.MajorAxis()
	return east, north
}

type Mercator struct {
	Lonf, Scale, Eastf, Northf float64
	GeodeticDatum              GeodeticDatum
}

func (crs Mercator) To(to CoordinateReferenceSystem) func(a, b, c float64) (a2, b2, c2 float64) {
	return Transform(crs, to)
}

func (crs Mercator) MajorAxis() float64 {
	return spheroid(crs.GeodeticDatum).MajorAxis()
}

func (crs Mercator) InverseFlattening() float64 {
	return spheroid(crs.GeodeticDatum).InverseFlattening()
}

func (crs Mercator) ToWGS84(x, y, z float64) (x0, y0, z0 float64) {
	return toWGS84(crs.GeodeticDatum, x, y, z)
}

func (crs Mercator) FromWGS84(x0, y0, z0 float64) (x, y, z float64) {
	return fromWGS84(crs.GeodeticDatum, x0, y0, z0)
}

func (crs Mercator) ToXYZ(a, b, c float64, gs GeodeticSpheroid) (x, y, z float64) {
	s := spheroid(gs, crs.GeodeticDatum)
	return Projection{
		GeodeticDatum:        crs.GeodeticDatum,
		CoordinateProjection: crs,
	}.ToXYZ(a, b, c, s)
}

func (crs Mercator) FromXYZ(x, y, z float64, gs GeodeticSpheroid) (a, b, c float64) {
	s := spheroid(gs, crs.GeodeticDatum)
	return Projection{
		GeodeticDatum:        crs.GeodeticDatum,
		CoordinateProjection: crs,
	}.FromXYZ(x, y, z, s)
}

func (crs Mercator) ToLonLat(east, north float64, gs GeodeticSpheroid) (lon, lat float64) {
	s := spheroid(gs, crs.GeodeticDatum)
	east = (east - crs.Eastf) / crs.Scale
	north = (north - crs.Northf) / crs.Scale
	t := math.Exp(-north * s.MajorAxis())
	φ := math.Pi/2 - 2*math.Atan(t)
	for i := 0; i < 5; i++ {
		φ = math.Pi/2 - 2*math.Atan(t*math.Pow((1-s.E()*math.Sin(φ))/(1+s.E()*math.Sin(φ)), s.E()/2))
	}
	return east/s.MajorAxis() + crs.Lonf, degree(φ)
}

func (crs Mercator) FromLonLat(lon, lat float64, gs GeodeticSpheroid) (east, north float64) {
	s := spheroid(gs, crs.GeodeticDatum)
	east = crs.Scale * s.MajorAxis() * (radian(lon) - radian(crs.Lonf))
	north = crs.Scale * s.MajorAxis() / 2 *
		math.Log(1+math.Sin(radian(lat))/(1-math.Sin(radian(lat)))*
			math.Pow((1-s.E()*math.Sin(radian(lat)))/(1+s.E()*math.Sin(radian(lat))), math.E))
	return east, north
}

func UTM(zone float64, northern bool) TransverseMercator {
	northf := 0.0
	if !northern {
		northf = 10000000
	}
	return TransverseMercator{
		Lonf:   zone*6 - 183,
		Latf:   0,
		Scale:  0.9996,
		Eastf:  500000,
		Northf: northf,
	}
}

type TransverseMercator struct {
	Lonf, Latf, Scale, Eastf, Northf float64
	GeodeticDatum                    GeodeticDatum
}

func (crs TransverseMercator) To(to CoordinateReferenceSystem) func(a, b, c float64) (a2, b2, c2 float64) {
	return Transform(crs, to)
}

func (crs TransverseMercator) MajorAxis() float64 {
	return spheroid(crs.GeodeticDatum).MajorAxis()
}

func (crs TransverseMercator) InverseFlattening() float64 {
	return spheroid(crs.GeodeticDatum).InverseFlattening()
}

func (crs TransverseMercator) ToWGS84(x, y, z float64) (x0, y0, z0 float64) {
	return toWGS84(crs.GeodeticDatum, x, y, z)
}

func (crs TransverseMercator) FromWGS84(x0, y0, z0 float64) (x, y, z float64) {
	return fromWGS84(crs.GeodeticDatum, x0, y0, z0)
}

func (crs TransverseMercator) ToXYZ(a, b, c float64, gs GeodeticSpheroid) (x, y, z float64) {
	s := spheroid(gs, crs.GeodeticDatum)
	return Projection{
		GeodeticDatum:        crs.GeodeticDatum,
		CoordinateProjection: crs,
	}.ToXYZ(a, b, c, s)
}

func (crs TransverseMercator) FromXYZ(x, y, z float64, gs GeodeticSpheroid) (a, b, c float64) {
	s := spheroid(gs, crs.GeodeticDatum)
	return Projection{
		GeodeticDatum:        crs.GeodeticDatum,
		CoordinateProjection: crs,
	}.FromXYZ(x, y, z, s)
}

func (crs TransverseMercator) ToLonLat(east, north float64, gs GeodeticSpheroid) (lon, lat float64) {
	s := spheroid(gs, crs.GeodeticDatum)
	east -= crs.Eastf
	north -= crs.Northf
	Mi := crs.M(radian(crs.Latf), s) + north/crs.Scale
	μ := Mi / (s.MajorAxis() * (1 - s.E2()/4 - 3*s.E4()/64 - 5*s.E6()/256))
	φ1 := μ + (3*s.Ei()/2-27*s.Ei3()/32)*math.Sin(2*μ) +
		(21*s.Ei2()/16-55*s.Ei4()/32)*math.Sin(4*μ) +
		(151*s.Ei3()/96)*math.Sin(6*μ) +
		(1097*s.Ei4()/512)*math.Sin(8*μ)
	R1 := s.MajorAxis() * (1 - s.E2()) / math.Pow(1-s.E2()*sin2(φ1), 3/2)
	D := east / (crs.N(φ1, s) * crs.Scale)
	φ := φ1 - (crs.N(φ1, s)*math.Tan(φ1)/R1)*(D*D/2-(5+3*crs.T(φ1)+10*crs.C(φ1, s)-4*crs.C(φ1, s)*crs.C(φ1, s)-9*s.Ei2())*
		math.Pow(D, 4)/24+(61+90*crs.T(φ1)+298*crs.C(φ1, s)+45*crs.T(φ1)*crs.T(φ1)-252*s.Ei2()-3*crs.C(φ1, s)*crs.C(φ1, s))*
		math.Pow(D, 6)/720)
	λ := radian(crs.Lonf) + (D-(1+2*crs.T(φ1)+crs.C(φ1, s))*D*D*D/6+(5-2*crs.C(φ1, s)+
		28*crs.T(φ1)-3*crs.C(φ1, s)*crs.C(φ1, s)+8*s.Ei2()+24*crs.T(φ1)*crs.T(φ1))*
		math.Pow(D, 5)/120)/math.Cos(φ1)
	return degree(λ), degree(φ)
}

func (crs TransverseMercator) FromLonLat(lon, lat float64, gs GeodeticSpheroid) (east, north float64) {
	s := spheroid(gs, crs.GeodeticDatum)
	φ := radian(lat)
	A := (radian(lon) - radian(crs.Lonf)) * math.Cos(φ)
	east = crs.Scale*crs.N(φ, s)*(A+(1-crs.T(φ)+crs.C(φ, s))*
		math.Pow(A, 3)/6+(5-18*crs.T(φ)+crs.T(φ)*crs.T(φ)+72*crs.C(φ, s)-58*s.Ei2())*
		math.Pow(A, 5)/120) + crs.Eastf
	north = crs.Scale*(crs.M(φ, s)-crs.M(radian(crs.Latf), s)+crs.N(φ, s)*math.Tan(φ)*
		(A*A/2+(5-crs.T(φ)+9*crs.C(φ, s)+4*crs.C(φ, s)*crs.C(φ, s))*
			math.Pow(A, 4)/24+(61-58*crs.T(φ)+crs.T(φ)*crs.T(φ)+600*crs.C(φ, s)-330*s.Ei2())*math.Pow(A, 6)/720)) + crs.Northf
	return east, north
}

func (TransverseMercator) M(φ float64, s Spheroid) float64 {
	return s.MajorAxis() * ((1-s.E2()/4-3*s.E4()/64-5*s.E6()/256)*φ -
		(3*s.E2()/8+3*s.E4()/32+45*s.E6()/1024)*math.Sin(2*φ) +
		(15*s.E4()/256+45*s.E6()/1024)*math.Sin(4*φ) -
		(35*s.E6()/3072)*math.Sin(6*φ))
}

func (TransverseMercator) N(φ float64, s Spheroid) float64 {
	return s.MajorAxis() / math.Sqrt(1-s.E2()*sin2(φ))
}

func (TransverseMercator) T(φ float64) float64 {
	return tan2(φ)
}

func (TransverseMercator) C(φ float64, s Spheroid) float64 {
	return s.Ei2() * cos2(φ)
}

type LambertConformalConic1SP struct {
	Lonf, Latf, Scale, Eastf, Northf float64
	GeodeticDatum                    GeodeticDatum
}

func (crs LambertConformalConic1SP) To(to CoordinateReferenceSystem) func(a, b, c float64) (a2, b2, c2 float64) {
	return Transform(crs, to)
}

func (crs LambertConformalConic1SP) MajorAxis() float64 {
	return spheroid(crs.GeodeticDatum).MajorAxis()
}

func (crs LambertConformalConic1SP) InverseFlattening() float64 {
	return spheroid(crs.GeodeticDatum).InverseFlattening()
}

func (crs LambertConformalConic1SP) ToWGS84(x, y, z float64) (x0, y0, z0 float64) {
	return toWGS84(crs.GeodeticDatum, x, y, z)
}

func (crs LambertConformalConic1SP) FromWGS84(x0, y0, z0 float64) (x, y, z float64) {
	return fromWGS84(crs.GeodeticDatum, x0, y0, z0)
}

func (crs LambertConformalConic1SP) ToXYZ(a, b, c float64, gs GeodeticSpheroid) (x, y, z float64) {
	s := spheroid(gs, crs.GeodeticDatum)
	return Projection{
		GeodeticDatum:        crs.GeodeticDatum,
		CoordinateProjection: crs,
	}.ToXYZ(a, b, c, s)
}

func (crs LambertConformalConic1SP) FromXYZ(x, y, z float64, gs GeodeticSpheroid) (a, b, c float64) {
	s := spheroid(gs, crs.GeodeticDatum)
	return Projection{
		GeodeticDatum:        crs.GeodeticDatum,
		CoordinateProjection: crs,
	}.FromXYZ(x, y, z, s)
}

func (crs LambertConformalConic1SP) ToLonLat(east, north float64, gs GeodeticSpheroid) (lon, lat float64) {
	s := spheroid(gs, crs.GeodeticDatum)
	ρi := math.Sqrt(math.Pow(east-crs.Eastf, 2) + math.Pow(crs.ρ(radian(crs.Latf), s)-(north-crs.Northf), 2))
	if crs.n() < 0 {
		ρi = -ρi
	}
	ti := math.Pow(ρi/(s.MajorAxis()*crs.Scale*crs.F(s)), 1/crs.n())
	φ := math.Pi/2 - 2*math.Atan(ti)
	for i := 0; i < 5; i++ {
		φ = math.Pi/2 - 2*math.Atan(ti*math.Pow((1-s.E()*math.Sin(φ))/(1+s.E()*math.Sin(φ)), s.E()/2))
	}
	λ := math.Atan((east-crs.Eastf)/(crs.ρ(radian(crs.Latf), s)-(north-crs.Northf)))/crs.n() + radian(crs.Lonf)
	return degree(λ), degree(φ)
}

func (crs LambertConformalConic1SP) FromLonLat(lon, lat float64, gs GeodeticSpheroid) (east, north float64) {
	s := spheroid(gs, crs.GeodeticDatum)
	θ := crs.n() * (radian(lon) - radian(crs.Lonf))
	east = crs.Eastf + crs.ρ(radian(lat), s)*math.Sin(θ)
	north = crs.Northf + crs.ρ(radian(crs.Latf), s) - crs.ρ(radian(lat), s)*math.Cos(θ)
	return east, north
}

func (LambertConformalConic1SP) t(φ float64, s Spheroid) float64 {
	return math.Tan(math.Pi/4-φ/2) /
		math.Pow((1-s.E()*math.Sin(φ))/(1+s.E()*math.Sin(φ)), s.E()/2)
}

func (LambertConformalConic1SP) m(φ float64, s Spheroid) float64 {
	return math.Cos(φ) / math.Sqrt(1-s.E2()*sin2(φ))
}

func (crs LambertConformalConic1SP) n() float64 {
	return math.Sin(radian(crs.Latf))
}

func (crs LambertConformalConic1SP) F(s Spheroid) float64 {
	return crs.m(radian(crs.Latf), s) / (crs.n() * math.Pow(crs.t(radian(crs.Latf), s), crs.n()))
}

func (crs LambertConformalConic1SP) ρ(φ float64, s Spheroid) float64 {
	return s.MajorAxis() * crs.F(s) * math.Pow(crs.t(φ, s)*crs.Scale, crs.n())
}

type LambertConformalConic2SP struct {
	Lonf, Latf, Lat1, Lat2, Eastf, Northf float64
	GeodeticDatum                         GeodeticDatum
}

func (crs LambertConformalConic2SP) To(to CoordinateReferenceSystem) func(a, b, c float64) (a2, b2, c2 float64) {
	return Transform(crs, to)
}

func (crs LambertConformalConic2SP) MajorAxis() float64 {
	return spheroid(crs.GeodeticDatum).MajorAxis()
}

func (crs LambertConformalConic2SP) InverseFlattening() float64 {
	return spheroid(crs.GeodeticDatum).InverseFlattening()
}

func (crs LambertConformalConic2SP) ToWGS84(x, y, z float64) (x0, y0, z0 float64) {
	return toWGS84(crs.GeodeticDatum, x, y, z)
}

func (crs LambertConformalConic2SP) FromWGS84(x0, y0, z0 float64) (x, y, z float64) {
	return fromWGS84(crs.GeodeticDatum, x0, y0, z0)
}

func (crs LambertConformalConic2SP) ToXYZ(a, b, c float64, gs GeodeticSpheroid) (x, y, z float64) {
	s := spheroid(gs, crs.GeodeticDatum)
	return Projection{
		GeodeticDatum:        crs.GeodeticDatum,
		CoordinateProjection: crs,
	}.ToXYZ(a, b, c, s)
}

func (crs LambertConformalConic2SP) FromXYZ(x, y, z float64, gs GeodeticSpheroid) (a, b, c float64) {
	s := spheroid(gs, crs.GeodeticDatum)
	return Projection{
		GeodeticDatum:        crs.GeodeticDatum,
		CoordinateProjection: crs,
	}.FromXYZ(x, y, z, s)
}

func (crs LambertConformalConic2SP) ToLonLat(east, north float64, gs GeodeticSpheroid) (lon, lat float64) {
	s := spheroid(gs, crs.GeodeticDatum)
	ρi := math.Sqrt(math.Pow(east-crs.Eastf, 2) + math.Pow(crs.ρ(radian(crs.Latf), s)-(north-crs.Northf), 2))
	if crs.n(s) < 0 {
		ρi = -ρi
	}
	ti := math.Pow(ρi/(s.MajorAxis()*crs.F(s)), 1/crs.n(s))
	φ := math.Pi/2 - 2*math.Atan(ti)
	for i := 0; i < 5; i++ {
		φ = math.Pi/2 - 2*math.Atan(ti*math.Pow((1-s.E()*math.Sin(φ))/(1+s.E()*math.Sin(φ)), s.E()/2))
	}
	λ := math.Atan((east-crs.Eastf)/(crs.ρ(radian(crs.Latf), s)-(north-crs.Northf)))/crs.n(s) + radian(crs.Lonf)
	return degree(λ), degree(φ)
}

func (crs LambertConformalConic2SP) FromLonLat(lon, lat float64, gs GeodeticSpheroid) (east, north float64) {
	s := spheroid(gs, crs.GeodeticDatum)
	θ := crs.n(s) * (radian(lon) - radian(crs.Lonf))
	east = crs.Eastf + crs.ρ(radian(lat), s)*math.Sin(θ)
	north = crs.Northf + crs.ρ(radian(crs.Latf), s) - crs.ρ(radian(lat), s)*math.Cos(θ)
	return east, north
}

func (crs LambertConformalConic2SP) t(φ float64, s Spheroid) float64 {
	return math.Tan(math.Pi/4-φ/2) /
		math.Pow((1-s.E()*math.Sin(φ))/(1+s.E()*math.Sin(φ)), s.E()/2)
}

func (crs LambertConformalConic2SP) m(φ float64, s Spheroid) float64 {
	return math.Cos(φ) / math.Sqrt(1-s.E2()*sin2(φ))
}

func (crs LambertConformalConic2SP) n(s Spheroid) float64 {
	if radian(crs.Lat1) == radian(crs.Lat2) {
		return math.Sin(radian(crs.Lat1))
	}
	return (math.Log(crs.m(radian(crs.Lat1), s)) - math.Log(crs.m(radian(crs.Lat2), s))) /
		(math.Log(crs.t(radian(crs.Lat1), s)) - math.Log(crs.t(radian(crs.Lat2), s)))
}

func (crs LambertConformalConic2SP) F(s Spheroid) float64 {
	return crs.m(radian(crs.Lat1), s) / (crs.n(s) * math.Pow(crs.t(radian(crs.Lat1), s), crs.n(s)))
}

func (crs LambertConformalConic2SP) ρ(φ float64, s Spheroid) float64 {
	return s.MajorAxis() * crs.F(s) * math.Pow(crs.t(φ, s), crs.n(s))
}

type AlbersEqualAreaConic struct {
	Lonf, Latf, Lat1, Lat2, Eastf, Northf float64
	GeodeticDatum                         GeodeticDatum
}

func (crs AlbersEqualAreaConic) To(to CoordinateReferenceSystem) func(a, b, c float64) (a2, b2, c2 float64) {
	return Transform(crs, to)
}

func (crs AlbersEqualAreaConic) MajorAxis() float64 {
	return spheroid(crs.GeodeticDatum).MajorAxis()
}

func (crs AlbersEqualAreaConic) InverseFlattening() float64 {
	return spheroid(crs.GeodeticDatum).InverseFlattening()
}

func (crs AlbersEqualAreaConic) ToWGS84(x, y, z float64) (x0, y0, z0 float64) {
	return toWGS84(crs.GeodeticDatum, x, y, z)
}

func (crs AlbersEqualAreaConic) FromWGS84(x0, y0, z0 float64) (x, y, z float64) {
	return fromWGS84(crs.GeodeticDatum, x0, y0, z0)
}

func (crs AlbersEqualAreaConic) ToXYZ(a, b, c float64, gs GeodeticSpheroid) (x, y, z float64) {
	s := spheroid(gs, crs.GeodeticDatum)
	return Projection{
		GeodeticDatum:        crs.GeodeticDatum,
		CoordinateProjection: crs,
	}.ToXYZ(a, b, c, s)
}

func (crs AlbersEqualAreaConic) FromXYZ(x, y, z float64, gs GeodeticSpheroid) (a, b, c float64) {
	s := spheroid(gs, crs.GeodeticDatum)
	return Projection{
		GeodeticDatum:        crs.GeodeticDatum,
		CoordinateProjection: crs,
	}.FromXYZ(x, y, z, s)
}

func (crs AlbersEqualAreaConic) ToLonLat(east, north float64, gs GeodeticSpheroid) (lon, lat float64) {
	s := spheroid(gs, crs.GeodeticDatum)
	east -= crs.Eastf
	north -= crs.Northf
	ρi := math.Sqrt(east*east + math.Pow(crs.ρ(radian(crs.Latf), s)-north, 2))
	qi := (crs.C(s) - ρi*ρi*crs.n(s)*crs.n(s)/s.A2()) / crs.n(s)
	φ := math.Asin(qi / 2)
	for i := 0; i < 5; i++ {
		φ += math.Pow(1-s.E2()*sin2(φ), 2) /
			(2 * math.Cos(φ)) * (qi/(1-s.E2()) -
			math.Sin(φ)/(1-s.E2()*sin2(φ)) +
			1/(2*s.E())*math.Log((1-s.E()*math.Sin(φ))/(1+s.E()*math.Sin(φ))))
	}
	θ := math.Atan(east / (crs.ρ(radian(crs.Latf), s) - north))
	return degree(radian(crs.Lonf) + θ/crs.n(s)), degree(φ)
}

func (crs AlbersEqualAreaConic) FromLonLat(lon, lat float64, gs GeodeticSpheroid) (east, north float64) {
	s := spheroid(gs, crs.GeodeticDatum)
	θ := crs.n(s) * (radian(lon) - radian(crs.Lonf))
	east = crs.Eastf + crs.ρ(radian(lat), s)*math.Sin(θ)
	north = crs.Northf + crs.ρ(radian(crs.Latf), s) - crs.ρ(radian(lat), s)*math.Cos(θ)
	return east, north
}

func (crs AlbersEqualAreaConic) m(φ float64, s Spheroid) float64 {
	return math.Cos(φ) / math.Sqrt(1-s.E2()*sin2(φ))
}

func (crs AlbersEqualAreaConic) q(φ float64, s Spheroid) float64 {
	return (1 - s.E2()) * (math.Sin(φ)/(1-s.E2()*sin2(φ)) -
		(1/(2*s.E()))*math.Log((1-s.E()*math.Sin(φ))/(1+s.E()*math.Sin(φ))))
}

func (crs AlbersEqualAreaConic) n(s Spheroid) float64 {
	if radian(crs.Lat1) == radian(crs.Lat2) {
		return math.Sin(radian(crs.Lat1))
	}
	return (crs.m(radian(crs.Lat1), s)*crs.m(radian(crs.Lat1), s) - crs.m(radian(crs.Lat2), s)*crs.m(radian(crs.Lat2), s)) /
		(crs.q(radian(crs.Lat2), s) - crs.q(radian(crs.Lat1), s))
}

func (crs AlbersEqualAreaConic) C(s Spheroid) float64 {
	return crs.m(radian(crs.Lat1), s)*crs.m(radian(crs.Lat1), s) + crs.n(s)*crs.q(radian(crs.Lat1), s)
}

func (crs AlbersEqualAreaConic) ρ(φ float64, s Spheroid) float64 {
	return s.MajorAxis() * math.Sqrt(crs.C(s)-crs.n(s)*crs.q(φ, s)) / crs.n(s)
}

type EquidistantConic struct {
	Lonf, Latf, Lat1, Lat2, Eastf, Northf float64
	GeodeticDatum                         GeodeticDatum
}

func (crs EquidistantConic) To(to CoordinateReferenceSystem) func(a, b, c float64) (a2, b2, c2 float64) {
	return Transform(crs, to)
}

func (crs EquidistantConic) MajorAxis() float64 {
	return spheroid(crs.GeodeticDatum).MajorAxis()
}

func (crs EquidistantConic) InverseFlattening() float64 {
	return spheroid(crs.GeodeticDatum).InverseFlattening()
}

func (crs EquidistantConic) ToWGS84(x, y, z float64) (x0, y0, z0 float64) {
	return toWGS84(crs.GeodeticDatum, x, y, z)
}

func (crs EquidistantConic) FromWGS84(x0, y0, z0 float64) (x, y, z float64) {
	return fromWGS84(crs.GeodeticDatum, x0, y0, z0)
}

func (crs EquidistantConic) ToXYZ(a, b, c float64, gs GeodeticSpheroid) (x, y, z float64) {
	s := spheroid(gs, crs.GeodeticDatum)
	return Projection{
		GeodeticDatum:        crs.GeodeticDatum,
		CoordinateProjection: crs,
	}.ToXYZ(a, b, c, s)
}

func (crs EquidistantConic) FromXYZ(x, y, z float64, gs GeodeticSpheroid) (a, b, c float64) {
	s := spheroid(gs, crs.GeodeticDatum)
	return Projection{
		GeodeticDatum:        crs.GeodeticDatum,
		CoordinateProjection: crs,
	}.FromXYZ(x, y, z, s)
}

func (crs EquidistantConic) ToLonLat(east, north float64, gs GeodeticSpheroid) (lon, lat float64) {
	s := spheroid(gs, crs.GeodeticDatum)
	east -= crs.Eastf
	north -= crs.Northf
	ρi := math.Sqrt(east*east + math.Pow(crs.ρ(radian(crs.Latf), s)-north, 2))
	if crs.n(s) < 0 {
		ρi = -ρi
	}
	Mi := s.MajorAxis()*crs.G(s) - ρi
	μ := Mi / (s.MajorAxis() * (1 - s.E2()/4 - 3*s.E4()/64 - 5*s.E6()/256))
	φ := μ + (3*s.Ei()/2-27*s.Ei3()/32)*math.Sin(2*μ) +
		(21*s.Ei2()/16-55*s.Ei4()/32)*math.Sin(4*μ) +
		(151*s.Ei3()/96)*math.Sin(6*μ) +
		(1097*s.Ei4()/512)*math.Sin(8*μ)
	θ := math.Atan(east / (crs.ρ(radian(crs.Latf), s) - north))
	return degree((radian(crs.Lonf) + θ/crs.n(s))), degree(φ)
}

func (crs EquidistantConic) FromLonLat(lon, lat float64, gs GeodeticSpheroid) (east, north float64) {
	s := spheroid(gs, crs.GeodeticDatum)
	θ := crs.n(s) * (radian(lon) - radian(crs.Lonf))
	east = crs.Eastf + crs.ρ(radian(lat), s)*math.Sin(θ)
	north = crs.Northf + crs.ρ(radian(crs.Latf), s) - crs.ρ(radian(lat), s)*math.Cos(θ)
	return east, north
}

func (crs EquidistantConic) M(φ float64, s Spheroid) float64 {
	return s.MajorAxis() * ((1-s.E2()/4-3*s.E4()/64-5*s.E6()/256)*φ -
		(3*s.E2()/8+3*s.E4()/32+45*s.E6()/1024)*math.Sin(2*φ) +
		(15*s.E4()/256+45*s.E6()/1024)*math.Sin(4*φ) -
		(35*s.E6()/3072)*math.Sin(6*φ))
}

func (crs EquidistantConic) m(φ float64, s Spheroid) float64 {
	return math.Cos(φ) / math.Sqrt(1-s.E2()*sin2(φ))
}

func (crs EquidistantConic) n(s Spheroid) float64 {
	if radian(crs.Lat1) == radian(crs.Lat2) {
		return math.Sin(radian(crs.Lat1))
	}
	return s.MajorAxis() * (crs.m(radian(crs.Lat1), s) - crs.m(radian(crs.Lat2), s)) / (crs.M(radian(crs.Lat2), s) - crs.M(radian(crs.Lat1), s))
}

func (crs EquidistantConic) G(s Spheroid) float64 {
	return crs.m(radian(crs.Lat1), s)/crs.n(s) + crs.M(radian(crs.Lat1), s)/s.MajorAxis()
}

func (crs EquidistantConic) ρ(φ float64, s Spheroid) float64 {
	return s.MajorAxis()*crs.G(s) - crs.M(φ, s)
}