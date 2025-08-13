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

func UpdateUserFromDB(db *sql.DB, id int, name, email string) error {
	// Validate ID
	if id <= 0 {
		return sql.ErrNoRows
	}

	// Update user
	result, err := db.Exec("UPDATE users SET name = $1, email = $2 WHERE id = $3", name, email, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func DeleteUserFromDB(db *sql.DB, id int) error {
	// Validate ID
	if id <= 0 {
		return sql.ErrNoRows
	}

	// Soft delete user
	result, err := db.Exec("UPDATE users SET deleted_at = NOW() WHERE id = $1", id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func GetUserByIDFromDB(db *sql.DB, id int) (User, error) {
	// Validate ID
	if id <= 0 {
		return User{}, sql.ErrNoRows
	}

	var user User
	err := db.QueryRow(`SELECT id, name, email FROM users 
		WHERE id = $1 AND deleted_at IS NULL`, id).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, sql.ErrNoRows
		}
		return User{}, err
	}
	return user, nil
}

func CreateUserInDB(db *sql.DB, user *User) error {
	if user.Name == "" || user.Email == "" {
		return sql.ErrNoRows
	}

	err := db.QueryRow("INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id",
		user.Name, user.Email).Scan(&user.ID)
	if err != nil {
		return err
	}
	return nil
}
