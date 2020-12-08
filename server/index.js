const express = require('express');
const bodyParser = require('body-parser');
const { 
    create_record, 
    get_container, 
    get_all_containers, 
    change_record, 
    delete_record,
    get_record_history
} = require('./controller.js');

// Save our port
const port = process.env.PORT || 3001;

const app = express();
app.use(bodyParser.urlencoded({ extended: true }));
app.use(bodyParser.json());

app.get('/get/:id', (req, res) => {
    get_container(req, res);
});

app.get('/get_all', (req, res) => {
    get_all_containers(req, res);
});

app.put('/change/:key', (req, res) => {
    change_record(req, res);
});

app.post('/create', (req, res) => {
    create_record(req, res);
});

app.delete('/delete/:key', (req, res) => {
    delete_record(req, res);
});

app.get('/get_history/:key', (req, res) => {
    get_record_history(req, res);
});

// set up a static file server that points to the "client" directory
// app.use(express.static(path.join(__dirname, '../client')));

// Start the server and listen on port 
app.listen(port, () => console.log("Live on port: " + port));

