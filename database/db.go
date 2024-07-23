package database

import (
	"database/sql"
	"github.com/DavidAChris/s3lite/services"
	_ "github.com/glebarez/go-sqlite"
	"sync"
)

var (
	s3db *sql.DB
	mu   sync.RWMutex
)

func InitDb(s3Scv *services.S3Serivce, s3Details *services.S3Details) error {
	if s3db != nil {
		return nil
	}
	mu.Lock()
	defer mu.Unlock()
	s3Scv.DownloadFile(s3Details.BucketName, s3Details.Key, s3Details.FileName)
	newDb, err := sql.Open("sqlite", s3Details.FileName)
	if err != nil {
		return err
	}
	s3db = newDb
	return nil
}

func Write(query string, s3Scv *services.S3Serivce, s3Details *services.S3Details) error {
	mu.Lock()
	defer mu.Unlock()
	_, err := s3db.Exec(query)
	if err != nil {
		return err
	}
	s3Scv.UploadFile(s3Details.BucketName, s3Details.Key, s3Details.FileName)
	return err
}

func Query(query string) (*sql.Rows, error) {
	mu.RLock()
	defer mu.RUnlock()

	return s3db.Query(query)
}

func QueryRowScan(query string, dest ...any) error {
	mu.RLock()
	defer mu.Unlock()
	row := s3db.QueryRow(query)

	return row.Scan(dest)
}
