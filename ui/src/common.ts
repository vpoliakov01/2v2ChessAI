export const BOARD_SIZE = 14;
export const CORNER_SIZE = 3;

export enum Color {
  Red = 'red',
  Blue = 'blue',
  Yellow = 'yellow',
  Green = 'green',
  Black = 'black',
  LightGray = 'light-gray',
  Gray = 'gray',
  DarkGray = 'dark-gray',
}

export const colorCode = (color: Color) => `var(--color-${color})`;

export const PlayerColors = [Color.Red, Color.Blue, Color.Yellow, Color.Green];

export enum PieceType {
  Pawn = 'p',
  Rook = 'r',
  Knight = 'n',
  Bishop = 'b',
  Queen = 'q',
  King = 'k'
}

export const pieceName = {
  [PieceType.Pawn]: 'pawn',
  [PieceType.Rook]: 'rook',
  [PieceType.Knight]: 'knight',
  [PieceType.Bishop]: 'bishop',
  [PieceType.Queen]: 'queen',
  [PieceType.King]: 'king',
}

export interface Piece {
  type: PieceType;
  color: Color;
}

export interface Position {
  row: number;
  col: number;
}

export type PGNMove = string;

export class Move {
  constructor(public from: Position, public to: Position) {}

  static fromPGN(pgn: PGNMove): Move {
    const [from, to] = pgn.split('-').map(pos => ({
      col: pos.charCodeAt(0) - 'a'.charCodeAt(0),
      row: BOARD_SIZE - parseInt(pos.slice(1)),
    }));
    return new Move(from, to);
  }

  toPGN(): PGNMove {
    const [from, to] = [this.from, this.to].map(pos => `${String.fromCharCode(pos.col + 'a'.charCodeAt(0))}${BOARD_SIZE - pos.row}`);
    return `${from}-${to}`;
  }
}

export class MoveInfo extends Move {
  constructor(public from: Position, public to: Position, public piece: Piece, public capturedPiece: Piece | null) {
    super(from, to);
  }
}


export function positionsEqual(a: Position, b: Position): boolean {
  return a.row === b.row && a.col === b.col;
}

export function movesEqual(a: Move, b: Move): boolean {
  return a.from.row === b.from.row &&
         a.from.col === b.from.col && 
         a.to.row === b.to.row && 
         a.to.col === b.to.col;
}

export function movesToPGN(moves: Move[]): string {
  let pgn = "";

  for (let i = 0; i < moves.length; i += 4) {
    if (i > 0 && i % 4 === 0) {
      pgn += "\n";
    }
    pgn += `${i / 4 + 1}.`;
    for (let j = 0; j < 4 && i + j < moves.length; j++) {
      pgn += ` ${moves[i + j].toPGN()}`;
    }
  }

  return pgn;
}
