package main

import (
	"fmt"
	"errors"
	"github.com/inancgumus/screen"
)

const x = 1
const o = -1
const blank = 0

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

	fmt.Println("DEBUG: {")
	fmt.Println("	Rows:")
	fmt.Println("		",ws.row_wins[:])
	fmt.Println("	Cols:")
	fmt.Println("		",ws.col_wins[:])
	fmt.Println("	Diag:")
	fmt.Println("		",ws.diag_wins[:])
	fmt.Println("}")

	return slc_or(ws.row_wins[:]) || slc_or(ws.col_wins[:]) || slc_or(ws.diag_wins[:])
}

func initializeBrd() board{
	var out = board{}
	out.turn = x
	out.winner = 0

	return out
}

func persistent_input(options ...string) (out string) {
	fmt.Scanln(&out)
	for i := 0 ; i<len(options); i++{
		if options[i] == out {
			return out
		}
	}
	out = persistent_input(options...)
	return out
}

func generate_turn() turn {
	out := turn{}
	x_word_options := [3]string{"l","m","r"}
	y_word_options := [3]string{"b","m","t"}

	fmt.Println("Enter the x turn (l,m,r):")
	x_word := persistent_input(x_word_options[:]...)

	fmt.Println("Enter the y turn (b,m,t):")
	y_word := persistent_input(y_word_options[:]...)

	if x_word == "l"{
		out.x = 0
	} else if x_word == "m"{
		out.x = 1
	} else if x_word == "r"{
		out.x = 2
	}

	if y_word == "b"{
		out.y = 2
	} else if y_word == "m"{
	out.y = 1
	} else if y_word == "t"{
		out.y = 0
	}

	return out
}

func slc_sum(in []int) (sum int) {
	for _,e := range in {
		sum += e
	}
	return
}


func (b board) is_won() (winner int, err error){
	x_winSum := initializeWinSum() 
	o_winSum := initializeWinSum() 

	for col := 0 ; col<3 ; col++{
		for row := 0 ; row<3 ; row++{
			//Check Linear solutions
		 	if b.brd[row][col] == blank{
				x_winSum.row_wins[row] = false
				o_winSum.row_wins[row] = false
				x_winSum.col_wins[col] = false
				o_winSum.col_wins[col] = false
			} else if b.brd[col][row] == x{
				o_winSum.row_wins[row] = false
				o_winSum.col_wins[col] = false
			} else {
				x_winSum.row_wins[row] = false
				x_winSum.col_wins[col] = false
			}

			//Check Diagonals
			if 3-row == col {
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
	fmt.Println("x dbg:")
	x_win := x_winSum.compile()
	fmt.Println("o dbg:")
	o_win := o_winSum.compile()
	if x_win && o_win{
		return blank, errors.New("Both players have won")
	}
	if x_win {
		return x, nil
	}
	if o_win {
		return o, nil
	}
	return blank, nil
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
	if b.winner == 0{
	} else if b.winner == x{
		fmt.Println("x won")
	} else if b.winner == o{
		fmt.Println("o won")
	}
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
	winner, err := b.is_won() 
	if err != nil {
		return b, err
	}

	if winner != blank{
		return b, errors.New("Game has already been won")
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
	screen.Clear()
	brd := initializeBrd()
	brd.print_board()
	tn := generate_turn()
	var err error = nil
    for true{
		brd, err = brd.apply_turn(tn)
		if err != nil {
			fmt.Println("ERROR: ",err)
		}
		brd.print_board()
		tn = generate_turn()
		screen.Clear()
    }
}
