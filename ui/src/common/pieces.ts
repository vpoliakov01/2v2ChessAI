import { Color } from './colors';

export enum PieceType {
	Pawn = '',
	Rook = 'r',
	Knight = 'n',
	Bishop = 'b',
	Queen = 'q',
	King = 'k',
}

export const pieceName = {
	[PieceType.Pawn]: 'pawn',
	[PieceType.Rook]: 'rook',
	[PieceType.Knight]: 'knight',
	[PieceType.Bishop]: 'bishop',
	[PieceType.Queen]: 'queen',
	[PieceType.King]: 'king',
};

export const pieceFANCharacter = {
	[PieceType.Pawn]: '♙', // ♟ is less readable
	[PieceType.Rook]: '♜',
	[PieceType.Knight]: '♞',
	[PieceType.Bishop]: '♝',
	[PieceType.Queen]: '♛',
	[PieceType.King]: '♔', // ♚
};

export const pieceSANCharacter = {
	[PieceType.Pawn]: '',
	[PieceType.Rook]: 'R',
	[PieceType.Knight]: 'N',
	[PieceType.Bishop]: 'B',
	[PieceType.Queen]: 'Q',
	[PieceType.King]: 'K',
};

export interface Piece {
	type: PieceType;
	color: Color;
}

export const getPieceImage = (piece: Piece) => {
	const color = piece.color;
	const type = piece.type;
	return `/${color}_${pieceName[type]}.svg`;
};
