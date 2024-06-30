package input

import (
	"embed"
	"gogoboardgame/utils"
	"log"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Touch struct {
	OriginX, OriginY int
	CurrX, CurrY     int
	Duration         int
	WasPinch, IsPan  bool
}
type Tap struct {
	X, Y int
}
type pinch struct {
	Id1, Id2 ebiten.TouchID
	originH  float64
	prevH    float64
}

type pan struct {
	Id ebiten.TouchID

	prevX, prevY     int
	originX, originY int
}

type TouchInput struct {
	X, Y float64
	Xoom float64

	TouchIDs        []ebiten.TouchID
	Touches         map[ebiten.TouchID]*Touch
	ShowTouchInput  bool
	OnScreenButtons []OnScreenButton
	Pinch           *pinch
	Pan             *pan
	Taps            []Tap
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
		log.Println("x: ", tx, " y: ", ty)
		switch {
		// check menu min max
		case float64(tx) >= (float64(screenWidth)*0.025)-25.0 && float64(tx) <= (float64(screenWidth)*0.025)+150.0 &&
			float64(ty) >= (float64(screenHeight)*0.025)-25.0 && float64(ty) <= (float64(screenHeight)*0.025)+150.0:
			//detected menu
			log.Println("MENU")
		case float64(tx) >= (float64(screenWidth)*0.1)-25.0 && float64(tx) <= (float64(screenWidth)*0.1)+150.0 &&
			float64(ty) >= (float64(screenHeight)*0.625)-25.0 && float64(ty) <= (float64(screenHeight)*0.625)+150.0:
			//detected menu
			log.Println("UP")
		case float64(tx) >= (float64(screenWidth)*0.1)-25.0 && float64(tx) <= (float64(screenWidth)*0.1)+150.0 &&
			float64(ty) >= float64(screenHeight)*0.825-25.0 && float64(ty) <= float64(screenHeight)*0.825+150:
			log.Println("DOWN")
		case float64(tx) >= (float64(screenWidth)*0.05)-25.0 && float64(tx) <= (float64(screenWidth)*0.05)+90.0 &&
			float64(ty) >= float64(screenHeight)*0.725-25.0 && float64(ty) <= float64(screenHeight)*0.725+150:
			log.Println("LEFT")
		case float64(tx) >= (float64(screenWidth)*0.15)-25.0 && float64(tx) <= (float64(screenWidth)*0.15)+90.0 &&
			float64(ty) >= float64(screenHeight)*0.725-25.0 && float64(ty) <= float64(screenHeight)*0.725+150:
			log.Println("RIGHT")
		default:
		}

	}
}

// distance between points a and b.
func Distance(xa, ya, xb, yb int) float64 {
	x := math.Abs(float64(xa - xb))
	y := math.Abs(float64(ya - yb))
	return math.Sqrt(x*x + y*y)
}

