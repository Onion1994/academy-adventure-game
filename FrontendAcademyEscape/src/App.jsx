import { useState } from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import './App.css'

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
    <h1> {title} </h1>
  )
}

export default App
