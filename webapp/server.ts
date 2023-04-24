import express, { NextFunction, Request, Response } from 'express'
import bodyParser from 'body-parser'
import pino from 'pino'
import expressPinoLogger from 'express-pino-logger'
import session from 'express-session'
import passport from 'passport'
import path from 'path';
import { Server } from 'http';


// set up mysql
const dbConfig = {
  host: 'localhost',
  user: 'root',
  password: 'Wadqq3.23',
  database: 'upsDB'
};

const LocalStrategy = require('passport-local').Strategy;
const bcrypt = require('bcryptjs');
const mysql = require('mysql');

// set up Express
const app = express()
app.use(bodyParser.json())
app.use(bodyParser.urlencoded({ extended: true }))
// app.use(express.static(path.join(__dirname, 'ui')));


// set up Pino logging
const logger = pino({
  transport: {
    target: 'pino-pretty'
  }
})

app.use(expressPinoLogger({ logger }))
app.use(session({
  // some seceret?
  secret: 'session_secret',
  resave: false,
  saveUninitialized: false,
}));

app.use(passport.initialize());
app.use(passport.session());

passport.serializeUser((user: any, done: any) => {
  logger.info("serializeUser " + JSON.stringify(user))
  done(null, user)
})

passport.deserializeUser((user: any, done: any) => {
  logger.info("deserializeUser " + JSON.stringify(user))
  done(null, user)
})

function checkAuthenticated(req: Request, res: Response, next: NextFunction) {
  if (!req.isAuthenticated()) {
    res.sendStatus(401)
    return
  }
  next()
}

// Registration route
app.get('/register', (req, res) => {
  res.send('<form method="POST"><input type="text" name="username" placeholder="Username"><input type="password" name="password" placeholder="Password"><button type="submit">Register</button></form>');
});

app.post('/register', async (req, res) => {
  const { username, password } = req.body;

  // Check if the user already exists
  const query = 'SELECT * FROM users WHERE username = ?';
  db.query(query, [username], async (error, results) => {
    if (error) {
      console.error('Error executing query:', error);
      res.status(500).send('Error executing query');
      return;
    }

    if (results.length > 0) {
      res.status(400).send('User already exists');
      return;
    }

    // Hash the password and insert the new user
    const hashedPassword = await bcrypt.hash(password, 10);
    const insertQuery = 'INSERT INTO users (username, password) VALUES (?, ?)';
    db.query(insertQuery, [username, hashedPassword], (error, results) => {
      if (error) {
        console.error('Error executing query:', error);
        res.status(500).send('Error executing query');
        return;
      }

      res.redirect('/login');
    });
  });
});

// Login route
app.get('/login', (req, res) => {
  res.send('<form method="POST"><input type="text" name="username" placeholder="Username"><input type="password" name="password" placeholder="Password"><button type="submit">Log in</button></form>');
});

app.post('/login', passport.authenticate('local', {
  successRedirect: '/',
  failureRedirect: '/login'
}));

// Logout route
app.get('/logout', (req, res, next) => {
  req.logout(function(err) {
    if (err) { return next(err); }
    res.redirect('/');
  });
});

// Home route
app.get('/', (req, res) => {
  if (req.isAuthenticated()) {
    res.send(`Welcome, ${req.user.username}! <a href="/logout">Logout</a>`);
    // add some button to go to the tracking page

  } else {
    res.send('Please <a href="/login">log in</a> or <a href="/register">register</a> to access this page.');
    // TODO: use tracking number search for packages

  }
});

// Example of a protected route
app.get('/protected', checkAuthenticated, (req, res) => {
  res.send(`This is a protected page. Welcome, ${req.user.username}!`);
});

// /api/packages route, which returns the packages for the logged in user
app.get('/api/packages', checkAuthenticated, (req, res) => {
  const userId = req.user.id;

  // Replace this query with the actual query for your database
  const query = 'SELECT * FROM package WHERE user_id = ?';

  db.query(query, [userId], (error, results) => {
    if (error) {
      console.error('Error executing query:', error);
      res.status(500).send('Error executing query');
      return;
    }

    res.json(results);
  });
});

// /api/packages/:packageId route, which returns the package with the given ID for the logged in user's specific package
app.get('/api/packages/:packageId', checkAuthenticated, (req, res) => {
  const userId = req.user.id;
  const packageId = req.params.packageId;

  console.log('userId', userId);
  console.log('packageId', packageId);
  // Replace this query with the actual query for your database
  const query = 'SELECT * FROM package WHERE packageid = ? AND user_id = ?';

  db.query(query, [packageId, userId], (error, results) => {
    if (error) {
      console.error('Error executing query:', error);
      res.status(500).send('Error executing query');
      return;
    }

    if (results.length === 0) {
      res.status(404).send('Package not found');
      return;
    }

    res.json(results[0]);
  });
});

// /api/packages/:packageId/update-address route, which updates the address of the package with the given ID for the logged in user's specific package
app.put('/api/packages/:packageId/update-address', checkAuthenticated, (req, res) => {
  const userId = req.user.id;
  const packageId = req.params.packageId;
  const newAddress = req.body.targetaddr;

  if (!newAddress) {
    res.status(400).send('Missing target address');
    return;
  }

  // Replace this query with the actual query for your database
  const query = 'UPDATE package SET targetaddr = ? WHERE packageid = ? AND user_id = ?';

  db.query(query, [newAddress, packageId, userId], (error, results) => {
    if (error) {
      console.error('Error executing query:', error);
      res.status(500).send('Error executing query');
      return;
    }

    if (results.affectedRows === 0) {
      res.status(404).send('Package not found');
      return;
    }

    res.status(200).send('Address updated successfully');
  });
});



const PORT = process.env.PORT || 3000;
app.listen(PORT, () => {
  console.log(`Server is running on port ${PORT}`);
});




function createConnection() {
  const connection = mysql.createConnection(dbConfig);

  connection.connect((error) => {
    if (error) {
      console.error('Error connecting to the database:', error);
      setTimeout(createConnection, 2000);
    } else {
      console.log('Connected to the MySQL database.');
    }
  });

  connection.on('error', (error) => {
    console.error('Database error:', error);
    if (error.code === 'PROTOCOL_CONNECTION_LOST') {
      createConnection();
    } else {
      throw error;
    }
  });

  return connection;
}


const db = createConnection();

app.use(express.urlencoded({ extended: false }));

app.use(session({
  secret: 'your_session_secret',
  resave: false,
  saveUninitialized: false,
}));

app.use(passport.initialize());
app.use(passport.session());

passport.use(new LocalStrategy(
  (username, password, done) => {
    const query = 'SELECT * FROM users WHERE username = ?';
    db.query(query, [username], (error, results) => {
      if (error) return done(error);
      if (results.length === 0) {
        return done(null, false, { message: 'Incorrect username.' });
      }
      
      const user = results[0];
      bcrypt.compare(password, user.password, (error, isMatch) => {
        if (error) return done(error);
        if (!isMatch) {
          return done(null, false, { message: 'Incorrect password.' });
        }
        return done(null, user);
      });
    });
  }
));

passport.serializeUser((user, done) => {
  done(null, user.id);
});

passport.deserializeUser((id, done) => {
  const query = 'SELECT * FROM users WHERE id = ?';
  db.query(query, [id], (error, results) => {
    if (error) return done(error);
    done(null, results[0]);
  });
});
