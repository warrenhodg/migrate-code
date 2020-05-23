package code

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4/source"
)

// A Migration returns an array of migrations to run
type Migration func() (map[string]string, error)

// Errors
var (
	ErrFactoryNotFound = errors.New("migrate-code factory not found")
	ErrNotFound        = errors.New("not found")
)

// Factories
var (
	factories = map[string]Migration{}
)

// Register a function callback
func Register(key string, c Migration) {
	factories[key] = c
}

// Code handles migrations returned from a function
type Code struct {
	ms         *source.Migrations
	path       string
	migrations map[string]string
}

// Open calls the relevant function to get its migrations
// The url passed in should be of the form "code://xyz"
// where xyz is a registered code factory in your own code.
func (c *Code) Open(fullURL string) (source.Driver, error) {
	u, err := url.Parse(fullURL)
	if err != nil {
		return nil, err
	}

	migrateFactory := u.Opaque

	f, found := factories[migrateFactory]
	if !found {
		return nil, fmt.Errorf("%w: %s", ErrFactoryNotFound, migrateFactory)
	}

	migrations, err := f()
	if err != nil {
		return nil, err
	}

	ms := source.NewMigrations()
	for k := range migrations {
		m, err := source.Parse(k)
		if err != nil {
			return nil, err
		}

		if !ms.Append(m) {
			return nil, source.ErrDuplicateMigration{}
		}
	}

	return &Code{
		ms:         ms,
		migrations: migrations,
		path:       migrateFactory,
	}, nil
}

// First gets the first migration
func (c *Code) First() (uint, error) {
	v, ok := c.ms.First()
	if !ok {
		return 0, &os.PathError{
			Op:   "first",
			Path: c.path,
			Err:  os.ErrNotExist,
		}
	}
	return v, nil
}

// Prev gets the previous migration
func (c *Code) Prev(version uint) (uint, error) {
	v, ok := c.ms.Prev(version)
	if !ok {
		return 0, &os.PathError{
			Op:   fmt.Sprintf("prev for version %v", version),
			Path: c.path,
			Err:  os.ErrNotExist,
		}
	}
	return v, nil
}

// Next gets the next migration
func (c *Code) Next(version uint) (uint, error) {
	v, ok := c.ms.Next(version)
	if !ok {
		return 0, &os.PathError{
			Op:   fmt.Sprintf("next for version %v", version),
			Path: c.path,
			Err:  os.ErrNotExist,
		}
	}
	return v, nil
}

// ReadUp reads the up function
func (c *Code) ReadUp(version uint) (io.ReadCloser, string, error) {
	if m, ok := c.ms.Up(version); ok {
		sql, found := c.migrations[m.Raw]
		if !found {
			return nil, "", fmt.Errorf("migrate-code readup: %w", ErrNotFound)
		}

		return ioutil.NopCloser(strings.NewReader(sql)), m.Identifier, nil
	}

	return nil, "", &os.PathError{
		Op:   fmt.Sprintf("read version %v", version),
		Path: c.path,
		Err:  os.ErrNotExist,
	}
}

// ReadDown reads the down function
func (c *Code) ReadDown(version uint) (io.ReadCloser, string, error) {
	if m, ok := c.ms.Down(version); ok {
		sql, found := c.migrations[m.Raw]
		if !found {
			return nil, "", fmt.Errorf("migrate-code readdown: %w", ErrNotFound)
		}

		return ioutil.NopCloser(strings.NewReader(sql)), m.Identifier, nil
	}

	return nil, "", &os.PathError{
		Op:   fmt.Sprintf("read version %v", version),
		Path: c.path,
		Err:  os.ErrNotExist,
	}
}

// Close does nothing
func (c *Code) Close() error {
	return nil
}

func init() {
	source.Register("code", &Code{})
}
