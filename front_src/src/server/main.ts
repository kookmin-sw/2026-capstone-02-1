
import express from "express";
import ViteExpress from "vite-express";

import multer from "multer";

import fs from "fs";
import path from "path";
import url from "url";
import os from "os";

import { exec } from "child_process";

const app = express();
const port = 3000;

// Specifiy file paths
const __filename = url.fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const UPLOAD_DIR = path.resolve(__dirname, "../../go");
const PROGRAM_DIR = path.resolve(__dirname, "../../bin");
const OUTPUT_DIR = path.resolve(__dirname, "../../json");

const UPLOAD_NAME = "source.go";
const PROGRAM_NAME = (os.platform() == "win32") ? "traceinspector.exe" : "traceinspector.o";

const OUTPUT_NAME = "output.json";

const UPLOAD_PATH = path.join(UPLOAD_DIR, UPLOAD_NAME);
const PROGRAM_PATH = path.join(PROGRAM_DIR, PROGRAM_NAME);
const OUTPUT_PATH = path.join(OUTPUT_DIR, OUTPUT_NAME);

if (!fs.existsSync(UPLOAD_DIR))
    fs.mkdirSync(UPLOAD_DIR, { recursive: true });

if (!fs.existsSync(PROGRAM_DIR))
    fs.mkdirSync(PROGRAM_DIR, { recursive: true });

if (!fs.existsSync(OUTPUT_DIR))
    fs.mkdirSync(OUTPUT_DIR, { recursive: true });

const storage = multer.diskStorage({
    destination: (_req, _file, cb) => cb(null, UPLOAD_DIR),
    filename: (_req, file, cb) => {
        if (file.originalname.toLowerCase().endsWith(".go")) {
            cb(null, UPLOAD_NAME);
        } else {
            cb(null, "rubbish")
        }
    },
});

const upload = multer({ storage });

// Get file from client
app.post("/upload", upload.single("file"), (req, res) => {
    if (!req.file)
        return res.status(400).json({ error: "No file uploaded" });

    const savedName = req.file.filename;
    const savedPath = req.file.path;

    if (!savedName.toLowerCase().endsWith(".go")) {
        fs.unlink(savedPath, () => { });
        return res.status(400).json({ error: "Only .go files are allowed" });
    }

    res.json({ savedAs: savedName, path: savedPath });
});

// Run inspection and return to client
app.get("/run", async (_req, res) => {
    try {
        await new Promise<void>((resolve, reject)=>  {
            exec(`${PROGRAM_PATH} --print-cfg-json --gofile ${UPLOAD_PATH} > ${OUTPUT_PATH} 2>&1`, null, (error, stdout, stderr) => {
                if (error) {
                    reject(new Error(`${stderr || error.message}`));
                } else {
                    resolve();
                }
            });
        });

        if (!fs.existsSync(OUTPUT_PATH)) {
            return res.status(400).json({ error: "Program did not generate output.json" });
        }

        const jsonData = fs.readFileSync(OUTPUT_PATH, "utf-8");
        const parsedData = JSON.parse(jsonData);

        res.json({
            savedAs: OUTPUT_NAME,
            output: parsedData
        });

    } catch (err: any) {
        res.status(500).json({
            error: "Execution failed",
            message: err.message
        });
    }
});

ViteExpress.listen(app, port, () =>
    console.log(`Server is listening on ${port}...`),
);
