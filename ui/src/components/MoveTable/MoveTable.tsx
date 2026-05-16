import React, { useState } from 'react';
import { Move, PlayerColors } from '../../common';
import { useBoardStateContext } from '../../context/BoardStateContext';
import { OnMoveHover } from '../../utils';
import styles from './MoveTable.module.css';

interface MoveTableProps {
	moves: Move[];
	currentMove: number;
	handleSetCurrentMove: (moveIndex: number) => void;
	startOffset?: number;
	overrideHoverMode?: OnMoveHover;
}

export function MoveTable(
	{ moves, currentMove, handleSetCurrentMove, startOffset = 0, overrideHoverMode }: MoveTableProps,
) {
	const { setHoveredMove, hoveredMove } = useBoardStateContext();
	const [selectedMove, setSelectedMove] = useState<number | null>(null);

	const hoverMode = overrideHoverMode ?? 'set board';

	const handleMouseEnter = (moveIndex: number) => {
		if (selectedMove !== null) {
			return;
		}

		switch (hoverMode) {
			case 'arrow':
			case 'highlight':
			case 'highlight+':
				setHoveredMove({ move: moves[moveIndex], color: PlayerColors[(moveIndex + startOffset) % 4] });
				break;
			case 'set board':
				handleSetCurrentMove(moveIndex);
				break;
			case 'none':
				setHoveredMove(null);
				break;
		}
	};

	const handleMouseLeave = (moveIndex: number) => {
		switch (hoverMode) {
			case 'arrow':
			case 'highlight':
			case 'highlight+':
				setHoveredMove(null);
				break;
			case 'set board':
				// Handled by handleTableMouseLeave.
				break;
			case 'none':
				break;
		}
	};

	const handleTableMouseLeave = () => {
		if (hoverMode !== 'set board') {
			return;
		}
		const target = selectedMove ?? moves.length - 1;
		if (target !== currentMove) {
			handleSetCurrentMove(target);
		}
	};

	const handleClick = (moveIndex: number) => {
		setSelectedMove(moveIndex);
		handleSetCurrentMove(moveIndex);
	};

	// Total cells (including the leading inactive padding cells).
	const totalCells = startOffset + moves.length;
	// Round down the first row to a multiple of 4 so column colors align to players.
	const firstAbsoluteIndex = Math.floor(startOffset / 4) * 4;

	const rows = [];
	for (let i = firstAbsoluteIndex; i < totalCells; i += 4) {
		const cells = Array.from({ length: 4 }).map((_, j) => {
			const cellIndex = i + j;
			const moveIndex = cellIndex - startOffset;
			if (moveIndex < 0 || moveIndex >= moves.length) {
				return <td key={`${i}-${j}-empty`} className={styles.inactiveCell}></td>;
			}

			const isHovered = hoveredMove?.move === moves[moveIndex];
			return (
				<td
					className={[styles.moveCell, moveIndex === currentMove || isHovered ? styles.currentMove : ''].filter(Boolean)
						.join(' ')}
					key={`${i}-${moves[moveIndex].toPGN()}`}
					onClick={() => handleClick(moveIndex)}
					onMouseEnter={() => handleMouseEnter(moveIndex)}
					onMouseLeave={() => handleMouseLeave(moveIndex)}
				>
					{moves[moveIndex].toPGN()}
				</td>
			);
		});
		rows.push(
			<tr key={`${i / 4 + 1}-row`}>
				<td className={styles.moveNumber} key={`${i}-number`}>{i / 4 + 1}.</td>
				{cells}
			</tr>,
		);
	}

	return (
		<table className={styles.moveTable} onMouseLeave={handleTableMouseLeave}>
			<thead>{rows}</thead>
		</table>
	);
}
