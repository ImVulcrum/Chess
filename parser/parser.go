package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"../pieces"
)

func clean_pgn(input_string string) string {

	//make spaces after the points
	parts := strings.Split(input_string, ".")
	input_string = strings.Join(parts, ". ")

	//remove tag section
	var string_without_tags string = input_string
	index := strings.LastIndex(input_string, "]")                    //find the last instance of a closed bracket
	if index != -1 && !strings.Contains(input_string[index:], "[") { //Check if there are no more open brackets after the first closed bracket
		string_without_tags = input_string[index+1:] //Extract the substring starting from the character after the first closed bracket
	} //else: this string does not have a comment section

	//bring everything to one line and remove additional spaces between characters
	reg := regexp.MustCompile(`\s+`)
	cleaned := reg.ReplaceAllString(string_without_tags, " ")
	cleaned = strings.TrimSpace(cleaned)

	//remove exvery instance of x (indicates a capture) and plus (indicates a check) from the string cuz otherwise the parsing would be overly complicated
	var result strings.Builder
	for i := 0; i < len(cleaned); i++ {
		if !(string(cleaned[i]) == "x" || string(cleaned[i]) == "+") {
			result.WriteByte(cleaned[i])
		}
	}
	//remove comments
	re := regexp.MustCompile(`\{[^}]*\}`)
	cleaned_string := re.ReplaceAllString(result.String(), "")
	return cleaned_string
}

func Create_Array_Of_Moves(input_string string) []string {

	var cleaned_string string = clean_pgn(input_string)

	//create a list with every move
	moves := strings.Split(cleaned_string, " ")
	var cleanedMoves []string
	for _, move := range moves {
		if move != "" && !unicode.IsDigit(rune(move[0])) { //wenn der string mit einer nummer startet, das heißt enweder wenn es sich um die zahl des moves handelt oder wenn das endergebnis aufgeschrieben wird
			cleanedMoves = append(cleanedMoves, move)
		}
	}
	// for i := 0; i < len(cleanedMoves); i++ {
	// 	fmt.Println("-" + cleanedMoves[i] + "-")
	// }
	// fmt.Println("------------")
	return cleanedMoves
}

func Get_Correct_Move(move string, pieces_a [64]pieces.Piece, current_king_index int) (int, int, string) {
	var field [2]uint16
	var piece_executing_move int
	var index_of_correct_legal_move int
	var pawn_promotion_to_piece string = "A" //A indicates that there is no pawn promotion

	if move[len(move)-1] != 'O' { //normal move
		if move[len(move)-2:len(move)-1] == "=" { //Pawn promotion
			pawn_promotion_to_piece = move[len(move)-1:]
			move = move[:len(move)-2]
		}

		field = Get_Field_From_Move(move)

		if firstChar := rune(move[0]); !unicode.IsUpper(firstChar) { //pawn move cuz the string does not start with an uppercase letter
			if len(move) == 2 { // Simple Pawn Move
				piece_executing_move, index_of_correct_legal_move = Get_Piece_Index_And_Move_Index(pieces_a, field, pieces_a[current_king_index].Is_White_Piece(), "P", "0")
			} else if len(move) == 3 {
				piece_executing_move, index_of_correct_legal_move = Get_Piece_Index_And_Move_Index(pieces_a, field, pieces_a[current_king_index].Is_White_Piece(), "P", string(move[0]))
			}
		} else { //Piece move cuz the move string starts with an uppercase letter
			if len(move) == 3 { // simple piece move, the first character of the move string indicates the piece that is moving
				piece_executing_move, index_of_correct_legal_move = Get_Piece_Index_And_Move_Index(pieces_a, field, pieces_a[current_king_index].Is_White_Piece(), string(move[0]), "0")
			} else if len(move) == 4 { //piece move with position, the first character of the move string indicates the piece that is moving, the second the starting position
				piece_executing_move, index_of_correct_legal_move = Get_Piece_Index_And_Move_Index(pieces_a, field, pieces_a[current_king_index].Is_White_Piece(), string(move[0]), string(move[1]))
			}
		}

	} else {
		if move == "O-O" { //short castle
			piece_executing_move, index_of_correct_legal_move = Get_Piece_Index_And_Move_Index(pieces_a, [2]uint16{7, pieces_a[current_king_index].Give_Pos()[1]}, pieces_a[current_king_index].Is_White_Piece(), "K", "0")
		} else if move == "O-O-O" { //long castle
			piece_executing_move, index_of_correct_legal_move = Get_Piece_Index_And_Move_Index(pieces_a, [2]uint16{0, pieces_a[current_king_index].Give_Pos()[1]}, pieces_a[current_king_index].Is_White_Piece(), "K", "0")
		} else {
			fmt.Println("Error while Reading Premove File: Expected either (O-O) or (O-O-O), got", move, "instead")
		}
	}
	return piece_executing_move, index_of_correct_legal_move, pawn_promotion_to_piece
}

