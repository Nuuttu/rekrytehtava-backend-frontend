import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Paper from '@mui/material/Paper';
import DeleteIcon from '@mui/icons-material/Delete';
import IconButton from '@mui/material/IconButton';



export default function CoffeeList({ coffeeList, deleteCoffee }) {

  const handleDeleteClick = (coffee) => {
    if (window.confirm('Do you want to delete this coffee: ' + coffee.Name + '?')) {
      deleteCoffee(coffee)
    }
  }

  return (
    <TableContainer component={Paper}>
      <Table>
        <TableHead>
          <TableRow>
            <TableCell align="right">Name</TableCell>
            <TableCell align="right">Weight&nbsp;(g)</TableCell>
            <TableCell align="right">Price&nbsp;(€)</TableCell>
            <TableCell align="right">RoastLevel&nbsp;(1-5)</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {coffeeList.map((c, i) =>
            <TableRow key={c.Id}>
              <TableCell component="th" scope="row">
                {c.Name}
              </TableCell>
              <TableCell align="right">
                {c.Weight}
              </TableCell>
              <TableCell align="right">
                {c.Price.toFixed(2)}
              </TableCell>
              <TableCell align="right">
                {c.RoastLevel}
              </TableCell>
              <TableCell align="right">
                <IconButton color="inherit" onClick={() => handleDeleteClick(c)}><DeleteIcon /></IconButton>
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>
    </TableContainer>
  )
}