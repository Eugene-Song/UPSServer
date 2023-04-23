"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const express_1 = __importDefault(require("express"));
const body_parser_1 = __importDefault(require("body-parser"));
const pino_1 = __importDefault(require("pino"));
const express_pino_logger_1 = __importDefault(require("express-pino-logger"));
const express_session_1 = __importDefault(require("express-session"));
const passport_1 = __importDefault(require("passport"));
const dbConfig = {
    host: 'localhost:3306',
    user: 'root',
    password: 'Wadqq3.23',
    database: 'upsDB'
};
const LocalStrategy = require('passport-local').Strategy;
const bcrypt = require('bcryptjs');
const mysql = require('mysql');
const app = (0, express_1.default)();
const port = 8095;
app.use(body_parser_1.default.json());
app.use(body_parser_1.default.urlencoded({ extended: true }));
const logger = (0, pino_1.default)({
    transport: {
        target: 'pino-pretty'
    }
});
app.use((0, express_pino_logger_1.default)({ logger }));
app.use(express_1.default.urlencoded({ extended: false }));
app.use((0, express_session_1.default)({
    secret: 'session_secret',
    resave: false,
    saveUninitialized: false,
}));
app.use(passport_1.default.initialize());
app.use(passport_1.default.session());
passport_1.default.serializeUser((user, done) => {
    logger.info("serializeUser " + JSON.stringify(user));
    done(null, user);
});
passport_1.default.deserializeUser((user, done) => {
    logger.info("deserializeUser " + JSON.stringify(user));
    done(null, user);
});
function checkAuthenticated(req, res, next) {
    if (!req.isAuthenticated()) {
        res.sendStatus(401);
        return;
    }
    next();
}
app.get('/register', (req, res) => {
    res.send('<form method="POST"><input type="text" name="username" placeholder="Username"><input type="password" name="password" placeholder="Password"><button type="submit">Register</button></form>');
});
app.post('/register', async (req, res) => {
    const { username, password } = req.body;
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
app.get('/login', (req, res) => {
    res.send('<form method="POST"><input type="text" name="username" placeholder="Username"><input type="password" name="password" placeholder="Password"><button type="submit">Log in</button></form>');
});
app.post('/login', passport_1.default.authenticate('local', {
    successRedirect: '/',
    failureRedirect: '/login'
}));
app.get('/logout', (req, res, next) => {
    req.logout(function (err) {
        if (err) {
            return next(err);
        }
        res.redirect('/');
    });
});
app.get('/', (req, res) => {
    if (req.isAuthenticated()) {
        res.send(`Welcome, ${req.user.username}! <a href="/logout">Logout</a>`);
    }
    else {
        res.send('Please <a href="/login">log in</a> or <a href="/register">register</a> to access this page.');
    }
});
app.get('/protected', checkAuthenticated, (req, res) => {
    res.send(`This is a protected page. Welcome, ${req.user.username}!`);
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
        }
        else {
            console.log('Connected to the MySQL database.');
        }
    });
    connection.on('error', (error) => {
        console.error('Database error:', error);
        if (error.code === 'PROTOCOL_CONNECTION_LOST') {
            createConnection();
        }
        else {
            throw error;
        }
    });
    return connection;
}
const db = createConnection();
app.use(express_1.default.urlencoded({ extended: false }));
app.use((0, express_session_1.default)({
    secret: 'your_session_secret',
    resave: false,
    saveUninitialized: false,
}));
app.use(passport_1.default.initialize());
app.use(passport_1.default.session());
passport_1.default.use(new LocalStrategy((username, password, done) => {
    const query = 'SELECT * FROM users WHERE username = ?';
    db.query(query, [username], (error, results) => {
        if (error)
            return done(error);
        if (results.length === 0) {
            return done(null, false, { message: 'Incorrect username.' });
        }
        const user = results[0];
        bcrypt.compare(password, user.password, (error, isMatch) => {
            if (error)
                return done(error);
            if (!isMatch) {
                return done(null, false, { message: 'Incorrect password.' });
            }
            return done(null, user);
        });
    });
}));
passport_1.default.serializeUser((user, done) => {
    done(null, user.id);
});
passport_1.default.deserializeUser((id, done) => {
    const query = 'SELECT * FROM users WHERE id = ?';
    db.query(query, [id], (error, results) => {
        if (error)
            return done(error);
        done(null, results[0]);
    });
});
//# sourceMappingURL=server.js.map