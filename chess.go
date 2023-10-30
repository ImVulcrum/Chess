package main

import (
	"fmt"
	"image"
	"os"
	"strconv"
	"time"

	"./buttons"
	"./clipboard"
	gfx "./gfxw"
	"./imaging"
	"./parser"
	"./pieces"
	"./sliders"
	"./textbox"
	"./time_counter"
)

func main() {
	fmt.Println("start game")
	var restart_window bool = false //decides if a new window is created on initialize --> false by default so that one window will be created

	var deselect_piece_after_clicking = false
	var start_window_width uint16 = 600 //the start window is technically scalable but with resolutions higher than 800 graphical bugs are occuring due to the font size
	game_timer, friendly_game, duration_of_premove_animation, use_clipboard_as_premoves, name_player_white, name_player_black, a, troll_mode := start_menu(start_window_width / 3 * 2)

restart_marker:
	var image_location string = set_image_string(troll_mode) //troll mode decides which picture should be used for the pieces
	var game_history_can_be_changed bool = friendly_game
	var w_x, w_y uint16 = 10 * a, 8 * a
	var white_is_current_player bool = false
	var player_change bool = true
	var current_king_index int
	var current_player_has_no_legal_moves bool
	var check bool
	var current_piece pieces.Piece
	var temp_current_piece pieces.Piece
	var piece_index int
	var current_legal_moves [][3]uint16
	var moves_counter int16
	var current_field [2]uint16
	var dragging bool
	var ending_premoves bool = true
	var piece_is_selected uint16 = 64
	var premoves string = get_clipboard_if_asked(use_clipboard_as_premoves)
	var promotion uint16 = 0
	var move_string string
	var take string = ""
	m_channel := make(chan [4]int16, 1)
	gor_status := make(chan bool, 1) //needs to be a buffered channel, indicated by the one, otherwise the program will hold until the channel is empty after puting something to it

	pieces_a, white_king_index, black_king_index, one_move_back, one_move_forward, restart_button, pause_button, save_button, moves_a, white_time_counter, black_time_counter, pgn_moves_a := initialize(w_x, w_y, a, restart_window, game_timer, name_player_white, name_player_black, image_location)
	premoves_array := parser.Create_Array_Of_Moves(premoves) //get the premoves array

	//draw_pieces(pieces_a, w_x, w_y, a) //draw pieces for the first time, actually not really needed
	display_message(2, a, false) //starting message

	go mouse_handler(m_channel, gor_status) //activate the mouse handler so that mouse input can be obtained

	for { //gameloop

		if player_change { //after a player change do the following:
			pieces_a, moves_a, pgn_moves_a = game_history_handler(moves_counter, moves_a, pieces_a, game_history_can_be_changed, move_string, pgn_moves_a) //add the current move to pgn moves array as well as to the normal moves array (important for back and forward)

			white_is_current_player, current_king_index, player_change, moves_counter = change_player(white_is_current_player, white_king_index, black_king_index, moves_counter) //this function sets the current_king_index
			pieces_a, current_player_has_no_legal_moves = pieces.Calc_Moves_With_Check(pieces_a, moves_counter, current_king_index)                                               //calc the move
			check = pieces_a[current_king_index].(*pieces.King).Is_In_Check(pieces_a, moves_counter)                                                                              //check if the current king is in check

			if !use_clipboard_as_premoves || duration_of_premove_animation != 0 { //wenn keine premoves mehr da sind oder premoves da sind und diese gezeichnet werden sollen dann wird das board gezeichnet
				draw_board(a, w_x, w_y, current_piece, current_legal_moves, pieces_a, false, current_king_index, check)
				draw_moves_sidebar(a, moves_counter-1, pgn_moves_a)
			}

			if restart_window = restart_handler(current_player_has_no_legal_moves, check, a, white_is_current_player, gor_status); restart_window { //check if the re is a checkmate or stalemate
				goto restart_marker
			}

			if !use_clipboard_as_premoves && !friendly_game && int(moves_counter) == len(moves_a) { //only start the timers if we're not premoving and it isn'nt a friendly game and no older game state is reviewed
				time_handler(white_is_current_player, white_time_counter, black_time_counter, pause_button)
			}
		}

		if len(premoves_array) > 0 { //if there are premoves the program premoves
			var piece_promoted_to string
			piece_executing_move, index_of_move, piece_promoting_to := parser.Get_Correct_Move(premoves_array[0], pieces_a, current_king_index) //get the move from the first element in the premoves array

			if piece_executing_move != 64 { //if there is no move matching the specifications, this code won't be excuted, instead the premove sequence will end at this point (else statement)
				var original_pos [2]uint16 = pieces_a[piece_executing_move].Give_Pos() //get the position so that the pgn string can be recreated for the sidebar

				pieces_a, promotion, take = pieces.Move_Piece_To(pieces_a[piece_executing_move], pieces_a[piece_executing_move].Give_Legal_Moves()[index_of_move], moves_counter, pieces_a)
				if promotion != 64 {
					pieces_a, piece_promoted_to = pawn_promotion(w_x, w_y, a, piece_executing_move, pieces_a, piece_promoting_to, m_channel)
					piece_promoted_to = "=" + piece_promoted_to
				}
				move_string = get_move_string(piece_executing_move, original_pos, piece_promoted_to, take, pieces_a) //get the pgn string for the sidebar

				premoves_array = premoves_array[1:] //remove the first element of the premoves array
				player_change = true
				time.Sleep(time.Duration(duration_of_premove_animation) * time.Millisecond)
			} else {
				premoves_array = nil
			}
		} else if ending_premoves { //soll im optimalfall nur einmal ausgeführt werden
			draw_board(a, w_x, w_y, current_piece, current_legal_moves, pieces_a, false, current_king_index, check)
			use_clipboard_as_premoves = false
			ending_premoves = false
			draw_moves_sidebar(a, moves_counter-1, pgn_moves_a)
			//<-m_channel //there is a bug where you can the mouse_handler already reads when the program is still premoving
		}

		if !use_clipboard_as_premoves { //wenn die andere restart methode ausgeführt wird, dann sollt hier noch noch überprüft werden ob restart == false ist, da es sonst zu fehlern kommt
			select {
			case mouse_input := <-m_channel:

				if dragging && !(mouse_input[0] == 1 && mouse_input[1] == 0) { //wenn nichts gehalten wird dann löschen
					draw_board(a, w_x, w_y, current_piece, current_legal_moves, pieces_a, true, current_king_index, check)
					dragging = false
				}

				if mouse_input[0] == 1 && mouse_input[1] == 1 {
					if one_move_back.Is_Clicked(uint16(mouse_input[2]), uint16(mouse_input[3])) {
						if moves_counter >= 2 {
							moves_counter = moves_counter - 2
							pieces_a = pieces.Copy_Array(moves_a[moves_counter])
							player_change = true
						}
					} else if one_move_forward.Is_Clicked(uint16(mouse_input[2]), uint16(mouse_input[3])) {
						if int(moves_counter) < len(moves_a) {
							pieces_a = pieces.Copy_Array(moves_a[moves_counter])
							player_change = true
						}
					} else if restart_button.Is_Clicked(uint16(mouse_input[2]), uint16(mouse_input[3])) {
						restart_window = true
						gor_status <- true
						goto restart_marker
					} else if save_button.Is_Clicked(uint16(mouse_input[2]), uint16(mouse_input[3])) {
						parser.Write_PGN_File(pgn_moves_a, name_player_white, name_player_black)
					} else if int(moves_counter) == len(moves_a) && pause_button.Is_Clicked(uint16(mouse_input[2]), uint16(mouse_input[3])) {
						if pause_button.Switch(38, 37, 34) {
							white_time_counter.Stop_Counting()
							black_time_counter.Stop_Counting()
						} else if white_is_current_player {
							white_time_counter.Init_Counting()
						} else if !white_is_current_player {
							black_time_counter.Init_Counting()
						}
					} else {
						current_field = calc_field(a, uint16(mouse_input[2]), uint16(mouse_input[3]), 0)
						temp_current_piece, piece_index = get_current_piece(pieces_a, current_field)

						//auswählen eines pieces
						if temp_current_piece != nil && temp_current_piece.Is_White_Piece() == white_is_current_player && (!deselect_piece_after_clicking || (current_piece == nil || temp_current_piece.Give_Pos() != current_piece.Give_Pos())) {
							current_piece = temp_current_piece
							current_legal_moves = current_piece.Give_Legal_Moves()
							promotion = 0
							draw_board(a, w_x, w_y, current_piece, current_legal_moves, pieces_a, true, current_king_index, check)
							piece_is_selected = uint16(piece_index)

						} else if piece_is_selected != 64 {
							//überprüfen ob auf ein Feld in Legal Moves geklickt wurde
							pieces_a, piece_is_selected, player_change, move_string = move_if_current_field_is_in_legal_moves(current_field, pieces_a, piece_is_selected, a, w_x, w_y, current_king_index, check, moves_counter, m_channel)

							//Deselect sobald ein Piece ausgewählt ist: wenn auf kein Piece geklickt wurde, oder auf ein gegnerisches Piece geklickt wurde, wird das ausgewählte Piece deselected
							//mit der Option deselect_piece_after_clicking kann eingestellt werde, ob auch deselected werden soll wenn auf dasslebe piece nocheinmal geklickt wurde
							if (temp_current_piece == nil) || (temp_current_piece != nil && (temp_current_piece.Give_Pos() == current_piece.Give_Pos() || temp_current_piece.Is_White_Piece() != white_is_current_player)) {
								current_piece = nil
								draw_board(a, w_x, w_y, current_piece, current_legal_moves, pieces_a, false, current_king_index, check)
								piece_is_selected = 64
							}
						}
					}

				} else if mouse_input[1] == -1 && mouse_input[0] == 1 {
					current_field = calc_field(a, uint16(mouse_input[2]), uint16(mouse_input[3]), 0)
					temp_current_piece, _ = get_current_piece(pieces_a, current_field)

					if piece_is_selected != 64 {
						//überprüfen ob in current legal moves
						pieces_a, piece_is_selected, player_change, move_string = move_if_current_field_is_in_legal_moves(current_field, pieces_a, piece_is_selected, a, w_x, w_y, current_king_index, check, moves_counter, m_channel)

						//Deselect sobald ein Piece ausgewählt ist: wenn auf keinem Piece losgelassen wurde oder auf ein generisches Piece losgelassen wurde oder wenn auf einem eigenen piece losgelassen wurde
						if (temp_current_piece == nil) || (temp_current_piece != nil && ((temp_current_piece.Is_White_Piece() != white_is_current_player) || (temp_current_piece.Is_White_Piece() == white_is_current_player && temp_current_piece.Give_Pos() != current_piece.Give_Pos()))) {
							current_piece = nil
							draw_board(a, w_x, w_y, current_piece, current_legal_moves, pieces_a, false, current_king_index, check)
							piece_is_selected = 64
						}
					}
				} else if mouse_input[1] == 0 && mouse_input[0] == 1 && current_piece != nil && current_piece.Is_White_Piece() == white_is_current_player {
					//wenn die maustaste gehalten wird, wird ein ghosttpiece gemalt, welches der maus folgt
					dragging = true
					pieces.Draw_To_Point(current_piece, w_x, w_y, a, uint16(mouse_input[2]), uint16(mouse_input[3]), -int16(a/2), -int16(a/2), 50, uint16(mouse_input[2]))
				}

			default:
				time.Sleep(5 * time.Millisecond)
			}
			if !friendly_game && draw_timers(white_time_counter, black_time_counter, a) {
				display_message(0, a, white_is_current_player)
				restart_window = true
				gor_status <- true
				goto restart_marker
			}
		}
	}
}

