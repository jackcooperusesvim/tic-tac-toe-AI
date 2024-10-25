package main

import (
	"errors"
	"fmt"
	"github.com/inancgumus/screen"
	"strconv"
	// "sync"
)

const TEST = false

const x = 1
const o = -1
const blank = 0
const bad_piece = 3

type board struct {
	brd         [3][3]int
	turn        int
	winner      int
	is_terminal bool
}

type vboard struct {
	brd    [3][3]string
	turn   string
	winner string
}

var TestCases []board = []board{}

func (b board) int_to_vb_piece(in int, place_num int) (string, error) {
	if in == x {
		return "x", nil
	}
	if in == o {
		return "o", nil
	}
	if in == blank {
		return fmt.Sprintf("%d", place_num), nil
	}
	return " ", errors.New("Invalid piece")
}

func (b board) to_vboard(nums bool) (vb vboard, err error) {
	for i := 1; i <= 9; i++ {
		t := int_to_turn(i)
		p := b.get(t)
		if p == blank {
			if nums {
				vb = vb.set(t, fmt.Sprintf("%v", i))
			} else {
				vb = vb.set(t, " ")
			}
		} else if p == x {
			vb = vb.set(t, "x")
		} else {
			vb = vb.set(t, "o")
		}
		if err != nil {
			return vb, err
		}
	}
	if b.is_terminal && b.winner == blank {
		vb.winner = "Tie Game"
	}
	return vb, nil
}

type turn struct {
	x    int
	y    int
	null bool
}

type winSum struct {
	row_wins  [3]bool
	col_wins  [3]bool
	diag_wins [2]bool
}

func initializeWinSum() (ws winSum) {
	ws.row_wins = [3]bool{true, true, true}
	ws.col_wins = [3]bool{true, true, true}
	ws.diag_wins = [2]bool{true, true}
	return ws

}

func slc_or(input []bool) bool {
	for _, e := range input {
		if e {
			return true
		}
	}
	return false
}

func (ws winSum) compile() bool {
	// fmt.Printf("	dbg: cols: %v,%v,%v\n",ws.row_wins[0],ws.row_wins[1],ws.row_wins[2])
	return (slc_or(ws.row_wins[:]) || slc_or(ws.col_wins[:]) || slc_or(ws.diag_wins[:]))
}

func initializeBrd() board {
	var out = board{}
	out.turn = x
	out.winner = 0
	out.is_terminal = false

	return out
}

func persistent_input(options ...string) int {
	var in string
	fmt.Scanln(&in)
	for i := 0; i < len(options); i++ {
		if options[i] == in {
			return i
		}
	}
	out := persistent_input(options...)
	return out
}
func u_input_int(retry func(board), b board) int {
	var in string
	fmt.Printf("Num : ")
	fmt.Scanln(&in)
	intIn, err := strconv.Atoi(in)
	if err != nil || intIn > 9 || intIn < 1 {
		retry(b)
		intIn = u_input_int(retry, b)
	}

	return intIn
}

func u_input_turn(retry func(board), b board) turn {
	num := u_input_int(retry, b)
	fmt.Printf("\n")
	return int_to_turn(num)

}
func int_to_turn(num int) turn {
	tn := turn{}
	num = 9 - num
	tn.x = 2 - (num % 3)
	tn.y = (num - (num % 3)) / 3

	return tn
}

func slc_sum(in []int) (sum int) {
	for _, e := range in {
		sum += e
	}
	return
}

func (b board) update_win_state() (board, error) {
	x_winSum := initializeWinSum()
	o_winSum := initializeWinSum()

	for col := 0; col < 3; col++ {
		for row := 0; row < 3; row++ {
			// fmt.Printf("dbg: val: %v col:%v row:%v \n",b.brd[row][col],col,row)
			//Check Linear solutions
			if b.brd[row][col] == blank {
				// fmt.Printf("dbg: blank\n")
				x_winSum.row_wins[row] = false
				x_winSum.col_wins[col] = false

				o_winSum.row_wins[row] = false
				o_winSum.col_wins[col] = false

			} else if b.brd[row][col] == x {
				// fmt.Printf("dbg: x\n")
				o_winSum.row_wins[row] = false
				o_winSum.col_wins[col] = false

			} else {
				// fmt.Printf("dbg: o\n")
				x_winSum.row_wins[row] = false
				x_winSum.col_wins[col] = false
			}

			//Check Diagonals
			if 2-row == col {
				if b.brd[row][col] == blank {
					x_winSum.diag_wins[0] = false
					o_winSum.diag_wins[0] = false

				} else if b.brd[row][col] == x {
					o_winSum.diag_wins[0] = false

				} else {
					x_winSum.diag_wins[0] = false
				}
			}

			if row == col {
				if b.brd[row][col] == blank {
					x_winSum.diag_wins[1] = false
					o_winSum.diag_wins[1] = false

				} else if b.brd[row][col] == x {
					o_winSum.diag_wins[1] = false

				} else {
					x_winSum.diag_wins[1] = false
				}
			}
		}
	}
	x_win := x_winSum.compile()
	o_win := o_winSum.compile()
	b.is_terminal = false
	if x_win && o_win {
		b.winner = bad_piece
		b.is_terminal = true
		return b, errors.New("Both players have won")
	}
	if x_win {
		b.winner = x
		b.is_terminal = true
		return b, nil
	}
	if o_win {
		b.winner = o
		b.is_terminal = true
		return b, nil
	}
	if len(b.legal_moves()) == 0 {
		b.is_terminal = true
	}
	b.winner = blank
	return b, nil
}

