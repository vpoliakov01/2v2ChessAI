import React from 'react';
import { BOARD_SIZE, Color, CORNER_SIZE, colorCode } from '../../common';
import styles from './PlayerIndicator.module.css';

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
          right: offsetLength,
          top: color === Color.Yellow ? -offset : `calc(100% - ${lineWidth - offset}px)`,
        };
      case Color.Blue:
      case Color.Green:
        return {
          backgroundColor: colorCode(color),
          bottom: offsetLength,
          left: color === Color.Blue ? -offset : `calc(100% - ${lineWidth - offset}px)`,
          top: offsetLength,
          width: lineWidth,
        };
      default:
        return {};
    }
  }
  return <div className={styles.playerIndicator} style={getPlayerIndicatorStyle(color)} />;
}