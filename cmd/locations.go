package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"

	"github.com/mmcloughlin/globe"
	"io"
	"image/png"
)

// locationsCmd represents the locations command
var locationsCmd = &cobra.Command{
	Use:   "locations",
	Short: "TODO",
	Long: `TODO
.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		out, err := cmd.Flags().GetString("output")
		if err != nil {
			return fmt.Errorf("cannot determine output: %s", err)
		}
		g := globe.New()
		g.DrawGraticule(10.0)
		g.DrawLandBoundaries()
		g.DrawCountryBoundaries()
		g.CenterOn(51.453349, -2.588323)
		writer, err := openWriter(out)
		if err != nil {
			return fmt.Errorf("cannot open output '%s': %s", out, err)
		}
		defer writer.Close()
		image := g.Image(400)
		return png.Encode(writer, image)
	},

}

func openWriter(out string) (io.WriteCloser, error) {
	if out == "" {
		return os.Stdout, nil
	}
	return os.Create(out)
}

func init() {
	RootCmd.AddCommand(locationsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// locationsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	locationsCmd.Flags().StringP("output", "o", "", "Output")
}
