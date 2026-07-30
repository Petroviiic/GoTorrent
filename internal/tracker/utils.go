package tracker

import (
	"bufio"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

const CUSTOM_TRACKERS_FOLDER = "./internal/tracker/custom_trackers"

func loadTrackers() []string {
	res := []string{}

	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return res
	}
	if err := filepath.WalkDir(filepath.Join(dir, CUSTOM_TRACKERS_FOLDER), func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		trackers, err := readTrackerFile(path)
		if err != nil {
			log.Printf("Error reading file %s: %v", path, err)
			return nil
		}

		res = append(res, trackers...)
		return nil
	}); err != nil {
		fmt.Println(err)
	}
	return res
}

func readTrackerFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var list []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			list = append(list, line)
		}
	}

	return list, scanner.Err()
}
