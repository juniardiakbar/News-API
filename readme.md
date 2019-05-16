# News API - A JSON API for live news and blog articles

## Spesification
- Using Chi as a Router
- Parse database.json into 2 table
- First table is API Key tabel (keyTabele)
- Second table is articles table (articlesTable)
- API receive string q as query anf string sortBy for sorting 
Request:
```
http://localhost:8080/v2/everything?q=bitcoin&sortBy=publishedAt
```