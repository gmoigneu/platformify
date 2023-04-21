package platformifiers

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path"
	"text/template"

	"github.com/Masterminds/sprig/v3"

	"github.com/platformsh/platformify/internal/models"
)

const laravelTemplatesPath = "templates/laravel"

type LaravelPlatformifier struct {
	*UserInput
}

func (p *LaravelPlatformifier) Platformify(ctx context.Context) error {
	if p.Stack != models.Laravel.String() {
		return fmt.Errorf("cannot platformify non-laravel stack: %s", p.Stack)
	}

	// Get working directory.
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not get current working directory: %w", err)
	}
	err = fs.WalkDir(templatesFs, laravelTemplatesPath, func(filePath string, d fs.DirEntry, walkErr error) error {
		if d.IsDir() {
			return nil
		}
		tpl, parseErr := template.New(d.Name()).Funcs(sprig.FuncMap()).ParseFS(templatesFs, filePath)
		if parseErr != nil {
			return fmt.Errorf("could not parse template: %w", parseErr)
		}

		filePath = path.Join(cwd, filePath[len(laravelTemplatesPath):])
		if writeErr := writeTemplate(ctx, filePath, tpl, p.UserInput); writeErr != nil {
			return fmt.Errorf("could not write template: %w", writeErr)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}