package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"

	"github.com/mmcloughlin/globe"
	"io"
	"image/png"
	"github.com/spf13/pflag"
)

type Point struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// locationsCmd represents the locations command
var locationsCmd = &cobra.Command{
	Use:   "locations",
	Short: "TODO",
	Long: `TODO
.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		g, err := createGlobe(cmd.Flags())
		if err != nil {
			return fmt.Errorf("cannot create globe: %s", err)
		}
		err = appendPointsFromFiles(g, args)
		if err != nil {
			return fmt.Errorf("cannot draw points: %s", err)
		}
		out, err := cmd.Flags().GetString("output")
		if err != nil {
			return fmt.Errorf("cannot determine output: %s", err)
		}
		writer, err := openWriter(out)
		if err != nil {
			return fmt.Errorf("cannot open output '%s': %s", out, err)
		}
		defer writer.Close()
		image := g.Image(400)
		return png.Encode(writer, image)
	},

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
	return nil
}

func createGlobe(flags *pflag.FlagSet) (*globe.Globe, error) {
	g := globe.New()
	graticule, _ := flags.GetFloat64("graticule")
	if graticule <= 0 {
		return nil, fmt.Errorf("invalid value: graticule must be greater than zero")
	}
	g.DrawGraticule(graticule)
	if drawLand, _ := flags.GetBool("land-boundaries"); drawLand {
		g.DrawLandBoundaries()
	}
	if drawCountry, _ := flags.GetBool("country-boundaries"); drawCountry {
		g.DrawCountryBoundaries()
	}
	lat, _ := flags.GetFloat64("center-latitude")
	lon, _ := flags.GetFloat64("center-longitude")
	g.CenterOn(lat, lon)
	return g, nil
}

func openWriter(out string) (io.WriteCloser, error) {
	if out == "" {
		return os.Stdout, nil
	}
	return os.Create(out)
}

func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// locationsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	locationsCmd.Flags().StringP("output", "o", "", "specify the output file, leave empty for stdout")
	locationsCmd.Flags().Float64P("graticule", "g", 10, "specify the graticule for the globe")
	locationsCmd.Flags().BoolP("land-boundaries", "l", true, "specify if land boundaries shall be drawn")
	locationsCmd.Flags().BoolP("country-boundaries", "c", true, "specify if country boundaries shall be drawn")
	locationsCmd.Flags().Float64P("center-latitude", "f", 51.509865, "specify the center latitude of the view")
	locationsCmd.Flags().Float64P("center-longitude", "t", -0.118092, "specify the center longitude of the view")

	RootCmd.AddCommand(locationsCmd)
}
