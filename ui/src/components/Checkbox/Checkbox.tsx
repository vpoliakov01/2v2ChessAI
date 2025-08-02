import React from 'react';
import { Color, colorCode } from '../../common';

interface CheckboxProps {
  background: string;
  borderColor: string;
  checked: boolean;
  onChange: (checked: boolean) => void;
}

export function Checkbox({ background, borderColor, checked, onChange }: CheckboxProps) {
  return (
    <div
      className={`custom-checkbox ${checked ? 'checked' : ''}`}
      onClick={() => onChange(!checked)}
      style={{
        background: checked ? background : 'transparent',
        borderWidth: 2,
        borderStyle: 'solid',
        borderColor,
        alignItems: 'center',
        borderRadius: 5,
        cursor: 'pointer',
        display: 'flex',
        height: 20,
        justifyContent: 'center',
        width: 20,
      }}
    >
      {checked && (
        <span style={{
          userSelect: 'none',
          color: colorCode(Color.LightGray), 
          fontSize: 18, 
          fontWeight: 'bold',
          fontFamily: 'Roboto',
          marginRight: 1,
        }}>✓</span>
      )}
    </div>
  );
}
