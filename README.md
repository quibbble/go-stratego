# Go-stratego

Go-stratego is a [Go](https://golang.org) implementation of the board game [Stratego](https://boardgamegeek.com/boardgame/1917/stratego). Please note that this repo only includes game logic and a basic API to interact with the game but does NOT include any form of GUI.

Check out [stratego.quibbble.com](stratego.quibbble.com) if you wish to view and play a live version of this game which utilizes this project along with a separate custom UI.

[![Quibbble Stratego](https://i.imgur.com/iekrcod.png)](stratego.quibbble.com)

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
