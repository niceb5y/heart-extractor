package db

import (
	"github.com/dghubble/oauth1"
	"github.com/niceb5y/heart-extractor/pkg/crypto"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"path/filepath"
)

type Token struct {
	gorm.Model
	Token       string
	TokenSecret string
}

type DownloadedItem struct {
	gorm.Model
	FileName string
}

var db *gorm.DB

func InitDatabase() {
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalln(err)
	}

	appConfigPath := filepath.Join(configDir, "Heart Extractor")
	if _, err := os.Stat(appConfigPath); os.IsNotExist(err) {
		err = os.Mkdir(appConfigPath, 0755)
		if err != nil {
			log.Fatalln(err)
		}
	}

	dbPath := filepath.Join(appConfigPath, "hx.db")

	db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&Token{}, &DownloadedItem{})
	if err != nil {
		panic("failed to migrate database")
	}
}

func SetAuthToken(accessToken *oauth1.Token) {
	var err error
	var token, tokenSecret string

	token, err = crypto.Encrypt(accessToken.Token)
	if err != nil {
		log.Fatalln(err)
	}

	tokenSecret, err = crypto.Encrypt(accessToken.TokenSecret)
	if err != nil {
		log.Fatalln(err)
	}

	db.Create(&Token{Token: token, TokenSecret: tokenSecret})
}

func GetAuthToken() *oauth1.Token {
	var accessToken Token
	result := db.First(&accessToken)
	if result.Error != nil {
		return nil
	}

	var err error
	var token, tokenSecret string

	token, err = crypto.Decrypt(accessToken.Token)
	if err != nil {
		log.Fatalln(err)
	}

	tokenSecret, err = crypto.Decrypt(accessToken.TokenSecret)
	if err != nil {
		log.Fatalln(err)
	}

	return oauth1.NewToken(token, tokenSecret)
}

func AddDownloadHistory(filename string) {
	db.Create(&DownloadedItem{FileName: filename})
}

func IsExistInDownloadHistory(filename string) bool {
	downloadedItem := DownloadedItem{}
	result := db.Where(&DownloadedItem{FileName: filename}).First(&downloadedItem)
	if result.Error != nil {
		return false
	}

	return true
}

func ResetDownloadHistory() {
	db.Where("1 = 1").Delete(&DownloadedItem{})
}

func ResetAccessToken() {
	db.Where("1 = 1").Delete(&Token{})
}
