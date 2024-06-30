package board

import (
	"bytes"
	"embed"
	"image"
	_ "image/png"
	"log"
	"strings"
	"time"

	"os"

	"github.com/hajimehoshi/ebiten/v2"
)

type GridType int

var Resources embed.FS
var UseEmbeded bool
var Build = "DESKTOP"

const (
	GridCenter GridType = iota + 1
	GridTopLeft
	GridTopEdge
	GridTopRight
	GridLeftEdge
	GridRightEdge
	GridBottomLeft
	GridBottomEdge
	GridBottomRight
)

type HasLibertiesBool struct {
	HasTop    bool
	HasLeft   bool
	HasRight  bool
	HasBottom bool
}

type LibertiesUpdateResult struct {
	index            int
	HasLibertiesBool HasLibertiesBool
	Liberties        int
}

type BoardGrid struct {
	Image        *ebiten.Image
	Stone        *ebiten.Image
	HoverStone   *ebiten.Image
	Type         GridType
	IsStoneHover bool
	Row          int
	Col          int
	PosX         float32
	PosY         float32
	State        int
	HoverState   int
	HasLiberties HasLibertiesBool
	Liberties    int
}

type Board struct {
	Grid            []BoardGrid
	StoneW          *ebiten.Image
	StoneB          *ebiten.Image
	Turn            string
	LastMove        time.Time
	LastPlacedIndex int
}

func (b *Board) NewBoard(screenWidth, screenHeight, boardSize int, rendererScale float64) error {
	var newBoard Board
	boardTopLeft, _, err := LoadImage("./resources/board/board_topleft.png")
	if err != nil {
		return err
	}
	boardTopEdge, _, err := LoadImage("./resources/board/board_topedge.png")
	if err != nil {
		return err
	}
	boardTopRight, _, err := LoadImage("./resources/board/board_topright.png")
	if err != nil {
		return err
	}
	boardLeftEdge, _, err := LoadImage("./resources/board/board_leftedge.png")
	if err != nil {
		return err
	}
	boardRightEdge, _, err := LoadImage("./resources/board/board_rightedge.png")
	if err != nil {
		return err
	}
	boardCenter, _, err := LoadImage("./resources/board/board_center.png")
	if err != nil {
		return err
	}
	boardBottomLeft, _, err := LoadImage("./resources/board/board_bottomleft.png")
	if err != nil {
		return err
	}
	boardBottomRight, _, err := LoadImage("./resources/board/board_bottomright.png")
	if err != nil {
		return err
	}
	boardBottomEdge, _, err := LoadImage("./resources/board/board_bottomedge.png")
	if err != nil {
		return err
	}

	for i := 0; i < boardSize; i++ {
		for j := 0; j < boardSize; j++ {
			var boardImg *ebiten.Image
			//typeOfGrid := "./resources/board/board_center.png"
			boardImg = boardCenter
			newType := GridCenter
			if i == 0 && j == 0 {
				//typeOfGrid = "./resources/board/board_topleft.png"
				boardImg = boardTopLeft
				newType = GridTopLeft
			}
			if i == boardSize-1 && j == 0 {
				//typeOfGrid = "./resources/board/board_topright.png"
				boardImg = boardTopRight
				newType = GridTopRight
			}
			if i == boardSize-1 && j == boardSize-1 {
				//typeOfGrid = "./resources/board/board_bottomright.png"
				boardImg = boardBottomRight
				newType = GridBottomRight
			}
			if i == 0 && j == boardSize-1 {
				//typeOfGrid = "./resources/board/board_bottomleft.png"
				boardImg = boardBottomLeft
				newType = GridBottomLeft
			}
			if i == 0 && j != 0 && j != boardSize-1 {
				//typeOfGrid = "./resources/board/board_leftedge.png"
				boardImg = boardLeftEdge
				newType = GridLeftEdge
			}
			if i == boardSize-1 && j != 0 && j != boardSize-1 {
				//typeOfGrid = "./resources/board/board_rightedge.png"
				boardImg = boardRightEdge
				newType = GridRightEdge
			}
			if j == 0 && i != 0 && i != boardSize-1 {
				//typeOfGrid = "./resources/board/board_topedge.png"
				boardImg = boardTopEdge
				newType = GridTopEdge
			}
			if j == boardSize-1 && i != 0 && i != boardSize-1 {
				//typeOfGrid = "./resources/board/board_bottomedge.png"
				boardImg = boardBottomEdge
				newType = GridBottomEdge
			}

			var scaleFactor int = int(1 / rendererScale)
			//println("scale factor :", scaleFactor)
			newGrid, err := SpawnGrid(float32((i*1024)+((screenWidth-(boardSize*1024/scaleFactor))/2*scaleFactor)), float32((j*1024)+((screenHeight-(boardSize*1024/scaleFactor))/2*scaleFactor)), boardImg, newType, i, j)
			if err != nil {
				return err
			}
			newBoard.Grid = append(newBoard.Grid, newGrid)
		}

	}
	newStoneW, _, err := LoadImage("./resources/stone/stone_w_l.png")
	if err != nil {
		return err
	}
	newStoneB, _, err := LoadImage("./resources/stone/stone_b_l.png")
	if err != nil {
		return err
	}
	b.StoneW = newStoneW
	b.StoneB = newStoneB
	b.Turn = "black"
	b.LastMove = time.Now()
	b.LastPlacedIndex = 0
	b.Grid = append(b.Grid, newBoard.Grid...)
	return nil
}

func SpawnGrid(posX float32, posY float32, image *ebiten.Image, gridType GridType, row int, col int) (BoardGrid, error) {
	var newGrid BoardGrid
	newGrid.Row = row
	newGrid.Col = col
	newGrid.PosX = posX
	newGrid.PosY = posY
	newGrid.State = 0
	newGrid.HoverState = 0

	newGrid.Stone = nil
	newGrid.Image = image
	newGrid.Type = gridType
	newGrid.Liberties = 4

	return newGrid, nil
}

func LoadImage(path string) (*ebiten.Image, *image.Alpha, error) {
	var imageByte []byte
	var err error
	if path == "" {
		if UseEmbeded {
			log.Println("====== ", Build, "==========")
			imageByte, err = Resources.ReadFile("resources/board/board_center.png")
		} else {
			log.Println("====== ", Build, "==========")
			imageByte, err = os.ReadFile("./resources/board/board_center.png")
		}
	} else {
		if UseEmbeded {
			log.Println("====== ", Build, "==========")
			path = strings.TrimPrefix(path, "./")
			imageByte, err = Resources.ReadFile(path)
		} else {
			log.Println("====== ", Build, "==========")
			imageByte, err = os.ReadFile(path)
		}

	}
	if err != nil {
		println("Failed to load test image, REASON :" + err.Error())
		return nil, nil, err
	}
	img, fileFormat, err := image.Decode(bytes.NewReader(imageByte))
	if err != nil {
		println(path)
		println("Failed to decode image : " + err.Error() + "  file format : " + fileFormat)
		return nil, nil, err
	}
	//println(fileFormat)
	ebitenImage := ebiten.NewImageFromImage(img)
	bound := img.Bounds()
	ebitenAlphaImage := image.NewAlpha(bound)

	for j := bound.Min.Y; j < bound.Max.Y; j++ {
		for i := bound.Min.X; i < bound.Max.X; i++ {
			ebitenAlphaImage.Set(i, j, img.At(i, j))
		}
	}
	return ebitenImage, ebitenAlphaImage, nil
}

func (b *Board) CheckLiberties(state int, index int, boardSize int) bool {
	gridType := b.Grid[index].Type
	libertiesCounter := 0
	TopIndex := index - 1
	BottomIndex := index + 1
	LeftIndex := index - boardSize
	RightIndex := index + boardSize
	if gridType == GridCenter {
		if b.Grid[TopIndex].State == state || b.Grid[TopIndex].State == 0 {
			libertiesCounter++
		}
		if b.Grid[BottomIndex].State == state || b.Grid[BottomIndex].State == 0 {
			libertiesCounter++
		}
		if b.Grid[LeftIndex].State == state || b.Grid[LeftIndex].State == 0 {
			libertiesCounter++
		}
		if b.Grid[RightIndex].State == state || b.Grid[RightIndex].State == 0 {
			libertiesCounter++
		}
	} else if gridType == GridTopLeft {
		if b.Grid[BottomIndex].State == state || b.Grid[BottomIndex].State == 0 {
			libertiesCounter++
		}
		if b.Grid[RightIndex].State == state || b.Grid[RightIndex].State == 0 {
			libertiesCounter++
		}
	} else if gridType == GridTopEdge {
		if b.Grid[BottomIndex].State == state || b.Grid[BottomIndex].State == 0 {
			libertiesCounter++
		}
		if b.Grid[LeftIndex].State == state || b.Grid[LeftIndex].State == 0 {
			libertiesCounter++
		}
		if b.Grid[RightIndex].State == state || b.Grid[RightIndex].State == 0 {
			libertiesCounter++
		}
	} else if gridType == GridTopRight {
		if b.Grid[BottomIndex].State == state || b.Grid[BottomIndex].State == 0 {
			libertiesCounter++
		}
		if b.Grid[LeftIndex].State == state || b.Grid[LeftIndex].State == 0 {
			libertiesCounter++
		}
	} else if gridType == GridLeftEdge {
		if b.Grid[TopIndex].State == state || b.Grid[TopIndex].State == 0 {
			libertiesCounter++
		}
		if b.Grid[BottomIndex].State == state || b.Grid[BottomIndex].State == 0 {
			libertiesCounter++
		}
		if b.Grid[RightIndex].State == state || b.Grid[RightIndex].State == 0 {
			libertiesCounter++
		}
	} else if gridType == GridRightEdge {
		if b.Grid[TopIndex].State == state || b.Grid[TopIndex].State == 0 {
			libertiesCounter++
		}
		if b.Grid[BottomIndex].State == state || b.Grid[BottomIndex].State == 0 {
			libertiesCounter++
		}
		if b.Grid[LeftIndex].State == state || b.Grid[LeftIndex].State == 0 {
			libertiesCounter++
		}
	} else if gridType == GridBottomLeft {
		if b.Grid[TopIndex].State == state || b.Grid[TopIndex].State == 0 {
			libertiesCounter++
		}
		if b.Grid[RightIndex].State == state || b.Grid[RightIndex].State == 0 {
			libertiesCounter++
		}
	} else if gridType == GridBottomEdge {
		if b.Grid[TopIndex].State == state || b.Grid[TopIndex].State == 0 {
			libertiesCounter++
		}
		if b.Grid[LeftIndex].State == state || b.Grid[LeftIndex].State == 0 {
			libertiesCounter++
		}
		if b.Grid[RightIndex].State == state || b.Grid[RightIndex].State == 0 {
			libertiesCounter++
		}
	} else if gridType == GridBottomRight {
		if b.Grid[TopIndex].State == state || b.Grid[TopIndex].State == 0 {
			libertiesCounter++
		}
		if b.Grid[LeftIndex].State == state || b.Grid[LeftIndex].State == 0 {
			libertiesCounter++
		}

	}
	println("LIBERTIES COUNTER ", libertiesCounter)
	if libertiesCounter != 0 {
		return true
	} else {
		if gridType == GridCenter {
			if b.LastPlacedIndex == TopIndex || b.LastPlacedIndex == LeftIndex ||
				b.LastPlacedIndex == RightIndex || b.LastPlacedIndex == BottomIndex {
				return false
			}
			canCapture := b.canCaptureWithZeroLiberties(TopIndex, state, boardSize)
			if canCapture {
				println(b.LastPlacedIndex, TopIndex)
				b.CaptureStone(TopIndex)

				return true
			}

			canCapture = b.canCaptureWithZeroLiberties(LeftIndex, state, boardSize)
			if canCapture {
				println(b.LastPlacedIndex, LeftIndex)
				b.CaptureStone(LeftIndex)

				return true
			}

			canCapture = b.canCaptureWithZeroLiberties(RightIndex, state, boardSize)
			if canCapture {
				println(b.LastPlacedIndex, RightIndex)
				b.CaptureStone(RightIndex)
				return true
			}

			canCapture = b.canCaptureWithZeroLiberties(BottomIndex, state, boardSize)
			if canCapture {
				println(b.LastPlacedIndex, BottomIndex)
				b.CaptureStone(BottomIndex)
				return true
			}
		}
		return false
	}
}

