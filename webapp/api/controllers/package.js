import { db } from "../db.js";

export const allPakcages = (req, res) => {
    const token = req.cookies.access_token;
    if (!token) return res.status(401).json("No authenticated user");

    jwt.verify(token, "jwtkey", (err, userInfo) => {
        if (err) return res.status(403).json("JWT token is not valid");

        const userId = userInfo.username;

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

export const updateAddress = (req, res) => { }