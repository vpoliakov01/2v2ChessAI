import React from 'react';
import { useBoardStateContext } from '../context/BoardStateContext';

type ShowLabels = 'all' | 'border' | 'pieces' | 'moves' | 'moves+' | 'none';
export interface DisplaySettingsState {
  showLabels: ShowLabels;
}

const showLabelsOptions: ShowLabels[] = ['all', 'border', 'pieces', 'moves', 'moves+', 'none'];

export const defaultDisplaySettings: DisplaySettingsState = {
  showLabels: 'border',
};

export function DisplaySettings() {
  const { displaySettings, setDisplaySettings } = useBoardStateContext();

  return (
    <div id="display-settings" style={{
      width: '100%',
      marginTop: 10,
    }}>
      <div id="settings-table">
        <table>
          <tbody>
            <tr>
              <td>Show labels:</td>
              <td>
                <select value={displaySettings.showLabels} onChange={(e) => setDisplaySettings({ ...displaySettings, showLabels: e.target.value as ShowLabels })}>
                  {showLabelsOptions.map(option => (
                    <option value={option}>{option.charAt(0).toUpperCase() + option.slice(1)}</option>
                  ))}
                </select>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  );
}
