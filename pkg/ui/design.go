package ui

import (
	"errors"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	_ "github.com/orfby/ui/pkg/ui/builtin"
	"github.com/orfby/ui/pkg/ui/element"
	"log"
	"sync"
	"time"
)

//A UI design, built from an XML file
type Design struct {
	sync.Mutex

	//The design's window
	window *pixelgl.Window
	//The window's previous
	//bounds
	prevWindowBounds pixel.Rect
	//The window's canvas

	//The path to the root element
	path string
	//The root element of the design
	root *element.Root

	//Condition variable for waiting
	//for the design to be closed
	waitCondVar *sync.Cond
}

//Function to create a new design from
//an XML string. This function must be
//called within pixelgl.Run
func NewDesign(path string, windowConfig pixelgl.WindowConfig) (d *Design, err error) {
	//Create a new design struct
	d = new(Design)
	//Create the condition variable
	d.waitCondVar = sync.NewCond(d)
	//The path
	d.path = path

	//Create the root
	log.Printf("Loading XML design from '" + path + "'...")
	d.root, err = element.NewRoot(nil, path)
	if err != nil {
		return nil, err
	}

	//Create the window
	log.Print("Creating pixelgl.Window...")
	d.window, err = pixelgl.NewWindow(windowConfig)
	if err != nil {
		return nil, err
	}

	//Do an initial design update
	log.Printf("Initialising design...")
	err = d.update(d.root)

	//Start the routine to poll events
	log.Printf("Starting window event routine")
	go d.pollEvents()

	return
}

//Function to update the design
func (d *Design) update(root *element.Root) error {
	//Update the window's bounds
	d.prevWindowBounds = d.window.Bounds()

	//Reset the design (from the root)
	root.Reset()

	//Keep initialising the elements into they're all done
	//or until the root node has been initialised for the 1000th time
	for i := 0; i < 1000 && !root.IsInitialised(); i++ {
		//Initialise the root element (and therefore all its children)
		err := root.Init(d.window, &d.prevWindowBounds)
		if err != nil {
			return err
		}
	}

	//If the root still isn't initialised (because the
	//loop limit was reached) return an error
	if !root.IsInitialised() {
		//Make an array of uninitialised elements
		uninitialisedElements := make([]element.Element, 0)

		//Recursive function to find uninitialised elements
		var getUninitialisedElements func(element.Element)
		getUninitialisedElements = func(e element.Element) {
			//If the element isn't initialised
			if !e.IsInitialised() {
				//Add it to the array
				uninitialisedElements = append(uninitialisedElements, e)

				//Try to convert to a layout
				layout, ok := e.(element.Layout)
				//If it is a layout, iterate over the children
				if ok {
					for i := 0; i < layout.NumChildren(); i++ {
						//Get the child element's uninitialised elements
						getUninitialisedElements(layout.GetChild(i))
					}
				}
			}
		}

		//Call the function on root
		getUninitialisedElements(root.Element)
		//Create a string
		uninitialisedElementsStr := ""
		for i, elem := range uninitialisedElements {
			//Add the element to the string
			uninitialisedElementsStr +=
				element.FullName(elem, ".", true)
			//Add a comma unless this is the last element
			if i < len(uninitialisedElements)-1 {
				uninitialisedElementsStr += ", "
			}
		}

		//Return an error with the uninitialised elements
		return errors.New("infinite loop in element init detected. " +
			"The following element(s) are still uninitialised: " + uninitialisedElementsStr)
	}

	//Draw the design
	root.Draw()
	//Draw the root onto the window
	element.DrawCanvasOntoParent(root.Element.GetCanvas(), d.window.Canvas())
	//Swap the window's buffers
	d.window.SwapBuffers()

	return nil
}

//Function to poll the events of the window
func (d *Design) pollEvents() {
	//While the window is still open
	for !d.window.Closed() {
		//Wait for a new event
		d.window.UpdateInputWait(time.Second)

		//If ctrl + shift + r was pressed
		if d.window.Pressed(pixelgl.KeyLeftControl) &&
			d.window.Pressed(pixelgl.KeyLeftShift) &&
			d.window.JustPressed(pixelgl.KeyR) {
			//Create a new root root
			log.Printf("Loading XML design from '" + d.path + "'...")
			newRoot, err := element.NewRoot(nil, d.path)
			if err != nil {
				log.Printf("Error reloading XML: %+v", err)
			}

			//Do an initial design update
			log.Printf("Initialising design...")
			err = d.update(newRoot)

			//Set the new root
			d.root.Element = newRoot.Element

			//If the window bounds changed
		} else if d.prevWindowBounds != d.window.Bounds() {
			//todo  update periodically (in case the
			//todo  ui needs to be updated for some reason)
			//Update the design
			err := d.update(d.root)
			if err != nil {
				log.Fatalf("Fatal error: %+v", err)
			}
		}

		//Tell the root element
		//there was a new event
		/*go */
		d.root.NewEvent(d.window)

		//Wait a bit before the next event
		//time.Sleep(time.Second / 50)
	}

	//Broadcast to any threads
	//waiting for the design to close
	d.waitCondVar.Broadcast()
}

//Function to wait for the design to close
func (d *Design) Wait() error {
	//Wait for the design to close
	d.Lock()
	d.waitCondVar.Wait()
	return nil
}
