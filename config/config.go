package config

type PlayConfig struct {
	Port           string `yaml:"port"`
	BinariesFolder string `yaml:"binaries_folder"`
	FilesFolder    string `yaml:"files_folder"`
	HTMLName       string `yaml:"html_name"`
}
