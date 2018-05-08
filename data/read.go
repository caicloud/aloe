package data

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/ghodss/yaml"

	"github.com/caicloud/aloe/types"
)

// Dir defines a directory to store test data
type Dir struct {
	// Context defines the context config
	Context types.ContextConfig

	// Name defines the dir name
	Name string

	// CaseNum defines the case num under the dir
	CaseNum int

	// Dirs defines the children dirs
	Dirs map[string]Dir

	// Files defines the children files
	Files map[string]File
}

// File defines a file to store a test case
type File struct {
	// Case defines the case config
	Case types.Case

	// Name defines the file name
	Name string
}

// Walk walks a dir and return Dir struct
func Walk(path string) (*Dir, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	ctxConfig, err := readContext(path)
	if err != nil {
		return nil, fmt.Errorf("read context config %v error: %v", path, err)
	}
	dir := Dir{
		Context: *ctxConfig,
		Name:    filepath.Base(path),
		Dirs:    map[string]Dir{},
		Files:   map[string]File{},
	}
	for _, file := range files {
		name := file.Name()
		childPath := filepath.Join(path, name)
		if file.IsDir() {
			childDir, err := Walk(childPath)
			if err != nil {
				return nil, err
			}
			dir.Dirs[name] = *childDir
			dir.CaseNum += childDir.CaseNum
		} else if !isIgnored(name) {
			c, err := readCase(childPath)
			if err != nil {
				return nil, fmt.Errorf("read test case %v error: %v", childPath, err)
			}
			dir.Files[name] = File{
				Case: *c,
				Name: file.Name(),
			}
			dir.CaseNum++
		}
	}
	return &dir, nil
}

func readContext(dir string) (*types.ContextConfig, error) {
	contextFile := filepath.Join(dir, types.ContextFile)

	contextBody, err := ioutil.ReadFile(contextFile)
	if err != nil {
		return nil, err
	}
	context := types.ContextConfig{}
	if err := yaml.Unmarshal(contextBody, &context); err != nil {
		return nil, fmt.Errorf("can't unmarshal %v, err: %v", contextFile, err)
	}
	if err := ValidateContext(&context); err != nil {
		return nil, err
	}
	return &context, nil
}

func isIgnored(name string) bool {
	if filepath.Base(name) == types.ContextFile {
		return true
	}
	ext := filepath.Ext(name)
	// TODO(liubog2008): add json support
	return ext != ".yaml"
}

func readCase(file string) (*types.Case, error) {
	body, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	c := types.Case{}
	if err := yaml.Unmarshal(body, &c); err != nil {
		return nil, err
	}
	if err := ValidateCase(&c); err != nil {
		return nil, err
	}
	return &c, nil
}
