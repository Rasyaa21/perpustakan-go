package unit_test

import (
	"base-gin/config"
	"base-gin/domain"
	"base-gin/domain/dao"
	"base-gin/repository"
	"base-gin/storage"
	"base-gin/util"
	"fmt"
	"log"
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

	dummyAdmin  *dao.Person
	dummyMember *dao.Person

	accountRepo *repository.AccountRepository
	personRepo  *repository.PersonRepository
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

	a := createDummyAccount()
	dummyAdmin = createDummyProfile(a)
	dummyMember = createDummyProfile(nil)
	createDummyProfile(nil)
}

func teardownDB() {
	_ = db.Migrator().DropTable(
		&dao.Account{},
		&dao.Person{},
	)
}

func setupDB() {
	_ = db.AutoMigrate(
		&dao.Account{},
		&dao.Person{},
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