func start_menu(start_window_size uint16) (int64, bool, int, bool, string, string, uint16, bool) { //this function is supposed to define the buttons, slider and text boxes and return the values afterwards to the section of self definable variables
	gfx.Fenster(start_window_size/2*3, start_window_size)
	gfx.Fenstertitel("Chess")

	var bg_color [3]uint8 = [3]uint8{24, 24, 24}
	var primary_color [3]uint8 = [3]uint8{70, 70, 70}
	var secondary_color [3]uint8 = [3]uint8{204, 204, 204}
	var highlight_color1 [3]uint8 = [3]uint8{150, 0, 4}
	var highlight_color2 [3]uint8 = [3]uint8{0, 94, 4}

	gfx.Stiftfarbe(bg_color[0], bg_color[1], bg_color[2])
	gfx.Vollrechteck(0, 0, start_window_size/2*3, start_window_size)

	gfx.Stiftfarbe(secondary_color[0], secondary_color[1], secondary_color[2])
	gfx.SetzeFont("./resources/fonts/punk.ttf", int(start_window_size/12))
	gfx.SchreibeFont(start_window_size/2*3/2/12*10, start_window_size/90, "CHESS")

	var name_player_white textbox.Box = textbox.New(start_window_size/20, start_window_size/6, start_window_size/14, start_window_size/17*10+start_window_size/40, primary_color, secondary_color, highlight_color1, 19, "enter white name...")
	var name_player_black textbox.Box = textbox.New(start_window_size/20, start_window_size/4, start_window_size/14, start_window_size/17*10+start_window_size/40, primary_color, secondary_color, highlight_color1, 19, "enter black name...")
	name_player_white.Draw()
	name_player_black.Draw()

	gfx.SetzeFont("./resources/fonts/unispace.ttf", int(start_window_size/20))
	gfx.Stiftfarbe(secondary_color[0], secondary_color[1], secondary_color[2])
	gfx.SchreibeFont(start_window_size/17*10+start_window_size/8, start_window_size/6, "name of white player")
	gfx.SchreibeFont(start_window_size/17*10+start_window_size/8, start_window_size/4, "name of black player")

	var m_time sliders.Slider = sliders.New(start_window_size/20, start_window_size/25*10, start_window_size/17*10, start_window_size/20, start_window_size/40, 0, 59, 10, "game time in minutes", true, primary_color, secondary_color, bg_color)
	var s_time sliders.Slider = sliders.New(start_window_size/20, start_window_size/21*10, start_window_size/17*10, start_window_size/20, start_window_size/40, 0, 60, 0, "game time in seconds", true, primary_color, secondary_color, bg_color)
	var w_size sliders.Slider = sliders.New(start_window_size/20, start_window_size/18*10, start_window_size/17*10, start_window_size/20, start_window_size/40, 1, 150, 100, "squaresize in pixels", true, primary_color, secondary_color, bg_color)
	var p_time sliders.Slider = sliders.New(start_window_size/20, start_window_size/16*10, start_window_size/17*10, start_window_size/20, start_window_size/40, 0, 200, 1, "premove time in ms", true, primary_color, secondary_color, bg_color)
	m_time.Draw()
	s_time.Draw()
	w_size.Draw()
	p_time.Draw()

	var friendly_game buttons.Button = buttons.New(start_window_size/20, start_window_size/13*10, start_window_size/22*10, start_window_size/14, "friendly game", highlight_color1[0], highlight_color1[1], highlight_color1[2], secondary_color[0], secondary_color[1], secondary_color[2], start_window_size/100, int(start_window_size)/19)
	var use_premoves buttons.Button = buttons.New(start_window_size/18*10, start_window_size/13*10, start_window_size/25*10, start_window_size/14, "use premoves", highlight_color2[0], highlight_color2[1], highlight_color2[2], secondary_color[0], secondary_color[1], secondary_color[2], start_window_size/100, int(start_window_size)/19)
	var troll_mode buttons.Button = buttons.New(start_window_size, start_window_size/13*10, start_window_size/22*10, start_window_size/14, "extreme mode", highlight_color1[0], highlight_color1[1], highlight_color1[2], secondary_color[0], secondary_color[1], secondary_color[2], start_window_size/100, int(start_window_size)/17)
	var start buttons.Button = buttons.New(start_window_size/15*10, start_window_size/11*10, start_window_size/50*10, start_window_size/14, "Start", primary_color[0], primary_color[1], primary_color[2], secondary_color[0], secondary_color[1], secondary_color[2], start_window_size/50, int(start_window_size)/19)
	friendly_game.Draw()
	use_premoves.Draw()
	troll_mode.Draw()
	start.Draw()

	for {
		button, status, m_x, m_y := gfx.MausLesen1()

		if button == 1 && status == 1 {
			name_player_white.If_Clicked_Write(m_x, m_y)
			name_player_black.If_Clicked_Write(m_x, m_y)

			m_time.If_Clicked_Draw(m_x, m_y)
			s_time.If_Clicked_Draw(m_x, m_y)
			w_size.If_Clicked_Draw(m_x, m_y)
			p_time.If_Clicked_Draw(m_x, m_y)

			if friendly_game.Is_Clicked(m_x, m_y) {
				if friendly_game.Switch(highlight_color2[0], highlight_color2[1], highlight_color2[2]) {
					m_time.Deactivate()
					s_time.Deactivate()
				} else {
					m_time.Activate()
					s_time.Activate()
				}
			} else if use_premoves.Is_Clicked(m_x, m_y) {
				if use_premoves.Switch(highlight_color1[0], highlight_color1[1], highlight_color1[2]) {
					p_time.Deactivate()
				} else {
					p_time.Activate()
				}
			} else if troll_mode.Is_Clicked(m_x, m_y) {
				troll_mode.Switch(highlight_color2[0], highlight_color2[1], highlight_color2[2])
			} else if start.Is_Clicked(m_x, m_y) {
				var game_time int64 = int64(m_time.Get_Value())*60*1000 + int64(s_time.Get_Value())*1000
				var name_player_white_string string = name_player_white.Get_Text()
				var name_player_black_string string = name_player_black.Get_Text()

				if friendly_game.Give_State() {
					game_time = 0
				}

				if name_player_white_string == "" {
					name_player_white_string = "White"
				}
				if name_player_black_string == "" {
					name_player_black_string = "Black"
				}

				gfx.FensterAus()
				return game_time, friendly_game.Give_State(), int(p_time.Get_Value()), !use_premoves.Give_State(), name_player_white_string, name_player_black_string, uint16(w_size.Get_Value()), troll_mode.Give_State()
			}

		}
	}
}