func CalculateTouchInput(screenWidth, screenHeight int, tx, ty int) string {
	log.Println("x: ", tx, " y: ", ty)
	switch {
	// check menu min max
	case float64(tx) >= (float64(screenWidth)*0.025)-25.0 && float64(tx) <= (float64(screenWidth)*0.025)+150.0 &&
		float64(ty) >= (float64(screenHeight)*0.025)-25.0 && float64(ty) <= (float64(screenHeight)*0.025)+150.0:
		//detected menu
		log.Println("MENU")
		return "MENU"
	case float64(tx) >= (float64(screenWidth)*0.1)-25.0 && float64(tx) <= (float64(screenWidth)*0.1)+150.0 &&
		float64(ty) >= (float64(screenHeight)*0.625)-25.0 && float64(ty) <= (float64(screenHeight)*0.625)+150.0:
		//detected menu
		log.Println("UP")
		return "UP"
	case float64(tx) >= (float64(screenWidth)*0.1)-25.0 && float64(tx) <= (float64(screenWidth)*0.1)+150.0 &&
		float64(ty) >= float64(screenHeight)*0.825-25.0 && float64(ty) <= float64(screenHeight)*0.825+150:
		log.Println("DOWN")
		return "DOWN"
	case float64(tx) >= (float64(screenWidth)*0.05)-25.0 && float64(tx) <= (float64(screenWidth)*0.05)+90.0 &&
		float64(ty) >= float64(screenHeight)*0.725-25.0 && float64(ty) <= float64(screenHeight)*0.725+150:
		log.Println("LEFT")
		return "LEFT"
	case float64(tx) >= (float64(screenWidth)*0.15)-25.0 && float64(tx) <= (float64(screenWidth)*0.15)+90.0 &&
		float64(ty) >= float64(screenHeight)*0.725-25.0 && float64(ty) <= float64(screenHeight)*0.725+150:
		log.Println("RIGHT")
		return "RIGHT"
	case float64(tx) >= (float64(screenWidth)*0.8)-25.0 && float64(tx) <= (float64(screenWidth)*0.8)+150.0 &&
		float64(ty) >= float64(screenHeight)*0.825-25.0 && float64(ty) <= float64(screenHeight)*0.825+150:
		log.Println("ENTER")
		return "ENTER"
	case float64(tx) >= (float64(screenWidth)*0.1)+128.0-25.0 && float64(tx) <= (float64(screenWidth)*0.1)+128.0+150.0 &&
		float64(ty) >= float64(screenHeight)*0.825-25.0 && float64(ty) <= float64(screenHeight)*0.825+150:
		log.Println("REMOVE")
		return "REMOVE"
	default:
		return "NONE"
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
		//width, height := img.Size()
		//_ = screenWidth - width
		//_ := screenHeight - height
		op := &ebiten.DrawImageOptions{}
		op.ColorScale.Scale(1, 1, 1, 0.5)
		op.GeoM.Reset()
		op.GeoM.Scale(0.5, 0.5)
		switch i {
		case 0:
			//up
			op.GeoM.Translate((float64(screenWidth) * 0.1), float64(screenHeight)*0.625)
			//op.GeoM.Translate(float64(dx+0*width-int(float64(screenWidth)*0.8)), float64(dy-256))
		case 1:
			//down rotate it 180deg and move downward
			theta := 180 * 3.141592 / 180
			op.GeoM.Rotate(theta)
			//offset +125 in x, +125 in y because of rotation
			op.GeoM.Translate((float64(screenWidth)*0.1)+(125.0), float64(screenHeight)*0.825+(125.0))
			//op.GeoM.Translate(float64(dx+0*width-int(float64(screenWidth)*0.8)+128), float64(dy+128))

		case 2:
			//Left rotate it 270deg
			theta := 270 * 3.141592 / 180
			op.GeoM.Rotate(theta)
			//offset +125 in y because of rotation
			op.GeoM.Translate((float64(screenWidth) * 0.05), float64(screenHeight)*0.725+125.0)
			//op.GeoM.Translate(float64(dx+0*width-int(float64(screenWidth)*0.8)-128), float64(dy))
		case 3:
			//Right rotate it 90deg
			theta := 90 * 3.141592 / 180
			op.GeoM.Rotate(theta)
			//offset +125 in x because of rotation
			op.GeoM.Translate((float64(screenWidth)*0.15 + 125.0), float64(screenHeight)*0.725)
			//op.GeoM.Translate(float64(dx+0*width-int(float64(screenWidth)*0.8)+256), float64(dy-128))
		case 4:
			//menu
			op.GeoM.Translate((float64(screenWidth) * 0.025), float64(screenHeight)*0.025)
		case 5:
			//place stone
			op.GeoM.Translate((float64(screenWidth) * 0.8), float64(screenHeight)*0.825)
		case 6:
			//remove stone
			op.GeoM.Translate((float64(screenWidth)*0.8)+128, float64(screenHeight)*0.825)

		}

		//op.GeoM.Translate(float64(screenWidth/8+i*100-width/2), float64(screenHeight-screenHeight/5-height/2))
		screen.DrawImage(img, op)
	}
}
