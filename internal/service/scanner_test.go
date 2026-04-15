package service

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScannerService_IsMusicFile(t *testing.T) {
	svc := NewScannerService(nil, nil)

	tests := []struct {
		path     string
		expected bool
	}{
		{"/music/song.mp3", true},
		{"/music/song.flac", true},
		{"/music/song.ogg", true},
		{"/music/song.m4a", true},
		{"/music/song.wav", true},
		{"/music/song.txt", false},
		{"/music/song.jpg", false},
		{"/music/song", false},
	}

	for _, tt := range tests {
		result := svc.IsMusicFile(tt.path)
		if result != tt.expected {
			t.Errorf("IsMusicFile(%s) = %v, want %v", tt.path, result, tt.expected)
		}
	}
}

func TestScannerService_GetMusicFiles(t *testing.T) {
	// 创建临时目录结构
	tmpDir, err := os.MkdirTemp("", "music-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// 创建测试文件
	musicFiles := []string{
		filepath.Join(tmpDir, "song1.mp3"),
		filepath.Join(tmpDir, "song2.flac"),
		filepath.Join(tmpDir, "folder", "song3.ogg"),
	}
	nonMusicFiles := []string{
		filepath.Join(tmpDir, "readme.txt"),
		filepath.Join(tmpDir, "cover.jpg"),
	}

	for _, f := range append(musicFiles, nonMusicFiles...) {
		dir := filepath.Dir(f)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create dir: %v", err)
		}
		if err := os.WriteFile(f, []byte("test"), 0644); err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
	}

	svc := NewScannerService(nil, nil)

	files, err := svc.GetMusicFiles(tmpDir)
	if err != nil {
		t.Fatalf("GetMusicFiles failed: %v", err)
	}

	if len(files) != len(musicFiles) {
		t.Errorf("Expected %d music files, got %d", len(musicFiles), len(files))
	}
}

func TestScannerService_GetMusicFilesWithErrors(t *testing.T) {
	// 创建临时目录
	tmpDir, err := os.MkdirTemp("", "music-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// 创建测试音乐文件
	testFile := filepath.Join(tmpDir, "test.mp3")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	svc := NewScannerService(nil, nil)

	files, errors, err := svc.GetMusicFilesWithErrors(tmpDir)
	if err != nil {
		t.Fatalf("GetMusicFilesWithErrors failed: %v", err)
	}

	if len(files) != 1 {
		t.Errorf("Expected 1 file, got %d", len(files))
	}

	if len(errors) != 0 {
		t.Errorf("Expected no errors, got %d", len(errors))
	}
}

func TestScannerService_GetFileModTime(t *testing.T) {
	svc := NewScannerService(nil, nil)

	// 创建临时文件
	tmpFile, err := os.CreateTemp("", "modtime-test-")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	modTime, err := svc.GetFileModTime(tmpFile.Name())
	if err != nil {
		t.Fatalf("GetFileModTime failed: %v", err)
	}

	if modTime == 0 {
		t.Error("Expected non-zero modification time")
	}
}

func TestScannerService_GetFileModTime_NonExistent(t *testing.T) {
	svc := NewScannerService(nil, nil)

	_, err := svc.GetFileModTime("/non/existent/file.mp3")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestScannerService_SetLastScanTime(t *testing.T) {
	svc := NewScannerService(nil, nil)

	svc.SetLastScanTime(1234567890)
	// Verify by checking internal state is set (间接验证)
	if svc.lastScanTime != 1234567890 {
		t.Errorf("Expected lastScanTime 1234567890, got %d", svc.lastScanTime)
	}
}

func TestScannerService_ScanMode_Constants(t *testing.T) {
	if ScanModeFull != "full" {
		t.Errorf("Expected ScanModeFull to be 'full', got '%s'", ScanModeFull)
	}
	if ScanModeIncremental != "incremental" {
		t.Errorf("Expected ScanModeIncremental to be 'incremental', got '%s'", ScanModeIncremental)
	}
}
