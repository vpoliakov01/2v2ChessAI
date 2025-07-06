import React, { useState } from 'react';

interface NumberInputProps {
  value: number;
  onChange: (value: number) => void;
  min?: number;
  max?: number;
  step?: number;
  editable?: boolean;
  disableButtons?: boolean;
  width?: number;
}

export function NumberInput({ value, onChange, min = 0, max = Infinity, step = 1, editable = false, disableButtons = false, width = 28 }: NumberInputProps) {
  const [editValue, setEditValue] = useState(value.toString());

  const increment = () => {
    if (value + step <= max) {
      setEditValue((value + step).toString());
      onChange(value + step);
    }
  };

  const decrement = () => {
    if (value - step >= min) {
      setEditValue((value - step).toString());
      onChange(value - step);
    }
  };

  const validateValue = (value: string) => {
    const newValue = parseInt(editValue);
    if (!isNaN(newValue) && newValue.toString().length === value.length && newValue >= min && newValue <= max) {
      setEditValue(value.toString());
      onChange(newValue);
    } else {
      setEditValue(min.toString());
      onChange(min);
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      validateValue(editValue);
    } else if (e.key === 'Escape') {
      e.preventDefault();
      setEditValue(value.toString());
    }
  };

  const handleBlur = () => {
    validateValue(editValue);
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setEditValue(e.target.value);
  };

  const handleFocus = () => {
    if (value === 0) {
      setEditValue('');
    }
  };

  return (
    <div className="number-input-container">
      <div className="number-input-controls">
        {!disableButtons && <button 
          className="number-input-button" 
          onClick={decrement}
          disabled={value <= min}
        >
          -
        </button>}
        <input
          type="text"
          value={editValue}
          onChange={handleChange}
          onFocus={handleFocus}
          onKeyDown={handleKeyDown}
          onBlur={handleBlur}
          className="number-input-display"
          style={{
              width,
              borderRadius: disableButtons ? '5px' : 0,
          }}
          disabled={!editable}
          onWheel={(e) => e.currentTarget.blur()}
        />
        {!disableButtons && <button 
          className="number-input-button"
          onClick={increment}
          disabled={value >= max}
        >
          +
        </button>}
      </div>
    </div>
  );
} 