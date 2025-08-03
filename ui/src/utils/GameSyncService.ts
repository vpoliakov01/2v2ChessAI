import { MessageType } from './ws';
import { GameStateManager } from './GameStateManager';
import { movesToPGN } from '../common';

export class GameSyncService {
    static async syncWithEngine(ws: WebSocket): Promise<void> {
        if (ws.readyState !== WebSocket.OPEN) return;

        const state = GameStateManager.load();
        const settings = GameStateManager.loadSettings();

        const engineSettings = {
            depth: settings.depth,
            captureDepth: settings.captureDepth,
            humanPlayers: settings.humanPlayers,
            evalLimit: settings.evalLimit,
        };

        ws.send(JSON.stringify({ type: MessageType.SetSettings, data: engineSettings }));

        if (state.allMoves.length > 0) {
            const pgn = state.pgn || movesToPGN(state.allMoves);
            if (pgn.trim()) {
                ws.send(JSON.stringify({ type: MessageType.LoadGame, data: pgn }));
                if (state.currentMove !== state.allMoves.length) {
                    ws.send(JSON.stringify({ type: MessageType.SetCurrentMove, data: state.currentMove }));
                }
            } else {
                ws.send(JSON.stringify({ type: MessageType.NewGame, data: {} }));
            }
        } else {
            ws.send(JSON.stringify({ type: MessageType.NewGame, data: {} }));
        }

        ws.send(JSON.stringify({ type: MessageType.GetAvailableMoves, data: {} }));
    }
}