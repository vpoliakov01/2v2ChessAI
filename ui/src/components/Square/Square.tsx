import React from 'react';
import { Color, colorCode, Piece, pieceName } from '../../common';
import styles from './Square.module.css';

interface SquareProps {
  isPlayable: boolean;
  isLight: boolean;
  piece: Piece | null | undefined;
  higlighted: Color | null;
  possibleMove: boolean;
  label: string;
  onClick: () => void;
}

export function Square({ isPlayable, isLight, piece, higlighted, possibleMove, label, onClick }: SquareProps) {
  if (!isPlayable) {
    return <div
      className={styles.squareNonPlayable}
      style={{
        backgroundColor: colorCode(Color.Black),
        border: `1px solid ${colorCode(Color.Black)}`,
      }}
    />;
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
      className={styles.square}
      style={{
        backgroundColor: higlighted && !possibleMove ? higlightedBackgroundColor : backgroundColor,
        cursor: piece ? 'pointer' : 'default',
      }}
      onMouseEnter={onMouseEnter}
      onMouseLeave={onMouseLeave}
      onClick={onClick}
    >
      {possibleMove && (
        <div className={styles.possibleMoveIndicator} />
      )}
      {piece && (
        <img
          alt={`${pieceName[piece.type]}`}
          src={getPieceImage(piece)}
          className={styles.pieceImage}
        />
      )}
      {label && (
        <span className={styles.squareLabel}>{label}</span>
      )}
    </div>
  );
}