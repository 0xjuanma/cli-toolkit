package dirs_test

import (
	"fmt"

	"github.com/0xjuanma/cli-toolkit/dirs"
)

// Example_golazo shows how golazo's internal/data/storage.go would migrate
// off its hardcoded ".golazo" config/cache logic onto this package.
func Example_golazo() {
	app := dirs.New("golazo")

	configDir, err := app.ConfigDir()
	if err != nil {
		fmt.Println("config dir error:", err)
		return
	}

	cacheDir, err := app.CacheDir()
	if err != nil {
		fmt.Println("cache dir error:", err)
		return
	}

	_ = configDir
	_ = cacheDir

	// Safe to call from init() or flag help text — no I/O, cannot fail.
	fmt.Println(app.ConfigFileDisplay("golazo_debug.log"))
}