func game_history_handler(moves_counter int16, moves_a [][64]pieces.Piece, pieces_a [64]pieces.Piece, game_history_can_be_changed bool, move_string string, pgn_moves_a []string) ([64]pieces.Piece, [][64]pieces.Piece, []string) {
	if int(moves_counter) == len(moves_a) { //apppend the moves_array when a new move was made
		moves_a = append_moves_array(moves_a, pieces_a)
		pgn_moves_a = append(pgn_moves_a, move_string)
	} else if !array_one_is_equal_to_array_two(pieces_a, moves_a[moves_counter]) { //if there was a move changed
		if game_history_can_be_changed { //if it's allowed to change a move in the history, this block replaces the move in the moves array with the current one and cuts of the rest of the array
			fmt.Println("there was a move changed in the game history")
			moves_a[moves_counter] = pieces.Copy_Array(pieces_a)
			moves_a = moves_a[:moves_counter+1]

			pgn_moves_a[moves_counter] = move_string
			pgn_moves_a = pgn_moves_a[:moves_counter+1]
		} else { //if it's not allowed it just deletes the move made and resets
			fmt.Println("changed game history although this is not allowed")
			pieces_a = pieces.Copy_Array(moves_a[moves_counter])
		}
	}
	return pieces_a, moves_a, pgn_moves_a
}

