const AWSXRay = require("aws-xray-sdk");
const XRayExpress = AWSXRay.express;
const express = require("express");

AWSXRay.setDaemonAddress("xray-daemon:2000");

// Capture all AWS clients we create
const AWS = AWSXRay.captureAWS(require("aws-sdk"));
AWS.config.update({
  region: process.env.AWS_REGION || "us-west-2",
});

// Capture all outgoing https requests
AWSXRay.captureHTTPsGlobal(require("https"));
const https = require("https");

const app = express();
const port = 3000;

app.use(XRayExpress.openSegment("Tracetest"));

app.get("/", (req, res) => {
  const seg = AWSXRay.getSegment();
  const sub = seg.addNewSubsegment("customSubsegment");
  setTimeout(() => {
    sub.close();
    res.sendFile(`${process.cwd()}/index.html`);
  }, 500);
});

app.get("/http-request/", (req, res) => {
  const endpoint = "https://amazon.com/";
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

app.use(XRayExpress.closeSegment());
app.listen(port, () => console.log(`Example app listening on port ${port}!`));
