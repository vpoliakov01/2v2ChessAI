import React, { useState } from 'react';
import { Color, colorCode, movesToPGN } from '../common';
import { useBoardStateContext } from '../context/BoardStateContext';
import { Message, MessageType } from '../ws';
import { CollapsibleBlock, BORDER_RADIUS } from './CollapsibleBlock';
import { MoveTable } from './MoveTable';
import { Settings } from './Settings';
import { DisplaySettings } from './DisplaySettings';

export function Menu() {
  const { allMoves, currentMove, setCurrentMove, sendMessage } = useBoardStateContext();
  const [pgnBlockCollapsed, setPGNBlockCollapsed] = useState(false);
  const [userPGN, setUserPGN] = useState<string | null>(null);

  const pgn = userPGN != null ? userPGN : movesToPGN(allMoves);

  function handleNewGame(event: React.MouseEvent<HTMLButtonElement>) {
    sendMessage(new Message(MessageType.NewGame, null));
    event.stopPropagation();
  }

  function handleCopy(event: React.MouseEvent<HTMLButtonElement>) {
    navigator.clipboard.writeText(movesToPGN(allMoves));
    setPGNBlockCollapsed(false);
    event.stopPropagation();
  }

  function handleLoad(event: React.MouseEvent<HTMLButtonElement>) {
    sendMessage(new Message(MessageType.LoadGame, pgn));
    setUserPGN(null);
    setPGNBlockCollapsed(true);
    event.stopPropagation();
  }

  function handleSetCurrentMove(moveIndex: number) {
    setCurrentMove(moveIndex);
    sendMessage(new Message(MessageType.SetCurrentMove, moveIndex));
  }

  return (
    <div className="menu-container" style={{
      boxSizing: 'border-box',
      height: '100%',
      padding: 10,
      position: 'relative',
      width: 420,
    }}>
      <div className="menu" style={{
        backgroundColor: colorCode(Color.DarkGray),
        boxSizing: 'border-box',
        width: '100%',
        padding: 5,
        borderRadius: BORDER_RADIUS + 5,
      }}>
        <CollapsibleBlock 
          collapsed={pgnBlockCollapsed}
          header={
            <div style={{
              display: 'flex',
              gap: 10,
            }}>
              <button id="button-new-game" onClick={handleNewGame}>New Game</button>
              <button id="button-copy" onClick={handleCopy}>Copy</button>
              <button id="button-load" onClick={handleLoad}>Load</button>
            </div>
          }>
          <textarea id="game-save-text"
            value={pgn}
            onChange={(e) => setUserPGN(e.target.value)}
            onBlur={() => setUserPGN(userPGN || movesToPGN(allMoves))} // Reset on empty userPGN.
            style={{
              backgroundColor: colorCode(Color.Black),
              borderRadius: BORDER_RADIUS / 2,
              boxSizing: 'border-box',
              color: colorCode(Color.Gray),
              fontSize: 15,
              minHeight: 100,
              padding: 10,
              resize: 'vertical',
              width: '100%',
          }} />
        </CollapsibleBlock>
        <CollapsibleBlock header="Settings" collapsed={false}>
          <Settings />
        </CollapsibleBlock>
        <CollapsibleBlock header="Display Settings" collapsed={false}>
          <DisplaySettings />
        </CollapsibleBlock>
        <CollapsibleBlock header="Moves" collapsed={false}>
          {allMoves.length > 0 && <MoveTable moves={allMoves} currentMove={currentMove} handleSetCurrentMove={handleSetCurrentMove} />}
        </CollapsibleBlock>
      </div>
    </div>
  );
} 