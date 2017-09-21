package match

import (
	"regexp"

	"github.com/tchap/go-patricia/patricia"
)

type Match struct {
	Regex            *regexp.Regexp
	Prefix           string
	PrefixIsComplete bool
	Data             interface{}
}

func (m Match) TriePrefix() patricia.Prefix {
	return patricia.Prefix(m.Prefix)
}

type Matches []Match

func MustNewMatch(regex string, data interface{}) Match {
	match := regexp.MustCompile(regex)
	prefix, prefixIsComplete := match.LiteralPrefix()

	return Match{
		Regex:            match,
		Prefix:           prefix,
		PrefixIsComplete: prefixIsComplete,
		Data:             data,
	}
}

func NewMatch(regex string, data interface{}) (*Match, error) {
	match, err := regexp.Compile(regex)
	if err != nil {
		return nil, err
	}

	prefix, prefixIsComplete := match.LiteralPrefix()

	return &Match{
		Regex:            match,
		Prefix:           prefix,
		PrefixIsComplete: prefixIsComplete,
		Data:             data,
	}, nil
}

type Filter struct {
	Trie *patricia.Trie
}

func NewFilter() Filter {
	return Filter{Trie: patricia.NewTrie()}
}

func (f Filter) Add(m Match) {
	var existing *Matches

	prefix := m.TriePrefix()
	iExisting := f.Trie.Get(prefix)

	if iExisting == nil {
		existing = &Matches{}
		f.Trie.Insert(prefix, existing)
	} else {
		existing = iExisting.(*Matches)
	}
	*existing = append(*existing, m)
}

func (f Filter) Match(data string) Matches {
	prefix := patricia.Prefix(data)
	ret := Matches{}

	f.Trie.VisitPrefixes(prefix, func(prefix patricia.Prefix, item patricia.Item) error {
		matches := item.(*Matches)
		for _, match := range *matches {
			if match.Regex.FindAllStringIndex(data, -1) == nil {
				continue
			}
			ret = append(ret, match)
		}
		return nil
	})

	return ret
}
