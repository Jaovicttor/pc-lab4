package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

type FileSum struct {
	totalSum    int
	fingerprint []int
	path        string
}

// read a file from a filepath and return a slice of bytes
func readFile(filePath string) ([]byte, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file %s: %v", filePath, err)
		return nil, err
	}
	return data, nil
}

// sum all bytes of a file
func sum(filePath string) (FileSum, error) {
	var chunks []int

	data, err := readFile(filePath)
	if err != nil {
		return FileSum{
			0,
			chunks,
			filePath,
		}, err
	}

	_sum := 0
	current := 0
	total := 0

	for _, b := range data {
		current += 1

		_sum += int(b)
		total += int(b)

		if current%100 == 0 {
			chunks = append(chunks, _sum)
			_sum = 0
		}
	}

	result := FileSum{
		total,
		chunks,
		filePath,
	}

	return result, err
}

func sumWrapper(filePath string, sumChan chan FileSum) {
	fileSum, _ := sum(filePath)
	sumChan <- fileSum
}

func similiarity(base, target []int) float64 {
	counter := 0
	targetCopy := append([]int(nil), target...)

	for _, value := range base {
		for i, t := range targetCopy {
			if value == t {
				counter++
				targetCopy = append(targetCopy[:i], targetCopy[i+1:]...)
				break
			}
		}
	}

	return float64(counter) / float64(len(base))
}

// print the totalSum for all files and the files with equal sum
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <file1> <file2> ...")
		return
	}

	var totalSum int64
	size := len(os.Args[1:])
	sumChannel := make(chan FileSum, size)
	sums := make(map[int][]string)
	fileFingerprints := make(map[string][]int)

	for _, path := range os.Args[1:] {
		go sumWrapper(path, sumChannel)
	}

	for i := 0; i < size; i++ {
		result := <-sumChannel
		totalSum += int64(result.totalSum)
		sums[result.totalSum] = append(sums[result.totalSum], result.path)
		fileFingerprints[result.path] = result.fingerprint
	}

	fmt.Println(totalSum)

	for sum, files := range sums {
		if len(files) > 1 {
			fmt.Printf("Sum %d: %v\n", sum, files)
		}
	}

	for _, i := range os.Args[1:] {
		for _, j := range os.Args[1:] {
			filesSimilarity := similiarity(fileFingerprints[i], fileFingerprints[j])
			fmt.Printf("Similarity between %s and %s is: %v\n", i, j, (filesSimilarity * 100))
		}
	}
}
