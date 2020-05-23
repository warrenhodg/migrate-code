# migrate-code

A code-based source for https://github.com/golang-migrate/migrate

## Example

When running the example migration below, its source url would be `code:sqlite3`.
### Define the migrations for a particular database type

```
package sqlitemigration

import (
  "github.com/warrenhodg/migrate-code"

  // Some linters would complain about the need to put these anonymous imports
  // into a main package, rather than a database-specific migration.
  // You can choose, so long as they're imported somewhere
  _ "github.com/golang-migrate/migrate/v4/database/sqlite3"
  _ "github.com/mattn/go-sqlite3" // Simply import driver
)

func get() (map[string]string, error) {
  return map[string]string{
    "1_create channels table.up.sql": "CREATE TABLE IF NOT EXISTS channels ( id TEXT PRIMARY KEY, value TEXT NOT NULL );",
    "1_create channels table.down.sql": "DROP TABLE IF EXISTS channels",
  }, nil
}

func init() {
  code.Register("sqlite3", get)
}
```

### Run the migrations

```
import (
  "github.com/golang-migrate/migrate/v4"
)

func doMigrations() error {
  m, err := migrate.New("code:sqlite3", "sqlite3://./mydatabase.db")
  if err != nil {
    return err
  }

  err = m.Up()
  if err != nil && err != migrate.ErrNoChange {
    return err
  }

  return nil
}
```
