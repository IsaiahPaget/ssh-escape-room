package main

import (
	"fmt"
	"log"
	"math"
	"os"

	"github.com/gdamore/tcell/v2"
)

const PI = 3.1415926535
const virtualWidth = 1024
const virtualHeight = 512

func debugLog(v any) {
	f, err := os.OpenFile("debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("logToFile error: %v", err)
		return
	}
	defer f.Close()

	_, err = fmt.Fprintf(f, "%v\n", v)
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

// TODO: frame buffer that is 1024 x 512 this will solve flickering, and will allow the rendering
// to downscale to the size of a terminal but while allowing the game math to be more like a normal game

func player_actions(key rune, env *Env) {
	if key == 'e' {
		env.player.postion.x += int(env.player.rotation.delta.x * 2)
		env.player.postion.y += int(env.player.rotation.delta.y * 2)
	}
	if key == 's' {
		env.player.rotation.angle -= 0.1
		if env.player.rotation.angle < 0 {
			env.player.rotation.angle += 2 * PI
		}
		env.player.rotation.delta.x = float32(math.Cos(float64(env.player.rotation.angle) * 5))
		env.player.rotation.delta.y = float32(math.Sin(float64(env.player.rotation.angle) * 5))
	}
	if key == 'd' {
		env.player.postion.x -= int(env.player.rotation.delta.x * 2)
		env.player.postion.y -= int(env.player.rotation.delta.y * 2)
	}
	if key == 'f' {
		env.player.rotation.angle += 0.1
		if env.player.rotation.angle > 2*PI {
			env.player.rotation.angle -= 2 * PI
		}
		env.player.rotation.delta.x = float32(math.Cos(float64(env.player.rotation.angle) * 5))
		env.player.rotation.delta.y = float32(math.Sin(float64(env.player.rotation.angle) * 5))
	}
}
func game_actions(key rune, _x, _y int, env *Env) {
	player_actions(key, env)
}
func drawBox(s VirtualScreen, position Vector2DInt, width int, height int, style tcell.Style) {
	x2 := position.x + (width * 2) - 1 // times 2 because of the terminals taller pixels
	y2 := position.y + height - 1

	// Fill background
	for row := position.y; row <= y2; row++ {
		for col := position.x; col <= x2; col++ {
			s.SetContent(col, row, style)
		}
	}
}

type Map struct {
	width        int
	height       int
	block_width  int
	block_height int
	area         int
	level_data   []int
}

type Rotation struct {
	delta Vector2DFloat32
	angle float32
}

type Vector2DFloat32 struct {
	x float32
	y float32
}
type Vector2DInt struct {
	x int
	y int
}
type Player struct {
	postion  Vector2DInt
	rotation Rotation
}
type VirtualScreen struct {
	buffer []Pixel
	width  int
	height int
}

func (virtual_screen VirtualScreen) SetContent(x int, y int, style tcell.Style) {
	virtual_screen.buffer[y*virtual_screen.width+x] = Pixel{
		style: style,
	}
}

const (
	EntityTypePlayer = "player_type"
)

type Entity struct {
	entity_type int
}

type Env struct {
	ox       int
	oy       int
	v_screen VirtualScreen
	t_screen tcell.Screen
	player   Player
	game_map Map
}

type Pixel struct {
	style tcell.Style // Color, bold, etc.
}

func drawMap(env Env) {
	wall_style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorGreen)
	floor_style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorWhite)
	map_width := env.game_map.width
	map_height := env.game_map.height
	map_block_width := env.game_map.block_width
	map_block_height := env.game_map.block_height
	map_level_data := env.game_map.level_data

	for row := range map_width {
		for col := range map_height {
			block_style := floor_style

			if map_level_data[row*map_width+col] == 1 {
				block_style = wall_style
			}

			drawBox(env.v_screen, Vector2DInt{col * 2 * map_block_height, row * map_block_width}, map_block_width, map_block_height, block_style)
		}
	}
}
func drawPlayer(env Env) {
	style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlueViolet)
	debug_style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorRed)
	drawBox(env.v_screen, env.player.postion, 2, 2, style)

	drawBox(
		env.v_screen,
		Vector2DInt{
			x: env.player.postion.x + int(env.player.rotation.delta.x)*5,
			y: env.player.postion.y + int(env.player.rotation.delta.y)*5,
		},
		1,
		1,
		debug_style,
	)
}

func renderBuffer(env Env) {
	for row := range env.v_screen.width {
		for col := range env.v_screen.height {
			env.t_screen.SetContent(col, row, ' ', nil, env.v_screen.buffer[row*env.v_screen.width+col].style)
		}
	}
}

func display(env Env) {
	env.t_screen.Clear()
	drawMap(env)
	drawPlayer(env)

	renderBuffer(env)
	env.t_screen.Sync()
}

func main() {
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)

	player := Player{
		postion: Vector2DInt{
			x: 30,
			y: 30,
		},
		rotation: Rotation{
			delta: Vector2DFloat32{
				x: .5,
				y: .5,
			},
		},
	}

	game_map := Map{
		width:  8,
		height: 8,
		area:   64,
		level_data: []int{
			1, 1, 1, 1, 1, 1, 1, 1,
			1, 0, 0, 0, 0, 0, 0, 1,
			1, 0, 1, 1, 1, 0, 0, 1,
			1, 0, 1, 0, 1, 0, 0, 1,
			1, 0, 1, 0, 1, 0, 0, 1,
			1, 0, 0, 0, 0, 0, 0, 1,
			1, 0, 0, 0, 0, 0, 0, 1,
			1, 1, 1, 1, 1, 1, 1, 1,
		},
		block_width:  4,
		block_height: 4,
	}
	// Initialize screen target size 1024 by 512
	t_screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := t_screen.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	xmax, ymax := t_screen.Size()
	virtual_screen := VirtualScreen{
		buffer: make([]Pixel, virtualWidth*virtualHeight),
		width:  xmax,
		height: ymax,
	}

	// Event loop
	env := Env{
		ox:       -1,
		oy:       -1,
		player:   player,
		game_map: game_map,
		v_screen: virtual_screen,
		t_screen: t_screen,
	}

	env.t_screen.SetStyle(defStyle)
	env.t_screen.EnableMouse()
	env.t_screen.EnablePaste()
	env.t_screen.Clear()

	quit := func() {
		// You have to catch panics in a defer, clean up, and
		// re-raise them - otherwise your application can
		// die without leaving any diagnostic trace.
		maybePanic := recover()
		env.t_screen.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()

	// Here's how to get the screen size when you need it.

	// Here's an example of how to inject a keystroke where it will
	// be picked up by the next PollEvent call.  Note that the
	// queue is LIFO, it has a limited length, and PostEvent() can
	// return an error.
	// s.PostEvent(tcell.NewEventKey(tcell.KeyRune, rune('a'), 0))

	for {
		// Update screen
		env.t_screen.Show()

		// Poll event
		ev := env.t_screen.PollEvent()
		// Process event
		switch ev := ev.(type) {
		case *tcell.EventResize:
			env.t_screen.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				return
			} else if ev.Key() == tcell.KeyCtrlL {
				env.t_screen.Sync()
			} else if ev.Rune() == 'C' || ev.Rune() == 'c' {
				env.t_screen.Clear()
			} else {
				game_actions(ev.Rune(), 0, 0, &env)
			}
		}
		display(env)
	}
}
