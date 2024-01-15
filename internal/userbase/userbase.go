package userbase

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

    _ "github.com/mattn/go-sqlite3"
	sql "github.com/jmoiron/sqlx"
)

const filename = "users.db"

const (
    initilize = `CREATE TABLE IF NOT EXISTS users (
    user_id INTEGER PRIMARY KEY,
    total_requests INTEGER NOT NULL DEFAULT 1,
    date_joined DATETIME NOT NULL,
    last_visited DATETIME NOT NULL
);
`

    updateUser = `INSERT INTO users (user_id, date_joined, last_visited)
    VALUES (?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
    ON CONFLICT(user_id) DO UPDATE SET
        total_requests = total_requests + 1,
        last_visited = CURRENT_TIMESTAMP;`

    countUsers = `SELECT COUNT(*) FROM users;`
)

func getDatabaseFilepath() (string, error) {
    homeDir, err := os.UserHomeDir()

    if err != nil {
        return "", fmt.Errorf("GetDatabaseFilepath: %w\n", err)
    }

    return filepath.Join(homeDir, filename), nil
}

func createDatabaseFile() error {

    fullpath, err := getDatabaseFilepath()

    if err != nil {
        return fmt.Errorf("CreateDatabaseFile: %w\n", err)
    }


    _, err = os.Create(fullpath)

    if err != nil {
        return fmt.Errorf("CreateDatabaseFile: %w\n", err)
    }

    return nil
}

func databaseFileExists() (bool, error){
    var responce bool

    fullpath, err := getDatabaseFilepath()

    if err != nil  {
        return responce, fmt.Errorf("DatabaseFileNotExists: %w\n", err)
    }

    _, err = os.Stat(fullpath)

    if err != nil {
        return responce, fmt.Errorf("DatabaseFileNotExists: %w\n", err)
    }

    return true, nil
}

func connect() (*sql.DB, error){

    fullpath, err := getDatabaseFilepath()

    if err != nil {
        return nil, fmt.Errorf("Coonect: %w\n", err)
    }

    if _, err := databaseFileExists(); err != nil && errors.Is(err, fs.ErrNotExist) {

        if err := createDatabaseFile(); err != nil {
            return nil,  fmt.Errorf("Connect: %w\n", err)
        }
    }

    db, err := sql.Open("sqlite3", fullpath)

    if err != nil  {
        return nil, fmt.Errorf("Connect: %w\n", err)
    }

    _, err = db.Exec(initilize)

    if err != nil  {
        return nil, err
    }

    return db, nil
}

func UpdateUser(id int64) error {

    db, err := connect()

    if err != nil {
        return fmt.Errorf("UpdateUser failed: %w", err)
    }

    _, err = db.Exec(updateUser, id)

    if err != nil {
        return fmt.Errorf("UpdateUser failed: %w", err)
    }

    return nil }

func UsersCount() (int, error) {
    var count int
    db, err := connect()

    if err != nil {
        return count, fmt.Errorf("UsersCount failed: %w", err)
    }

    err = db.Get(&count, countUsers)

    if err != nil {
        return count, fmt.Errorf("UsersCount failed: %w", err)
    }

    return count, nil
}
