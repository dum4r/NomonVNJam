package util

import (
	//"encoding/json"
	"encoding/json"
	"fmt"
	"regexp"
)

var (
	currentLenguaje *string
	langs           []string
	tables          map[string]any
)

// Look on folder the Languages and verify struct, the table don't load
func initLenguaje(pointerLang *string, vtables map[string]any) error {
	tables = vtables
	currentLenguaje = pointerLang
	if err := ReadDir(folderLang, func(isDir bool, nameLang string) bool {
		if isDir {
			if err := ReadDir(folderLang+slash+nameLang, func(isDir2 bool, nameTable string) bool {
				if isDir2 {
					return isDir2 // return true because is directory
				}
				if _, bln := tables[nameTable[:len(nameTable)-5]]; bln {
					if match, err := regexp.MatchString(formatLang, nameTable); match && err == nil {
						return false
					}
				}
				return isDir2
			}); err != nil {
				print("Error resources: The tables of '" + nameLang + "' not avalible => ")
				panic(err)
			}
			langs = append(langs, nameLang)
		}
		return false
	}); err != nil {
		return err
	}
	return PushLanguage()
}

type Intables interface{}

// update the tables of current language
func PushLanguage() error {
	for _, lenguage := range langs {
		if lenguage == (*currentLenguaje) { // Current Lang
			if err := ReadDir(folderLang+slash+lenguage, func(i bool, nameTable string) bool {
				dat, err := assets.ReadFile(folderLang + slash + lenguage + slash + nameTable)
				if err != nil {
					return true
				}
				if err := json.Unmarshal([]byte(dat), tables[nameTable[:len(nameTable)-5]]); err != nil {
					print("Error novel: corrupted file")
					panic(err)
				}
				return false
			}); err != nil {
				print("Error resources: The table can't load of '" + lenguage + "'=> ")
				panic(err)
			}
			return nil
		}
	}
	return fmt.Errorf("Error resources: the language was not identified.")
}

func Languages() []string { return langs }
