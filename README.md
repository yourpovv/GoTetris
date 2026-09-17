<div align="center">
  
# GoTetris

**Terminal tetris game built with Go, Bubbletea and Lipgloss.**

[![GitHub](https://img.shields.io/github/stars/yourpovv/GOOB?style=social)](https://github.com/yourpovv/GOOB)


</div>

## Requirements

- Go 1.21+

## Running it

```bash
go run .
```

or build it first:

```bash
go build -o goteTris.exe
.\goTetris.exe
```

## Controls

- Left/Right or A/D to move
- Down or S to soft drop
- Up, W or X to rotate clockwise, Z for counter-clockwise
- Space to hard drop
- P to pause
- Q to quit
- R to restart after game over

## How it works

Blocks fall on a 10x20 grid. Clear the lines to earn points based on your current level

| Clear  | Base Points |
| ------ | ----------: |
| Single |         100 |
| Double |         300 |
| Triple |         500 |
| Tetris |         800 |

Soft and hard drops also give bonus points. Every 10 lines cleared increases your level and speeds up the gravity. The game ends when a new piece has nowhere to go

## License

[MIT](LICENSE) © [YourPOVV](https://github.com/yourpovv)
