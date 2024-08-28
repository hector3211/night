![Night image](./public/night.png)

# Night

### Night a minimal utra fast database seed tool

Get up and running with night is seconds, the ultimate database seed tool out there. Night uses Golang for its core implementaion and some of your favortie database drivers including `sqlite3` and `postgres`.

### Install

```go
go install github.com/hector3211/night@latest
```

### Get Started

Create a file and call it `seed.sql` or `seed.go` fill it out to your liking and point Night to it.

### SQL Example

```sql
--  Sqlite Version
CREATE TABLE IF NOT EXISTS user (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Postgres Version
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_usersname ON users(username);

INSERT INTO users (username) VALUES ( 'maddog' );
```

### SQL Example

```go

type Users struct {
    ID night.Int `orm:"primary_key"`
    Name night.String
    Email night.VarChar `orm:"unique"`
    EmailVerified night.Bool `orm:"nullable"`
}

```

Then run this command to start seeding

```bash
night seed
```

### Database Drivers Supported

✅ sqlite3
✅ postgres
🟧 mysql

more incoming...
