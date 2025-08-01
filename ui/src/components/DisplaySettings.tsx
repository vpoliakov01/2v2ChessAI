import React from 'react';
import { useBoardStateContext } from '../context/BoardStateContext';

const showLabelsOptions = ['all', 'border', 'pieces', 'moves', 'moves+', 'none'] as const;
const onMoveHoverOptions = ['set board', 'arrow', 'highlight', 'none'] as const;

type ShowLabels = typeof showLabelsOptions[number];
type OnMoveHover = typeof onMoveHoverOptions[number];

export interface DisplaySettingsState {
  showLabels: ShowLabels;
  onMoveHover: OnMoveHover;
}

export const defaultDisplaySettings: DisplaySettingsState = {
  showLabels: 'moves',
  onMoveHover: 'set board',
};

function capitalize(s: string): string {
  return s.charAt(0).toUpperCase() + s.slice(1);
}

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
                <select
                  value={displaySettings.showLabels}
                  onChange={(e) => setDisplaySettings({ ...displaySettings, showLabels: e.target.value as ShowLabels })}
                >
                  {showLabelsOptions.map(option => (
                    <option key={option} value={option}>{capitalize(option)}</option>
                  ))}
                </select>
              </td>
            </tr>
            <tr>
              <td>On move hover:</td>
              <td>
                <select
                  value={displaySettings.onMoveHover}
                  onChange={(e) => setDisplaySettings({ ...displaySettings, onMoveHover: e.target.value as OnMoveHover })}
                >
                  {onMoveHoverOptions.map(option => (
                    <option key={option} value={option}>{capitalize(option)}</option>
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
