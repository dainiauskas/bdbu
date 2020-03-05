package app

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

type diskSize struct {
	start   time.Time
	since   time.Duration
	size    int64
	written int64
}

func (ds *diskSize) writeFile() error {
	fName := `/tmp/diskio` // test file
	defer os.Remove(fName)
	f, err := os.Create(fName)
	if err != nil {
		return err
	}
	const defaultBufSize = 4096
	buf := make([]byte, defaultBufSize)
	buf[len(buf)-1] = '\n'
	w := bufio.NewWriterSize(f, len(buf))

	ds.start = time.Now()
	written := int64(0)
	for i := int64(0); i < ds.size; i += int64(len(buf)) {
		nn, err := w.Write(buf)
		written += int64(nn)
		if err != nil {
			return err
		}
	}
	err = w.Flush()
	if err != nil {
		return err
	}
	err = f.Sync()
	if err != nil {
		return err
	}
	ds.since = time.Since(ds.start)

	err = f.Close()
	if err != nil {
		return err
	}

	ds.written = written

	return nil
}

func (ds *diskSize) Print() {
	fmt.Printf("written [%d]: %dB %dns %.2fGB %.2fs %.2fMB/s\n",
		ds.size/(1024*1024*1024),
		ds.written, ds.since,
		float64(ds.written)/1000000000,
		float64(ds.since)/float64(time.Second),
		(float64(ds.written)/1000000)/(float64(ds.since)/float64(time.Second)),
	)
}

type benchDisk struct {
	test1  *diskSize
	test8  *diskSize
	test16 *diskSize
	test32 *diskSize
}

func newBenchDisk() *benchDisk {
	return &benchDisk{
		test1:  &diskSize{size: 1 * (1024 * 1024 * 1024)},
		test8:  &diskSize{size: 8 * (1024 * 1024 * 1024)},
		test16: &diskSize{size: 16 * (1024 * 1024 * 1024)},
		test32: &diskSize{size: 32 * (1024 * 1024 * 1024)},
	}
}

func (bd *benchDisk) start() {
	// Testing 1
	err := bd.test1.writeFile()
	if err == nil {
		bd.test1.Print()
	}

	// Testing 8
	err = bd.test8.writeFile()
	if err == nil {
		bd.test8.Print()
	}

	// Testing 16
	err = bd.test16.writeFile()
	if err == nil {
		bd.test16.Print()
	}
}

// BenchDiskWrite start disk benchmark
func BenchDiskWrite() {
	bd := newBenchDisk()
	bd.start()
}
