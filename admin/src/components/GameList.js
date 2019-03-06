import React from 'react';
import Games from '../games.json';

const GameList = () => {
  return <div>
    <h3> Game List...</h3>
    {Games.map((x, i) => {
      const { GameID, StartTime, Duration, Variant, Team } = x;
      return <div key={i}>
        <ul style={{ listStyle: "none" }}>
          <li> GameID: {GameID} </li>
          <li> StartTime: {StartTime} </li>
          <li> Duration: {Duration} </li>
          <li> Variant: {Variant} </li>
          <li> Team: {Team} </li>
        </ul>
      </div>
    })}
  </div>
}

export default GameList;