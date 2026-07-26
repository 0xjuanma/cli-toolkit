# cli-toolkit

Small, independent, reusable Go packages for CLI projects.

| Package | What it does | Why it's useful |
|---|---|---|
| [`dirs`](#dirs) | Platform-aware config/cache directory resolution for a named app | Avoids reimplementing XDG/macOS/Windows path conventions per project |
| [`cache`](#cache) | Generic, thread-safe in-memory cache with TTL + max-size eviction | Drop-in caching for any key/value type, no hand-rolled eviction logic |
| [`ratelimit`](#ratelimit) | Interval-based blocking rate limiter | Throttle outbound API calls without pulling in a scheduling library |

<details>
<summary><h2 id="dirs">dirs</h2></summary>

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

</details>

<details>
<summary><h2 id="cache">cache</h2></summary>

Generic, thread-safe, in-memory cache with TTL and max-size eviction.

```go
import "github.com/0xjuanma/cli-toolkit/cache"

c := cache.NewMap[string, int](time.Minute, 100) // TTL, max entries

c.Set("key", 42)                        // stored with the default TTL
c.SetWithTTL("key", 42, 10*time.Second) // stored with a custom TTL

val, ok := c.Get("key")
c.Delete("key")
c.Clear()
keys := c.Keys() // non-expired keys only
```

When `maxSize` is reached, expired entries are purged first, then the oldest remaining entry is evicted.

</details>

<details>
<summary><h2 id="ratelimit">ratelimit</h2></summary>

Interval-based rate limiting for outbound requests.

```go
import "github.com/0xjuanma/cli-toolkit/ratelimit"

l := ratelimit.New(200 * time.Millisecond) // minimum interval between requests
// or: l := ratelimit.NewFromRate(60)      // requests per minute

l.Wait() // blocks until the minimum interval has elapsed since the last call
```

</details>
