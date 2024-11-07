package service

import (
	"base-gin/config"
	"base-gin/repository"
)

var (
	accountService   *AccountService
	personService    *PersonService
	publisherService *PublisherService
	authorService    *AuthorService
	bookService      *BookService
	borrowService    *BorrowService
)

func SetupServices(cfg *config.Config) {
	accountService = NewAccountService(cfg, repository.GetAccountRepo())
	personService = NewPersonService(repository.GetPersonRepo())
	publisherService = NewPublisherService(repository.GetPublisherRepo())
	authorService = NewAuthorService(repository.GetAuthorRepo())
	bookService = NewBookService(repository.GetBookRepo())
	borrowService = NewBorrowService(repository.GetBorrowRepo())
}

func GetAccountService() *AccountService {
	if accountService == nil {
		panic("account service is not initialised")
	}

	return accountService
}

func GetPersonService() *PersonService {
	return personService
}

func GetPublisherService() *PublisherService {
	return publisherService
}

func GetAuthorService() *AuthorService {
	return authorService
}

func GetBookService() *BookService {
	return bookService
}

func GetBorrowService() *BorrowService {
	return borrowService
}
