import { useState } from 'react'
import './App.css'
import Game from './Game';

function App() {
  const [title, setTitle] = useState("")

  const fetchTitle = async () => {
    try {
      const response = await fetch('http://localhost:8080/');
      if (!response.ok) {
        throw new Error('Failed to fetch title');
      }
      const data = await response.text();
      setTitle(data);
    } catch (error) {
      console.error('Error fetching title', error);
      throw error;
    }
  };
  fetchTitle()

  return (
    <div>
      <h1> {title} </h1>
      <Game />
    </div>

  )
}

export default App
