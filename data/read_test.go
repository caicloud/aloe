package data

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/caicloud/aloe/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	caseContent = []byte(`
summary: "test case"
`)

	contextContent = []byte(`
summary: "test context"
`)
)

func constructDir(path string, dir *fakeDir) error {
	if dir == nil {
		return nil
	}
	dirpath := filepath.Join(path, dir.name)
	if err := os.Mkdir(dirpath, 0700); err != nil {
		return err
	}
	for _, f := range dir.files {
		caseFile := filepath.Join(dirpath, f)
		if err := ioutil.WriteFile(caseFile, caseContent, 0666); err != nil {
			return err
		}
	}
	contextFile := filepath.Join(dirpath, types.ContextFile)
	if err := ioutil.WriteFile(contextFile, contextContent, 0666); err != nil {
		return err
	}
	for _, d := range dir.dirs {
		if err := constructDir(dirpath, &d); err != nil {
			return err
		}
	}
	return nil
}

type fakeDir struct {
	name  string
	dirs  []fakeDir
	files []string
}

func TestWalk(t *testing.T) {
	cases := []struct {
		description string
		dir         *fakeDir
		expectedDir *Dir
		hasErr      bool
	}{
		{
			description: "read all contexts and cases",
			dir: &fakeDir{
				name: "test",
				dirs: []fakeDir{
					{
						name: "nested",
						files: []string{
							"eee.yaml",
							"fff.yaml",
						},
					},
				},
				files: []string{
					"aaa.yaml",
					"bbb.yaml",
					"ccc.yaml",
					"ddd.yaml",
				},
			},
			expectedDir: &Dir{
				Context: types.Context{
					Summary: "test context",
				},
				Name:    "test",
				CaseNum: 6,
				Dirs: map[string]Dir{
					"nested": Dir{
						Context: types.Context{
							Summary: "test context",
						},
						Name:    "nested",
						CaseNum: 2,
						Dirs:    map[string]Dir{},
						Files: map[string]File{
							"eee.yaml": {
								Case: types.Case{
									Summary: "test case",
								},
								Name: "eee.yaml",
							},
							"fff.yaml": {
								Case: types.Case{
									Summary: "test case",
								},
								Name: "fff.yaml",
							},
						},
					},
				},
				Files: map[string]File{
					"aaa.yaml": {
						Case: types.Case{
							Summary: "test case",
						},
						Name: "aaa.yaml",
					},
					"bbb.yaml": {
						Case: types.Case{
							Summary: "test case",
						},
						Name: "bbb.yaml",
					},
					"ccc.yaml": {
						Case: types.Case{
							Summary: "test case",
						},
						Name: "ccc.yaml",
					},
					"ddd.yaml": {
						Case: types.Case{
							Summary: "test case",
						},
						Name: "ddd.yaml",
					},
				},
			},
		},
		{
			description: "empty path",
			dir:         nil,
			expectedDir: nil,
			hasErr:      true,
		},
	}

	for _, c := range cases {
		path, err := ioutil.TempDir("", "test")
		require.NoError(t, err, c.description)
		defer os.RemoveAll(path)
		require.NoError(t, constructDir(path, c.dir), c.description)
		dirpath := filepath.Join(path, "empty")
		if c.dir != nil {
			dirpath = filepath.Join(path, c.dir.name)
		}
		dir, err := Walk(dirpath)
		if !c.hasErr {
			assert.NoError(t, err, c.description)
		}
		assert.Equal(t, c.expectedDir, dir, c.description)
	}
}