func time_handler(white_is_current_player bool, white_time_counter time_counter.Counter, black_time_counter time_counter.Counter, pause_button buttons.Button) {
	// if pause_button.Is_Active() {
	if pause_button.Give_State() {
		pause_button.Switch(0, 0, 0)
	}
	if white_is_current_player {
		white_time_counter.Stop_Counting()
		black_time_counter.Stop_Counting()
		white_time_counter.Init_Counting()

	} else if !white_is_current_player {
		black_time_counter.Stop_Counting()
		white_time_counter.Stop_Counting()
		black_time_counter.Init_Counting()
	}
	// }
}

func restart_handler(current_player_has_no_legal_moves bool, check bool, a uint16, white_is_current_player bool, gor_status chan bool) bool {
	if current_player_has_no_legal_moves { //game end / restart
		var review bool
		if check {
			review = display_message(0, a, white_is_current_player)
		} else {
			review = display_message(1, a, white_is_current_player)
		}
		if review {
			return false
		}
		gor_status <- true
		return true
	}
	return false
}

func mouse_handler(m_channel chan [4]int16, gor_status chan bool) {
	for {
		select {
		case quit := <-gor_status:
			if quit {
				return
			}
		default:
			button, status, m_x, m_y := gfx.MausLesen1()
			if !(button == 0 && status == 0) {
				select {
				case temp := <-m_channel: //stellt sicher dass leer ist
					if temp[0] == 1 && temp[1] == 1 { //wenn es ein klicken befehl ist, dann soll dieser nicht überschrieben werden, da es sonst zu bugs kommen kann und "verschluckt" wird, dass ein piece ausgewählt wurde
						m_channel <- temp
					} else {
						m_channel <- [4]int16{int16(button), int16(status), int16(m_x), int16(m_y)} //schreibt nur wenn leer ist
					}
				default:
					m_channel <- [4]int16{int16(button), int16(status), int16(m_x), int16(m_y)}
				}
			}
		}
	}
}

func get_clipboard_if_asked(use_clipboard bool) string {
	if use_clipboard {
		err := clipboard.Init()
		if err != nil {
			panic(err)
		}
		return string(clipboard.Read(clipboard.FmtText))
	} else {
		return ""
	}
}

func draw_timers(white_time_counter, black_time_counter time_counter.Counter, a uint16) bool {
	gfx.SetzeFont("./resources/fonts/firamono.ttf", int((a*10)/38))

	gfx.UpdateAus()
	gfx.Stiftfarbe(86, 82, 77)
	gfx.Vollrechteck(81*a/10, 57*a/10, 2*a-2*a/10, a-65*a/100)
	gfx.Stiftfarbe(48, 46, 43)
	gfx.Vollrechteck(89*a/10, 57*a/10, a-8*a/10, a-65*a/100)
	gfx.Stiftfarbe(255, 255, 255)
	white_string, white_has_no_time := white_time_counter.Return_Current_Counter()
	black_string, black_has_no_time := black_time_counter.Return_Current_Counter()
	black_time_counter.Return_Current_Counter()
	gfx.SchreibeFont(81*a/10, 57*a/10, white_string)
	gfx.SchreibeFont(91*a/10, 57*a/10, black_string)
	gfx.UpdateAn()

	if white_has_no_time || black_has_no_time {
		return true
	} else {
		return false
	}
}

