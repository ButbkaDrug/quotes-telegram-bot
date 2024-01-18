package userbase

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	sql "github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
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

    updateLastVisited = `INSERT INTO users (user_id, date_joined, last_visited)
    VALUES (?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
    ON CONFLICT(user_id) DO UPDATE SET
        last_visited = CURRENT_TIMESTAMP;`

    updateTotalRequests = `INSERT INTO users (user_id, date_joined, last_visited)
    VALUES (?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
    ON CONFLICT(user_id) DO UPDATE SET
    total_requests = total_requests + 1;`

    countUsers = `SELECT COUNT(*) FROM users;`

    getUser = `SELECT * FROM users WHERE user_id = ?`
)

type User struct {
    User_id         int64
    Total_requests  int
    Date_joined     time.Time
    Last_visited    time.Time
}

type Title struct {
    Title string
    Message string
    Request int
}

var Titles = []Title{
    {
        Title: "Энтузиаст",
        Message: "",
        Request: 0,
    },
    {
        Title: "Искатель мудрости",
        Message: "Отличная работа! Теперь ты 'Искатель мудрости'. Продолжай углублять свои познания!",
        Request: 25,

    },
    {
        Title: "Философ",
        Message: "Поздравляю с достижением уровня 'Философ'. Твои мысли становятся глубже!",
        Request: 50,

    },
    {
        Title: "Знаток цитат",
        Message: "Браво! Теперь ты 'Знаток цитат'. Твои познания в мире цитат впечатляют!",
        Request: 100,

    },
    {
        Title: "Словесник",
        Message: "Вы словесный мастер! Теперь вы 'Словесник'. Продолжай изучать литературные шедевры!",
        Request: 200,

    },
    {
        Title: "Мастер цитат",
        Message: "Мастерство цитат на твоей стороне! Добро пожаловать на уровень 'Мастер цитат'.",
        Request: 400,

    },
    {
        Title: "Глубокий мыcлитель",
        Message: "Поздравляю с достижением уровня 'Глубокий мыслитель'. Твои мысли становятся все более глубокими!",
        Request: 800,

    },
    {
        Title: "Мудрец",
        Message: "Ты стал 'Мудрецом'. Поздравляю с этим выдающимся достижением!",
        Request: 1600,

    },
    {
        Title: "Эпичный литератор",
        Message: "Поздравляем с достижением уровня 'Эпичный литератор'. Ваши слова приобретают новый уровень изысканности!",
        Request: 3200,

    },
    {
        Title: "Великий Архивариус",
        Message: "Поздравляем с достижением вершины мудрости! Теперь вы 'Великий Архивариус'. Ваши слова — светоч для всех!",
        Request: 6400,

    },
}

func(u *User) Graduate() Title {

    for _, t := range Titles {
        if u.Total_requests == t.Request {
            return t
        }
    }

    return Title{}
}

func(u User) CurrentTitle() Title{

    var title Title

    for _, t := range Titles {

        if t.Request < u.Total_requests {
            title = t
        }
    }

    return title
}


func(u User) NextTitle() Title {
    var result Title

    for _, t := range Titles {
        if t.Request > u.Total_requests {
            result = t
            break
        }
    }

    return result
}

func(u User) TimeFromLastVisit() time.Duration {
    return time.Since(u.Last_visited)
}

func(u User) TimeJoined() time.Duration {
    return time.Since(u.Date_joined)
}

func(u User) TimeJoinedString() string {
    var result, unit  string
    duration := time.Since(u.Date_joined)
// Calculate years, days, hours, and minutes
	years := int(duration.Hours() / 24 / 365)
	days := int(duration.Hours() / 24) % 365
	hours := int(duration.Hours()) % 24
	minutes := int(duration.Minutes()) % 60

    if years > 0 {
        switch years % 10 {
        case 1:
            unit = "год"
        case 2, 3, 4:
            unit = "года"
        default:
            unit = "лет"
        }

        if years > 10 && years < 15 { unit = "лет" }
        result += fmt.Sprintf("%d %s", years, unit)
    }
    if days > 0 {
        switch days % 10 {
        case 1:
            unit = "день"
        case 2, 3, 4:
            unit = "дня"
        default:
            unit = "дней"
        }
        if days > 10 && days < 15 { unit = "дней" }
        result += fmt.Sprintf("%d %s ", days, unit)
    }
    if hours > 0 {
        switch hours % 10 {
        case 1:
            unit = "час"
        case 2, 3, 4:
            unit = "часа"
        default:
            unit = "часов"
        }
        if hours > 10 && hours < 15 { unit = "часов" }
        result += fmt.Sprintf("%d %s ", hours, unit)
    }
    if minutes > 0 {
        switch minutes % 10 {
        case 1:
            unit = "минуту"
        case 2, 3, 4:
            unit = "минуты"
        default:
            unit = "минут"
        }
        if minutes > 10 && minutes < 15 { unit = "минут" }
        result += fmt.Sprintf("%d %s", minutes, unit)
    }

    return result

}

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

func(u *User) UpdateLastVisited(id int64) error {

    db, err := connect()

    if err != nil {
        return fmt.Errorf("UpdateUser failed: %w", err)
    }

    _, err = db.Exec(updateLastVisited, id)

    if err != nil {
        return fmt.Errorf("UpdateUser failed: %w", err)
    }

    return nil
}

func(u *User) UpdateTotalRequests(id int64) error {

    db, err := connect()

    if err != nil {
        return fmt.Errorf("UpdateUser failed: %w", err)
    }

    _, err = db.Exec(updateTotalRequests, id)

    if err != nil {
        return fmt.Errorf("UpdateUser failed: %w", err)
    }

    u.Total_requests++

    return nil
}

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

func GetUser(id int64) (*User, error) {
    var u = &User{}
    db, err := connect()

    if err != nil {
        return u, err
    }

    err = db.Get(u, getUser, id)

    if err != nil {
        return u, err
    }

    return u, nil
}
