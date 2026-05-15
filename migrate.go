// Package migrate provides database migration functionality.
// It is a fork of golang-migrate/migrate with additional features and fixes.
package migrate

import (
	"errors"
	"fmt"
	"os"
	"sync"
)

// DefaultPrefetchMigrations is the default number of migrations to prefetch.
const DefaultPrefetchMigrations = 10

// ErrNoChange is returned when no migration is needed.
var ErrNoChange = errors.New("no change")

// ErrNilVersion is returned when the version is nil.
var ErrNilVersion = errors.New("no migration")

// ErrLocked is returned when the database is locked.
var ErrLocked = errors.New("database locked")

// ErrLockTimeout is returned when the lock times out.
var ErrLockTimeout = errors.New("lock timeout")

// Migrate is the main struct for managing database migrations.
type Migrate struct {
	// sourceName is the registered source driver name.
	sourceName string
	// sourceDrv is the source driver instance.
	sourceDrv Source

	// databaseName is the registered database driver name.
	databaseName string
	// databaseDrv is the database driver instance.
	databaseDrv Database

	// PrefetchMigrations is the number of migrations to prefetch.
	PrefetchMigrations uint

	// LockTimeout is the timeout for acquiring the database lock.
	LockTimeout int

	// Log is an optional logger.
	Log Logger

	// GracefulStop is a channel to signal graceful stop.
	GracefulStop chan bool
	isGracefulStop bool

	stateMu sync.Mutex
	isRunning bool
}

// Logger is the interface for logging migration events.
type Logger interface {
	Printf(format string, v ...interface{})
	Verbose() bool
}

// New creates a new Migrate instance from source and database URL strings.
func New(sourceURL, databaseURL string) (*Migrate, error) {
	m := &Migrate{
		PrefetchMigrations: DefaultPrefetchMigrations,
		GracefulStop:       make(chan bool, 1),
	}

	sourceDrv, err := Open(sourceURL)
	if err != nil {
		return nil, fmt.Errorf("source: %w", err)
	}
	m.sourceDrv = sourceDrv

	databaseDrv, err := OpenDatabase(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("database: %w", err)
	}
	m.databaseDrv = databaseDrv

	return m, nil
}

// NewWithSourceInstance creates a new Migrate instance with an existing source driver.
func NewWithSourceInstance(sourceName string, sourceInstance Source, databaseURL string) (*Migrate, error) {
	m := &Migrate{
		sourceName:         sourceName,
		sourceDrv:          sourceInstance,
		PrefetchMigrations: DefaultPrefetchMigrations,
		GracefulStop:       make(chan bool, 1),
	}

	databaseDrv, err := OpenDatabase(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("database: %w", err)
	}
	m.databaseDrv = databaseDrv

	return m, nil
}

// Close closes the source and database connections.
func (m *Migrate) Close() (sourceErr error, databaseErr error) {
	ch := make(chan error)

	go func() {
		ch <- m.sourceDrv.Close()
	}()

	go func() {
		ch <- m.databaseDrv.Close()
	}()

	return <-ch, <-ch
}

// logPrintf logs a message if a logger is set.
func (m *Migrate) logPrintf(format string, v ...interface{}) {
	if m.Log != nil {
		m.Log.Printf(format, v...)
	}
}

// logVerbosePrintf logs a verbose message if a logger is set and verbose mode is enabled.
func (m *Migrate) logVerbosePrintf(format string, v ...interface{}) {
	if m.Log != nil && m.Log.Verbose() {
		m.Log.Printf(format, v...)
	}
}

// newLogger creates a simple logger that writes to stderr.
func newLogger() Logger {
	return &defaultLogger{}
}

type defaultLogger struct{}

func (l *defaultLogger) Printf(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, format, v...)
}

func (l *defaultLogger) Verbose() bool {
	return false
}
