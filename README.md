# ClockTowerAPI

This is the backend for a (incredibly still in progress!) version of the game [Blood on the Clock Tower](https://bloodontheclocktower.com/), best described as a social deception, role-playing game of Mafia, but it's way more complicated than that. The game *does* have [an open-source implementation](https://clocktower.online/) that's now deprecated, and [a newer version which requires payment](https://online.bloodontheclocktower.com). Ideally, this would be closer to the first version.

**This should not be used to play a game,** it's still very much under construction.

This repository is a Go backend intended for three main purposes:
- Executing game logic for the entire game (the "game" package)
- Writing information about the game's whereabouts in a PostgreSQL database (powered by [GORM](https://gorm.io/))
- Communicating with clients, whether playing the game or as the host (or, storyteller) with [Gin](https://gin-gonic.com/) for a REST API for non-real time communication, and [Melody](https://github.com/olahol/melody)-powered websockets

At the moment this is nowhere near complete, and generally this is just be a long-term side project for me (and my friends') own enjoyment.

## Build

To attempt to serve a backend:
1. Create an `.env` file with `PG_URL=[postgres URL]`, `PGUSER=[postgres username]`, and `PGPASS=[postgres password]`
2. Run `go install`
3. Run `go build ClockTowerAPI`
4. Run the main executable.

While tests are still being made, you can query both the API and the Websocket. Most websocket messages are of the form `{type: MESSAGE, message: { [any JSON here] }`


TODOs:
- [x] Reliable websocket communication
- [x] Role shuffling
- [ ] Client execution of requests (for night phases)
- [ ] Nominations & voting
- [ ] Trouble Brewing script
- [ ] More scripts! More characters!