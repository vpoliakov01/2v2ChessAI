export enum Color {
  Red = 'red',
  Blue = 'blue',
  Yellow = 'yellow',
  Green = 'green',
  Black = 'black',
  LightGray = 'light-gray',
  DarkGray = 'dark-gray',
}

export const colorCode = {
  [Color.Red]: '#bf3B43',
  [Color.Blue]: '#4185bf',
  [Color.Yellow]: '#c09526',
  [Color.Green]: '#4e9161',
  [Color.Black]: '#302e2b',
  [Color.LightGray]: '#dadada',
  [Color.DarkGray]: '#adadad',
}

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

export interface Move {
  from: Position;
  to: Position;
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

export interface MoveInfo extends Move {
  piece: Piece;
  capturedPiece: Piece | null;
}
