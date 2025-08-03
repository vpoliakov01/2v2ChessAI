import { PGNMove } from '../common';

export enum MessageType {
  SetSettings = 'setSettings',
  GetAvailableMoves = 'getAvailableMoves',
  AvailableMoves = 'availableMoves',
  PlayerMove = 'playerMove',
  EngineMove = 'engineMove',
  InvalidMove = 'invalidMove',
  NewGame = 'newGame',
  SaveGame = 'saveGame',
  SaveGameResponse = 'saveGameResponse',
  LoadGame = 'loadGame',
  LoadGameResponse = 'loadGameResponse',
  SetCurrentMove = 'setCurrentMove',
}

export interface BestMoveResponse {
  move: PGNMove;
  score: number;
  time: number;
  evaluations: number;
}

export interface SaveGameResponse {
  pgn: string;
}

export interface LoadGameResponse {
  pastMoves: PGNMove[];
  currentMove: number;
}

export interface GameSettings {
  humanPlayers: number[];
  depth: number;
  captureDepth: number;
  evalLimit: number;
}

type MessageData = PGNMove | PGNMove[] | BestMoveResponse | SaveGameResponse | LoadGameResponse | GameSettings | string | number | null;

export class Message {
  constructor(public type: MessageType, public data: MessageData) { }

  json() {
    return JSON.stringify({
      type: this.type,
      data: this.data,
    });
  }
}
