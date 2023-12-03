package main //added some commentary and fixed a bug where the user could make mouse inputs during the premove phase which led to unexpected behaviour

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

	var deselect_piece_after_clicking = false //decides whether a selected piece should be deselected after executing a second click directly on it
	var start_window_width uint16 = 600       //the start window is technically scalable but with resolutions higher than 800 graphical bugs are occuring due to the font size implementation of gfx
	game_timer, friendly_game, duration_of_premove_animation, use_clipboard_as_premoves, name_player_white, name_player_black, a, troll_mode := start_menu(start_window_width / 3 * 2)

restart_marker: //jump marker for the restart
	var image_location string = set_image_string(troll_mode) //troll mode decides which picture should be used for the pieces, if activated a picture of a cloud will be used
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
	var review_the_game bool = false
	var end_of_game bool = false
	m_channel := make(chan [4]int16, 1) //channel gets the latest mouse input, structured as described in the following: button, status, x_cord, y_cord
	gor_status := make(chan bool, 1)    //needs to be a buffered channel, indicated by the one, otherwise the program will hold until the channel is empty after puting something to it

	// execute the initialize to define important vars and create the window
	pieces_a, white_king_index, black_king_index, one_move_back, one_move_forward, restart_button, pause_button, save_button, moves_a, white_time_counter, black_time_counter, pgn_moves_a := initialize(w_x, w_y, a, restart_window, game_timer, name_player_white, name_player_black, image_location)
	premoves_array := parser.Create_Array_Of_Moves(premoves) //get the premoves array

	display_message(2, a, false) //starting message

	go mouse_handler(m_channel, gor_status) //activate the mouse handler so that mouse input can be obtained

	for { //gameloop

		if player_change { //after a player change do the following:
			pieces_a, moves_a, pgn_moves_a = game_history_handler(moves_counter, moves_a, pieces_a, game_history_can_be_changed, move_string, pgn_moves_a)                        //add the current move to pgn moves array as well as to the normal moves array (important for back and forward)
			white_is_current_player, current_king_index, player_change, moves_counter = change_player(white_is_current_player, white_king_index, black_king_index, moves_counter) //this function sets the current_king_index
			pieces_a, current_player_has_no_legal_moves = pieces.Calc_Moves_With_Check(pieces_a, moves_counter, current_king_index)                                               //calc the move
			check = pieces_a[current_king_index].(*pieces.King).Is_In_Check(pieces_a, moves_counter)                                                                              //check if the current king is in check

			if !use_clipboard_as_premoves || duration_of_premove_animation != 0 { //wenn keine premoves mehr da sind oder premoves da sind und diese gezeichnet werden sollen dann wird das board gezeichnet
				draw_board(a, w_x, w_y, current_piece, current_legal_moves, pieces_a, false, current_king_index, check)
				draw_moves_sidebar(a, moves_counter-1, pgn_moves_a)
			}
			if !end_of_game {
				if restart_window = restart_handler(current_player_has_no_legal_moves, check, a, white_is_current_player, gor_status, false); restart_window { //restart if the premove handler calculated the end of the game
					goto restart_marker
				} else if current_player_has_no_legal_moves {
					review_the_game = true
					end_of_game = true
				}
			}
			if review_the_game && !friendly_game && int(moves_counter) == len(moves_a) && !end_of_game {
				review_the_game = false //reset review game if the player is reviewing the most recent move
				if pause_button.Give_State() {
					pause_button.Switch(38, 37, 34)
				}
			} else if review_the_game {
				review_mask(a) //creates a review mask, to visulize that the game is in review state
			}
			if !use_clipboard_as_premoves && !friendly_game { //only start the timers if we're not premoving and it is a competetive game
				time_handler(white_is_current_player, white_time_counter, black_time_counter, pause_button, review_the_game)
			}
			gfx.UpdateAn()
		}

		if len(premoves_array) > 0 { //if there are premoves the program premoves
			var piece_promoted_to string = ""
			piece_executing_move, index_of_move, piece_promoting_to := parser.Get_Correct_Move(premoves_array[0], pieces_a, current_king_index) //get the move (so that the enigine can handle it)

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
			} else { //immmediately end the premove phase if there is one error in the premove sequence
				premoves_array = nil
			}
		} else if ending_premoves { //should only be executed once after the premove phase
			draw_board(a, w_x, w_y, current_piece, current_legal_moves, pieces_a, false, current_king_index, check)
			use_clipboard_as_premoves = false
			ending_premoves = false
			draw_moves_sidebar(a, moves_counter-1, pgn_moves_a)
			gfx.UpdateAn()
			clear_m_channel(m_channel) //this is necessary cuz otherwise pieces would be selected unexpextetly after the premove phase if the user made mlus inputs during this phase
		}

		if !use_clipboard_as_premoves { //wenn die andere restart methode ausgeführt wird, dann sollt hier noch noch überprüft werden ob restart == false ist, da es sonst zu fehlern kommt
			select {
			case mouse_input := <-m_channel: //get the mouse input if there is any from the mouse channel

				if dragging && !(mouse_input[0] == 1 && mouse_input[1] == 0) { //wenn nichts gehalten wird dann löschen
					draw_board(a, w_x, w_y, current_piece, current_legal_moves, pieces_a, true, current_king_index, check)
					gfx.UpdateAn()
					dragging = false
				}

				if mouse_input[0] == 1 && mouse_input[1] == 1 { //if left clicked was pressed
					//check the buttons
					if one_move_back.Is_Clicked(uint16(mouse_input[2]), uint16(mouse_input[3])) {
						if moves_counter >= 2 {
							moves_counter = moves_counter - 2
							pieces_a = pieces.Copy_Array(moves_a[moves_counter]) //set the pieces array to the moves counter reduced by 2 (cuz it will be increased by one in the next cycle so that the button is actually going back one move and not two)
							player_change = true
							if !friendly_game {
								review_the_game = true
							}
						}
					} else if one_move_forward.Is_Clicked(uint16(mouse_input[2]), uint16(mouse_input[3])) {
						if int(moves_counter) < len(moves_a) {
							pieces_a = pieces.Copy_Array(moves_a[moves_counter]) //set the pieces array to the moves couhnter, which is at this point exactly one higher than the index of the current pieces array
							player_change = true
						}
					} else if restart_button.Is_Clicked(uint16(mouse_input[2]), uint16(mouse_input[3])) {
						restart_window = true //prevents the initialize function from creating a new window
						gor_status <- true    //kills the mouse_handler which is running in concurrency
						goto restart_marker
					} else if save_button.Is_Clicked(uint16(mouse_input[2]), uint16(mouse_input[3])) {
						parser.Write_PGN_File(pgn_moves_a, name_player_white, name_player_black)
					} else if int(moves_counter) == len(moves_a) && pause_button.Is_Clicked(uint16(mouse_input[2]), uint16(mouse_input[3])) { //the pause button will only work if the user is reviewing the current move
						if pause_button.Switch(38, 37, 34) { //if pause was pressed
							white_time_counter.Stop_Counting()
							black_time_counter.Stop_Counting()
						} else if white_is_current_player { //if pause was released and white is current player
							white_time_counter.Init_Counting()
						} else if !white_is_current_player { //if pause was released and white is current player
							black_time_counter.Init_Counting()
						}
					} else if !review_the_game { //otherwise the program checks if a piece was clicked
						current_field = calc_field(a, uint16(mouse_input[2]), uint16(mouse_input[3]), 0)
						temp_current_piece, piece_index = get_current_piece(pieces_a, current_field) //get the piece

						//if a piece was selected before that is in the correct color and deselect_piece_after_clicking is off or no piece was selected before or there was a piece selected but it was not the same
						if temp_current_piece != nil && temp_current_piece.Is_White_Piece() == white_is_current_player && (!deselect_piece_after_clicking || (current_piece == nil || temp_current_piece.Give_Pos() != current_piece.Give_Pos())) {
							//select the piece
							current_piece = temp_current_piece
							current_legal_moves = current_piece.Give_Legal_Moves()
							promotion = 0
							draw_board(a, w_x, w_y, current_piece, current_legal_moves, pieces_a, true, current_king_index, check)
							gfx.UpdateAn()
							piece_is_selected = uint16(piece_index)
						} else if piece_is_selected != 64 { //if there is a piece selected already
							//überprüfen ob auf ein Feld in Legal Moves geklickt wurde
							pieces_a, piece_is_selected, player_change, move_string = move_if_current_field_is_in_legal_moves(current_field, pieces_a, piece_is_selected, a, w_x, w_y, current_king_index, check, moves_counter, m_channel)

							//Deselect sobald ein Piece ausgewählt ist: wenn auf kein Piece geklickt wurde, oder auf ein gegnerisches Piece geklickt wurde, wird das ausgewählte Piece deselected
							//mit der Option deselect_piece_after_clicking kann eingestellt werde, ob auch deselected werden soll wenn auf dasslebe piece nocheinmal geklickt wurde
							if (temp_current_piece == nil) || (temp_current_piece != nil && (temp_current_piece.Give_Pos() == current_piece.Give_Pos() || temp_current_piece.Is_White_Piece() != white_is_current_player)) {
								current_piece = nil
								draw_board(a, w_x, w_y, current_piece, current_legal_moves, pieces_a, false, current_king_index, check)
								gfx.UpdateAn()
								piece_is_selected = 64
							}
						}
					}

				} else if mouse_input[1] == -1 && mouse_input[0] == 1 && !review_the_game { //if left click was released
					current_field = calc_field(a, uint16(mouse_input[2]), uint16(mouse_input[3]), 0)
					temp_current_piece, _ = get_current_piece(pieces_a, current_field)

					if piece_is_selected != 64 { //if there is a piece selected
						//überprüfen ob in current legal moves
						pieces_a, piece_is_selected, player_change, move_string = move_if_current_field_is_in_legal_moves(current_field, pieces_a, piece_is_selected, a, w_x, w_y, current_king_index, check, moves_counter, m_channel)

						//Deselect sobald ein Piece ausgewählt ist: wenn auf keinem Piece losgelassen wurde oder auf ein generisches Piece losgelassen wurde oder wenn auf einem eigenen piece losgelassen wurde
						if (temp_current_piece == nil) || (temp_current_piece != nil && ((temp_current_piece.Is_White_Piece() != white_is_current_player) || (temp_current_piece.Is_White_Piece() == white_is_current_player && temp_current_piece.Give_Pos() != current_piece.Give_Pos()))) {
							current_piece = nil
							draw_board(a, w_x, w_y, current_piece, current_legal_moves, pieces_a, false, current_king_index, check)
							gfx.UpdateAn()
							piece_is_selected = 64
						}
					}
				} else if mouse_input[1] == 0 && mouse_input[0] == 1 && current_piece != nil && current_piece.Is_White_Piece() == white_is_current_player && !review_the_game { //if the left button is pressed and there is a current piece in the correct color
					//wenn die maustaste gehalten wird, wird ein ghosttpiece gemalt, welches der maus folgt
					dragging = true
					pieces.Draw_To_Point(current_piece, w_x, w_y, a, uint16(mouse_input[2]), uint16(mouse_input[3]), -int16(a/2), -int16(a/2), 50, uint16(mouse_input[2]))
				}

			default: //prevents lag or high cpu ussage
				time.Sleep(5 * time.Millisecond)
			}
			if !friendly_game && draw_timers(white_time_counter, black_time_counter, a) && !end_of_game { //restart if one player has no time left and this message wasn't displayed before indicated by the end_of_game statement
				if restart_window = restart_handler(current_player_has_no_legal_moves, check, a, white_is_current_player, gor_status, true); restart_window { //restart if the premove handler calculated the end of the game
					goto restart_marker
				} else { //if the player decides to review by clicking on any other button than enter or esc, the game is put into the review state
					review_the_game = true
					end_of_game = true
					review_mask(a)
				}
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

	//background and title
	gfx.Stiftfarbe(bg_color[0], bg_color[1], bg_color[2])
	gfx.Vollrechteck(0, 0, start_window_size/2*3, start_window_size)
	gfx.Stiftfarbe(secondary_color[0], secondary_color[1], secondary_color[2])
	gfx.SetzeFont("./resources/fonts/punk.ttf", int(start_window_size/12))
	gfx.SchreibeFont(start_window_size/2*3/2/12*10, start_window_size/90, "CHESS")

	//initialize the textboxes
	var name_player_white textbox.Box = textbox.New(start_window_size/20, start_window_size/6, start_window_size/14, start_window_size/17*10+start_window_size/40, primary_color, secondary_color, highlight_color1, 19, "enter white name...")
	var name_player_black textbox.Box = textbox.New(start_window_size/20, start_window_size/4, start_window_size/14, start_window_size/17*10+start_window_size/40, primary_color, secondary_color, highlight_color1, 19, "enter black name...")
	name_player_white.Draw()
	name_player_black.Draw()
	gfx.SetzeFont("./resources/fonts/unispace.ttf", int(start_window_size/20))
	gfx.Stiftfarbe(secondary_color[0], secondary_color[1], secondary_color[2])
	gfx.SchreibeFont(start_window_size/17*10+start_window_size/8, start_window_size/6, "name of white player")
	gfx.SchreibeFont(start_window_size/17*10+start_window_size/8, start_window_size/4, "name of black player")

	//create the sliders
	var m_time sliders.Slider = sliders.New(start_window_size/20, start_window_size/25*10, start_window_size/17*10, start_window_size/20, start_window_size/40, 0, 59, 10, "game time in minutes", true, primary_color, secondary_color, bg_color)
	var s_time sliders.Slider = sliders.New(start_window_size/20, start_window_size/21*10, start_window_size/17*10, start_window_size/20, start_window_size/40, 0, 60, 0, "game time in seconds", true, primary_color, secondary_color, bg_color)
	var w_size sliders.Slider = sliders.New(start_window_size/20, start_window_size/18*10, start_window_size/17*10, start_window_size/20, start_window_size/40, 1, 150, 100, "squaresize in pixels", true, primary_color, secondary_color, bg_color)
	var p_time sliders.Slider = sliders.New(start_window_size/20, start_window_size/16*10, start_window_size/17*10, start_window_size/20, start_window_size/40, 0, 200, 1, "premove time in ms", true, primary_color, secondary_color, bg_color)
	m_time.Draw()
	s_time.Draw()
	w_size.Draw()
	p_time.Draw()

	//create the buttons
	var friendly_game buttons.Button = buttons.New(start_window_size/20, start_window_size/13*10, start_window_size/22*10, start_window_size/14, "friendly game", highlight_color1[0], highlight_color1[1], highlight_color1[2], secondary_color[0], secondary_color[1], secondary_color[2], start_window_size/100, int(start_window_size)/19)
	var use_premoves buttons.Button = buttons.New(start_window_size/18*10, start_window_size/13*10, start_window_size/25*10, start_window_size/14, "use premoves", highlight_color2[0], highlight_color2[1], highlight_color2[2], secondary_color[0], secondary_color[1], secondary_color[2], start_window_size/100, int(start_window_size)/19)
	var troll_mode buttons.Button = buttons.New(start_window_size, start_window_size/13*10, start_window_size/22*10, start_window_size/14, "extreme mode", highlight_color1[0], highlight_color1[1], highlight_color1[2], secondary_color[0], secondary_color[1], secondary_color[2], start_window_size/100, int(start_window_size)/17)
	var start buttons.Button = buttons.New(start_window_size/15*10, start_window_size/11*10, start_window_size/50*10, start_window_size/14, "Start", primary_color[0], primary_color[1], primary_color[2], secondary_color[0], secondary_color[1], secondary_color[2], start_window_size/50, int(start_window_size)/19)
	friendly_game.Draw()
	use_premoves.Draw()
	troll_mode.Draw()
	start.Draw()

	//start menu cycle
	for {
		button, status, m_x, m_y := gfx.MausLesen1()

		if button == 1 && status == 1 {
			//go into the white cycles when the textboxes are clicked
			name_player_white.If_Clicked_Write(m_x, m_y)
			name_player_black.If_Clicked_Write(m_x, m_y)

			//redraw the sliders if they were clicked
			m_time.If_Clicked_Draw(m_x, m_y)
			s_time.If_Clicked_Draw(m_x, m_y)
			w_size.If_Clicked_Draw(m_x, m_y)
			p_time.If_Clicked_Draw(m_x, m_y)

			if friendly_game.Is_Clicked(m_x, m_y) {
				if friendly_game.Switch(highlight_color2[0], highlight_color2[1], highlight_color2[2]) { //if friendly game is active, timers are not needed --> temporarily remove the time sliders
					m_time.Deactivate()
					s_time.Deactivate()
				} else {
					m_time.Activate()
					s_time.Activate()
				}
			} else if use_premoves.Is_Clicked(m_x, m_y) {
				if use_premoves.Switch(highlight_color1[0], highlight_color1[1], highlight_color1[2]) { //if there are no premoves, the slider for the premove time is not needed --> temporarily remove the premove time slider
					p_time.Deactivate()
				} else {
					p_time.Activate()
				}
			} else if troll_mode.Is_Clicked(m_x, m_y) {
				troll_mode.Switch(highlight_color2[0], highlight_color2[1], highlight_color2[2])

			} else if start.Is_Clicked(m_x, m_y) { //extract the parameters from the user control units
				var game_time int64 = int64(m_time.Get_Value())*60*1000 + int64(s_time.Get_Value())*1000
				var name_player_white_string string = name_player_white.Get_Text()
				var name_player_black_string string = name_player_black.Get_Text()

				if friendly_game.Give_State() {
					game_time = 0
				}
				//default names for the players if no names are specified
				if name_player_white_string == "" {
					name_player_white_string = "White"
				}
				if name_player_black_string == "" {
					name_player_black_string = "Black"
				}

				gfx.FensterAus() //close the start menue and return the received parameters to main
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
			//there was a move changed in the game history
			moves_a[moves_counter] = pieces.Copy_Array(pieces_a)
			moves_a = moves_a[:moves_counter+1]

			pgn_moves_a[moves_counter] = move_string
			pgn_moves_a = pgn_moves_a[:moves_counter+1]
		} else { //if it's not allowed it just deletes the move made and resets
			//changed game history although this is not allowed
			pieces_a = pieces.Copy_Array(moves_a[moves_counter])
		}
	}
	return pieces_a, moves_a, pgn_moves_a
}

func time_handler(white_is_current_player bool, white_time_counter time_counter.Counter, black_time_counter time_counter.Counter, pause_button buttons.Button, review_mode bool) {
	if !review_mode { //when not in review mode a switch of the timers is needed
		if pause_button.Give_State() { //additionally the case that the pause button is still pressed should be taken care of --> if so, unpause should be initiated
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
	} else { //pause the timers when in review mode and indictae that via the pause button
		if !pause_button.Give_State() {
			pause_button.Switch(38, 37, 34)
			white_time_counter.Stop_Counting()
			black_time_counter.Stop_Counting()
		}
	}
}

func restart_handler(current_player_has_no_legal_moves bool, check bool, a uint16, white_is_current_player bool, gor_status chan bool, time_is_up bool) bool {
	var review bool

	if current_player_has_no_legal_moves { //game is over cuz there are no legal moves
		if check {
			review = display_message(0, a, white_is_current_player) //checkmate
		} else {
			review = display_message(1, a, white_is_current_player) //stalemate
		}
		if review { //if the game is supposed to be reviewed return false
			return false
		}
		gor_status <- true //otherwise kill the mouse handler via this channel and return true
		return true
	} else if time_is_up { //if the time is up the game is over anyways
		review = display_message(0, a, white_is_current_player) //display losing message
		if review {                                             //if the game is supposed to be reviewed return false
			return false
		}
		gor_status <- true //otherwise kill the mouse handler via this channel and return true
		return true
	}
	return false //otherwise return false
}

func mouse_handler(m_channel chan [4]int16, gor_status chan bool) {
	for {
		select {
		case quit := <-gor_status: //kill this thread if there is a true in this channel
			if quit {
				return //kill
			}
		default:
			button, status, m_x, m_y := gfx.MausLesen1()
			if !(button == 0 && status == 0) {
				select {
				case temp := <-m_channel: //stellt sicher dass leer ist
					if temp[0] == 1 && temp[1] == 1 { //wenn es ein klicken befehl ist, dann soll dieser nicht überschrieben werden, da es sonst zu bugs kommen kann und "verschluckt" wird, dass ein piece ausgewählt wurde
						m_channel <- temp
					} else { //otherwise just overwrite
						m_channel <- [4]int16{int16(button), int16(status), int16(m_x), int16(m_y)} //schreibt nur wenn leer ist
					}
				default: //if the channel is empty anyway, just push in the mouse information
					m_channel <- [4]int16{int16(button), int16(status), int16(m_x), int16(m_y)}
				}
			}
		}
	}
}

func review_mask(a uint16) { //creates a transparent mask over the board to inciate that nothing can be moved and the game is in the reviewing state
	gfx.Stiftfarbe(0, 0, 0)
	gfx.Transparenz(120)
	gfx.Vollrechteck(0, 0, a*8, a*8)
	gfx.Transparenz(0)
}

func clear_m_channel(m_channel chan [4]int16) {
	select {
	case <-m_channel:
		//there was something in the channel --> cleared
	default:
		//channel empty anyway
	}
}

func get_clipboard_if_asked(use_clipboard bool) string { //get the latest entry of the clipboard
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
	//draw timers
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

	// return true if any of the players ran out of time
	if white_has_no_time || black_has_no_time {
		return true
	} else {
		return false
	}
}

func array_one_is_equal_to_array_two(array_a [64]pieces.Piece, array_b [64]pieces.Piece) bool {
	for i := 0; i < len(array_a); i++ {
		if fmt.Sprint(array_a[i]) != fmt.Sprint(array_b[i]) || fmt.Sprintf("%T", array_a[i]) != fmt.Sprintf("%T", array_b[i]) {
			return false //return false if either the cordinates of the coresponding pieces are different or the pieces are of different type in general (this is important if a pawn promoted to the same field but to a different piece type)
		}
	}
	return true
}

func get_current_piece(pieces_a [64]pieces.Piece, current_field [2]uint16) (pieces.Piece, int) {
	//returns a piece and the corresponding index in the pieces array that matches the given field
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
	for k := 0; k < len(current_legal_moves); k++ { //iterates through the legal moves of the given piece
		if current_field == [2]uint16{current_legal_moves[k][0], current_legal_moves[k][1]} { //if the correct legal move is found, execute it
			//moved something
			var piece_promoting_to string = ""
			var original_field [2]uint16 = current_piece.Give_Pos()
			var take string = ""
			var promotion uint16

			pieces_a, promotion, take = pieces.Move_Piece_To(current_piece, current_legal_moves[k], moves_counter, pieces_a)

			if promotion != 64 { //execute the promotion if needed
				draw_board(a, w_x, w_y, current_piece, current_legal_moves, pieces_a, false, current_king_index, check)
				gfx.UpdateAn()
				pieces_a, piece_promoting_to = pawn_promotion(w_x, w_y, a, int(piece_is_selected), pieces_a, "A", m_channel)
				piece_promoting_to = "=" + piece_promoting_to
			}
			move_string := get_move_string(int(piece_is_selected), original_field, piece_promoting_to, take, pieces_a) //create the pgn string for the sidebar

			return pieces_a, 64, true, move_string //return the changend array, a 64 (prevents the game in the main loop from executing this function several times), true on pos 3 indicates that a player change should be exectuted
		}
	}
	return pieces_a, piece_is_selected, false, "" //returning piece_is_selected esures that the piece is still selected after this block
}

func get_move_string(current_piece_index int, original_pos [2]uint16, piece_promoting string, take string, pieces_a [64]pieces.Piece) string {
	piece_name := pieces_a[current_piece_index].Give_Piece_Type()

	if piece_name == "K" && original_pos[0] == 4 && pieces_a[current_piece_index].Give_Pos()[0] == 6 { //short castle
		return "O-O"
	} else if piece_name == "K" && original_pos[0] == 4 && pieces_a[current_piece_index].Give_Pos()[0] == 2 { //long castle
		return "O-O-O"
	} else {

		var original_field string = ""
		var new_field string = parser.Get_Move_As_String_From_Field(pieces_a[current_piece_index].Give_Pos()) //get the current pos (which is the new pos) and translate this to pgn notation

		if piece_promoting != "" { //if the pieces isn't promoting, the piece name is not needed
			piece_name = ""
		}

		if take == "x" && piece_name == "" { //if a pawn takes someting, the original x cord of this pawn must be included
			original_field = parser.Translate_Field_Cord_To_PGN_String(original_pos[0], true)
		}

		//überprüfen ob es noch ein piece gibt welches auf das gleiche feld moven kann und vom selben typ ist
		if index_of_other_piece, _ := parser.Get_Piece_Index_And_Move_Index(pieces_a, pieces_a[current_piece_index].Give_Pos(), pieces_a[current_piece_index].Is_White_Piece(), piece_name, "0", current_piece_index); index_of_other_piece != 64 {
			//there is another piece that can execute the same move which means the pgn string is supposed to include more detailed information about the original position of the moved piece
			if pieces_a[index_of_other_piece].Give_Pos()[0] != original_pos[0] {
				original_field = parser.Translate_Field_Cord_To_PGN_String(original_pos[0], true)
			} else if pieces_a[index_of_other_piece].Give_Pos()[1] != original_pos[1] {
				original_field = parser.Translate_Field_Cord_To_PGN_String(original_pos[1], false)
			}
		}
		return piece_name + original_field + take + new_field + piece_promoting
	}
}

func display_message(message_type uint8, a uint16, white_is_current_player bool) bool {
	gfx.Archivieren() //needed if the user decides to go to the review mode

	gfx.Stiftfarbe(0, 0, 0)
	gfx.Transparenz(70)
	gfx.Vollrechteck(0, 0, a*8, a*8)
	gfx.Transparenz(50)
	gfx.Stiftfarbe(8, 8, 8)
	gfx.Vollrechteck((a), 3*a, 6*a, 2*a)

	gfx.Stiftfarbe(220, 220, 220)
	gfx.SetzeFont("./resources/fonts/junegull.ttf", int(5*a/10))

	if message_type == 0 { //checkmate or time is up
		if !white_is_current_player {
			gfx.SchreibeFont(18*a/10, 31*a/10, "player black has")
		} else {
			gfx.SchreibeFont(17*a/10, 31*a/10, "player white has")
		}
		gfx.Stiftfarbe(136, 8, 8)
		gfx.SetzeFont("./resources/fonts/punk.ttf", int(a))
		gfx.SchreibeFont(29*a/10, 37*a/10, "lost")
	} else if message_type == 1 { //stalemate
		gfx.SchreibeFont(295*a/100, 31*a/10, "That's a ")

		gfx.Stiftfarbe(0, 143, 230)
		gfx.SetzeFont("./resources/fonts/punk.ttf", int(a))
		gfx.SchreibeFont(144*a/100, 37*a/10, "stalemate")

	} else if message_type == 2 { //start message
		gfx.SchreibeFont(180*a/100, 31*a/10, "Press any key to")

		gfx.Stiftfarbe(28, 205, 60)
		gfx.SetzeFont("./resources/fonts/punk.ttf", int(a))
		gfx.SchreibeFont(260*a/100, 37*a/10, "star t")
	}

	gfx.Transparenz(0)
	gfx.UpdateAn()
	key, _, _ := gfx.TastaturLesen1()

	if key == 13 && message_type != 2 { //check if enter was pressed if not start message--> normal restart is initiated
		return false
	} else if key != 13 && message_type != 2 { //otherwise review mode is triggered --> restoring the board
		gfx.Restaurieren(0, 0, 8*a, 8*a)
		return true
	}
	return false //if the start message is displayed false will be returned with any key
}

func pawn_promotion(w_x, w_y, a uint16, pawn_index int, pieces_a [64]pieces.Piece, premoved string, m_channel chan [4]int16) ([64]pieces.Piece, string) {
	var piece_promoted_to string
	var queen pieces.Piece = pieces.NewQueen(pieces_a[pawn_index].Give_Pos()[0], pieces_a[pawn_index].Give_Pos()[1], pieces_a[pawn_index].Is_White_Piece())
	var knight pieces.Piece = pieces.NewKnight(pieces_a[pawn_index].Give_Pos()[0], pieces_a[pawn_index].Give_Pos()[1], pieces_a[pawn_index].Is_White_Piece())
	var rook pieces.Piece = pieces.NewRook(pieces_a[pawn_index].Give_Pos()[0], pieces_a[pawn_index].Give_Pos()[1], pieces_a[pawn_index].Is_White_Piece())
	var bishop pieces.Piece = pieces.NewBishop(pieces_a[pawn_index].Give_Pos()[0], pieces_a[pawn_index].Give_Pos()[1], pieces_a[pawn_index].Is_White_Piece())

	if premoved == "A" { //there is no premove, which means the user should decide the piece that is promoted to
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
		for { //loop as long as the user hasn't clicked on a piece to promote to
			mouse_input := <-m_channel
			if mouse_input[0] == 1 && mouse_input[1] == 1 && uint16(mouse_input[2]) >= 2*a && uint16(mouse_input[2]) <= 6*a && uint16(mouse_input[3]) >= (4*a)-(5*a)/10 && uint16(mouse_input[3]) <= (4*a)+(5*a)/10 {

				x_field_pos := (calc_field(a, uint16(mouse_input[2]), uint16(mouse_input[3]), (5*a)/10))[0] //calculate the field the user clicked on --> works cuz the displayed pieces are centered over known fields

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

	} else { //there is a premove --> no need for the user to decide
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
	return pieces_a, piece_promoted_to //return the piece_promoted_to string for the pgn sidebar notation
}

func draw_board(a, w_x, w_y uint16, current_piece pieces.Piece, current_legal_moves [][3]uint16, pieces_a [64]pieces.Piece, highlighting_is_activated bool, current_king_index int, check bool) {
	gfx.UpdateAus()
	draw_background(a)
	if check { //highlight the king that is in check
		highlight(a, pieces_a[current_king_index].Give_Pos(), 255, 0, 0)
	}
	if highlighting_is_activated {
		highlight(a, current_piece.Give_Pos(), 0, 50, 255) //highlight the current piece
		for k := 0; k < len(current_legal_moves); k++ {    //highlight the legal moves of the current piece
			highlight(a, [2]uint16{current_legal_moves[k][0], current_legal_moves[k][1]}, 0, 255, 0)
		}
	}
	draw_pieces(pieces_a, w_x, w_y, a)
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
	const display_limit int16 = 18 //how many rows of pgn notation can be displayed
	const lower_bound uint16 = 52  //indicates how long the move sidebar can be

	//draw the bg
	gfx.Stiftfarbe(124, 119, 111)
	gfx.Vollrechteck(8*a+a/10, 8*a/10, 9*a/10, (lower_bound-4)*a/10)
	gfx.Stiftfarbe(31, 32, 33)
	gfx.Vollrechteck(9*a, 8*a/10, 9*a/10, (lower_bound-4)*a/10)

	for i := moves_counter; moves_counter-i <= display_limit && i != 0; i-- {
		if i%2 != 0 { //white's move
			var move_number string = strconv.Itoa(int(i+1) / 2)

			//display move number on every move of white
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
			} else { //display the move
				gfx.SetzeFont("./resources/fonts/firamono.ttf", int(a/4))
				gfx.Stiftfarbe(124, 119, 111)
				gfx.SchreibeFont(90*a/10, lower_bound*a/10-5*a/10*uint16((moves_counter-i+1)/2), pgn_moves_a[i])
			}
		}
	}
}

func draw_player_names(name_player_white, name_player_black string, a uint16) {
	var max_name_lenght int

	//display bg
	gfx.Stiftfarbe(200, 191, 179)
	gfx.Vollrechteck(81*a/10, a/10, 18*a/10, 4*a/10)
	gfx.Vollrechteck(80*a/10, 6*a/10, 20*a/10, 1*a/10)
	gfx.Stiftfarbe(48, 46, 43)
	gfx.Vollrechteck(89*a/10, a/10, 2*a/10, 5*a/10)
	gfx.Stiftfarbe(48, 46, 43)

	if len(name_player_white) >= len(name_player_black) { //calculate the max length of the names that are supposed to be displayed, cuz every name should be displayed in the same size
		max_name_lenght = len(name_player_white)
	} else {
		max_name_lenght = len(name_player_black)
	}
	gfx.SetzeFont("./resources/fonts/firamono.ttf", int(a)/max_name_lenght)

	//write the names
	gfx.SchreibeFont(82*a/10, 2*a/10, name_player_white)
	gfx.SchreibeFont(92*a/10, 2*a/10, name_player_black)
}

func initialize(w_x, w_y, a uint16, restart bool, game_timer int64, name_player_white string, name_player_black string, image_location string) ([64]pieces.Piece, int, int, buttons.Button, buttons.Button, buttons.Button, buttons.Button, buttons.Button, [][64]pieces.Piece, time_counter.Counter, time_counter.Counter, []string) {
	var moves_a [][64]pieces.Piece
	var pgn_moves_a []string

	//init buttons
	var one_move_back buttons.Button = buttons.New(81*a/10, 7*a+a/10, 8*a/10, a-a/5, "<", 38, 37, 34, 200, 200, 200, (a / 4), int(a/2))
	var one_move_forward buttons.Button = buttons.New(91*a/10, 7*a+a/10, 8*a/10, a-a/5, ">", 38, 37, 34, 200, 200, 200, (a / 4), int(a/2))
	var restart_button buttons.Button = buttons.New(8*a+a/7, 62*a/10, 17*a/10, 3*a/10, "restart game", 86, 82, 77, 200, 200, 200, (a / 10), int(a/5))
	var pause_button buttons.Button = buttons.New(91*a/10, 66*a/10, 8*a/10, 3*a/10, "pause", 86, 82, 77, 200, 200, 200, (a / 15), int(a/5))
	var save_button buttons.Button = buttons.New(81*a/10, 66*a/10, 8*a/10, 3*a/10, "save", 86, 82, 77, 200, 200, 200, (a / 7), int(a/5))

	if !restart { //only create a new window if the game isn't restarted
		gfx.Fenster(w_x, w_y)
		gfx.Fenstertitel("Chess")
		rescale_image(a, image_location)
	}

	if game_timer == 0 { //if it is a friendly game (indicated by the game timer beeing set to 0) the pause button is not needed
		pause_button.Deactivate()
	}

	//draw sidebar lower bg
	gfx.Stiftfarbe(86, 82, 77)
	gfx.Vollrechteck(8*a, 7*a, 2*a, a)
	gfx.Stiftfarbe(48, 46, 43)
	gfx.Vollrechteck(8*a, 6*a, 2*a, a)

	//draw buttons
	one_move_back.Draw()
	one_move_forward.Draw()
	restart_button.Draw()
	pause_button.Draw()
	save_button.Draw()

	//draw sidebar upper background
	gfx.Stiftfarbe(48, 46, 43)
	gfx.Vollrechteck(8*a, 0, 2*a, 6*a)

	draw_player_names(name_player_white, name_player_black, a)
	draw_background(a)

	var pieces_a [64]pieces.Piece
	var white_king_index int = -1 //creates an error if the game is started without a white king
	var black_king_index int = -1 //creates an error if the game is started without a black king

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

	//find the king indexes
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

	//init the timers
	white_time_counter := time_counter.New(game_timer)
	black_time_counter := time_counter.New(game_timer)

	return pieces_a, white_king_index, black_king_index, one_move_back, one_move_forward, restart_button, pause_button, save_button, moves_a, white_time_counter, black_time_counter, pgn_moves_a
}

func append_moves_array(moves_a [][64]pieces.Piece, pieces_a [64]pieces.Piece) [][64]pieces.Piece {
	moves_a = append(moves_a, pieces.Copy_Array(pieces_a)) //create a deep copy of the array of pieces
	return moves_a
}

func draw_pieces(pieces_a [64]pieces.Piece, w_x, w_y, a uint16) {
	gfx.UpdateAus()
	for i := 0; i < len(pieces_a); i++ {
		if pieces_a[i] != nil {
			gfx.Archivieren()                     //save the current board state
			pieces.Draw(pieces_a[i], w_x, w_y, a) //temporarily paste the whole image of all pieces and copy it to clipboard--> unfortuanetely this is needed in gfx
			gfx.Restaurieren(0, 0, w_x, w_y)      //reset the board

			gfx.Clipboard_einfuegenMitColorKey(pieces_a[i].Give_Pos()[0]*a, pieces_a[i].Give_Pos()[1]*a, 5, 5, 5) //paste the piece from the clipboard
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
	gfx.Vollrechteck(cord_x, cord_y, a, a) //transparent square over the piece
	gfx.Transparenz(0)
}

func draw_background(a uint16) { //create the board tiles
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
	current_field[1] = (m_y + y_offset) / a //y_offset was part of a feature that is removed now, but it doesn't affect the code so it just stays there
	return current_field
}

func set_image_string(troll_mode bool) string { //decides which sprite inage should be used for the pieces
	if troll_mode {
		return "./resources/images/Troll.bmp"
	} else {
		return "./resources/images/Pieces_Source_Original.bmp"
	}
}
