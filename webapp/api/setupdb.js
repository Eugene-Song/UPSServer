import mysql from 'mysql2/promise'
// MySQL connection configuration
const dbConfig = {
    host: 'localhost',
    user: 'root',
    password: 'Wadqq3.23',
    database: 'upsdb'
};



const sampleUserRecords1 = 'INSERT INTO package (packageID, status, currentX, currentY, destinationX, destinationY, username, date) VALUES (1, \'truck en route to warehouse\', 10, 20, 100, 200, \'wadqq\', \'2023-04-22 12:30:45\');'
const sampleUserRecords2 = 'INSERT INTO package (packageID, status, currentX, currentY, destinationX, destinationY, username, date) VALUES (2, \'truck en route to warehouse\', 75, 95, 75, 95, \'wadqq\', \'2023-04-20 16:15:30\');'

const sampleUserRecords3 = 'INSERT INTO package (packageID, status, currentX, currentY, destinationX, destinationY, username, date) VALUES (3, \'out for delivery\', 50, 45, 150, 180, \'wadqq\', \'2023-04-23 11:00:00\');'

const sampleUserRecords4 = 'INSERT INTO package (packageID, status, currentX, currentY, destinationX, destinationY, username, date) VALUES (4, \'truck waiting for package\', 30, 70, 200, 300, \'wadqq\', \'2023-04-24 14:20:10\');'

const sampleUserRecords5 = 'INSERT INTO package (packageID, status, currentX, currentY, destinationX, destinationY, username, date) VALUES (5, \'out for delivery\', NULL, NULL, 250, 350, \'wadqq\', \'2023-04-26 09:45:25\');'

const createDatabaseQuery = 'CREATE DATABASE upsdb;';
const useDatabaseQuery = 'USE upsdb;';
const createUserTableQuery = `
  CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255),
    email VARCHAR(255) UNIQUE,
    password VARCHAR(255) NOT NULL
  );
`;
const createPackageTableQuery = `
CREATE TABLE package (
  packageID BIGINT PRIMARY KEY,
  status VARCHAR(255) NOT NULL,
  currentX INT,
  currentY INT,
  destinationX INT,
  destinationY INT,
  username VARCHAR(255),
  item VARCHAR(255),
  date DATETIME NOT NULL,
  FOREIGN KEY (username) REFERENCES users(username)
);
`;
// SQL statements to drop existing tables
const dropPackageTableQuery = 'DROP TABLE IF EXISTS package;';
const dropUserTableQuery = 'DROP TABLE IF EXISTS users;';


async function checkDatabaseExists(databaseName, connection) {
    const checkDatabaseQuery = `
      SELECT COUNT(*)
      FROM information_schema.schemata
      WHERE schema_name = ?
    `;

    const [rows] = await connection.query(checkDatabaseQuery, [databaseName]);
    const count = Object.values(rows[0])[0];
    return count > 0;
}

async function setupDatabase() {
    try {
        // Connect to the MySQL server
        const connection = await mysql.createConnection(dbConfig);
        // Check if the database exists
        const databaseExists = await checkDatabaseExists('myDatabase', connection);

        if (databaseExists) {
            console.log('Database already exists');
        } else {
            console.log('Database does not exist');
            await connection.query(createDatabaseQuery);
            console.log('Database created');
        }
        // Execute the SQL statements
        await connection.query(dropPackageTableQuery);
        console.log('Package table dropped');
        // await connection.query(dropUserTableQuery);
        // console.log('User table dropped');
        await connection.query(useDatabaseQuery);
        console.log('Database selected');
        await connection.query(createPackageTableQuery).then(() => {
            console.log('Package table created');
            connection.query(sampleUserRecords1);
            connection.query(sampleUserRecords2);
            connection.query(sampleUserRecords3);
            connection.query(sampleUserRecords4);
            connection.query(sampleUserRecords5);
        });
        // Close the connection
        await connection.end();
    } catch (error) {
        console.error('Error setting up the database:', error);
    }
}

// Run the setupDatabase function
setupDatabase();
