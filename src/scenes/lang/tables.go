package lang

var (
	// in this package been structure the word-tables of Language of App
	Version         = "WinterJam Dic 2022 ??? -version"
	Menu    tblMenu = tblMenu{}
)

type tblMenu struct {
	ChooseLang, ChooseResolution, ChooseFullScreen, BtnConfig, WriteName, QuestExit, Star, Config, Extra, Credits string
}