func array_one_is_equal_to_array_two(array_a [64]pieces.Piece, array_b [64]pieces.Piece) bool {
	for i := 0; i < len(array_a); i++ {
		if fmt.Sprint(array_a[i]) != fmt.Sprint(array_b[i]) || fmt.Sprintf("%T", array_a[i]) != fmt.Sprintf("%T", array_b[i]) {
			return false
		}
	}
	return true
}

func get_current_piece(pieces_a [64]pieces.Piece, current_field [2]uint16) (pieces.Piece, int) {
	//gibt das
	for piece_index := 0; piece_index < len(pieces_a); piece_index++ {
		if pieces_a[piece_index] != nil && current_field == pieces_a[piece_index].Give_Pos() {
			return pieces_a[piece_index], piece_index
		}
	}
	return nil, 64
}

func move_if_current_field_is_in_legal_moves(current_field [2]uint16, pieces_a [64]pieces.Piece, piece_is_selected uint16, a, w_x, w_y uint16, current_king_index int, check bool, moves_counter int16, m_channel chan [4]int16) ([64]pieces.Piece, uint16, bool, string) {
	current_piece := pieces_a[piece_is_selected]
	current_legal_moves := current_piece.Give_Legal_Moves()
	for k := 0; k < len(current_legal_moves); k++ {
		if current_field == [2]uint16{current_legal_moves[k][0], current_legal_moves[k][1]} { //wenn das der Fall ist, wird das Piece bewegt
			//moved something
			var piece_promoting_to string = ""
			var original_field [2]uint16 = current_piece.Give_Pos()
			var take string = ""
			var promotion uint16

			pieces_a, promotion, take = pieces.Move_Piece_To(current_piece, current_legal_moves[k], moves_counter, pieces_a)

			if promotion != 64 {
				draw_board(a, w_x, w_y, current_piece, current_legal_moves, pieces_a, false, current_king_index, check)
				pieces_a, piece_promoting_to = pawn_promotion(w_x, w_y, a, int(piece_is_selected), pieces_a, "A", m_channel)
				piece_promoting_to = "=" + piece_promoting_to
			}
			move_string := get_move_string(int(piece_is_selected), original_field, piece_promoting_to, take, pieces_a)

			return pieces_a, 64, true, move_string
		}
	}
	return pieces_a, piece_is_selected, false, ""
}

func get_move_string(current_piece_index int, original_pos [2]uint16, piece_promoting string, take string, pieces_a [64]pieces.Piece) string {
	piece_name := pieces_a[current_piece_index].Give_Piece_Type()

	if piece_name == "K" && original_pos[0] == 4 && pieces_a[current_piece_index].Give_Pos()[0] == 6 { //short castle
		return "O-O"
	} else if piece_name == "K" && original_pos[0] == 4 && pieces_a[current_piece_index].Give_Pos()[0] == 2 { //long castle
		return "O-O-O"
	} else {

		var original_field string = ""
		var new_field string = parser.Get_Move_As_String_From_Field(pieces_a[current_piece_index].Give_Pos())

		if piece_promoting != "" {
			piece_name = ""
		}

		if take == "x" && piece_name == "" { //if a pawn takes someting, the original x cord of this pawn must be included
			original_field = parser.Translate_Field_Cord_To_PGN_String(original_pos[0], true)
		}

		//überprüfen ob es noch ein piece gibt welches auf das gleiche feld moven kann und vom selben typ ist
		if index_of_other_piece, _ := parser.Get_Piece_Index_And_Move_Index(pieces_a, pieces_a[current_piece_index].Give_Pos(), pieces_a[current_piece_index].Is_White_Piece(), piece_name, "0", current_piece_index); index_of_other_piece != 64 {
			//there is another piece that can execute the same move which means the pgn string is supposed to include more detailed information about the original position of the moved piece
			if pieces_a[index_of_other_piece].Give_Pos()[0] != original_pos[0] {
				//fmt.Println("the other piece has a different x cord --> move string will include the x cord")
				original_field = parser.Translate_Field_Cord_To_PGN_String(original_pos[0], true)
			} else if pieces_a[index_of_other_piece].Give_Pos()[1] != original_pos[1] {
				//fmt.Println("the other piece has a different y cord --> move string will include the y cord")
				original_field = parser.Translate_Field_Cord_To_PGN_String(original_pos[1], false)
			}
		}
		return piece_name + original_field + take + new_field + piece_promoting
	}
}

