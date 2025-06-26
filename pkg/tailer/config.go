package tailer

type Config struct {
	ScrapePath string `yaml:"scrape_path"`
	FileRegex  string `yaml:"file_regex"`
}
