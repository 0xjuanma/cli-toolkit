package dirs

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveConfigDir_LinuxWithXDGConfigHome(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "/scratch/xdg")

	got, err := New("myapp").resolveConfigDir("linux")
	if err != nil {
		t.Fatalf("resolveConfigDir() error = %v", err)
	}

	want := filepath.Join("/scratch/xdg", "myapp")
	if got != want {
		t.Errorf("resolveConfigDir() = %q, want %q", got, want)
	}
}

func TestResolveConfigDir_LinuxWithoutXDGConfigHome(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "")
	home := t.TempDir()
	t.Setenv("HOME", home)

	got, err := New("myapp").resolveConfigDir("linux")
	if err != nil {
		t.Fatalf("resolveConfigDir() error = %v", err)
	}

	want := filepath.Join(home, ".config", "myapp")
	if got != want {
		t.Errorf("resolveConfigDir() = %q, want %q", got, want)
	}
}

func TestResolveConfigDir_NonLinux(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	for _, goos := range []string{"darwin", "windows"} {
		t.Run(goos, func(t *testing.T) {
			got, err := New("myapp").resolveConfigDir(goos)
			if err != nil {
				t.Fatalf("resolveConfigDir() error = %v", err)
			}

			want := filepath.Join(home, ".myapp")
			if got != want {
				t.Errorf("resolveConfigDir() = %q, want %q", got, want)
			}
		})
	}
}

func TestResolveConfigDir_HomeDirError(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "")
	t.Setenv("HOME", "")

	_, err := New("myapp").resolveConfigDir("linux")
	if err == nil {
		t.Fatal("resolveConfigDir() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "get home directory") {
		t.Errorf("resolveConfigDir() error = %v, want it to wrap \"get home directory\"", err)
	}
}

func TestConfigDirDisplay_LinuxWithXDGConfigHome(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "/scratch/xdg")

	got := New("myapp").configDirDisplay("linux")

	want := filepath.Join("/scratch/xdg", "myapp")
	if got != want {
		t.Errorf("configDirDisplay() = %q, want %q", got, want)
	}
}

func TestConfigDirDisplay_LinuxWithoutXDGConfigHome(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "")

	got := New("myapp").configDirDisplay("linux")

	want := filepath.Join("~", ".config", "myapp")
	if got != want {
		t.Errorf("configDirDisplay() = %q, want %q", got, want)
	}
}

func TestConfigDirDisplay_NonLinux(t *testing.T) {
	for _, goos := range []string{"darwin", "windows"} {
		t.Run(goos, func(t *testing.T) {
			got := New("myapp").configDirDisplay(goos)

			want := filepath.Join("~", ".myapp")
			if got != want {
				t.Errorf("configDirDisplay() = %q, want %q", got, want)
			}
		})
	}
}

func TestConfigFileDisplay(t *testing.T) {
	a := New("myapp")

	got := a.ConfigFileDisplay("debug.log")

	want := filepath.Join(a.ConfigDirDisplay(), "debug.log")
	if got != want {
		t.Errorf("ConfigFileDisplay() = %q, want %q", got, want)
	}
}

func TestConfigDir_CreatesDirectory(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CONFIG_HOME", "")

	got, err := New("myapp").ConfigDir()
	if err != nil {
		t.Fatalf("ConfigDir() error = %v", err)
	}

	if !strings.Contains(got, "myapp") {
		t.Errorf("ConfigDir() = %q, want path containing \"myapp\"", got)
	}
	if info, err := os.Stat(got); err != nil || !info.IsDir() {
		t.Errorf("ConfigDir() did not create directory at %q", got)
	}
}

func TestCacheDir_CreatesDirectory(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	got, err := New("myapp").CacheDir()
	if err != nil {
		t.Fatalf("CacheDir() error = %v", err)
	}

	if filepath.Base(got) != "myapp" {
		t.Errorf("CacheDir() = %q, want path ending in \"myapp\"", got)
	}
	if info, err := os.Stat(got); err != nil || !info.IsDir() {
		t.Errorf("CacheDir() did not create directory at %q", got)
	}
}
