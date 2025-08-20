package users

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrInvalidSort  = errors.New("invalid sort field")
	ErrInvalidOrder = errors.New("invalid sort order")
)

type ListOptions struct {
	Search string
	Limit  int
	Offset int
	SortBy string
	Order  string
}

func ListUsers(db *sql.DB, opt ListOptions) ([]User, int, error) {
	if opt.Limit <= 0 {
		opt.Limit = 10
	}
	if opt.Offset < 0 {
		opt.Offset = 0
	}
	if opt.SortBy == "" {
		opt.SortBy = "id"
	}
	if opt.Order == "" {
		opt.Order = "ASC"
	}

	allowedSort := map[string]bool{
		"id":         true,
		"name":       true,
		"email":      true,
		"created_at": true,
	}
	if !allowedSort[strings.ToLower(opt.SortBy)] {
		return nil, 0, ErrInvalidSort
	}

	sortCol := strings.ToLower(opt.SortBy)

	order := strings.ToUpper(opt.Order)
	if order != "ASC" && order != "DESC" {
		return nil, 0, ErrInvalidOrder
	}

	//search term
	search := "%"
	if strings.TrimSpace(opt.Search) != "" {
		search = "%" + strings.ToLower(strings.TrimSpace(opt.Search)) + "%"
	}

	var total int
	if err := db.QueryRow(`
		SELECT COUNT(*)
		FROM users
		WHERE deleted_at IS NULL
		AND (LOWER(name) LIKE $1 OR LOWER(email) LIKE $2)
	`, search, search).Scan(&total); err != nil {
		return nil, 0, err
	}

	// Page data
	// ORDER BY must be injected *after* validation (no placeholders allowed for identifiers)
	query := fmt.Sprintf(`
		SELECT id, name, email
		FROM users
		WHERE deleted_at IS NULL
		  AND (LOWER(name) LIKE $1 OR LOWER(email) LIKE $1)
		ORDER BY %s %s
		LIMIT $2 OFFSET $3
	`, sortCol, order)

	rows, err := db.Query(query, search, opt.Limit, opt.Offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var out []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			return nil, 0, err
		}
		out = append(out, u)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return out, total, nil
}

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

func UpdateUserFromDB(db *sql.DB, id int, user *User) error {
	//check if email exists
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND id != $2)", user.Email, id).Scan(&exists)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("email %s already exists", user.Email)
	}
	// Validate ID
	if id <= 0 {
		return fmt.Errorf("invalid user ID: %d", id)
	}

	// Update user
	result, err := db.Exec("UPDATE users SET name = $1, email = $2 WHERE id = $3 AND deleted_at IS NULL", user.Name, user.Email, id)
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
	result, err := db.Exec("UPDATE users SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL", id)
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
