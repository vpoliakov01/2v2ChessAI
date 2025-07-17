import React, { createContext, useContext, ReactNode, useState } from 'react';
import { useBoardState } from '../hooks/useBoardState';
import { DisplaySettingsState, defaultDisplaySettings } from '../components/DisplaySettings';

type BoardStateContextType = ReturnType<typeof useBoardState> & {
  displaySettings: DisplaySettingsState;
  setDisplaySettings: React.Dispatch<React.SetStateAction<DisplaySettingsState>>;
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
  const [displaySettings, setDisplaySettings] = useState<DisplaySettingsState>(defaultDisplaySettings);
  return (
    <BoardStateContext.Provider value={{ ...boardState, displaySettings, setDisplaySettings }}>
      {children}
    </BoardStateContext.Provider>
  );
};