func (b *Board) CanCurrentPlacedStoneCapture(chanNewLibertiesResult chan []LibertiesUpdateResult, gridType GridType, index int, currentState int, boardSize int) bool {
	var canCapture bool = false
	var state = currentState
	var hasCaptured bool = false
	if gridType == GridCenter {
		opposingColorTopIndex := index - 1
		if b.LastPlacedIndex == opposingColorTopIndex {
			RepetitionDetected()
			return false
		}
		opposingColorLeftIndex := index - boardSize
		if b.LastPlacedIndex == opposingColorLeftIndex {
			RepetitionDetected()
			return false
		}
		opposingColorRightIndex := index + boardSize
		if b.LastPlacedIndex == opposingColorRightIndex {
			RepetitionDetected()
			return false
		}
		opposingBottomIndex := index + 1
		if b.LastPlacedIndex == opposingBottomIndex {
			RepetitionDetected()
			return false
		}
		canCapture = b.CanCapture(opposingColorTopIndex, state, boardSize)
		if canCapture {
			b.CaptureStone(opposingColorTopIndex)
			hasCaptured = true
			println(b.LastPlacedIndex, opposingColorTopIndex)
		}
		canCapture = b.CanCapture(opposingColorLeftIndex, state, boardSize)
		if canCapture {
			b.CaptureStone(opposingColorLeftIndex)
			hasCaptured = true
			println(b.LastPlacedIndex, opposingColorLeftIndex)

		}
		canCapture = b.CanCapture(opposingColorRightIndex, state, boardSize)
		if canCapture {
			b.CaptureStone(opposingColorRightIndex)
			hasCaptured = true
			println(b.LastPlacedIndex, opposingColorRightIndex)
		}
		canCapture = b.CanCapture(opposingBottomIndex, state, boardSize)
		if canCapture {
			b.CaptureStone(opposingBottomIndex)
			hasCaptured = true
			println(b.LastPlacedIndex, opposingBottomIndex)
		}
	}
	if gridType == GridTopLeft {
		opposingColorRightIndex := index + boardSize
		canCapture = b.CanCapture(opposingColorRightIndex, state, boardSize)
		if canCapture {
			b.CaptureStone(opposingColorRightIndex)
			hasCaptured = true
		}
		opposingBottomIndex := index + 1
		canCapture = b.CanCapture(opposingBottomIndex, state, boardSize)
		if canCapture {
			b.CaptureStone(opposingBottomIndex)
			hasCaptured = true
		}
	}

	if gridType == GridTopEdge {
		opposingColorLeftIndex := index - boardSize
		canCapture = b.CanCapture(opposingColorLeftIndex, state, boardSize)
		if canCapture {
			b.CaptureStone(opposingColorLeftIndex)
			hasCaptured = true
		}
		opposingColorRightIndex := index + boardSize
		canCapture = b.CanCapture(opposingColorRightIndex, state, boardSize)
		if canCapture {
			b.CaptureStone(opposingColorRightIndex)
			hasCaptured = true
		}
		opposingBottomIndex := index + 1
		canCapture = b.CanCapture(opposingBottomIndex, state, boardSize)
		if canCapture {
			b.CaptureStone(opposingBottomIndex)
			hasCaptured = true
		}
	}
	if gridType == GridTopRight {
		opposingColorLeftIndex := index - boardSize
		canCapture = b.CanCapture(opposingColorLeftIndex, state, boardSize)
		if canCapture {
			b.CaptureStone(opposingColorLeftIndex)
			hasCaptured = true
		}
		opposingBottomIndex := index + 1
		canCapture = b.CanCapture(opposingBottomIndex, state, boardSize)
		if canCapture {
			b.CaptureStone(opposingBottomIndex)
			hasCaptured = true
		}
	}
	if gridType == GridLeftEdge {
		opposingColorTopIndex := index - 1
		canCapture = b.CanCapture(opposingColorTopIndex, state, boardSize)
		if canCapture {
			b.CaptureStone(opposingColorTopIndex)
			hasCaptured = true
		}
		opposingColorRightIndex := index + boardSize
		canCapture = b.CanCapture(opposingColorRightIndex, state, boardSize)
		if canCapture {
			b.CaptureStone(opposingColorRightIndex)
			hasCaptured = true
		}
		opposingBottomIndex := index + 1
		canCapture = b.CanCapture(opposingBottomIndex, state, boardSize)
		if canCapture {
			b.CaptureStone(opposingBottomIndex)
			hasCaptured = true
		}
	}
	if gridType == GridRightEdge {
		opposingColorTopIndex := index - 1
		canCapture = b.CanCapture(opposingColorTopIndex, state, boardSize)
		if canCapture {
			b.CaptureStone(opposingColorTopIndex)
			hasCaptured = true
		}
		opposingColorLeftIndex := index - boardSize
		canCapture = b.CanCapture(opposingColorLeftIndex, state, boardSize)
		if canCapture {
			b.CaptureStone(opposingColorLeftIndex)
			hasCaptured = true
		}
		opposingBottomIndex := index + 1
		canCapture = b.CanCapture(opposingBottomIndex, state, boardSize)
		if canCapture {
			b.CaptureStone(opposingBottomIndex)
			hasCaptured = true
		}
	}
	if gridType == GridBottomLeft {
		opposingColorTopIndex := index - 1
		canCapture = b.CanCapture(opposingColorTopIndex, state, boardSize)
		if canCapture {
			b.CaptureStone(opposingColorTopIndex)
			hasCaptured = true
		}
		opposingColorRightIndex := index + boardSize
		canCapture = b.CanCapture(opposingColorRightIndex, state, boardSize)
		if canCapture {
			b.CaptureStone(opposingColorRightIndex)
			hasCaptured = true
		}
	}
	if gridType == GridBottomEdge {
		opposingColorTopIndex := index - 1
		canCapture = b.CanCapture(opposingColorTopIndex, state, boardSize)
		if canCapture {
			b.CaptureStone(opposingColorTopIndex)
			hasCaptured = true
		}
		opposingColorLeftIndex := index - boardSize
		canCapture = b.CanCapture(opposingColorLeftIndex, state, boardSize)
		if canCapture {
			b.CaptureStone(opposingColorLeftIndex)
			hasCaptured = true
		}
		opposingColorRightIndex := index + boardSize
		canCapture = b.CanCapture(opposingColorRightIndex, state, boardSize)
		if canCapture {
			b.CaptureStone(opposingColorRightIndex)
			hasCaptured = true
		}
	}
	if gridType == GridBottomRight {
		opposingColorTopIndex := index - 1
		canCapture = b.CanCapture(opposingColorTopIndex, state, boardSize)
		if canCapture {
			b.CaptureStone(opposingColorTopIndex)
			hasCaptured = true
		}
		opposingColorLeftIndex := index - boardSize
		canCapture = b.CanCapture(opposingColorLeftIndex, state, boardSize)
		if canCapture {
			b.CaptureStone(opposingColorLeftIndex)
			hasCaptured = true
		}
	}
	println("has captured :", hasCaptured)
	return hasCaptured
}

func (b *Board) CanCapture(canCaptureIndex int, CurrentTurn int, boardSize int) bool {
	gridType := b.Grid[canCaptureIndex].Type
	maxLiberties := 0
	surroundedCounter := 0
	TopIndex := canCaptureIndex - 1
	BottomIndex := canCaptureIndex + 1
	LeftIndex := canCaptureIndex - boardSize
	RightIndex := canCaptureIndex + boardSize
	if gridType == GridCenter {
		maxLiberties = 4
		if b.Grid[TopIndex].State != b.Grid[canCaptureIndex].State && b.Grid[TopIndex].State != 0 {
			surroundedCounter++
		}

		if b.Grid[BottomIndex].State != b.Grid[canCaptureIndex].State && b.Grid[BottomIndex].State != 0 {
			surroundedCounter++
		}

		if b.Grid[LeftIndex].State != b.Grid[canCaptureIndex].State && b.Grid[LeftIndex].State != 0 {
			surroundedCounter++
		}

		if b.Grid[RightIndex].State != b.Grid[canCaptureIndex].State && b.Grid[RightIndex].State != 0 {
			surroundedCounter++
		}

	}
	if gridType == GridTopLeft {
		maxLiberties = 2
		if b.Grid[BottomIndex].State != b.Grid[canCaptureIndex].State && b.Grid[BottomIndex].State != 0 {
			surroundedCounter++
		}
		if b.Grid[RightIndex].State != b.Grid[canCaptureIndex].State && b.Grid[RightIndex].State != 0 {
			surroundedCounter++
		}
	}
	if gridType == GridTopEdge {
		maxLiberties = 3
		if b.Grid[BottomIndex].State != b.Grid[canCaptureIndex].State && b.Grid[BottomIndex].State != 0 {
			surroundedCounter++
		}
		if b.Grid[LeftIndex].State != b.Grid[canCaptureIndex].State && b.Grid[LeftIndex].State != 0 {
			surroundedCounter++
		}
		if b.Grid[RightIndex].State != b.Grid[canCaptureIndex].State && b.Grid[RightIndex].State != 0 {
			surroundedCounter++
		}
	}
	if gridType == GridTopRight {
		maxLiberties = 2
		if b.Grid[BottomIndex].State != b.Grid[canCaptureIndex].State && b.Grid[BottomIndex].State != 0 {
			surroundedCounter++
		}
		if b.Grid[LeftIndex].State != b.Grid[canCaptureIndex].State && b.Grid[LeftIndex].State != 0 {
			surroundedCounter++
		}
	}
	if gridType == GridLeftEdge {
		maxLiberties = 3
		if b.Grid[TopIndex].State != b.Grid[canCaptureIndex].State && b.Grid[TopIndex].State != 0 {
			surroundedCounter++
		}
		if b.Grid[BottomIndex].State != b.Grid[canCaptureIndex].State && b.Grid[BottomIndex].State != 0 {
			surroundedCounter++
		}
		if b.Grid[RightIndex].State != b.Grid[canCaptureIndex].State && b.Grid[RightIndex].State != 0 {
			surroundedCounter++
		}
	}
	if gridType == GridRightEdge {
		maxLiberties = 3
		if b.Grid[TopIndex].State != b.Grid[canCaptureIndex].State && b.Grid[TopIndex].State != 0 {
			surroundedCounter++
		}
		if b.Grid[BottomIndex].State != b.Grid[canCaptureIndex].State && b.Grid[BottomIndex].State != 0 {
			surroundedCounter++
		}
		if b.Grid[LeftIndex].State != b.Grid[canCaptureIndex].State && b.Grid[LeftIndex].State != 0 {
			surroundedCounter++
		}
	}
	if gridType == GridBottomLeft {
		maxLiberties = 2
		if b.Grid[TopIndex].State != b.Grid[canCaptureIndex].State && b.Grid[TopIndex].State != 0 {
			surroundedCounter++
		}
		if b.Grid[RightIndex].State != b.Grid[canCaptureIndex].State && b.Grid[RightIndex].State != 0 {
			surroundedCounter++
		}
	}
	if gridType == GridBottomEdge {
		maxLiberties = 3
		if b.Grid[TopIndex].State != b.Grid[canCaptureIndex].State && b.Grid[TopIndex].State != 0 {
			surroundedCounter++
		}
		if b.Grid[LeftIndex].State != b.Grid[canCaptureIndex].State && b.Grid[LeftIndex].State != 0 {
			surroundedCounter++
		}
		if b.Grid[RightIndex].State != b.Grid[canCaptureIndex].State && b.Grid[RightIndex].State != 0 {
			surroundedCounter++
		}
	}
	if gridType == GridBottomRight {
		maxLiberties = 2
		if b.Grid[TopIndex].State != b.Grid[canCaptureIndex].State && b.Grid[TopIndex].State != 0 {
			surroundedCounter++
		}
		if b.Grid[LeftIndex].State != b.Grid[canCaptureIndex].State && b.Grid[LeftIndex].State != 0 {
			surroundedCounter++
		}
	}
	println("max liberties ", maxLiberties)
	println(surroundedCounter)
	if surroundedCounter == maxLiberties {
		return true
	} else {
		return false
	}
}
func (b *Board) canCaptureWithZeroLiberties(canCaptureIndex int, CurrentTurn int, boardSize int) bool {
	gridType := b.Grid[canCaptureIndex].Type
	maxLiberties := 0
	surroundedCounter := 0
	TopIndex := canCaptureIndex - 1
	BottomIndex := canCaptureIndex + 1
	LeftIndex := canCaptureIndex - boardSize
	RightIndex := canCaptureIndex + boardSize
	if gridType == GridCenter {
		maxLiberties = 4
		if b.Grid[TopIndex].State != b.Grid[canCaptureIndex].State && b.Grid[TopIndex].State != 0 {
			surroundedCounter++
		}

		if b.Grid[BottomIndex].State != b.Grid[canCaptureIndex].State && b.Grid[BottomIndex].State != 0 {
			surroundedCounter++
		}

		if b.Grid[LeftIndex].State != b.Grid[canCaptureIndex].State && b.Grid[LeftIndex].State != 0 {
			surroundedCounter++
		}

		if b.Grid[RightIndex].State != b.Grid[canCaptureIndex].State && b.Grid[RightIndex].State != 0 {
			surroundedCounter++
		}
		if surroundedCounter > maxLiberties-1 {
			if b.Grid[TopIndex].State == 0 || b.Grid[LeftIndex].State == 0 ||
				b.Grid[RightIndex].State == 0 || b.Grid[BottomIndex].State == 0 {
				maxLiberties--
			}
		}

	}
	if gridType == GridTopLeft {
		maxLiberties = 2
		if b.Grid[BottomIndex].State != b.Grid[canCaptureIndex].State && b.Grid[BottomIndex].State != 0 {
			surroundedCounter++
		}
		if b.Grid[RightIndex].State != b.Grid[canCaptureIndex].State && b.Grid[RightIndex].State != 0 {
			surroundedCounter++
		}
		if surroundedCounter > maxLiberties-1 {
			if b.Grid[RightIndex].State == 0 || b.Grid[BottomIndex].State == 0 {
				maxLiberties--
			}
		}
	}
	if gridType == GridTopEdge {
		maxLiberties = 3
		if b.Grid[BottomIndex].State != b.Grid[canCaptureIndex].State && b.Grid[BottomIndex].State != 0 {
			surroundedCounter++
		}
		if b.Grid[LeftIndex].State != b.Grid[canCaptureIndex].State && b.Grid[LeftIndex].State != 0 {
			surroundedCounter++
		}
		if b.Grid[RightIndex].State != b.Grid[canCaptureIndex].State && b.Grid[RightIndex].State != 0 {
			surroundedCounter++
		}
		if surroundedCounter > maxLiberties-1 {
			if b.Grid[LeftIndex].State == 0 ||
				b.Grid[RightIndex].State == 0 || b.Grid[BottomIndex].State == 0 {
				maxLiberties--
			}
		}
	}
	if gridType == GridTopRight {
		maxLiberties = 2
		if b.Grid[BottomIndex].State != b.Grid[canCaptureIndex].State && b.Grid[BottomIndex].State != 0 {
			surroundedCounter++
		}
		if b.Grid[LeftIndex].State != b.Grid[canCaptureIndex].State && b.Grid[LeftIndex].State != 0 {
			surroundedCounter++
		}
		if surroundedCounter > maxLiberties-1 {
			if b.Grid[LeftIndex].State == 0 || b.Grid[BottomIndex].State == 0 {
				maxLiberties--
			}
		}
	}
	if gridType == GridLeftEdge {
		maxLiberties = 3
		if b.Grid[TopIndex].State != b.Grid[canCaptureIndex].State && b.Grid[TopIndex].State != 0 {
			surroundedCounter++
		}
		if b.Grid[BottomIndex].State != b.Grid[canCaptureIndex].State && b.Grid[BottomIndex].State != 0 {
			surroundedCounter++
		}
		if b.Grid[RightIndex].State != b.Grid[canCaptureIndex].State && b.Grid[RightIndex].State != 0 {
			surroundedCounter++
		}
		if surroundedCounter > maxLiberties-1 {
			if b.Grid[TopIndex].State == 0 || b.Grid[RightIndex].State == 0 ||
				b.Grid[BottomIndex].State == 0 {
				maxLiberties--
			}
		}
	}
	if gridType == GridRightEdge {
		maxLiberties = 3
		if b.Grid[TopIndex].State != b.Grid[canCaptureIndex].State && b.Grid[TopIndex].State != 0 {
			surroundedCounter++
		}
		if b.Grid[BottomIndex].State != b.Grid[canCaptureIndex].State && b.Grid[BottomIndex].State != 0 {
			surroundedCounter++
		}
		if b.Grid[LeftIndex].State != b.Grid[canCaptureIndex].State && b.Grid[LeftIndex].State != 0 {
			surroundedCounter++
		}
		if surroundedCounter > maxLiberties-1 {
			if b.Grid[TopIndex].State == 0 || b.Grid[LeftIndex].State == 0 ||
				b.Grid[BottomIndex].State == 0 {
				maxLiberties--
			}
		}
	}
	if gridType == GridBottomLeft {
		maxLiberties = 2
		if b.Grid[TopIndex].State != b.Grid[canCaptureIndex].State && b.Grid[TopIndex].State != 0 {
			surroundedCounter++
		}
		if b.Grid[RightIndex].State != b.Grid[canCaptureIndex].State && b.Grid[RightIndex].State != 0 {
			surroundedCounter++
		}
		if surroundedCounter > maxLiberties-1 {
			if b.Grid[TopIndex].State == 0 || b.Grid[RightIndex].State == 0 {
				maxLiberties--
			}
		}
	}
	if gridType == GridBottomEdge {
		maxLiberties = 3
		if b.Grid[TopIndex].State != b.Grid[canCaptureIndex].State && b.Grid[TopIndex].State != 0 {
			surroundedCounter++
		}
		if b.Grid[LeftIndex].State != b.Grid[canCaptureIndex].State && b.Grid[LeftIndex].State != 0 {
			surroundedCounter++
		}
		if b.Grid[RightIndex].State != b.Grid[canCaptureIndex].State && b.Grid[RightIndex].State != 0 {
			surroundedCounter++
		}
		if surroundedCounter > maxLiberties-1 {
			if b.Grid[TopIndex].State == 0 || b.Grid[LeftIndex].State == 0 ||
				b.Grid[RightIndex].State == 0 {
				maxLiberties--
			}
		}
	}
	if gridType == GridBottomRight {
		maxLiberties = 2
		if b.Grid[TopIndex].State != b.Grid[canCaptureIndex].State && b.Grid[TopIndex].State != 0 {
			surroundedCounter++
		}
		if b.Grid[LeftIndex].State != b.Grid[canCaptureIndex].State && b.Grid[LeftIndex].State != 0 {
			surroundedCounter++
		}
		if surroundedCounter > maxLiberties-1 {
			if b.Grid[TopIndex].State == 0 || b.Grid[LeftIndex].State == 0 {
				maxLiberties--
			}
		}

	}
	println("max liberties ", maxLiberties)
	println(surroundedCounter)
	if surroundedCounter >= maxLiberties-1 {
		return true
	} else {
		return false
	}
}

