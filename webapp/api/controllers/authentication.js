import { db } from "../db.js";
import jwt from "jsonwebtoken";
import bcrypt from "bcryptjs";

export const login = (req, res) => {
    const q = "SELECT * FROM users WHERE username = ?";


    db.query(q, [req.body.username], (err, data) => {
        // if query error, return 500
        if (err) return res.status(500).json(err);
        // if user does not exist, return 404
        if (data.length === 0) return res.status(404).json("User not exists!");
        // compare the input password
        const isPasswordCorrect = bcrypt.compareSync(
        req.body.password,
        data[0].password
        );
        if (!isPasswordCorrect) {
            // if the password is not correct, return 400
            return res.status(400).json("Wrong password / username!!!");
        }
        const token = jwt.sign({ id: data[0].id, username: data[0].username, email: data[0].email }, "secretkey");
        const { password, ...other } = data[0];

        res
        .cookie("jwttoken", token, {
            httpOnly: true,
        })
        .status(200)
        .json(other);
    });
}

export const register = (req, res) => {
    const q = "SELECT * FROM users WHERE email = ? OR username = ?";

    db.query(q, [req.body.email, req.body.username], (err, data) => {
      if (err) return res.status(500).json(err);
      if (data.length) return res.status(409).json("User already exists!");
  
      //Hash the password, add salt, and create a user
      const salt = bcrypt.genSaltSync(10);
      const hash = bcrypt.hashSync(req.body.password, salt);

      const q = "INSERT INTO users(`username`,`name`, `email`,`password`) VALUES (?)";
      const values = [req.body.username, req.body.name, req.body.email, hash];
  
      db.query(q, [values], (err, data) => {
        if (err) return res.status(500).json(err);
        return res.status(200).json("User has been created.");
      });
    });
}

export const logout = (req, res) => {
    res.clearCookie("jwttoken",{
    sameSite:"none",
    secure:true
}).status(200).json("User is logged out.");
}