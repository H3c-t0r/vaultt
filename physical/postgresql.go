package physical

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	log "github.com/mgutz/logxi/v1"

	"github.com/armon/go-metrics"
	"github.com/lib/pq"
)

// PostgreSQL Backend is a physical backend that stores data
// within a PostgreSQL database.
type PostgreSQLBackend struct {
	table        string
	client       *sql.DB
	put_query    string
	get_query    string
	delete_query string
	list_query   string
	logger       log.Logger
}

// newPostgreSQLBackend constructs a PostgreSQL backend using the given
// API client, server address, credentials, and database.
func newPostgreSQLBackend(conf map[string]string, logger log.Logger) (Backend, error) {
	// Get the PostgreSQL credentials to perform read/write operations.
	connURL, ok := conf["connection_url"]
	if !ok || connURL == "" {
		return nil, fmt.Errorf("missing connection_url")
	}

	unquoted_table, ok := conf["table"]
	if !ok {
		unquoted_table = "vault_kv_store"
	}
	quoted_table := pq.QuoteIdentifier(unquoted_table)

	// Create PostgreSQL handle for the database.
	db, err := sql.Open("postgres", connURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %v", err)
	}

	// Determine if we should use an upsert function (versions < 9.5)
	var upsert_required bool
	upsert_required_query := "SELECT string_to_array(setting, '.')::int[] < '{9,5}' FROM pg_settings WHERE name = 'server_version'"
	if err := db.QueryRow(upsert_required_query).Scan(&upsert_required); err != nil {
		return nil, fmt.Errorf("failed to check for native upsert: %v", err)
	}

	// Setup our put strategy based on the presence or absence of a native
	// upsert.
	var put_query string
	if upsert_required {
		put_query = "SELECT vault_kv_put($1, $2, $3, $4)"
	} else {
		put_query = "INSERT INTO " + quoted_table + " VALUES($1, $2, $3, $4)" +
			" ON CONFLICT (path, key) DO " +
			" UPDATE SET (parent_path, path, key, value) = ($1, $2, $3, $4)"
	}

	// Setup the backend.
	m := &PostgreSQLBackend{
		table:        quoted_table,
		client:       db,
		put_query:    put_query,
		get_query:    "SELECT value FROM " + quoted_table + " WHERE path = $1 AND key = $2",
		delete_query: "DELETE FROM " + quoted_table + " WHERE path = $1 AND key = $2",
		list_query: "SELECT key FROM " + quoted_table + " WHERE path = $1" +
			"UNION SELECT DISTINCT substring(substr(path, length($1)+1) from '^.*?/') FROM " +
			quoted_table + " WHERE parent_path LIKE concat($1, '%')",
		logger: logger,
	}

	return m, nil
}

// splitKey is a helper to split a full path key into individual
// parts: parentPath, path, key
func (m *PostgreSQLBackend) splitKey(fullPath string) (string, string, string) {
	var parentPath string
	var path string

	pieces := strings.Split(fullPath, "/")
	depth := len(pieces)
	key := pieces[depth-1]

	if depth == 1 {
		parentPath = ""
		path = "/"
	} else if depth == 2 {
		parentPath = "/"
		path = "/" + pieces[0] + "/"
	} else {
		parentPath = "/" + strings.Join(pieces[:depth-2], "/") + "/"
		path = "/" + strings.Join(pieces[:depth-1], "/") + "/"
	}

	return parentPath, path, key
}

// Put is used to insert or update an entry.
func (m *PostgreSQLBackend) Put(entry *Entry) error {
	defer metrics.MeasureSince([]string{"postgres", "put"}, time.Now())

	parentPath, path, key := m.splitKey(entry.Key)

	_, err := m.client.Exec(m.put_query, parentPath, path, key, entry.Value)
	if err != nil {
		return err
	}
	return nil
}

// Get is used to fetch and entry.
func (m *PostgreSQLBackend) Get(fullPath string) (*Entry, error) {
	defer metrics.MeasureSince([]string{"postgres", "get"}, time.Now())

	_, path, key := m.splitKey(fullPath)

	var result []byte
	err := m.client.QueryRow(m.get_query, path, key).Scan(&result)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	ent := &Entry{
		Key:   key,
		Value: result,
	}
	return ent, nil
}

// Delete is used to permanently delete an entry
func (m *PostgreSQLBackend) Delete(fullPath string) error {
	defer metrics.MeasureSince([]string{"postgres", "delete"}, time.Now())

	_, path, key := m.splitKey(fullPath)

	_, err := m.client.Exec(m.delete_query, path, key)
	if err != nil {
		return err
	}
	return nil
}

// List is used to list all the keys under a given
// prefix, up to the next prefix.
func (m *PostgreSQLBackend) List(prefix string) ([]string, error) {
	defer metrics.MeasureSince([]string{"postgres", "list"}, time.Now())

	rows, err := m.client.Query(m.list_query, "/"+prefix)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []string
	for rows.Next() {
		var key string
		err = rows.Scan(&key)
		if err != nil {
			return nil, fmt.Errorf("failed to scan rows: %v", err)
		}

		keys = append(keys, key)
	}

	return keys, nil
}
