import React, { useState } from 'react';
import { Color, colorCode } from '../../common';
import styles from './CollapsibleBlock.module.css';

interface CollapsibleBlockProps {
  header: React.ReactNode;
  children: React.ReactNode;
  collapsed?: boolean;
}

export function CollapsibleBlock({ header, children, collapsed = true }: CollapsibleBlockProps) {
  const [isCollapsed, setIsCollapsed] = useState(collapsed);

  return (
    <div
      className={styles.collapsibleBlock}
      style={{
        backgroundColor: colorCode(Color.Black),
      }}
    >
      <div
        className={styles.collapsibleHeader}
        onClick={() => setIsCollapsed(!isCollapsed)}
      >
        {header}
        <span className={styles.toggleIcon}>
          {isCollapsed ? '+' : '-'}
        </span>
      </div>
      <div
        className={styles.collapsibleContent}
        style={{
          maxHeight: isCollapsed ? 0 : 'fit-content',
          padding: isCollapsed ? 0 : '0 10px 10px 10px',
        }}
      >
        {children}
      </div>
    </div>
  );
}