package game

import (
	"gogoboardgame/board"
	"gogoboardgame/input"
	"gogoboardgame/ui"
	"log"
	"math/rand/v2"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var timeSampler int
var debugUpdateTime bool
var update_counter int
var msg string
var chanString = make(chan string, 5)
var stateLen int
var chanNewState = make(chan []int, 3)
var chanNewHoverState = make(chan []int, 3)
var chanNewLibertiesResult = make(chan []board.LibertiesUpdateResult, 2)
var chanNewUIAsset = make(chan ui.UI, 3)
var chanNewOnScreenButton = make(chan input.OnScreenButton, 10)
var (
	screenWidth  = 1920
	screenHeight = 1080
)

const rendererScale float64 = 0.1

type GameState int

const (
	UIMainMenu GameState = iota + 1
	UIPauseMenu
	Gameplay
)

type Game struct {
	Game       ebiten.Game
	GameState  GameState
	boardSize  int // 9 = 9x9, 13 = 13x13, 19 = 19x19, etc
	board      board.Board
	LevelState string
	UIs        []ui.UI
	TouchInput input.TouchInput
}

func RunGame() {
	println("Starting... please wait")
	debugUpdateTime = false
	screenWidth, screenHeight = ebiten.Monitor().Size()
	ebiten.SetVsyncEnabled(true)
	ebiten.SetRunnableOnUnfocused(true)
	log.Printf("height : %d ", screenHeight)
	log.Printf("Width : %d ", screenWidth)
	ebiten.SetFullscreen(true)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Go")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}

func NewGame() *Game {
	var newGame Game
	if screenHeight == 2160 {
		newGame.boardSize = 19
	} else if screenHeight == 1440 {
		newGame.boardSize = 13
	} else {
		newGame.boardSize = 9
	}

	newGame.GameState = Gameplay
	newGame.TouchInput.ShowTouchInput = true
	newGame.board.NewBoard(screenWidth, screenHeight, newGame.boardSize, rendererScale)
	//go randomStoneLooper(newGame.boardSize * newGame.boardSize)
	go spawnHoverStone(newGame.boardSize * newGame.boardSize)
	go ui.PreloadUIAssets(chanNewUIAsset, screenWidth, screenHeight)
	if board.UseEmbeded == true {
		newGame.TouchInput.ShowTouchInput = true
		go input.LoadOnScreenButton(chanNewOnScreenButton, "ANDROID")
	} else {
		go input.LoadOnScreenButton(chanNewOnScreenButton, "DESKTOP")
	}

	return &newGame
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) Update() error {
	//println(g.GameState)
	start := time.Now()
	timeSampler++

	select {
	case newState := <-chanNewState:
		gridlen := len(g.board.Grid)
		newStatelen := len(newState)
		stateLen = gridlen
		//println("new state len :", gridlen, newStatelen)
		if newStatelen == g.boardSize*g.boardSize {
			for i := range g.board.Grid {
				g.board.Grid[i].State = newState[i]
				if g.board.Grid[i].State == 0 {
					g.board.Grid[i].Stone = nil
				}
				if g.board.Grid[i].State == 1 {
					g.board.Grid[i].Stone = g.board.StoneW
				}
				if g.board.Grid[i].State == 2 {
					g.board.Grid[i].Stone = g.board.StoneB
				}
				//println(g.board.Grid[i].State)
			}
		}
	default:
	}

	select {
	case newHoverState := <-chanNewHoverState:
		gridlen := len(g.board.Grid)
		newStatelen := len(newHoverState)
		stateLen = gridlen
		//println("new state len :", gridlen, newStatelen)
		if newStatelen == g.boardSize*g.boardSize {
			for i := range g.board.Grid {
				g.board.Grid[i].HoverState = newHoverState[i]
				if g.board.Grid[i].HoverState == 0 {
					g.board.Grid[i].HoverStone = nil
					g.board.Grid[i].IsStoneHover = false
				}
				if g.board.Grid[i].HoverState == 1 {
					g.board.Grid[i].HoverStone = nil
					g.board.Grid[i].IsStoneHover = false
				}
				if g.board.Grid[i].HoverState == 2 {
					g.board.Grid[i].HoverStone = nil
					g.board.Grid[i].IsStoneHover = false
				}
				if g.board.Grid[i].HoverState == 3 {
					g.board.Grid[i].HoverStone = g.board.StoneW
					g.board.Grid[i].IsStoneHover = true
				}
				if g.board.Grid[i].HoverState == 4 {
					g.board.Grid[i].HoverStone = g.board.StoneB
					g.board.Grid[i].IsStoneHover = true
				}
				//println(g.board.Grid[i].HoverState)
			}
		}
	default:
	}
	select {
	case newResult := <-chanNewLibertiesResult:
		for i := range newResult {
			g.board.Grid[i].Liberties = newResult[i].Liberties
			g.board.Grid[i].HasLiberties = newResult[i].HasLibertiesBool
			if g.board.Grid[i].Liberties == 0 {
				g.board.Grid[i].Stone = nil
				g.board.Grid[i].State = 0
			}
		}
		println("DONE UPDATING LIBERTIES RESULT")
	default:
	}

	select {
	case newUIAsset := <-chanNewUIAsset:
		g.UIs = append(g.UIs, newUIAsset)
		println("RECEIVED UI :", g.UIs[0].Name, newUIAsset.Name)

	default:
	}
	select {
	case newOnScreenButton := <-chanNewOnScreenButton:
		g.TouchInput.OnScreenButtons = append(g.TouchInput.OnScreenButtons, newOnScreenButton)
		println("RECEIVED OnScreenButton from channel")

	default:
	}

	switch g.GameState {
	case Gameplay:
		if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
			g.MoveUp()
		}

		if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
			g.MoveDown()
		}
		if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
			g.MoveLeft()
		}
		if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
			g.MoveRight()
		}

		if ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsKeyPressed(ebiten.KeyEnter) {
			g.PlaceStone()

		}
		if ebiten.IsKeyPressed(ebiten.KeyP) {
			timeNow := time.Now()
			if timeNow.Sub(g.board.LastMove) > time.Millisecond*200 {
				g.board.LastMove = timeNow
				var moved bool = false
				for i := range g.board.Grid {
					if g.board.Grid[i].IsStoneHover {
						if g.board.Turn == "black" && !moved {
							g.board.Turn = "white"
							g.board.Grid[i].HoverStone = g.board.StoneW
							g.board.Grid[i].HoverState = 3
							moved = true
						}
						if g.board.Turn == "white" && !moved {
							g.board.Turn = "black"
							g.board.Grid[i].HoverStone = g.board.StoneB
							g.board.Grid[i].HoverState = 4
							moved = true
						}
					}
				}
			}
		}
		if ebiten.IsKeyPressed(ebiten.KeyBackspace) {
			g.RemoveStone()
		}
		if ebiten.IsKeyPressed(ebiten.KeyEscape) || ebiten.IsKeyPressed(ebiten.KeyTab) || ebiten.IsKeyPressed(ebiten.KeyM) {
			g.TogglePauseMenu()
		}
	case UIMainMenu:
		//handle main menu ui

	case UIPauseMenu:
		//handle pause menu ui
		if ebiten.IsKeyPressed(ebiten.KeyW) && len(g.UIs) > 1 ||
			ebiten.IsKeyPressed(ebiten.KeyUp) && len(g.UIs) > 1 {
			g.UIUp()
		}
		if ebiten.IsKeyPressed(ebiten.KeyS) && len(g.UIs) > 1 ||
			ebiten.IsKeyPressed(ebiten.KeyDown) && len(g.UIs) > 1 {
			g.UIDown()
		}
		if ebiten.IsKeyPressed(ebiten.KeySpace) && len(g.UIs) > 1 ||
			ebiten.IsKeyPressed(ebiten.KeyEnter) && len(g.UIs) > 1 {
			err := g.UIConfirm()
			if err != nil {
				return err
			}
		}
		if ebiten.IsKeyPressed(ebiten.KeyTab) ||
			ebiten.IsKeyPressed(ebiten.KeyEscape) ||
			ebiten.IsKeyPressed(ebiten.KeyM) {
			g.TogglePauseMenu()
		}

	default:
	}

	if g.TouchInput.ShowTouchInput == true {
		g.TouchInput.Taps = g.TouchInput.Taps[:0]
		for id, t := range g.TouchInput.Touches {
			if inpututil.IsTouchJustReleased(id) {
				if g.TouchInput.Pinch != nil && (id == g.TouchInput.Pinch.Id1 || id == g.TouchInput.Pinch.Id2) {
					g.TouchInput.Pinch = nil
				}
				if g.TouchInput.Pan != nil && id == g.TouchInput.Pan.Id {
					g.TouchInput.Pan = nil
				}

				diff := input.Distance(t.OriginX, t.OriginY, t.CurrX, t.CurrY)
				if !t.WasPinch && !t.IsPan && (t.Duration <= 30 || diff < 2) {
					g.TouchInput.Taps = append(g.TouchInput.Taps, input.Tap{
						X: t.CurrX,
						Y: t.CurrY,
					})
				}
				delete(g.TouchInput.Touches, id)
			}
		}
		g.TouchInput.TouchIDs = inpututil.AppendJustPressedTouchIDs(g.TouchInput.TouchIDs[:0])
		for _, id := range g.TouchInput.TouchIDs {
			tx, ty := ebiten.TouchPosition(id)
			input := input.CalculateTouchInput(screenWidth, screenHeight, tx, ty)
			switch input {
			case "UP":
				if g.GameState == Gameplay {
					g.MoveUp()
				}
				if g.GameState == UIPauseMenu {
					g.UIUp()
				}
				if g.GameState == UIMainMenu {

				}
			case "DOWN":
				if g.GameState == Gameplay {
					g.MoveDown()
				}
				if g.GameState == UIPauseMenu {
					g.UIDown()
				}
				if g.GameState == UIMainMenu {

				}
			case "LEFT":
				if g.GameState == Gameplay {
					g.MoveLeft()
				}
				if g.GameState == UIPauseMenu {
					g.UIUp()
				}
				if g.GameState == UIMainMenu {

				}
			case "RIGHT":
				if g.GameState == Gameplay {
					g.MoveRight()
				}
				if g.GameState == UIPauseMenu {
					g.UIDown()
				}
				if g.GameState == UIMainMenu {

				}
			case "MENU":
				if g.GameState == Gameplay {
					g.TogglePauseMenu()
				}
				if g.GameState == UIPauseMenu {
					g.TogglePauseMenu()
				}
				if g.GameState == UIMainMenu {
				}
			case "ENTER":
				if g.GameState == Gameplay {
					g.PlaceStone()
				}
				if g.GameState == UIPauseMenu {
					g.UIConfirm()
				}
				if g.GameState == UIMainMenu {
				}
			case "REMOVE":
				if g.GameState == Gameplay {
					g.RemoveStone()
				}
				if g.GameState == UIPauseMenu {
				}
				if g.GameState == UIMainMenu {
				}
			default:
			}
		}

	}
	if timeSampler >= 600 {
		timeSampler = 0
		end := time.Now()
		duration := end.Sub(start).Milliseconds()
		go println("UPDATE() TOOK : ", duration, " MILLISEC")
	}

	return nil
}

