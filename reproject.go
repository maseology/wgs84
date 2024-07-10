package wgs84

func ReprojectMap(coords map[int][]float64, fromEpsg, toEpsg int) map[int][]float64 {
	m := make(map[int][]float64, len(coords))
	fromCode, toCode := EPSG().Code(fromEpsg), EPSG().Code(toEpsg)
	for k, cxy := range coords {
		longitude, latitude, _ := From(fromCode)(cxy[0], cxy[1], 0)
		m[k] = []float64{longitude, latitude}

		if toEpsg != 4326 {
			x, y, _ := To(toCode)(m[k][0], m[k][1], 0)
			m[k] = []float64{x, y}
		}
	}
	return m
}

func Reproject(x, y float64, fromEpsg, toEpsg int) (float64, float64) {
	fromCode, toCode := EPSG().Code(fromEpsg), EPSG().Code(toEpsg)
	longitude, latitude, _ := From(fromCode)(x, y, 0)
	if toEpsg != 4326 {
		x2, y2, _ := To(toCode)(longitude, latitude, 0)
		return x2, y2
	}
	return longitude, latitude
}

func ReprojectSlice(xys [][]float64, fromEpsg, toEpsg int) [][]float64 {
	xysnew := make([][]float64, len(xys))
	for i, xy := range xys {
		x, y := Reproject(xy[0], xy[1], fromEpsg, toEpsg)
		xysnew[i] = []float64{x, y}
	}
	return xysnew
}
