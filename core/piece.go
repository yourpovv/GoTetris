package core

type pieceKind int

const (
	empty pieceKind = iota
	I
	O
	T
	S
	Z
	J
	L
)

// baseCells hold the spawn orientation for each piece (0 is empty)
var baseCells = [8][]point{
	{},
	{{0, 1}, {1, 1}, {2, 1}, {3, 1}},
	{{0, 0}, {1, 0}, {0, 1}, {1, 1}},
	{{1, 0}, {0, 1}, {1, 1}, {2, 1}},
	{{1, 0}, {2, 0}, {0, 1}, {1, 1}},
	{{0, 0}, {1, 0}, {1, 1}, {2, 1}},
	{{0, 0}, {0, 1}, {1, 1}, {2, 1}},
	{{2, 0}, {0, 1}, {1, 1}, {2, 1}},
}

var pieceColors = [8]string{
	"",
	"#00F0F0",
	"#F0F000",
	"#A000F0",
	"#00F000",
	"#F00000",
	"#0000F0",
	"#F0A000",
}

const ghostColor = "#555555"

type piece struct {
	kind     pieceKind
	cells    []point
	position point
}

func newPiece(kind pieceKind, x, y int) *piece {
	cells := make([]point, len(baseCells[kind]))
	copy(cells, baseCells[kind])
	return &piece{kind: kind, cells: cells, position: point{x, y}}
}

func (p *piece) rotatedCells(clockwise bool) []point {
	rotated := make([]point, len(p.cells))
	for i, c := range p.cells {
		if clockwise {
			rotated[i] = point{-c.y, c.x}
		} else {
			rotated[i] = point{c.y, -c.x}
		}
	}
	return normalizeCells(rotated)
}

func normalizeCells(cells []point) []point {
	minX, minY := cells[0].x, cells[0].y
	for _, c := range cells[1:] {
		if c.x < minX {
			minX = c.x
		}
		if c.y < minY {
			minY = c.y
		}
	}
	shifted := make([]point, len(cells))
	for i, c := range cells {
		shifted[i] = point{c.x - minX, c.y - minY}
	}
	return shifted
}