func printDuration(start time.Time, end time.Time) {
	duration := end.Sub(start)
	println(duration.Milliseconds())
}

func (g *Game) Draw(screen *ebiten.Image) {
	//x, y := ebiten.CursorPosition()
	//println("   CURSOR : ", x, " , ", y)
	ebitenutil.DebugPrint(screen, "WIP. WASD/ARROW KEY TO MOVE, SPACEBAR/ENTER TO PLACE STONE, P KEY TO PASS, BACKSCAPE TO MANUALLY CAPTURE STONE.")
	for _, grid := range g.board.Grid {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(grid.PosX), float64(grid.PosY))
		op.GeoM.Scale(rendererScale, rendererScale)
		screen.DrawImage(grid.Image, op)
		if grid.Stone != nil {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(grid.PosX), float64(grid.PosY))
			op.GeoM.Scale(rendererScale, rendererScale)
			screen.DrawImage(grid.Stone, op)
		}
		if grid.IsStoneHover && grid.HoverStone != nil {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(grid.PosX-100), float64(grid.PosY-100))
			op.GeoM.Scale(rendererScale, rendererScale)
			op.ColorScale.ScaleAlpha(0.5)
			screen.DrawImage(grid.HoverStone, op)
		}
	}

	if g.GameState != Gameplay {
		if g.GameState == UIMainMenu && g.UIs[0].BackgroundImage != nil {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(0, 0)
			op.ColorScale.ScaleAlpha(0.55)
			screen.DrawImage(g.UIs[0].BackgroundImage, op)
			for _, ui := range g.UIs[1].Elements {
				//println("drawing UI :", ui.Name)
				if ui.CurrentImage != nil {
					op := &ebiten.DrawImageOptions{}
					op.GeoM.Translate(ui.PosX, ui.PosY)
					op.GeoM.Scale(0.25, 0.25)
					screen.DrawImage(ui.CurrentImage, op)
				}

			}
		}
		if len(g.UIs) > 1 {
			if g.GameState == UIPauseMenu && g.UIs[1].BackgroundImage != nil {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(0, 0)
				op.ColorScale.ScaleAlpha(0.55)
				screen.DrawImage(g.UIs[1].BackgroundImage, op)
				for _, ui := range g.UIs[1].Elements {
					//println("drawing UI :", ui.Name)
					if ui.CurrentImage != nil {
						op := &ebiten.DrawImageOptions{}
						op.GeoM.Translate(ui.PosX, ui.PosY)
						op.GeoM.Scale(0.25, 0.25)
						screen.DrawImage(ui.CurrentImage, op)
					}

				}
			}
		}
	}
	if g.TouchInput.ShowTouchInput == true && len(g.TouchInput.OnScreenButtons) >= 1 {
		g.TouchInput.Draw(screen, screenHeight, screenWidth)
	}
}

func randomStoneLooper(arrayLen int) {
	for range time.Tick(time.Second * 2) {
		if arrayLen != stateLen {
			arrayLen = stateLen
		}
		var newState []int
		for i := 0; i < arrayLen; i++ {
			number := rand.IntN(3)
			newState = append(newState, number)
		}
		sent := false
		for !sent {
			select {
			case chanNewState <- newState:
				sent = true
			default:
				time.Sleep(time.Millisecond * 20)
				sent = false
			}
		}
	}

}

func spawnHoverStone(arrayLen int) {
	// if arrayLen != stateLen {
	// 	arrayLen = stateLen
	// }
	var newGameState []int
	for i := 0; i < arrayLen; i++ {
		if i == ((arrayLen - 1) / 2) {
			startingLoc := 4 // 3 = hover white stone, 4 = hover black stone, which start first
			newGameState = append(newGameState, startingLoc)
		} else {
			startingEmptyState := 0
			newGameState = append(newGameState, startingEmptyState)
		}

	}
	var sent bool = false
	time.Sleep(time.Millisecond * 16)
	for !sent {
		select {
		case chanNewHoverState <- newGameState:
			sent = true
		default:
			time.Sleep(time.Millisecond * 16)
			sent = false
		}
	}
}
