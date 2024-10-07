const opentelemetry = require("@opentelemetry/api");
const tracer = opentelemetry.trace.getTracer("node-app");

const express = require("express");
const app = express();
const axios = require("axios");
const cors = require("cors");

const allowedOrigins = ["http://localhost:3000", "http://localhost:8080", "http://host.docker.internal:3000"];
app.use(
  cors({
    origin: function (origin, callback) {
      if (!origin) return callback(null, true);
      if (allowedOrigins.indexOf(origin) === -1) {
        var msg = "The CORS policy for this site does not " + "allow access from the specified Origin.";
        return callback(new Error(msg), false);
      }
      return callback(null, true);
    },
  })
);

const { AVAILABILITY_HOST = "localhost" } = process.env;

app.get("/", (req, res) => {
  const span = tracer.startSpan("Hello World");
  span.end();
  res.send("Hello World");
});

app.get("/books", booksListHandler);

async function booksListHandler(req, res) {
  const span = tracer.startSpan("Books List");
  const books = await getAvailableBooks();
  span.setAttribute("books.list.count", books.length);
  span.end();

  res.json(books);
}

async function getAvailableBooks() {
  const books = getBooks();
  const availableBooks = await Promise.all(
    books.map(async (book) => {
      const endpoint = `http://${AVAILABILITY_HOST}:8082/availability/${book.id}`;
      const {
        data: { isAvailable },
      } = await axios.get(`${endpoint}`);
      return { ...book, isAvailable };
    })
  );
  return availableBooks;
}

function getBooks() {
  return [
    {
      id: 1,
      title: "Harry Potter",
    },
    {
      id: 2,
      title: "Foundation",
    },
    {
      id: 3,
      title: "Moby Dick",
    },
  ];
}

app.listen(8081, () => {
  console.log(`Listening for requests on http://localhost:8081`);
});
