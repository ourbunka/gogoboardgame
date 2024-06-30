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

	newGame.GameState = UIPauseMenu
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

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) || inpututil.IsKeyJustPressed(ebiten.KeyTab) || inpututil.IsKeyJustPressed(ebiten.KeyM) {
		timeNow := time.Now()
		if timeNow.Sub(g.board.LastMove) > time.Millisecond*100 {
			g.board.LastMove = timeNow
			var moved bool = false
			if g.GameState == UIMainMenu && !moved {
				board.ToDo()
			}

			//remove ui and return to gameplay while in ui gamestate
			if g.GameState == UIPauseMenu && !moved {
				g.GameState = Gameplay
				println("ESC : ", g.GameState)
				moved = true
			}
			//spawn pause menu while in gameplay gamestate
			if g.GameState == Gameplay && !moved {
				g.GameState = UIPauseMenu
				println("ESC : ", g.GameState)
				moved = true
			}
		}

		//return ebiten.Termination
	}

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

	if g.GameState == Gameplay {
		if inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyW) {
			timeNow := time.Now()
			if timeNow.Sub(g.board.LastMove) > time.Millisecond*120 {
				g.board.LastMove = timeNow
				for i := range g.board.Grid {
					if g.board.Grid[i].IsStoneHover && i != 0 && g.board.Grid[i].Col != 0 {
						g.board.Grid[i].IsStoneHover = false
						g.board.Grid[i].HoverStone = nil
						if i != 0 {
							g.board.Grid[i-1].IsStoneHover = true
						}
						if g.board.Turn == "black" && i != 0 {
							g.board.Grid[i-1].HoverStone = g.board.StoneB
						} else if i != 0 {
							g.board.Grid[i-1].HoverStone = g.board.StoneW
						}
					}
				}
			}

		}

		if inpututil.IsKeyJustPressed(ebiten.KeyDown) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
			timeNow := time.Now()
			if timeNow.Sub(g.board.LastMove) > time.Millisecond*120 {
				g.board.LastMove = timeNow
				var moved bool = false
				for i := range g.board.Grid {
					if !moved {
						if g.board.Grid[i].IsStoneHover && i != len(g.board.Grid)-1 && g.board.Grid[i].Col != g.boardSize-1 {
							g.board.Grid[i].IsStoneHover = false
							g.board.Grid[i].HoverStone = nil
							g.board.Grid[i+1].IsStoneHover = true
							if g.board.Turn == "black" {
								g.board.Grid[i+1].HoverStone = g.board.StoneB
							} else {
								g.board.Grid[i+1].HoverStone = g.board.StoneW
							}
							moved = true

						}
					}

				}
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
			timeNow := time.Now()
			if timeNow.Sub(g.board.LastMove) > time.Millisecond*120 {
				g.board.LastMove = timeNow
				var moved bool = false
				for i := range g.board.Grid {
					if !moved {
						if g.board.Grid[i].IsStoneHover && i != 0 && g.board.Grid[i].Row != 0 {
							g.board.Grid[i].IsStoneHover = false
							g.board.Grid[i].HoverStone = nil
							g.board.Grid[i-g.boardSize].IsStoneHover = true
							if g.board.Turn == "black" {
								g.board.Grid[i-g.boardSize].HoverStone = g.board.StoneB
							} else {
								g.board.Grid[i-g.boardSize].HoverStone = g.board.StoneW
							}
							moved = true

						}
					}

				}
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) {
			timeNow := time.Now()
			if timeNow.Sub(g.board.LastMove) > time.Millisecond*120 {
				g.board.LastMove = timeNow
				var moved bool = false
				for i := range g.board.Grid {
					if !moved {
						if g.board.Grid[i].IsStoneHover && i != len(g.board.Grid)-1 && g.board.Grid[i].Row != g.boardSize-1 {
							g.board.Grid[i].IsStoneHover = false
							g.board.Grid[i].HoverStone = nil
							g.board.Grid[i+g.boardSize].IsStoneHover = true
							if g.board.Turn == "black" {
								g.board.Grid[i+g.boardSize].HoverStone = g.board.StoneB
							} else {
								g.board.Grid[i+g.boardSize].HoverStone = g.board.StoneW
							}
							moved = true
						}
					}

				}
			}
		}

		if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			//prototype todo check liberties,check if can place stone, check capture etc
			timeNow := time.Now()
			if timeNow.Sub(g.board.LastMove) > time.Millisecond*120 {
				g.board.LastMove = timeNow
				moved := false
				for i := range g.board.Grid {
					if g.board.Grid[i].IsStoneHover && g.board.Grid[i].Stone == nil {
						currentTurn := g.board.Turn
						var currentTurnState int
						if currentTurn == "black" {
							currentTurnState = 2
						} else {
							currentTurnState = 1
						}
						hasLiberties := g.board.CheckLiberties(currentTurnState, i, g.boardSize)
						if hasLiberties {
							if g.board.Turn == "black" && !moved {
								g.board.Grid[i].Stone = g.board.StoneB
								g.board.Grid[i].State = 2
								g.board.LastPlacedIndex = i
								moved = true
								g.board.CanCurrentPlacedStoneCapture(chanNewLibertiesResult, g.board.Grid[i].Type, i, g.board.Grid[i].State, g.boardSize)
								g.board.Turn = "white"
								g.board.Grid[i].HoverStone = g.board.StoneW
								g.board.Grid[i].HoverState = 3

							}
							if g.board.Turn == "white" && !moved {
								g.board.Grid[i].Stone = g.board.StoneW
								g.board.Grid[i].State = 1
								g.board.LastPlacedIndex = i
								moved = true
								g.board.CanCurrentPlacedStoneCapture(chanNewLibertiesResult, g.board.Grid[i].Type, i, g.board.Grid[i].State, g.boardSize)
								g.board.Turn = "black"
								g.board.Grid[i].HoverStone = g.board.StoneB
								g.board.Grid[i].HoverState = 4

							}
							go g.board.CalculateIndividualGridLiberties(g.boardSize, chanNewLibertiesResult)

						}

					}
				}
			}

		}
		if inpututil.IsKeyJustPressed(ebiten.KeyP) {
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
		if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) {
			timeNow := time.Now()
			if timeNow.Sub(g.board.LastMove) > time.Millisecond*200 {
				g.board.LastMove = timeNow
				for i := range g.board.Grid {
					if g.board.Grid[i].IsStoneHover {
						g.board.CaptureStone(i)
					}
				}
			}
		}
	}

	if g.TouchInput.ShowTouchInput == true {
		g.TouchInput.ProcessTouchInput(screenWidth, screenHeight)
	}

	if g.GameState == UIMainMenu {
		//handle main menu ui
	}
	if g.GameState == UIPauseMenu {
		//handle pause menu ui
		timeNow := time.Now()
		if timeNow.Sub(g.board.LastMove) > time.Millisecond*10 {
			g.board.LastMove = timeNow
			if inpututil.IsKeyJustPressed(ebiten.KeyW) || inpututil.IsKeyJustPressed(ebiten.KeyUp) && len(g.UIs) > 1 {
				for i, element := range g.UIs[1].Elements {
					if element.CurrentState == ui.Selected && element.Name == "quit button" {
						g.UIs[1].Elements[i].CurrentState = ui.Deselected
						g.UIs[1].Elements[i].CurrentImage = g.UIs[1].Elements[i].DeselectedImage
						g.UIs[1].Elements[i-1].CurrentState = ui.Selected
						g.UIs[1].Elements[i-1].CurrentImage = g.UIs[1].Elements[i-1].SelectedImage
					}
				}
			}
			if inpututil.IsKeyJustPressed(ebiten.KeyS) || inpututil.IsKeyJustPressed(ebiten.KeyDown) && len(g.UIs) > 1 {
				for i, element := range g.UIs[1].Elements {
					if element.CurrentState == ui.Selected && element.Name == "resume button" {
						g.UIs[1].Elements[i].CurrentState = ui.Deselected
						g.UIs[1].Elements[i].CurrentImage = g.UIs[1].Elements[i].DeselectedImage
						g.UIs[1].Elements[i+1].CurrentState = ui.Selected
						g.UIs[1].Elements[i+1].CurrentImage = g.UIs[1].Elements[i+1].SelectedImage
					}
				}
			}
			if inpututil.IsKeyJustPressed(ebiten.KeySpace) || inpututil.IsKeyJustPressed(ebiten.KeyEnter) && len(g.UIs) > 1 {
				for _, element := range g.UIs[1].Elements {
					if element.CurrentState == ui.Selected && element.Name == "resume button" {
						g.GameState = Gameplay
					}
					if element.CurrentState == ui.Selected && element.Name == "quit button" {
						return ebiten.Termination
					}
				}
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
	if g.TouchInput.ShowTouchInput == true && len(g.TouchInput.OnScreenButtons) >= 1 {
		g.TouchInput.Draw(screen, screenHeight, screenWidth)
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
