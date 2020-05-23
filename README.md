# migrate-code

A code-based source for https://github.com/golang-migrate/migrate

## Example

```
package sqlitemigration

import (
	"github.com/warrenhodg/migrate-code"

	"github.com/golang-migrate/migrate/v4/source"
)

func get() ([]*source.Migration, error) {
	return []*source.Migration{
		{
			Version:    1,
			Identifier: "Create channels table",
			Direction:  source.Up,
			Raw:        "CREATE TABLE IF NOT EXISTS channels ( id TEXT PRIMARY KEY, value TEXT NOT NULL );",
		},
		{
			Version:    1,
			Identifier: "Delete channels table",
			Direction:  source.Down,
			Raw:        "DROP TABLE IF EXISTS channels",
		},
	}, nil
}

func init() {
	code.Register("sqlite", get)
}
```
