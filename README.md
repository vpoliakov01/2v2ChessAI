# 2v2ChessAI

## What is 2v2 Chess?

A variant of chess with 4 players (1 at every side of the board), where the moves are done in a clockwise order and teammates are on opposite sides of the board, i.e. Red and Yellow vs Blue and Green.
![image](https://user-images.githubusercontent.com/53489500/168638482-0886ab3a-a565-452b-9a94-80c3531cb19b.png)

## Why can't you play normal chess?

I can and I do! But 2v2 chess is more dynamic and requires good teamwork and communication, which can be very rewarding! Also, unlike for 2v2 chess, there are already thousands of engines for normal chess.

## What stage is the project in?

Since it's been started only a few days ago, it is in the MVP stage. It plays generally sound moves, but more optimization / testing is to be done.

![image](https://user-images.githubusercontent.com/53489500/169457551-9ab1c224-d676-4c19-ab04-6b76f1828257.png)

An example of a position reached by the engine playing against itself. Pretty similar to the kind of positions reached by human players.

## How does it work?

To pick a move, it uses negamax with alpha-beta pruning to arrive to the most favorable forced position at a specified depth. How favorable a position is is evaluated based on the team's pieces' positions, progression of the game, and number of available moves. It uses [multithreading by running the position evaluation on all availabe CPUs](https://github.com/vpoliakov01/2v2ChessAI/blob/dev/ai/ai.go#L78-L93) (GPU acceleration is planned for the future).

## What is the ELO estimate for this engine

On depth 5, it is around 1700-1800 ELO

## What are the main components that are worth checking out?
* [ai/ai.go](https://github.com/vpoliakov01/2v2ChessAI/blob/main/ai/ai.go)
* [game/game.go](https://github.com/vpoliakov01/2v2ChessAI/blob/main/ai/game.go)
* [game/board.go](https://github.com/vpoliakov01/2v2ChessAI/blob/main/ai/board.go)
* [game/piece.go](https://github.com/vpoliakov01/2v2ChessAI/blob/main/ai/piece.go)
* [game/ in general](https://github.com/vpoliakov01/2v2ChessAI/tree/main/game)
* [dev branch](https://github.com/vpoliakov01/2v2ChessAI/tree/dev)
* [~~PRs~~](https://github.com/vpoliakov01/2v2ChessAI/pulls?q=+)

Some more positions reached by the engine playing itself:

![image](https://user-images.githubusercontent.com/53489500/169458751-f20fe24b-2372-4ced-937b-75d575195e10.png)
![image](https://user-images.githubusercontent.com/53489500/169458772-539fa726-ffde-4f65-abb7-9e5271950d29.png)

## To play against the AI:
`go build -o cmd/ai cmd/main.go && ./cmd/ai`

## TODO:
### UI:
* Add toggle for game / analysis
* Add more settings

### Engine:
* Support castling
* Support forced calculation for checks
* Test with very sophisticated position evaluation
    * Fully tune piece position strength
    * Incorporate threat / liability

### Other:
* Dockerize (1 for ui, 1 for the engine)
* UI/Engine desync handling
* Update readme


FIX THE ISSUE WITH THE RECKLESS CAPTURES AND LACK OF TAKE BACKS

1. e2-e3 b10-c10 e13-e12 m5-l5
2. f1-b5 a6-b5 f14-a9 m8-l8
3. j1-i3 a8-a9 j13-j12 n6-j2
4. i1-j2 a5-c4 g13-g12 n7-h13
5. e1-f3 a9-f14 g14-h13 m10-l10
6. d2-d3 f14-f13 h13-h12 n9-i14
7. d3-c4 f13-g12 h12-i11 n10-l11
8. d1-d11 b5-k14 h14-h3 l11-k13
9. k2-k4 g12-h11 h3-h11 k13-l11
10. d11-b11 k14-h11 i11-h11 l11-j12
11. b11-a11 b8-c8 h11-h12 j12-i10
12. a11-a10 a7-b8 h12-g13 i10-g11
13. a10-a4 b6-c6 e14-f12 g11-e12
14. a4-b4 b8-a7 d13-e12 i14-d9
15. k4-l5 a7-a6 d14-d9 m4-l5
16. b4-a4 a6-b6 d9-b9 n4-i4
17. k1-k10 b6-c5 b9-b7 i4-i13
18. a4-a5 c5-d6 g13-h14 i13-i14
19. k10-c10 c8-d8 h14-i14 n5-l6
20. c10-c6 d6-c6 b7-b8 l6-k4
21. i3-k4 c6-c7 j14-i12 l5-k4
22. j2-d8 c7-b8 f12-d11 n11-n9
23. h2-h3 b8-c8 i14-j14 n8-n7
24. g1-d1 c8-b9 d11-e9 m6-l6
25. d1-d6 b9-a10 e9-f11 n7-n6
26. a5-a10