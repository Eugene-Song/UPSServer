import express from "express";
import { allPakcages, detailPackage, updateAddress } from "../controllers/package.js";

const router = express.Router();

router.get("/all", allPakcages)
router.get("/:id", detailPackage)
router.put("/:id", updateAddress)

export default router;