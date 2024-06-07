package repositories_test

import (
	"ClubTennis/initializers"
	"ClubTennis/models"
	"ClubTennis/repositories"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

type AnnouncementTestSuite struct {
	suite.Suite
	repo *repositories.AnnouncementRepository
}

func (suite *AnnouncementTestSuite) SetupTest() {
	//drop the schema
	db := initializers.GetTestDatabase()
	if db == nil {
		panic("error in setup!")
	}
	db.Exec("DROP SCHEMA " + initializers.TestDBName + ";")

	//get the schema back. Burning my flash transistors.
	db = initializers.GetTestDatabase()

	err := db.AutoMigrate(&Match{}, &User{}, &models.Announcement{})
	if err != nil {
		panic(err)
	}

	suite.repo = repositories.NewAnnouncementRepository(db)
}

func TestAnnouncementRepo(t *testing.T) {
	suite.Run(t, new(AnnouncementTestSuite))
}

func (suite *AnnouncementTestSuite) TestSubmitGetAnnouncements() {
	post := string("<h1>Test</h1><p>this is a test blah blah blah whatever </p>")
	ann := models.NewAnnouncement(post, "")
	suite.Require().NoError(suite.repo.SubmitAnnouncement(ann))

	f, e := suite.repo.GetAnnouncementPage(0, 5)

	suite.Assert().NoError(e)
	suite.Require().Len(f, 1)

	suite.Require().Equal(post, f[0].Data)
}

func (suite *AnnouncementTestSuite) TestPages() {
	for i := 11; i >= 0; i-- {
		post := string(fmt.Sprintf("<h1>Test</h1><p>this is page %d </p>", i))
		ann := models.NewAnnouncement(post, "")
		suite.Require().NoError(suite.repo.SubmitAnnouncement(ann))
	}

	f, e := suite.repo.GetAnnouncementPage(0, 5)

	suite.Assert().NoError(e)
	suite.Require().Len(f, 5)

	suite.Require().Contains(f[0].Data, "0")
	suite.Require().Contains(f[1].Data, "1")
	suite.Require().Contains(f[2].Data, "2")
	suite.Require().Contains(f[3].Data, "3")
	suite.Require().Contains(f[4].Data, "4")

	f, e = suite.repo.GetAnnouncementPage(1, 5)

	suite.Assert().NoError(e)
	suite.Require().Len(f, 5)

	suite.Require().Contains(f[0].Data, "5")
	suite.Require().Contains(f[1].Data, "6")
	suite.Require().Contains(f[2].Data, "7")
	suite.Require().Contains(f[3].Data, "8")
	suite.Require().Contains(f[4].Data, "9")

	f, e = suite.repo.GetAnnouncementPage(2, 5)

	suite.Assert().NoError(e)
	suite.Require().Len(f, 2)

}
