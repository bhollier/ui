package element

import (
	"encoding/xml"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/orfby/ui/pkg/ui/util"
	"log"
)

//A button's state
type ButtonState string

//The default button state
const ButtonDefaultState = "default"

//The button state when the mouse is over it
const ButtonHoveredState = "hovered"

//The button state for when the
//mouse is pressing the button
const ButtonPressedState = "pressed"

//Interface for a button element
type HasButton interface {
	//A button is an element
	Element

	//Function to get the button's
	//current state
	GetState() ButtonState
	//Function to set the button's
	//current state
	SetState(s ButtonState)

	//Function to get the button's
	//background field from XML for
	//the given state
	GetButtonBkgField(s ButtonState) string

	//Function to get the button's
	//background sprite for the given
	//state
	GetButtonBkg(s ButtonState) *pixel.Sprite
	//Function to set the button's
	//background sprite for the given
	//state
	SetButtonBkg(state ButtonState, sprite *pixel.Sprite)

	//Function to call the press callback
	CallPressCallback() error
}

//Type for a button. Note, structs
//must either include ButtonImpl or
//Impl, not both
type ButtonImpl struct {
	//The button is an element
	Impl

	//The button's current state
	state ButtonState

	//The button's background when
	//being hovered over from XML
	HoveredBackground string `uixml:"http://github.com/orfby/ui/api/schema bkg-hovered,optional"`
	//The button's background when
	//being pressed from XML
	PressedBackground string `uixml:"http://github.com/orfby/ui/api/schema bkg-pressed,optional"`

	//The button's background
	//sprites for each state
	backgrounds map[ButtonState]*pixel.Sprite

	//The element's press callback
	PressCallback string `uixml:"http://github.com/orfby/ui/api/schema press-callback,optional"`
}

//Function to create a button element
func NewButton(name xml.Name, parent Layout) ButtonImpl {
	e := ButtonImpl{Impl: NewElement(name, parent)}
	//Set the state as the default
	e.state = ButtonDefaultState
	//Create the backgrounds map
	e.backgrounds = map[ButtonState]*pixel.Sprite{
		ButtonDefaultState: nil,
		ButtonHoveredState: nil,
		ButtonPressedState: nil,
	}
	return e
}

//Function to get the button's
//current state
func (e *ButtonImpl) GetState() ButtonState { return e.state }

//Function to set the button's
//current state
func (e *ButtonImpl) SetState(s ButtonState) {
	e.state = s
	//Update the background
	e.SetBkgSprite(e.backgrounds[e.state])
}

//Function to get the button's
//background field from XML for
//the given state
func (e *ButtonImpl) GetButtonBkgField(s ButtonState) string {
	if s == ButtonDefaultState {
		return e.GetBkgField()
	} else if s == ButtonHoveredState {
		return e.HoveredBackground
	} else {
		return e.PressedBackground
	}
}

//Function to get the button's
//background sprite for the given
//state
func (e *ButtonImpl) GetButtonBkg(s ButtonState) *pixel.Sprite { return e.backgrounds[s] }

//Function to set the button's
//background sprite for the given
//state
func (e *ButtonImpl) SetButtonBkg(state ButtonState, sprite *pixel.Sprite) {
	e.backgrounds[state] = sprite
}

//Function to call the press callback
func (e *ButtonImpl) CallPressCallback() error {
	if e.PressCallback != "" {
		return Call(e.PressCallback, e)
	}
	return nil
}

//Function to determine whether
//the element is initialised
func (e *ButtonImpl) IsInitialised() bool {
	return e.Impl.IsInitialised() &&
		(e.HoveredBackground == "" || e.backgrounds[ButtonHoveredState] != nil) &&
		(e.PressedBackground == "" || e.backgrounds[ButtonPressedState] != nil)
}

//Function to initialise the element
func (e *ButtonImpl) Init(window *pixelgl.Window, bounds *pixel.Rect) error {
	//Initialise the element
	err := e.Impl.Init(window, bounds)
	if err != nil {
		return err
	}

	//If the bounds are known
	if bounds != nil {
		//If the default background hasn't been set
		if e.backgrounds[ButtonDefaultState] == nil {
			e.backgrounds[ButtonDefaultState] = e.Impl.GetBkgSprite()
		}
		//If the hovered background hasn't been made
		if e.backgrounds[ButtonHoveredState] == nil {
			//Load the background picture
			picture, err := util.CreatePictureFromField(e.HoveredBackground)
			if err != nil {
				return err
			}
			if picture != nil {
				//Create a sprite
				e.backgrounds[ButtonHoveredState] = pixel.NewSprite(picture, picture.Bounds())
			} else {
				e.backgrounds[ButtonHoveredState] = e.backgrounds[ButtonDefaultState]
			}
		}
		//If the pressed background hasn't been made
		if e.backgrounds[ButtonPressedState] == nil {
			//Load the background picture
			picture, err := util.CreatePictureFromField(e.PressedBackground)
			if err != nil {
				return err
			}
			if picture != nil {
				e.backgrounds[ButtonPressedState] = pixel.NewSprite(picture, picture.Bounds())
			} else {
				e.backgrounds[ButtonPressedState] = e.backgrounds[ButtonDefaultState]
			}
		}
	}
	return nil
}

//Function to handle a button's new event
func ButtonNewEvent(e HasButton, window *pixelgl.Window) {
	//Whether the button's state changed
	stateChange := false
	//If the mouse is actually in the window
	if window.MouseInsideWindow() &&
		e.GetCanvas().Bounds().Contains(window.MousePosition()) {
		//If the mouse button is being pressed
		if window.Pressed(pixelgl.MouseButtonLeft) {
			stateChange = e.GetState() != ButtonPressedState
			e.SetState(ButtonPressedState)

			//If the state changed
			if stateChange {
				//Call the press callback
				err := e.CallPressCallback()
				if err != nil {
					//could be better
					log.Printf("Error from button callback: %+v", err)
				}
			}
		} else {
			stateChange = e.GetState() != ButtonHoveredState
			e.SetState(ButtonHoveredState)
		}
	} else {
		stateChange = e.GetState() != ButtonDefaultState
		e.SetState(ButtonDefaultState)
	}

	//If the state was changed
	if stateChange {
		//Redraw the whole UI
		DrawUI(e, window)
	}
}
