package tests

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/thadUra/Go-ML/cluster"
)

/**
 * TestMultithread()
 * Tests singlethread vs multhreading
 */
func TestMultithread(t *testing.T) {
	// Open iris data file (WIP CLOSE FILE)
	file, err := os.Open("./misc/iris_data.csv")
	if err != nil {
		t.Fatalf(`TestCluster(): failed to open file -> "%s"`, err)
	}
	r := csv.NewReader(file)

	// Parse data and store into slices
	var data [][]float64
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf(`TestCluster(): failed in file parsing -> "%s"`, err)
		}
		next := make([]float64, 4)
		for i := 0; i < 4; i++ {
			next[i], err = strconv.ParseFloat(record[i], 64)
			if err != nil {
				t.Fatalf(`TestCluster(): failed in file parsing -> "%s"`, err)
			}
		}
		data = append(data, next)
	}

	// Test PCA
	result, err := cluster.PCA(data, 2)
	if err != nil {
		t.Fatalf(`TestCluster(): failed at building PCA -> "%s"`, err)
	}

	// Test KMeans with a single thread
	num_iterations := 100000
	start := time.Now()
	for i := 0; i < num_iterations; i++ {
		km := cluster.NewKMeans(3, 500)
		_ = km.Train(result)
		_, _, _ = km.Evaluate(result)
	}
	duration := time.Since(start)
	fmt.Printf("single thread: %v\n", duration)

	// Test KMeans with multithreading
	start = time.Now()
	var wg sync.WaitGroup
	groups := 10
	wg.Add(groups)
	for i := 0; i < groups; i++ {
		go func() {
			for i := 0; i < num_iterations/groups; i++ {
				km := cluster.NewKMeans(3, 500)
				_ = km.Train(result)
				_, _, _ = km.Evaluate(result)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	duration = time.Since(start)
	fmt.Printf("multithread: %v\n", duration)
}
