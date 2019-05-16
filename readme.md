# News API - A JSON API for live news and blog articles

## Spesification
- Using Chi as a Router
- Parse database.json into 2 table
- First table is API Key tabel (keyTabele)
- Second table is articles table (articlesTable)
- API receive string q as query and string sortBy for sorting 
Request:
```
http://localhost:8080/v2/everything?q=bitcoin&sortBy=publishedAt
```

## DATABASE MODEL

### Article
```
CREATE TABLE Articles(
    id serial PRIMARY KEY,
    source JSONB,
    author VARCHAR(500) NOT NULL,
    title VARCHAR(500) NOT NULL,
    description VARCHAR(500) NOT NULL,
    url VARCHAR(500) NOT NULL,
    urlToImage VARCHAR(500) NOT NULL,
    publishedAt TIMESTAMP NOT NULL,
    content VARCHAR(500) NOT NULL 
);
```

### Key
```
CREATE TABLE key(
    id serial PRIMARY KEY,
    key VARCHAR(500) NOT NULL
);
```
