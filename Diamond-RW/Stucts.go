package main

type Win32Process struct {
	Name string
}

type Field struct {
	s    [][]bool
	w, h int
}

type Life struct {
	a, b *Field
	w, h int
}

type Board struct {
	grid        [3][3]string
	turn        int
	movesPlayed int
	size        int
}
