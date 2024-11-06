package integration_test

import (
	"base-gin/config"
	"base-gin/domain"
	"base-gin/domain/dao"
	"base-gin/repository"
	"base-gin/rest"
	"base-gin/server"
	"base-gin/service"
	"base-gin/storage"
	"base-gin/util"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

const (
	password = "Paswd123"
)

var (
	cfg config.Config
	db  *gorm.DB
	app *gin.Engine

	dummyAdmin  *dao.Person
	dummyMember *dao.Person

	accountRepo   *repository.AccountRepository
	personRepo    *repository.PersonRepository
	publisherRepo *repository.PublisherRepository
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	setup()

	os.Exit(m.Run())
}

func setup() {
	if err := godotenv.Load("./../../.env.test"); err != nil {
		log.Fatal(fmt.Errorf("Test.Integration: Can not find .env.test on root dir"))
	}

	cfg = config.NewConfig()

	storage.InitDB(cfg)
	db = storage.GetDB()
	teardownDB()
	setupDB()

	repository.SetupRepositories()
	accountRepo = repository.GetAccountRepo()
	personRepo = repository.GetPersonRepo()
	publisherRepo = repository.GetPublisherRepo()

	a := createDummyAccount()
	dummyAdmin = createDummyProfile(a)
	dummyMember = createDummyProfile(nil)
	createDummyProfile(nil)

	service.SetupServices(&cfg)

	app = server.Init(&cfg, accountRepo)
	rest.SetupRestHandlers(app)
}

func teardownDB() {
	_ = db.Migrator().DropTable(
		&dao.Account{},
		&dao.Person{},
		&dao.Publisher{},
	)
}

func setupDB() {
	_ = db.AutoMigrate(
		&dao.Account{},
		&dao.Person{},
		&dao.Publisher{},
	)
}

func createDummyAccount() *dao.Account {
	account, _ := dao.NewUser("admin", password, cfg.AuthN.PasswordEncryptionSecret)
	accountRepo.Create(&account)
	return &account
}

func createDummyProfile(account *dao.Account) *dao.Person {
	birthDate, _ := time.Parse("2006-01-02", "1995-04-05")
	male := domain.GenderMale
	person := dao.Person{
		Fullname:  util.RandomStringAlpha(5) + " " + util.RandomStringAlpha(6),
		Gender:    &male,
		BirthDate: &birthDate,
	}
	if account != nil {
		person.AccountID = &account.ID
		person.Account = account
	}

	personRepo.Create(&person)

	return &person
}

func createAuthAccessToken(username string) string {
	token, err := util.CreateAuthAccessToken(cfg, username)
	if err != nil {
		log.Fatal(fmt.Errorf("main_test.createAuthAccessToken %w", err))
	}
	return token
}

func doTest(
	method, url string,
	body interface{},
	authAccessToken string,
) *httptest.ResponseRecorder {
	requestBody, _ := json.Marshal(body)
	r, _ := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	if authAccessToken != "" {
		r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", authAccessToken))
	}
	w := httptest.NewRecorder()
	app.ServeHTTP(w, r)
	if w.Code >= 400 {
		fmt.Printf("[REQUEST] %s \n", string(requestBody)) //nolint:forbidigo //debug
		fmt.Printf("[RESPONSE] %s \n", w.Body.String())    //nolint:forbidigo //debug
	}
	return w
}
