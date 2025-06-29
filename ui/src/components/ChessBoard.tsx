import React from 'react';

import { Color, colorCode, Piece, pieceName } from '../common';
import { useBoardState, BOARD_SIZE, CORNER_SIZE } from '../hooks/useBoardState';

interface SquareProps {
  isPlayable: boolean;
  isLight: boolean;
  piece: Piece | null | undefined;
  higlighted: Color | null;
  onClick: () => void;
}

function Square({ isPlayable, isLight, piece, higlighted, onClick }: SquareProps) {
  if (!isPlayable) {
    return <div style={{
      aspectRatio: '1',
      backgroundColor: colorCode[Color.Black],
      border: `1px solid ${colorCode[Color.Black]}`,
      width: '100%',
    }} />;
  }

  let backgroundColor = isLight ? colorCode[Color.LightGray] : colorCode[Color.DarkGray];
  if (higlighted) {
    backgroundColor = `color-mix(in srgb, ${colorCode[higlighted]} 45%, ${backgroundColor})`;
  }

  const getPieceImage = (piece: Piece) => {
    const color = piece.color;
    const type = piece.type;
    return `/${color}_${pieceName[type]}.svg`;
  };

  return (
    <div
      style={{
        alignItems: 'center',
        aspectRatio: '1',
        backgroundColor,
        cursor: piece ? 'pointer' : 'default',
        display: 'flex',
        justifyContent: 'center',
        position: 'relative',
        userSelect: 'none',
        width: '100%',
      }}
      onClick={onClick}
    >
      {piece && (
        <img 
          src={getPieceImage(piece)} 
          alt={`${piece.color} ${piece.type}`}
          style={{
            pointerEvents: 'none',
            position: 'absolute',
          }}
        />
      )}
    </div>
  );
}

export function PlayerIndicator({ color }: { color: Color }) {
  const getPlayerIndicatorStyle = (color: Color): React.CSSProperties => {
    const offsetLength = `calc(${CORNER_SIZE / BOARD_SIZE} * 100%)`;
    const lineWidth = '3px';
    const offset = 6;

    switch (color) {
      case Color.Red:
      case Color.Yellow:
        return {
          backgroundColor: colorCode[color],
          height: lineWidth,
          left: offsetLength,
          position: 'absolute',
          right: offsetLength,
          top: color === Color.Yellow ? -offset : `calc(100% - ${lineWidth} + ${offset}px)`,
        };
      case Color.Blue:
      case Color.Green:
        return {
          backgroundColor: colorCode[color],
          bottom: offsetLength,
          left: color === Color.Blue ? -offset : `calc(100% - ${lineWidth} + ${offset}px)`,
          position: 'absolute',
          top: offsetLength,
          width: lineWidth,
        };
      default:
        return {};
    }
  }
  return <div className="player-indicator" style={getPlayerIndicatorStyle(color)} />;
}

export function ChessBoard() {
  const { board, activePlayer, moves, selectedPiece, movePiece, setSelectedPiece } = useBoardState();

  const handleSquareClick = (row: number, col: number) => {
    if (board[row][col] === undefined) return;

    if (selectedPiece) {
      if (movePiece(selectedPiece.row, selectedPiece.col, row, col)) {
        setSelectedPiece(null);
      } else if (board[row][col]?.color === activePlayer) {
        setSelectedPiece({ row, col });
      }
    } else if (board[row][col]?.color === activePlayer) {
      setSelectedPiece({ row, col });
    }
  };

  let higlightedSquares: {row: number, col: number, color: Color}[] = [];
  for (let i = moves.length - 1; i >= 0 && i > moves.length - 5; i--) {
    higlightedSquares.push({ ...moves[i].from, color: moves[i].piece.color});
    higlightedSquares.push({ ...moves[i].to, color: moves[i].piece.color});
  }
  if (selectedPiece) {
    higlightedSquares.push({ ...selectedPiece, color: activePlayer });
  }

  return (
    <div className="board-container" style={{
      boxSizing: 'border-box',
      padding: '10px',
      position: 'relative',
      height: '100%',
    }}>
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
            position: 'fixed',
            top: '50%',
            transform: 'translate(-50%, -50%)',
          }} />
        </div>
        <PlayerIndicator color={activePlayer} />
        </div>
    </div>
  );
} 