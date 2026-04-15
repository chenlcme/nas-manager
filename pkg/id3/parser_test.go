package id3

import (
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestParser_ParseFile_NonExistent(t *testing.T) {
	p := NewParser()

	_, err := p.ParseFile("/non/existent/file.mp3")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestParser_IsMusicFile(t *testing.T) {
	p := NewParser()

	tests := []struct {
		ext    string
		isMusic bool
	}{
		{".mp3", true},
		{".flac", true},
		{".ogg", true},
		{".m4a", true},
		{".wav", true},
		{".txt", false},
		{".jpg", false},
	}

	for _, tt := range tests {
		// 创建临时文件
		tmpFile, err := os.CreateTemp("", "test*"+tt.ext)
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		tmpFile.WriteString("test")
		tmpFile.Close()
		defer os.Remove(tmpFile.Name())

		metadata, err := p.ParseFile(tmpFile.Name())
		if err != nil {
			t.Errorf("ParseFile(%s) failed: %v", tt.ext, err)
			continue
		}

		if metadata.FileSize == 0 {
			t.Errorf("Expected non-zero file size for %s", tt.ext)
		}
	}
}

func TestParser_ParseMP3(t *testing.T) {
	// 创建一个最小的 MP3 文件（带 ID3v2 标签）
	tmpDir, err := os.MkdirTemp("", "id3-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tmpFile := filepath.Join(tmpDir, "test.mp3")

	// ID3v2.4 标签头
	id3Header := []byte{
		0x49, 0x44, 0x33, // "ID3"
		0x04, 0x00,       // 版本 2.4.0
		0x00,             // 标志
		0x00, 0x00, 0x00, // 标签大小 (0)
	}

	// 创建文件
	file, err := os.Create(tmpFile)
	if err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	file.Write(id3Header)
	file.Close()

	p := NewParser()

	metadata, err := p.ParseFile(tmpFile)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	if metadata.FileHash == "" {
		t.Error("Expected non-empty file hash")
	}

	if metadata.FileSize == 0 {
		t.Error("Expected non-zero file size")
	}
}

func TestParser_CalculateHash(t *testing.T) {
	p := NewParser()

	tmpFile, err := os.CreateTemp("", "hash-test-")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	content := []byte("test content")
	tmpFile.Write(content)
	tmpFile.Close()

	file, err := os.Open(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	hash, err := p.CalculateHash(file)
	if err != nil {
		t.Fatalf("CalculateHash failed: %v", err)
	}

	if hash == "" {
		t.Error("Expected non-empty hash")
	}

	// 验证哈希一致性
	file.Seek(0, io.SeekStart)
	hash2, _ := p.CalculateHash(file)
	if hash != hash2 {
		t.Error("Expected same hash for same content")
	}
}
