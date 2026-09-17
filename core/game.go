package core

import (
	"fmt"
	"math/rand"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tickMsg time.Time

type gameState int

const (
	playing gameState = iota
	paused
	gameOver
)

const (
	baseGravity   = 800 * time.Millisecond
	gravityStep   = 70 * time.Millisecond
	minGravity    = 100 * time.Millisecond
	linesPerLevel = 10
)

var clearPoints = [...]int{0, 100, 300, 500, 800}

var rotationKicks = []point{{0, 0}, {-1, 0}, {1, 0}, {0, -1}}

type Game struct {
	board  *board
	active *piece
	next   pieceKind
	bag    []pieceKind
	score  int
	lines  int
	level  int
	state  gameState
}

func NewGame() *Game {
	g := &Game{board: &board{}, level: 1, state: playing}
	g.next = g.draw()
	g.spawn()
	return g
}

func (g *Game) Init() tea.Cmd {
	return tick(gravityInterval(g.level))
}

func (g *Game) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		switch key {
		case "ctrl+c", "q":
			return g, tea.Quit
		case "r":
			if g.state == gameOver {
				newG := NewGame()
				return newG, tick(gravityInterval(newG.level))
			}
		case "p":
			switch g.state {
			case playing:
				g.state = paused
			case paused:
				g.state = playing
			}
		}
		if g.state != playing {
			return g, nil
		}
		switch key {
		case "left", "a", "h":
			g.move(-1)
		case "right", "d", "l":
			g.move(1)
		case "down", "s", "j":
			g.softDrop()
		case "up", "w", "k", "x":
			g.rotate(true)
		case "z":
			g.rotate(false)
		case " ":
			g.hardDrop()
		}
	case tickMsg:
		if g.state == playing {
			g.step()
		}
		return g, tick(gravityInterval(g.level))
	}

	return g, nil
}

func (g *Game) View() string {
	title := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFD700")).
		Bold(true).
		MarginBottom(1).
		Render("🧱 GoTetris")

	stats := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00FFFF")).
		MarginBottom(1).
		Render(fmt.Sprintf("Score: %d | Level: %d | Lines: %d", g.score, g.level, g.lines))

	var ghost point
	if g.active != nil {
		ghost = g.ghostPosition()
	}
	gameBoard := g.board.render(g.active, ghost)

	side := lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Foreground(lipgloss.Color("#888888")).Render("Next:"),
		renderNext(g.next),
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			MarginTop(1).
			Render("Arrows/WASD: Move\nSpace: Hard drop\nZ/X: Rotate | P: Pause\nQ: Quit | R: Restart"),
	)

	body := lipgloss.JoinHorizontal(lipgloss.Top, gameBoard, "  "+side)

	var statusBar string
	switch g.state {
	case paused:
		statusBar = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFF00")).
			Bold(true).
			Render("⏸ PAUSED")
	case gameOver:
		statusBar = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true).
			Render("💀 GAME OVER | Press R to restart")
	default:
		statusBar = ""
	}

	view := fmt.Sprintf("%s\n%s\n\n%s", title, stats, body)

	if statusBar != "" {
		view = fmt.Sprintf("%s\n\n%s", view, statusBar)
	}

	return lipgloss.NewStyle().
		MarginLeft(2).
		MarginTop(1).
		Render(view)
}

func (g *Game) spawn() {
	cells := baseCells[g.next]
	width := 0
	for _, c := range cells {
		if c.x > width {
			width = c.x
		}
	}
	active := newPiece(g.next, (boardWidth-(width+1))/2, 0)
	g.next = g.draw()
	if g.board.collides(active.cells, active.position) {
		g.state = gameOver
		g.active = nil
		return
	}
	g.active = active
}

func (g *Game) draw() pieceKind {
	if len(g.bag) == 0 {
		kinds := []pieceKind{I, O, T, S, Z, J, L}
		for _, i := range rand.Perm(len(kinds)) {
			g.bag = append(g.bag, kinds[i])
		}
	}
	kind := g.bag[0]
	g.bag = g.bag[1:]
	return kind
}

func (g *Game) step() {
	down := g.active.position.add(point{0, 1})
	if g.board.collides(g.active.cells, down) {
		g.lockActive()
		return
	}
	g.active.position = down
}

func (g *Game) move(dx int) {
	next := g.active.position.add(point{dx, 0})
	if !g.board.collides(g.active.cells, next) {
		g.active.position = next
	}
}

func (g *Game) softDrop() {
	down := g.active.position.add(point{0, 1})
	if g.board.collides(g.active.cells, down) {
		g.lockActive()
		return
	}
	g.active.position = down
	g.score++
}

func (g *Game) hardDrop() {
	for {
		down := g.active.position.add(point{0, 1})
		if g.board.collides(g.active.cells, down) {
			break
		}
		g.active.position = down
		g.score += 2
	}
	g.lockActive()
}

func (g *Game) rotate(clockwise bool) {
	cells := g.active.rotatedCells(clockwise)
	for _, kick := range rotationKicks {
		pos := g.active.position.add(kick)
		if !g.board.collides(cells, pos) {
			g.active.cells = cells
			g.active.position = pos
			return
		}
	}
}

func (g *Game) lockActive() {
	g.board.lock(g.active)
	g.awardLines(g.board.clearLines())
	g.spawn()
}

func (g *Game) awardLines(cleared int) {
	g.lines += cleared
	g.score += clearPoints[cleared] * g.level
	g.level = g.lines/linesPerLevel + 1
}

func (g *Game) ghostPosition() point {
	pos := g.active.position
	for !g.board.collides(g.active.cells, pos.add(point{0, 1})) {
		pos = pos.add(point{0, 1})
	}
	return pos
}

func gravityInterval(level int) time.Duration {
	interval := baseGravity - time.Duration(level-1)*gravityStep
	if interval < minGravity {
		return minGravity
	}
	return interval
}

func renderNext(kind pieceKind) string {
	occupied := map[point]bool{}
	for _, c := range baseCells[kind] {
		occupied[c] = true
	}
	style := lipgloss.NewStyle().Foreground(lipgloss.Color(pieceColors[kind])).Bold(true)
	preview := ""
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			if occupied[point{x, y}] {
				preview += style.Render("[]")
			} else {
				preview += "  "
			}
		}
		preview += "\n"
	}
	return preview
}

func tick(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
