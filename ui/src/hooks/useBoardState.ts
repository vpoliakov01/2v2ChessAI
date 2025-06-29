import { useState } from 'react';
import { Color, PlayerColors, Piece, PieceType, Move } from '../common';

export type BoardPosition = (Piece | null | undefined)[][];

export const BOARD_SIZE = 14;
export const CORNER_SIZE = 3;

// Initialize the 14x14 board with cut corners
function createEmptyBoard(): BoardPosition {
  const board: BoardPosition = Array(BOARD_SIZE).fill(null)
    .map(() => Array(BOARD_SIZE).fill(null));

  // Cut out corners
  for (let i = 0; i < CORNER_SIZE; i++) {
    for (let j = 0; j < CORNER_SIZE; j++) {
      for (const k of [i, BOARD_SIZE - 1 - i]) {
        for (const l of [j, BOARD_SIZE - 1 - j]) {
          board[k][l] = undefined;
        }
      }
    }
  }

  return board;
}

function setPieces(board: BoardPosition): BoardPosition {
  const placePieces = (color: Color) => {
    const pieces: PieceType[][] = [
      [
        PieceType.Rook,
        PieceType.Knight,
        PieceType.Bishop,
        PieceType.Queen,
        PieceType.King,
        PieceType.Bishop,
        PieceType.Knight,
        PieceType.Rook
      ],
      Array(8).fill(PieceType.Pawn)
    ];

    const transformIndex = (color: Color, i: number, j: number) => {
      switch (color) {
        case Color.Yellow:
          return [i, j + CORNER_SIZE];
        case Color.Green:
          return [j + CORNER_SIZE, BOARD_SIZE - 1 - i];
        case Color.Red:
          return [BOARD_SIZE - 1 - i, j + CORNER_SIZE];
        case Color.Blue:
          return [j + CORNER_SIZE, i];
        default:
          throw new Error(`Invalid color: ${color}`);
      }
    }

    for (let i = 0; i < pieces.length; i++) {
      for (let j = 0; j < pieces[i].length; j++) {
        const piece = pieces[i][j];
        const [k, l] = transformIndex(color, i, j);
        board[k][l] = { type: piece, color };
      }
    }
  };

  for (const color of PlayerColors) {
    placePieces(color);
  }

  return board;
}

export function useBoardState() {
  const [board, setBoard] = useState<BoardPosition>(() => {
    const emptyBoard = createEmptyBoard();
    return setPieces(emptyBoard);
  });
  const [activePlayer, setActivePlayer] = useState<Color>(Color.Red);
  const [selectedPiece, setSelectedPiece] = useState<{row: number, col: number} | null>(null);
  const [moves, setMoves] = useState<Move[]>([]);

  const turns: Color[] = PlayerColors;

  const movePiece = (fromRow: number, fromCol: number, toRow: number, toCol: number) => {
    if (!isValidMove(fromRow, fromCol, toRow, toCol)) {
      return false;
    }

    setMoves([...moves, {
      from: {row: fromRow, col: fromCol},
      to: {row: toRow, col: toCol},
      piece: board[fromRow][fromCol]!,
      capturedPiece: board[toRow][toCol] ?? null,
    }]);

    const newBoard = [...board.map(row => [...row])];
    newBoard[toRow][toCol] = board[fromRow][fromCol];
    newBoard[fromRow][fromCol] = null;
    setBoard(newBoard);

    setActivePlayer(turns[(turns.indexOf(activePlayer) + 1) % turns.length]);
    return true;
  };

  const isValidMove = (fromRow: number, fromCol: number, toRow: number, toCol: number): boolean => {
    if (
      board?.[fromRow]?.[fromCol] === undefined ||
      board?.[toRow]?.[toCol] === undefined
    ) {
      return false;
    }

    const piece = board[fromRow][fromCol];
    if (!piece || piece.color !== activePlayer) {
      return false;
    }

    const targetPiece = board[toRow][toCol];
    if (targetPiece && targetPiece.color === piece.color) {
      return false;
    }

    // TODO: Add backend control.
    return true;
  };

  return {
    board,
    activePlayer,
    moves,
    selectedPiece,
    movePiece,
    setSelectedPiece,
  };
} 