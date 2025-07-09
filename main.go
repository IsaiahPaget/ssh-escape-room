package main

import (
	"errors"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/google/uuid"
	"log"
	"os"
)

const PI = 3.1415926535
const VIRTUAL_WIDTH = 1024
const VIRTUAL_HEIGHT = 512
const MAX_ENTITIES = 2048

type Rotation struct {
	delta Vector2D
	angle float32
}

type Vector2D struct {
	x float32
	y float32
}

type VirtualScreen struct {
	buffer []Pixel
	width  int
	height int
}

func (virtual_screen VirtualScreen) SetContent(x int, y int, style tcell.Style) {
	virtual_screen.buffer[GetFlatMapIndex(y, virtual_screen.width, x)] = Pixel{
		style: style,
	}
}

type Entity struct {
	id          uuid.UUID
	entity_type string
	on_init     func()
	on_update   func()
	on_destroy  func()
	on_draw     func()
	postion     Vector2D
	rotation    Rotation
	width       int
	height      int
	// map stuff
	block_width  int
	block_height int
	level_data   []int
}

type Environment struct {
	v_screen   VirtualScreen
	t_screen   tcell.Screen
	entities   []Entity
	input_rune rune      // this is the charactor code like "w" as in wasd
	input_key  tcell.Key // this can be used for ctrl+c and such
}

func (env Environment) GetEntities(entityType string) ([]Entity, error) {
	if entityType == "" {
		return nil, errors.New("entityType cannot be empty")
	}

	entities := []Entity{}

	for _, e := range env.entities {
		if e.entity_type == entityType {
			entities = append(entities, e)
		}
	}

	return entities, nil
}

func (env *Environment) CreateEntity(entity_type string, setup func(*Entity)) {
	env.entities = append(env.entities, Entity{
		id:          uuid.New(),
		entity_type: entity_type,
	})
	idx := len(env.entities) - 1
	if idx < 0 || idx > MAX_ENTITIES {
		panic("Index out of bounds")
	}
	setup(&env.entities[idx])
	env.entities[idx].on_init()
}

type Pixel struct {
	style tcell.Style // Color, bold, etc.
}

func DebugLog(v any) {
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

// func DrawText(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
// 	row := y1
// 	col := x1
// 	for _, r := range []rune(text) {
// 		s.SetContent(col, row, r, nil, style)
// 		col++
// 		if col >= x2 {
// 			row++
// 			col = x1
// 		}
// 		if row > y2 {
// 			break
// 		}
// 	}
// }

func GameActions(env *Environment) {
	for _, entity := range env.entities {
		entity.on_update()
	}
}
func DrawBox(s VirtualScreen, position Vector2D, width int, height int, style tcell.Style) {
	x2 := int(position.x) + (width * 2) - 1 // times 2 because of the terminals taller pixels
	y2 := int(position.y) + height - 1

	// Fill background
	for row := int(position.y); row <= y2; row++ {
		for col := int(position.x); col <= x2; col++ {
			s.SetContent(col, row, style)
		}
	}
}

// func DrawMap(env Environment) {
// 	wall_style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorGreen)
// 	floor_style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorWhite)
// 	map_width := env.game_map.width
// 	map_height := env.game_map.height
// 	map_block_width := env.game_map.block_width
// 	map_block_height := env.game_map.block_height
// 	map_level_data := env.game_map.level_data
//
// 	for row := range map_width {
// 		for col := range map_height {
// 			block_style := floor_style
//
// 			if map_level_data[GetFlatMapIndex(row, map_width, col)] == 1 {
// 				block_style = wall_style
// 			}
//
// 			position := Vector2D{
// 				x: float32(col * 2 * map_block_height),
// 				y: float32(row * map_block_width),
// 			}
//
// 			DrawBox(env.v_screen, position, map_block_width, map_block_height, block_style)
// 		}
// 	}
// }

func GetFlatMapIndex(row, width, col int) int {
	if row < 0 {
		row = 0
	}
	if width < 0 {
		width = 0
	}
	return row*width + col
}

func RenderBuffer(env Environment) {
	for row := range env.v_screen.width {
		for col := range env.v_screen.height {
			env.t_screen.SetContent(col, row, ' ', nil, env.v_screen.buffer[GetFlatMapIndex(row, env.v_screen.width, col)].style)
		}
	}
}

func DrawEntities(env Environment) {
	for _, entity := range env.entities {
		entity.on_draw()
	}
}

func Display(env Environment) {
	env.t_screen.Clear()
	DrawEntities(env)

	RenderBuffer(env)
	env.t_screen.Sync()
}

func InitEntities(env *Environment) {
	env.entities = []Entity{}
	// The order of these function determines the 'z-index' of the stuff
	InitGameMap(env)
	InitPlayer(env)
}

func InitGame(env *Environment) {

	// Initialize screen
	t_screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := t_screen.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	xmax, ymax := t_screen.Size()
	virtual_screen := VirtualScreen{
		buffer: make([]Pixel, VIRTUAL_WIDTH*VIRTUAL_HEIGHT),
		width:  xmax,
		height: ymax,
	}

	env.v_screen = virtual_screen
	env.t_screen = t_screen

	InitEntities(env)

	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	env.t_screen.SetStyle(defStyle)
	env.t_screen.EnableMouse()
	env.t_screen.EnablePaste()
	env.t_screen.Clear()
}

func main() {
	// Event loop
	env := Environment{}

	InitGame(&env)

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
				env.input_rune = ev.Rune()
				env.input_key = ev.Key()
			}
		}
		GameActions(&env)
		Display(env)
	}
}
