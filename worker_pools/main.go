package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sync"
)

type Worker struct {
	Jobs <-chan Job
	Done <-chan struct{}
	Errs chan<- error
}
type Job struct {
	File    string
	Pattern *regexp.Regexp
	Result  chan<- Result
}

type Result struct {
	File       string
	Line       []byte
	LineNumber int
}

func main() {
	nWorker := 5
	jobs := make(chan Job, nWorker)
	results := make(chan Result, nWorker)
	errCh := make(chan error, nWorker)
	done := make(chan struct{})

	var wg sync.WaitGroup
	for i := 0; i < nWorker; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			w := Worker{
				Jobs: jobs,
				Done: done,
				Errs: errCh,
			}
			w.Run()
		}()
	}

	rex, err := regexp.Compile(os.Args[2])
	if err != nil {
		log.Fatalln("error: ", err)
		return
	}

	go func() {
		err := filepath.WalkDir(os.Args[1], func(path string, f fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !f.IsDir() {
				jobs <- Job{
					File:    path,
					Pattern: rex,
					Result:  results,
				}
			}

			return nil
		})
		if err != nil {
			errCh <- err
		}

		close(done)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	for r := range results {
		fmt.Printf("%s|%d|%s\n", r.File, r.LineNumber, r.Line)
	}

	close(errCh)
	for err := range errCh {
		fmt.Printf("error: %s\n", err)
	}
}

func (w *Worker) Run() {
	for {
		select {
		case j := <-w.Jobs:
			err := w.Work(j)
			if err != nil {
				w.Errs <- err
			}
		case <-w.Done:
			return
		}
	}
}

func (w *Worker) Work(job Job) error {
	file, err := os.Open(job.File)
	if err != nil {
		return err
	}
	defer file.Close()

	n := 1
	scnr := bufio.NewScanner(file)
	for scnr.Scan() {
		res := job.Pattern.Find(scnr.Bytes())
		if len(res) > 0 {
			job.Result <- Result{
				File:       job.File,
				Line:       res,
				LineNumber: n,
			}
		}
		n++
	}

	return nil
}
