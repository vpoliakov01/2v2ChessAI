import React, { createContext, useContext, ReactNode, useState } from 'react';
import { useBoardState } from '../hooks/useBoardState';
import { DisplaySettingsState, loadDisplaySettingsFromStorage } from '../components/DisplaySettings';
import { Color, Move } from '../common';

type HoveredMove = {
  move: Move;
  color: Color;
};

type BoardStateContextType = ReturnType<typeof useBoardState> & {
  displaySettings: DisplaySettingsState;
  setDisplaySettings: React.Dispatch<React.SetStateAction<DisplaySettingsState>>;
  hoveredMove: HoveredMove | null;
  setHoveredMove: React.Dispatch<React.SetStateAction<HoveredMove | null>>;
};

const BoardStateContext = createContext<BoardStateContextType | null>(null);

export const useBoardStateContext = () => {
  const context = useContext(BoardStateContext);
  if (!context) {
    throw new Error('useBoardStateContext must be used within a BoardStateProvider');
  }
  return context;
};

export const BoardStateProvider = ({ children }: { children: ReactNode }) => {
  const boardState = useBoardState();
  const [displaySettings, setDisplaySettings] = useState<DisplaySettingsState>(loadDisplaySettingsFromStorage);
  const [hoveredMove, setHoveredMove] = useState<HoveredMove | null>(null);
  return (
    <BoardStateContext.Provider value={{ ...boardState, displaySettings, setDisplaySettings, hoveredMove, setHoveredMove }}>
      {children}
    </BoardStateContext.Provider>
  );
};
