package services

import (
	"ClubTennis/models"
	"ClubTennis/repositories"
	"sync"

	"gorm.io/gorm"
)

type cachedSnippet struct {
	snippet models.Snippet
	dirty   bool
}

// service for handling rules like challenge match rules and Ladder rules.
// since theres currently a discrete no of them they are cached.
type SnippetService struct {
	repo       repositories.SnippetRepository
	cacheMutex sync.RWMutex             // RW mutex for cache
	cache      map[string]cachedSnippet //map is to snip directly since infrequently updated
}

const HOMEPAGE_CATEGORY string = "homepage"
const LADDER_CATEGORY string = "ladder"
const CHALLENGE_CATEGORY string = "challenge"

func NewSnippetService(db *gorm.DB) *SnippetService {
	var s SnippetService

	s.repo = *repositories.NewSnippetRepository(db)
	s.cache = make(map[string]cachedSnippet)
	snips := s.repo.FindAll()

	for _, snippet := range snips {
		s.cache[snippet.Category] = cachedSnippet{snippet: snippet, dirty: false}
	}

	if err := s.ensureHomepage(); err != nil {
		print(err.Error())
		return nil
	}
	if err := s.ensureChallenge(); err != nil {
		print(err.Error())
		return nil
	}
	if err := s.ensureLadder(); err != nil {
		print(err.Error())
		return nil
	}

	return &s
}

// ensure homepage has something
func (s *SnippetService) ensureHomepage() error {
	snip := s.Get(HOMEPAGE_CATEGORY)
	if snip != nil {
		return nil
	}
	return s.Set(HOMEPAGE_CATEGORY, models.NewSnippet("", "<h2>NC State Club Tennis</h2><p>Custom homepage has not been set up yet!</p>"))
}

// ensure ladder has something
func (s *SnippetService) ensureLadder() error {
	snip := s.Get(LADDER_CATEGORY)
	if snip != nil {
		return nil
	}
	return s.Set(LADDER_CATEGORY, models.NewSnippet("", "<h1>Ladder Rules (TOC Format)</h1><h3>Who can I challenge?</h3><ul><li>Top 10: Can challenge up to 2 spots above</li><li>11-20: Can challenge up to 3 spots above</li><li>21-35: Can challenge up to 5 spots above</li><li>36+: Can challenge up to 8 spots above</li></ul><h3>Challenge match rules</h3><ul><li><a href=\"/club/challengerules\">click here</a></li></ul><h3>Practice Schedule</h3><ul><li>Red team: Monday/Wednesday Practice</li><li>White team: Tuesday/Thursday Practice</li></ul>"))
}

// ensure challenge has something
func (s *SnippetService) ensureChallenge() error {
	snip := s.Get(CHALLENGE_CATEGORY)
	if snip != nil {
		return nil
	}
	return s.Set(CHALLENGE_CATEGORY, models.NewSnippet("", "<h1>Challenge Match Rules</h1><ul><li>1 set, First to 6 games</li><li><strong>No Ads</strong></li><li>Switch courts every 4 games</li><li>Tiebreaker at 5-5</li><li>Tiebreaker is sudden death first to 5 points</li><li>Tiebreak serving is 2, 2, 2, 3. Serving starts on deuce side</li></ul>"))

}

func (s *SnippetService) Get(category string) *models.Snippet {

	s.cacheMutex.RLock()
	val, ok := s.cache[category]
	s.cacheMutex.RUnlock()

	if ok && !val.dirty {
		return &val.snippet
	}

	snip := s.repo.FindByCategory(category)
	if snip == nil {
		return nil
	}

	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()
	s.cache[category] = cachedSnippet{snippet: *snip, dirty: false}
	return snip
}

func (s *SnippetService) Set(category string, snippet *models.Snippet) error {
	s.cacheMutex.Lock()
	s.cache[category] = cachedSnippet{dirty: true}
	s.cacheMutex.Unlock()

	return s.repo.Save(category, snippet)
}
