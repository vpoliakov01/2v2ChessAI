import React, { useState } from 'react';
import { Color, colorCode, movesToPGN } from '../common';
import { useBoardStateContext } from '../context/BoardStateContext';
import { Message, MessageType } from '../ws';
import { MoveTable } from './MoveTable';

interface CollapsibleBlockProps {
  header: React.ReactNode;
  children: React.ReactNode;
  collapsed?: boolean;
}

const BORDER_RADIUS = 10;

function CollapsibleBlock({ header, children, collapsed=true }: CollapsibleBlockProps) {
  const [isCollapsed, setIsCollapsed] = useState(collapsed);

  return (
    <div className="collapsible-block" style={{
      backgroundColor: colorCode(Color.Black),
      borderRadius: BORDER_RADIUS,
      marginBottom: 5,
      width: '100%',
    }}>
      <div 
        className="collapsible-header"
        onClick={() => setIsCollapsed(!isCollapsed)}
        style={{
          alignItems: 'center',
          cursor: 'pointer',
          display: 'flex',
          justifyContent: 'space-between',
          padding: '10px 12px 10px 10px',
        }}
      >
        {header}
        <span style={{
          fontSize: 20,
          position: 'relative',
          bottom: 2,
        }}>{isCollapsed ? '+' : '-'}</span>
      </div>
      <div 
        className="collapsible-content"
        style={{
          maxHeight: isCollapsed ? 0 : 'fit-content',
          overflow: 'hidden',
          padding: isCollapsed ? 0 : '0 10px 10px 10px',
          transition: 'all 0.02s linear',
        }}
      >
        {children}
      </div>
    </div>
  );
}

export function Menu() {
  const { wsRef, allMoves, currentMove, setCurrentMove } = useBoardStateContext();
  const [pgnBlockCollapsed, setPGNBlockCollapsed] = useState(false);
  const [userPGN, setUserPGN] = useState<string | null>(null);

  function handleNewGame(event: React.MouseEvent<HTMLButtonElement>) {
    if (wsRef.current) {
      wsRef.current.send(new Message(
        MessageType.NewGame,
        null,
      ).json());
    }
    event.stopPropagation();
  }

  function handleCopy(event: React.MouseEvent<HTMLButtonElement>) {
    navigator.clipboard.writeText(movesToPGN(allMoves));
    setPGNBlockCollapsed(false);
    event.stopPropagation();
  }

  function handleLoad(event: React.MouseEvent<HTMLButtonElement>) {
    if (wsRef.current) {
      wsRef.current.send(new Message(
        MessageType.LoadGame,
        userPGN || movesToPGN(allMoves),
      ).json());
    }
    setUserPGN(null);
    setPGNBlockCollapsed(true);
    event.stopPropagation();
  }

  function handleSetCurrentMove(moveIndex: number) {
    setCurrentMove(moveIndex);
    if (wsRef.current) {
      wsRef.current.send(new Message(
        MessageType.SetCurrentMove,
        moveIndex,
      ).json());
    }
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
        height: '100%',
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
            value={userPGN || movesToPGN(allMoves)}
            onChange={(e) => setUserPGN(e.target.value)}
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
        <CollapsibleBlock header="Moves" collapsed={false}>
          {allMoves.length > 0 && <MoveTable moves={allMoves} currentMove={currentMove} handleSetCurrentMove={handleSetCurrentMove} />}
        </CollapsibleBlock>
      </div>
    </div>
  );
} 