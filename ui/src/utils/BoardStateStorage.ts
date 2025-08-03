import { Piece, Color, MoveInfo } from '../common';
import { ArrowProps } from '../components/Arrow';

export type BoardPosition = (Piece | null | undefined)[][];

export interface SavedBoardState {
    board: BoardPosition;
    activePlayer: Color;
    allMoves: MoveInfo[];
    currentMove: number;
    score: number;
    pgn: string;
    drawnArrows: ArrowProps[];
}

/**
 * Handles complete storage and loading of board state to/from localStorage.
 * 
 * Features:
 * - Preserves undefined vs null distinction in board positions
 * - Properly reconstructs Move and MoveInfo class instances
 * - Handles serialization of complex nested objects
 */
export class BoardStateStorage {
    private static readonly STORAGE_KEY = 'chess-board-state';
    private static readonly UNDEFINED_MARKER = '__UNDEFINED__';

    // Stores the complete board state to localStorage.
    static store(state: SavedBoardState): void {
        try {
            const serializedState = BoardStateStorage.serializeState(state);
            localStorage.setItem(BoardStateStorage.STORAGE_KEY, JSON.stringify(serializedState));
        } catch (error) {
            console.warn('Failed to save board state to localStorage:', error);
        }
    }

    // Loads the complete board state from localStorage.
    // Returns null if no state exists or if loading fails.
    static load(): SavedBoardState | null {
        try {
            const stored = localStorage.getItem(BoardStateStorage.STORAGE_KEY);
            if (!stored) {
                return null;
            }

            const parsedState = JSON.parse(stored);
            return BoardStateStorage.deserializeState(parsedState);
        } catch (error) {
            console.warn('Failed to load board state from localStorage:', error);
            return null;
        }
    }

    // Clears the stored board state from localStorage.
    static clear(): void {
        try {
            localStorage.removeItem(BoardStateStorage.STORAGE_KEY);
        } catch (error) {
            console.warn('Failed to clear board state from localStorage:', error);
        }
    }

    private static serializeState(state: SavedBoardState): any {
        return {
            ...state,
            board: BoardStateStorage.serializeBoardPosition(state.board),
            allMoves: state.allMoves.map(move => BoardStateStorage.serializeMoveInfo(move)),
        };
    }

    private static deserializeState(serializedState: any): SavedBoardState {
        return {
            ...serializedState,
            board: BoardStateStorage.deserializeBoardPosition(serializedState.board || []),
            allMoves: (serializedState.allMoves || []).map((move: any) =>
                BoardStateStorage.deserializeMoveInfo(move)
            ),
        };
    }

    private static serializeBoardPosition(board: BoardPosition): any[][] {
        return board.map(row =>
            row.map(cell => {
                if (cell === undefined) return BoardStateStorage.UNDEFINED_MARKER;
                if (cell === null) return null;
                return cell;
            })
        );
    }

    private static deserializeBoardPosition(serializedBoard: any[][]): BoardPosition {
        if (!Array.isArray(serializedBoard)) {
            return [];
        }

        return serializedBoard.map(row =>
            row.map(cell => {
                if (cell === BoardStateStorage.UNDEFINED_MARKER) return undefined;
                if (cell === null) return null;
                return cell as Piece;
            })
        );
    }

    private static serializeMoveInfo(moveInfo: MoveInfo): any {
        return {
            from: moveInfo.from,
            to: moveInfo.to,
            piece: moveInfo.piece,
            capturedPiece: moveInfo.capturedPiece,
        };
    }

    private static deserializeMoveInfo(serializedMove: any): MoveInfo {
        return new MoveInfo(
            serializedMove.from,
            serializedMove.to,
            serializedMove.piece,
            serializedMove.capturedPiece
        );
    }
}