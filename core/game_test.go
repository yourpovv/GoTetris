package core

import (
	"testing"
	"time"
)

func TestAwardLines(t *testing.T) {
	tests := []struct {
		name      string
		lines     int
		level     int
		score     int
		cleared   int
		wantScore int
		wantLines int
		wantLevel int
	}{
		{"single at level 1", 0, 1, 0, 1, 100, 1, 1},
		{"double at level 1", 0, 1, 0, 2, 300, 2, 1},
		{"triple scored at level 2", 0, 2, 0, 3, 1000, 3, 1},
		{"tetris at level 3", 0, 3, 0, 4, 2400, 4, 1},
		{"level up at ten lines", 9, 1, 0, 1, 100, 10, 2},
		{"level up keeps score", 19, 2, 500, 2, 500 + 600, 21, 3},
		{"no lines no points", 5, 1, 250, 0, 250, 5, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Game{lines: tt.lines, level: tt.level, score: tt.score}
			g.awardLines(tt.cleared)
			if g.score != tt.wantScore {
				t.Errorf("score = %d, want %d", g.score, tt.wantScore)
			}
			if g.lines != tt.wantLines {
				t.Errorf("lines = %d, want %d", g.lines, tt.wantLines)
			}
			if g.level != tt.wantLevel {
				t.Errorf("level = %d, want %d", g.level, tt.wantLevel)
			}
		})
	}
}

func TestGravityInterval(t *testing.T) {
	tests := []struct {
		level int
		want  time.Duration
	}{
		{1, 800 * time.Millisecond},
		{2, 730 * time.Millisecond},
		{10, 170 * time.Millisecond},
		{11, 100 * time.Millisecond},
		{99, 100 * time.Millisecond},
	}

	for _, tt := range tests {
		if got := gravityInterval(tt.level); got != tt.want {
			t.Errorf("gravityInterval(%d) = %v, want %v", tt.level, got, tt.want)
		}
	}
}

func TestSpawnGameOver(t *testing.T) {
	g := &Game{board: &board{}, level: 1, state: playing}
	for x := 0; x < boardWidth; x++ {
		g.board.cells[0][x] = T
		g.board.cells[1][x] = T
	}
	g.bag = []pieceKind{I}
	g.next = g.draw()
	g.spawn()

	if g.state != gameOver {
		t.Errorf("state = %v, want gameOver after blocked spawn", g.state)
	}
	if g.active != nil {
		t.Errorf("active piece should be nil after blocked spawn")
	}
}

func TestHardDropLands(t *testing.T) {
	g := NewGame()
	if g.active == nil {
		t.Fatal("expected an active piece after NewGame")
	}
	g.hardDrop()

	bottomFilled := false
	for x := 0; x < boardWidth; x++ {
		if g.board.cells[boardHeight-1][x] != empty {
			bottomFilled = true
		}
	}
	if !bottomFilled {
		t.Error("expected locked cells on the bottom row after hard drop")
	}
	if g.score == 0 {
		t.Error("expected hard drop bonus points, got 0")
	}
}
