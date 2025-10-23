package bench

import (
	"database/sql"
	"fmt"
	_ "github.com/tursodatabase/turso-go"
	"log"
	"os"
	"sync"
	"testing"
)

var concurrentDB *sql.DB

func TestMain(m *testing.M) {
	// Clean up any existing test database
	os.Remove("concurrent_test.db")

	var err error
	concurrentDB, err = sql.Open("turso", "concurrent_test.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	// Configure connection pool for concurrent access
	concurrentDB.SetMaxOpenConns(20)
	concurrentDB.SetMaxIdleConns(10)

	// Create table once
	_, err = concurrentDB.Exec("CREATE TABLE bench_test (id INTEGER PRIMARY KEY, data TEXT, value INTEGER)")
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}

	code := m.Run()
	concurrentDB.Close()
	os.Remove("concurrent_test.db")
	os.Exit(code)
}

func BenchmarkConcurrentWrite1(b *testing.B) {
	benchmarkConcurrentWrite(b, 1)
}

func BenchmarkConcurrentWrite2(b *testing.B) {
	benchmarkConcurrentWrite(b, 2)
}

func BenchmarkConcurrentWrite4(b *testing.B) {
	benchmarkConcurrentWrite(b, 4)
}

func BenchmarkConcurrentWrite8(b *testing.B) {
	benchmarkConcurrentWrite(b, 8)
}

func BenchmarkConcurrentWrite10(b *testing.B) {
	benchmarkConcurrentWrite(b, 10)
}

func benchmarkConcurrentWrite(b *testing.B, numWorkers int) {
	var wg sync.WaitGroup
	opsPerWorker := 50 // Fixed number for fair comparison

	b.ResetTimer()

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < opsPerWorker; j++ {
				_, err := concurrentDB.Exec("INSERT INTO bench_test (data, value) VALUES (?, ?)",
					fmt.Sprintf("w%d_o%d", workerID, j), workerID*1000+j)
				if err != nil {
					b.Errorf("Failed to execute insert: %v", err)
				}
			}
		}(i)
	}

	wg.Wait()
}

func BenchmarkConcurrentRead1(b *testing.B) {
	benchmarkConcurrentRead(b, 1)
}

func BenchmarkConcurrentRead2(b *testing.B) {
	benchmarkConcurrentRead(b, 2)
}

func BenchmarkConcurrentRead4(b *testing.B) {
	benchmarkConcurrentRead(b, 4)
}

func BenchmarkConcurrentRead8(b *testing.B) {
	benchmarkConcurrentRead(b, 8)
}

func BenchmarkConcurrentRead10(b *testing.B) {
	benchmarkConcurrentRead(b, 10)
}

func benchmarkConcurrentRead(b *testing.B, numWorkers int) {
	// Insert test data first
	b.StopTimer()
	for i := 0; i < 1000; i++ {
		concurrentDB.Exec("INSERT INTO bench_test (data, value) VALUES (?, ?)", fmt.Sprintf("data_%d", i), i)
	}
	b.StartTimer()

	var wg sync.WaitGroup
	opsPerWorker := 50 // Fixed number for fair comparison

	b.ResetTimer()

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < opsPerWorker; j++ {
				var data string
				var value int
				id := (workerID*opsPerWorker+j)%1000 + 1
				err := concurrentDB.QueryRow("SELECT data, value FROM bench_test WHERE id = ?", id).Scan(&data, &value)
				if err != nil {
					b.Errorf("Failed to query: %v", err)
				}
			}
		}(i)
	}

	wg.Wait()
}
