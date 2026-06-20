package chunk

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/dkr290/go-advanced-projects/go-rag-api/llm"
	"github.com/dkr290/go-advanced-projects/go-rag-api/vector"
	"github.com/fsnotify/fsnotify"
)

// watch this document directory and see if new file is there
// using the package  "github.com/fsnotify/fsnotify"

// documentDelay is the debounce window before processing a newly detected file.
// This gives the writer time to finish writing / flushing to disk.
const documentDelay = 500 * time.Millisecond

// implement function here tyhe directory is ./documents and preocessed files are moved in documents/processed

// Watch watches the source directory for new or modified files.
// When a file appears, it waits for a short delay (debounce), processes it
// through the ingestion pipeline, and moves it to the processed directory.
// It blocks until ctx is cancelled.
func Watch(
	ctx context.Context,
	opts Options,
	embedder llm.Embedder,
	store vector.Store,
	logger *log.Logger,
) error {

	sourceDir, err := filepath.Abs(opts.SourceDir)
	if err != nil {
		return fmt.Errorf("resolve source dir: %w", err)
	}
	processedDir, err := filepath.Abs(opts.ProcessedDir)
	if err != nil {
		return fmt.Errorf("resolve processed dir: %w", err)
	}

	if sourceDir == processedDir {
		return errors.New("source and processed directories must be different")
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("new watcher: %w", err)
	}

	defer watcher.Close()

	// Ensure both directories exist
	if err := os.MkdirAll(opts.SourceDir, 0o755); err != nil {
		return fmt.Errorf("mkdir source: %w", err)
	}
	if err := os.MkdirAll(opts.ProcessedDir, 0o755); err != nil {
		return fmt.Errorf("mkdir processed: %w", err)
	}

	// Start watching the source directory
	if err := watcher.Add(opts.SourceDir); err != nil {
		return fmt.Errorf("add source dir to watcher: %w", err)
	}

	logger.Printf("watching %s for new documents", opts.SourceDir)

// Process any files that were already in sourceDir before we started watching
	if err := batchProcessExisting(ctx, sourceDir, processedDir,opts, embedder, store, logger); err != nil {
		logger.Printf("batch process existing files: %v", err)
	}

	// Buffered channel to queue files for processing.
	// The buffer allows the watcher goroutine to keep up even if
	// the processor goroutine is busy embedding.
	pending := make(chan string, 100)

	// Worker goroutine: reads from the pending queue, debounces,
	// processes, and moves the file to the processed directory.
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case path := <-pending:
				// Debounce: give the writer time to finish
				select {
				case <-ctx.Done():
					return
				case <-time.After(documentDelay):
				}

				// File might have been removed before we got to it
				if _, err := os.Stat(path); err != nil {
					continue
				}

				// Run the full ingestion pipeline (READ → CHUNK → EMBED → DELETE → UPSERT)
				if err := processOne(ctx, path, opts, embedder, store); err != nil {
					logger.Printf("processOne %s: %v", path, err)
					continue
				}

				// Move the processed file to the "processed" directory
				base := filepath.Base(path)
				dest := filepath.Join(opts.ProcessedDir, base)
				if err := os.Rename(path, dest); err != nil {
					logger.Printf("rename %s → %s: %v", path, dest, err)
					continue
				}

				logger.Printf("processed %s → %s", path, dest)
			}
		}
	}()

	// Main event loop
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			// React to Create or Write events for regular files (not directories)
			if event.Op&(fsnotify.Create|fsnotify.Write) != 0 {
				if info, err := os.Stat(event.Name); err == nil && !info.IsDir() {
					select {
					case pending <- event.Name:
					default:
						logger.Printf("pending queue full, skipping %s", event.Name)
					}
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			logger.Printf("watcher error: %v", err)
		}
	}
}

// batchProcessExisting scans the source directory for files that were already
// there before the watcher started, processes them concurrently, and moves
// successful ones to the processed directory.
func batchProcessExisting(
	ctx context.Context,
	sourceDir string,
	processedDir string,
	opts Options,
	embedder llm.Embedder,
	store vector.Store,
	logger *log.Logger,
) error {
	entries, err := os.ReadDir(sourceDir)
	if err != nil {
		return fmt.Errorf("read dir %s: %w", sourceDir, err)
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var errs []error

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		path := filepath.Join(sourceDir, entry.Name())
		if !supportedFormat(path) {
			continue
		}

		wg.Add(1)
		go func(p string) {
			defer wg.Done()

			if err := processOne(ctx, p, opts, embedder, store); err != nil {
				mu.Lock()
				errs = append(errs, fmt.Errorf("%s: %w", p, err))
				mu.Unlock()
				return
			}

			base := filepath.Base(p)
			dest := filepath.Join(processedDir, base)
			if err := os.Rename(p, dest); err != nil {
				mu.Lock()
				errs = append(errs, fmt.Errorf("rename %s: %w", p, err))
				mu.Unlock()
				return
			}

			logger.Printf("batch processed %s → %s", p, dest)
		}(path)
	}

	wg.Wait()

	if len(errs) > 0 {
		return fmt.Errorf("%d files failed to process: %v", len(errs), errs[0])
	}

	return nil
}
