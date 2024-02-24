package main

import (
	"fmt"
	"errors"
	"github.com/inancgumus/screen"
	"strconv"
)
const TEST = false 

const x = 1
const o = -1
const blank = 0
const bad_piece = 3

type board struct {
	brd [3][3]int
	turn int
	winner int
}

type vboard struct {
	brd [3][3]string
	turn string
	winner string
}

var TestCases []board = []board{}

func int_to_piece(in int) (string, error) {
	if in == x{
		return "x", nil
	}
	if in == o{
		return "o", nil
	}
	if in == blank{
		return " ", nil
	}
	return " ", errors.New("Invalid piece")
}

func (b board) to_vboard() (vb vboard, err error) {
	for i := 0 ; i<3 ; i++{
		for j := 0 ; j<3 ; j++{
			vb.brd[i][j], err = int_to_piece(b.brd[i][j])
			if err != nil {
				return vb , err
			}
		}
	}
	return vb, nil
}
type turn struct {
	x int
	y int
}

type winSum struct {
	row_wins [3]bool
	col_wins [3]bool
	diag_wins [2]bool

}

func initializeWinSum() (ws winSum) {
	ws.row_wins = [3]bool{true,true,true}
	ws.col_wins = [3]bool{true,true,true}
	ws.diag_wins = [2]bool{true,true}
	return ws

}

func slc_or(input []bool) bool {
	for _,e := range input {
		if e{
			return true
		}
	}
	return false
}

func (ws winSum) compile() bool{
	// fmt.Printf("	dbg: cols: %v,%v,%v\n",ws.row_wins[0],ws.row_wins[1],ws.row_wins[2])
	return (slc_or(ws.row_wins[:]) || slc_or(ws.col_wins[:]) ||slc_or(ws.diag_wins[:]))
}

func initializeBrd() board{
	var out = board{}
	out.turn = x
	out.winner = 0

	return out
}

func persistent_input(options ...string) int {
	var in string
	fmt.Scanln(&in)
	for i := 0 ; i<len(options); i++{
		if options[i] == in{
			return i
		}
	}
	out := persistent_input(options...)
	return out
}
func int_move_input(retry func(board), b board) int {
	var in string
	fmt.Printf("Num : ")
	fmt.Scanln(&in)
	intIn, err := strconv.Atoi(in)
	if err != nil || intIn>9 || intIn<1{
		retry(b)
		intIn = int_move_input(retry, b)
	}

	return intIn
}
	

func generate_turn_from_int(num int) turn{
	fmt.Printf("\n")


	tn := turn{}
	num = 9-num
	tn.x = 2-(num%3)
	tn.y = (num-(num%3))/3

	return tn

}


func generate_turn() turn {
	out := turn{}
	x_word_options := [3]string{"l","m","r"}
	y_word_options := [3]string{"t","m","b"}

	fmt.Println("Enter the x turn (l,m,r):")
	out.x = persistent_input(x_word_options[:]...)
	fmt.Println("Enter the y turn (b,m,t):")
	out.y = persistent_input(y_word_options[:]...)


	return out
}

func slc_sum(in []int) (sum int) {
	for _,e := range in {
		sum += e
	}
	return
}


func (b board) update_win_state() (board, error){
	x_winSum := initializeWinSum() 
	o_winSum := initializeWinSum() 

	for col := 0 ; col<3 ; col++{
		for row := 0 ; row<3 ; row++{
			// fmt.Printf("dbg: val: %v col:%v row:%v \n",b.brd[row][col],col,row)
			//Check Linear solutions
		 	if b.brd[row][col] == blank{
				// fmt.Printf("dbg: blank\n")
				x_winSum.row_wins[row] = false
				x_winSum.col_wins[col] = false

				o_winSum.row_wins[row] = false
				o_winSum.col_wins[col] = false

			} else if b.brd[row][col] == x{
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
				if b.brd[row][col] == blank{
					x_winSum.diag_wins[0] = false
					o_winSum.diag_wins[0] = false

				} else if b.brd[row][col] == x{
					o_winSum.diag_wins[0] = false

				} else {
					x_winSum.diag_wins[0] = false
				}
			}

			if row == col {
				if b.brd[row][col] == blank{
					x_winSum.diag_wins[1] = false
					o_winSum.diag_wins[1] = false

				} else if b.brd[row][col] == x{
					o_winSum.diag_wins[1] = false

				} else {
					x_winSum.diag_wins[1] = false
				}
			}
		}
	}
	x_win := x_winSum.compile()
	o_win := o_winSum.compile()
	if x_win && o_win{
		b.winner = bad_piece
		return b, errors.New("Both players have won")
	}
	if x_win {
		b.winner = x
		return b, nil
	}
	if o_win {
		b.winner = o
		return b, nil
	}
	b.winner = blank
	return b, nil
}

func (b board) is_legal(t turn) bool{
	if b.brd[t.y][t.x] == blank{
		return true
	}
	return false
}

func (b board) print_board(){
	vb, err := b.to_vboard()
	if err != nil{
		println("ERROR: ",err)
		return
	}
	fmt.Println("    |   |    ")
	for i :=0; i<3; i++ {
		fmt.Println(" ",vb.brd[i][0],"|",vb.brd[i][1],"|",vb.brd[i][2]," ")
		if i==2 {
			fmt.Println("    |   |    ")
		} else {
			fmt.Println("----+---+----")
		}
	}
	fmt.Println("~~~~~~~~~~~~")
	if b.winner == x{
	fmt.Println("x has won")
	} else if b.winner == o{
	fmt.Println("o has won")
	} else {
	fmt.Println("no one has won")
	}
}

