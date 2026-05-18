import { useEffect, useState } from 'react';
import { GameSettings, GameStateManager, Message, MessageType } from '../utils';

export function useGameSettings(sendMessage: (message: Message) => void) {
	const [settings, setSettings] = useState<GameSettings>(GameStateManager.loadSettings);

	useEffect(() => {
		const engineSettings = {
			...settings,
			evalLimit: settings.evalLimit * 1000,
		};
		sendMessage(new Message(MessageType.SetSettings, engineSettings));
		GameStateManager.saveSettings(settings);
	}, [settings, sendMessage]);

	return { settings, setSettings };
}
