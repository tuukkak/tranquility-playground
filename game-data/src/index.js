const express = require("express");
const app = express();

app.get("/", (req, res) => {
    console.log("Request received!");
    res.send("OK");
});

app.listen(3001, () => {
    console.log("Server started");
});
