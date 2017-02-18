const express = require('express');
const app = express();
const MongoClient = require('mongodb').MongoClient;
const bodyParser = require('body-parser');

app.use(bodyParser.json())

function execApp(db) {
    const col = db.collection("testCollection");

    col.remove({}, () => {});

    col.insert([
        {data: 'someData'},
        {data: 'similarData'},
        ], (err, result) => {
            if (err) {
                console.log('could not insert data');
            }
        });

        app.get('/', (req, res) => {
            res.sendFile('/tmp/images/eduard.png')
        });

        app.get('/data/', (req, res) => {
            col.find({}).toArray((err, result) => {
                if (err) {
                    console.log('could not fetch data');
                    res.json({succes: false});
                }
                res.json(result);
            });
        });

        app.post('/data/', (req, res) => {
            col.insert(req.body, (err, result) => {
                if (err) {
                    console.log('could not insert data');
                    res.json({succes: false});
                }
                res.json({success: true});
            });
        });

        app.listen(8080, () => {
            console.log("Server started on port: 8080");
        });
}

const connectAndRetry = () => {
    return MongoClient.connect("mongodb://db:27017/test", (err, db) => {
        if (err) {
            console.log('error connecting to db...retrying in 2 sec...');
            setTimeout(connectAndRetry, 2000);
        } else {
            execApp(db);
        }
    });
};

connectAndRetry();
