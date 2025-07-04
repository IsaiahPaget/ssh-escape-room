package main

import (
	"github.com/gdamore/tcell/v2"
	"log"
	"os"
	"fmt"
)

func debugLog(v any) {
	f, err := os.OpenFile("debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("logToFile error: %v", err)
		return
	}
	defer f.Close()

	_, err = fmt.Fprintf(f,"%v\n", v)
	if err != nil {
		log.Printf("logToFile write error: %v", err)
	}
}

// sig_winch for scaling

func drawText(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	row := y1
	col := x1
	for _, r := range []rune(text) {
		s.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
}

func buttons(key rune, _x, _y int, env *Env) {
	if key == 'w' { env.player.postion.y -= 1 }
	if key == 'a' { env.player.postion.x -= 1 }
	if key == 's' { env.player.postion.y += 1 }
	if key == 'd' { env.player.postion.x += 1 }
}
func drawBox(s tcell.Screen, position Vector2D, width int, height int, style tcell.Style) {
	x2 := position.x + (width * 2) - 1 // times 2 because of the terminals taller pixels
	y2 := position.y + height - 1

	// Fill background
	for row := position.y; row <= y2; row++ {
		for col := position.x; col <= x2; col++ {
			s.SetContent(col, row, ' ', nil, style)
		}
	}
}

type Vector2D struct {
	x int
	y int
}
type Player struct {
	postion Vector2D
}
type Env struct {
	ox     int
	oy     int
	screen tcell.Screen
	player Player
}

func drawPlayer(env Env) {
	style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlueViolet)
	debugLog(env.player.postion)
	drawBox(env.screen, env.player.postion, 1, 1, style)
}
func display(env Env) {
	env.screen.Clear()
	drawPlayer(env)

	env.screen.Sync()
}

func main() {
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)

	player := Player{
		postion: Vector2D{
			x: 30,
			y: 30,
		},
	}
	// Event loop
	env := Env{
		ox:     -1,
		oy:     -1,
		player: player,
	}

	// Initialize screen target size 1024 by 512
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	env.screen = s
	s.SetStyle(defStyle)
	s.EnableMouse()
	s.EnablePaste()
	s.Clear()

	quit := func() {
		// You have to catch panics in a defer, clean up, and
		// re-raise them - otherwise your application can
		// die without leaving any diagnostic trace.
		maybePanic := recover()
		s.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()

	// Here's how to get the screen size when you need it.
	// xmax, ymax := s.Size()

	// Here's an example of how to inject a keystroke where it will
	// be picked up by the next PollEvent call.  Note that the
	// queue is LIFO, it has a limited length, and PostEvent() can
	// return an error.
	// s.PostEvent(tcell.NewEventKey(tcell.KeyRune, rune('a'), 0))

	for {
		// Update screen
		s.Show()

		// Poll event
		ev := s.PollEvent()
		// Process event
		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				return
			} else if ev.Key() == tcell.KeyCtrlL {
				s.Sync()
			} else if ev.Rune() == 'C' || ev.Rune() == 'c' {
				s.Clear()
			} else {
				buttons(ev.Rune(), 0, 0, &env)
			}
		}
		display(env)
	}
}
