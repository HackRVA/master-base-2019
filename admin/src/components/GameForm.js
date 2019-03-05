import React from 'react';
import TimePicker from './TimePicker';
import NativeSelect from '@material-ui/core/NativeSelect';
import Input from '@material-ui/core/Input';
import Button from '@material-ui/core/Button';
import { withStyles } from '@material-ui/core/styles';

const styles = theme => ({
  btn: {
    background: theme.palette.primary.main,
  }
})

const GameForm = props => {
  const { classes } = props;
  return <div style={{ justifyContent: 'left' }}>
    <NativeSelect
      // value={this.state.age}
      // onChange={this.handleChange('gameType')}
      input={<Input name="gameType" id="gameType" />}
    >
      <option value={"DEATMATCH"}>DeathMatch</option>
      <option value={"TEAMDEATHMATCH"}>Team DeathMatch</option>
    </NativeSelect>
    <TimePicker label={"Game Start Time"} />
    <TimePicker label={"Game Duration"} />

    <Button className={classes.btn} variant="outlined" onClick={props.handleSchedule}>Submit</Button>
  </div>
}

export default withStyles(styles)(GameForm);