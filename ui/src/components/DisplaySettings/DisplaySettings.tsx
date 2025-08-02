import React from 'react';
import { useBoardStateContext } from '../../context/BoardStateContext';
import styles from './DisplaySettings.module.css';

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
  onMoveHover: 'arrow',
};

function capitalize(s: string): string {
  return s.charAt(0).toUpperCase() + s.slice(1);
}

export function DisplaySettings() {
  const { displaySettings, setDisplaySettings } = useBoardStateContext();

  return (
    <div className={styles.displaySettings}>
      <div className={styles.settingsTable}>
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
