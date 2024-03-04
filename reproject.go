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
