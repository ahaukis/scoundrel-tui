# Scoundrel TUI

`scoundrel-tui` is a Go CLI implementation of Scoundrel, a rogue-like single-player card game. The player progresses through rooms - sets of cards dealt from the dungeon deck - until the dungeon is exhausted, or the player runs out of HP.

The Scoundrel card game was developed by Zach Gage and Kurt Bieg. For detailed rules and more information on the card game itself, see the [BoardGameGeek page](https://boardgamegeek.com/boardgame/191095/scoundrel) and the [rules document](http://stfj.net/art/2011/Scoundrel.pdf).

## Installation

For installation, use `go install`:

```
go install github.com/ahaukis/scoundrel-tui@latest
```

> Alternatively, you can clone the repository and run `make install`.

After installation, launch the game with `scoundrel-tui`.

## License

[MIT](LICENSE)
