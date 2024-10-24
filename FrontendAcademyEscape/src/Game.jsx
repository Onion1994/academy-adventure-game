import { useState } from 'react'

function Game() {
  // const [input, setInput] = useState("")
  const [output,setOutput] = useState("")

  const fetchResponse = async () => {
    try {
        const requestData = {
            method :  "POST",
            headers : {"Content-Type": "application/json"},
            body : JSON.stringify({command : "start" , args: [""]})
        };
      const response = await fetch('http://localhost:8080/GameResponse');
      if (!response.ok) {
        throw new Error('Failed to fetch response');
      }
      const data = await response.json();
      console.log(data);
      setOutput(data);
    } catch (error) {
      console.error('Error fetching response', error);
      throw error;
    }
  };
  fetchResponse()

  return (
    <p> {output} </p>
  )
}

export default Game
