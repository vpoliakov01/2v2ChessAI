import React from 'react';
import { ChessBoard } from './components/ChessBoard';
import { Color, colorCode } from './common';

function App() {
  return (
    <div className="App" style={{
      backgroundColor: colorCode[Color.Black],
      color: 'white',
      height: '100vh',
      textAlign: 'center',
      width: '100vw',
      display: 'flex',
      justifyContent: 'center',
      alignItems: 'center',
    }}>
      <ChessBoard />
    </div>
  );
}

export default App;
