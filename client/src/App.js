import React, { Component } from 'react';
import outrun from './outrun-background.jpg';
import './App.css';

class App extends Component {
  render() {
    return (
      <div style={{ backgroundImage: `url(${outrun})`, height: '100vh' }} className="App">
        <div className="transparent-background">
          <div className="header">
            <h1>Master Base Station</h1>
          </div>
        </div>
      </div>
    );
  }
}

export default App;
