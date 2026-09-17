package core

import "testing"

func cellSet(cells []point) map[point]bool {
	set := make(map[point]bool, len(cells))
	for _, c := range cells {
		set[c] = true
	}
	return set
}

func equalCellSets(a, b []point) bool {
	if len(a) != len(b) {
		return false
	}
	set := cellSet(a)
	for _, c := range b {
		if !set[c] {
			return false
		}
	}
	return true
}

func TestRotateFullCircle(t *testing.T) {
	for kind := I; kind <= L; kind++ {
		start := normalizeCells(append([]point(nil), baseCells[kind]...))
		rotated := append([]point(nil), start...)
		for i := 0; i < 4; i++ {
			p := &piece{cells: rotated}
			rotated = p.rotatedCells(true)
		}
		if !equalCellSets(rotated, start) {
			t.Errorf("kind %d: 4 clockwise rotations did not return to start: got %v, want %v", kind, rotated, start)
		}

		rotated = append([]point(nil), start...)
		for i := 0; i < 4; i++ {
			p := &piece{cells: rotated}
			rotated = p.rotatedCells(false)
		}
		if !equalCellSets(rotated, start) {
			t.Errorf("kind %d: 4 counter-clockwise rotations did not return to start: got %v, want %v", kind, rotated, start)
		}
	}
}

func TestRotateClockwise(t *testing.T) {
	tests := []struct {
		name      string
		kind      pieceKind
		clockwise bool
		want      []point
	}{
		{"T clockwise", T, true, []point{{1, 1}, {0, 0}, {0, 1}, {0, 2}}},
		{"S clockwise", S, true, []point{{1, 1}, {1, 2}, {0, 0}, {0, 1}}},
		{"L clockwise", L, true, []point{{1, 2}, {0, 0}, {0, 1}, {0, 2}}},
		{"I clockwise", I, true, []point{{0, 0}, {0, 1}, {0, 2}, {0, 3}}},
		{"T counter-clockwise", T, false, []point{{0, 1}, {1, 0}, {1, 1}, {1, 2}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &piece{cells: append([]point(nil), baseCells[tt.kind]...)}
			if got := p.rotatedCells(tt.clockwise); !equalCellSets(got, tt.want) {
				t.Errorf("rotatedCells(%v) = %v, want %v", tt.clockwise, got, tt.want)
			}
		})
	}
}
