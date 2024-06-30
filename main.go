package main

import (
	"embed"
	"gogoboardgame/board"
	"gogoboardgame/game"
)

//go:embed resources
var Resources embed.FS

func main() {
	board.UseEmbeded = false
	game.RunGame()
}
