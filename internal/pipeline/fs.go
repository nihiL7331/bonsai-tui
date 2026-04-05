package pipeline

import (
	"io"
	"os"
	"path/filepath"
)

func copyFile(src string, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

// returns 'true' if 'dstPath' exists
// and is newer than 'srcPath'
func isUpToDate(srcPath string, dstPath string) bool {
	srcStat, err := os.Stat(srcPath)
	if err != nil {
		return false
	}

	dstStat, err := os.Stat(dstPath)
	if err != nil {
		return false
	}

	return dstStat.ModTime().After(srcStat.ModTime())
}

// checks if any file inside 'srcPath'
// was modified more recently than 'dstFile'
func isDirNewer(srcDir string, dstFile string) bool {
	dstStat, err := os.Stat(dstFile)
	if err != nil {
		return true
	}
	dstTime := dstStat.ModTime()

	needsRepack := false
	filepath.WalkDir(srcDir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		stat, err := d.Info()
		if err == nil && stat.ModTime().After(dstTime) {
			needsRepack = true
			return filepath.SkipDir
		}
		return nil
	})

	return needsRepack
}
