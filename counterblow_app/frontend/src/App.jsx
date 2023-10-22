import {useState, useEffect} from 'react';
import logo from './assets/images/logo-universal.png';
import './App.css';
import {StartBalancer, StopBalancer, OnDOMContentLoaded, AddRule, RemoveRule, RefreshRules} from "../wailsjs/go/main/App";

import React from 'react';
import ReactDOM from 'react-dom';

runtime.EventsOn("rcv:update_served_pages", (msg) => document.getElementById("served_pages").innerText = msg)
runtime.EventsOn("rcv:clear_rules_listbox", () => document.getElementById("rules_list_box").innerHTML = "")


runtime.EventsOn("rcv:add_served_rule", function(rule_id, rule_type, rule_ipaddr, rule_subnetmask, rule_servers, rule_source, rule_dest) {
  var listbox = document.getElementById("rules_list_box");
  const opt1 = document.createElement("option");
  opt1.value = rule_id;
  opt1.text = "Type: " + rule_type + " From: " + rule_ipaddr + "/" + rule_subnetmask + " - Servers: " + rule_servers + " [rewrite from '" + rule_source + "' to '" + rule_dest + "']";
  listbox.add(opt1);
})
runtime.EventsOn("rcv:add_log_string", function(msg) { 
  var el = document.getElementById("textAreaLog");
  if (el != null) document.getElementById("textAreaLog").value += msg;
  else console.log("Cannot print msg" + msg);
 }
)

/*
document.addEventListener('DOMContentLoaded', (event) => {
  //alert(document.getElementById("rules_list_box"));
  alert("Welcome to CounterBlow load balancer"); // there are some problems with this function!
                // bug: OnDOMContentLoaded called but DOM is not actually ready
  var listbox = document.getElementById("rules_list_box");
  if (listbox == null) {
    console.log("DOM not ready.");
  }
  OnDOMContentLoaded('Loading...').then((result) => {  if (result != null) alert(result) ; } ); 
});
*/


function App() {
  const updateResultStartBalancer = (result) => {
    alert(result);
  } ;


  function startBalancer(el) {
      var port = parseInt (document.getElementById('bindPort').value);
      // todo: add bing ip
      StartBalancer('0.0.0.0', port).then(updateResultStartBalancer); // will not return
      document.getElementById('buttonStart').disabled = true;
      document.getElementById('buttonStop').disabled = false;
      document.getElementById('loaderSpinner').style.visibility = 'visible';
  }

  function addRule() {
    // aggiungi una regola
    var rule_type = document.getElementById('rule_type').value;
    var rule_ipaddr = '0.0.0.0';
    var rule_subnetmask = 0;
    var rule_servers =  document.getElementById('rule_servers').value;
    if (rule_servers == 0) {
      alert("Write at least one server!");
      return;
    }
    var rule_source =  document.getElementById('rule_source').value;
    var rule_dest =  document.getElementById('rule_dest').value;

    AddRule(parseInt(rule_type), rule_ipaddr, parseInt(rule_subnetmask), rule_servers, rule_source, rule_dest).then(refreshRules);
    document.getElementById('rule_servers').value = '';
    document.getElementById('rule_source').value = '';
    document.getElementById('rule_dest').value = '';
  }

  function removeRule() {
    var rule_id = document.getElementById('rules_list_box').value;

    RemoveRule( parseInt(rule_id)).then(refreshRules);
  }

  function refreshRules() {
    RefreshRules();
  }
  function stopBalancer() {
    alert('Sorry! Not yet implemented. Restart the application.');
    //document.getElementById('buttonStart').disabled = false;
    //document.getElementById('buttonStop').disabled = true;
  }

  window.StartUpdateServedPagesEvent = () => {
    // noinspection JSIgnoredPromiseFromCall
    window.go.main.App.StartUpdateServedPagesEvent();
  }
    return (
        <div id="App">
            <div className="square">
              <h2>Rules config</h2>
              
              <select multiple id="rules_list_box">
                <option value="-1">Loading...</option>
              </select>
              
              <div className="formRow">
              <div><span>IP Addr filter (not implemented):</span></div>
              </div>

              <div className="formRow">
              <span>Backend servers (comma-separated)</span>
              <input className="servers" id="rule_servers" type="text" placeholder="google.it:80,microsoft.it:80" />
              </div>
              <span>Algorythm</span>
              <select id="rule_type">
                <option value="1">1: Round robin</option>
                <option value="2">2: IP hash (not implemented)</option>
              </select>

              <div className="formRow">
              <span>Source url filter (regex)</span>
              <input className="servers" id="rule_source" type="text" placeholder="(.*)" />
              </div>
              <div className="formRow">
              <span>Destination URL (with placeholders)</span>
              <input className="servers" id="rule_dest" type="text" placeholder="$1" />
              </div>

              <button className="btn" title="Add rule" onClick={addRule}>&#x2795;</button>
              <button className="btn"  title="Remove rule" onClick={removeRule}>&#x2796;</button>
              {/*  <button className="btn"  title="Remove rule" onClick={refreshRules}>&#x2796;</button> --> */}

            </div>
            <div className="square">
              <h2>Served connections: <div id="served_pages">0</div></h2>

                <div id="input" className="input-box">
                    <span className="hitCounter"> Listen to IP/port: </span>
                    <input id="bindPort" className="portInput" autoComplete="off" name="input" type="text" onChange={e => this.setState({ text: e.target.value })} defaultValue="8080"/>
                    <br /><br />
                    <button className="btn" id="buttonStart" onClick={startBalancer}>Start</button>

                    <button className="btn" id="buttonStop" onClick={stopBalancer}>Stop</button>
                    <div className="">
                    <div className="loader" id="loaderSpinner"></div>

                    </div>
                </div>
              </div>
              <div className="clearer"></div>

              <div className="square-log">
                <h2>Log</h2>
                <textarea className="textAreaField" id="textAreaLog"></textarea>
              </div>
          </div>
    )
}

export default App
