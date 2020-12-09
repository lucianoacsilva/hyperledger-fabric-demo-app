import React, { Component } from 'react';
import { toast } from 'react-toastify';
import Notification from './Notification'; 
import './App.css';

class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      sampleKey: '',
      newHolderSampleKey: '',
      newHolderName: null,
      newForce: null,
      newStretching: null,
      newSampleKey: '',
      containerDescription: '',
      holderName: '',
      allContainers: [],
      force: '',
      stretching: ''
    };

    this.changeHolder = this.changeHolder.bind(this);
    this.createRecord = this.createRecord.bind(this);
    this.handleTextChange = this.handleTextChange.bind(this);
    this.queryContainer = this.queryContainer.bind(this);
    this.queryAllContainers = this.queryAllContainers.bind(this);
  }

  notifySuccess = message => toast.success(message);
  notifyError = message => toast.error(message);

  createRecord(event) {
    event.preventDefault();
    const { 
      holderName,
      force,
      stretching 
    } = this.state;
    if (force && stretching && holderName) {
      fetch('create', {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          holder: holderName,
          force,
          stretching
        }),
      })
      .then(response => response.json())
      .then(data => {
        if (data.success && data.result) {
          this.notifySuccess('New record was created');
        } else {
          console.error(data.error);
          this.notifyError('Something went wrong');
        }
      });
    }
  }

  changeHolder(event) {
    event.preventDefault();
    const { 
      allContainers, 
      sampleKey, 
      newHolderName,
      newForce,
      newStretching
    } = this.state;

    const reqBody = {};

    if (newHolderName) {
      reqBody.holder = newHolderName
    }

    if (newForce) {
      reqBody.force = newForce
    }

    if (newStretching) {
      reqBody.stretching = newStretching
    }

    fetch(`change/${encodeURIComponent(sampleKey)}`, {
      method: "PUT",
      
      headers: {
        "Content-Type": "application/json",
      },
        body: JSON.stringify({
          id: newHoldersampleKey,
          holder: newHolderName
        }),
      })
      .then(response => response.json())
      .then(data => {
        if (data.success && data.result) {
        this.notifySuccess('Sample was changed');
      } else {
        console.error(data.error);
        this.notifyError('Something went wrong');
      }
    });
  }
  }

  handleTextChange = event => {
    this.setState({ [event.target.id]: event.target.value });
  };

  queryContainer(event) {
    event.preventDefault();
    const { sampleKey } = this.state;
    if (sampleKey) {
      fetch(`get/${encodeURIComponent(sampleKey)}`)
        .then(response => response.json())
        .then(data => {
          if (data.success && data.result) {
            const result = JSON.parse(data.result)

            this.setState({ 
              allContainers: [{ Key: sampleKey, Record: result.sample }]
            });
            console.log("Wallet Address", result.wallet)
            console.log("MAM Root", result.mamstate.root)
            console.log("MAM payload: ")
            console.log(result.messages)
            console.log("======================")
          } else {
            console.error(data.error);
          }
      });
    }
  }

  queryAllContainers(event) {
    event.preventDefault();
    fetch('get_all')
      .then(response => response.json())
      .then(data => {
        if (data.success && data.result) {
          this.setState({ allContainers: data.result });
        } else {
          console.error(data.error);
        }
      });
  }

  render() {
    return (
      <div className="App">
        <header>
          <div id="left_header">Hyperledger Fabric Demo Application</div>
          <i id="right_header">Example Blockchain Application for Hyperledger Fabric</i>
        </header>
        <div className="queryContainer">
          <form onSubmit={this.queryContainer}>
            <label>Query a Specific Sample</label><br />
            Enter a sample ID: <br />
            <input
              id="sampleKey"
              type="string"
              placeholder="Ex: Sample_1"
              value={this.state.sampleKey}
              onChange={this.handleTextChange}
            />
            <br />
            <button type="submit" className="btn btn-primary">Query Sample Record</button>
          </form>
        </div>
        <br />
        <br />

        <div className="queryAllContainers">
          <div className="form-group">
            <label>Query All Containers</label><br />
            <button type="button" className="btn btn-primary" onClick={this.queryAllContainers}>Query All Containers</button>
          </div>
 
          {
            this.state.allContainers.length ? (
              <table id="all_containers" className="table" align="center">
                <thead>
                  <tr>
                    <th>ID</th>
                    <th>Timestamp</th>
                    <th>Holder</th>
                    <th>Stretching</th>
                    <th>Force</th>
                  </tr>
                </thead>
                <tbody>
                  {
                    this.state.allContainers
                    .sort((a, b) => parseFloat(a.Key) - parseFloat(b.Key))
                    .map(container => (
                      <tr key={container.Key}>
                        <td>{container.Key}</td>
                        <td>{container.Record.timestamp}</td>
                        <td>{container.Record.holder}</td>
                        <td>{container.Record.stretching}</td>
                        <td>{container.Record.force}</td>
                      </tr>
                    ))
                  }
                </tbody>
              </table>
            ) : null
          }
        </div>

        <br />
        <br />

        <div className="createRecord">
          <form onSubmit={this.createRecord}>
            <label>Create Sample Record</label>
            <br />
            Enter sample force:
            <input
              className="form-control" 
              id="force"
              name="force" 
              type="text" 
              placeholder="Ex: 11" 
              value={this.state.force}
              onChange={this.handleTextChange}
            />

            Enter sample stretching: 
            <input 
              className="form-control" 
              id="stretching"
              name="stretching" 
              type="number" 
              placeholder="Ex: 7"
              value={this.state.stretching}
              onChange={this.handleTextChange}
            /> 
            
            Enter name of holder: 
            <input 
              className="form-control" 
              id="holderName"
              name="holderName" 
              type="text" 
              placeholder="Ex: A" 
              value={this.state.holderName}
              onChange={this.handleTextChange}
            />
            <button type="submit" className="btn btn-primary">Create record</button>
          </form>
        </div>

        <br />
        <br />

        <div className="changeContainerHolder">
          <form onSubmit={this.changeHolder}>
            <label>Change Sample Holder</label><br />

            Enter a sample ID:
            <input
              className="form-control"
              id="sampleKey"
              type="string"
              placeholder="Ex: Sample_1"
              value={this.state.sampleKey}
              onChange={this.handleTextChange}
            />
            Enter name of new holder:
            <input
              className="form-control"
              id="newHolderName"
              name="newHolderName"
              placeholder="Ex: Barry"
              type="text"
              value={this.state.newHolderName}
              onChange={this.handleTextChange}
            />
            <button type="submit" className="btn btn-primary">Change</button>
          </form>
        </div>
        
        <br /><br /><br />
        <Notification />
      </div>
    );
  }
}

export default App;
