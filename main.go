package main

import (
	"fmt"
	"log"

	"github.com/faiface/pixel/pixelgl"
	"github.com/spf13/cobra"
)

var drawCmd = &cobra.Command{
	Use:   "fractal",
	Short: "Draw a fractal image",
	Long:  `Draw a fractal image`,
	Run: func(cmd *cobra.Command, args []string) {
		err := Image()
		if err != nil {
			log.Fatal(err)
		}
	}}

var animCmd = &cobra.Command{
	Use:   "anim",
	Short: "Draw a fractal image",
	Long:  `Draw a fractal image`,
	Run: func(cmd *cobra.Command, args []string) {
		err := Anim()
		if err != nil {
			log.Fatal(err)
		}
	}}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "draw",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		switch cmd.CalledAs() {
		case "init", "version":
			return
		}
	},
}

func main() {
	fmt.Println("~~~~~~~~~~~~~~~")
	fmt.Printf("window: %v\n", window)
	if window {
		pixelgl.Run(run)
	} else {
		run()
	}

}

func run() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.AddCommand(drawCmd)
	rootCmd.AddCommand(animCmd)

	// local flags
	drawCmd.Flags().IntVarP(&width, "width", "W", 1000, "The width of the output image")
	drawCmd.Flags().IntVarP(&height, "height", "H", 700, "The height of the output image")
	drawCmd.Flags().IntVarP(&maxIterations, "maxIteration", "i", 100, "The maximum number of iterations to cycle the fractal")
	drawCmd.Flags().Float64VarP(&scale, "scale", "s", 50, "The scale of the generated image. (level of zoom within fractal)")
	drawCmd.Flags().StringVarP(&outLoc, "outputLocation", "o", "./output/output.png", "The output path. (must also contain final filename)")
	drawCmd.Flags().BoolVarP(&window, "window", "w", false, "Output window")

	animCmd.Flags().IntVarP(&width, "width", "W", 1000, "The width of the output image")
	animCmd.Flags().IntVarP(&height, "height", "H", 700, "The height of the output image")
	animCmd.Flags().IntVarP(&maxIterations, "maxIteration", "i", 100, "The maximum number of iterations to cycle the fractal")
	animCmd.Flags().IntVarP(&frameCount, "frameCount", "f", 10, "The maximum number of frames to animate")
	animCmd.Flags().Float64VarP(&scale, "scale", "s", 50, "The scale of the generated image. (level of zoom within fractal)")
	animCmd.Flags().Float64VarP(&depth, "depth", "d", 2.5, "The depth of the generated image. (level of zoom within fractal)")
	animCmd.Flags().StringVarP(&outLoc, "outputLocation", "o", "./output/output.gif", "The output path. (must also contain final filename)")
	animCmd.Flags().BoolVarP(&window, "window", "w", false, "Output window")

}
