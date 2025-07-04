import React from 'react';
import { Move } from '../common';
import { BOARD_SIZE } from '../common';

interface MoveTableProps {
  moves: Move[];
  currentMove: number;
  handleSetCurrentMove: (moveIndex: number) => void;
}

export function MoveTable({moves, currentMove, handleSetCurrentMove}: MoveTableProps) {
  const toPGN = (move: Move) => {
    const toFile = (col: number) => String.fromCharCode(col + 'a'.charCodeAt(0));
    const toRank = (row: number) => BOARD_SIZE - row;
    return `${toFile(move.from.col)}${toRank(move.from.row)}-${toFile(move.to.col)}${toRank(move.to.row)}`;
  };

  const rows = [];
  for (let i = 0; i < moves.length; i += 4) {
    const cells = Array.from({length: 4}).map((_, j) => (
      i + j < moves.length ? (
        <td
          className={`move-cell ${i + j === currentMove ? 'current-move' : ''}`}
          key={`${i}-${toPGN(moves[i+j])}`}
          onClick={() => handleSetCurrentMove(i + j)}
        >
          {toPGN(moves[i+j])}
        </td>
      ) : (
        <td key={`${i}-${j}`}></td>
      )
    ));
    rows.push(
      <tr key={`${i / 4 + 1}-row`}>
        <td className="move-number" key={`${i}-number`}>{i / 4 + 1}.</td>
        {cells}
      </tr>
    );
  }

  return (
    <table id="move-table" style={{
      width: '100%',
      marginTop: 10,
    }}>
      <thead>
        {rows}
      </thead>
    </table>
  );
}