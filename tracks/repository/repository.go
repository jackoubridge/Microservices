package repository

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Repository struct {
	DB *sql.DB
}

var repo Repository

func Init() {
	if db, err := sql.Open("sqlite3", "/tmp/enterprise.db"); err == nil {
		repo = Repository{DB: db}
	} else {
		log.Fatal("Database initialisation")
	}
}

func Create() int {
	const sql = "CREATE TABLE IF NOT EXISTS Tracks" +
		"(Id TEXT PRIMARY KEY, Audio TEXT)"
	if _, err := repo.DB.Exec(sql); err == nil {
		return 0
	} else {
		return -1
	}
}

func Clear() int {
	const sql = "DELETE FROM Tracks"
	if _, err := repo.DB.Exec(sql); err == nil {
		return 0
	} else {
		return -1
	}
}

func Read(id string) (Track, int64) {
	const sql = "SELECT * FROM Tracks WHERE Id = ?"
	if stmt, err := repo.DB.Prepare(sql); err == nil {
		defer stmt.Close()
		var t Track
		row := stmt.QueryRow(id)
		if err := row.Scan(&t.Id, &t.Audio); err == nil {
			return t, 1
		} else {
			return Track{}, 0
		}
	}
	return Track{}, -1
}

func Update(t Track) int64 {
	const sql = "UPDATE Tracks SET Audio = ? " + "WHERE Id = ?"
	if stmt, err := repo.DB.Prepare(sql); err == nil {
		defer stmt.Close()
		if res, err := stmt.Exec(t.Audio, t.Id); err == nil {
			if n, err := res.RowsAffected(); err == nil {
				return n
			}
		}
	}
	return -1
}

func Insert(t Track) int64 {
	var sql = "INSERT INTO Tracks(Id, Audio) VALUES (?, ?)"
	if stmt, err := repo.DB.Prepare(sql); err == nil {
		defer stmt.Close()
		if res, err := stmt.Exec(t.Id, t.Audio); err == nil {
			if n, err := res.RowsAffected(); err == nil {
				return n
			}
		}
	}
	return -1
}

func Delete(id string) int64 {
	const sql = "DELETE FROM Tracks WHERE Id = ?"
	if stmt, err := repo.DB.Prepare(sql); err == nil {
		defer stmt.Close()
		if res, err := stmt.Exec(id); err == nil {
			if n, err := res.RowsAffected(); err == nil {
				return n
			}
		}
	}
	return -1
}

func List() ([]string, int64) {
	const sql = "SELECT Id FROM TRACKS"
	var tracklist []string
	if stmt, err := repo.DB.Prepare(sql); err == nil {
		defer stmt.Close()
		if rows, err := stmt.Query(); err == nil {
			for rows.Next() {
				var t Track
				if err := rows.Scan(&t.Id); err != nil {
					return tracklist, -1
				} else {
					tracklist = append(tracklist, t.Id)
				}
			}
			return tracklist, 1
		}
	}
	return tracklist, -1
}