func (b board) is_legal(t turn) bool {
	if b.get(t) == blank {
		return true
	}
	return false
}
func (b board) get(t turn) int {
	return b.brd[t.y][t.x]
}
func (b board) set(t turn, set_to int) board {
	b.brd[t.y][t.x] = set_to
	return b
}
func (vb vboard) set(t turn, set_to string) vboard {
	vb.brd[t.y][t.x] = set_to
	return vb
}

func (b board) print_board(nums bool) {
	vb, err := b.to_vboard(nums)
	if err != nil {
		println("ERROR: ", err)
		return
	}
	fmt.Println("    |   |    ")
	for i := 0; i < 3; i++ {
		fmt.Println(" ", vb.brd[i][0], "|", vb.brd[i][1], "|", vb.brd[i][2], " ")
		if i == 2 {
			fmt.Println("    |   |    ")
		} else {
			fmt.Println("----+---+----")
		}
	}
	if b.winner == x {
		fmt.Println("x has won")
	} else if b.winner == o {
		fmt.Println("o has won")
	} else if b.is_terminal {
		fmt.Println("Tie Game")
	} else {
		fmt.Println("~~~~~~~~~~~~")
		util, err := b.utility_assign_rec()
		if util == x {
			fmt.Println("o screwed up and will probably lose")
		} else if util == o {
			fmt.Println("x screwed up and will probably lose")
		} else {
			fmt.Println("perfect play thus far")
		}

		if err != nil {
			fmt.Println("Utility Could Not Be Calculated")
		}
		fmt.Println("~~~~~~~~~~~~")
		fmt.Println("no one has won")
	}
}

func (b board) apply_turn(t turn) (board, error) {
	if !b.is_legal(t) {
		return b, errors.New("non-legal move")
	}
	if b.is_terminal {
		return b, errors.New("game is over")
	}
	if b.winner != blank {
		return b, errors.New("game is already won (also, terminal state was not updated)")
	}

	b.brd[t.y][t.x] = b.turn

	if b.turn == x {
		b.turn = o
	} else {
		b.turn = x
	}

	return b, nil
}
func main() {
	if TEST {
		test()
	} else {
		mn()
	}
}

func test() {
	const turnsPerTest = 5
	const numOfTests = 4
	turns := [numOfTests][turnsPerTest]turn{
		//checks vertical x
		{turn{x: 0, y: 0}, turn{x: 1, y: 1}, turn{x: 0, y: 1}, turn{x: 2, y: 2}, turn{x: 0, y: 2}},
		//checks horizontal x
		{turn{x: 0, y: 0}, turn{x: 1, y: 1}, turn{x: 1, y: 0}, turn{x: 2, y: 2}, turn{x: 2, y: 0}},
		//checks diag 1 x
		{turn{x: 0, y: 0}, turn{x: 0, y: 1}, turn{x: 1, y: 1}, turn{x: 0, y: 2}, turn{x: 2, y: 2}},
		//checks diag 2 x
		{turn{x: 0, y: 2}, turn{x: 0, y: 1}, turn{x: 1, y: 1}, turn{x: 2, y: 1}, turn{x: 2, y: 0}}}

	//each test will be done with both x and o

	numFailed := 0
	numPassed := 0
	var testBrd board = initializeBrd()
	passed := true
	for i, test := range turns {
		testBrd = initializeBrd()
		passed = true
		var err error = nil

		fmt.Printf("TEST #%v\n", (i + 1))

		for _, turn := range test {

			fmt.Println("	APPLYING TURN")
			testBrd, err = testBrd.apply_turn(turn)
			if err != nil {
				fmt.Println("		FAIL")
				fmt.Printf("			%v \n", err)
			} else {
				fmt.Println("		PASS")
			}

			fmt.Println("	UPDATING WIN STATE")
			testBrd, err = testBrd.update_win_state()
			if err != nil {
				passed = false
				fmt.Println("		FAIL")
				fmt.Printf("			%v \n", err)
			} else {
				fmt.Println("		PASS")
			}

		}
		fmt.Println("	CHECKING WIN STATE")
		testBrd.print_board(false)
		if testBrd.winner == x {
			fmt.Println("		PASS")
		} else {
			passed = false
			fmt.Println("		FAIL")
			fmt.Printf("			Expected: %v \n", x)
			fmt.Printf("			Recieved: %v \n", testBrd.winner)
		}
		if passed {
			numPassed++
		} else {
			numFailed++
		}
		testBrd = initializeBrd()

	}
	fmt.Printf("\n\n\n")
	fmt.Printf("SUMMARY:\n")
	fmt.Printf("	PASSED: %v", numPassed)
	fmt.Printf("	FAILED: %v", numFailed)
}

