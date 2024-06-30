package mobile

import (
	"embed"
	"gogoboardgame/board"
	"gogoboardgame/game"

	"github.com/hajimehoshi/ebiten/v2/mobile"
)

//go:embed resources
var Resources embed.FS

// run this to build for android
// ebitenmobile bind -target android -javapkg com.ourbunka.goboardgame -o ./android/goboardgame.aar .
// comment this when build for desktop
func init() {
	board.Resources = Resources
	board.UseEmbeded = true
	board.Build = "ANDROID"
	mobile.SetGame(game.NewGame())
}

func Dummy() {}
