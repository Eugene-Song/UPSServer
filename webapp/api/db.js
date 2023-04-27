import mysql from 'mysql2'

export const db = mysql.createConnection({
    host: 'localhost',
    user: 'root',
    password: 'Wadqq3.23',
    database: 'upsdb'
})