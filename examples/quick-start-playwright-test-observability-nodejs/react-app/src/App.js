import './App.css';
import { useState, useEffect } from 'react'
const { APP_HOST = 'localhost' } = process.env 

function App() {
  const [books, setBooks] = useState([]);
  useEffect(() => {
    fetch(`http://${APP_HOST}:8081/books`)
      .then((res) => {
        return res.json();
      })
      .then((data) => {
        console.log(data);
        setBooks(data);
      })
      .catch((err) => {
        console.log(err)
      })
  }, []);
  return (
    <div>
      <h1>Bookstore</h1>
      {books.map((book) => (
        <li key={book.id}>
          {book.title}
          <span>
            &nbsp;{book.isAvailable === true ? "✅": "❌" }
          </span>
        </li>
      ))}
    </div>
  );
}

export default App;
