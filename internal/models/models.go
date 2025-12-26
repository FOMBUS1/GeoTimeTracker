package models

type GeoModel struct {
	RequestID      uint64
	UserID         uint64
	Longitude      float32
	Latitude       float32
	Departure      bool
	CustomLocation string
}