func display_message(message_type uint8, a uint16, white_is_current_player bool) bool {
	gfx.Archivieren()

	gfx.Stiftfarbe(0, 0, 0)
	gfx.Transparenz(70)
	gfx.Vollrechteck(0, 0, a*8, a*8)
	gfx.Transparenz(50)
	gfx.Stiftfarbe(8, 8, 8)
	gfx.Vollrechteck((a), 3*a, 6*a, 2*a)

	gfx.Stiftfarbe(220, 220, 220)
	gfx.SetzeFont("./resources/fonts/junegull.ttf", int(5*a/10))

	if message_type == 0 {
		if !white_is_current_player {
			gfx.SchreibeFont(18*a/10, 31*a/10, "player black has")
		} else {
			gfx.SchreibeFont(17*a/10, 31*a/10, "player white has")
		}
		gfx.Stiftfarbe(136, 8, 8)
		gfx.SetzeFont("./resources/fonts/punk.ttf", int(a))
		gfx.SchreibeFont(29*a/10, 37*a/10, "lost")
	} else if message_type == 1 {
		gfx.SchreibeFont(295*a/100, 31*a/10, "That's a ")

		gfx.Stiftfarbe(0, 143, 230)
		gfx.SetzeFont("./resources/fonts/punk.ttf", int(a))
		gfx.SchreibeFont(144*a/100, 37*a/10, "stalemate")

	} else if message_type == 2 {
		gfx.SchreibeFont(180*a/100, 31*a/10, "Press any key to")

		gfx.Stiftfarbe(28, 205, 60)
		gfx.SetzeFont("./resources/fonts/punk.ttf", int(a))
		gfx.SchreibeFont(260*a/100, 37*a/10, "star t")
	}

	gfx.Transparenz(0)
	gfx.UpdateAn()
	key, _, _ := gfx.TastaturLesen1()

	if key == 13 && message_type != 2 {
		return false
	} else if key != 13 && message_type != 2 {
		gfx.Restaurieren(0, 0, 8*a, 8*a)
		return true
	}
	return false
}

func pawn_promotion(w_x, w_y, a uint16, pawn_index int, pieces_a [64]pieces.Piece, premoved string, m_channel chan [4]int16) ([64]pieces.Piece, string) {
	var piece_promoted_to string
	var queen pieces.Piece = pieces.NewQueen(pieces_a[pawn_index].Give_Pos()[0], pieces_a[pawn_index].Give_Pos()[1], pieces_a[pawn_index].Is_White_Piece())
	var knight pieces.Piece = pieces.NewKnight(pieces_a[pawn_index].Give_Pos()[0], pieces_a[pawn_index].Give_Pos()[1], pieces_a[pawn_index].Is_White_Piece())
	var rook pieces.Piece = pieces.NewRook(pieces_a[pawn_index].Give_Pos()[0], pieces_a[pawn_index].Give_Pos()[1], pieces_a[pawn_index].Is_White_Piece())
	var bishop pieces.Piece = pieces.NewBishop(pieces_a[pawn_index].Give_Pos()[0], pieces_a[pawn_index].Give_Pos()[1], pieces_a[pawn_index].Is_White_Piece())

	if premoved == "A" {
		gfx.Stiftfarbe(0, 0, 0)
		gfx.Transparenz(70)
		gfx.Vollrechteck(0, 0, a*8, a*8)
		gfx.Transparenz(0)
		gfx.Stiftfarbe(221, 221, 221)
		gfx.Vollrechteck((2 * a), (4*a)-(5*a)/10, 4*a, a)

		pieces.Draw_To_Point(queen, w_x, w_y, a, (2 * a), (4*a)-(5*a)/10, 0, 0, 0, 0)
		pieces.Draw_To_Point(knight, w_x, w_y, a, (3 * a), (4*a)-(5*a)/10, 0, 0, 0, 0)
		pieces.Draw_To_Point(rook, w_x, w_y, a, (4 * a), (4*a)-(5*a)/10, 0, 0, 0, 0)
		pieces.Draw_To_Point(bishop, w_x, w_y, a, (5 * a), (4*a)-(5*a)/10, 0, 0, 0, 0)
		for {
			mouse_input := <-m_channel
			if mouse_input[0] == 1 && mouse_input[1] == 1 && uint16(mouse_input[2]) >= 2*a && uint16(mouse_input[2]) <= 6*a && uint16(mouse_input[3]) >= (4*a)-(5*a)/10 && uint16(mouse_input[3]) <= (4*a)+(5*a)/10 {
				x_field_pos := (calc_field(a, uint16(mouse_input[2]), uint16(mouse_input[3]), (5*a)/10))[0]
				if x_field_pos == 2 {
					piece_promoted_to = "Q"
					pieces_a[pawn_index] = queen
				} else if x_field_pos == 3 {
					piece_promoted_to = "N"
					pieces_a[pawn_index] = knight
				} else if x_field_pos == 4 {
					piece_promoted_to = "R"
					pieces_a[pawn_index] = rook
				} else if x_field_pos == 5 {
					piece_promoted_to = "B"
					pieces_a[pawn_index] = bishop
				}
				break
			}
		}

	} else {
		piece_promoted_to = premoved

		if premoved == "Q" {
			pieces_a[pawn_index] = queen
		} else if premoved == "N" {
			pieces_a[pawn_index] = knight
		} else if premoved == "R" {
			pieces_a[pawn_index] = rook
		} else if premoved == "B" {
			pieces_a[pawn_index] = bishop
		}
	}
	return pieces_a, piece_promoted_to
}

func draw_board(a, w_x, w_y uint16, current_piece pieces.Piece, current_legal_moves [][3]uint16, pieces_a [64]pieces.Piece, highlighting_is_activated bool, current_king_index int, check bool) {
	gfx.UpdateAus()
	draw_background(a)
	if check {
		highlight(a, pieces_a[current_king_index].Give_Pos(), 255, 0, 0)
	}
	if highlighting_is_activated {
		highlight(a, current_piece.Give_Pos(), 0, 50, 255)
		for k := 0; k < len(current_legal_moves); k++ {
			highlight(a, [2]uint16{current_legal_moves[k][0], current_legal_moves[k][1]}, 0, 255, 0)
		}
	}
	draw_pieces(pieces_a, w_x, w_y, a)
	gfx.UpdateAn()
}

