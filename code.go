package code

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"strings"

	"github.com/golang-migrate/migrate/v4/source"
)

// A Migration returns an array of migrations to run
type Migration func() ([]*source.Migration, error)

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
	migrations *source.Migrations
}

// Open calls the relevant function to get its migrations
// The url passed in should be of the form "code://xyz"
// where xyz is a registered code factory in your own code.
func (c *Code) Open(fullUrl string) (source.Driver, error) {
	u, err := url.Parse(fullUrl)
	if err != nil {
		return nil, err
	}

	migrateFactory := u.Opaque

	f, found := factories[migrateFactory]
	if !found {
		return nil, fmt.Errorf("%w: %s", ErrFactoryNotFound, migrateFactory)
	}

	cms, err := f()
	if err != nil {
		return nil, err
	}

	ms := source.NewMigrations()
	for _, m := range cms {
		if !ms.Append(m) {
			return nil, source.ErrDuplicateMigration{}
		}
	}

	return &Code{
		migrations: ms,
	}, nil
}

// First gets the first migration
func (c *Code) First() (uint, error) {
	v, ok := c.migrations.First()
	if !ok {
		return 0, ErrNotFound
	}
	return v, nil
}

// Prev gets the previous migration
func (c *Code) Prev(version uint) (uint, error) {
	v, ok := c.migrations.Prev(version)
	if !ok {
		return 0, ErrNotFound
	}
	return v, nil
}

// Next gets the next migration
func (c *Code) Next(version uint) (uint, error) {
	v, ok := c.migrations.Next(version)
	if !ok {
		return 0, ErrNotFound
	}
	return v, nil
}

// ReadUp reads the up function
func (c *Code) ReadUp(version uint) (io.ReadCloser, string, error) {
	if m, ok := c.migrations.Up(version); ok {
		return ioutil.NopCloser(strings.NewReader(m.Raw)), m.Identifier, nil
	}
	return nil, "", ErrNotFound
}

// ReadDown reads the down function
func (c *Code) ReadDown(version uint) (io.ReadCloser, string, error) {
	if m, ok := c.migrations.Down(version); ok {
		return ioutil.NopCloser(strings.NewReader(m.Raw)), m.Identifier, nil
	}
	return nil, "", ErrNotFound
}

// Close does nothing
func (c *Code) Close() error {
	return nil
}

func init() {
	source.Register("code", &Code{})
}
