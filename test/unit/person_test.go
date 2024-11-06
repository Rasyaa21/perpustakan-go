package unit_test

import (
	"base-gin/domain"
	"base-gin/domain/dto"
	"base-gin/util"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPerson_Update_Success(t *testing.T) {
	birthDate, _ := time.Parse("2006-01-02", "1993-09-13")
	gender := domain.GenderFemale
	params := dto.PersonUpdateReq{
		ID:           dummyMember.ID,
		Fullname:     util.RandomStringAlpha(4) + " " + util.RandomStringAlpha(6) + " " + util.RandomStringAlpha(6),
		Gender:       string(gender),
		BirthDateStr: birthDate.Format("2006-01-02"),
		BirthDate:    birthDate,
	}

	err := personRepo.Update(&params)
	assert.Nil(t, err)

	item, _ := personRepo.GetByID(dummyMember.ID)
	assert.Equal(t, params.Fullname, item.Fullname)
	assert.EqualValues(t, params.Gender, string(*item.Gender))
	assert.EqualValues(t, params.BirthDateStr, item.BirthDate.Format("2006-01-02"))
}
