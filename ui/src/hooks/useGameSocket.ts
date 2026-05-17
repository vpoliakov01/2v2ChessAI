import { useCallback, useEffect, useRef } from 'react';
import { GameSyncService, Message } from '../utils';

export function useGameSocket(onMessage: (message: Message) => void) {
	const wsRef = useRef<WebSocket | null>(null);
	const handlerRef = useRef(onMessage);
	handlerRef.current = onMessage;

	const sendMessage = useCallback((message: Message) => {
		if (wsRef.current?.readyState === WebSocket.OPEN) {
			console.debug(`Sending    ${message.type}`.padEnd(30), message);
			wsRef.current.send(message.json());
			return;
		}

		setTimeout(() => {
			if (wsRef.current?.readyState === WebSocket.OPEN) {
				console.debug(`Sending    ${message.type}`.padEnd(30), message);
				wsRef.current.send(message.json());
			} else {
				alert('No WebSocket connection');
			}
		}, 10);
	}, []);

	useEffect(() => {
		const ws = new WebSocket('ws://localhost:8080/ws');

		ws.onopen = () => {
			wsRef.current = ws;
			console.log('WS connected');
			GameSyncService.syncWithEngine(ws);
		};

		ws.onmessage = event => {
			const message = JSON.parse(event.data) as Message;
			console.debug(`Received   ${message.type}`.padEnd(30), message);
			handlerRef.current(message);
		};

		ws.onclose = () => {
			if (wsRef.current === ws) {
				wsRef.current = null;
			}
			console.log('WS disconnected');
		};

		return () => ws.close();
	}, []);

	return sendMessage;
}
