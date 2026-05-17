import { useEffect, useState } from 'react';
import { DisplaySettingsState, GameStateManager } from '../utils';

export function useDisplaySettings() {
	const [displaySettings, setDisplaySettings] = useState<DisplaySettingsState>(GameStateManager.loadDisplaySettings);

	useEffect(() => {
		GameStateManager.saveDisplaySettings(displaySettings);
	}, [displaySettings]);

	return { displaySettings, setDisplaySettings };
}
