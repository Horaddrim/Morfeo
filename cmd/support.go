package cmd

// Langs is a stuct that holds out language support
type Langs struct {
	langs map[string]string
}

var instace Langs

func init() {
	instace.langs = make(map[string]string)

	instace.langs["python"] = "--version"
	instace.langs["java"] = "-version"
	instace.langs["go"] = "version"
	instace.langs["ruby"] = "--version"
	instace.langs["rails"] = "ruby --version"
	instace.langs["javac"] = "version"
	instace.langs["node"] = "--version"
	instace.langs["dotnet"] = "--version"
}

func langVersion(lang string) (string, string) {
	version := instace.langs[lang]
	return lang, version
}
