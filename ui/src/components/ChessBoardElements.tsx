import React from 'react';

import { BOARD_SIZE, Color, CORNER_SIZE, colorCode, Piece, pieceName } from '../common';

interface SquareProps {
  isPlayable: boolean;
  isLight: boolean;
  piece: Piece | null | undefined;
  higlighted: Color | null;
  possibleMove: boolean;
  onClick: () => void;
}

export function Square({ isPlayable, isLight, piece, higlighted, possibleMove, onClick }: SquareProps) {
  if (!isPlayable) {
    return <div style={{
      aspectRatio: '1',
      backgroundColor: colorCode(Color.Black),
      border: `1px solid ${colorCode(Color.Black)}`,
      width: '100%',
    }} />;
  }

  let backgroundColor = isLight ? colorCode(Color.LightGray) : colorCode(Color.Gray);
  const originalBackgroundColor = backgroundColor;
  let higlightedBackgroundColor = originalBackgroundColor;
  if (higlighted) {
    higlightedBackgroundColor = `color-mix(in srgb, ${colorCode(higlighted)} 45%, ${backgroundColor})`;
  }

  const getPieceImage = (piece: Piece) => {
    const color = piece.color;
    const type = piece.type;
    return `/${color}_${pieceName[type]}.svg`;
  };

  function onMouseEnter(e: React.MouseEvent<HTMLDivElement>) {
    if (possibleMove) {
      e.currentTarget.style.backgroundColor = higlightedBackgroundColor;
    }
  }
  function onMouseLeave(e: React.MouseEvent<HTMLDivElement>) {
    if (possibleMove) {
      e.currentTarget.style.backgroundColor = originalBackgroundColor;
    }
  }

  return (
    <div
      style={{
        alignItems: 'center',
        aspectRatio: '1',
        backgroundColor: higlighted && !possibleMove ? higlightedBackgroundColor : backgroundColor,
        cursor: piece ? 'pointer' : 'default',
        display: 'flex',
        justifyContent: 'center',
        position: 'relative',
        userSelect: 'none',
        width: '100%',
      }}
      onMouseEnter={onMouseEnter}
      onMouseLeave={onMouseLeave}
      onClick={onClick}
    >
      {possibleMove && (
        <div style={{
          backgroundColor: 'rgba(0, 0, 0, 0.2)',
          borderRadius: '50%',
          height: '30%',
          left: '35%',
          top: '35%',
          width: '30%',
        }} />
      )}
      {piece && (
        <img 
          src={getPieceImage(piece)} 
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
    const lineWidth = 3;
    const offset = 6;

    switch (color) {
      case Color.Red:
      case Color.Yellow:
        return {
          backgroundColor: colorCode(color),
          height: lineWidth,
          left: offsetLength,
          position: 'absolute',
          right: offsetLength,
          top: color === Color.Yellow ? -offset : `calc(100% - ${lineWidth-offset}px)`,
        };
      case Color.Blue:
      case Color.Green:
        return {
          backgroundColor: colorCode(color),
          bottom: offsetLength,
          left: color === Color.Blue ? -offset : `calc(100% - ${lineWidth-offset}px)`,
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

export function ScoreDisplay({ score }: { score: number }) {
  const maxScore = 10;
  const offsetLength = `calc(${CORNER_SIZE / BOARD_SIZE} * 100%)`;
  const height = Math.max(Math.min(50 + score / maxScore / 2 * 100, 100), 0);

  return <div className="score-display" style={{
    width: 20,
    marginRight: 10,
    position: 'relative',
    top: offsetLength,
    height: `${(BOARD_SIZE - 2 * CORNER_SIZE) / BOARD_SIZE * 100}%`,
  }}>
    <div style={{
      backgroundColor: colorCode(Color.Blue),
      height: `${100 -height}%`,
    }} />
    <div style={{
      backgroundColor: colorCode(Color.Red),
      height: `${height}%`,
    }} />
  </div>;
}
