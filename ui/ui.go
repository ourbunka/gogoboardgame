package ui

import (
	"gogoboardgame/board"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type UI struct {
	Name            string
	BackgroundImage *ebiten.Image
	Elements        []UIElement
}

type ButtonSelected int

const (
	Selected ButtonSelected = iota + 1
	Deselected
)

type UIElement struct {
	Name            string
	CurrentState    ButtonSelected
	IsSelectable    bool
	SelectableIndex int
	CurrentImage    *ebiten.Image
	SelectedImage   *ebiten.Image
	DeselectedImage *ebiten.Image
	PosX            float64
	PosY            float64
}

func PreloadUIAssets(chanNewUIAsset chan UI, screenWidth, screenHeight int) error {
	println("Preloading UI assets")
	var mainUI UI
	var pauseUI UI

	BtnStart, _, err := board.LoadImage("./resources/ui/btn_start.png")
	if err != nil {
		println(err.Error())
		return err
	}
	BtnStartDeselected, _, err := board.LoadImage("./resources/ui/btn_start_deselected.png")
	if err != nil {
		println(err.Error())
		return err
	}
	BtnResume, _, err := board.LoadImage("./resources/ui/btn_resume.png")
	if err != nil {
		println(err.Error())
		return err
	}
	BtnResumeDeselected, _, err := board.LoadImage("./resources/ui/btn_resume_deselected.png")
	if err != nil {
		println(err.Error())
		return err
	}
	BtnQuit, _, err := board.LoadImage("./resources/ui/btn_quit.png")
	if err != nil {
		println(err.Error())
		return err
	}
	BtnQuitDeselected, _, err := board.LoadImage("./resources/ui/btn_quit_deselected.png")
	if err != nil {
		println(err.Error())
		return err
	}
	Background, _, err := board.LoadImage("./resources/ui/bg.png")
	if err != nil {
		println(err.Error())
		return err
	}

	mainUI.BackgroundImage = Background
	mainUI.Name = "main menu"

	var BtnMainMenuStart UIElement
	BtnMainMenuStart.Name = "start button"
	BtnMainMenuStart.IsSelectable = true
	BtnMainMenuStart.CurrentState = Selected
	BtnMainMenuStart.PosX = float64(screenWidth)/2*4 - 512.0
	BtnMainMenuStart.PosY = float64(screenHeight)/2*4 - 180.0
	BtnMainMenuStart.CurrentImage = BtnStart
	BtnMainMenuStart.SelectedImage = BtnStart
	BtnMainMenuStart.DeselectedImage = BtnStartDeselected

	mainUI.Elements = append(mainUI.Elements, BtnMainMenuStart)
	var BtnMainMenuQuit UIElement
	BtnMainMenuQuit.Name = "quit button"
	BtnMainMenuQuit.IsSelectable = true
	BtnMainMenuQuit.CurrentState = Deselected
	BtnMainMenuQuit.PosX = float64(screenWidth)/2*4 - 512.0
	BtnMainMenuQuit.PosY = float64(screenHeight)/2*4 + 180.0 + 50.0
	BtnMainMenuQuit.CurrentImage = BtnQuitDeselected
	BtnMainMenuQuit.SelectedImage = BtnQuit
	BtnMainMenuQuit.DeselectedImage = BtnQuitDeselected

	mainUI.Elements = append(mainUI.Elements, BtnMainMenuQuit)

	sent := false
	for !sent {
		select {
		case chanNewUIAsset <- mainUI:
			println("SENT UI :", mainUI.Name)
			sent = true
		default:
			time.Sleep(time.Millisecond * 15)
			sent = false
		}
	}
	pauseUI.BackgroundImage = Background
	pauseUI.Name = "pause menu"
	sent = false
	var PauseMenuResume UIElement
	PauseMenuResume.Name = "resume button"
	PauseMenuResume.IsSelectable = true
	PauseMenuResume.CurrentState = Selected
	PauseMenuResume.PosX = float64(screenWidth)/2*4 - 512.0
	PauseMenuResume.PosY = float64(screenHeight)/2*4 - 180.0
	PauseMenuResume.CurrentImage = BtnResume
	PauseMenuResume.SelectedImage = BtnResume
	PauseMenuResume.DeselectedImage = BtnResumeDeselected

	pauseUI.Elements = append(pauseUI.Elements, PauseMenuResume)

	pauseUI.Elements = append(pauseUI.Elements, BtnMainMenuQuit)

	for !sent {
		select {
		case chanNewUIAsset <- pauseUI:
			println("SENT UI :", pauseUI.Name)
			sent = true
		default:
			time.Sleep(time.Millisecond * 15)
			sent = false
		}
	}

	return nil
}
