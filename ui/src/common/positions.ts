import { BOARD_SIZE } from "./constants";

export interface Position {
    row: number;
    col: number;
}

export function positionToPGN(position: Position): string {
    return `${String.fromCharCode(position.col + 'a'.charCodeAt(0))}${BOARD_SIZE - position.row}`;
}

export function positionsEqual(a: Position, b: Position): boolean {
    return a.row === b.row && a.col === b.col;
}
