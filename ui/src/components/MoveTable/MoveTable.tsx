import React, { useEffect } from 'react';
import { Move, PlayerColors } from '../../common';
import { BOARD_SIZE } from '../../common';
import { useBoardStateContext } from '../../context/BoardStateContext';
import styles from './MoveTable.module.css';

interface MoveTableProps {
  moves: Move[];
  currentMove: number;
  handleSetCurrentMove: (moveIndex: number) => void;
}

export function MoveTable({ moves, currentMove, handleSetCurrentMove }: MoveTableProps) {
  const { displaySettings, setHoveredMove } = useBoardStateContext();

  useEffect(() => {
    const handleArrowPress = (e: KeyboardEvent) => {
      if ((e.key === 'ArrowLeft' || e.key === 'ArrowUp') && currentMove > 0) {
        handleSetCurrentMove(currentMove - 1);
      } else if ((e.key === 'ArrowRight' || e.key === 'ArrowDown') && currentMove < moves.length - 1) {
        handleSetCurrentMove(currentMove + 1);
      }
    };
    window.addEventListener('keydown', handleArrowPress);
    return () => {
      window.removeEventListener('keydown', handleArrowPress);
    };
  }, [moves, currentMove, handleSetCurrentMove]);

  const handleMouseEnter = (moveIndex: number) => {
    switch (displaySettings.onMoveHover) {
      case 'arrow':
      case 'highlight':
        setHoveredMove({ move: moves[moveIndex], color: PlayerColors[moveIndex % 4] });
        break;
      case 'set board':
        handleSetCurrentMove(moveIndex);
        break;
      case 'none':
        setHoveredMove(null);
        break;
    }
  }

  const handleMouseLeave = (moveIndex: number) => {
    switch (displaySettings.onMoveHover) {
      case 'arrow':
      case 'highlight':
        setHoveredMove(null);
        break;
      case 'set board':
        // Handled onMouseLeave the table.
        break;
      case 'none':
        break;
    }
  }

  const rows = [];
  for (let i = 0; i < moves.length; i += 4) {
    const cells = Array.from({ length: 4 }).map((_, j) => (
      i + j < moves.length ? (
        <td
          className={`${styles.moveCell} ${i + j === currentMove ? styles.currentMove : ''}`}
          key={`${i}-${moves[i + j].toPGN()}`}
          onClick={() => handleSetCurrentMove(i + j)}
          onMouseEnter={() => handleMouseEnter(i + j)}
          onMouseLeave={() => handleMouseLeave(i + j)}
        >
          {moves[i + j].toPGN()}
        </td>
      ) : (
        <td key={`${i}-${j}`}></td>
      )
    ));
    rows.push(
      <tr key={`${i / 4 + 1}-row`}>
        <td className={styles.moveNumber} key={`${i}-number`}>{i / 4 + 1}.</td>
        {cells}
      </tr>
    );
  }

  return (
    <table className={styles.moveTable}
      onMouseLeave={() => handleSetCurrentMove(moves.length - 1)}
    >
      <thead>
        {rows}
      </thead>
    </table>
  );
}