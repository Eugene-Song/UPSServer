import mysql from 'mysql2/promise';

// MySQL connection configuration
const dbConfig = {
    host: 'localhost',
    user: 'root',
    password: 'Wadqq3.23',
    database: 'upsdb'
};

const sampleUserRecords1 = 'INSERT INTO users (username, password) VALUES (\'Yuxin\', \'$2a$10$SZqPcbLFSM60JSlS6SGD4OAIT0F2v22VGo/P81JrnHDCEaaWzJFv2\')';
const sampleRecords1 = 'INSERT INTO package (status, tracking_number, date, targetaddr, user_id) VALUES (\'In Transit\', \'ABC1234567\', \'2023-04-01\', \'123 Main St, New York, NY 10001\', 1)';
const sampleRecords2 = 'INSERT INTO package (status, tracking_number, date, targetaddr, user_id) VALUES (\'Delivered\', \'DEF1234567\', \'2023-03-15\', \'456 Elm St, Los Angeles, CA 90001\', 1)';
const sampleRecords3 = 'INSERT INTO package (status, tracking_number, date, targetaddr, user_id) VALUES (\'Returned\', \'GHI1234567\', \'2023-04-10\', \'789 Oak St, Chicago, IL 60601\', 1)';
const sampleRecords4 = 'INSERT INTO package (status, tracking_number, date, targetaddr, user_id) VALUES (\'Cancelled\', \'JKL1234567\', \'2023-04-05\', \'321 Pine St, San Francisco, CA 94101\', 1)';
const sampleRecords5 = 'INSERT INTO package (status, tracking_number, date, targetaddr, user_id) VALUES (\'In Transit\', \'MNO1234567\', \'2023-04-20\', \'654 Maple St, Miami, FL 33101\', 1)';

// SQL statements to create database and tables
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
    FOREIGN KEY (username) REFERENCES users(username)
    date DATETIME NOT NULL,
  );
`;
// SQL statements to drop existing tables
const dropPackageTableQuery = 'DROP TABLE IF EXISTS package;';
const dropUserTableQuery = 'DROP TABLE IF EXISTS users;';

async function checkDatabaseExists(databaseName: string, connection: mysql.Connection): Promise<boolean> {
    const checkDatabaseQuery = `
      SELECT COUNT(*)
      FROM information_schema.schemata
      WHERE schema_name = ?
    `;

    const [rows] = await connection.query(checkDatabaseQuery, [databaseName]);
    const count = Object.values(rows[0])[0] as number;
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
        await connection.query(dropUserTableQuery);
        console.log('User table dropped');


        await connection.query(useDatabaseQuery);
        console.log('Database selected');
        await connection.query(createUserTableQuery).then(() => {
            console.log('User table created');
            connection.query(sampleUserRecords1);
            connection.query(createPackageTableQuery);
            connection.query(sampleRecords1);
            connection.query(sampleRecords2);
            connection.query(sampleRecords3);
            connection.query(sampleRecords4);
            connection.query(sampleRecords5);
        });
        
        
       
        // Close the connection
        await connection.end();
    } catch (error) {
        console.error('Error setting up the database:', error);
    }
}

// Run the setupDatabase function
setupDatabase();
