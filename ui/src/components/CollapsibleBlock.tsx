import React, { useState } from 'react';
import { Color, colorCode } from '../common';

export const BORDER_RADIUS = 10;

interface CollapsibleBlockProps {
  header: React.ReactNode;
  children: React.ReactNode;
  collapsed?: boolean;
}

export function CollapsibleBlock({ header, children, collapsed=true }: CollapsibleBlockProps) {
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
          fontWeight: 'bold',
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
          textAlign: 'left',
        }}
      >
        {children}
      </div>
    </div>
  );
}