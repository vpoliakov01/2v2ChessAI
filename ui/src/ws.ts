import { Move } from './common';

export enum MessageType {
  GetBoardState = 'getBoardState',
  BoardState = 'boardState',
  GetMoves = 'getMoves',
  Moves = 'moves',
  PlayerMove = 'playerMove',
  EngineMove = 'engineMove',
  InvalidMove = 'invalidMove',
}

type PlayerMove = Move;

export interface BestMoveResponse {
  move: Move;
  score: number;
}

type MessageData = PlayerMove | Move[] | BestMoveResponse | null;

export class Message {
  constructor(public type: MessageType, public data: MessageData) {}

  json() {
    return JSON.stringify({
      type: this.type,
      data: this.data,
    });
  }
}
