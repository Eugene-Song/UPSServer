import express from 'express';
import authRouters from "./routers/authentication.js";
import packageRouters from "./routers/package.js";
import cookieParser from 'cookie-parser';


const app = express();

app.use(express.json());

app.use(cookieParser());

app.use("/api/auth", authRouters);
app.use("/api/package", packageRouters);

app.listen(8089, () => {
    console.log('Server started on port 8089');
})