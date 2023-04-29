import { db } from "../db.js";
import jwt from 'jsonwebtoken';
import net from 'net';

export const allPakcages = (req, res) => {
    const token = req.cookies.access_token;
    if (!token) return res.status(401).json("No authenticated user");

    jwt.verify(token, "jwtkey", (err, userInfo) => {
        if (err) return res.status(403).json("JWT token is not valid");

        // Replace this query with the actual query for your database
        const query = 'SELECT * FROM package WHERE username = ?';

        db.query(query, [userInfo.username], (error, results) => {
            if (error) {
                console.error('Error executing query:', error);
                res.status(500).send('Error executing query');
                return;
            }

            res.json(results);
        });
    });
}

export const detailPackage = (req, res) => {
  const q = "SELECT * FROM package WHERE `packageID` = ?";

  db.query(q, [req.params.id], (err, data) => {
    if (err) return res.status(500).json(err);

    return res.status(200).json(data[0]);
  });
}

export const updateAddress = (req, res) => { 
  const token = req.cookies.access_token;
  if (!token) return res.status(401).json("Not authenticated!");

  jwt.verify(token, "jwtkey", (err, userInfo) => {
    if (err) return res.status(403).json("Token is not valid!");

    const shipID = req.params.id;
    console.log("shipID: dwadad", shipID)

    const client = new net.Socket();

    const message = {
      ShipID: shipID,
      X: req.body.X,
      Y: req.body.Y,
    };
    

    client.connect(8090, 'localhost', () => {
      const messageString = JSON.stringify(message);
      const messageLength = Buffer.byteLength(messageString);

      // Create a buffer to send the message length and data
      const buffer = Buffer.alloc(4 + messageLength);
      buffer.writeInt32BE(messageLength, 0);
      buffer.write(messageString, 4);

      client.write(buffer);
    });

    client.on('data', (data) => {
      // Data received from the server
      const message = data.toString();
      console.log('Received message:', message);
      // End the connection after receiving the data
      client.end();
      return res.status(200).json(message)
    });

    client.on('close', () => {
      console.log('Connection closed');
    });

  });
}