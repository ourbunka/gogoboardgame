package game

import (
	"gogoboardgame/board"
	"gogoboardgame/ui"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

func (g *Game) MoveUp() {
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

func (g *Game) MoveDown() {
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

func (g *Game) MoveLeft() {
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

func (g *Game) MoveRight() {
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

func (g *Game) PlaceStone() {
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

func (g *Game) RemoveStone() {
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

func (g *Game) TogglePauseMenu() {
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
}

func (g *Game) UIConfirm() error {
	for _, element := range g.UIs[1].Elements {
		if element.CurrentState == ui.Selected && element.Name == "resume button" {
			g.GameState = Gameplay
		}
		if element.CurrentState == ui.Selected && element.Name == "quit button" {
			return ebiten.Termination
		}
	}
	return nil
}

func (g *Game) UIUp() {
	for i, element := range g.UIs[1].Elements {
		if element.CurrentState == ui.Selected && element.Name == "quit button" {
			g.UIs[1].Elements[i].CurrentState = ui.Deselected
			g.UIs[1].Elements[i].CurrentImage = g.UIs[1].Elements[i].DeselectedImage
			g.UIs[1].Elements[i-1].CurrentState = ui.Selected
			g.UIs[1].Elements[i-1].CurrentImage = g.UIs[1].Elements[i-1].SelectedImage
		}
	}
}

func (g *Game) UIDown() {
	for i, element := range g.UIs[1].Elements {
		if element.CurrentState == ui.Selected && element.Name == "resume button" {
			g.UIs[1].Elements[i].CurrentState = ui.Deselected
			g.UIs[1].Elements[i].CurrentImage = g.UIs[1].Elements[i].DeselectedImage
			g.UIs[1].Elements[i+1].CurrentState = ui.Selected
			g.UIs[1].Elements[i+1].CurrentImage = g.UIs[1].Elements[i+1].SelectedImage
		}
	}
}
