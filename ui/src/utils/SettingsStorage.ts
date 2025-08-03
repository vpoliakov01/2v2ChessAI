import { GameSettings } from './ws';

const SETTINGS_STORAGE_KEY = 'chess-game-settings';

export const defaultSettings: GameSettings = {
    humanPlayers: [0, 2],
    depth: 6,
    captureDepth: 8,
    evalLimit: 0,
};

export function loadSettingsFromStorage(): GameSettings {
    const stored = localStorage.getItem(SETTINGS_STORAGE_KEY);
    let storedSettings: Partial<GameSettings> = {};

    if (stored) {
        storedSettings = JSON.parse(stored);
    }

    return {
        ...defaultSettings,
        ...storedSettings,
    };
}

export function saveSettingsToStorage(settings: GameSettings): void {
    localStorage.setItem(SETTINGS_STORAGE_KEY, JSON.stringify(settings));
}
