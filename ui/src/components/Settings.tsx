import React, { useState, useEffect } from 'react';
import { Message, MessageType, GameSettings } from '../ws';
import { useBoardStateContext } from '../context/BoardStateContext';
import { Color, PlayerColors, colorCode } from '../common';
import { Checkbox } from './Checkbox';
import { NumberInput } from './NumberInput';

export function Settings() {
  const { sendMessage } = useBoardStateContext();
  const [settings, setSettings] = useState<GameSettings>({
    humanPlayers: [0, 2],
    depth: 6,
    captureDepth: 8,
    evalLimit: 0,
  });

  useEffect(() => {
    sendMessage(new Message(MessageType.SetSettings, settings));
  }, [settings, sendMessage]);
  
  return (
    <div id="settings" style={{
      width: '100%',
      marginTop: 10,
    }}>
      <div className="settings-row" style={{
        display: 'flex',
        alignItems: 'center',
        gap: '10px',
      }}>
        <label>Human Players:</label>
        <div style={{ display: 'flex', gap: '5px' }}>
          {PlayerColors.map(color => (
            <Checkbox
              key={`human-player-checkbox-${color}`}
              checked={settings.humanPlayers.includes(PlayerColors.indexOf(color))}
              onChange={checked => setSettings({
                ...settings,
                humanPlayers: checked ?
                  [...settings.humanPlayers, PlayerColors.indexOf(color)].sort()
                  : settings.humanPlayers.filter(c => c !== PlayerColors.indexOf(color)),
              })}
              background={colorCode(color)}
              borderColor={`color-mix(in srgb, ${colorCode(color)} 80%, ${colorCode(Color.DarkGray)})`}
            />
          ))}
          <Checkbox
            key="human-player-checkbox-all"
            checked={settings.humanPlayers.length === PlayerColors.length}
            background={`conic-gradient(
              from 135deg,
              ${colorCode(Color.Red)} 0deg 90deg,
              ${colorCode(Color.Blue)} 90deg 180deg,
              ${colorCode(Color.Yellow)} 180deg 270deg,
              ${colorCode(Color.Green)} 270deg 360deg
            )`}
            borderColor={`${colorCode(Color.Yellow)} ${colorCode(Color.Green)} ${colorCode(Color.Red)} ${colorCode(Color.Blue)}`}
            onChange={checked => setSettings({
              ...settings,
              humanPlayers: checked ?
                PlayerColors.map((_, i) => i).sort()
                : [],
            })}
          />
        </div>
      </div>
      <div id="settings-table">
        <table>
          <tbody>
            <tr>
              <td>Depth:</td>
              <td>
                <NumberInput
                  value={settings.depth}
                  onChange={(value) => setSettings({ ...settings, depth: value })}
                  min={1}
                />
              </td>
              <td>Forcing:</td>
              <td>
                <NumberInput
                  value={settings.captureDepth}
                  onChange={(value) => setSettings({ ...settings, captureDepth: value })}
                  min={settings.depth}
                />
              </td>
            </tr>
            <tr>
              <td>Eval Limit, k:</td>
              <td>
                <NumberInput
                  value={settings.evalLimit}
                  onChange={(value) => setSettings({ ...settings, evalLimit: value * 1000 })}
                  min={0}
                  editable={true}
                  disableButtons={true}
                  width={84}
                />
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  );
}