func (b board) apply_turn(t turn) (board, error){
	if !b.is_legal(t){
		return b, errors.New("Non-Legal move")
	}
	if b.winner != blank{
		return b, errors.New("Game already won")
	}

	b.brd[t.y][t.x] = b.turn

	if b.turn == x {
		b.turn = o
	} else {
		b.turn = x
	}

	return b, nil
}
func main(){
	if TEST{test()} else {mn()}
}

func test(){
	const turnsPerTest = 5
	const numOfTests = 4
	turns := [numOfTests][turnsPerTest]turn{
		//checks vertical x
		{turn{x:0,y:0},turn{x:1,y:1},turn{x:0,y:1},turn{x:2,y:2},turn{x:0,y:2}},
		//checks horizontal x 
		{turn{x:0,y:0},turn{x:1,y:1},turn{x:1,y:0},turn{x:2,y:2},turn{x:2,y:0}},
		//checks diag 1 x 
		{turn{x:0,y:0},turn{x:0,y:1},turn{x:1,y:1},turn{x:0,y:2},turn{x:2,y:2}},
		//checks diag 2 x
		{turn{x:0,y:2},turn{x:0,y:1},turn{x:1,y:1},turn{x:2,y:1},turn{x:2,y:0}}}
	
	//each test will be done with both x and o

	numFailed := 0
	numPassed := 0
	var testBrd board = initializeBrd()
	passed := true
	for i, test := range turns{
		testBrd = initializeBrd()	
		passed = true
		var err error = nil

		fmt.Printf("TEST #%v\n",(i+1))

		for _, turn := range test {

			fmt.Println("	APPLYING TURN")
			testBrd, err = testBrd.apply_turn(turn)
			if err != nil {
				fmt.Println("		FAIL")
				fmt.Printf("			%v \n",err)
			} else {
				fmt.Println("		PASS")
			}


			fmt.Println("	UPDATING WIN STATE")
			testBrd, err = testBrd.update_win_state()
			if err != nil {
				passed = false
				fmt.Println("		FAIL")
				fmt.Printf("			%v \n",err)
			} else {
				fmt.Println("		PASS")
			}

		}
		fmt.Println("	CHECKING WIN STATE")
		testBrd.print_board()
		if testBrd.winner == x {
			fmt.Println("		PASS")
		} else {
			passed = false
			fmt.Println("		FAIL")
			fmt.Printf("			Expected: %v \n",x)
			fmt.Printf("			Recieved: %v \n",testBrd.winner)
		}
		if passed {numPassed++} else {numFailed++}
		testBrd = initializeBrd()

	}
	fmt.Printf("\n\n\n")
	fmt.Printf("SUMMARY:\n")
	fmt.Printf("	PASSED: %v", numPassed)
	fmt.Printf("	FAILED: %v", numFailed)
}

func retry_screen(board board){
	screen.Clear()
	screen.MoveTopLeft()
	board.print_board()
	fmt.Printf("Not a valid input (0-9)\n")
}	
func reset_screen(board board){
	screen.Clear()
	screen.MoveTopLeft()
	board.print_board()
}

func mn(){
	brd := initializeBrd()
	var err error = nil
    for true{
		reset_screen(brd)
		tn := generate_turn_from_int(int_move_input(retry_screen,brd))

		brd, err = brd.apply_turn(tn)
		if err != nil {
			fmt.Println("ERROR: ",err)
		}

		brd, err = brd.update_win_state()
		if err != nil {
			fmt.Println("ERROR: ",err)
		}
		if brd.winner != blank {
			brd.print_board()
			return 
		}
    }
}
//~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~


func sear









































































func search_tree_multi(b board, outchan chan int, channel_close_notification chan struct{}, min_goal bool){

	var err error
	b, err = b.update_win_state()
	if err != nil {
		fmt.Printf("ERROR: %v",err)
	}

	if b.winner != blank {
		outchan <- b.winner
	}

	legal_moves := make([]turn,0,9)
	subroutine_count := 0
	var tn turn

	for i:=0 ; i<9; i++ {
		tn = generate_turn_from_int(i)
		if b.is_legal(tn){
			subroutine_count++
			legal_moves = append(legal_moves,tn)
		}
	}

	inchan := make(chan int)
	//TODO: START GOROUTINES HERE

	var closest_to_goal int
	if min_goal {
		closest_to_goal = 1
	} else {
		closest_to_goal = -1
	}
	var rec_value int
	for num_recieved := 0 ; num_recieved < subroutine_count ; num_recieved++ {

		if min_goal{
			rec_value = <- inchan
			if rec_value < closest_to_goal{
			}

				//TODO: RETURN THE -1 IF THE RETURN CHANNEL IS OPEN
			} else {
				rec_value := <- 
				
			}



		} else {
			//TODO: RETURN THE -1 IF THE RETURN CHANNEL IS OPEN
		}
	}
	//TODO: CLOSE CHANNEL AND ACCOUNT FOR POSSIBLE CLOSURE OF SUPERIOR CHANNEL 
	close(inchan)
}
func dependant_close(inchan chan int, outchan chan_int){
	_, ok := <- outchan
	if !ok{
		close(inchan)
		return
	}
}
