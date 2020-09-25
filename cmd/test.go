package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/orfby/ui/pkg/ui"
	"github.com/orfby/ui/pkg/ui/element"
	"log"
)

func main() {
	//Run everything on pixelgl's thread
	pixelgl.Run(func() {
		//Register a callback
		element.RegisterCallback("testcallback", func(element.Element) error {
			log.Printf("*fart noise*")
			return nil
		})

		//Create a new test design
		design, err := ui.NewDesign("./assets/ui/designs/test.xml",
			pixelgl.WindowConfig{
				Bounds: pixel.R(0, 0, 800, 600),
				Title:  "test",
			})
		if err != nil {
			log.Fatalf("Fatal error: %+v", err)
		}

		//Wait for it to finish
		err = design.Wait()
		if err != nil {
			log.Fatalf("Fatal error: %+v", err)
		}
	})
}
