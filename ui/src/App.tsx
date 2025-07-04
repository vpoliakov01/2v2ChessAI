import React from 'react';
import { ChessBoard } from './components/ChessBoard';
import { Menu } from './components/Menu';
import { Color, colorCode } from './common';
import { BoardStateProvider } from './context/BoardStateContext';

function App() {
  return (
    <div className="App" style={{
      backgroundColor: colorCode(Color.Black),
      color: 'white',
      height: '100vh',
      textAlign: 'center',
      width: '100vw',
      display: 'flex',
      justifyContent: 'center',
      alignItems: 'center',
    }}>
      <BoardStateProvider>
        <ChessBoard />
        <Menu />
      </BoardStateProvider>
    </div>
  );
}

export default App;
