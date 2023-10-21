import {useState, useEffect} from 'react';
import logo from './assets/images/logo-universal.png';
import './App.css';
import {StartBalancer, StopBalancer} from "../wailsjs/go/main/App";

import React from 'react';

const MyComponent = () => {
    const [count, setCount] = useState(0);
  
    useEffect(() => {
      document.title = `Count: ${count}`;
    }, [count]);
  
    return (
      <div>
        <h1>Balanced: {count}</h1>
        <button onClick={() => setCount(count + 1)}>Increment</button>
      </div>
    );
  };

function App() {
    const [resultText, setResultText] = useState("Please enter your name below 👇");
    const [name, setName] = useState('');
    const updateName = (e) => setName(e.target.value);
    const updateResultText = (result) => setResultText(result);

    function startBalancer() {
        StartBalancer(name).then(updateResultText);
    }

    function stopBalancer() {
        StopBalancer(name).then(updateResultText);
    }

    return (
        <div id="App">
            <MyComponent />
            <h1>CounterBlow Load Balancer</h1>
            <p></p>
            <div id="result" className="result">{resultText}</div>
            <div id="input" className="input-box">
                <span className="hitCounter"> Listen to port:</span>
                <input id="name" className="input" onChange={updateName} autoComplete="off" name="input" type="text" value="2345"/><br /><br />
                <input id="name" className="input" onChange={updateName} autoComplete="off" name="input" type="text"/>
                <br /><br />
                <button className="btn" onClick={startBalancer}>Start</button>

                <button className="btn" disabled onClick={stopBalancer}>Stop</button>

            </div>
            <div>
            <h2>Log</h2>
              <textarea className="textAreaField"></textarea>

            </div>
        </div>
    )
}

export default App
