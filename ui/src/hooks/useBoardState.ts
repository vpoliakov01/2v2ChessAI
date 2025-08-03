import { useState, useEffect, useRef, useCallback } from 'react';
import { Color, PlayerColors, PieceType, Position, MoveInfo, Move, movesEqual, movesToPGN, BOARD_SIZE, CORNER_SIZE, PGNMove, formatNumber } from '../common';
import { ArrowProps } from '../components/Arrow';
import { BoardStateStorage, BoardPosition, SavedBoardState, MessageType, Message, BestMoveResponse, SaveGameResponse, LoadGameResponse, loadSettingsFromStorage } from '../utils';

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

  return setPieces(board);
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



function getDefaultBoardState(): SavedBoardState {
  return {
    board: createEmptyBoard(),
    activePlayer: Color.Red,
    allMoves: [],
    currentMove: -1,
    score: 0,
    pgn: '',
    drawnArrows: [],
  };
}

function loadBoardStateFromStorage(): SavedBoardState {
  const loadedState = BoardStateStorage.load();
  return loadedState || getDefaultBoardState();
}

function saveBoardStateToStorage(state: SavedBoardState) {
  BoardStateStorage.store(state);
}

export function useBoardState() {
  // Load initial state from localStorage
  const savedState = loadBoardStateFromStorage();

  const [board, setBoard] = useState<BoardPosition>(savedState.board);
  const [activePlayer, setActivePlayer] = useState<Color>(savedState.activePlayer);
  const [selectedSquare, setSelectedSquare] = useState<Position | null>(null);
  const [allMoves, setAllMoves] = useState<MoveInfo[]>(savedState.allMoves);
  const [currentMove, setCurrentMove] = useState<number>(savedState.currentMove);
  const [availableMoves, setAvailableMoves] = useState<Move[]>([]);
  const [score, setScore] = useState<number>(savedState.score);
  const [pgn, setPgn] = useState<string>(savedState.pgn);
  const [drawnArrows, setDrawnArrows] = useState<ArrowProps[]>(savedState.drawnArrows);
  const [isDrawingArrow, setIsDrawingArrow] = useState<boolean>(false);
  const [arrowStart, setArrowStart] = useState<Position | null>(null);
  const [arrowEnd, setArrowEnd] = useState<Position | null>(null);
  const moves = allMoves.slice(0, currentMove + 1);

  const wsRef = useRef<WebSocket | null>(null);

  const sendMessage = useCallback((message: Message) => {
    if (wsRef.current) {
      wsRef.current.send(message.json());
    } else {
      console.log('No WebSocket connection');
    }
  }, [wsRef]);

  function resetBoard() {
    const defaultState = getDefaultBoardState();
    setBoard(defaultState.board);
    setActivePlayer(defaultState.activePlayer);
    setSelectedSquare(null);
    setAvailableMoves([]);
    setPgn(defaultState.pgn);
    setScore(defaultState.score);
    setCurrentMove(defaultState.currentMove);
    setDrawnArrows(defaultState.drawnArrows);
    setIsDrawingArrow(false);
    setArrowStart(null);
    setArrowEnd(null);
    // Do not setAllMoves.

    // Clear localStorage when resetting board
    BoardStateStorage.clear();
  }

  const replayMoves = useCallback((pastMoves: Move[], currentMove: number) => {
    resetBoard();
    const newMoves = [];
    let newBoard = createEmptyBoard();

    const lastMoveIndex = Math.min(currentMove, pastMoves.length - 1);

    for (let i = 0; i < pastMoves.length; i++) {
      const move = pastMoves[i];
      const { from, to } = move;
      newMoves.push(new MoveInfo(from, to, newBoard[from.row][from.col]!, newBoard[to.row][to.col] ?? null));
      newBoard[to.row][to.col] = newBoard[from.row][from.col];
      newBoard[from.row][from.col] = null;

      if (i === lastMoveIndex) { // Set the board and proceed with move calculation.
        setBoard(newBoard);
        newBoard = [...newBoard.map(row => [...row])];
      }
    }
    setAllMoves(newMoves);
    setActivePlayer(PlayerColors[(lastMoveIndex + 1) % PlayerColors.length]);
    setCurrentMove(lastMoveIndex);
  }, []);

  const isValidMove = useCallback((move: Move): boolean => {
    return availableMoves.some(m => movesEqual(m, move));
  }, [availableMoves]);

  const movePiece = useCallback((move: Move, playerMove: boolean = false) => {
    const { from, to } = move;

    if (playerMove && wsRef.current) {
      if (!isValidMove(move)) {
        return false;
      }

      sendMessage(new Message(MessageType.PlayerMove, move.toPGN()));
    }

    setAllMoves([...moves, new MoveInfo(from, to, board[from.row][from.col]!, board[to.row][to.col] ?? null)]);

    if (currentMove < moves.length - 1 && !playerMove) { // Engine move during game review.
      return true;
    }

    setCurrentMove(currentMove + 1);

    const newBoard = [...board.map(row => [...row])];
    newBoard[to.row][to.col] = newBoard[from.row][from.col];
    newBoard[from.row][from.col] = null;
    setBoard(newBoard);

    setActivePlayer(PlayerColors[(PlayerColors.indexOf(activePlayer) + 1) % PlayerColors.length]);
    return true;
  }, [board, isValidMove, moves, activePlayer, currentMove, sendMessage]);

  // Arrow drawing functions
  const handleSquareRightMouseDown = useCallback((position: Position) => {
    setIsDrawingArrow(true);
    setArrowStart(position);
    setArrowEnd(position);
  }, []);

  const handleSquareMouseEnter = useCallback((position: Position) => {
    if (isDrawingArrow) {
      setArrowEnd(position);
    }
  }, [isDrawingArrow]);

  const handleSquareRightMouseUp = useCallback((position: Position) => {
    if (!isDrawingArrow || !arrowStart) {
      return;
    }

    const newArrow: ArrowProps = {
      move: new Move(arrowStart, position),
      color: activePlayer
    };

    // If the new arrow is the same as any existing arrow,
    // erase the old one instead of adding a new one.
    setDrawnArrows(arrows => {
      const existingIndex = arrows.findIndex(arrow => movesEqual(arrow.move, newArrow.move));

      if (existingIndex === -1) {
        return [...arrows, newArrow];
      }

      return arrows.filter((_, index) => index !== existingIndex);
    });

    setIsDrawingArrow(false);
    setArrowStart(null);
    setArrowEnd(null);
  }, [isDrawingArrow, arrowStart, activePlayer]);

  const handleSquareLeftClick = useCallback(() => {
    // Clear all arrows on left click
    setDrawnArrows([]);
  }, []);

  // Save board state to localStorage whenever relevant state changes
  useEffect(() => {
    const currentBoardState: SavedBoardState = {
      board,
      activePlayer,
      allMoves,
      currentMove,
      score,
      pgn,
      drawnArrows,
    };
    saveBoardStateToStorage(currentBoardState);
  }, [board, activePlayer, allMoves, currentMove, score, pgn, drawnArrows]);

  useEffect(() => {
    const ws = new WebSocket('ws://localhost:8080/ws');
    wsRef.current = ws;

    ws.onopen = () => {
      console.log('connected');
      console.log(
        'move'.padEnd(8),
        'time'.padStart(6),
        'score'.padStart(9),
        'evals'.padStart(8),
        'avg'.padStart(6),
      );


      if (ws.readyState === WebSocket.OPEN) {
        const settings = loadSettingsFromStorage();
        const engineSettings = {
          ...settings,
          evalLimit: settings.evalLimit * 1000, // Convert to milliseconds for engine
        };

        ws.send(new Message(MessageType.SetSettings, engineSettings).json());

        const savedState = loadBoardStateFromStorage();
        const pgnToLoad = savedState.pgn || movesToPGN(savedState.allMoves);
        ws.send(new Message(MessageType.LoadGame, pgnToLoad).json());

        if (savedState.currentMove !== savedState.allMoves.length - 1) {
          ws.send(new Message(MessageType.SetCurrentMove, savedState.currentMove).json());
        }
      }

      console.log('Engine state synchronized with localStorage data');
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
      const message = JSON.parse(event.data) as Message;
      switch (message.type) {
        case MessageType.AvailableMoves:
          setAvailableMoves((message.data as PGNMove[]).map(Move.fromPGN));
          break;
        case MessageType.EngineMove:
          const moveData = message.data as BestMoveResponse;
          console.log(
            moveData.move.padEnd(8),
            formatNumber(moveData.time, 3, 2, 's'),
            formatNumber(moveData.score, 5, 2),
            formatNumber(moveData.evaluations / 1000, 5, 2, 'k'),
            formatNumber(moveData.time / moveData.evaluations * 1e6, 4, 0, 'μs'),
          );
          movePiece(Move.fromPGN(moveData.move));
          setScore(moveData.score);
          break;
        case MessageType.SaveGameResponse:
          const saveData = message.data as SaveGameResponse;
          setPgn(saveData.pgn);
          break;
        case MessageType.LoadGameResponse:
          const loadData = message.data as LoadGameResponse;
          replayMoves(loadData.pastMoves.map(Move.fromPGN), loadData.currentMove);
          break;
        default:
          console.log('unknown message', message);
          break;
      }
    };
  }, [movePiece, setScore, currentMove, replayMoves]);

  return {
    activePlayer,
    allMoves,
    availableMoves,
    board,
    currentMove,
    moves,
    pgn,
    score,
    selectedSquare,
    drawnArrows,
    isDrawingArrow,
    arrowStart,
    arrowEnd,
    setCurrentMove,
    movePiece,
    setPgn,
    setSelectedSquare,
    sendMessage,
    handleSquareRightMouseDown,
    handleSquareMouseEnter,
    handleSquareRightMouseUp,
    handleSquareLeftClick,
  };
} 
