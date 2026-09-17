package core

type point struct {
	x, y int
}

func (p point) add(q point) point {
	return point{p.x + q.x, p.y + q.y}
}

func (p point) equals(q point) bool {
	return p.x == q.x && p.y == q.y
}
