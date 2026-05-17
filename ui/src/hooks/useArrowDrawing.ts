import { useCallback, useState } from 'react';
import { Color, Move, movesEqual, Position } from '../common';
import { ArrowProps } from '../components/Arrow';
import { GameStateManager } from '../utils';

export function useArrowDrawing(activePlayer: Color) {
	const [drawnArrows, setDrawnArrows] = useState<ArrowProps[]>(() => GameStateManager.load().drawnArrows);
	const [isDrawingArrow, setIsDrawingArrow] = useState<boolean>(false);
	const [arrowStart, setArrowStart] = useState<Position | null>(null);
	const [arrowEnd, setArrowEnd] = useState<Position | null>(null);

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
			color: activePlayer,
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

	return {
		drawnArrows,
		isDrawingArrow,
		arrowStart,
		arrowEnd,
		handleSquareRightMouseDown,
		handleSquareMouseEnter,
		handleSquareRightMouseUp,
		handleSquareLeftClick,
	};
}
