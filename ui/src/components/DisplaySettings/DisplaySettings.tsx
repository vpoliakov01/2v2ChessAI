import React from 'react';
import { useBoardStateContext } from '../../context/BoardStateContext';
import { GameStateManager, ShowLabels, OnMoveHover } from '../../utils';
import styles from './DisplaySettings.module.css';

const showLabelsOptions: ShowLabels[] = ['all', 'border', 'pieces', 'moves', 'moves+', 'none'];
const onMoveHoverOptions: OnMoveHover[] = ['set board', 'arrow', 'highlight', 'highlight+', 'none'];



function capitalize(s: string): string {
  return s.charAt(0).toUpperCase() + s.slice(1);
}

export function DisplaySettings() {
  const { displaySettings, setDisplaySettings } = useBoardStateContext();

  // Update localStorage whenever settings change
  React.useEffect(() => {
    GameStateManager.saveDisplaySettings(displaySettings);
  }, [displaySettings]);

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
