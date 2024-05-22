package main

import (
	"bytes"
	"math"
	"math/rand"
	"time"
)

func LetsPlaySomeGames() bool {
	time.Sleep(1 * time.Second)
	prime := PrimeSieve()
	time.Sleep(1 * time.Second)
	life := GameOfLife()
	time.Sleep(1 * time.Second)
	go Pi(5000)
	ticTacToe := TicTacToe()
	if prime && life && ticTacToe {
		return true
	}
	return false
}

func PrimeSieve() bool {
	ch := make(chan int)
	go Generate(ch)
	for i := 0; i < 10; i++ {
		prime := <-ch
		ch1 := make(chan int)
		go Filter(ch, ch1, prime)
		ch = ch1
	}
	return true
}

func GameOfLife() bool {
	l := NewLife(40, 15)
	for i := 0; i < 300; i++ {
		l.Step()
		time.Sleep(time.Second / 30)
	}
	return true
}

func Pi(n int) float64 {
	ch := make(chan float64)
	for k := 0; k < n; k++ {
		go term(ch, float64(k))
	}
	f := 0.0
	for k := 0; k < n; k++ {
		f += <-ch
	}
	return f
}

func Generate(ch chan<- int) {
	for i := 2; ; i++ {
		ch <- i
	}
}

func Filter(in <-chan int, out chan<- int, prime int) {
	for {
		i := <-in
		if i%prime != 0 {
			out <- i
		}
	}
}

func term(ch chan float64, k float64) {
	ch <- 4 * math.Pow(-1, k) / (2*k + 1)
}

func NewField(w, h int) *Field {
	s := make([][]bool, h)
	for i := range s {
		s[i] = make([]bool, w)
	}
	return &Field{s: s, w: w, h: h}
}

func (f *Field) Set(x, y int, b bool) {
	f.s[y][x] = b
}

func (f *Field) Alive(x, y int) bool {
	x += f.w
	x %= f.w
	y += f.h
	y %= f.h
	return f.s[y][x]
}

func (f *Field) Next(x, y int) bool {
	alive := 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if (j != 0 || i != 0) && f.Alive(x+i, y+j) {
				alive++
			}
		}
	}
	return alive == 3 || alive == 2 && f.Alive(x, y)
}

func NewLife(w, h int) *Life {
	a := NewField(w, h)
	for i := 0; i < (w * h / 4); i++ {
		a.Set(rand.Intn(w), rand.Intn(h), true)
	}
	return &Life{
		a: a, b: NewField(w, h),
		w: w, h: h,
	}
}

func (l *Life) Step() {
	for y := 0; y < l.h; y++ {
		for x := 0; x < l.w; x++ {
			l.b.Set(x, y, l.a.Next(x, y))
		}
	}
	l.a, l.b = l.b, l.a
}

func (l *Life) String() string {
	var buf bytes.Buffer
	for y := 0; y < l.h; y++ {
		for x := 0; x < l.w; x++ {
			b := byte(' ')
			if l.a.Alive(x, y) {
				b = '*'
			}
			buf.WriteByte(b)
		}
		buf.WriteByte('\n')
	}
	return buf.String()
}

func TicTacToe() bool {
	board := NewBoard()
	player1 := NewPlayer("Player 1", "X")
	player2 := NewPlayer("Player 2", "O")

	for !board.IsGameOver() {
		currentPlayer := player1
		if board.GetTurn() == 2 {
			currentPlayer = player2
		}

		move := currentPlayer.GetMove(board)
		board.MakeMove(move)

		time.Sleep(300 * time.Millisecond)
	}

	return true
}



func NewBoard() *Board {
	return &Board{
		grid:        [3][3]string{},
		turn:        1,
		movesPlayed: 0,
		size:        3,
	}
}

func (b *Board) GetTurn() int {
	return b.turn
}

func (b *Board) MakeMove(move Move) {
	b.grid[move.Row][move.Column] = move.PlayerSymbol
	b.turn++
	b.movesPlayed++
}

func (b *Board) IsGameOver() bool {
	if b.checkForWin() || b.checkForDraw() {
		return true
	}
	return false
}

func (b *Board) checkForWin() bool {
	// Check for winning conditions
	if b.grid[0][0] == b.grid[0][1] && b.grid[0][1] == b.grid[0][2] && b.grid[0][0] != "" {
		return true
	}
	if b.grid[1][0] == b.grid[1][1] && b.grid[1][1] == b.grid[1][2] && b.grid[1][0] != "" {
		return true
	}
	if b.grid[2][0] == b.grid[2][1] && b.grid[2][1] == b.grid[2][2] && b.grid[2][0] != "" {
		return true
	}
	if b.grid[0][0] == b.grid[1][0] && b.grid[1][0] == b.grid[2][0] && b.grid[0][0] != "" {
		return true
	}
	if b.grid[0][1] == b.grid[1][1] && b.grid[1][1] == b.grid[2][1] && b.grid[0][1] != "" {
		return true
	}
	if b.grid[0][2] == b.grid[1][2] && b.grid[1][2] == b.grid[2][2] && b.grid[0][2] != "" {
		return true
	}
	if b.grid[0][0] == b.grid[1][1] && b.grid[1][1] == b.grid[2][2] && b.grid[0][0] != "" {
		return true
	}
	if b.grid[0][2] == b.grid[1][1] && b.grid[1][1] == b.grid[2][0] && b.grid[0][2] != "" {
		return true
	}
	return false
}

func (b *Board) checkForDraw() bool {
	return b.movesPlayed == b.size*b.size
}

type Move struct {
	Row          int
	Column       int
	PlayerSymbol string
}

type Player struct {
	Name         string
	PlayerSymbol string
}

func NewPlayer(name, symbol string) *Player {
	return &Player{
		Name:         name,
		PlayerSymbol: symbol,
	}
}

func (p *Player) GetMove(board *Board) Move {
	// Get all available moves
	availableMoves := []Move{}
	for row := 0; row < board.size; row++ {
		for col := 0; col < board.size; col++ {
			if board.grid[row][col] == "" {
				availableMoves = append(availableMoves, Move{Row: row, Column: col, PlayerSymbol: p.PlayerSymbol})
			}
		}
	}

	// Choose a random move from the available moves
	randomMove := availableMoves[rand.Intn(len(availableMoves))]

	return randomMove
}
