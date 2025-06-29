import { useState, useEffect, useRef, useCallback } from 'react';
import { Color, PlayerColors, Piece, PieceType, Position, MoveInfo, Move, movesEqual } from '../common';
import { MessageType, Message, BestMoveResponse } from '../ws';

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
          return [i, BOARD_SIZE - CORNER_SIZE - 1 - j];
        case Color.Green:
          return [BOARD_SIZE - CORNER_SIZE - 1 - j, BOARD_SIZE - 1 - i];
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
  const [selectedPiece, setSelectedPiece] = useState<Position | null>(null);
  const [moves, setMoves] = useState<MoveInfo[]>([]);
  const [availableMoves, setAvailableMoves] = useState<Move[]>([]);

  const wsRef = useRef<WebSocket | null>(null);

  const isValidMove = useCallback((move: Move): boolean => {
    return availableMoves.some(m => movesEqual(m, move));
  }, [availableMoves]);

  const movePiece = useCallback((move: Move, playerMove: boolean = true) => {
    const {from, to} = move;

    if (playerMove && wsRef.current) {
      if (!isValidMove(move)) {
        return false;
      }

      wsRef.current.send(new Message(
        MessageType.PlayerMove,
        move,
      ).json());
    }

    setMoves([...moves, {
      ...move,
      piece: board[from.row][from.col]!,
      capturedPiece: board[to.row][to.col] ?? null,
    }]);

    const newBoard = [...board.map(row => [...row])];
    newBoard[to.row][to.col] = board[from.row][from.col];
    newBoard[from.row][from.col] = null;
    setBoard(newBoard);

    setActivePlayer(PlayerColors[(PlayerColors.indexOf(activePlayer) + 1) % PlayerColors.length]);
    return true;
  }, [board, isValidMove, moves, activePlayer]);

  useEffect(() => {
    const ws = new WebSocket('ws://localhost:8080/ws');
    wsRef.current = ws;

    ws.onopen = () => {
      console.log('connected');
    };

    ws.onclose = () => {
      console.log('disconnected');
    };

    return () => {
      ws.close();
    };
  }, []);

  useEffect(() => {
    if (!wsRef.current) {
      return;
    }

    wsRef.current.onmessage = (event) => {
      console.log('message', event.data);
      const message = JSON.parse(event.data) as Message;
      switch (message.type) {
        case MessageType.Moves:
          setAvailableMoves([...(message.data as Move[])]);
          break;
        case MessageType.EngineMove:
          const response = message.data as BestMoveResponse;
          movePiece(response.move, false);
          break;
        default:
          console.log('unknown message', message);
          break;
      }
    };
  }, [movePiece]);

  return {
    board,
    activePlayer,
    moves,
    availableMoves,
    selectedPiece,
    movePiece,
    setSelectedPiece,
  };
} 