package storage

import (
	"fmt"
)

type StorageClient struct {
	quorum *Quorum
}

// Creates a new StorageClient with the provided srvAddresses as the configuration
func NewStorageClient(srvAddresses []string) *StorageClient {
	return &StorageClient{
		quorum: newQuorum(srvAddresses),
	}
}

// Writes the provided value to a random server
func (sc *StorageClient) WriteValue(value string) error {
	success := sc.quorum.Write(value)
	if !success {
		return fmt.Errorf("write failed")
	}
	return nil
}

// Returns a slice of values stored on all servers
func (sc *StorageClient) ReadValue() (string, error) {
	reply := sc.quorum.Read()
	if reply == "ERROR" {
		return reply, fmt.Errorf("read failed")
	}
	return reply, nil
}
