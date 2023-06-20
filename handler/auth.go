package handler

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

const SALT_SIZE = 16
var SecretKey string = os.Getenv("SECRET_KEY") 

// Generate 16 bytes randomly and securely using the
// Cryptographically secure pseudorandom number generator (CSPRNG)
// in the crypto.rand package
func generateRandomSalt(saltSize int) []byte {
	var salt = make([]byte, saltSize)

	_, err := rand.Read(salt[:])

	if err != nil {
		panic(err)
	}

	return salt
}

// Combine password and salt then hash them using the SHA-512
// hashing algorithm and then return the hashed password
// as a hex string
func hashPassword(password string, salt []byte) string {
	// Convert password string to byte slice
	var passwordBytes = []byte(password)

	// Create sha-512 hasher
	var sha512Hasher = sha512.New()

	// Append salt to password
	passwordBytes = append(passwordBytes, salt...)

	// Write password bytes to the hasher
	sha512Hasher.Write(passwordBytes)

	// Get the SHA-512 hashed password
	var hashedPasswordBytes = sha512Hasher.Sum(nil)

	// Convert the hashed password to a hex string
	var hashedPasswordHex = hex.EncodeToString(hashedPasswordBytes)

	return hashedPasswordHex
}

// Check if two passwords match
func doPasswordsMatch(hashedPassword, currPassword string,
	salt []byte) bool {
	var currPasswordHash = hashPassword(currPassword, salt)

	return hashedPassword == currPasswordHash
}

func loadPrivateKey(path string) *rsa.PrivateKey {
	keyBytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyBytes)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	return key
}

func loadPublicKey(path string) *rsa.PublicKey {
	keyBytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to load public key: %v", err)
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(keyBytes)
	if err != nil {
		log.Fatalf("Failed to parse public key: %v", err)
	}

	return key
}

func generateToken(id string) (string, error) {
	privateKey := []byte(SecretKey)
	// Define the claims for the JWT token
	claims := jwt.StandardClaims{
		Subject:   id,    // Subject (sub) claim
		ExpiresAt: time.Now().Unix() + 3600, // Expiration time (exp) claim
		IssuedAt:  time.Now().Unix(),         // Issued at (iat) claim
	}

	// Create the JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	// Sign the token using the private key
	return token.SignedString(privateKey)
}

// verify token then return id if success
func verifyToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			  return nil, fmt.Errorf("there's an error with the signing method")
		}
		return []byte(SecretKey), nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid token")
	}

	return claims["sub"].(string), nil
}

