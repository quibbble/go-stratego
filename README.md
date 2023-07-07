# Go-stratego

Go-stratego is a [Go](https://golang.org) implementation of the board game [Stratego](https://boardgamegeek.com/boardgame/1917/stratego). Please note that this repo only includes game logic and a basic API to interact with the game but does NOT include any form of GUI.

## Usage

To play a game create a new Connect4 instance:
```go
builder := Builder{}
game, err := builder.Create(&bg.BoardGameOptions{
    Teams: []string{"TeamA", "TeamB"}, // must contain at least 2 and at most 2 teams
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
