# 2v2ChessAI

## What is 2v2 Chess?

A variant of chess with 4 players (1 at every side of the board), where the moves are done in a clockwise order and teammates are on opposite sides of the board, i.e. Red and Yellow vs Blue and Green.
![image](https://user-images.githubusercontent.com/53489500/168638482-0886ab3a-a565-452b-9a94-80c3531cb19b.png)

## Why can't you play normal chess?

I can and I do! But 2v2 chess is more dynamic and requires good teamwork and communication, which can be very rewarding! Also, unlike for 2v2 chess, there are already thousands of engines for normal chess.

## What stage is the project in?

Since it's been started only a few days ago, it is in the MVP stage. It plays generally sound moves, but more optimization / testing is to be done.
![image](https://user-images.githubusercontent.com/53489500/168727993-bc9eb1a6-4c32-4994-8278-fb4cd57bccf5.png)
An example of a position reached by the engine playing against itself. Pretty similar to the kind of positions reached by human players.

## How does it work?

To pick a move, it uses negamax with alpha-beta pruning to arrive to the most favorable forced position at a specified depth. How favorable a position is is evaluated based on the team's pieces' positions, progression of the game, and number of available moves. It uses [multithreading by running the position evaluation on all availabe CPUs](https://github.com/vpoliakov01/2v2ChessAI/blob/dev/ai/ai.go#L82-L97) (GPU acceleration is planned for the future).

## What are the main components that are worth checking out?
* [ai/ai.go](https://github.com/vpoliakov01/2v2ChessAI/blob/main/ai/ai.go)
* [game/game.go](https://github.com/vpoliakov01/2v2ChessAI/blob/main/ai/game.go)
* [game/board.go](https://github.com/vpoliakov01/2v2ChessAI/blob/main/ai/board.go)
* [game/piece.go](https://github.com/vpoliakov01/2v2ChessAI/blob/main/ai/piece.go)
* [game/ in general](https://github.com/vpoliakov01/2v2ChessAI/tree/main/game)
* [dev branch](https://github.com/vpoliakov01/2v2ChessAI/tree/dev)
* [~~PRs~~](https://github.com/vpoliakov01/2v2ChessAI/pulls?q=+)

Some more positions reached by the engine playing itself:

![image](https://user-images.githubusercontent.com/53489500/168729637-f39da27a-744d-4229-9807-efcd3c516a0c.png)
![image](https://user-images.githubusercontent.com/53489500/168729546-68150198-a880-42b3-b38d-a92300a6f5b2.png)
![image](https://user-images.githubusercontent.com/53489500/168729819-535f804d-3136-4240-95c9-c1947319d8fa.png)