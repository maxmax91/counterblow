import {useState, useEffect} from 'react';
import logo from './assets/images/logo-universal.png';
import './App.css';
import {StartBalancer, StopBalancer} from "../wailsjs/go/main/App";

import React from 'react';
import ReactDOM from 'react-dom';
import IPut from 'iput';

runtime.EventsOn("rcv:update_served_pages", (msg) => document.getElementById("served_pages").innerText = msg)
runtime.EventsOn("rcv:add_log_string", (msg) => document.getElementById("textAreaLog").value += msg)

const MyComponent = () => {
 
    return (
      <div>
        
        <button onClick={() => setCount(count + 1)}>Increment</button>
      </div>
    );
  };

  const GridComponent = () => {
    return (
      <div className="square">left top </div>
    );
  };

function App() {
    const [resultText, setResultText] = useState("Please enter your name below 👇");
    const [name, setName] = useState('');
    const updateName = (e) => setName(e.target.value);
    const updateResultText = (result) => setResultText(result);

    function startBalancer(el) {
        document.getElementById('buttonStart').disabled = true;
        document.getElementById('buttonStop').disabled = false;
        StartBalancer(name).then(updateResultText);
    }

    function stopBalancer() {
      document.getElementById('buttonStart').disabled = false;
      document.getElementById('buttonStop').disabled = true;
        StopBalancer(name).then(updateResultText);
    }

    runtime.EventsOn("rcv:greet", (msg) => document.getElementById("result").innerText = msg)
    window.StartUpdateServedPagesEvent = () => {
      // noinspection JSIgnoredPromiseFromCall
      window.go.main.App.StartUpdateServedPagesEvent();
    }


    return (
        <div id="App">
           <h1>CounterBlow Load Balancer</h1>
            <div className="square bordered">
              <h2>Rules config</h2>
              
              <select multiple>
                <option value="1">IP: 1.1.1.1/21 Rule: RoundRobin Servers: boh1.com boh2.com</option>
                <option value="2">Rule 2</option>
                <option value="3">Rule 3</option>
              </select>
              
              <div className="formRow">
              <span>IP Addr:</span>
              <IPut className="IPut" defaultValue="0.0.0.0"/> / <input className="portInput" type="text" value="0" />
              </div>

              <div className="formRow">
              <span>Backend servers (comma-separated)</span>
              </div>
              <div className="formRow">
              <input className="servers" type="text" placeholder="google.it:80,microsoft.it:80" />
              </div>
              <span>Algorythm</span>
              <select>
                <option value="1">1: Round robin</option>
                <option value="2">2: IP hash</option>
              </select>
              <button>&#x2795;</button>
              <button>&#x2796;</button>
            </div>
            
            <h2>Served connections: <div id="served_pages">0</div></h2>
            <p></p>
            <div id="result" className="result">{resultText}</div>
            <div id="input" className="input-box">
                <span className="hitCounter"> Listen to IP/port:</span>
                <IPut className="IPut" id="bindIp" defaultValue="0.0.0.0" />
                <input id="name" className="input" onChange={updateName} autoComplete="off" name="input" type="text" value="2345"/>
                <br /><br />
                <button className="btn" id="buttonStart" onClick={startBalancer}>Start</button>

                <button className="btn" id="buttonStop" disabled onClick={stopBalancer}>Stop</button>

            </div>
            <div>
              <div className="square">
              <h2>Log</h2>
              <textarea className="textAreaField" id="textAreaLog"></textarea>
              </div>
            </div>
        </div>
    )
}

export default App
