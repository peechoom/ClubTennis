package repositories_test

import (
	"ClubTennis/repositories"
	"fmt"
	"testing"
	"time"

	"github.com/mazen160/go-random"
	"github.com/stretchr/testify/suite"
)

type RBTreeTestSuite struct {
	suite.Suite
	db *repositories.TTLRbTree
}

func (suite *RBTreeTestSuite) SetupTest() {
	suite.db = repositories.NewTTLRbTree()
}

func TestRBTreeSuite(t *testing.T) {
	suite.Run(t, new(RBTreeTestSuite))
}

func (suite *RBTreeTestSuite) TestSetDelete() {
	suite.Require().False(suite.db.Del("poop"))

	suite.db.Set("poop", time.Hour)
	suite.Require().True(suite.db.Del("poop"))
	suite.Require().False(suite.db.Del("poop"))

	var randStrings []string
	for i := 0; i < 1000; i++ {
		str, err := random.String(100)
		str = fmt.Sprintf("%s%d", str, i)
		suite.Require().NoError(err)
		randStrings = append(randStrings, str)
		suite.db.Set(str, time.Minute)
	}

	for _, str := range randStrings {
		suite.Require().True(suite.db.Del(str))
	}

	for _, str := range randStrings {
		suite.Require().False(suite.db.Del(str))
	}

}

func (suite *RBTreeTestSuite) TestExpiredEntry() {
	suite.db.Set("poop", time.Millisecond)
	time.Sleep(time.Millisecond * 2)

	suite.Require().False(suite.db.Del("poop"))

	var expiredStrings []string
	for i := 0; i < 300; i++ {
		str, err := random.String(100)
		str = fmt.Sprintf("%s%d", str, i)
		suite.Require().NoError(err)
		expiredStrings = append(expiredStrings, str)
		suite.db.Set(str, time.Millisecond)
	}
	time.Sleep(time.Millisecond * 2)

	var validStrings []string
	for i := 300; i < 600; i++ {
		str, err := random.String(100)
		str = fmt.Sprintf("%s%d", str, i)
		suite.Require().NoError(err)
		validStrings = append(validStrings, str)
		suite.db.Set(str, time.Minute)
	}

	for i := 0; i < 300; i++ {
		exp := expiredStrings[i]
		val := validStrings[i]
		//interleaved
		suite.Require().True(suite.db.Del(val))
		suite.Require().False(suite.db.Del(exp))
		suite.Require().False(suite.db.Del(val))
	}
}

func (suite *RBTreeTestSuite) TestClean() {
	suite.Require().Zero(suite.db.Clean())
	suite.db.Set("poop", time.Nanosecond)
	time.Sleep(time.Nanosecond * 2)
	suite.Require().Equal(1, suite.db.Clean())

	for i := 0; i < 300; i++ {
		str, err := random.String(100)
		str = fmt.Sprintf("%s%d", str, i)
		suite.Require().NoError(err)
		suite.db.Set(str, time.Millisecond)
	}
	time.Sleep(time.Millisecond * 2)

	suite.Require().Equal(300, suite.db.Clean())

	var count int = 0
	suite.db.Set("b", time.Millisecond)
	suite.db.Set("c", 5*time.Millisecond)
	suite.db.Set("a", 10*time.Millisecond)

	time.Sleep(time.Millisecond * 2)
	count += suite.db.Clean()
	time.Sleep(time.Millisecond * 5)
	count += suite.db.Clean()
	time.Sleep(time.Millisecond * 5)
	count += suite.db.Clean()

	suite.Require().Equal(3, count)

}
