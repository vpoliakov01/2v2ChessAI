import { useState, useEffect, useRef, useCallback, useReducer } from 'react';
import { Position, Move, movesEqual, PGNMove } from '../common';
import { ArrowProps } from '../components/Arrow';
import { MessageType, Message, BestMoveResponse, SaveGameResponse, LoadGameResponse, GameStateManager, GameSyncService, GameEndedResponse } from '../utils';
import { gameReducer, loadInitialState } from './gameReducer';

export function useBoardState() {
  const [state, dispatch] = useReducer(gameReducer, undefined, loadInitialState);
  const { board, activePlayer, allMoves, currentMove, availableMoves, score, pgn } = state;
  const moves = allMoves.slice(0, currentMove + 1);

  const [selectedSquare, setSelectedSquare] = useState<Position | null>(null);
  const [drawnArrows, setDrawnArrows] = useState<ArrowProps[]>(GameStateManager.load().drawnArrows);
  const [isDrawingArrow, setIsDrawingArrow] = useState<boolean>(false);
  const [arrowStart, setArrowStart] = useState<Position | null>(null);
  const [arrowEnd, setArrowEnd] = useState<Position | null>(null);

  const wsRef = useRef<WebSocket | null>(null);

  const sendMessage = useCallback((message: Message) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      console.debug(`Sending    ${message.type}`.padEnd(30), message);
      wsRef.current.send(message.json());
      return;
    }
    
    setTimeout(() => {
      if (wsRef.current?.readyState === WebSocket.OPEN) {
        console.debug(`Sending    ${message.type}`.padEnd(30), message);
        wsRef.current.send(message.json());
      } else {
        alert('No WebSocket connection');
      }
    }, 10);
  }, []);

  const setCurrentMove = useCallback((value: number) => {
    dispatch({ type: 'setCurrentMove', currentMove: value });
  }, []);

  const setPgn = useCallback((value: string) => {
    dispatch({ type: 'setPgn', pgn: value });
  }, []);

  const isValidMove = useCallback((move: Move): boolean => {
    return availableMoves.some(m => movesEqual(m, move));
  }, [availableMoves]);

  const movePiece = useCallback((move: Move, playerMove: boolean = false) => {
    if (playerMove && wsRef.current) {
      if (!isValidMove(move)) {
        return false;
      }
      sendMessage(new Message(MessageType.PlayerMove, move.toPGN()));
    }
    dispatch({ type: 'movePiece', move, playerMove });
    return true;
  }, [isValidMove, sendMessage]);

  // Arrow drawing
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
    setDrawnArrows([]);
  }, []);

  useEffect(() => {
    GameStateManager.save({ board, activePlayer, allMoves, currentMove, score, pgn, drawnArrows });
  }, [board, activePlayer, allMoves, currentMove, score, pgn, drawnArrows]);

  // WebSocket message handler.
  useEffect(() => {
    const ws = new WebSocket('ws://localhost:8080/ws');

    ws.onopen = () => {
      wsRef.current = ws;
      console.log('WS connected');
      GameSyncService.syncWithEngine(ws);
    };

    ws.onmessage = (event) => {
      const message = JSON.parse(event.data) as Message;
      console.debug(`Received   ${message.type}`.padEnd(30), message);

      switch (message.type) {
        case MessageType.AvailableMoves:
          dispatch({ type: 'setAvailableMoves', moves: (message.data as PGNMove[]).map(Move.fromPGN) });
          break;
        case MessageType.EngineMove:
          dispatch({ type: 'engineMove', moveData: message.data as BestMoveResponse });
          break;
        case MessageType.SaveGameResponse:
          dispatch({ type: 'setPgn', pgn: (message.data as SaveGameResponse).pgn });
          break;
        case MessageType.LoadGameResponse: {
          const loadData = message.data as LoadGameResponse;
          dispatch({ type: 'replayMoves', pastMoves: loadData.pastMoves.map(Move.fromPGN), currentMove: loadData.currentMove });
          break;
        }
        case MessageType.GameEnded: {
          const gameEndedData = message.data as GameEndedResponse;
          alert(`${gameEndedData.king} king has fallen! ${gameEndedData.winner} is victorious!`);
          break;
        }
        default:
          console.log('unknown message', message);
          break;
      }
    };

    ws.onclose = () => {
      if (wsRef.current === ws) {
        wsRef.current = null;
      }
      console.log('WS disconnected');
    };

    return () => ws.close();
  }, [dispatch]);

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
