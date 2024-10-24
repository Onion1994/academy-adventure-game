import { useState } from 'react'

function Game() {
    const [input, setInput] = useState("")
    const [output, setOutput] = useState("")

    const fetchResponse = async (e) => {
        e.preventDefault()
        try {
            const commandArgsToSend = parseInput();
            console.log(commandArgsToSend);
            const requestData = {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify(commandArgsToSend),
            };
            const response = await fetch('http://localhost:8080/GameResponse', requestData);
            if (!response.ok) {
                throw new Error('Failed to fetch response');
            }
            const data = await response.json();
            setOutput(data.message);
        } catch (error) {
            console.error('Error fetching response', error);
            throw error;
        }
    };

    const parseInput = () => {
        const splitCommands = input.split(" ");
        if (splitCommands.length > 1) {
            return { command: splitCommands[0], args: splitCommands.slice(1) };
        } else {
            return { command: splitCommands[0], args: [""] };
        }
    }
    const handleInput = (event) => {
        const value = event.target.value;
        setInput(value);

    }

    return (
        <>
            <input name="gameInput" type='text' value={input} onChange={handleInput} />
            <button type='submit' onClick={fetchResponse} > Submit </button>
            <p> {output} </p>
        </>
    )
}

export default Game
