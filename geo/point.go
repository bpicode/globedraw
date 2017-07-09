package geo

// Point is the simple model for a location on the surface of a sphere.
type Point struct {
	Latitude  float64 `json:"latitude"`  // The latitude of the point.
	Longitude float64 `json:"longitude"` // The longitude of the point.
}
