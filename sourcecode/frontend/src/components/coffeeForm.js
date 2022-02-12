import AddCircleOutlineIcon from '@mui/icons-material/AddCircleOutline';
import IconButton from '@mui/material/IconButton';
import Box from '@mui/material/Box';
import Slider from '@mui/material/Slider';
import TextField from '@mui/material/TextField';
import InputAdornment from '@mui/material/InputAdornment';
import Button from '@mui/material/Button';
import Stack from '@mui/material/Stack';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogTitle from '@mui/material/DialogTitle';
import Typography from '@mui/material/Typography';

import React, { useState } from 'react';

export default function CoffeeForm(props) {
  const [open, setOpen] = React.useState(false);

  const handleClickOpen = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
  };

  const [newCoffee, setNewCoffee] = useState({
    Name: '',
    Weight: '',
    Price: '',
    RoastLevel: 3
  })

  const handleChange = (prop) => (event) => {
    setNewCoffee({ ...newCoffee, [prop]: event.target.value });
  };

  const handleAdd = () => {
    setOpen(false);
    console.log('newCoffeee', newCoffee)
    props.addNewCoffee(newCoffee)
    setNewCoffee(
      {
        Name: '',
        Weight: '',
        Price: '',
        RoastLevel: 3
      }
    )
  }

  return (
    <div>
      <IconButton color="inherit" onClick={() => handleClickOpen()}><AddCircleOutlineIcon sx={{ fontSize: 50 }} /></IconButton>

      <Dialog open={open} onClose={handleClose}>
        <DialogTitle>New Coffee</DialogTitle>
        <DialogContent>
          <Box>
            <Box sx={{ width: '30ch' }}>
              <TextField
                required
                fullWidth
                margin='normal'
                value={newCoffee.Name}
                onChange={handleChange('Name')}
                id="nimi"
                label="Name"
                variant="outlined" />
              <TextField
                required
                fullWidth
                type='number'
                margin='normal'
                value={newCoffee.Weight}
                onChange={handleChange('Weight')}
                InputProps={{
                  endAdornment: <InputAdornment position="end">grams</InputAdornment>,
                }}
                id="paino"
                label="Weight"
                variant="outlined" />
              <TextField
                required
                fullWidth
                type='number'
                margin='normal'
                InputProps={{
                  endAdornment: <InputAdornment position="end">€</InputAdornment>,
                }}
                value={newCoffee.Price}
                onChange={handleChange('Price')}
                id="hinta"
                label="Price"
                variant="outlined"
              />

              <Slider
                sx={{ marginTop: '3ch' }}
                color='primary'
                aria-label="RoastLevel"
                // defaultValue={3}
                value={newCoffee.RoastLevel}
                onChange={handleChange('RoastLevel')}
                valueLabelDisplay="on"
                step={1}
                //marks
                min={1}
                max={5}
              />
              <Typography id="RoastLevel" gutterBottom>
                Roast Level
              </Typography>
            </Box>
            <Stack spacing={2} direction="row">

            </Stack>

          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => handleAdd()} variant="contained">Add Coffee</Button>
          <Button onClick={handleClose} variant="text" color="error">Cancel</Button>
        </DialogActions>
      </Dialog>
    </div>
  )
}