package Helpers

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"real-time-app/Database"
	"real-time-app/Models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func EnableCORS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173") // Frontend origin
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}
}

func SaveFile(file multipart.File, handler *multipart.FileHeader) (string, error) {
	// Create the uploads directory if it doesn't exist
	err := os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		return "", err
	}
	// Get the original file extension
	ext := path.Ext(handler.Filename)

	// Generate a unique filename (UUID or timestamp)
	uniqueFilename := fmt.Sprintf("%s%s", uuid.NewString(), ext)

	// Create the destination file path
	dst, err := os.Create(fmt.Sprintf("./uploads/%s", uniqueFilename))
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Copy the file contents
	_, err = io.Copy(dst, file)
	if err != nil {
		return "", err
	}

	return uniqueFilename, nil // Return the unique filename instead of the original
}

func GetUserByEmail(email string) (Models.User, error) {
	var user Models.User
	query := `SELECT id, email, password, nickname, first_name, last_name, age, gender FROM users WHERE email = ?`
	err := Database.DB.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.Password, &user.Nickname,
		&user.FirstName, &user.LastName, &user.Age, &user.Gender,
	)
	return user, err
}

// Add to your Helpers package
func GetUserByID(userID int) (*Models.User, error) {
	var user Models.User
	err := Database.DB.QueryRow(`
        SELECT id, nickname, email, first_name, last_name, gender, age 
        FROM users WHERE id = ?`, userID).Scan(
		&user.ID, &user.Nickname, &user.Email,
		&user.FirstName, &user.LastName, &user.Gender, &user.Age)

	if err != nil {
		return nil, err
	}
	return &user, nil
}
