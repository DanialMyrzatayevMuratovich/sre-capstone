package scripts

import (
	"cinema-booking/config"
	"cinema-booking/models"
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func CreateAdmin() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	usersCollection := config.GetCollection("users")

	email := "test@gmail.com"

	var existing models.User
	err := usersCollection.FindOne(ctx, bson.M{"email": email}).Decode(&existing)
	if err == nil {
		log.Printf("⚠️  Admin user already exists: %s", email)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("test1234"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("❌ Failed to hash password:", err)
	}

	admin := models.User{
		Email:    email,
		Password: string(hashedPassword),
		FullName: "Test Admin",
		Role:     "admin",
		Wallet: models.Wallet{
			Balance:  0,
			Currency: "KZT",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = usersCollection.InsertOne(ctx, admin)
	if err != nil {
		log.Fatal("❌ Failed to create admin:", err)
	}

	log.Printf("✅ Admin created: %s / test1234 (role: admin)", email)
}
