# 2v2ChessAI

What is 2v2 Chess?

A variant of chess with 4 players (1 at every side of the board), where the moves are done in a clockwise order and teammates are on opposite sides of the board, i.e. Red and Yellow vs Blue and Green.
![image](https://user-images.githubusercontent.com/53489500/168638482-0886ab3a-a565-452b-9a94-80c3531cb19b.png)

Why can't you play normal chess?

I can and I do! But 2v2 chess is more dynamic and requires good teamwork and communication, which can be very rewarding! Also, unlike for 2v2 chess, there are already thousands of engines for normal chess.

What stage is the project in?

Since it's been started only a few days ago, it is in the MVP stage. It plays generally sound moves, but more optimization / testing is to be done.
![image](https://user-images.githubusercontent.com/53489500/168645011-80c7446c-d72d-47ed-a834-99f5ee948fdb.png)
An example of a position reached by the engine. Pretty similar to the kind of positions reached by human players.

How does it work?
To pick a move, it uses negamax with alpha-beta pruning to arrive to the most favorable forced position at a specified depth. A rough evaluation of the player's pieces' strength is made at the final depth (based on the piece's position, progression of the game, and number of available moves). It uses multithreading by running the position evaluation on all availabe CPUs (GPU acceleration is planned for the future).
