package main

import (
	"github.com/gdamore/tcell/v2"
	"math"
)

const (
	EntityTypePlayer = "player_type"
)

type Player struct {
	postion  Vector2D
	rotation Rotation
}

func DrawPlayer(env Environment) {
	style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlueViolet)
	debug_style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorRed)

	player_arr, _ := env.GetEntities(EntityTypePlayer)
	if len(player_arr) < 1 {
		panic("No Player!")
	}
	player := player_arr[0]

	DrawBox(env.v_screen, player.postion, 2, 2, style)

	DrawBox(
		env.v_screen,
		Vector2D{
			x: player.postion.x + player.rotation.delta.x*5,
			y: player.postion.y + player.rotation.delta.y*5,
		},
		1,
		1,
		debug_style,
	)
}

func InitPlayer(env *Environment) {

	env.AddEntity(
		Entity{
			entity_type: EntityTypePlayer,
			postion: Vector2D{
				x: 30,
				y: 30,
			},
			rotation: Rotation{
				delta: Vector2D{
					x: .5,
					y: .5,
				},
			},
		},
		func(player *Entity) {

			player.on_init = func() {
				player.postion = Vector2D{
					x: 30,
					y: 30,
				}
				player.rotation.delta = Vector2D{
					x: .5,
					y: .5,
				}
			}

			player.on_update = func() {
				if env.input_rune == 'e' {
					player.postion.x += player.rotation.delta.x * 2
					player.postion.y += player.rotation.delta.y * 2
				}
				if env.input_rune == 's' {
					player.rotation.angle -= 0.1
					if player.rotation.angle < 0 {
						player.rotation.angle += 2 * PI
					}
					player.rotation.delta.x = float32(math.Cos(float64(player.rotation.angle) * 5))
					player.rotation.delta.y = float32(math.Sin(float64(player.rotation.angle) * 5))
				}
				if env.input_rune == 'd' {
					player.postion.x -= player.rotation.delta.x * 2
					player.postion.y -= player.rotation.delta.y * 2
				}
				if env.input_rune == 'f' {
					player.rotation.angle += 0.1
					if player.rotation.angle > 2*PI {
						player.rotation.angle -= 2 * PI
					}
					player.rotation.delta.x = float32(math.Cos(float64(player.rotation.angle) * 5))
					player.rotation.delta.y = float32(math.Sin(float64(player.rotation.angle) * 5))
				}
			}
		},
	)
}
