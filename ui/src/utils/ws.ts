import { Color, PGNMove } from '../common';

export enum MessageType {
	SetSettings = 'setSettings',
	GetAvailableMoves = 'getAvailableMoves',
	AvailableMoves = 'availableMoves',
	PlayerMove = 'playerMove',
	EngineMove = 'engineMove',
	InvalidMove = 'invalidMove',
	NewGame = 'newGame',
	SaveGame = 'saveGame',
	SaveGameResponse = 'saveGameResponse',
	LoadGame = 'loadGame',
	LoadGameResponse = 'loadGameResponse',
	SetCurrentMove = 'setCurrentMove',
	GameEnded = 'gameEnded',
	Processing = 'processing',
	StoppedProcessing = 'stoppedProcessing',
}

export interface BestMoveResponse {
	continuation: PGNMove[];
	score: number;
	time: number;
	evaluations: number;
	moveNumber: number;
}

export interface SaveGameResponse {
	pgn: string;
}

export interface LoadGameResponse {
	pastMoves: PGNMove[];
	currentMove: number;
}

export interface GameSettings {
	humanPlayers: number[];
	depth: number;
	spread: number;
	spreadDrop: number;
	evalLimit: number;
}

export interface GameEndedResponse {
	king: string;
	winner: string;
}

type MessageData =
	| PGNMove
	| PGNMove[]
	| BestMoveResponse
	| SaveGameResponse
	| LoadGameResponse
	| GameSettings
	| GameEndedResponse
	| string
	| number
	| null;

export class Message {
	constructor(public type: MessageType, public data: MessageData) {}

	json() {
		return JSON.stringify({
			type: this.type,
			data: this.data,
		});
	}
}
