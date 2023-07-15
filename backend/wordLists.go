package backend

import (
	"embed"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
)

type (
	WordLists interface {
		GetUnknownHskWords(version string, level int) ([]string, error)
		GetUnknownListWords(list string) ([]string, error)
		GetFrequencyList(list string) map[string]int
		GetWordData(words []string,
			occuranceSource string, frequencySource string) []UnknownWordRow
		ExportUnknownWordRow() UnknownWordRow
		AddList(name string, path string) error
		DeleteList(name string)
		SetPrimaryList(name string)

		GetPrimaryList() string
		GetLists() []string
	}

	wordLists struct {
		backend *Backend
	}
)

func NewWordLists(
	backend *Backend,
) *wordLists {
	lists := &wordLists{
		backend: backend,
	}
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

//go:embed assets/HSK
var hskWords embed.FS

func (lists *wordLists) GetUnknownHskWords(version string, level int) ([]string, error) {
	// ensure string == 2.0 or 3.0
	// ensure level == 1 - 6
	hskPath := path.Join(
		"assets",
		"HSK",
		version,
		fmt.Sprintf(`L%v.txt`, level),
	)
	rows := []string{}

	fileBytes, err := hskWords.ReadFile(hskPath)
	if err != nil {
		return rows, err
	}
	fileString := string(fileBytes)

	words := strings.Split(fileString, "\n")
	for _, word := range words {
		trimmed := strings.TrimSpace(word)
		trimmed = strings.Trim(trimmed, "\uFEFF")
		if !lists.backend.KnownWords.IsKnown(trimmed) && len(trimmed) > 0 {
			// Its fine if occurance is just 0
			rows = append(rows, trimmed)
		}
	}
	return rows, nil
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

func (lists *wordLists) GetUnknownListWords(list string) ([]string, error) {
	rows := []string{}
	listFile := fmt.Sprintf("%s", list)
	fileBytes, err := os.ReadFile(ConfigDir("userLists", listFile))
	if err != nil {
		return rows, err
	}
	fileString := string(fileBytes)

	words := strings.Split(fileString, "\n")
	for _, word := range words {
		trimmed := strings.TrimSpace(word)
		if !lists.backend.KnownWords.IsKnown(trimmed) && len(trimmed) > 0 {
			// Its fine if occurance is just 0
			rows = append(rows, trimmed)
			if len(rows) > 2000 {
				break
			}
		}
	}
	return rows, nil
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

func (lists *wordLists) GetWordData(
	words []string,
	// occuranceSource can be:
	// "all" : take from all books
	// "favorites" : take from favorite books
	// "{number}" : take from a specific book id
	occuranceSource string,
	// which frequency list do you want to use
	// Value is either a specific list, or "combined"
	frequencySource string,
) []UnknownWordRow {
	rows := []UnknownWordRow{}

	var occuranceMap map[string]int
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

	var frequencyMap map[string]int
	frequencyMap = lists.GetFrequencyList(frequencySource)

	for _, word := range words {
		pinyin := lists.backend.Dictionaries.getPinyin(word)
		definition := lists.backend.Dictionaries.getDefaultDefinition(word)
		occurance := occuranceMap[word]
		frequency, found := frequencyMap[word]
		if !found {
			frequency = len(frequencyMap) + 1
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
