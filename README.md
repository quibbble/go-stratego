# Go-stratego

Go-stratego is a [Go](https://golang.org) implementation of the board game [Stratego](https://en.wikipedia.org/wiki/Stratego).

Check out [stratego.quibbble.com](https://stratego.quibbble.com) to play a live version of this game. This website utilizes [stratego](https://github.com/quibbble/stratego) frontend code, [go-stratego](https://github.com/quibbble/go-stratego) game logic, and [go-quibbble](https://github.com/quibbble/go-quibbble) server logic.

[![Quibbble Stratego](https://raw.githubusercontent.com/quibbble/stratego/main/screenshot.png)](https://stratego.quibbble.com)

## Usage

To play a game create a new Stratego instance:
```go
builder := Builder{}
game, err := builder.Create(&bg.BoardGameOptions{
    Teams: []string{"TeamA", "TeamB"}, // must contain at least 2 and at most 2 teams
})
```

To reoganize your board and switch units before the game starts do the following action:
```go
err := game.Do(&bg.BoardGameAction{
    Team: "TeamA",
    ActionType: "SwitchUnits",
    MoreDetails: SwitchUnitsActionDetails{
        UnitRow: 0,
        UnitColumn: 0,
        SwitchUnitRow: 1,
        SwitchUnitColumn: 0,
    },
})
```

To move a unit do the following action:
```go
err := game.Do(&bg.BoardGameAction{
    Team: "TeamA",
    ActionType: "MoveUnit",
    MoreDetails: MoveUnitActionDetails{
        UnitRow: 0,
        UnitColumn: 0,
        MoveRow: 1,
        MoveColumn: 0,
    },
})
```

To get the current state of the game call the following:
```go
snapshot, err := game.GetSnapshot("TeamA")
```
