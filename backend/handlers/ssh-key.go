package handlers

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"samey/config"
	"samey/helpers"
	"samey/models"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/bson"
)

func AddUser(c fiber.Ctx) error {
	type Request struct {
		PublicKey string `json:"public_key"`
		Password  string `json:"password"` // Added password field
	}

	var req Request
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Validate required fields
	if req.PublicKey == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Public key is required"})
	}

	username := extractUsernameFromSSHKey(req.PublicKey)
	if username == "" {
		// Note: This approach doesn't work well in a web server context
		// Consider requiring username in the request instead
		fmt.Print("No username detected in SSH key. Please enter a username: ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		username = strings.TrimSpace(input)
	}

	// Check if username is still empty
	if username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Username could not be determined"})
	}

	usersCollection := config.DB.Collection("users")
	var existingUser models.User
	err := usersCollection.FindOne(c.Context(), bson.M{"username": username}).Decode(&existingUser)
	if err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "User already exists"})
	}

	// Pass both username and password to the helper function
	err = helpers.CreateLinuxUserWithPassword(username, req.Password, req.PublicKey)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user on base system"})
	}

	newUser := models.User{
		Username:  username,
		PublicKey: req.PublicKey,
	}
	_, err = usersCollection.InsertOne(c.Context(), newUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save user to database"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User created successfully", "username": username})
}

func extractUsernameFromSSHKey(publicKey string) string {
	keyParts := strings.Split(publicKey, " ")
	if len(keyParts) < 3 {
		return ""
	}
	comment := keyParts[2]
	if strings.Contains(comment, "@") {
		return strings.Split(comment, "@")[0]
	}
	return ""
}
