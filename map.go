package main

import (
	"github.com/gdamore/tcell/v2"
)

const EntityTypeMap = "map_type"

func InitGameMap(env *Environment) {

	env.CreateEntity(
		EntityTypeMap,
		func(game_map *Entity) {

			game_map.on_init = func() {
				game_map.width = 8
				game_map.height = 8
				game_map.level_data = []int{
					1, 1, 1, 1, 1, 1, 1, 1,
					1, 0, 0, 0, 0, 0, 0, 1,
					1, 0, 1, 1, 1, 0, 0, 1,
					1, 0, 1, 0, 1, 0, 0, 1,
					1, 0, 1, 0, 1, 0, 0, 1,
					1, 0, 0, 0, 0, 0, 0, 1,
					1, 0, 0, 0, 0, 0, 0, 1,
					1, 1, 1, 1, 1, 1, 1, 1,
				}
				game_map.block_width = 4
				game_map.block_height = 4

			}

			game_map.on_update = func() {
			}
			game_map.on_draw = func() {
				wall_style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorGreen)
				floor_style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorWhite)
				map_width := game_map.width
				map_height := game_map.height
				map_block_width := game_map.block_width
				map_block_height := game_map.block_height
				map_level_data := game_map.level_data

				for row := range map_width {
					for col := range map_height {
						block_style := floor_style

						if map_level_data[GetFlatMapIndex(row, map_width, col)] == 1 {
							block_style = wall_style
						}

						position := Vector2D{
							x: float32(col * 2 * map_block_height),
							y: float32(row * map_block_width),
						}

						DrawBox(env.v_screen, position, map_block_width, map_block_height, block_style)
					}
				}
			}
		},
	)
}
