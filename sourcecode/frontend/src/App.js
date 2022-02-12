import './App.css';
import React, { useState, useEffect } from 'react';
import CoffeeForm from './components/coffeeForm';
import CoffeeList from './components/coffeeList';

const Notify = ({ errorMessage }) => {
  if (!errorMessage) {
    return null
  }
  return (
    <div style={{ color: 'red' }}>
      {errorMessage}
    </div>
  )
}

function App() {
  const [errorMessage, setErrorMessage] = useState(null)
  const [coffeeList, setCoffeeList] = useState([])

  const getCoffeeListData = () => {
    try {
      fetch('http://localhost:10000/coffeelist')
        .then(response => {
          if (response.ok) {
            return response.json()
          }
          throw response;
        })
        .then(data => setCoffeeList(data))
        .catch(err => {
          notify(err)
          console.log('error happened', err)
        })
    } catch (e) {
      console.log('error', e)
    }
  }

  useEffect(() => {
    getCoffeeListData()
  }, [])

  const addNewCoffee = (coffee) => {
    const newCoffee = {
      Name: coffee.Name,
      Weight: parseInt(coffee.Weight),
      Price: parseFloat(coffee.Price),
      RoastLevel: parseInt(coffee.RoastLevel)
    }
    console.log('Adding new coffee: ', newCoffee)
    let request = {
      method: 'POST',
      mode: 'no-cors',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(newCoffee)
    };
    try {
      fetch(`http://localhost:10000/coffee/add`, request)
        .then(response => response.json())
    } catch (e) {
      console.log('error', e)

    }
    setTimeout(() => {
      getCoffeeListData()
      notify("Coffee Added")
    }, 1000)

  }

  const deleteCoffee = (coffee) => {
    let request = { method: 'DELETE' };
    fetch(`http://localhost:10000/coffee/delete/${coffee.Id}`, request)
      .then(response => response.json())
      .then(data => {
        console.log('response data, deleted: ', data)

        setTimeout(() => {
          notify("deleted")
          getCoffeeListData()
        }, 500)
      })
      .catch(err => {
        notify(err)
        console.log('error happened', err)
      })

  }

  const notify = (message) => {
    setErrorMessage(message)
    setTimeout(() => {
      setErrorMessage(null)
    }, 5000)
  }

  return (
    <div className="App">
      <header className="App-header">
        <div>
          <Notify errorMessage={errorMessage} />
          <h2>Coffees</h2>
          <CoffeeForm addNewCoffee={addNewCoffee} />
          {coffeeList != null ?
            <CoffeeList coffeeList={coffeeList} deleteCoffee={deleteCoffee} />
            :
            <h3>Coffee List is empty!</h3>
          }
        </div>
      </header>
    </div>
  );
}

export default App;
