# cli-toolkit

Small, independent, reusable Go packages for CLI projects.

## dirs

Platform-aware config and cache directory resolution for a named application.

```go
import "github.com/0xjuanma/cli-toolkit/dirs"

app := dirs.New("golazo")

configDir, err := app.ConfigDir() // creates the directory if needed
cacheDir, err := app.CacheDir()   // creates the directory if needed

// No I/O, cannot fail — safe to call from init() or flag help text.
app.ConfigDirDisplay()                  // "~/.config/golazo"
app.ConfigFileDisplay("debug.log")      // "~/.config/golazo/debug.log"
```

### Platform behavior

| OS      | `ConfigDir()`                                             | `CacheDir()`                    |
|---------|-------------------------------------------------------------|----------------------------------|
| Linux   | `$XDG_CONFIG_HOME/<name>`, else `~/.config/<name>`         | `$XDG_CACHE_HOME/<name>`, else `~/.cache/<name>` |
| macOS   | `~/.<name>`                                                | `~/Library/Caches/<name>`       |
| Windows | `~/.<name>`                                                | `%LocalAppData%/<name>`         |

`CacheDir()` delegates to `os.UserCacheDir()`.
