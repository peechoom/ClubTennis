package services_test

import (
	"ClubTennis/initializers"
	"ClubTennis/models"
	"ClubTennis/services"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

type EmailServiceTestSuite struct {
	suite.Suite
	s     *services.EmailService
	us    *services.UserService
	userA *models.User
	userB *models.User
	userC *models.User
	userD *models.User
	userE *models.User
}

func (suite *EmailServiceTestSuite) SetupTest() {
	db := initializers.GetTestDatabase()
	if db == nil {
		panic("error in setup!")
	}
	db.Exec("DROP SCHEMA " + initializers.TestDBName + ";")

	db = initializers.GetTestDatabase()
	if db == nil {
		panic("error in setup!")
	}
	err := db.AutoMigrate(models.User{}, models.Match{})
	if err != nil {
		panic(err)
	}

	suite.s = services.NewEmailService("/home/alec/go/src/ClubTennis/templates", "test@test.com", "")

	suite.us = services.NewUserService(db)

	suite.userA, _ = models.NewUser("shboil4", "ncsu", "Sam", "Boiland", "shboil4@ncsu.edu", models.MENS_LADDER)
	suite.userB, _ = models.NewUser("jbeno5", "ncsu", "James", "Benolli", "jbeno5@ncsu.edu", models.MENS_LADDER)
	suite.userC, _ = models.NewUser("pdiddy4", "ncsu", "Puff", "Daddy", "pdiddy@ncsu.edu", models.MENS_LADDER)
	suite.userD, _ = models.NewUser("jobitch2", "ncsu", "Joel", "Embitch", "jobitch@ncsu.edu", models.MENS_LADDER)
	suite.userE, _ = models.NewUser("myprince2", "ncsu", "Lebron", "James", "myprince2@ncsu.edu", models.MENS_LADDER)

	suite.userA.Rank = 1
	suite.userB.Rank = 2
	suite.userC.Rank = 3
	suite.userD.Rank = 4
	suite.userE.Rank = 5

	suite.userA.Matches = make([]*models.Match, 0)
	suite.userB.Matches = make([]*models.Match, 0)
	suite.userC.Matches = make([]*models.Match, 0)
	suite.userD.Matches = make([]*models.Match, 0)
	suite.userE.Matches = make([]*models.Match, 0)

	suite.us.Save(suite.userA)
	suite.us.Save(suite.userB)
	suite.us.Save(suite.userC)
	suite.us.Save(suite.userD)
	suite.us.Save(suite.userE)

}

func TestEmailServiceSuite(t *testing.T) {
	suite.Run(t, new(EmailServiceTestSuite))
}

func (suite *EmailServiceTestSuite) TestChallengeEmailHTML() {
	e1, e2 := suite.s.MakeChallengeEmails(suite.userB, suite.userA)

	//assert beginning bits are there
	suite.Require().Contains(string(e2.HTML), "<title>Challenge Email</title>")
	//assert middle bits are there and replaced
	suite.Require().Contains(string(e2.HTML), "<h3>James Benolli <span style=\"color: #CC0000;\">(2)</span> has issued a challenge for your spot at <span style=\"color: #CC0000;\">(1)</span>.</h3>")
	//assert end bit still there
	suite.Require().Contains(string(e2.HTML), "If the match is not played within 7 days it is automatically considered a forefit")

	suite.Require().Contains(string(e1.HTML), "<title>Challenge Email</title>")
	suite.Require().Contains(string(e1.HTML), "You have challenged Sam Boiland for their rank of <span style=\"color: #CC0000;\">(1)</span>")
	suite.Require().Contains(string(e1.HTML), "If the match is not played within 7 days it is automatically considered a forefit")
}

func (suite *EmailServiceTestSuite) TestChallengeEmailBody() {
	e1, e2 := suite.s.MakeChallengeEmails(suite.userB, suite.userA)

	suite.Require().Equal(e1.To[0], suite.userB.Email)
	suite.Require().Equal(e2.To[0], suite.userA.Email)

	suite.Require().Equal(e1.From, "NC State Club Tennis <test@test.com>")
	suite.Require().Equal(e2.From, "NC State Club Tennis <test@test.com>")

	suite.Require().Contains(e1.Cc, "test@test.com")
	suite.Require().Contains(e2.Cc, "test@test.com")

	suite.Require().Equal(e2.Text, []byte(fmt.Sprintf("You have been challenged by %s %s (%s). Reply to this email to contact them for scheduling.", suite.userB.FirstName, suite.userB.LastName, suite.userB.Email)))
	suite.Require().Equal(e1.Text, []byte(fmt.Sprintf("You successfully challenged %s %s (%s). You should expect an email from them soon regarding scheduling.", suite.userA.FirstName, suite.userA.LastName, suite.userA.Email)))

}
