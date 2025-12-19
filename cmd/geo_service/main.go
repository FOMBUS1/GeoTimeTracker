package main

import (
	"fmt"
	"os"

	"github.com/FOMBUS1/GeoTimeTracker/config"
)

func main() {
	cfg, err := config.LoadConfig(os.Getenv("configPath"))
	if err != nil {
		panic(fmt.Sprintf("ошибка парсинга конфига, %v", err))
	}

	fmt.Println(cfg)
}
