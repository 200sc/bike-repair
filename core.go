package main

import (
	"github.com/oakmound/oak"
)

func main() {
	oak.AddScene("day", sceneStart, sceneLoop, sceneEnd)
	oak.Init("day")
}

func sceneStart(prev string, input interface{}) {

}

func sceneLoop() bool {
	return true
}

func sceneEnd() (string, *oak.SceneResult) {
	return "day", nil
}
