package detector

type Detection struct {
	Signature string
	Score     int
}

type Result struct {
	Path       string
	Err        error
	Detections []Detection
}
