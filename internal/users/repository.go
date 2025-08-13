package users

import (
	"database/sql"
	"strings"
)

func GetUsersFromDB(db *sql.DB, search string, limit, offset int, sortBy, order string) ([]User, int, error) {
	if search != "" {
		search = "%" + strings.ToLower(search) + "%"
	} else {
		search = "%"
	}

	//Total count
	var total int
	err := db.QueryRow(`
		SELECT COUNT(*)
		FROM users
		WHERE deleted_at IS NULL
		AND (LOWER(name) LIKE $1 OR LOWER(email) LIKE $2)
	`, search, search).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	//Data retrieval
	query := `
		SELECT id, name, email
		FROM users
		WHERE deleted_at IS NULL
		AND (LOWER(name) LIKE $1 OR LOWER(email) LIKE $1)
		ORDER BY ` + sortBy + ` ` + order + `
		LIMIT $2 OFFSET $3
	`
	rows, err := db.Query(query, search, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var usersList []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			return nil, 0, err
		}
		usersList = append(usersList, user)
	}
	return usersList, total, nil
}
