# migrate-code

A code-based source for https://github.com/golang-migrate/migrate

## Example

```
package sqlitemigration

import (
	"github.com/warrenhodg/migrate-code"
)

func get() (map[string]string, error) {
	return map[string]string{
    "1_create channels table.up.sql": "CREATE TABLE IF NOT EXISTS channels ( id TEXT PRIMARY KEY, value TEXT NOT NULL );",
    "1_create channels table.down.sql": "DROP TABLE IF EXISTS channels",
	}, nil
}

func init() {
	code.Register("sqlite", get)
}
```
