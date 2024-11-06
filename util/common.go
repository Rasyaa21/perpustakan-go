package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func UUIDFromString(stringUUID string) (uuid.UUID, error) {
	return uuid.Parse(stringUUID)
}

func RandomNumber(length int) string {
	otpChars := "1234567890"
	buffer := make([]byte, length)
	_, _ = rand.Read(buffer)
	charsLength := len(otpChars)
	for i := 0; i < length; i++ {
		buffer[i] = otpChars[int(buffer[i])%charsLength]
	}
	return string(buffer)
}

func RandomString(length int) string {
	chars := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := make([]byte, length)
	_, _ = rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = chars[b%byte(len(chars))]
	}
	return string(bytes)
}

func RandomStringAlpha(length int) string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := make([]byte, length)
	_, _ = rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = chars[b%byte(len(chars))]
	}
	return string(bytes)
}

func ValidatePhoneNumber(text string) bool {
	re := regexp.MustCompile(`^[1-9][0-9]*$`)
	return re.MatchString(text)
}

// TruncateString separate string every four characters with delimiter.
func TruncateString(text string, delimiter string) string {
	for i := 4; i < len(text); i += 5 {
		text = text[:i] + delimiter + text[i:]
	}
	return text
}

// EncryptAESGCM is using 12 byte iv instead of 16.
func EncryptAESGCM(plain string, secret string) (string, error) {
	keyByte := []byte(secret)
	block, err := aes.NewCipher(keyByte)
	if err != nil {
		return "", err
	}

	iv := make([]byte, 12)
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	cipherByte := aesGCM.Seal(nil, iv, []byte(plain), nil)
	tag := cipherByte[len(cipherByte)-16:]
	encrypted := cipherByte[:len(cipherByte)-16]

	encryptedEncoded := base64.StdEncoding.EncodeToString(encrypted)
	ivEncoded := base64.StdEncoding.EncodeToString(iv)
	tagEncoded := base64.StdEncoding.EncodeToString(tag)
	result := fmt.Sprintf("%s$@%s$@%s", encryptedEncoded, ivEncoded, tagEncoded)
	return result, nil
}

func DecryptAESGCM(cipherText string, secret string) (string, error) {
	key := []byte(secret)
	splitCipherText := strings.Split(cipherText, "$@")
	encrypted, _ := base64.StdEncoding.DecodeString(splitCipherText[0])
	iv, _ := base64.StdEncoding.DecodeString(splitCipherText[1])
	tag, _ := base64.StdEncoding.DecodeString(splitCipherText[2])
	encryptedWithTag := encrypted
	encryptedWithTag = append(encryptedWithTag, tag...)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plainText, err := gcm.Open(nil, iv, encryptedWithTag, nil)
	if err != nil {
		return "", err
	}
	return string(plainText), nil
}

// PasswordHash Hash password with bcrypt (default cost, 10 rounds).
func PasswordHash(plain string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

// VerifyPasswordHash Verify password created with PasswordHash.
func VerifyPasswordHash(hash string, plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
	return err == nil
}

// MaskEmailUsername that takes an email address as input and returns a masked
// version of the username part of the email. If the email does not contain an
// "@" symbol, it returns the original email. The username is masked by replacing
// all but the first and last characters with asterisks. The masked username is
// then concatenated with the domain part of the email and returned.
func MaskEmailUsername(email string) string {
	atIndex := strings.Index(email, "@")
	if atIndex == -1 {
		return email
	}

	username := email[:atIndex]
	var maskedUsername string
	if len(username) > 2 {
		maskedUsername = username[:1] + strings.Repeat("*", len(username)-2) + username[len(username)-1:]
	} else {
		maskedUsername = strings.Repeat("*", len(username))
	}

	return maskedUsername + email[atIndex:]
}

func RemoveDash(str string) string {
	return strings.ReplaceAll(str, "-", "")
}

func SanitiseName(str string) string {
	re := regexp.MustCompile(`[:;|~!@#$%^*+={}\[\]\/"]+`)
	return re.ReplaceAllString(str, "")
}
