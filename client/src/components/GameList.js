import React from 'react';
import Games from '../games.json';

const GameList = () => {
  return <div>
    <h3> Game List...</h3>
    {Games.map((x, i) => {
      return <div key={i}>
        <ul style={{ listStyle: "none" }}>
          <li> GameID: {x.GameID} </li>
          <li> StartTime: {x.StartTime} </li>
          <li> Duration: {x.Duration} </li>
          <li> Variant: {x.Variant} </li>
          <li> Team: {x.Team} </li>
        </ul>
      </div>
    })}
  </div>
}

export default GameList;