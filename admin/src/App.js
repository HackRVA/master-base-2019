import React, { Component } from 'react';
import { withStyles } from '@material-ui/core/styles';
import outrun from './outrun-background.jpg';
import GameForm from './components/GameForm';
import GameList from './components/GameList';
import Button from '@material-ui/core/Button';
import './App.css';

const styles = theme => ({
  btn: {
    background: theme.palette.primary.main,
  }
});

class App extends Component {
  state = { showGameForm: false }
  handleSchedule = () => {
    this.setState({
      showGameForm: !this.state.showGameForm
    })
  }
  render() {
    console.log('app props: ', this.props)
    return (
      <div style={{ backgroundImage: `url(${outrun})`, height: '100vh' }} className="App">
        <div className="transparent-background">
          <div className="header">
            <h1>Master Base Station</h1>
            {this.state.showGameForm && <GameForm handleSchedule={this.handleSchedule} />}
            {!this.state.showGameForm && <Button className={this.props.classes.btn} variant="contained" onClick={this.handleSchedule}>Schedule Game</Button>}
            <GameList />
          </div>
        </div>
      </div>
    );
  }
}

export default withStyles(styles)(App);
