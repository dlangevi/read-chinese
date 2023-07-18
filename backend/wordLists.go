package backend

import (
	"embed"
	"fmt"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type (
	WordLists interface {
		GetFrequencyList(list string) map[string]int
		ExportUnknownWordRow() UnknownWordRow
		AddList(name string, path string) error
		DeleteList(name string)
		SetPrimaryList(name string)

		GetPrimaryList() string
		GetLists() []string

		// For inifinte scrolling
		SortWordData()
		SortByOccurance()
		SortByFrequency()

		SetWordSourceFromBook(bookId int)
		SetWordSourceFromSearch(search string)
		SetWordSourceFromFavorites()
		SetWordSourceFromAll()
		SetWordSourceFromList(listname string) error
		SetWordSourceFromHsk(ver string, lvl int) error
		SetFrequencySource(listname string)
		SetOccuranceSource(source string)

		RowCount() int
		GetRows(startRow int, endRow int) []UnknownWordRow
		Destroy()
	}

	tableSession struct {
		words           []string
		occuranceMap    map[string]int
		frequencyMap    map[string]int
		sortByFrequency bool
		sortByOccurance bool
	}

	wordLists struct {
		backend       *Backend
		activeSession tableSession
	}
)

func NewWordLists(
	backend *Backend,
) *wordLists {
	lists := &wordLists{
		backend: backend,
	}
	lists.activeSession.sortByFrequency = true
	return lists
}

// For now expect file with new word on each line
func (lists *wordLists) AddList(name string, path string) error {
	savedPath := ConfigDir("userLists", name)
	contents, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = os.WriteFile(savedPath, contents, 0666)
	if err != nil {
		return err
	}
	lists.backend.UserSettings.SaveList(name, savedPath)
	return nil
}

func (lists *wordLists) DeleteList(name string) {
	lists.backend.UserSettings.DeleteList(name)
	// TODO delete old lists
}

func (lists *wordLists) SetPrimaryList(name string) {
	// TODO this could fail
	lists.backend.UserSettings.SetPrimaryList(name)
}

type WordListInfo struct {
	Name      string `json:"name"`
	Path      string `json:"path"`
	IsPrimary bool   `json:"isPrimary"`
}

func (lists *wordLists) ExportWordListInfo() WordListInfo {
	return WordListInfo{}
}

func (lists *wordLists) GetWordListsInfo() map[string]WordListInfo {
	listInfoMap := map[string]WordListInfo{}
	primaryList := lists.backend.UserSettings.WordLists.PrimaryWordList
	for name, list := range lists.backend.UserSettings.WordLists.WordLists {
		listInfoMap[name] = WordListInfo{
			Name:      name,
			Path:      list,
			IsPrimary: name == primaryList,
		}
	}
	return listInfoMap
}

func (lists *wordLists) GetPrimaryList() string {
	return lists.backend.UserSettings.WordLists.PrimaryWordList
}

func (lists *wordLists) GetLists() []string {
	listMap := lists.backend.UserSettings.WordLists.WordLists
	listsSlice := []string{}
	for list := range listMap {
		listsSlice = append(listsSlice, list)
	}
	return listsSlice

}

func (lists *wordLists) GetFrequencyList(list string) map[string]int {
	words := map[string]int{}
	listFile := fmt.Sprintf("%s", list)
	fileBytes, err := os.ReadFile(ConfigDir("userLists", listFile))
	if err != nil {
		return words
	}
	fileString := string(fileBytes)

	wordsList := strings.Split(fileString, "\n")
	for i, word := range wordsList {
		trimmed := strings.TrimSpace(word)
		words[trimmed] = i
	}

	return words
}

type UnknownWordRow struct {
	Word       string `json:"word"`
	Pinyin     string `json:"pinyin"`
	Occurance  int    `json:"occurance"`
	Frequency  int    `json:"frequency"`
	Definition string `json:"definition"`
}

func (lists *wordLists) ExportUnknownWordRow() UnknownWordRow {
	return UnknownWordRow{}
}

// TODO is this the way?
func (lists *wordLists) SortByFrequency() {
	lists.activeSession.sortByOccurance = false
	lists.activeSession.sortByFrequency = true
	lists.SortWordData()
}

func (lists *wordLists) SortByOccurance() {
	lists.activeSession.sortByFrequency = false
	lists.activeSession.sortByOccurance = true
	lists.SortWordData()
}

func (lists *wordLists) SortWordData() {
	fmt.Println("Resorting Data",
		lists.activeSession.sortByFrequency,
		lists.activeSession.sortByOccurance)
	words := &lists.activeSession.words

	if lists.activeSession.sortByFrequency {
		sort.Slice(*words, func(i, j int) bool {
			wordI := (*words)[i]
			wordJ := (*words)[j]
			wordIPriority, ok := lists.activeSession.frequencyMap[wordI]
			if !ok {
				wordIPriority = len(lists.activeSession.frequencyMap)
			}
			wordJPriority, ok := lists.activeSession.frequencyMap[wordJ]
			if !ok {
				wordJPriority = len(lists.activeSession.frequencyMap)
			}
			return wordIPriority < wordJPriority
		})
	} else if lists.activeSession.sortByOccurance {
		sort.Slice(*words, func(i, j int) bool {
			wordI := (*words)[i]
			wordJ := (*words)[j]
			wordIPriority := lists.activeSession.occuranceMap[wordI]
			wordJPriority := lists.activeSession.occuranceMap[wordJ]
			return wordIPriority > wordJPriority
		})
	}

	// Emit Event // Wipe Board
	if lists.backend.ctx != nil {
		runtime.EventsEmit(lists.backend.ctx, "ResetBoard")
	}
}

func (lists *wordLists) SetWordSourceFromBook(bookId int) {
	words := lists.backend.BookLibrary.LearningTargetBook(bookId)
	lists.activeSession.words = words
	lists.SetOccuranceSource(fmt.Sprint(bookId))
	lists.SortWordData()
}
func (lists *wordLists) SetWordSourceFromSearch(search string) {
	words := lists.backend.Dictionaries.GetPossibleWords(search)
	lists.activeSession.words = words
	lists.SortWordData()
}

func (lists *wordLists) SetWordSourceFromFavorites() {
	words := lists.backend.BookLibrary.LearningTargetFavorites()
	lists.activeSession.words = words
	lists.SetOccuranceSource("favorites")
	lists.SortWordData()
}

func (lists *wordLists) SetWordSourceFromAll() {
	words := lists.backend.BookLibrary.LearningTarget()
	lists.activeSession.words = words
	lists.SetOccuranceSource("all")
	lists.SortWordData()
}

func (lists *wordLists) SetWordSourceFromList(listname string) error {
	listFile := fmt.Sprintf("%s", listname)
	fileBytes, err := os.ReadFile(ConfigDir("userLists", listFile))
	if err != nil {
		return err
	}
	fileString := string(fileBytes)

	words := []string{}
	rawWords := strings.Split(fileString, "\n")
	for _, word := range rawWords {
		trimmed := strings.TrimSpace(word)
		if len(trimmed) == 0 {
			continue
		}
		inDictionary := lists.backend.Dictionaries.IsInDictionary(trimmed)
		isKnown := lists.backend.KnownWords.IsKnown(trimmed)
		if !isKnown && inDictionary {
			// Its fine if occurance is just 0
			words = append(words, trimmed)
		}
	}
	lists.activeSession.words = words
	lists.SortWordData()
	return nil
}

//go:embed assets/HSK
var hskWords embed.FS

func (lists *wordLists) SetWordSourceFromHsk(version string, level int) error {
	// ensure string == 2.0 or 3.0
	// ensure level == 1 - 6
	hskPath := path.Join(
		"assets",
		"HSK",
		version,
		fmt.Sprintf(`L%v.txt`, level),
	)

	fileBytes, err := hskWords.ReadFile(hskPath)
	if err != nil {
		return err
	}
	fileString := string(fileBytes)

	words := []string{}
	rawWords := strings.Split(fileString, "\n")
	for _, word := range rawWords {
		trimmed := strings.TrimSpace(word)
		trimmed = strings.Trim(trimmed, "\uFEFF")
		if !lists.backend.KnownWords.IsKnown(trimmed) && len(trimmed) > 0 {
			// Its fine if occurance is just 0
			words = append(words, trimmed)
		}
	}
	lists.activeSession.words = words
	lists.SortWordData()
	return nil
}

func (lists *wordLists) SetFrequencySource(listname string) {
	lists.activeSession.frequencyMap = lists.GetFrequencyList(listname)
}

// This is wack
func (lists *wordLists) SetOccuranceSource(occuranceSource string) {
	var occuranceMap map[string]int
	words := lists.activeSession.words
	if occuranceSource == "all" {
		occuranceMap = lists.backend.KnownWords.GetOccurances(words)
	} else if occuranceSource == "favorites" {
		var err error
		occuranceMap, err = lists.backend.BookLibrary.GetFavoriteFrequencies()
		if err != nil {
			fmt.Println("error loading favorites", err)
		}
	} else {
		var err error
		bookInt, err := strconv.Atoi(occuranceSource)
		if err != nil {
			fmt.Println("error loading book", err)
		}
		occuranceMap, err = lists.backend.BookLibrary.GetBookFrequencies(bookInt)
		if err != nil {
			fmt.Println("error loading book", err)
		}
	}
	lists.activeSession.occuranceMap = occuranceMap
}

func (lists *wordLists) RowCount() int {
	return len(lists.activeSession.words)
}

func (lists *wordLists) GetRows(startRow int, endRow int) []UnknownWordRow {
	rows := []UnknownWordRow{}
	words := lists.activeSession.words
	if startRow > len(words) {
		startRow = len(words)
	}
	if endRow > len(words) {
		endRow = len(words)
	}
	for _, word := range words[startRow:endRow] {
		if lists.backend.KnownWords.IsKnown(word) {
			continue
		}
		pinyin := lists.backend.Dictionaries.getPinyin(word)
		definition := lists.backend.Dictionaries.getDefaultDefinition(word)
		occurance := lists.activeSession.occuranceMap[word]
		frequency, found := lists.activeSession.frequencyMap[word]
		if !found {
			frequency = len(lists.activeSession.frequencyMap) + 1
		}
		rows = append(rows,
			UnknownWordRow{
				Word:       word,
				Pinyin:     pinyin,
				Occurance:  occurance,
				Frequency:  frequency,
				Definition: definition,
			})

	}
	return rows
}

func (lists *wordLists) Destroy() {
	lists.activeSession.words = []string{}
}
