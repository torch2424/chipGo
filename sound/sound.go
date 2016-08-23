//Using this lib: https://godoc.org/golang.org/x/mobile/exp/audio#Player.Play
//And this example https://github.com/golang/mobile/blob/master/example/audio/main.go

package sound

import (
	"os"
	"golang.org/x/mobile/exp/audio"
)

//Our sound file path
const soundPath = "sound/Blip.wav"

//Debug mode
var debugMode bool

//Our AudioPlayer struct for accessing the class in a state
type AudioPlayer struct {
	//Our player object
	sound audio.Player
}

//Function to intialize our audioplayer
func NewAudioPlayer(debug bool) AudioPlayer {

	//Print that audio is being initalized
	print("\nInitializing Audio...\n")

	rc, err := os.Open(soundPath)
	if err != nil {
		panic(err)
	}
	player, err := audio.NewPlayer(rc, 0, 0)
	if err != nil {
		panic(err)
	}

	//Set our audio player
	audioPlayer := AudioPlayer{sound: *player}

	//Set debug mode
	debugMode = debug

	return audioPlayer
}

func PlayBlip(audioPlayer AudioPlayer) {
	//Check if we are already beeping
	//Using mobile audio source: https://sourcegraph.com/github.com/golang/mobile/-/def/GoPackage/github.com/golang/mobile/exp/audio/-/Playing
	if audioPlayer.sound.State() != audio.Playing {
		//Play the audio sound
		audioPlayer.sound.Play()
	}
}
