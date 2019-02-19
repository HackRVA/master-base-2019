import React, { Component } from 'react';
import outrun from './outrun-background.jpg';
import GameForm from './components/GameForm';
import GameList from './components/GameList';
import './App.css';

class App extends Component {
  state = { showGameForm: false }
  handleSchedule = () => {
    this.setState({
      showGameForm: !this.state.showGameForm
    })
  }
  render() {
    return (
      <div style={{ backgroundImage: `url(${outrun})`, height: '100vh' }} className="App">
        <div className="transparent-background">
          <div className="header">
            <h1>Master Base Station</h1>
            {this.state.showGameForm && <GameForm handleSchedule={this.handleSchedule} />}
            {!this.state.showGameForm && <button onClick={this.handleSchedule}>Schedule Game</button>}
            <GameList />
          </div>
        </div>
      </div>
    );
  }
}

export default App;
