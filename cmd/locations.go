package cmd

import (
	"encoding/json"
	"fmt"
	"image/color"
	"image/png"
	"io"
	"os"

	"github.com/mmcloughlin/globe"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Point is the simple model for a location on the surface of a sphere.
type Point struct {
	Latitude  float64 `json:"latitude"`  // The latitude of the point.
	Longitude float64 `json:"longitude"` // The longitude of the point.
}

// locationsCmd represents the "locations" CLI command.
var locationsCmd = &cobra.Command{
	Use:   "locations [files]",
	Short: "Draw a globe with locations as points read from files",
	Long: `Draw a globe with locations as points read from files.

The files need to contain json structured data compatible with

[
  {
    "longitude": 114.2,
    "latitude": 22.3,
  },
  {
    "longitude": 114.2,
    "latitude": 22.3,
  },
  ...
]

Example: draw Starbucks locations

 $ wget https://raw.githubusercontent.com/mmcloughlin/starbucks/master/locations.json
 $ globedraw locations locations.json -o starbucks_locations.png
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		g, err := createGlobe(cmd.Flags(), args)
		if err != nil {
			return fmt.Errorf("cannot create globe: %s", err)
		}
		err = writeGlobe(g, cmd.Flags())
		return err
	},
}

func writeGlobe(g *globe.Globe, flags *pflag.FlagSet) error {
	out, _ := flags.GetString("output")
	writer, err := openWriter(out)
	if err != nil {
		return fmt.Errorf("cannot open output '%s': %s", out, err)
	}
	defer writer.Close()

	size, _ := flags.GetInt("size")
	image := g.Image(size)
	return png.Encode(writer, image)
}

func createGlobe(flags *pflag.FlagSet, paths []string) (*globe.Globe, error) {
	g := globe.New()
	err := createGrid(g, flags)
	if err != nil {
		return nil, fmt.Errorf("cannot draw grid: %s", err)
	}
	err = appendPointsFromFiles(g, paths)
	if err != nil {
		return nil, fmt.Errorf("cannot draw points: %s", err)
	}
	lat, _ := flags.GetFloat64("center-latitude")
	lon, _ := flags.GetFloat64("center-longitude")
	g.CenterOn(lat, lon)
	return g, err
}

func appendPointsFromFiles(g *globe.Globe, paths []string) error {
	for _, path := range paths {
		err := appendPointsFromFile(g, path)
		if err != nil {
			return fmt.Errorf("error appending points from '%s': %s", path, err)
		}
	}
	return nil
}

func appendPointsFromFile(g *globe.Globe, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	pts := []Point{}
	err = decoder.Decode(&pts)
	if err != nil {
		return err
	}
	green := color.NRGBA{R: 0x00, G: 0x64, B: 0x3c, A: 192}
	for _, p := range pts {
		g.DrawDot(p.Latitude, p.Longitude, 0.025, globe.Color(green))
	}
	return nil
}

func createGrid(g *globe.Globe, flags *pflag.FlagSet) error {
	graticule, _ := flags.GetFloat64("graticule")
	if graticule <= 0 {
		return fmt.Errorf("invalid value: graticule must be greater than zero")
	}
	g.DrawGraticule(graticule)
	if drawLand, _ := flags.GetBool("land-boundaries"); drawLand {
		g.DrawLandBoundaries()
	}
	if drawCountry, _ := flags.GetBool("country-boundaries"); drawCountry {
		g.DrawCountryBoundaries()
	}
	return nil
}

func openWriter(out string) (io.WriteCloser, error) {
	if out == "" || out == "-" {
		return os.Stdout, nil
	}
	return os.Create(out)
}

func init() {
	locationsCmd.Flags().StringP("output", "o", "", "specify the output file, leave empty for stdout")
	locationsCmd.Flags().Float64P("graticule", "g", 10, "specify the graticule for the globe")
	locationsCmd.Flags().BoolP("land-boundaries", "l", true, "specify if land boundaries shall be drawn")
	locationsCmd.Flags().BoolP("country-boundaries", "c", true, "specify if country boundaries shall be drawn")
	locationsCmd.Flags().Float64P("center-latitude", "f", 51.509865, "specify the center latitude of the view")
	locationsCmd.Flags().Float64P("center-longitude", "t", -0.118092, "specify the center longitude of the view")
	locationsCmd.Flags().IntP("size", "s", 400, "specify the size of the image in pixels")

	RootCmd.AddCommand(locationsCmd)
}