func change_player(white_is_current_player bool, white_king_index, black_king_index int, moves_counter int16) (bool, int, bool, int16) {
	var current_king_index int
	moves_counter = moves_counter + 1
	if white_is_current_player {
		current_king_index = black_king_index
		white_is_current_player = false
	} else {
		current_king_index = white_king_index
		white_is_current_player = true
	}
	return white_is_current_player, current_king_index, false, moves_counter
}

func draw_moves_sidebar(a uint16, moves_counter int16, pgn_moves_a []string) {
	gfx.UpdateAus()
	const display_limit int16 = 18
	const lower_bound uint16 = 52

	gfx.Stiftfarbe(124, 119, 111)
	gfx.Vollrechteck(8*a+a/10, 8*a/10, 9*a/10, (lower_bound-4)*a/10)
	gfx.Stiftfarbe(31, 32, 33)
	gfx.Vollrechteck(9*a, 8*a/10, 9*a/10, (lower_bound-4)*a/10)

	for i := moves_counter; moves_counter-i <= display_limit && i != 0; i-- {
		if i%2 != 0 { //white's move
			var move_number string = strconv.Itoa(int(i+1) / 2)

			//display move number
			gfx.Stiftfarbe(48, 46, 43)
			gfx.Vollrechteck(81*a/10, (lower_bound-2)*a/10-5*a/10*uint16((moves_counter-i)/2), 18*a/10, 2*a/10)
			gfx.SetzeFont("./resources/fonts/firamono.ttf", int(a/5))
			gfx.Stiftfarbe(200, 200, 200)
			if len(move_number) == 1 {
				gfx.SchreibeFont(89*a/10, (lower_bound-2)*a/10-5*a/10*uint16((moves_counter-i)/2), move_number)
			} else if len(move_number) == 2 {
				gfx.SchreibeFont(89*a/10, (lower_bound-2)*a/10-5*a/10*uint16((moves_counter-i)/2), move_number)
			} else if len(move_number) == 3 {
				gfx.SchreibeFont(88*a/10, (lower_bound-2)*a/10-5*a/10*uint16((moves_counter-i)/2), move_number)
			}

			//display move
			gfx.SetzeFont("./resources/fonts/firamono.ttf", int(a/4))
			gfx.Stiftfarbe(31, 32, 33)
			gfx.SchreibeFont(81*a/10, lower_bound*a/10-5*a/10*uint16((moves_counter-i)/2), pgn_moves_a[i])
		} else { //black's move
			if moves_counter-i >= display_limit-1 { //break if the upper bound has been reached
				break
			} else {
				gfx.SetzeFont("./resources/fonts/firamono.ttf", int(a/4))
				gfx.Stiftfarbe(124, 119, 111)
				gfx.SchreibeFont(90*a/10, lower_bound*a/10-5*a/10*uint16((moves_counter-i+1)/2), pgn_moves_a[i])
			}
		}
	}
	gfx.UpdateAn()
}

func draw_player_names(name_player_white, name_player_black string, a uint16) {
	gfx.Stiftfarbe(200, 191, 179)
	gfx.Vollrechteck(81*a/10, a/10, 18*a/10, 4*a/10)
	gfx.Vollrechteck(80*a/10, 6*a/10, 20*a/10, 1*a/10)

	gfx.Stiftfarbe(48, 46, 43)
	gfx.Vollrechteck(89*a/10, a/10, 2*a/10, 5*a/10)

	gfx.Stiftfarbe(48, 46, 43)

	var max_name_lenght int = 1

	if len(name_player_white) >= len(name_player_black) {
		max_name_lenght = len(name_player_white)
	} else {
		max_name_lenght = len(name_player_black)
	}

	gfx.SetzeFont("./resources/fonts/firamono.ttf", int(a)/max_name_lenght)

	gfx.SchreibeFont(82*a/10, 2*a/10, name_player_white)
	gfx.SchreibeFont(92*a/10, 2*a/10, name_player_black)
}

