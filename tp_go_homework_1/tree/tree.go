package main

import (
	"fmt"
	"io"
	"os"
	"sort"
)

func countFiles(files *[]os.FileInfo, showFiles bool) int {
	count := 0
	for _, value := range *files {
		if value.IsDir() || showFiles {
			count++
		}
	}

	return count
}

func prepareToSort(files *[]os.FileInfo, filenames *[]string, filemap *map[string]os.FileInfo, showFiles bool) {
	i := 0 // для заполнения слайса с именами
	for _, value := range *files {
		if showFiles || value.IsDir() {
			(*filemap)[value.Name()] = value
			(*filenames)[i] = value.Name()
			i++
		}
	}

	return
}

func choosePrefix(key, count int) (current, subdir string) {
	if key != count-1 {
		current = "├───"
		subdir = "│\t"
	} else {
		current = "└───"
		subdir = "\t"
	}

	return
}

func nextLevel(out io.Writer, filename string, showFiles bool, prefix string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}

	files, err := f.Readdir(0)
	if err != nil {
		return err
	}

	count := countFiles(&files, showFiles) // считаем сколько будет элементов на данном уровне

	filenames := make([]string, count)
	filemap := make(map[string]os.FileInfo, count)

	prepareToSort(&files, &filenames, &filemap, showFiles) // разбирает слайс FileInfo на слайс с именами и мапу с ключами по именам
	sort.Strings(filenames)

	for key, value := range filenames {
		// currentPrefix - будет добавлен к этому файлу
		// subdirPrefix - будет добавлен к вложенным файлам (если этот - директория)
		currentPrefix, subdirPrefix := choosePrefix(key, count)

		fmt.Fprintf(out, prefix+currentPrefix+"%s", filemap[value].Name()) // имя

		if !filemap[value].IsDir() {
			if filemap[value].Size() == 0 {
				fmt.Fprintf(out, " (empty)")
			} else {
				fmt.Fprintf(out, " (%db)", filemap[value].Size())
			}
		}

		fmt.Fprintln(out)

		if filemap[value].IsDir() { // Смотрим вложенные элементы
			nextDirName := filename + string(os.PathSeparator) + filemap[value].Name()
			if err := nextLevel(out, nextDirName, showFiles, prefix+subdirPrefix); err != nil {
				return err
			}
		}
	}

	return nil
}

func dirTree(out io.Writer, filename string, showFiles bool) error {
	return nextLevel(out, filename, showFiles, "")
}