func Get_Piece_Index_And_Move_Index(pieces_a [64]pieces.Piece, field [2]uint16, white_is_current_player bool, piece_type string, position string) (int, int) {
	var current_piece_type string = "A"

	for i := 0; i < len(pieces_a); i++ {
		if pieces_a[i] != nil && pieces_a[i].Is_White_Piece() == white_is_current_player {
			for k := 0; k < len(pieces_a[i].Give_Legal_Moves()); k++ {
				if pieces_a[i].Give_Legal_Moves()[k][0] == field[0] && pieces_a[i].Give_Legal_Moves()[k][1] == field[1] { //there is a piece in the correct color with the given move
					cord, is_x_cord := Translate_PGN_Field_Notation(position)

					if (is_x_cord && cord == pieces_a[i].Give_Pos()[0]) || (!is_x_cord && cord == pieces_a[i].Give_Pos()[1]) || cord == 8 { //check if the piece has the given x or y cord or has no cord specifictaion indicated by cord beeing 8
						switch pieces_a[i].(type) {
						case *pieces.Rook:
							current_piece_type = "R"
						case *pieces.King:
							current_piece_type = "K"
						case *pieces.Pawn:
							current_piece_type = "P"
						case *pieces.Queen:
							current_piece_type = "Q"
						case *pieces.Bishop:
							current_piece_type = "B"
						case *pieces.Knight:
							current_piece_type = "N"
						default:
							fmt.Println("Error in Parser while iterating through pieces array: Unexpected piece type")
						}
						if current_piece_type == piece_type {
							return i, k
						}
					}
				}
			}
		}
	}
	panic("Error in Parser: there is no piece is the pieces array that matches the specifications given, which means that the given pgn file is corrupted")
}

func Translate_PGN_Field_Notation(cord_string string) (uint16, bool) {
	var cord uint16
	var is_x_cord bool

	if len(cord_string) != 1 {
		fmt.Println("Error: Unexpected lenght of string while trying to convert it from pgn field notation to a square notation")
		cord = 8
	} else {
		if unicode.IsDigit(rune(cord_string[0])) {
			is_x_cord = false
			num, _ := strconv.Atoi(cord_string)
			num = num - 8
			num = -num
			cord = uint16(num)
		} else {
			is_x_cord = true
			cord = uint16(cord_string[0] - 'a')
		}
	}
	return cord, is_x_cord
}

func Get_Field_From_Move(move string) [2]uint16 {
	var field [2]uint16

	var x_cord string = move[len(move)-2 : len(move)-1]
	var y_cord string = move[len(move)-1:]

	x, _ := Translate_PGN_Field_Notation(x_cord)
	y, _ := Translate_PGN_Field_Notation(y_cord)

	field = [2]uint16{x, y}
	return field
}

func Translate_Field_Cord_To_PGN_String(field_cord uint16, is_x_cord bool) string {
	if is_x_cord {
		return string(rune(int("a"[0]) + int(field_cord)))
	} else {
		return strconv.Itoa(-1 * (int(field_cord) - 8))
	}
}

func Get_Move_As_String_From_Field(field [2]uint16) string {
	return (Translate_Field_Cord_To_PGN_String(field[0], true) + Translate_Field_Cord_To_PGN_String(field[1], false))
}
