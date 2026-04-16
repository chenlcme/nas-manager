package repository

import (
	"testing"

	"nas-manager/internal/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func createFolderTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	db.AutoMigrate(&model.Song{})

	return db
}

func TestFolderRepository_GetAllFoldersWithSongCount(t *testing.T) {
	db := createFolderTestDB(t)
	repo := NewFolderRepository(db)

	// 创建测试数据 - 不同文件夹的歌曲
	songs := []*model.Song{
		{FilePath: "/music/rock/song1.mp3", Title: "Rock Song 1"},
		{FilePath: "/music/rock/song2.mp3", Title: "Rock Song 2"},
		{FilePath: "/music/rock/song3.mp3", Title: "Rock Song 3"},
		{FilePath: "/music/pop/song4.mp3", Title: "Pop Song 1"},
		{FilePath: "/music/pop/song5.mp3", Title: "Pop Song 2"},
		{FilePath: "/music/classical/song6.mp3", Title: "Classical Song 1"},
		{FilePath: "/downloads/song7.mp3", Title: "Downloaded Song 1"},
	}
	for _, s := range songs {
		db.Create(s)
	}

	// 测试获取文件夹列表（降序）
	folders, err := repo.GetAllFoldersWithSongCount(false)
	if err != nil {
		t.Fatalf("Failed to get folders: %v", err)
	}

	if len(folders) != 4 {
		t.Errorf("Expected 4 folders, got %d", len(folders))
	}

	// 验证文件夹及其歌曲数量
	folderCounts := make(map[string]int)
	for _, f := range folders {
		folderCounts[f.Path] = f.SongCount
	}

	if folderCounts["/music/rock"] != 3 {
		t.Errorf("Expected /music/rock to have 3 songs, got %d", folderCounts["/music/rock"])
	}
	if folderCounts["/music/pop"] != 2 {
		t.Errorf("Expected /music/pop to have 2 songs, got %d", folderCounts["/music/pop"])
	}
	if folderCounts["/music/classical"] != 1 {
		t.Errorf("Expected /music/classical to have 1 song, got %d", folderCounts["/music/classical"])
	}
	if folderCounts["/downloads"] != 1 {
		t.Errorf("Expected /downloads to have 1 song, got %d", folderCounts["/downloads"])
	}

	// 测试升序和降序返回的顺序不同（验证排序确实发生了）
	foldersAsc, err := repo.GetAllFoldersWithSongCount(true)
	if err != nil {
		t.Fatalf("Failed to get folders (asc): %v", err)
	}

	foldersDesc, err := repo.GetAllFoldersWithSongCount(false)
	if err != nil {
		t.Fatalf("Failed to get folders (desc): %v", err)
	}

	// 升序和降序的第一个文件夹应该不同（验证排序方向不同）
	if foldersAsc[0].Path == foldersDesc[0].Path {
		t.Error("Expected different first folder between asc and desc order")
	}
}

func TestFolderRepository_GetSongsByFolder(t *testing.T) {
	db := createFolderTestDB(t)
	repo := NewFolderRepository(db)

	// 创建测试数据
	songs := []*model.Song{
		{FilePath: "/music/rock/song1.mp3", Title: "Rock Song 1"},
		{FilePath: "/music/rock/song2.mp3", Title: "Rock Song 2"},
		{FilePath: "/music/pop/song3.mp3", Title: "Pop Song 1"},
	}
	for _, s := range songs {
		db.Create(s)
	}

	// 测试获取特定文件夹的歌曲
	found, err := repo.GetSongsByFolder("/music/rock", "title", "asc")
	if err != nil {
		t.Fatalf("Failed to get songs by folder: %v", err)
	}

	if len(found) != 2 {
		t.Errorf("Expected 2 songs for /music/rock, got %d", len(found))
	}

	// 验证返回的歌曲标题
	titles := make(map[string]bool)
	for _, s := range found {
		titles[s.Title] = true
	}

	if !titles["Rock Song 1"] || !titles["Rock Song 2"] {
		t.Error("Expected all /music/rock songs to be returned")
	}
}

func TestFolderRepository_GetSongsByFolder_NotFound(t *testing.T) {
	db := createFolderTestDB(t)
	repo := NewFolderRepository(db)

	found, err := repo.GetSongsByFolder("/non/existent/folder", "title", "asc")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(found) != 0 {
		t.Errorf("Expected 0 songs for non-existent folder, got %d", len(found))
	}
}

func TestFolderRepository_GetAllFoldersWithSongCount_EmptyDB(t *testing.T) {
	db := createFolderTestDB(t)
	repo := NewFolderRepository(db)

	// 空数据库应该返回空列表
	folders, err := repo.GetAllFoldersWithSongCount(false)
	if err != nil {
		t.Fatalf("Failed to get folders from empty DB: %v", err)
	}

	if len(folders) != 0 {
		t.Errorf("Expected 0 folders from empty DB, got %d", len(folders))
	}
}

