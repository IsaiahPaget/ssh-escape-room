package main

import (
	"math"

	"github.com/gdamore/tcell/v2"
)

const ENTITY_TYPE_PLAYER = "player_type"

func InitPlayer(env *Environment) {

	env.CreateEntity(
		ENTITY_TYPE_PLAYER,
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
				env.player = player
			}

			player.on_update = func() {
				if env.input_rune == 'e' {
					player.postion.x += player.rotation.delta.x * 1
					player.postion.y += player.rotation.delta.y * 1
				}
				if env.input_rune == 's' {
					player.rotation.angle -= 0.05
					if player.rotation.angle < 0 {
						player.rotation.angle += 2 * PI
					}
					player.rotation.delta.x = float32(math.Cos(float64(player.rotation.angle) * 5))
					player.rotation.delta.y = float32(math.Sin(float64(player.rotation.angle) * 5))
				}
				if env.input_rune == 'd' {
					player.postion.x -= player.rotation.delta.x * 1
					player.postion.y -= player.rotation.delta.y * 1
				}
				if env.input_rune == 'f' {
					player.rotation.angle += 0.05
					if player.rotation.angle > 2*PI {
						player.rotation.angle -= 2 * PI
					}
					player.rotation.delta.x = float32(math.Cos(float64(player.rotation.angle) * 5))
					player.rotation.delta.y = float32(math.Sin(float64(player.rotation.angle) * 5))
				}
			}
			player.on_draw = func() {
				style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlueViolet)
				debug_style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorRed)

				DrawBox(env.v_screen, player.postion, 2, 2, style)

				DrawBox(
					env.v_screen,
					Vector2D{
						x: player.postion.x + player.rotation.delta.x*2,
						y: player.postion.y + player.rotation.delta.y*2,
					},
					1,
					1,
					debug_style,
				)
			}
		},
	)
}
