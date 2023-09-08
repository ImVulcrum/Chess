package parser

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"../pieces"
)

func clean_pgn(input_string string) string {

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

func Create_Array_Of_Moves() []string {
	inputString := `[Event "F/S Return Match"]
	[Site "Belgrade, Serbia JUG"]
	[Date "1992.11.04"]
	[Round "29"]
	[White "Fischer, Robert J."]
	[Black "Spassky, Boris V."]
	[Result "1/2-1/2"]
	
	1. e4 e5 2. Nf3 Nc6 3. Bb5 a6 {This opening is called the Ruy Lopez.}
	4. Ba4 Nf3 5. O-O Be7 6. Re1 b5 7. Bb3 d6 8. c3 O-O 9. h3 Nb8 10. d4 Nbd7
	11. c4 c6 12. cxb5 axb5 13. Nc3 Bb7 14. Bg5 b4 15. Nb1 h6 16. Bh4 c5 17. dxe5
	Nxe4 18. Bxe7 Qxe7 19. exd6 Qf6 20. Nbd2 Nxd6 21. Nc4 Nxc4 22. Bxc4 Nb6
	23. Ne5 Rae8 24. Bxf7+ Rxf7 25. Nxf7 Rxe1+ 26. Qxe1 Kxf7 27. Qe3 Qg5 28. Qxg5
	hxg5 29. b3 Ke6 30. a3 Kd6 31. axb4 cxb4 32. Ra5 Nd5 33. f3 Bc8 34. Kf2 Bf5
	35. Ra7 g6 36. Ra6+ Kc5 37. Ke1 Nf4 38. g3 Nxh3 39. Kd2 Kb5 40. Rd6 Kc5 41. Ra6
	Nf2 42. g4 Bd3 43. Re6 1/2-1/2`

	var cleaned_string string = clean_pgn(inputString)

	//create a list with every move
	moves := strings.Split(cleaned_string, " ")
	var cleanedMoves []string
	for _, move := range moves {
		if move != "" && !unicode.IsDigit(rune(move[0])) { //wenn der string mit einer nummer startet, das heißt enweder wenn es sich um die zahl des moves handelt oder wenn das endergebnis aufgeschrieben wird
			cleanedMoves = append(cleanedMoves, move)
		}
	}

	fmt.Println("start of list")
	for i := 0; i < len(cleanedMoves); i++ {
		fmt.Println("-" + cleanedMoves[i] + "-")
	}
	fmt.Println("-----------------------------------------dfdsfgsdgsd----------------------------")
	return cleanedMoves
}

func Get_Correct_Move(move string, pieces_a [64]pieces.Piece, current_king_index int) ([3]uint16, int) {
	var long_rook int = 64
	var short_rook int = 64
	var correct_move [3]uint16
	var field [2]uint16
	var piece_executing_move int

	if move[len(move)-1] != 'O' {
		field = Get_Field_From_Move(move, current_king_index)
		fmt.Println(field)

	} else { //rochade
		for i := 0; i < len(pieces_a); i++ {
			if pieces_a[i] != nil {
				if pieces_a[i].Give_Pos()[0] == 0 && pieces_a[i].Give_Pos()[1] == pieces_a[current_king_index].Give_Pos()[1] {
					long_rook = i
				}
				if pieces_a[i].Give_Pos()[0] == 7 && pieces_a[i].Give_Pos()[1] == pieces_a[current_king_index].Give_Pos()[1] {
					short_rook = i
				}

			}
		}
	}
	if move == "O-O" { //short castle
		correct_move = [3]uint16{0, pieces_a[current_king_index].Give_Pos()[1], uint16(short_rook)}
	} else if move == "O-O-O" { //long castle
		correct_move = [3]uint16{7, pieces_a[current_king_index].Give_Pos()[1], uint16(long_rook)}
	} else {
		fmt.Println("Error while Reading Premove File: Expected either (O-O) or (O-O-O), got", move, "instead")
	}
	return correct_move, piece_executing_move
}

func Get_Field_From_Move(move string, current_king_index int) [2]uint16 {
	var Field [2]uint16

	return Field
}
