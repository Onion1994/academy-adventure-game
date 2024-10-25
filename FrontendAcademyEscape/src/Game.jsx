import { useState, useEffect } from 'react';
import './Game.css';

function Game() {
    const [commands, setCommands] = useState([
        "look",
        "exit",
        "commands",
        "take",
        "drop",
        "inventory",
        "approach",
        "use",
        "leave",
        "move",
        "map"
    ]);

    //Dummy actions for development must replace 
    
    const [availableActions, setAvailableActions] = useState({
        look: [],
        exit: [],
        commands: [],
        take: ["tea"],
        drop: [],
        inventory: [],
        approach: ["rosie", "kettle", "sofa", "cat"],
        use: ["tea"],
        leave: [],
        move: [],
        map: []
    });

    const [selectedCommand, setSelectedCommand] = useState('');
    const [selectedArgument, setSelectedArgument] = useState('');
    const [output, setOutput] = useState('');
    const [gameStarted, setGameStarted] = useState(false);

    const startGame = async () => {
        const commandArgsToSend = { command: "start", args: [] };
        try {
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
            setGameStarted(true); 
        } catch (error) {
            console.error('Error fetching response', error);
            setOutput('Error starting game');
        }
    };

    const fetchResponse = async (e) => {
        e.preventDefault();
        try {
            const commandArgsToSend = {
                command: selectedCommand,
                args: selectedArgument ? [selectedArgument] : []
            };
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
            setOutput('Error fetching response');
        }
    };

    return (
        <>
        {!gameStarted ? (
            <button onClick={startGame}>Start Game</button>
        ) : (
            <form onSubmit={fetchResponse}>
                <select onChange={(e) => setSelectedCommand(e.target.value)} value={selectedCommand}>
                    <option value="">Select a command</option>
                    {commands.map((command) => (
                        <option key={command} value={command}>{command}</option>
                    ))}
                </select>

                {selectedCommand && availableActions[selectedCommand].length > 0 && (
                    <select onChange={(e) => setSelectedArgument(e.target.value)} value={selectedArgument}>
                        <option value="">Select an argument</option>
                        {availableActions[selectedCommand].map((action) => (
                            <option key={action} value={action}>{action}</option>
                        ))}
                    </select>
                )}

                <button type='submit'>Submit</button>
            </form>
        )}
            <p>{output}</p>
        </>
    );
}

export default Game;
