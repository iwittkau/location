package forms

// Location is the form for locations
type Location struct {
	Name      string  `binding:"required"`
	Latitude  float64 `binding:"required"`
	Longitude float64 `binding:"required"`
}
