import React from 'react';
import { BOARD_SIZE, Color, colorCode, Move } from '../../common';
import styles from './Arrow.module.css';

export interface ArrowProps {
    move: Move;
    color: Color;
    short?: boolean;
}

export function Arrow({ move, color, short = false }: ArrowProps) {
    // Calculate the position and dimensions
    const squareSize = 100 / BOARD_SIZE; // Percentage

    // Center positions of squares
    const fromX = (move.from.col + 0.5) * squareSize;
    const fromY = (move.from.row + 0.5) * squareSize;
    const toX = (move.to.col + 0.5) * squareSize;
    const toY = (move.to.row + 0.5) * squareSize;

    // Calculate distance for zero-distance check
    const deltaX = toX - fromX;
    const deltaY = toY - fromY;
    const distance = Math.sqrt(deltaX * deltaX + deltaY * deltaY);

    // Arrow dimensions
    const arrowWidth = 0.18 * squareSize; // Line thickness
    const arrowHeadWidth = 0.25 * squareSize; // Width of arrowhead
    const arrowHeadLength = arrowHeadWidth * 1.4; // Base length of arrowhead

    // Normalize direction
    const dirX = deltaX / distance;
    const dirY = deltaY / distance;

    // Arrow base is the point where the arrowhead meets the stem
    const arrowBaseX = toX - dirX * arrowHeadLength;
    const arrowBaseY = toY - dirY * arrowHeadLength;

    const strokeColor = colorCode(color);
    const opacity = 0.8;

    if (distance === 0) { // Draw a circle if it's a single square highlight
        return (
            <div className={styles.arrowContainer}>
                <svg
                    className={styles.arrowSvg}
                    viewBox="0 0 100 100"
                    preserveAspectRatio="none"
                >
                    <circle
                        cx={fromX}
                        cy={fromY}
                        r={100 / BOARD_SIZE / 2 * 0.95}
                        fill="none"
                        stroke={strokeColor}
                        strokeWidth={0.5}
                        opacity={opacity}
                    />
                </svg>
            </div >
        );
    }

    return (
        <div className={styles.arrowContainer}>
            <svg
                className={styles.arrowSvg}
                viewBox="0 0 100 100"
                preserveAspectRatio="none"
            >
                <defs>
                    <marker
                        id={`arrowhead-${move.toPGN()}`}
                        orient="auto"
                        overflow="visible"
                        markerWidth={arrowHeadWidth}
                        markerHeight={arrowHeadLength}
                        refX={short ? squareSize * 0.16 : squareSize * 0.08}
                        refY={arrowHeadWidth}
                    >
                        <path
                            d={`M0,0 V${arrowHeadWidth * 2} L${arrowHeadLength},${arrowHeadWidth} Z`}
                            fill={strokeColor}
                        />
                    </marker>
                </defs>

                {/* Arrow line with marker */}
                <line
                    x1={fromX}
                    y1={fromY}
                    x2={arrowBaseX}
                    y2={arrowBaseY}
                    stroke={strokeColor}
                    strokeWidth={arrowWidth}
                    strokeLinecap="round"
                    opacity={opacity}
                    markerEnd={`url(#arrowhead-${move.toPGN()})`}
                />
            </svg>
        </div>
    );
}