package main

import (
	"backend-trainee-banner-avito/internal/config"
	"fmt"
)

func main() {
	cfg := config.MustLoadConfig()
	fmt.Println(cfg.Env)

	// TODO : LOGGER  : slog
	// TODO : STORAGE : postgresql
	// TODO : ROUTER  : chi, chi-render
	// TODO : SERVER

}
