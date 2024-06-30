package mobile

import (
	"embed"
	"gogoboardgame/board"
	"gogoboardgame/game"
	"gogoboardgame/input"

	"github.com/hajimehoshi/ebiten/v2/mobile"
)

//go:embed resources
var Resources embed.FS

// run this to build for android
// ebitenmobile bind -target android -javapkg com.ourbunka.goboardgame -o ./android/goboardgame.aar .
func init() {
	board.Resources = Resources
	board.UseEmbeded = true
	board.Build = "ANDROID"
	input.Resources = Resources
	input.UseEmbeded = true
	mobile.SetGame(game.NewGame())
}

func Dummy() {}