func retry_screen(board board) {
	screen.Clear()
	screen.MoveTopLeft()
	board.print_board(false)
	fmt.Printf("Not a valid input (0-9)\n")
}
func reset_screen(board board) {
	screen.Clear()
	screen.MoveTopLeft()
	board.print_board(false)
}

func first_or_second() (first bool) {
	fmt.Println("Would you like to go first or second? (1/2)(f/s)(first/second)")
	if persistent_input("first", "second", "f", "s", "1", "2")%2 == 0 {
		first = true
	}
	return
}
func (brd board) human_cycle() board {
	reset_screen(brd)
	var err error = errors.New("")
	for err != nil {
		tn := u_input_turn(retry_screen, brd)
		brd, err = brd.apply_turn(tn)
		reset_screen(brd)
		if err != nil {
			fmt.Println("Not a valid input: (", err, ")")
		}
	}

	brd, err = brd.update_win_state()
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	return brd
}

func (brd board) minimax_cycle() board {
	w_turn, err := brd.choose_move()
	if err != nil {
		fmt.Println("ERROR: ", err)
	}

	brd, err = brd.apply_turn(w_turn)
	if err != nil {
		fmt.Println("ERROR: ", err)
	}

	brd, err = brd.update_win_state()
	if err != nil {
		fmt.Println("ERROR: ", err)
	}
	return brd
}
func mn() {
	brd := initializeBrd()
	if !first_or_second() {
		brd = brd.minimax_cycle()
	}
	for true {
		brd = brd.human_cycle()

		//end game if ended
		if brd.is_terminal {
			reset_screen(brd)
			var asdf string
			fmt.Printf("Press Enter to Quit")
			fmt.Scanln(&asdf)
			screen.Clear()
			screen.MoveTopLeft()
			return
		}
		brd = brd.minimax_cycle()

		//end game if ended
		if brd.is_terminal {
			reset_screen(brd)
			var asdf string
			fmt.Printf("Press Enter to Quit")
			fmt.Scanln(&asdf)
			screen.Clear()
			screen.MoveTopLeft()
			return
		}
	}

}

//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

func (b board) legal_moves() (available_turns []turn) {
	var working_turn turn
	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			working_turn.y = row
			working_turn.x = col
			if b.is_legal(working_turn) {
				available_turns = append(available_turns, working_turn)
			}
		}
	}

	return available_turns
}
func (b board) choose_move() (move turn, err error) {
	b, err = b.update_win_state()

	if err != nil {
		return turn{}, err
	}

	if b.winner != blank {
		return turn{}, errors.New("ERROR: Game already won")
	}

	moves := b.legal_moves()
	if len(moves) == 0 {
		return turn{}, nil
	}
	var util []int
	var w_util int
	var w_brd board
	for _, move := range moves {
		w_brd, err = b.apply_turn(move)
		w_util, err = w_brd.utility_assign_rec()
		if err != nil {
			return turn{}, err
		}
		util = append(util, w_util)
	}
	var move_index int
	closest := util[0]
	if b.turn == x {
		for i, val := range util {
			if val == x {
				return moves[i], nil
			}
			if closest < val {
				closest = val
				move_index = i
			}
		}
	} else if b.turn == o {
		closest = util[0]
		for i, val := range util {
			if val == o {
				return moves[i], nil
			}
			if closest > val {
				closest = val
				move_index = i
			}
		}

	}
	return moves[move_index], nil

}
func (b board) utility_assign_rec() (out int, err error) {
	b, err = b.update_win_state()

	if err != nil {
		return blank, err
	}

	if b.winner != blank {
		return b.winner, nil
	}

	moves := b.legal_moves()
	if len(moves) == 0 {
		return blank, nil
	}
	var util []int
	var w_util int
	var w_brd board
	for _, move := range moves {
		w_brd, err = b.apply_turn(move)
		if err != nil {
			return blank, err
		}
		w_util, err = w_brd.utility_assign_rec()
		if err != nil {
			return blank, err
		}
		util = append(util, w_util)
	}
	closest := util[0]
	if b.turn == x {
		closest = o
		closest = util[0]
		for _, val := range util {
			if val == x {
				return x, nil
			}
			if closest < val {
				closest = val
			}
		}
	} else if b.turn == o {
		closest = util[0]
		closest = x
		for _, val := range util {
			if val == o {
				return o, nil
			}
			if closest > val {
				closest = val
			}
		}

	}
	return closest, nil

}

//
// func (b board) utility_assign_rec_multi_master(stop_util int) (util int){
// 	x_win_chan := make(chan struct{})
// 	o_win_chan := make(chan struct{})
// 	draw_chan := make(chan struct{})
// 	sub_threads_count := 9
// 	for _, v := range sub_threads_count {
//
// 	}
// 	results_recieved :=0
// 	for {
// 		select {
// 			case <- x_win_chan:
// 				if stop_util == x {
// 					return x
// 				}
// 				results_recieved++
// 			case <- o_win_chan:
// 				if stop_util == o {
// 					return x
// 				}
// 				results_recieved++
// 			case <- draw_chan:
// 				results_recieved++
// 		}
// 	}
// }
