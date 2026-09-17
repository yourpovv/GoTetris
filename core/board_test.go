package core

import (
	"strings"
	"testing"
)

func boardFromStrings(t *testing.T, rows []string) *board {
	t.Helper()
	if len(rows) != boardHeight {
		t.Fatalf("need %d rows, got %d", boardHeight, len(rows))
	}
	b := &board{}
	for y, row := range rows {
		if len(row) != boardWidth {
			t.Fatalf("row %d: need %d cells, got %q", y, boardWidth, row)
		}
		for x, cell := range row {
			if cell == 'X' {
				b.cells[y][x] = T
			}
		}
	}
	return b
}

func boardStrings(b *board) []string {
	rows := make([]string, boardHeight)
	for y := 0; y < boardHeight; y++ {
		var sb strings.Builder
		for x := 0; x < boardWidth; x++ {
			if b.cells[y][x] == empty {
				sb.WriteByte('.')
			} else {
				sb.WriteByte('X')
			}
		}
		rows[y] = sb.String()
	}
	return rows
}

func emptyRow() string {
	return strings.Repeat(".", boardWidth)
}

func fullRow() string {
	return strings.Repeat("X", boardWidth)
}

func TestCollides(t *testing.T) {
	occupied := &board{}
	occupied.cells[5][5] = T

	block := []point{{0, 0}, {1, 0}}

	tests := []struct {
		name  string
		board *board
		cells []point
		pos   point
		want  bool
	}{
		{"free space", &board{}, block, point{3, 3}, false},
		{"left wall", &board{}, block, point{-1, 3}, true},
		{"right wall", &board{}, block, point{9, 3}, true},
		{"floor", &board{}, block, point{3, boardHeight}, true},
		{"above board allowed", &board{}, block, point{3, -1}, false},
		{"occupied cell", occupied, block, point{4, 5}, true},
		{"adjacent cell free", occupied, block, point{6, 5}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.board.collides(tt.cells, tt.pos); got != tt.want {
				t.Errorf("collides(%v, %v) = %v, want %v", tt.cells, tt.pos, got, tt.want)
			}
		})
	}
}

func TestClearLines(t *testing.T) {
	partial := "X........."

	tests := []struct {
		name      string
		rows      func() []string
		cleared   int
		wantBoard func() []string
	}{
		{
			name: "no full rows",
			rows: func() []string {
				rows := make([]string, boardHeight)
				for i := range rows {
					rows[i] = emptyRow()
				}
				rows[boardHeight-1] = partial
				return rows
			},
			cleared: 0,
			wantBoard: func() []string {
				rows := make([]string, boardHeight)
				for i := range rows {
					rows[i] = emptyRow()
				}
				rows[boardHeight-1] = partial
				return rows
			},
		},
		{
			name: "single line",
			rows: func() []string {
				rows := make([]string, boardHeight)
				for i := range rows {
					rows[i] = emptyRow()
				}
				rows[boardHeight-2] = partial
				rows[boardHeight-1] = fullRow()
				return rows
			},
			cleared: 1,
			wantBoard: func() []string {
				rows := make([]string, boardHeight)
				for i := range rows {
					rows[i] = emptyRow()
				}
				rows[boardHeight-1] = partial
				return rows
			},
		},
		{
			name: "tetris",
			rows: func() []string {
				rows := make([]string, boardHeight)
				for i := range rows {
					rows[i] = emptyRow()
				}
				rows[boardHeight-5] = partial
				for i := boardHeight - 4; i < boardHeight; i++ {
					rows[i] = fullRow()
				}
				return rows
			},
			cleared: 4,
			wantBoard: func() []string {
				rows := make([]string, boardHeight)
				for i := range rows {
					rows[i] = emptyRow()
				}
				rows[boardHeight-1] = partial
				return rows
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := boardFromStrings(t, tt.rows())
			if got := b.clearLines(); got != tt.cleared {
				t.Errorf("clearLines() = %d, want %d", got, tt.cleared)
			}
			for y, want := range tt.wantBoard() {
				if got := boardStrings(b)[y]; got != want {
					t.Errorf("row %d = %q, want %q", y, got, want)
				}
			}
		})
	}
}

func TestLock(t *testing.T) {
	b := &board{}
	p := newPiece(O, 4, boardHeight-2)
	b.lock(p)

	for _, c := range baseCells[O] {
		if b.cells[boardHeight-2+c.y][4+c.x] != O {
			t.Errorf("cell (%d, %d) not locked as O piece", 4+c.x, boardHeight-2+c.y)
		}
	}
}
