package app

import (
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"

	"master-control/model"
)

var searchDirs = func() []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return []string{"/Applications"}
	}
	return []string{"/Applications", filepath.Join(home, "Applications")}
}()

func listApps() []model.App {
	seen := make(map[string]model.App)
	for _, dir := range searchDirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() && strings.HasSuffix(entry.Name(), ".app") {
				name := strings.TrimSuffix(entry.Name(), ".app")
				path := filepath.Join(dir, entry.Name())
				if _, ok := seen[name]; !ok {
					seen[name] = model.App{Name: name, Path: path}
				}
			}
		}
	}

	apps := make([]model.App, 0, len(seen))
	for _, a := range seen {
		apps = append(apps, a)
	}

	sort.Slice(apps, func(i, j int) bool {
		return strings.ToLower(apps[i].Name) < strings.ToLower(apps[j].Name)
	})

	return apps
}

func getAll(c *gin.Context) {
	c.JSON(http.StatusOK, listApps())
}

func findAppPath(name string) (string, bool) {
	for _, dir := range searchDirs {
		path := filepath.Join(dir, name+".app")
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			return path, true
		}
	}

	return "", false
}

func findIcon(appPath string) (string, error) {
	matches, err := filepath.Glob(filepath.Join(appPath, "Contents", "Resources", "*.icns"))
	if err != nil {
		return "", err
	}
	if len(matches) == 0 {
		return "", os.ErrNotExist
	}

	return matches[0], nil
}

func isValidAppPath(path string) bool {
	abs, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	if !strings.HasSuffix(abs, ".app") {
		return false
	}

	for _, dir := range searchDirs {
		dirAbs, err := filepath.Abs(dir)
		if err != nil {
			continue
		}

		if strings.HasPrefix(abs, dirAbs+string(filepath.Separator)) {
			info, err := os.Stat(abs)
			return err == nil && info.IsDir()
		}
	}

	return false
}

func serveIcon(c *gin.Context, appPath string) {
	iconPath, err := findIcon(appPath)
	if err != nil {
		c.JSON(http.StatusNotFound, model.Error{Error: "icon not found"})
		return
	}

	tmp, err := os.CreateTemp("", "app-icon-*.png")
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Error{Error: err.Error()})
		return
	}
	defer func() { _ = os.Remove(tmp.Name()) }()
	_ = tmp.Close()

	cmd := exec.Command("sips", "-s", "format", "png", iconPath, "--out", tmp.Name())
	if out, err := cmd.CombinedOutput(); err != nil {
		c.JSON(http.StatusInternalServerError, model.Error{Error: string(out)})
		return
	}

	c.File(tmp.Name())
}

func getIconByPath(c *gin.Context) {
	path := c.Query("path")
	if path == "" || !isValidAppPath(path) {
		c.JSON(http.StatusNotFound, model.Error{Error: "invalid app path"})
		return
	}

	serveIcon(c, path)
}

func getIcon(c *gin.Context) {
	path, ok := findAppPath(c.Param("name"))
	if !ok {
		c.JSON(http.StatusNotFound, model.Error{Error: "app not found"})
		return
	}

	serveIcon(c, path)
}