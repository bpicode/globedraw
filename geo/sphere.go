package geo

import (
	"image/jpeg"
	"image/png"
	"io"

	"github.com/mmcloughlin/globe"
)

type Sphere struct {
	*globe.Globe
}

func NewSphere() *Sphere {
	s := &Sphere{globe.New()}
	return s
}

func (s *Sphere) EncodePNG(size int, writer io.Writer) error {
	image := s.Image(size)
	return png.Encode(writer, image)
}

func (s *Sphere) EncodeJPEG(size int, quality int, writer io.Writer) error {
	image := s.Image(size)
	return jpeg.Encode(writer, image, &jpeg.Options{Quality: quality})
}
