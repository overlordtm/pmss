package scanner_test

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/overlordtm/pmss/internal/apiclient"
	"github.com/overlordtm/pmss/pkg/scanner"
)

func initDirStructure(root string, depth int, maxfile int, maxdir int) error {

	rnd := rand.New(rand.NewSource(42))

	if depth == 0 {
		return nil
	}
	if _, err := os.Stat(root); err == nil {
		// root already exists, exit
		return nil
	}

	os.MkdirAll(root, 0755)

	dirCount := rnd.Intn(maxdir) + 1
	fileCount := rnd.Intn(maxfile) + 1

	for i := 0; i < fileCount; i++ {
		f, err := os.Create(filepath.Join(root, fmt.Sprintf("file%d", i)))
		if err != nil {
			return err
		}
		fileSize := rnd.Intn(1024 * 100)

		_, err = io.Copy(f, io.LimitReader(rnd, int64(fileSize)))

		if err != nil {
			return err
		}
		f.Close()
	}

	for i := 0; i < dirCount; i++ {
		if err := initDirStructure(filepath.Join(root, fmt.Sprintf("dir%d", i)), depth-1, maxfile, maxdir); err != nil {
			return err
		}
	}
	return nil
}

func TestScanFile(t *testing.T) {

	testDir := "./testdata/gen1"

	now := time.Now()

	t.Logf("initializing test directory with random files")

	err := initDirStructure(testDir, 5, 100, 10)
	if err != nil {
		t.Error(err)
	}

	t.Logf("initDirStructure took %v", time.Since(now))

	s := scanner.New()

	results := make(chan apiclient.File, 1024)

	wg := sync.WaitGroup{}
	wg.Add(1)
	cnt := 0
	go func() {
		// pump results and count them
		defer wg.Done()
		for _ = range results {
			cnt++
		}
	}()

	now = time.Now()
	err = s.Scan(results, testDir)
	if err != nil {
		t.Error(err)
	}
	duration := time.Since(now)
	wg.Wait()
	t.Logf("Scan took %v, %v per file, %v files per second, %v file total", duration, time.Duration(int64(duration)/int64(cnt)), int64(cnt)/int64(duration/time.Second), cnt)
}