func TestFolderRepository_GetFolderPathByID_ZeroID(t *testing.T) {
	db := createFolderTestDB(t)
	repo := NewFolderRepository(db)

	// 创建测试数据
	songs := []*model.Song{
		{FilePath: "/music/rock/song1.mp3", Title: "Rock Song 1"},
	}
	for _, s := range songs {
		db.Create(s)
	}

	// ID=0 应该返回错误
	_, err := repo.GetFolderPathByID(0)
	if err == nil {
		t.Error("Expected error for ID=0, got nil")
	}
}

func TestFolderRepository_GetFolderPathByID_OutOfRange(t *testing.T) {
	db := createFolderTestDB(t)
	repo := NewFolderRepository(db)

	// 创建测试数据
	songs := []*model.Song{
		{FilePath: "/music/rock/song1.mp3", Title: "Rock Song 1"},
	}
	for _, s := range songs {
		db.Create(s)
	}

	// 超出范围的 ID 应该返回错误
	_, err := repo.GetFolderPathByID(999)
	if err == nil {
		t.Error("Expected error for out-of-range ID, got nil")
	}
}

func TestFolderRepository_GetFolderPathByID(t *testing.T) {
	db := createFolderTestDB(t)
	repo := NewFolderRepository(db)

	// 创建测试数据
	songs := []*model.Song{
		{FilePath: "/music/rock/song1.mp3", Title: "Rock Song 1"},
		{FilePath: "/music/pop/song2.mp3", Title: "Pop Song 1"},
	}
	for _, s := range songs {
		db.Create(s)
	}

	// 获取文件夹列表
	folders, err := repo.GetAllFoldersWithSongCount(false)
	if err != nil {
		t.Fatalf("Failed to get folders: %v", err)
	}

	if len(folders) != 2 {
		t.Fatalf("Expected 2 folders, got %d", len(folders))
	}

	// 验证通过 ID 获取的路径与列表中的路径一致
	for _, f := range folders {
		path, err := repo.GetFolderPathByID(f.ID)
		if err != nil {
			t.Errorf("Failed to get path by ID %d: %v", f.ID, err)
		}
		if path != f.Path {
			t.Errorf("Expected path %s for ID %d, got %s", f.Path, f.ID, path)
		}
	}
}

func TestFolderRepository_GetSongsByFolder_InvalidSortBy(t *testing.T) {
	db := createFolderTestDB(t)
	repo := NewFolderRepository(db)

	// 创建测试数据
	songs := []*model.Song{
		{FilePath: "/music/rock/song1.mp3", Title: "Rock Song 1"},
	}
	for _, s := range songs {
		db.Create(s)
	}

	// 使用无效的 sort_by 参数应该返回错误
	_, err := repo.GetSongsByFolder("/music/rock", "invalid_field", "asc")
	if err == nil {
		t.Error("Expected error for invalid sort_by parameter, got nil")
	}
}

func TestFolderRepository_GetSongsByFolder_InvalidOrder(t *testing.T) {
	db := createFolderTestDB(t)
	repo := NewFolderRepository(db)

	// 创建测试数据
	songs := []*model.Song{
		{FilePath: "/music/rock/song1.mp3", Title: "Rock Song 1"},
	}
	for _, s := range songs {
		db.Create(s)
	}

	// 使用无效的 order 参数应该返回错误
	_, err := repo.GetSongsByFolder("/music/rock", "title", "invalid_order")
	if err == nil {
		t.Error("Expected error for invalid order parameter, got nil")
	}
}

func TestFolderRepository_GetSongsByFolder_ValidSortParameters(t *testing.T) {
	db := createFolderTestDB(t)
	repo := NewFolderRepository(db)

	// 创建测试数据
	songs := []*model.Song{
		{FilePath: "/music/rock/song1.mp3", Title: "Rock Song 1", Duration: 180},
		{FilePath: "/music/rock/song2.mp3", Title: "Rock Song 2", Duration: 240},
	}
	for _, s := range songs {
		db.Create(s)
	}

	// 测试有效的 sort_by 和 order 参数组合
	testCases := []struct {
		sortBy string
		order  string
	}{
		{"title", "asc"},
		{"title", "desc"},
		{"duration", "asc"},
		{"duration", "desc"},
		{"created_at", "asc"},
		{"created_at", "desc"},
	}

	for _, tc := range testCases {
		found, err := repo.GetSongsByFolder("/music/rock", tc.sortBy, tc.order)
		if err != nil {
			t.Errorf("Unexpected error for sortBy=%s, order=%s: %v", tc.sortBy, tc.order, err)
		}
		if len(found) != 2 {
			t.Errorf("Expected 2 songs for sortBy=%s, order=%s, got %d", tc.sortBy, tc.order, len(found))
		}
	}
}
