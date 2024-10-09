const { APP_HOST = 'localhost' } = process.env 

export default async function Page() {
  let data = await fetch(`http://${APP_HOST}:8081/books`)
  let books = await data.json()
  return (
    <div>
      <h1>Bookstore</h1>
      <h2>http://{APP_HOST}:8081/books</h2>
      <ul>
        {books.map((book) => (
          <li key={book.id}>
            {book.title}
            <span>
              &nbsp;{book.isAvailable == true ? "✅": "❌" }
            </span>
          </li>
        ))}
      </ul>
    </div>
  )
}
