const express = require("express");
const app = express();
const https = require("https");

app.get("/", (req, res) => {
  setTimeout(() => {
    res.send("Hello World");
  }, 1000);
});

app.get("/http-request/", (req, res) => {
  const endpoint = "https://www.microsoft.com/";
  https.get(endpoint, (response) => {
    response.on("data", () => {});

    response.on("error", (err) => {
      res.send(`Encountered error while making HTTPS request: ${err}`);
    });

    response.on("end", () => {
      res.send(`Successfully reached ${endpoint}.`);
    });
  });
});

app.listen(8080, () => {
  console.log(`Listening for requests on http://localhost:8080`);
});
