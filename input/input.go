package input

import (
	"embed"
	"gogoboardgame/utils"
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type TouchInput struct {
	ShowTouchInput  bool
	OnScreenButtons []OnScreenButton
}

type OnScreenButton struct {
	Name string
	Img  *ebiten.Image
}

var Resources embed.FS
var UseEmbeded bool

func (ti TouchInput) ProcessTouchInput(screenWidth, screenHeight int) {
	for i, t := range ebiten.TouchIDs() {
		log.Println("PROCESSING TOUCH INDEX: ", i)
		tx, ty := ebiten.TouchPosition(t)
		switch {
		case screenWidth <= tx:

		// check menu min max
		case float64(tx) >= (float64(screenWidth)*0.025)-50.0 && float64(tx) <= (float64(screenWidth)*0.025)+150.0 &&
			float64(ty) >= (float64(screenHeight)*0.025)-50.0 && float64(ty) <= (float64(screenHeight)*0.025)+150.0:
			//detected menu
			log.Println("MENU")
		default:
		}

	}
}

func LoadOnScreenButton(chanNewOnScreenButton chan OnScreenButton,
	Build string) {
	var path string
	//movement button up
	path = "./resources/button/button_movement.png"
	img, _, err := utils.LoadImage(path, UseEmbeded, Build, Resources)
	if err != nil {
		log.Println(err.Error())
	}
	btn := OnScreenButton{
		Name: "movement_up",
		Img:  img,
	}
	var sent bool = false
	time.Sleep(time.Millisecond * 16)
	for !sent {
		select {
		case chanNewOnScreenButton <- btn:
			sent = true
		default:
			time.Sleep(time.Millisecond * 16)
			sent = false
		}
	}
	//movement button down, rotate this when draw()
	btn = OnScreenButton{
		Name: "movement_down",
		Img:  img,
	}
	chanNewOnScreenButton <- btn
	//movement button left, rotate this when draw()
	btn = OnScreenButton{
		Name: "movement_left",
		Img:  img,
	}
	chanNewOnScreenButton <- btn
	//movement button right, rotate this when draw()
	btn = OnScreenButton{
		Name: "movement_right",
		Img:  img,
	}
	sent = false
	time.Sleep(time.Millisecond * 16)
	for !sent {
		select {
		case chanNewOnScreenButton <- btn:
			sent = true
		default:
			time.Sleep(time.Millisecond * 16)
			sent = false
		}
	}

	//menu button
	path = "./resources/button/button_menu.png"
	img, _, err = utils.LoadImage(path, UseEmbeded, Build, Resources)
	if err != nil {
		log.Println(err.Error())
	}
	btn = OnScreenButton{
		Name: "button_menu",
		Img:  img,
	}
	sent = false
	time.Sleep(time.Millisecond * 16)
	for !sent {
		select {
		case chanNewOnScreenButton <- btn:
			sent = true
		default:
			time.Sleep(time.Millisecond * 16)
			sent = false
		}
	}

	//movement button place
	path = "./resources/button/button_add.png"
	img, _, err = utils.LoadImage(path, UseEmbeded, Build, Resources)
	if err != nil {
		log.Println(err.Error())
	}
	btn = OnScreenButton{
		Name: "button_place",
		Img:  img,
	}
	sent = false
	time.Sleep(time.Millisecond * 16)
	for !sent {
		select {
		case chanNewOnScreenButton <- btn:
			sent = true
		default:
			time.Sleep(time.Millisecond * 16)
			sent = false
		}
	}

	//movement button remove
	path = "./resources/button/button_minus.png"
	img, _, err = utils.LoadImage(path, UseEmbeded, Build, Resources)
	if err != nil {
		log.Println(err.Error())
	}
	btn = OnScreenButton{
		Name: "button_remove",
		Img:  img,
	}
	sent = false
	time.Sleep(time.Millisecond * 16)
	for !sent {
		select {
		case chanNewOnScreenButton <- btn:
			sent = true
		default:
			time.Sleep(time.Millisecond * 16)
			sent = false
		}
	}

}

func (ti *TouchInput) Draw(screen *ebiten.Image, screenHeight, screenWidth int) {
	for i, _ := range ti.OnScreenButtons {
		img := ti.OnScreenButtons[i].Img
		width, height := img.Size()
		dx := screenWidth - width
		dy := screenHeight - height
		op := &ebiten.DrawImageOptions{}
		op.ColorScale.Scale(1, 1, 1, 0.5)
		op.GeoM.Reset()
		op.GeoM.Scale(0.5, 0.5)
		switch i {
		case 0:
			//up
			op.GeoM.Translate(float64(dx+0*width-int(float64(screenWidth)*0.8)), float64(dy-256))
		case 1:
			//down rotate it 180deg and move downward
			theta := 180 * 3.141592 / 180
			op.GeoM.Rotate(theta)
			op.GeoM.Translate(float64(dx+0*width-int(float64(screenWidth)*0.8)+128), float64(dy+128))

		case 2:
			//Left rotate it 270deg
			theta := 270 * 3.141592 / 180
			op.GeoM.Rotate(theta)
			op.GeoM.Translate(float64(dx+0*width-int(float64(screenWidth)*0.8)-128), float64(dy))
		case 3:
			//Right rotate it 90deg
			theta := 90 * 3.141592 / 180
			op.GeoM.Rotate(theta)
			op.GeoM.Translate(float64(dx+0*width-int(float64(screenWidth)*0.8)+256), float64(dy-128))
		case 4:
			//menu
			op.GeoM.Translate((float64(screenWidth) * 0.025), float64(screenHeight)*0.025)
		case 5:
			//place stone
			op.GeoM.Translate((float64(screenWidth) * 0.8), float64(dy))
		case 6:
			//remove stone
			op.GeoM.Translate((float64(screenWidth)*0.8)+128+64, float64(dy))

		}

		//op.GeoM.Translate(float64(screenWidth/8+i*100-width/2), float64(screenHeight-screenHeight/5-height/2))
		screen.DrawImage(img, op)
	}
}
