package util

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"log"
	"math"
)

//Function to draw a sprite onto the
//given canvas, according to the given
//gravity and scale
func DrawSprite(canvas *pixelgl.Canvas, sprite *pixel.Sprite,
	scale ScaleOption, gravity Gravity) {
	//If the sprite exists
	if sprite != nil {
		//If the scale is zero
		if scale == ZeroScaleOption {
			//Set the scale as the default
			scale = DefaultScaleOption
		}

		//If the sprite should repeat
		if scale == Tiled {
			//Get the size of the sprite
			spriteSize := sprite.Frame().Size()
			//Iterate over the y coords of each tile
			for y := spriteSize.Y / 2; y < canvas.Bounds().Max.Y; y += spriteSize.Y {
				//Iterate over the x coords of each tile
				for x := spriteSize.X / 2; x < canvas.Bounds().Max.X; x += spriteSize.X {
					mat := pixel.IM
					//Move the tile to the position
					mat = mat.Moved(pixel.V(x, y))
					//Draw it
					sprite.Draw(canvas, mat)
				}
			}
		} else {
			mat := pixel.IM
			//Switch over the scale options
			switch scale {
			case NoScale:
				//Nothing needs to be done
			case ScaleToFill:
				mat = mat.Scaled(pixel.ZV, math.Max(
					canvas.Bounds().Size().X/sprite.Frame().Size().X,
					canvas.Bounds().Size().Y/sprite.Frame().Size().Y))
			case ScaleToFit:
				mat = mat.Scaled(pixel.ZV, math.Min(
					canvas.Bounds().Size().X/sprite.Frame().Size().X,
					canvas.Bounds().Size().Y/sprite.Frame().Size().Y))
			case Stretch:
				mat = mat.ScaledXY(pixel.ZV, pixel.V(
					canvas.Bounds().Size().X/sprite.Frame().Size().X,
					canvas.Bounds().Size().Y/sprite.Frame().Size().Y))
			default:
				log.Printf("unknown scale option '%s'", scale)
			}
			//Move it to the center of the canvas
			mat = mat.Moved(canvas.Bounds().Center())
			//todo use gravity
			//Draw the sprite
			sprite.Draw(canvas, mat)
		}
	}
}
