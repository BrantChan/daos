package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/daos-stack/daos/src/control/common/test"
)

var (
	testPoolDirs = [...]string{"a", "ab", "aac", "aaad"}
	testVosFiles = [...]string{"vos-0", "vos-1", "vos-2", "vos-10", "vos-a"}
)

func createFile(t *testing.T, filePath string) {
	t.Helper()

	fd, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("Failed to create test vos file %s: %v", filePath, err)
	}
	fd.Close()
}

func createDirAll(t *testing.T, dirPath string) {
	t.Helper()

	if err := os.MkdirAll(dirPath, 0755); err != nil {
		t.Fatalf("Failed to create test pool directory %s: %v", dirPath, err)
	}
}

func testSetup(t *testing.T) (tmpDir string, teardown func()) {
	t.Helper()

	tmpDir, teardown = test.CreateTestDir(t)

	for _, dir := range testPoolDirs {
		createDirAll(t, filepath.Join(tmpDir, dir))
		for _, file := range testVosFiles {
			createFile(t, filepath.Join(tmpDir, dir, file))
		}
	}

	createDirAll(t, filepath.Join(tmpDir, "foo"))
	createFile(t, filepath.Join(tmpDir, "foo", "bar"))

	createDirAll(t, filepath.Join(tmpDir, "bar"))
	createDirAll(t, filepath.Join(tmpDir, "bar", "foo"))
	createDirAll(t, filepath.Join(tmpDir, "bar", "baz"))
	createFile(t, filepath.Join(tmpDir, "bar", "baz", "no_vos"))

	return
}

func TestListVosFiles(t *testing.T) {
	tmpDir, teardown := testSetup(t)
	t.Cleanup(teardown)

	for name, tc := range map[string]struct {
		args   string
		expRes []string
	}{
		"unaccessible": {
			args:   "/root/",
			expRes: []string{},
		},
		"void director prefix": {
			args: tmpDir + string(os.PathSeparator),
			expRes: []string{
				filepath.Join(tmpDir, "a") + string(os.PathSeparator),
				filepath.Join(tmpDir, "ab") + string(os.PathSeparator),
				filepath.Join(tmpDir, "aac") + string(os.PathSeparator),
				filepath.Join(tmpDir, "aaad") + string(os.PathSeparator),
				filepath.Join(tmpDir, "foo") + string(os.PathSeparator),
				filepath.Join(tmpDir, "bar") + string(os.PathSeparator),
			},
		},
		"a pool directory prefix": {
			args: tmpDir + string(os.PathSeparator) + "a",
			expRes: []string{
				filepath.Join(tmpDir, "a") + string(os.PathSeparator),
				filepath.Join(tmpDir, "ab") + string(os.PathSeparator),
				filepath.Join(tmpDir, "aac") + string(os.PathSeparator),
				filepath.Join(tmpDir, "aaad") + string(os.PathSeparator),
			},
		},
		"aa pool directory prefix": {
			args: tmpDir + string(os.PathSeparator) + "aa",
			expRes: []string{
				filepath.Join(tmpDir, "aac") + string(os.PathSeparator),
				filepath.Join(tmpDir, "aaad") + string(os.PathSeparator),
			},
		},
		"all vos files": {
			args: tmpDir + string(os.PathSeparator) + "a" + string(os.PathSeparator),
			expRes: []string{
				filepath.Join(tmpDir, "a") + string(os.PathSeparator) + "vos-0",
				filepath.Join(tmpDir, "a") + string(os.PathSeparator) + "vos-1",
				filepath.Join(tmpDir, "a") + string(os.PathSeparator) + "vos-2",
				filepath.Join(tmpDir, "a") + string(os.PathSeparator) + "vos-10",
			},
		},
		"vos-1 prefix files": {
			args: tmpDir + string(os.PathSeparator) + "a" + string(os.PathSeparator) + "vos-1",
			expRes: []string{
				filepath.Join(tmpDir, "a") + string(os.PathSeparator) + "vos-1",
				filepath.Join(tmpDir, "a") + string(os.PathSeparator) + "vos-10",
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			results := listVosFiles(tc.args)
			test.AssertStringsEqual(t, tc.expRes, results, "listDirVos results do not match expected")
		})
	}
}

func TestListPoolDirs(t *testing.T) {
	tmpDir, teardown := testSetup(t)
	t.Cleanup(teardown)

	for name, tc := range map[string]struct {
		args   string
		expRes []string
	}{
		"unaccessible": {
			args:   "/root/",
			expRes: []string{},
		},
		"void director prefix": {
			args: tmpDir + string(os.PathSeparator),
			expRes: []string{
				filepath.Join(tmpDir, "a"),
				filepath.Join(tmpDir, "ab"),
				filepath.Join(tmpDir, "aac"),
				filepath.Join(tmpDir, "aaad"),
				filepath.Join(tmpDir, "bar") + string(os.PathSeparator),
			},
		},
		"a pool directory prefix": {
			args: tmpDir + string(os.PathSeparator) + "a",
			expRes: []string{
				filepath.Join(tmpDir, "a"),
				filepath.Join(tmpDir, "ab"),
				filepath.Join(tmpDir, "aac"),
				filepath.Join(tmpDir, "aaad"),
			},
		},
		"aa pool directory prefix": {
			args: tmpDir + string(os.PathSeparator) + "aa",
			expRes: []string{
				filepath.Join(tmpDir, "aac"),
				filepath.Join(tmpDir, "aaad"),
			},
		},
		"bar pool directory": {
			args:   tmpDir + string(os.PathSeparator) + "bar" + string(os.PathSeparator),
			expRes: []string{},
		},
	} {
		t.Run(name, func(t *testing.T) {
			results := listPoolDirs(tc.args)
			test.AssertStringsEqual(t, tc.expRes, results, "listDirPool results do not match expected")
		})
	}
}

func TestFilterSuggestions(t *testing.T) {
	t.Skip("TestFilterSuggestions not implemented yet")
}
