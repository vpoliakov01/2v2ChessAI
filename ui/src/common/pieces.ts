import { Color } from './colors';

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
