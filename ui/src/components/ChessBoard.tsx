import React from 'react';

import { BOARD_SIZE, Color, Move, movesEqual, positionsEqual } from '../common';
import { Square, ScoreDisplay, PlayerIndicator } from './ChessBoardElements';
import { useBoardStateContext } from '../context/BoardStateContext';

export function ChessBoard() {
  const { board, activePlayer, moves, availableMoves, selectedSquare, score, movePiece, setSelectedSquare } = useBoardStateContext();

  const handleSquareClick = (row: number, col: number) => {
    const newPosition = {row, col};
    if (board[row][col] === undefined) return;

    if (selectedSquare) {
      if (movePiece(new Move(selectedSquare, newPosition), true)) {
        setSelectedSquare(null);
      } else if (board[row][col]?.color === activePlayer) {
        setSelectedSquare(newPosition);
      }
    } else if (board[row][col]?.color === activePlayer) {
      setSelectedSquare(newPosition);
    }
  };

  let higlightedSquares: {row: number, col: number, color: Color}[] = [];
  for (let i = moves.length - 1; i >= 0 && i > moves.length - 5; i--) {
    const move = moves[i];
    if (move?.piece) {
      higlightedSquares.push({ ...move.from, color: move.piece.color});
      higlightedSquares.push({ ...move.to, color: move.piece.color});
    }
  }
  if (selectedSquare) {
    higlightedSquares.push({ ...selectedSquare, color: activePlayer });
    higlightedSquares.push(...availableMoves.filter(m => positionsEqual(m.from, selectedSquare)).map(m => ({ ...m.to, color: activePlayer })));
  }

  return (
    <div className="board-container" style={{
      boxSizing: 'border-box',
      padding: 10,
      position: 'relative',
      height: '100%',
      display: 'flex',
      flexDirection: 'row',
    }}>
      <ScoreDisplay score={score} />
      <div className="board-inner-container"style={{
        position: 'relative',
        width: '100%',
        height: '100%',
      }}>
        <div className="board" style={{
          position: 'relative',
          width: '100%',
          height: '100%',
        }}>
          {Array(BOARD_SIZE).fill(null).map((_, row) => (
            <div className="row" key={row} style={{ 
              display: 'flex',
              height: `${100 / BOARD_SIZE}%`
            }}>
              {Array(BOARD_SIZE).fill(null).map((_, col) => (
                <Square
                  key={`${row}-${col}`}
                  isPlayable={board[row][col] !== undefined}
                  isLight={(row + col) % 2 === 0}
                  piece={board[row][col]}
                  higlighted={higlightedSquares.find(square => square.row === row && square.col === col)?.color ?? null}
                  possibleMove={!!selectedSquare && availableMoves.some(m => movesEqual(m, new Move(selectedSquare, {row, col})))}
                  onClick={() => handleSquareClick(row, col)}
                />
              ))}
            </div>
          ))}
          <div className="center-marker" style={{
            border: '4px solid #666',
            borderRadius: '50%',
            left: 'calc(50%)',
            opacity: '0.5',
            position: 'absolute',
            top: '50%',
            transform: 'translate(-50%, -50%)',
          }} />
        </div>
        <PlayerIndicator color={activePlayer} />
        </div>
    </div>
  );
}
