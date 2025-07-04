import { createContext, useContext, ReactNode } from 'react';
import { useBoardState } from '../hooks/useBoardState';

type BoardStateContextType = ReturnType<typeof useBoardState> | null;

const BoardStateContext = createContext<BoardStateContextType>(null);

export const useBoardStateContext = () => {
  const context = useContext(BoardStateContext);
  if (!context) {
    throw new Error('useBoardStateContext must be used within a BoardStateProvider');
  }
  return context;
};

export const BoardStateProvider = ({ children }: { children: ReactNode }) => {
  const boardState = useBoardState();
  return (
    <BoardStateContext.Provider value={boardState}>
      {children}
    </BoardStateContext.Provider>
  );
};
