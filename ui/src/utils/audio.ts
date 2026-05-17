export enum Sound {
	Move = 'move',
	Capture = 'capture',
	Promote = 'promote',
	GameStart = 'game-start',
	GameEnd = 'game-end',
}

const preloaded: Record<Sound, HTMLAudioElement> = (() => {
	const map = {} as Record<Sound, HTMLAudioElement>;
	for (const sound of Object.values(Sound)) {
		const audio = new Audio(`/${sound}.mp3`);
		audio.preload = 'auto';
		audio.load();
		map[sound] = audio;
	}
	return map;
})();

export function playSound(sound: Sound): void {
	const audio = preloaded[sound];
	audio.currentTime = 0;
	audio.play().catch(() => {});
}
