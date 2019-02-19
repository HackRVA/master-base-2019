import React from 'react';

const GameForm = props => {
  return <div style={{ justifyContent: 'left' }}>
    <div>
      Game Type:
    <select>
        <option value="grapefruit">DeathMatch</option>
        <option value="lime">Team DeathMatch</option>
      </select>
    </div>
    <div>
      Game Start Time:
      <input />
    </div>
    <div>
      Game Duration:
      <input />
    </div>
    <button onClick={props.handleSchedule}>Submit</button>

  </div>
}

export default GameForm;