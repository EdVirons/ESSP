package mocks

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// MockMinIOClient is a mock implementation of MinIO/S3 client for testing.
type MockMinIOClient struct {
	mu      sync.RWMutex
	objects map[string][]byte // bucket/key -> content
	buckets map[string]bool
	calls   map[string]int
	errors  map[string]error // Method name -> error to return
}

// NewMockMinIOClient creates a new mock MinIO client.
func NewMockMinIOClient() *MockMinIOClient {
	return &MockMinIOClient{
		objects: make(map[string][]byte),
		buckets: make(map[string]bool),
		calls:   make(map[string]int),
		errors:  make(map[string]error),
	}
}

// BucketExists checks if a bucket exists.
func (m *MockMinIOClient) BucketExists(ctx context.Context, bucketName string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.trackCall("BucketExists")

	if err := m.errors["BucketExists"]; err != nil {
		return false, err
	}

	return m.buckets[bucketName], nil
}

// MakeBucket creates a new bucket.
func (m *MockMinIOClient) MakeBucket(ctx context.Context, bucketName string, opts interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.trackCall("MakeBucket")

	if err := m.errors["MakeBucket"]; err != nil {
		return err
	}

	m.buckets[bucketName] = true
	return nil
}

// PutObject uploads an object to the bucket.
func (m *MockMinIOClient) PutObject(ctx context.Context, bucketName, objectName string, data []byte, size int64, contentType string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.trackCall("PutObject")

	if err := m.errors["PutObject"]; err != nil {
		return err
	}

	if !m.buckets[bucketName] {
		return fmt.Errorf("bucket %s does not exist", bucketName)
	}

	key := fmt.Sprintf("%s/%s", bucketName, objectName)
	m.objects[key] = data
	return nil
}

// GetObject retrieves an object from the bucket.
func (m *MockMinIOClient) GetObject(ctx context.Context, bucketName, objectName string) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.trackCall("GetObject")

	if err := m.errors["GetObject"]; err != nil {
		return nil, err
	}

	key := fmt.Sprintf("%s/%s", bucketName, objectName)
	data, exists := m.objects[key]
	if !exists {
		return nil, fmt.Errorf("object %s does not exist", objectName)
	}

	return data, nil
}

// RemoveObject deletes an object from the bucket.
func (m *MockMinIOClient) RemoveObject(ctx context.Context, bucketName, objectName string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.trackCall("RemoveObject")

	if err := m.errors["RemoveObject"]; err != nil {
		return err
	}

	key := fmt.Sprintf("%s/%s", bucketName, objectName)
	delete(m.objects, key)
	return nil
}

// PresignedPutObject generates a presigned URL for uploading an object.
func (m *MockMinIOClient) PresignedPutObject(ctx context.Context, bucketName, objectName string, expiry time.Duration) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.trackCall("PresignedPutObject")

	if err := m.errors["PresignedPutObject"]; err != nil {
		return "", err
	}

	if !m.buckets[bucketName] {
		return "", fmt.Errorf("bucket %s does not exist", bucketName)
	}

	// Return a mock presigned URL
	return fmt.Sprintf("https://mock-minio.example.com/%s/%s?presigned=true&expires=%d",
		bucketName, objectName, expiry.Seconds()), nil
}

// PresignedGetObject generates a presigned URL for downloading an object.
func (m *MockMinIOClient) PresignedGetObject(ctx context.Context, bucketName, objectName string, expiry time.Duration, reqParams map[string][]string) (string, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.trackCall("PresignedGetObject")

	if err := m.errors["PresignedGetObject"]; err != nil {
		return "", err
	}

	// Return a mock presigned URL
	return fmt.Sprintf("https://mock-minio.example.com/%s/%s?presigned=true&expires=%d",
		bucketName, objectName, expiry.Seconds()), nil
}

// StatObject returns information about an object.
func (m *MockMinIOClient) StatObject(ctx context.Context, bucketName, objectName string) (int64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.trackCall("StatObject")

	if err := m.errors["StatObject"]; err != nil {
		return 0, err
	}

	key := fmt.Sprintf("%s/%s", bucketName, objectName)
	data, exists := m.objects[key]
	if !exists {
		return 0, fmt.Errorf("object %s does not exist", objectName)
	}

	return int64(len(data)), nil
}

// Reset clears all data in the mock client.
func (m *MockMinIOClient) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.objects = make(map[string][]byte)
	m.buckets = make(map[string]bool)
	m.calls = make(map[string]int)
	m.errors = make(map[string]error)
}

// SetError configures the mock to return an error for a specific method.
func (m *MockMinIOClient) SetError(method string, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.errors[method] = err
}

// ClearErrors removes all configured errors.
func (m *MockMinIOClient) ClearErrors() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.errors = make(map[string]error)
}

// GetCallCount returns the number of times a method was called.
func (m *MockMinIOClient) GetCallCount(method string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.calls[method]
}

// ObjectExists checks if an object exists in the bucket.
func (m *MockMinIOClient) ObjectExists(bucketName, objectName string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	key := fmt.Sprintf("%s/%s", bucketName, objectName)
	_, exists := m.objects[key]
	return exists
}

// GetObjectCount returns the number of objects in a bucket.
func (m *MockMinIOClient) GetObjectCount(bucketName string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	count := 0
	prefix := bucketName + "/"
	for key := range m.objects {
		if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
			count++
		}
	}
	return count
}

// trackCall increments the call counter for a method (must be called with lock held).
func (m *MockMinIOClient) trackCall(method string) {
	m.calls[method]++
}
