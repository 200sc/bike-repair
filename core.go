package main

import "github.com/oakmound/oak"
import "github.com/200sc/bike-repair/bike"

func main() {
	oak.AddScene("day", bike.DayStart, bike.DayLoop, bike.DayEnd)
	oak.Init("day")
}
