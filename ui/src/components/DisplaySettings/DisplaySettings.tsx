import React from 'react';
import { Color, colorCode } from '../../common';
import { useBoardStateContext } from '../../context/BoardStateContext';
import { ShowLabels } from '../../utils';
import { Checkbox } from '../Checkbox';
import styles from './DisplaySettings.module.css';

const showLabelsOptions: ShowLabels[] = ['all', 'border', 'pieces', 'moves', 'moves+', 'none'];

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
									onChange={e => setDisplaySettings({ ...displaySettings, showLabels: e.target.value as ShowLabels })}
								>
									{showLabelsOptions.map(option => <option key={option} value={option}>{capitalize(option)}</option>)}
								</select>
							</td>
						</tr>
						<tr>
							<td>Show Continuation:</td>
							<td>
								<Checkbox
									checked={displaySettings.showContinuation}
									onChange={checked => setDisplaySettings({ ...displaySettings, showContinuation: checked })}
									background={colorCode(Color.DarkGray)}
									borderColor={colorCode(Color.DarkGray)}
								/>
							</td>
						</tr>
					</tbody>
				</table>
			</div>
		</div>
	);
}