func (b *Board) CaptureStone(captureIndex int) {
	//b.lastPlacedIndex = captureIndex
	println("new capture index : ", b.LastPlacedIndex, captureIndex)
	b.Grid[captureIndex].State = 0
	b.Grid[captureIndex].Stone = nil
}

func ToDo() {

}

func RepetitionDetected() {
	ToDo() //ui
	println("REPETITON DETECTED. CANNOT REPEAT MOVE.")
}

// test
func (b *Board) CalculateIndividualGridLiberties(boardSize int, chanNewLibertiesResult chan []LibertiesUpdateResult) {
	start := time.Now()
	boardClones := b
	lenght := len(boardClones.Grid)
	var result []LibertiesUpdateResult
	for i := 0; i < lenght; i++ {
		var placeholder LibertiesUpdateResult
		result = append(result, placeholder)
	}
	println(len(result), lenght)

	for i := range boardClones.Grid {
		//var maxLiberties int
		result[i].index = i
		TopIndex := i - 1
		LeftIndex := i - boardSize
		RightIndex := i + boardSize
		BottomIndex := i + 1
		if boardClones.Grid[i].State != 0 {
			if boardClones.Grid[i].Type == GridCenter {
				//maxLiberties = 4
				result[i].Liberties = 0
				//check top
				result[i].HasLibertiesBool.HasTop = boardClones.CheckIndividualLiberties(TopIndex, i)
				if result[i].HasLibertiesBool.HasTop {
					result[i].Liberties++
				}
				result[i].HasLibertiesBool.HasLeft = boardClones.CheckIndividualLiberties(LeftIndex, i)
				if result[i].HasLibertiesBool.HasLeft {
					result[i].Liberties++
				}
				result[i].HasLibertiesBool.HasRight = boardClones.CheckIndividualLiberties(RightIndex, i)
				if result[i].HasLibertiesBool.HasRight {
					result[i].Liberties++
				}
				result[i].HasLibertiesBool.HasBottom = boardClones.CheckIndividualLiberties(BottomIndex, i)
				if result[i].HasLibertiesBool.HasBottom {
					result[i].Liberties++
				}
			} else if boardClones.Grid[i].Type == GridTopLeft || boardClones.Grid[i].Type == GridTopRight ||
				boardClones.Grid[i].Type == GridBottomLeft || boardClones.Grid[i].Type == GridBottomRight {
				//maxLiberties = 2
				if boardClones.Grid[i].Type == GridTopLeft {
					result[i].HasLibertiesBool.HasRight = boardClones.CheckIndividualLiberties(RightIndex, i)
					if result[i].HasLibertiesBool.HasRight {
						result[i].Liberties++
					}
					result[i].HasLibertiesBool.HasBottom = boardClones.CheckIndividualLiberties(BottomIndex, i)
					if result[i].HasLibertiesBool.HasBottom {
						result[i].Liberties++
					}
				}
				if boardClones.Grid[i].Type == GridTopRight {
					result[i].HasLibertiesBool.HasLeft = boardClones.CheckIndividualLiberties(LeftIndex, i)
					if result[i].HasLibertiesBool.HasLeft {
						result[i].Liberties++
					}
					result[i].HasLibertiesBool.HasBottom = boardClones.CheckIndividualLiberties(BottomIndex, i)
					if result[i].HasLibertiesBool.HasBottom {
						result[i].Liberties++
					}
				}
				if boardClones.Grid[i].Type == GridBottomLeft {
					result[i].HasLibertiesBool.HasTop = boardClones.CheckIndividualLiberties(TopIndex, i)
					if result[i].HasLibertiesBool.HasTop {
						result[i].Liberties++
					}
					result[i].HasLibertiesBool.HasRight = boardClones.CheckIndividualLiberties(RightIndex, i)
					if result[i].HasLibertiesBool.HasRight {
						result[i].Liberties++
					}
				}
				if boardClones.Grid[i].Type == GridBottomRight {
					result[i].HasLibertiesBool.HasLeft = boardClones.CheckIndividualLiberties(LeftIndex, i)
					if result[i].HasLibertiesBool.HasLeft {
						result[i].Liberties++
					}
					result[i].HasLibertiesBool.HasTop = boardClones.CheckIndividualLiberties(TopIndex, i)
					if result[i].HasLibertiesBool.HasTop {
						result[i].Liberties++
					}
				}
			} else {
				//maxLiberties = 3
				if boardClones.Grid[i].Type == GridTopEdge {
					result[i].HasLibertiesBool.HasLeft = boardClones.CheckIndividualLiberties(LeftIndex, i)
					if result[i].HasLibertiesBool.HasLeft {
						result[i].Liberties++
					}
					result[i].HasLibertiesBool.HasRight = boardClones.CheckIndividualLiberties(RightIndex, i)
					if result[i].HasLibertiesBool.HasRight {
						result[i].Liberties++
					}
					result[i].HasLibertiesBool.HasBottom = boardClones.CheckIndividualLiberties(BottomIndex, i)
					if result[i].HasLibertiesBool.HasBottom {
						result[i].Liberties++
					}
				}
				if boardClones.Grid[i].Type == GridLeftEdge {
					result[i].HasLibertiesBool.HasTop = boardClones.CheckIndividualLiberties(TopIndex, i)
					if result[i].HasLibertiesBool.HasTop {
						result[i].Liberties++
					}
					result[i].HasLibertiesBool.HasRight = boardClones.CheckIndividualLiberties(RightIndex, i)
					if result[i].HasLibertiesBool.HasRight {
						result[i].Liberties++
					}
					result[i].HasLibertiesBool.HasBottom = boardClones.CheckIndividualLiberties(BottomIndex, i)
					if result[i].HasLibertiesBool.HasBottom {
						result[i].Liberties++
					}
				}
				if boardClones.Grid[i].Type == GridRightEdge {
					result[i].HasLibertiesBool.HasTop = boardClones.CheckIndividualLiberties(TopIndex, i)
					if result[i].HasLibertiesBool.HasTop {
						result[i].Liberties++
					}
					result[i].HasLibertiesBool.HasLeft = boardClones.CheckIndividualLiberties(LeftIndex, i)
					if result[i].HasLibertiesBool.HasLeft {
						result[i].Liberties++
					}
					result[i].HasLibertiesBool.HasBottom = boardClones.CheckIndividualLiberties(BottomIndex, i)
					if result[i].HasLibertiesBool.HasBottom {
						result[i].Liberties++
					}
				}
				if boardClones.Grid[i].Type == GridBottomEdge {
					result[i].HasLibertiesBool.HasTop = boardClones.CheckIndividualLiberties(TopIndex, i)
					if result[i].HasLibertiesBool.HasTop {
						result[i].Liberties++
					}
					result[i].HasLibertiesBool.HasLeft = boardClones.CheckIndividualLiberties(LeftIndex, i)
					if result[i].HasLibertiesBool.HasLeft {
						result[i].Liberties++
					}
					result[i].HasLibertiesBool.HasRight = boardClones.CheckIndividualLiberties(RightIndex, i)
					if result[i].HasLibertiesBool.HasRight {
						result[i].Liberties++
					}
				}
			}
		}

	}
	//
	//                                                                  :...:-.       .  *.-.:.*..
	//                                                               ..=:%#-*.-%*-:.*#*=--.-..=%*+.
	//                                                              ..*+....:....:%:..  .    ....=%#:
	//                                                             .:* ..           .-=..         -*.
	//                                .-=====-..                   .+..         .          .      .@.
	//                              .#++++++++++%=.               .%.           .**.    :#=#:.    .+..
	//                            .+#++++++++++++++%:.            *:.     ....     .    .       ..-==
	//                            .*++++++++++**++%++##..      .*-..      .+*=#..     :-.=:. :-==-:.+..
	//                            .#++++++++++%*@#++++++%..  .*...          ....%= -:.=**=--%...    ==.
	//                            .#++++++++++@@++++++++++%=.%..                ..:*#*-...*..       .#.
	//                            .#++++++++++++++++++++++++*..        .:-:.    =*..     .%:. +.. .  +.
	//                            ..%++++++++++++++++++++++%.         .-:.@@@-.  .*.      :%:.#.:*@@@-+
	//                             .-*++++++++++++++++++++#..          .-%@@@+=..           :%..-++#*.#:
	//                              .:#+++++++++++++++++++*                   .:=*.          .#-     .#.
	//                                .%+++++++++++++++++*+      .            #....          +*.     :+.
	//                                 ##@#++++++++++++++#+                  .=%..  -%#.    ++..     *:
	//                                *%:...#*+++++++++++#+                    .....     ..#%.       #.
	//                                 .    ..-%+++++++++%-                      .  ...             :%..
	//                                          .=%+++++**.                  *-%%=... ...+%:.      =*..
	//                                            ...+@*%..                 .#.            .-...:#+..
	//                                                .:%.                                   .#*..
	//                                               ..%-.                             .+*++:@+**.
	//                                               :@-.                      ...+*.-.*+@@%+*+++#+..
	//                                               %.                       ....=+@@@@%++#+++++++#%..
	//                                            .-@#:.                          --@%*+++++++++++++++@=.
	//                                          ..%=====*%=.                        +-..-@*++++++++++++++#
	//                                       . :@#====#%*==*+.              .%.     .#=   .:#%++++++++++++
	//                                      .*#::%====%%#====#:.           .-%.       =*      ..*%++++++++
	//                               ..-%@#:.....%============+*..                    ..@@@%*+..  ..%*++++
	//                              =@-          %==+%%#=%%%#====*#..                      .. .:#*.  .-%*+
	//really dumb way of finding line of stones
	//to redo

	for i := range result {
		if result[i].Liberties == 1 {
			if result[i].HasLibertiesBool.HasTop {
				if result[i-1].Liberties == 1 && result[i-1].HasLibertiesBool.HasBottom {
					result[i].Liberties = 0
					result[i].HasLibertiesBool.HasTop = false
					result[i-1].Liberties = 0
					result[i-1].HasLibertiesBool.HasBottom = false
				}
			}
			if result[i].HasLibertiesBool.HasLeft {
				if result[i-boardSize].Liberties == 1 && result[i-boardSize].HasLibertiesBool.HasRight {
					result[i].Liberties = 0
					result[i].HasLibertiesBool.HasLeft = false
					result[i-boardSize].Liberties = 0
					result[i-boardSize].HasLibertiesBool.HasRight = false
				}
			}
			if result[i].HasLibertiesBool.HasRight {
				if result[i+boardSize].Liberties == 1 && result[i+boardSize].HasLibertiesBool.HasLeft {
					result[i].Liberties = 0
					result[i].HasLibertiesBool.HasRight = false
					result[i+boardSize].Liberties = 0
					result[i+boardSize].HasLibertiesBool.HasLeft = false
				}
			}
			if result[i].HasLibertiesBool.HasBottom {
				if result[i+1].Liberties == 1 && result[i+1].HasLibertiesBool.HasTop {
					result[i].Liberties = 0
					result[i].HasLibertiesBool.HasBottom = false
					result[i+1].Liberties = 0
					result[i+1].HasLibertiesBool.HasTop = false
				}
			}
		}
		if result[i].Liberties == 2 {
			if result[i].HasLibertiesBool.HasTop && result[i].HasLibertiesBool.HasBottom {
				if i-1 >= 0 && i+1 < len(result)-1 && result[i-1].Liberties == 1 && result[i+1].Liberties == 1 &&
					result[i-1].HasLibertiesBool.HasBottom && result[i+1].HasLibertiesBool.HasTop {
					result[i].Liberties = 0
					result[i].HasLibertiesBool.HasTop = false
					result[i].HasLibertiesBool.HasBottom = false
					result[i-1].Liberties = 0
					result[i-1].HasLibertiesBool.HasBottom = false
					result[i+1].Liberties = 0
					result[i+1].HasLibertiesBool.HasTop = false
				}
				if i-2 >= 0 && i+2 < len(result)-1 && result[i-1].Liberties == 2 && result[i+1].Liberties == 2 &&
					result[i-1].HasLibertiesBool.HasTop && result[i-1].HasLibertiesBool.HasBottom &&
					result[i-1-1].Liberties == 1 && result[i-1-1].HasLibertiesBool.HasBottom &&
					result[i+1].HasLibertiesBool.HasTop && result[i+1].HasLibertiesBool.HasBottom &&
					result[i+1+1].Liberties == 1 && result[i+1+1].HasLibertiesBool.HasTop {
					result[i].Liberties = 0
					result[i].HasLibertiesBool.HasTop = false
					result[i].HasLibertiesBool.HasBottom = false
					result[i-1].Liberties = 0
					result[i-1].HasLibertiesBool.HasTop = false
					result[i-1].HasLibertiesBool.HasBottom = false
					result[i+1].Liberties = 0
					result[i+1].HasLibertiesBool.HasTop = false
					result[i+1].HasLibertiesBool.HasBottom = false
					result[i-1-1].Liberties = 0
					result[i-1-1].HasLibertiesBool.HasBottom = false
					result[i+1+1].Liberties = 0
					result[i+1+1].HasLibertiesBool.HasTop = false
				}
				if i-3 >= 0 && i+3 < len(result)-1 && result[i-1].Liberties == 2 && result[i+1].Liberties == 2 &&
					result[i-1].HasLibertiesBool.HasTop && result[i-1].HasLibertiesBool.HasBottom &&
					result[i-2].Liberties == 2 && result[i-2].HasLibertiesBool.HasTop && result[i-2].HasLibertiesBool.HasBottom &&
					result[i-3].Liberties == 1 && result[i-3].HasLibertiesBool.HasBottom &&
					result[i+1].HasLibertiesBool.HasTop && result[i+1].HasLibertiesBool.HasBottom &&
					result[i+2].Liberties == 2 && result[i+2].HasLibertiesBool.HasTop && result[i+2].HasLibertiesBool.HasBottom &&
					result[i+3].Liberties == 1 && result[i+3].HasLibertiesBool.HasTop {
					result[i].Liberties = 0
					result[i].HasLibertiesBool.HasTop = false
					result[i].HasLibertiesBool.HasBottom = false
					result[i-1].Liberties = 0
					result[i-1].HasLibertiesBool.HasTop = false
					result[i-1].HasLibertiesBool.HasBottom = false
					result[i+1].Liberties = 0
					result[i+1].HasLibertiesBool.HasTop = false
					result[i+1].HasLibertiesBool.HasBottom = false
					result[i-2].Liberties = 0
					result[i-2].HasLibertiesBool.HasBottom = false
					result[i-2].HasLibertiesBool.HasTop = false
					result[i+2].Liberties = 0
					result[i+2].HasLibertiesBool.HasTop = false
					result[i+2].HasLibertiesBool.HasBottom = false
					result[i+3].Liberties = 0
					result[i+3].HasLibertiesBool.HasTop = false
					result[i-3].Liberties = 0
					result[i-3].HasLibertiesBool.HasBottom = false
				}
				if i-4 >= 0 && i+4 < len(result)-1 && result[i-1].Liberties == 2 && result[i+1].Liberties == 2 &&
					result[i-1].HasLibertiesBool.HasTop && result[i-1].HasLibertiesBool.HasBottom &&
					result[i-2].Liberties == 2 && result[i-2].HasLibertiesBool.HasTop && result[i-2].HasLibertiesBool.HasBottom &&
					result[i-3].Liberties == 2 && result[i-3].HasLibertiesBool.HasTop && result[i-3].HasLibertiesBool.HasBottom &&
					result[i-4].Liberties == 1 && result[i-4].HasLibertiesBool.HasBottom &&
					result[i+1].HasLibertiesBool.HasTop && result[i+1].HasLibertiesBool.HasBottom &&
					result[i+2].Liberties == 2 && result[i+2].HasLibertiesBool.HasTop && result[i+2].HasLibertiesBool.HasBottom &&
					result[i+3].Liberties == 2 && result[i+3].HasLibertiesBool.HasTop && result[i+3].HasLibertiesBool.HasBottom &&
					result[i+4].Liberties == 1 && result[i+4].HasLibertiesBool.HasTop {
					result[i].Liberties = 0
					result[i].HasLibertiesBool.HasTop = false
					result[i].HasLibertiesBool.HasBottom = false
					result[i-1].Liberties = 0
					result[i-1].HasLibertiesBool.HasTop = false
					result[i-1].HasLibertiesBool.HasBottom = false
					result[i+1].Liberties = 0
					result[i+1].HasLibertiesBool.HasTop = false
					result[i+1].HasLibertiesBool.HasBottom = false
					result[i-2].Liberties = 0
					result[i-2].HasLibertiesBool.HasBottom = false
					result[i-2].HasLibertiesBool.HasTop = false
					result[i+2].Liberties = 0
					result[i+2].HasLibertiesBool.HasTop = false
					result[i+2].HasLibertiesBool.HasBottom = false
					result[i+3].Liberties = 0
					result[i+3].HasLibertiesBool.HasTop = false
					result[i+3].HasLibertiesBool.HasBottom = false
					result[i-3].Liberties = 0
					result[i-3].HasLibertiesBool.HasTop = false
					result[i-3].HasLibertiesBool.HasBottom = false
					result[i+4].Liberties = 0
					result[i+4].HasLibertiesBool.HasTop = false
					result[i-4].Liberties = 0
					result[i-4].HasLibertiesBool.HasBottom = false
				}
				if i-5 >= 0 && i+5 < len(result)-1 && result[i-1].Liberties == 2 && result[i+1].Liberties == 2 &&
					result[i-1].HasLibertiesBool.HasTop && result[i-1].HasLibertiesBool.HasBottom &&
					result[i-2].Liberties == 2 && result[i-2].HasLibertiesBool.HasTop && result[i-2].HasLibertiesBool.HasBottom &&
					result[i-3].Liberties == 2 && result[i-3].HasLibertiesBool.HasTop && result[i-3].HasLibertiesBool.HasBottom &&
					result[i-4].Liberties == 2 && result[i-4].HasLibertiesBool.HasTop && result[i-4].HasLibertiesBool.HasBottom &&
					result[i-5].Liberties == 1 && result[i-5].HasLibertiesBool.HasBottom &&
					result[i+1].HasLibertiesBool.HasTop && result[i+1].HasLibertiesBool.HasBottom &&
					result[i+2].Liberties == 2 && result[i+2].HasLibertiesBool.HasTop && result[i+2].HasLibertiesBool.HasBottom &&
					result[i+3].Liberties == 2 && result[i+3].HasLibertiesBool.HasTop && result[i+3].HasLibertiesBool.HasBottom &&
					result[i+4].Liberties == 2 && result[i+4].HasLibertiesBool.HasTop && result[i+4].HasLibertiesBool.HasBottom &&
					result[i+5].Liberties == 1 && result[i+5].HasLibertiesBool.HasTop {
					result[i].Liberties = 0
					result[i].HasLibertiesBool.HasTop = false
					result[i].HasLibertiesBool.HasBottom = false
					result[i-1].Liberties = 0
					result[i-1].HasLibertiesBool.HasTop = false
					result[i-1].HasLibertiesBool.HasBottom = false
					result[i+1].Liberties = 0
					result[i+1].HasLibertiesBool.HasTop = false
					result[i+1].HasLibertiesBool.HasBottom = false
					result[i-2].Liberties = 0
					result[i-2].HasLibertiesBool.HasBottom = false
					result[i-2].HasLibertiesBool.HasTop = false
					result[i+2].Liberties = 0
					result[i+2].HasLibertiesBool.HasTop = false
					result[i+2].HasLibertiesBool.HasBottom = false
					result[i+3].Liberties = 0
					result[i+3].HasLibertiesBool.HasTop = false
					result[i+3].HasLibertiesBool.HasBottom = false
					result[i-3].Liberties = 0
					result[i-3].HasLibertiesBool.HasTop = false
					result[i-3].HasLibertiesBool.HasBottom = false
					result[i+4].Liberties = 0
					result[i+4].HasLibertiesBool.HasTop = false
					result[i+4].HasLibertiesBool.HasBottom = false
					result[i-4].Liberties = 0
					result[i-4].HasLibertiesBool.HasTop = false
					result[i-4].HasLibertiesBool.HasBottom = false
					result[i+5].Liberties = 0
					result[i+5].HasLibertiesBool.HasTop = false
					result[i-5].Liberties = 0
					result[i-5].HasLibertiesBool.HasBottom = false
				}
				if i-6 >= 0 && i+6 < len(result)-1 && result[i-1].Liberties == 2 && result[i+1].Liberties == 2 &&
					result[i-1].HasLibertiesBool.HasTop && result[i-1].HasLibertiesBool.HasBottom &&
					result[i-2].Liberties == 2 && result[i-2].HasLibertiesBool.HasTop && result[i-2].HasLibertiesBool.HasBottom &&
					result[i-3].Liberties == 2 && result[i-3].HasLibertiesBool.HasTop && result[i-3].HasLibertiesBool.HasBottom &&
					result[i-4].Liberties == 2 && result[i-4].HasLibertiesBool.HasTop && result[i-4].HasLibertiesBool.HasBottom &&
					result[i-5].Liberties == 2 && result[i-5].HasLibertiesBool.HasTop && result[i-5].HasLibertiesBool.HasBottom &&
					result[i-6].Liberties == 1 && result[i-6].HasLibertiesBool.HasBottom &&
					result[i+1].HasLibertiesBool.HasTop && result[i+1].HasLibertiesBool.HasBottom &&
					result[i+2].Liberties == 2 && result[i+2].HasLibertiesBool.HasTop && result[i+2].HasLibertiesBool.HasBottom &&
					result[i+3].Liberties == 2 && result[i+3].HasLibertiesBool.HasTop && result[i+3].HasLibertiesBool.HasBottom &&
					result[i+4].Liberties == 2 && result[i+4].HasLibertiesBool.HasTop && result[i+4].HasLibertiesBool.HasBottom &&
					result[i+5].Liberties == 2 && result[i+5].HasLibertiesBool.HasTop && result[i+5].HasLibertiesBool.HasBottom &&
					result[i+6].Liberties == 1 && result[i+6].HasLibertiesBool.HasTop {
					result[i].Liberties = 0
					result[i].HasLibertiesBool.HasTop = false
					result[i].HasLibertiesBool.HasBottom = false
					result[i-1].Liberties = 0
					result[i-1].HasLibertiesBool.HasTop = false
					result[i-1].HasLibertiesBool.HasBottom = false
					result[i+1].Liberties = 0
					result[i+1].HasLibertiesBool.HasTop = false
					result[i+1].HasLibertiesBool.HasBottom = false
					result[i-2].Liberties = 0
					result[i-2].HasLibertiesBool.HasBottom = false
					result[i-2].HasLibertiesBool.HasTop = false
					result[i+2].Liberties = 0
					result[i+2].HasLibertiesBool.HasTop = false
					result[i+2].HasLibertiesBool.HasBottom = false
					result[i+3].Liberties = 0
					result[i+3].HasLibertiesBool.HasTop = false
					result[i+3].HasLibertiesBool.HasBottom = false
					result[i-3].Liberties = 0
					result[i-3].HasLibertiesBool.HasTop = false
					result[i-3].HasLibertiesBool.HasBottom = false
					result[i+4].Liberties = 0
					result[i+4].HasLibertiesBool.HasTop = false
					result[i+4].HasLibertiesBool.HasBottom = false
					result[i-4].Liberties = 0
					result[i-4].HasLibertiesBool.HasTop = false
					result[i-4].HasLibertiesBool.HasBottom = false
					result[i+5].Liberties = 0
					result[i+5].HasLibertiesBool.HasTop = false
					result[i+5].HasLibertiesBool.HasBottom = false
					result[i-5].Liberties = 0
					result[i-5].HasLibertiesBool.HasTop = false
					result[i-5].HasLibertiesBool.HasBottom = false
					result[i+6].Liberties = 0
					result[i+6].HasLibertiesBool.HasTop = false
					result[i-6].Liberties = 0
					result[i-6].HasLibertiesBool.HasBottom = false
				}
				if i-7 >= 0 && i+7 < len(result)-1 && result[i-1].Liberties == 2 && result[i+1].Liberties == 2 &&
					result[i-1].HasLibertiesBool.HasTop && result[i-1].HasLibertiesBool.HasBottom &&
					result[i-2].Liberties == 2 && result[i-2].HasLibertiesBool.HasTop && result[i-2].HasLibertiesBool.HasBottom &&
					result[i-3].Liberties == 2 && result[i-3].HasLibertiesBool.HasTop && result[i-3].HasLibertiesBool.HasBottom &&
					result[i-4].Liberties == 2 && result[i-4].HasLibertiesBool.HasTop && result[i-4].HasLibertiesBool.HasBottom &&
					result[i-5].Liberties == 2 && result[i-5].HasLibertiesBool.HasTop && result[i-5].HasLibertiesBool.HasBottom &&
					result[i-6].Liberties == 2 && result[i-6].HasLibertiesBool.HasTop && result[i-6].HasLibertiesBool.HasBottom &&
					result[i-7].Liberties == 1 && result[i-7].HasLibertiesBool.HasBottom &&
					result[i+1].HasLibertiesBool.HasTop && result[i+1].HasLibertiesBool.HasBottom &&
					result[i+2].Liberties == 2 && result[i+2].HasLibertiesBool.HasTop && result[i+2].HasLibertiesBool.HasBottom &&
					result[i+3].Liberties == 2 && result[i+3].HasLibertiesBool.HasTop && result[i+3].HasLibertiesBool.HasBottom &&
					result[i+4].Liberties == 2 && result[i+4].HasLibertiesBool.HasTop && result[i+4].HasLibertiesBool.HasBottom &&
					result[i+5].Liberties == 2 && result[i+5].HasLibertiesBool.HasTop && result[i+5].HasLibertiesBool.HasBottom &&
					result[i+6].Liberties == 2 && result[i+6].HasLibertiesBool.HasTop && result[i+6].HasLibertiesBool.HasBottom &&
					result[i+7].Liberties == 1 && result[i+7].HasLibertiesBool.HasTop {
					result[i].Liberties = 0
					result[i].HasLibertiesBool.HasTop = false
					result[i].HasLibertiesBool.HasBottom = false
					result[i-1].Liberties = 0
					result[i-1].HasLibertiesBool.HasTop = false
					result[i-1].HasLibertiesBool.HasBottom = false
					result[i+1].Liberties = 0
					result[i+1].HasLibertiesBool.HasTop = false
					result[i+1].HasLibertiesBool.HasBottom = false
					result[i-2].Liberties = 0
					result[i-2].HasLibertiesBool.HasBottom = false
					result[i-2].HasLibertiesBool.HasTop = false
					result[i+2].Liberties = 0
					result[i+2].HasLibertiesBool.HasTop = false
					result[i+2].HasLibertiesBool.HasBottom = false
					result[i+3].Liberties = 0
					result[i+3].HasLibertiesBool.HasTop = false
					result[i+3].HasLibertiesBool.HasBottom = false
					result[i-3].Liberties = 0
					result[i-3].HasLibertiesBool.HasTop = false
					result[i-3].HasLibertiesBool.HasBottom = false
					result[i+4].Liberties = 0
					result[i+4].HasLibertiesBool.HasTop = false
					result[i+4].HasLibertiesBool.HasBottom = false
					result[i-4].Liberties = 0
					result[i-4].HasLibertiesBool.HasTop = false
					result[i-4].HasLibertiesBool.HasBottom = false
					result[i+5].Liberties = 0
					result[i+5].HasLibertiesBool.HasTop = false
					result[i+5].HasLibertiesBool.HasBottom = false
					result[i-5].Liberties = 0
					result[i-5].HasLibertiesBool.HasTop = false
					result[i-5].HasLibertiesBool.HasBottom = false
					result[i+6].Liberties = 0
					result[i+6].HasLibertiesBool.HasTop = false
					result[i+6].HasLibertiesBool.HasBottom = false
					result[i-6].Liberties = 0
					result[i-6].HasLibertiesBool.HasTop = false
					result[i-6].HasLibertiesBool.HasBottom = false
					result[i+7].Liberties = 0
					result[i+7].HasLibertiesBool.HasTop = false
					result[i-7].Liberties = 0
					result[i-7].HasLibertiesBool.HasBottom = false
				}
				if i-8 >= 0 && i+8 < len(result)-1 && result[i-1].Liberties == 2 && result[i+1].Liberties == 2 &&
					result[i-1].HasLibertiesBool.HasTop && result[i-1].HasLibertiesBool.HasBottom &&
					result[i-2].Liberties == 2 && result[i-2].HasLibertiesBool.HasTop && result[i-2].HasLibertiesBool.HasBottom &&
					result[i-3].Liberties == 2 && result[i-3].HasLibertiesBool.HasTop && result[i-3].HasLibertiesBool.HasBottom &&
					result[i-4].Liberties == 2 && result[i-4].HasLibertiesBool.HasTop && result[i-4].HasLibertiesBool.HasBottom &&
					result[i-5].Liberties == 2 && result[i-5].HasLibertiesBool.HasTop && result[i-5].HasLibertiesBool.HasBottom &&
					result[i-6].Liberties == 2 && result[i-6].HasLibertiesBool.HasTop && result[i-6].HasLibertiesBool.HasBottom &&
					result[i-7].Liberties == 2 && result[i-7].HasLibertiesBool.HasTop && result[i-7].HasLibertiesBool.HasBottom &&
					result[i-8].Liberties == 1 && result[i-8].HasLibertiesBool.HasBottom &&
					result[i+1].HasLibertiesBool.HasTop && result[i+1].HasLibertiesBool.HasBottom &&
					result[i+2].Liberties == 2 && result[i+2].HasLibertiesBool.HasTop && result[i+2].HasLibertiesBool.HasBottom &&
					result[i+3].Liberties == 2 && result[i+3].HasLibertiesBool.HasTop && result[i+3].HasLibertiesBool.HasBottom &&
					result[i+4].Liberties == 2 && result[i+4].HasLibertiesBool.HasTop && result[i+4].HasLibertiesBool.HasBottom &&
					result[i+5].Liberties == 2 && result[i+5].HasLibertiesBool.HasTop && result[i+5].HasLibertiesBool.HasBottom &&
					result[i+6].Liberties == 2 && result[i+6].HasLibertiesBool.HasTop && result[i+6].HasLibertiesBool.HasBottom &&
					result[i+7].Liberties == 2 && result[i+7].HasLibertiesBool.HasTop && result[i+7].HasLibertiesBool.HasBottom &&
					result[i+8].Liberties == 1 && result[i+8].HasLibertiesBool.HasTop {
					result[i].Liberties = 0
					result[i].HasLibertiesBool.HasTop = false
					result[i].HasLibertiesBool.HasBottom = false
					result[i-1].Liberties = 0
					result[i-1].HasLibertiesBool.HasTop = false
					result[i-1].HasLibertiesBool.HasBottom = false
					result[i+1].Liberties = 0
					result[i+1].HasLibertiesBool.HasTop = false
					result[i+1].HasLibertiesBool.HasBottom = false
					result[i-2].Liberties = 0
					result[i-2].HasLibertiesBool.HasBottom = false
					result[i-2].HasLibertiesBool.HasTop = false
					result[i+2].Liberties = 0
					result[i+2].HasLibertiesBool.HasTop = false
					result[i+2].HasLibertiesBool.HasBottom = false
					result[i+3].Liberties = 0
					result[i+3].HasLibertiesBool.HasTop = false
					result[i+3].HasLibertiesBool.HasBottom = false
					result[i-3].Liberties = 0
					result[i-3].HasLibertiesBool.HasTop = false
					result[i-3].HasLibertiesBool.HasBottom = false
					result[i+4].Liberties = 0
					result[i+4].HasLibertiesBool.HasTop = false
					result[i+4].HasLibertiesBool.HasBottom = false
					result[i-4].Liberties = 0
					result[i-4].HasLibertiesBool.HasTop = false
					result[i-4].HasLibertiesBool.HasBottom = false
					result[i+5].Liberties = 0
					result[i+5].HasLibertiesBool.HasTop = false
					result[i+5].HasLibertiesBool.HasBottom = false
					result[i-5].Liberties = 0
					result[i-5].HasLibertiesBool.HasTop = false
					result[i-5].HasLibertiesBool.HasBottom = false
					result[i+6].Liberties = 0
					result[i+6].HasLibertiesBool.HasTop = false
					result[i+6].HasLibertiesBool.HasBottom = false
					result[i-6].Liberties = 0
					result[i-6].HasLibertiesBool.HasTop = false
					result[i-6].HasLibertiesBool.HasBottom = false
					result[i+7].Liberties = 0
					result[i+7].HasLibertiesBool.HasTop = false
					result[i+7].HasLibertiesBool.HasBottom = false
					result[i-7].Liberties = 0
					result[i-7].HasLibertiesBool.HasTop = false
					result[i-7].HasLibertiesBool.HasBottom = false
					result[i+8].Liberties = 0
					result[i+8].HasLibertiesBool.HasTop = false
					result[i-8].Liberties = 0
					result[i-8].HasLibertiesBool.HasBottom = false
				}
				if i-9 >= 0 && i+9 < len(result)-1 && result[i-1].Liberties == 2 && result[i+1].Liberties == 2 &&
					result[i-1].HasLibertiesBool.HasTop && result[i-1].HasLibertiesBool.HasBottom &&
					result[i-2].Liberties == 2 && result[i-2].HasLibertiesBool.HasTop && result[i-2].HasLibertiesBool.HasBottom &&
					result[i-3].Liberties == 2 && result[i-3].HasLibertiesBool.HasTop && result[i-3].HasLibertiesBool.HasBottom &&
					result[i-4].Liberties == 2 && result[i-4].HasLibertiesBool.HasTop && result[i-4].HasLibertiesBool.HasBottom &&
					result[i-5].Liberties == 2 && result[i-5].HasLibertiesBool.HasTop && result[i-5].HasLibertiesBool.HasBottom &&
					result[i-6].Liberties == 2 && result[i-6].HasLibertiesBool.HasTop && result[i-6].HasLibertiesBool.HasBottom &&
					result[i-7].Liberties == 2 && result[i-7].HasLibertiesBool.HasTop && result[i-7].HasLibertiesBool.HasBottom &&
					result[i-8].Liberties == 2 && result[i-8].HasLibertiesBool.HasTop && result[i-8].HasLibertiesBool.HasBottom &&
					result[i-9].Liberties == 1 && result[i-9].HasLibertiesBool.HasBottom &&
					result[i+1].HasLibertiesBool.HasTop && result[i+1].HasLibertiesBool.HasBottom &&
					result[i+2].Liberties == 2 && result[i+2].HasLibertiesBool.HasTop && result[i+2].HasLibertiesBool.HasBottom &&
					result[i+3].Liberties == 2 && result[i+3].HasLibertiesBool.HasTop && result[i+3].HasLibertiesBool.HasBottom &&
					result[i+4].Liberties == 2 && result[i+4].HasLibertiesBool.HasTop && result[i+4].HasLibertiesBool.HasBottom &&
					result[i+5].Liberties == 2 && result[i+5].HasLibertiesBool.HasTop && result[i+5].HasLibertiesBool.HasBottom &&
					result[i+6].Liberties == 2 && result[i+6].HasLibertiesBool.HasTop && result[i+6].HasLibertiesBool.HasBottom &&
					result[i+7].Liberties == 2 && result[i+7].HasLibertiesBool.HasTop && result[i+7].HasLibertiesBool.HasBottom &&
					result[i+8].Liberties == 2 && result[i+8].HasLibertiesBool.HasTop && result[i+8].HasLibertiesBool.HasBottom &&
					result[i+9].Liberties == 1 && result[i+9].HasLibertiesBool.HasTop {
					result[i].Liberties = 0
					result[i].HasLibertiesBool.HasTop = false
					result[i].HasLibertiesBool.HasBottom = false
					result[i-1].Liberties = 0
					result[i-1].HasLibertiesBool.HasTop = false
					result[i-1].HasLibertiesBool.HasBottom = false
					result[i+1].Liberties = 0
					result[i+1].HasLibertiesBool.HasTop = false
					result[i+1].HasLibertiesBool.HasBottom = false
					result[i-2].Liberties = 0
					result[i-2].HasLibertiesBool.HasBottom = false
					result[i-2].HasLibertiesBool.HasTop = false
					result[i+2].Liberties = 0
					result[i+2].HasLibertiesBool.HasTop = false
					result[i+2].HasLibertiesBool.HasBottom = false
					result[i+3].Liberties = 0
					result[i+3].HasLibertiesBool.HasTop = false
					result[i+3].HasLibertiesBool.HasBottom = false
					result[i-3].Liberties = 0
					result[i-3].HasLibertiesBool.HasTop = false
					result[i-3].HasLibertiesBool.HasBottom = false
					result[i+4].Liberties = 0
					result[i+4].HasLibertiesBool.HasTop = false
					result[i+4].HasLibertiesBool.HasBottom = false
					result[i-4].Liberties = 0
					result[i-4].HasLibertiesBool.HasTop = false
					result[i-4].HasLibertiesBool.HasBottom = false
					result[i+5].Liberties = 0
					result[i+5].HasLibertiesBool.HasTop = false
					result[i+5].HasLibertiesBool.HasBottom = false
					result[i-5].Liberties = 0
					result[i-5].HasLibertiesBool.HasTop = false
					result[i-5].HasLibertiesBool.HasBottom = false
					result[i+6].Liberties = 0
					result[i+6].HasLibertiesBool.HasTop = false
					result[i+6].HasLibertiesBool.HasBottom = false
					result[i-6].Liberties = 0
					result[i-6].HasLibertiesBool.HasTop = false
					result[i-6].HasLibertiesBool.HasBottom = false
					result[i+7].Liberties = 0
					result[i+7].HasLibertiesBool.HasTop = false
					result[i+7].HasLibertiesBool.HasBottom = false
					result[i-7].Liberties = 0
					result[i-7].HasLibertiesBool.HasTop = false
					result[i-7].HasLibertiesBool.HasBottom = false
					result[i+8].Liberties = 0
					result[i+8].HasLibertiesBool.HasTop = false
					result[i+8].HasLibertiesBool.HasBottom = false
					result[i-8].Liberties = 0
					result[i-8].HasLibertiesBool.HasTop = false
					result[i-8].HasLibertiesBool.HasBottom = false
					result[i+9].Liberties = 0
					result[i+9].HasLibertiesBool.HasTop = false
					result[i-9].Liberties = 0
					result[i-9].HasLibertiesBool.HasBottom = false
				}

			}
			if result[i].HasLibertiesBool.HasLeft && result[i].HasLibertiesBool.HasRight {
				if i-(boardSize) >= 0 && i+(boardSize) < len(result)-1 &&
					result[i-boardSize].Liberties == 1 && result[i+boardSize].Liberties == 1 &&
					result[i-boardSize].HasLibertiesBool.HasRight && result[i+boardSize].HasLibertiesBool.HasLeft {
					result[i].Liberties = 0
					result[i].HasLibertiesBool.HasRight = false
					result[i+boardSize].Liberties = 0
					result[i+boardSize].HasLibertiesBool.HasLeft = false
					result[i-boardSize].Liberties = 0
					result[i-boardSize].HasLibertiesBool.HasRight = false
				}
				bS := boardSize
				if i-(bS*2) >= 0 && i+(bS*2) < len(result)-1 &&
					result[i-bS].Liberties == 2 && result[i+bS].Liberties == 2 &&
					result[i-bS].HasLibertiesBool.HasLeft && result[i-bS].HasLibertiesBool.HasRight &&
					result[i-(bS*2)].Liberties == 1 && result[i-(bS*2)].HasLibertiesBool.HasRight &&
					result[i+bS].HasLibertiesBool.HasLeft && result[i+bS].HasLibertiesBool.HasRight &&
					result[i+(bS*2)].Liberties == 1 && result[i+(bS*2)].HasLibertiesBool.HasLeft {
					result[i].Liberties = 0
					result[i].HasLibertiesBool.HasLeft = false
					result[i].HasLibertiesBool.HasRight = false
					result[i-bS].Liberties = 0
					result[i-bS].HasLibertiesBool.HasLeft = false
					result[i-bS].HasLibertiesBool.HasRight = false
					result[i+bS].Liberties = 0
					result[i+bS].HasLibertiesBool.HasLeft = false
					result[i+bS].HasLibertiesBool.HasRight = false
					result[i+(bS*2)].Liberties = 0
					result[i+(bS*2)].HasLibertiesBool.HasLeft = false
					result[i-(bS*2)].Liberties = 0
					result[i-(bS*2)].HasLibertiesBool.HasRight = false
				}
				if i-(bS*3) >= 0 && i+(bS*3) < len(result)-1 &&
					result[i-bS].Liberties == 2 && result[i+bS].Liberties == 2 &&
					result[i-bS].HasLibertiesBool.HasLeft && result[i-bS].HasLibertiesBool.HasRight &&
					result[i-(bS*2)].Liberties == 2 && result[i-(bS*2)].HasLibertiesBool.HasLeft && result[i-(bS*2)].HasLibertiesBool.HasRight &&
					result[i-(bS*3)].Liberties == 1 && result[i-(bS*3)].HasLibertiesBool.HasRight &&
					result[i+bS].HasLibertiesBool.HasLeft && result[i+bS].HasLibertiesBool.HasRight &&
					result[i+(bS*2)].Liberties == 2 && result[i+(bS*2)].HasLibertiesBool.HasLeft && result[i+(bS*2)].HasLibertiesBool.HasRight &&
					result[i+(bS*3)].Liberties == 1 && result[i+(bS*3)].HasLibertiesBool.HasLeft {
					result[i].Liberties = 0
					result[i].HasLibertiesBool.HasLeft = false
					result[i].HasLibertiesBool.HasRight = false
					result[i-bS].Liberties = 0
					result[i-bS].HasLibertiesBool.HasLeft = false
					result[i-bS].HasLibertiesBool.HasRight = false
					result[i+bS].Liberties = 0
					result[i+bS].HasLibertiesBool.HasLeft = false
					result[i+bS].HasLibertiesBool.HasRight = false
					result[i-(bS*2)].Liberties = 0
					result[i-(bS*2)].HasLibertiesBool.HasLeft = false
					result[i-(bS*2)].HasLibertiesBool.HasRight = false
					result[i+(bS*2)].Liberties = 0
					result[i+(bS*2)].HasLibertiesBool.HasLeft = false
					result[i+(bS*2)].HasLibertiesBool.HasRight = false
					result[i+(bS*3)].Liberties = 0
					result[i+(bS*3)].HasLibertiesBool.HasLeft = false
					result[i-(bS*3)].Liberties = 0
					result[i-(bS*3)].HasLibertiesBool.HasRight = false
				}
				if i-(bS*4) >= 0 && i+(bS*4) < len(result)-1 &&
					result[i-bS].Liberties == 2 && result[i+bS].Liberties == 2 &&
					result[i-bS].HasLibertiesBool.HasLeft && result[i-bS].HasLibertiesBool.HasRight &&
					result[i-(bS*2)].Liberties == 2 && result[i-(bS*2)].HasLibertiesBool.HasLeft && result[i-(bS*2)].HasLibertiesBool.HasRight &&
					result[i-(bS*3)].Liberties == 2 && result[i-(bS*3)].HasLibertiesBool.HasLeft && result[i-(bS*3)].HasLibertiesBool.HasRight &&
					result[i-(bS*4)].Liberties == 1 && result[i-(bS*4)].HasLibertiesBool.HasRight &&
					result[i+bS].HasLibertiesBool.HasLeft && result[i+bS].HasLibertiesBool.HasRight &&
					result[i+(bS*2)].Liberties == 2 && result[i+(bS*2)].HasLibertiesBool.HasLeft && result[i+(bS*2)].HasLibertiesBool.HasRight &&
					result[i+(bS*3)].Liberties == 2 && result[i+(bS*3)].HasLibertiesBool.HasLeft && result[i+(bS*3)].HasLibertiesBool.HasRight &&
					result[i+(bS*4)].Liberties == 1 && result[i+(bS*4)].HasLibertiesBool.HasLeft {
					result[i].Liberties = 0
					result[i].HasLibertiesBool.HasLeft = false
					result[i].HasLibertiesBool.HasRight = false
					result[i-bS].Liberties = 0
					result[i-bS].HasLibertiesBool.HasLeft = false
					result[i-bS].HasLibertiesBool.HasRight = false
					result[i+bS].Liberties = 0
					result[i+bS].HasLibertiesBool.HasLeft = false
					result[i+bS].HasLibertiesBool.HasRight = false
					result[i-(bS*2)].Liberties = 0
					result[i-(bS*2)].HasLibertiesBool.HasLeft = false
					result[i-(bS*2)].HasLibertiesBool.HasRight = false
					result[i+(bS*2)].Liberties = 0
					result[i+(bS*2)].HasLibertiesBool.HasLeft = false
					result[i+(bS*2)].HasLibertiesBool.HasRight = false
					result[i+(bS*3)].Liberties = 0
					result[i+(bS*3)].HasLibertiesBool.HasLeft = false
					result[i+(bS*3)].HasLibertiesBool.HasRight = false
					result[i-(bS*3)].Liberties = 0
					result[i-(bS*3)].HasLibertiesBool.HasLeft = false
					result[i-(bS*3)].HasLibertiesBool.HasRight = false
					result[i+(bS*4)].Liberties = 0
					result[i+(bS*4)].HasLibertiesBool.HasLeft = false
					result[i-(bS*4)].Liberties = 0
					result[i-(bS*4)].HasLibertiesBool.HasRight = false
				}
				if i-(bS*5) >= 0 && i+(bS*5) < len(result)-1 &&
					result[i-bS].Liberties == 2 && result[i+bS].Liberties == 2 &&
					result[i-bS].HasLibertiesBool.HasLeft && result[i-bS].HasLibertiesBool.HasRight &&
					result[i-(bS*2)].Liberties == 2 && result[i-(bS*2)].HasLibertiesBool.HasLeft && result[i-(bS*2)].HasLibertiesBool.HasRight &&
					result[i-(bS*3)].Liberties == 2 && result[i-(bS*3)].HasLibertiesBool.HasLeft && result[i-(bS*3)].HasLibertiesBool.HasRight &&
					result[i-(bS*4)].Liberties == 2 && result[i-(bS*4)].HasLibertiesBool.HasLeft && result[i-(bS*4)].HasLibertiesBool.HasRight &&
					result[i-(bS*5)].Liberties == 1 && result[i-(bS*5)].HasLibertiesBool.HasRight &&
					result[i+bS].HasLibertiesBool.HasLeft && result[i+bS].HasLibertiesBool.HasRight &&
					result[i+(bS*2)].Liberties == 2 && result[i+(bS*2)].HasLibertiesBool.HasLeft && result[i+(bS*2)].HasLibertiesBool.HasRight &&
					result[i+(bS*3)].Liberties == 2 && result[i+(bS*3)].HasLibertiesBool.HasLeft && result[i+(bS*3)].HasLibertiesBool.HasRight &&
					result[i+(bS*4)].Liberties == 2 && result[i+(bS*4)].HasLibertiesBool.HasLeft && result[i+(bS*4)].HasLibertiesBool.HasRight &&
					result[i+(bS*5)].Liberties == 1 && result[i+(bS*5)].HasLibertiesBool.HasLeft {
					result[i].Liberties = 0
					result[i].HasLibertiesBool.HasLeft = false
					result[i].HasLibertiesBool.HasRight = false
					result[i-bS].Liberties = 0
					result[i-bS].HasLibertiesBool.HasLeft = false
					result[i-bS].HasLibertiesBool.HasRight = false
					result[i+bS].Liberties = 0
					result[i+bS].HasLibertiesBool.HasLeft = false
					result[i+bS].HasLibertiesBool.HasRight = false
					result[i-(bS*2)].Liberties = 0
					result[i-(bS*2)].HasLibertiesBool.HasLeft = false
					result[i-(bS*2)].HasLibertiesBool.HasRight = false
					result[i+(bS*2)].Liberties = 0
					result[i+(bS*2)].HasLibertiesBool.HasLeft = false
					result[i+(bS*2)].HasLibertiesBool.HasRight = false
					result[i+(bS*3)].Liberties = 0
					result[i+(bS*3)].HasLibertiesBool.HasLeft = false
					result[i+(bS*3)].HasLibertiesBool.HasRight = false
					result[i-(bS*3)].Liberties = 0
					result[i-(bS*3)].HasLibertiesBool.HasLeft = false
					result[i-(bS*3)].HasLibertiesBool.HasRight = false
					result[i+(bS*4)].Liberties = 0
					result[i+(bS*4)].HasLibertiesBool.HasLeft = false
					result[i+(bS*4)].HasLibertiesBool.HasRight = false
					result[i-(bS*4)].Liberties = 0
					result[i-(bS*4)].HasLibertiesBool.HasLeft = false
					result[i-(bS*4)].HasLibertiesBool.HasRight = false
					result[i+(bS*5)].Liberties = 0
					result[i+(bS*5)].HasLibertiesBool.HasLeft = false
					result[i-(bS*5)].Liberties = 0
					result[i-(bS*5)].HasLibertiesBool.HasRight = false
				}
				if i-(bS*6) >= 0 && i+(bS*6) < len(result)-1 &&
					result[i-bS].Liberties == 2 && result[i+bS].Liberties == 2 &&
					result[i-bS].HasLibertiesBool.HasLeft && result[i-bS].HasLibertiesBool.HasRight &&
					result[i-(bS*2)].Liberties == 2 && result[i-(bS*2)].HasLibertiesBool.HasLeft && result[i-(bS*2)].HasLibertiesBool.HasRight &&
					result[i-(bS*3)].Liberties == 2 && result[i-(bS*3)].HasLibertiesBool.HasLeft && result[i-(bS*3)].HasLibertiesBool.HasRight &&
					result[i-(bS*4)].Liberties == 2 && result[i-(bS*4)].HasLibertiesBool.HasLeft && result[i-(bS*4)].HasLibertiesBool.HasRight &&
					result[i-(bS*5)].Liberties == 2 && result[i-(bS*5)].HasLibertiesBool.HasLeft && result[i-(bS*5)].HasLibertiesBool.HasRight &&
					result[i-(bS*6)].Liberties == 1 && result[i-(bS*6)].HasLibertiesBool.HasRight &&
					result[i+bS].HasLibertiesBool.HasLeft && result[i+bS].HasLibertiesBool.HasRight &&
					result[i+(bS*2)].Liberties == 2 && result[i+(bS*2)].HasLibertiesBool.HasLeft && result[i+(bS*2)].HasLibertiesBool.HasRight &&
					result[i+(bS*3)].Liberties == 2 && result[i+(bS*3)].HasLibertiesBool.HasLeft && result[i+(bS*3)].HasLibertiesBool.HasRight &&
					result[i+(bS*4)].Liberties == 2 && result[i+(bS*4)].HasLibertiesBool.HasLeft && result[i+(bS*4)].HasLibertiesBool.HasRight &&
					result[i+(bS*5)].Liberties == 2 && result[i+(bS*5)].HasLibertiesBool.HasLeft && result[i+(bS*5)].HasLibertiesBool.HasRight &&
					result[i+(bS*6)].Liberties == 1 && result[i+(bS*6)].HasLibertiesBool.HasLeft {
					result[i].Liberties = 0
					result[i].HasLibertiesBool.HasLeft = false
					result[i].HasLibertiesBool.HasRight = false
					result[i-bS].Liberties = 0
					result[i-bS].HasLibertiesBool.HasLeft = false
					result[i-bS].HasLibertiesBool.HasRight = false
					result[i+bS].Liberties = 0
					result[i+bS].HasLibertiesBool.HasLeft = false
					result[i+bS].HasLibertiesBool.HasRight = false
					result[i-(bS*2)].Liberties = 0
					result[i-(bS*2)].HasLibertiesBool.HasLeft = false
					result[i-(bS*2)].HasLibertiesBool.HasRight = false
					result[i+(bS*2)].Liberties = 0
					result[i+(bS*2)].HasLibertiesBool.HasLeft = false
					result[i+(bS*2)].HasLibertiesBool.HasRight = false
					result[i+(bS*3)].Liberties = 0
					result[i+(bS*3)].HasLibertiesBool.HasLeft = false
					result[i+(bS*3)].HasLibertiesBool.HasRight = false
					result[i-(bS*3)].Liberties = 0
					result[i-(bS*3)].HasLibertiesBool.HasLeft = false
					result[i-(bS*3)].HasLibertiesBool.HasRight = false
					result[i+(bS*4)].Liberties = 0
					result[i+(bS*4)].HasLibertiesBool.HasLeft = false
					result[i+(bS*4)].HasLibertiesBool.HasRight = false
					result[i-(bS*4)].Liberties = 0
					result[i-(bS*4)].HasLibertiesBool.HasLeft = false
					result[i-(bS*4)].HasLibertiesBool.HasRight = false
					result[i+(bS*5)].Liberties = 0
					result[i+(bS*5)].HasLibertiesBool.HasLeft = false
					result[i+(bS*5)].HasLibertiesBool.HasRight = false
					result[i-(bS*5)].Liberties = 0
					result[i-(bS*5)].HasLibertiesBool.HasLeft = false
					result[i-(bS*5)].HasLibertiesBool.HasRight = false
					result[i+(bS*6)].Liberties = 0
					result[i+(bS*6)].HasLibertiesBool.HasLeft = false
					result[i-(bS*6)].Liberties = 0
					result[i-(bS*6)].HasLibertiesBool.HasRight = false
				}
				if i-(bS*7) >= 0 && i+(bS*7) < len(result)-1 &&
					result[i-bS].Liberties == 2 && result[i+bS].Liberties == 2 &&
					result[i-bS].HasLibertiesBool.HasLeft && result[i-bS].HasLibertiesBool.HasRight &&
					result[i-(bS*2)].Liberties == 2 && result[i-(bS*2)].HasLibertiesBool.HasLeft && result[i-(bS*2)].HasLibertiesBool.HasRight &&
					result[i-(bS*3)].Liberties == 2 && result[i-(bS*3)].HasLibertiesBool.HasLeft && result[i-(bS*3)].HasLibertiesBool.HasRight &&
					result[i-(bS*4)].Liberties == 2 && result[i-(bS*4)].HasLibertiesBool.HasLeft && result[i-(bS*4)].HasLibertiesBool.HasRight &&
					result[i-(bS*5)].Liberties == 2 && result[i-(bS*5)].HasLibertiesBool.HasLeft && result[i-(bS*5)].HasLibertiesBool.HasRight &&
					result[i-(bS*6)].Liberties == 2 && result[i-(bS*6)].HasLibertiesBool.HasLeft && result[i-(bS*6)].HasLibertiesBool.HasRight &&
					result[i-(bS*7)].Liberties == 1 && result[i-(bS*7)].HasLibertiesBool.HasRight &&
					result[i+bS].HasLibertiesBool.HasLeft && result[i+bS].HasLibertiesBool.HasRight &&
					result[i+(bS*2)].Liberties == 2 && result[i+(bS*2)].HasLibertiesBool.HasLeft && result[i+(bS*2)].HasLibertiesBool.HasRight &&
					result[i+(bS*3)].Liberties == 2 && result[i+(bS*3)].HasLibertiesBool.HasLeft && result[i+(bS*3)].HasLibertiesBool.HasRight &&
					result[i+(bS*4)].Liberties == 2 && result[i+(bS*4)].HasLibertiesBool.HasLeft && result[i+(bS*4)].HasLibertiesBool.HasRight &&
					result[i+(bS*5)].Liberties == 2 && result[i+(bS*5)].HasLibertiesBool.HasLeft && result[i+(bS*5)].HasLibertiesBool.HasRight &&
					result[i+(bS*6)].Liberties == 2 && result[i+(bS*6)].HasLibertiesBool.HasLeft && result[i+(bS*6)].HasLibertiesBool.HasRight &&
					result[i+(bS*7)].Liberties == 1 && result[i+(bS*7)].HasLibertiesBool.HasLeft {
					result[i].Liberties = 0
					result[i].HasLibertiesBool.HasLeft = false
					result[i].HasLibertiesBool.HasRight = false
					result[i-bS].Liberties = 0
					result[i-bS].HasLibertiesBool.HasLeft = false
					result[i-bS].HasLibertiesBool.HasRight = false
					result[i+bS].Liberties = 0
					result[i+bS].HasLibertiesBool.HasLeft = false
					result[i+bS].HasLibertiesBool.HasRight = false
					result[i-(bS*2)].Liberties = 0
					result[i-(bS*2)].HasLibertiesBool.HasLeft = false
					result[i-(bS*2)].HasLibertiesBool.HasRight = false
					result[i+(bS*2)].Liberties = 0
					result[i+(bS*2)].HasLibertiesBool.HasLeft = false
					result[i+(bS*2)].HasLibertiesBool.HasRight = false
					result[i+(bS*3)].Liberties = 0
					result[i+(bS*3)].HasLibertiesBool.HasLeft = false
					result[i+(bS*3)].HasLibertiesBool.HasRight = false
					result[i-(bS*3)].Liberties = 0
					result[i-(bS*3)].HasLibertiesBool.HasLeft = false
					result[i-(bS*3)].HasLibertiesBool.HasRight = false
					result[i+(bS*4)].Liberties = 0
					result[i+(bS*4)].HasLibertiesBool.HasLeft = false
					result[i+(bS*4)].HasLibertiesBool.HasRight = false
					result[i-(bS*4)].Liberties = 0
					result[i-(bS*4)].HasLibertiesBool.HasLeft = false
					result[i-(bS*4)].HasLibertiesBool.HasRight = false
					result[i+(bS*5)].Liberties = 0
					result[i+(bS*5)].HasLibertiesBool.HasLeft = false
					result[i+(bS*5)].HasLibertiesBool.HasRight = false
					result[i-(bS*5)].Liberties = 0
					result[i-(bS*5)].HasLibertiesBool.HasLeft = false
					result[i-(bS*5)].HasLibertiesBool.HasRight = false
					result[i+(bS*6)].Liberties = 0
					result[i+(bS*6)].HasLibertiesBool.HasLeft = false
					result[i+(bS*6)].HasLibertiesBool.HasRight = false
					result[i-(bS*6)].Liberties = 0
					result[i-(bS*6)].HasLibertiesBool.HasLeft = false
					result[i-(bS*6)].HasLibertiesBool.HasRight = false
					result[i+(bS*7)].Liberties = 0
					result[i+(bS*7)].HasLibertiesBool.HasLeft = false
					result[i-(bS*7)].Liberties = 0
					result[i-(bS*7)].HasLibertiesBool.HasRight = false
				}
				if i-(bS*8) >= 0 && i+(bS*8) < len(result)-1 &&
					result[i-bS].Liberties == 2 && result[i+bS].Liberties == 2 &&
					result[i-bS].HasLibertiesBool.HasLeft && result[i-bS].HasLibertiesBool.HasRight &&
					result[i-(bS*2)].Liberties == 2 && result[i-(bS*2)].HasLibertiesBool.HasLeft && result[i-(bS*2)].HasLibertiesBool.HasRight &&
					result[i-(bS*3)].Liberties == 2 && result[i-(bS*3)].HasLibertiesBool.HasLeft && result[i-(bS*3)].HasLibertiesBool.HasRight &&
					result[i-(bS*4)].Liberties == 2 && result[i-(bS*4)].HasLibertiesBool.HasLeft && result[i-(bS*4)].HasLibertiesBool.HasRight &&
					result[i-(bS*5)].Liberties == 2 && result[i-(bS*5)].HasLibertiesBool.HasLeft && result[i-(bS*5)].HasLibertiesBool.HasRight &&
					result[i-(bS*6)].Liberties == 2 && result[i-(bS*6)].HasLibertiesBool.HasLeft && result[i-(bS*6)].HasLibertiesBool.HasRight &&
					result[i-(bS*7)].Liberties == 2 && result[i-(bS*7)].HasLibertiesBool.HasLeft && result[i-(bS*7)].HasLibertiesBool.HasRight &&
					result[i-(bS*8)].Liberties == 1 && result[i-(bS*8)].HasLibertiesBool.HasRight &&
					result[i+bS].HasLibertiesBool.HasLeft && result[i+bS].HasLibertiesBool.HasRight &&
					result[i+(bS*2)].Liberties == 2 && result[i+(bS*2)].HasLibertiesBool.HasLeft && result[i+(bS*2)].HasLibertiesBool.HasRight &&
					result[i+(bS*3)].Liberties == 2 && result[i+(bS*3)].HasLibertiesBool.HasLeft && result[i+(bS*3)].HasLibertiesBool.HasRight &&
					result[i+(bS*4)].Liberties == 2 && result[i+(bS*4)].HasLibertiesBool.HasLeft && result[i+(bS*4)].HasLibertiesBool.HasRight &&
					result[i+(bS*5)].Liberties == 2 && result[i+(bS*5)].HasLibertiesBool.HasLeft && result[i+(bS*5)].HasLibertiesBool.HasRight &&
					result[i+(bS*6)].Liberties == 2 && result[i+(bS*6)].HasLibertiesBool.HasLeft && result[i+(bS*6)].HasLibertiesBool.HasRight &&
					result[i+(bS*7)].Liberties == 2 && result[i+(bS*7)].HasLibertiesBool.HasLeft && result[i+(bS*7)].HasLibertiesBool.HasRight &&
					result[i+(bS*8)].Liberties == 1 && result[i+(bS*8)].HasLibertiesBool.HasLeft {
					result[i].Liberties = 0
					result[i].HasLibertiesBool.HasLeft = false
					result[i].HasLibertiesBool.HasRight = false
					result[i-bS].Liberties = 0
					result[i-bS].HasLibertiesBool.HasLeft = false
					result[i-bS].HasLibertiesBool.HasRight = false
					result[i+bS].Liberties = 0
					result[i+bS].HasLibertiesBool.HasLeft = false
					result[i+bS].HasLibertiesBool.HasRight = false
					result[i-(bS*2)].Liberties = 0
					result[i-(bS*2)].HasLibertiesBool.HasLeft = false
					result[i-(bS*2)].HasLibertiesBool.HasRight = false
					result[i+(bS*2)].Liberties = 0
					result[i+(bS*2)].HasLibertiesBool.HasLeft = false
					result[i+(bS*2)].HasLibertiesBool.HasRight = false
					result[i+(bS*3)].Liberties = 0
					result[i+(bS*3)].HasLibertiesBool.HasLeft = false
					result[i+(bS*3)].HasLibertiesBool.HasRight = false
					result[i-(bS*3)].Liberties = 0
					result[i-(bS*3)].HasLibertiesBool.HasLeft = false
					result[i-(bS*3)].HasLibertiesBool.HasRight = false
					result[i+(bS*4)].Liberties = 0
					result[i+(bS*4)].HasLibertiesBool.HasLeft = false
					result[i+(bS*4)].HasLibertiesBool.HasRight = false
					result[i-(bS*4)].Liberties = 0
					result[i-(bS*4)].HasLibertiesBool.HasLeft = false
					result[i-(bS*4)].HasLibertiesBool.HasRight = false
					result[i+(bS*5)].Liberties = 0
					result[i+(bS*5)].HasLibertiesBool.HasLeft = false
					result[i+(bS*5)].HasLibertiesBool.HasRight = false
					result[i-(bS*5)].Liberties = 0
					result[i-(bS*5)].HasLibertiesBool.HasLeft = false
					result[i-(bS*5)].HasLibertiesBool.HasRight = false
					result[i+(bS*6)].Liberties = 0
					result[i+(bS*6)].HasLibertiesBool.HasLeft = false
					result[i+(bS*6)].HasLibertiesBool.HasRight = false
					result[i-(bS*6)].Liberties = 0
					result[i-(bS*6)].HasLibertiesBool.HasLeft = false
					result[i-(bS*6)].HasLibertiesBool.HasRight = false
					result[i+(bS*7)].Liberties = 0
					result[i+(bS*7)].HasLibertiesBool.HasLeft = false
					result[i+(bS*7)].HasLibertiesBool.HasRight = false
					result[i-(bS*7)].Liberties = 0
					result[i-(bS*7)].HasLibertiesBool.HasLeft = false
					result[i-(bS*7)].HasLibertiesBool.HasRight = false
					result[i+(bS*8)].Liberties = 0
					result[i+(bS*8)].HasLibertiesBool.HasLeft = false
					result[i-(bS*8)].Liberties = 0
					result[i-(bS*8)].HasLibertiesBool.HasRight = false
				}
				if i-(bS*9) >= 0 && i+(bS*9) < len(result)-1 &&
					result[i-bS].Liberties == 2 && result[i+bS].Liberties == 2 &&
					result[i-bS].HasLibertiesBool.HasLeft && result[i-bS].HasLibertiesBool.HasRight &&
					result[i-(bS*2)].Liberties == 2 && result[i-(bS*2)].HasLibertiesBool.HasLeft && result[i-(bS*2)].HasLibertiesBool.HasRight &&
					result[i-(bS*3)].Liberties == 2 && result[i-(bS*3)].HasLibertiesBool.HasLeft && result[i-(bS*3)].HasLibertiesBool.HasRight &&
					result[i-(bS*4)].Liberties == 2 && result[i-(bS*4)].HasLibertiesBool.HasLeft && result[i-(bS*4)].HasLibertiesBool.HasRight &&
					result[i-(bS*5)].Liberties == 2 && result[i-(bS*5)].HasLibertiesBool.HasLeft && result[i-(bS*5)].HasLibertiesBool.HasRight &&
					result[i-(bS*6)].Liberties == 2 && result[i-(bS*6)].HasLibertiesBool.HasLeft && result[i-(bS*6)].HasLibertiesBool.HasRight &&
					result[i-(bS*7)].Liberties == 2 && result[i-(bS*7)].HasLibertiesBool.HasLeft && result[i-(bS*7)].HasLibertiesBool.HasRight &&
					result[i-(bS*8)].Liberties == 2 && result[i-(bS*8)].HasLibertiesBool.HasLeft && result[i-(bS*8)].HasLibertiesBool.HasRight &&
					result[i-(bS*9)].Liberties == 1 && result[i-(bS*9)].HasLibertiesBool.HasRight &&
					result[i+bS].HasLibertiesBool.HasLeft && result[i+bS].HasLibertiesBool.HasRight &&
					result[i+(bS*2)].Liberties == 2 && result[i+(bS*2)].HasLibertiesBool.HasLeft && result[i+(bS*2)].HasLibertiesBool.HasRight &&
					result[i+(bS*3)].Liberties == 2 && result[i+(bS*3)].HasLibertiesBool.HasLeft && result[i+(bS*3)].HasLibertiesBool.HasRight &&
					result[i+(bS*4)].Liberties == 2 && result[i+(bS*4)].HasLibertiesBool.HasLeft && result[i+(bS*4)].HasLibertiesBool.HasRight &&
					result[i+(bS*5)].Liberties == 2 && result[i+(bS*5)].HasLibertiesBool.HasLeft && result[i+(bS*5)].HasLibertiesBool.HasRight &&
					result[i+(bS*6)].Liberties == 2 && result[i+(bS*6)].HasLibertiesBool.HasLeft && result[i+(bS*6)].HasLibertiesBool.HasRight &&
					result[i+(bS*7)].Liberties == 2 && result[i+(bS*7)].HasLibertiesBool.HasLeft && result[i+(bS*7)].HasLibertiesBool.HasRight &&
					result[i+(bS*8)].Liberties == 2 && result[i+(bS*8)].HasLibertiesBool.HasLeft && result[i+(bS*8)].HasLibertiesBool.HasRight &&
					result[i+(bS*9)].Liberties == 1 && result[i+(bS*9)].HasLibertiesBool.HasLeft {
					result[i].Liberties = 0
					result[i].HasLibertiesBool.HasLeft = false
					result[i].HasLibertiesBool.HasRight = false
					result[i-bS].Liberties = 0
					result[i-bS].HasLibertiesBool.HasLeft = false
					result[i-bS].HasLibertiesBool.HasRight = false
					result[i+bS].Liberties = 0
					result[i+bS].HasLibertiesBool.HasLeft = false
					result[i+bS].HasLibertiesBool.HasRight = false
					result[i-(bS*2)].Liberties = 0
					result[i-(bS*2)].HasLibertiesBool.HasLeft = false
					result[i-(bS*2)].HasLibertiesBool.HasRight = false
					result[i+(bS*2)].Liberties = 0
					result[i+(bS*2)].HasLibertiesBool.HasLeft = false
					result[i+(bS*2)].HasLibertiesBool.HasRight = false
					result[i+(bS*3)].Liberties = 0
					result[i+(bS*3)].HasLibertiesBool.HasLeft = false
					result[i+(bS*3)].HasLibertiesBool.HasRight = false
					result[i-(bS*3)].Liberties = 0
					result[i-(bS*3)].HasLibertiesBool.HasLeft = false
					result[i-(bS*3)].HasLibertiesBool.HasRight = false
					result[i+(bS*4)].Liberties = 0
					result[i+(bS*4)].HasLibertiesBool.HasLeft = false
					result[i+(bS*4)].HasLibertiesBool.HasRight = false
					result[i-(bS*4)].Liberties = 0
					result[i-(bS*4)].HasLibertiesBool.HasLeft = false
					result[i-(bS*4)].HasLibertiesBool.HasRight = false
					result[i+(bS*5)].Liberties = 0
					result[i+(bS*5)].HasLibertiesBool.HasLeft = false
					result[i+(bS*5)].HasLibertiesBool.HasRight = false
					result[i-(bS*5)].Liberties = 0
					result[i-(bS*5)].HasLibertiesBool.HasLeft = false
					result[i-(bS*5)].HasLibertiesBool.HasRight = false
					result[i+(bS*6)].Liberties = 0
					result[i+(bS*6)].HasLibertiesBool.HasLeft = false
					result[i+(bS*6)].HasLibertiesBool.HasRight = false
					result[i-(bS*6)].Liberties = 0
					result[i-(bS*6)].HasLibertiesBool.HasLeft = false
					result[i-(bS*6)].HasLibertiesBool.HasRight = false
					result[i+(bS*7)].Liberties = 0
					result[i+(bS*7)].HasLibertiesBool.HasLeft = false
					result[i+(bS*7)].HasLibertiesBool.HasRight = false
					result[i-(bS*7)].Liberties = 0
					result[i-(bS*7)].HasLibertiesBool.HasLeft = false
					result[i-(bS*7)].HasLibertiesBool.HasRight = false
					result[i+(bS*8)].Liberties = 0
					result[i+(bS*8)].HasLibertiesBool.HasLeft = false
					result[i+(bS*8)].HasLibertiesBool.HasRight = false
					result[i-(bS*8)].Liberties = 0
					result[i-(bS*8)].HasLibertiesBool.HasLeft = false
					result[i-(bS*8)].HasLibertiesBool.HasRight = false
					result[i+(bS*9)].Liberties = 0
					result[i+(bS*9)].HasLibertiesBool.HasLeft = false
					result[i-(bS*9)].Liberties = 0
					result[i-(bS*9)].HasLibertiesBool.HasRight = false
				}

			}

		}
	}
	println("DONE CHECKING ALL STONE LIBERTIES", len(result))
	end := time.Now()
	duration := end.Sub(start).Microseconds()
	println("CALCULATE LIBERTIES TOOK MICROSEC : ", duration)
	var sent bool = false
	for !sent {
		select {
		case chanNewLibertiesResult <- result:
			sent = true
		default:
			time.Sleep(time.Millisecond * 16)
			sent = false
		}
	}
}

func (b *Board) CheckIndividualLiberties(SurroundingIndex int, i int) bool {
	if b.Grid[SurroundingIndex].State == b.Grid[i].State || b.Grid[SurroundingIndex].State == 0 {
		return true
	}
	return false
}
