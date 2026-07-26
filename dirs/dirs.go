// Package dirs resolves platform-aware config and cache directory
// locations for a named application.
package dirs

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

// App resolves config and cache directory locations for a named application.
type App struct {
	Name string
}

// New returns an App for the given application name.
func New(name string) App {
	return App{Name: name}
}

// ConfigDir returns the app's config directory, creating it if needed.
//   - Linux: $XDG_CONFIG_HOME/<name>, or ~/.config/<name>
//   - macOS/Windows: ~/.<name>
func (a App) ConfigDir() (string, error) {
	path, err := a.resolveConfigDir(runtime.GOOS)
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(path, 0755); err != nil {
		return "", fmt.Errorf("create config directory: %w", err)
	}

	return path, nil
}

func (a App) resolveConfigDir(goos string) (string, error) {
	if goos == "linux" {
		if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
			return filepath.Join(xdgConfig, a.Name), nil
		}

		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("get home directory: %w", err)
		}
		return filepath.Join(homeDir, ".config", a.Name), nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home directory: %w", err)
	}
	return filepath.Join(homeDir, "."+a.Name), nil
}

// CacheDir returns the app's cache directory, creating it if needed.
// Uses os.UserCacheDir(), which resolves to:
//   - Linux: ~/.cache/<name> (or $XDG_CACHE_HOME/<name>)
//   - macOS: ~/Library/Caches/<name>
//   - Windows: %LocalAppData%/<name>
func (a App) CacheDir() (string, error) {
	userCache, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("get user cache directory: %w", err)
	}

	path := filepath.Join(userCache, a.Name)
	if err := os.MkdirAll(path, 0755); err != nil {
		return "", fmt.Errorf("create cache directory: %w", err)
	}

	return path, nil
}

// ConfigDirDisplay returns a display-friendly config directory path with
// the home directory shown as "~" rather than resolved. Performs no I/O
// and cannot fail — safe to call from init() or flag usage/help text.
func (a App) ConfigDirDisplay() string {
	return a.configDirDisplay(runtime.GOOS)
}

func (a App) configDirDisplay(goos string) string {
	if goos == "linux" {
		if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
			return filepath.Join(xdg, a.Name)
		}
		return filepath.Join("~", ".config", a.Name)
	}

	return filepath.Join("~", "."+a.Name)
}

// ConfigFileDisplay returns a display-friendly path to a file within the
// config directory, e.g. a.ConfigFileDisplay("debug.log").
func (a App) ConfigFileDisplay(filename string) string {
	return filepath.Join(a.ConfigDirDisplay(), filename)
}
