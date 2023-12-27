package models

//
//import "database/sql"
//
//func userExists(db *sql.DB, user User) (bool, error) {
//	row := db.QueryRow("SELECT id FROM users WHERE email LIKE %?%", user.Email)
//	var id int64
//
//	err := row.Scan(&id)
//
//	if err == sql.ErrNoRows {
//		return false, nil
//	}
//
//	if err != nil {
//		return false, err
//	}
//
//	return true, nil
//}