func initialize(w_x, w_y, a uint16, restart bool, game_timer int64, name_player_white string, name_player_black string, image_location string) ([64]pieces.Piece, int, int, buttons.Button, buttons.Button, buttons.Button, buttons.Button, buttons.Button, [][64]pieces.Piece, time_counter.Counter, time_counter.Counter, []string) {
	var moves_a [][64]pieces.Piece
	var pgn_moves_a []string

	var one_move_back buttons.Button = buttons.New(81*a/10, 7*a+a/10, 8*a/10, a-a/5, "<", 38, 37, 34, 200, 200, 200, (a / 4), int(a/2))
	var one_move_forward buttons.Button = buttons.New(91*a/10, 7*a+a/10, 8*a/10, a-a/5, ">", 38, 37, 34, 200, 200, 200, (a / 4), int(a/2))
	var restart_button buttons.Button = buttons.New(8*a+a/7, 62*a/10, 17*a/10, 3*a/10, "restart game", 86, 82, 77, 200, 200, 200, (a / 10), int(a/5))
	var pause_button buttons.Button = buttons.New(91*a/10, 66*a/10, 8*a/10, 3*a/10, "pause", 86, 82, 77, 200, 200, 200, (a / 15), int(a/5))
	var save_button buttons.Button = buttons.New(81*a/10, 66*a/10, 8*a/10, 3*a/10, "save", 86, 82, 77, 200, 200, 200, (a / 7), int(a/5))

	if !restart {
		gfx.Fenster(w_x, w_y)
		gfx.Fenstertitel("Chess")
		rescale_image(a, image_location)
	}

	if game_timer == 0 {
		pause_button.Deactivate()
	}

	gfx.Stiftfarbe(86, 82, 77)
	gfx.Vollrechteck(8*a, 7*a, 2*a, a)
	gfx.Stiftfarbe(48, 46, 43)
	gfx.Vollrechteck(8*a, 6*a, 2*a, a)

	one_move_back.Draw()
	one_move_forward.Draw()
	restart_button.Draw()
	pause_button.Draw()
	save_button.Draw()

	gfx.Stiftfarbe(48, 46, 43)
	gfx.Vollrechteck(8*a, 0, 2*a, 6*a)

	//draw_moves_sidebar(a)
	draw_player_names(name_player_white, name_player_black, a)
	draw_background(a)

	var pieces_a [64]pieces.Piece
	var white_king_index int = -1
	var black_king_index int = -1

	pieces_a[0] = pieces.NewRook(0, 0, false)
	pieces_a[1] = pieces.NewKnight(1, 0, false)
	pieces_a[2] = pieces.NewBishop(2, 0, false)
	pieces_a[3] = pieces.NewQueen(3, 0, false)
	pieces_a[4] = pieces.NewKing(4, 0, false)
	pieces_a[5] = pieces.NewBishop(5, 0, false)
	pieces_a[6] = pieces.NewKnight(6, 0, false)
	pieces_a[7] = pieces.NewRook(7, 0, false)

	var i uint16
	for i = 0; i < 8; i++ {
		pieces_a[i+8] = pieces.NewPawn(i, 1, false)
	}
	for i = 0; i < 8; i++ {
		pieces_a[i+16] = pieces.NewPawn(i, 6, true)
	}

	pieces_a[24] = pieces.NewRook(0, 7, true)
	pieces_a[25] = pieces.NewKnight(1, 7, true)
	pieces_a[26] = pieces.NewBishop(2, 7, true)
	pieces_a[27] = pieces.NewQueen(3, 7, true)
	pieces_a[28] = pieces.NewKing(4, 7, true)
	pieces_a[29] = pieces.NewBishop(5, 7, true)
	pieces_a[30] = pieces.NewKnight(6, 7, true)
	pieces_a[31] = pieces.NewRook(7, 7, true)

	//moves_a = append_moves_array(moves_a, pieces_a)

	for i := 0; i < len(pieces_a); i++ {
		if pieces_a[i] != nil {
			if king, ok := pieces_a[i].(*pieces.King); ok {
				if king.Is_White_Piece() {
					white_king_index = i
				} else if !king.Is_White_Piece() {
					black_king_index = i
				}
			}
		}
	}

	white_time_counter := time_counter.New(game_timer)
	black_time_counter := time_counter.New(game_timer)

	return pieces_a, white_king_index, black_king_index, one_move_back, one_move_forward, restart_button, pause_button, save_button, moves_a, white_time_counter, black_time_counter, pgn_moves_a
}

func append_moves_array(moves_a [][64]pieces.Piece, pieces_a [64]pieces.Piece) [][64]pieces.Piece {
	moves_a = append(moves_a, pieces.Copy_Array(pieces_a))
	return moves_a
}

func draw_pieces(pieces_a [64]pieces.Piece, w_x, w_y, a uint16) {
	gfx.UpdateAus()
	for i := 0; i < len(pieces_a); i++ {
		if pieces_a[i] != nil {
			gfx.Archivieren()
			pieces.Draw(pieces_a[i], w_x, w_y, a)
			gfx.Restaurieren(0, 0, w_x, w_y)

			gfx.Clipboard_einfuegenMitColorKey(pieces_a[i].Give_Pos()[0]*a, pieces_a[i].Give_Pos()[1]*a, 5, 5, 5)
		}
	}
	gfx.Archivieren()
	gfx.UpdateAn()
}

func highlight(a uint16, pos [2]uint16, r, g, b uint8) {

	var cord_x uint16 = a * pos[0]
	var cord_y uint16 = a * pos[1]

	gfx.Transparenz(170)
	gfx.Stiftfarbe(r, g, b)
	gfx.Vollrechteck(cord_x, cord_y, a, a)
	gfx.Transparenz(0)
}

func draw_background(a uint16) {
	var f_x uint16 = 0
	var f_y uint16 = 0
	for i := 0; i <= 7; i++ {
		for k := 0; k <= 7; k++ {
			if k%2 == 0 {
				if i%2 == 0 {
					gfx.Stiftfarbe(240, 217, 181)
				} else {
					gfx.Stiftfarbe(181, 136, 99)
				}

			} else {
				if i%2 == 1 {
					gfx.Stiftfarbe(240, 217, 181)
				} else {
					gfx.Stiftfarbe(181, 136, 99)
				}
			}

			gfx.Vollrechteck(f_x, f_y, a, a)
			f_x = f_x + a
		}
		f_x = 0
		f_y = f_y + a
	}
}

func rescale_image(a uint16, image_location string) {

	// Open the BMP file
	file, err := os.Open(image_location)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Decode the BMP file
	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("Error decoding file:", err)
		return
	}

	// Specify the desired width and height
	var width int = 6 * int(a)
	var height int = 2 * int(a)

	// Resize the image to the specified dimensions
	resizedImg := imaging.Resize(img, width, height, imaging.NearestNeighbor)

	// Save the resized image to a new BMP file
	err = imaging.Save(resizedImg, "./resources/images/Pieces.bmp")
	if err != nil {
		fmt.Println("Error saving file:", err)
		return
	}
}

func calc_field(a, m_x, m_y, y_offset uint16) [2]uint16 {
	var current_field [2]uint16

	current_field[0] = m_x / a
	current_field[1] = (m_y + y_offset) / a
	return current_field
}

func set_image_string(troll_mode bool) string {
	if troll_mode {
		return "./resources/images/Troll.bmp"
	} else {
		return "./resources/images/Pieces_Source_Original.bmp"
	}
}
