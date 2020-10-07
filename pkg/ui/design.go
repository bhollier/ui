package ui

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	_ "github.com/orfby/ui/pkg/ui/builtin"
	"github.com/orfby/ui/pkg/ui/element"
	"log"
	"net/http"
	"sync"
	"time"
)

// A UI design, built from an XML file
type Design struct {
	sync.Mutex

	// The filesystem to use
	fs http.FileSystem

	// The design's window
	window *pixelgl.Window
	// The window's previous
	// bounds
	prevWindowBounds pixel.Rect

	// The path to the root element
	path string
	// The root element of the design
	root *element.Root

	// Condition variable for waiting
	// for the design to be closed
	waitCondVar *sync.Cond
}

// Function to create a new design from
// an XML string
func NewDesign(fs http.FileSystem, path string, windowConfig pixelgl.WindowConfig) (d *Design, err error) {
	// Create a new design struct
	d = new(Design)
	// The file system
	d.fs = fs
	// Create the condition variable
	d.waitCondVar = sync.NewCond(d)
	// The path
	d.path = path

	// Create the root
	d.root, err = element.NewRoot(fs, nil, path)
	if err != nil {
		return nil, err
	}

	// Create the window
	d.window, err = pixelgl.NewWindow(windowConfig)
	if err != nil {
		return nil, err
	}

	return
}

// Function to initialise (and draw) the
// design. This function must be called
// within pixelgl.Run
func (d *Design) Init() (err error) {
	// Update the root node
	return d.update(d.root)
}

// Function to start the design
func (d *Design) Start() { go d.pollEvents() }

// Function to wait for the design to close
func (d *Design) Wait() {
	d.Lock()
	// If the window is already closed
	if d.window.Closed() {
		return
	}
	// Wait for the design to close
	d.waitCondVar.Wait()
}

// Function to start the design then wait
// for it to finish
func (d *Design) StartThenWait() { d.Start(); d.Wait() }

// Function to get the design's window
func (d *Design) GetWindow() *pixelgl.Window { return d.window }

// Function to update the design
func (d *Design) update(root *element.Root) error {
	// Update the window's bounds
	d.prevWindowBounds = d.window.Bounds()

	// Initialise the design
	err := element.InitUI(root.Element, d.window, &d.prevWindowBounds)
	if err != nil {
		return err
	}

	// Draw the design
	element.DrawUI(root.Element, d.window)

	return nil
}

// Function to poll the events of the window
func (d *Design) pollEvents() {
	// While the window is still open
	for !d.window.Closed() {
		// Wait for a new event
		d.window.UpdateInputWait(time.Second)

		// Make sure the window is in focus
		if d.window.Focused() {
			// If ctrl + shift + r was pressed
			if d.window.Pressed(pixelgl.KeyLeftControl) &&
				d.window.Pressed(pixelgl.KeyLeftShift) &&
				d.window.JustPressed(pixelgl.KeyR) {
				// Create a new root root
				log.Printf("Loading XML design from '" + d.path + "'...")
				newRoot, err := element.NewRoot(d.fs, nil, d.path)
				if err != nil {
					log.Printf("Error reloading XML: %+v", err)
				}

				// Do an initial design update
				err = d.update(newRoot)
				if err != nil {
					log.Fatal(err)
				}

				// Set the new root
				d.Lock()
				d.root = newRoot
				d.Unlock()

				// If the window bounds changed
			} else if d.prevWindowBounds != d.window.Bounds() {
				// Update the design
				d.Lock()
				err := d.update(d.root)
				d.Unlock()
				if err != nil {
					log.Fatal(err)
				}
			}

			// Tell the root element
			// there was a new event
			d.Lock()
			/*go */
			d.root.NewEvent(d.window)

			// Draw the design
			element.DrawUI(d.root.Element, d.window)
			d.Unlock()
		}

		// Wait a bit before the next event
		// time.Sleep(time.Second / 50)
	}

	// Broadcast to any threads
	// waiting for the design to close
	d.waitCondVar.Broadcast()
}
