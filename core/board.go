package core

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	boardWidth  = 10
	boardHeight = 20
)

type board struct {
	cells [boardHeight][boardWidth]pieceKind
}

func (b *board) collides(cells []point, pos point) bool {
	for _, c := range cells {
		x, y := pos.x+c.x, pos.y+c.y
		if x < 0 || x >= boardWidth || y >= boardHeight {
			return true
		}
		if y >= 0 && b.cells[y][x] != empty {
			return true
		}
	}
	return false
}

func (b *board) lock(p *piece) {
	for _, c := range p.cells {
		x, y := p.position.x+c.x, p.position.y+c.y
		if y >= 0 && y < boardHeight && x >= 0 && x < boardWidth {
			b.cells[y][x] = p.kind
		}
	}
}

func (b *board) clearLines() int {
	cleared := 0
	write := boardHeight - 1
	for read := boardHeight - 1; read >= 0; read-- {
		if b.isFullRow(read) {
			cleared++
			continue
		}
		if write != read {
			b.cells[write] = b.cells[read]
		}
		write--
	}
	for ; write >= 0; write-- {
		b.cells[write] = [boardWidth]pieceKind{}
	}
	return cleared
}

func (b *board) isFullRow(y int) bool {
	for x := 0; x < boardWidth; x++ {
		if b.cells[y][x] == empty {
			return false
		}
	}
	return true
}

func (b *board) render(active *piece, ghostPos point) string {
	var sb strings.Builder

	activeCells := map[point]pieceKind{}
	ghostCells := map[point]bool{}
	if active != nil {
		for _, c := range active.cells {
			activeCells[active.position.add(c)] = active.kind
		}
		for _, c := range active.cells {
			p := ghostPos.add(c)
			if _, ok := activeCells[p]; !ok {
				ghostCells[p] = true
			}
		}
	}

	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
	ghostStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(ghostColor))

	sb.WriteString(borderStyle.Render("╔"+strings.Repeat("═", boardWidth*2)+"╗") + "\n")

	for y := 0; y < boardHeight; y++ {
		sb.WriteString(borderStyle.Render("║"))
		for x := 0; x < boardWidth; x++ {
			p := point{x, y}
			switch kind, ok := activeCells[p]; {
			case ok:
				sb.WriteString(blockStyle(pieceColors[kind]).Render("[]"))
			case ghostCells[p]:
				sb.WriteString(ghostStyle.Render("[]"))
			case b.cells[y][x] != empty:
				sb.WriteString(blockStyle(pieceColors[b.cells[y][x]]).Render("[]"))
			default:
				sb.WriteString("  ")
			}
		}
		sb.WriteString(borderStyle.Render("║") + "\n")
	}

	sb.WriteString(borderStyle.Render("╚"+strings.Repeat("═", boardWidth*2)+"╝") + "\n")

	return sb.String()
}

func blockStyle(color string) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Bold(true)
}
