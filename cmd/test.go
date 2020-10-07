package main

import (
	"flag"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/orfby/ui/pkg/ui"
	"log"
	"net/http"
	"os"
	"runtime/pprof"
)

func main() {
	// Define the flags
	cpuProfile := flag.String("cpu-profile", "", "Where to output the CPU profile")
	width := flag.Float64("window-width", 800, "The window's width")
	height := flag.Float64("window-height", 600, "The window's height")

	// Parse them
	flag.Parse()

	// Check the flags
	if len(flag.Args()) == 0 {
		log.Fatal("Error: no path given")
	} else if len(flag.Args()) > 1 {
		log.Fatal("Error: multiple paths given")
	}

	// Get the path
	path := flag.Arg(0)

	// Open the ui assets folder
	uiDir := http.Dir("./assets/ui")

	var design *ui.Design

	// Run the design stuff on the pixelgl thread
	pixelgl.Run(func() {
		var err error

		// If the CPU profile was given
		if *cpuProfile != "" {
			// Create the file
			file, err := os.Create(*cpuProfile)
			if err != nil {
				log.Fatal(err)
			}
			// Start the profile
			err = pprof.StartCPUProfile(file)
			if err != nil {
				log.Fatal(err)
			}
			// Finish it after the program ends
			defer pprof.StopCPUProfile()
		}

		// Create a new design
		design, err = ui.NewDesign(uiDir, path,
			pixelgl.WindowConfig{
				Bounds: pixel.R(0, 0, *width, *height),
				Title:  path,
			})
		if err != nil {
			log.Fatal(err)
		}

		// Initialise the design
		err = design.Init()
		if err != nil {
			log.Fatal(err)
		}

		// Start then wait for it to finish
		design.StartThenWait()
	})
}
