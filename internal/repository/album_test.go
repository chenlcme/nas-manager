package repository

import (
	"testing"

	"nas-manager/internal/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func createAlbumTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	db.AutoMigrate(&model.Song{})

	return db
}

func TestAlbumRepository_GetAllAlbumsWithSongCount(t *testing.T) {
	db := createAlbumTestDB(t)
	repo := NewAlbumRepository(db)

	// 创建测试数据
	songs := []*model.Song{
		{FilePath: "/test/song1.mp3", Title: "晴天", Album: "叶惠美", Artist: "周杰伦"},
		{FilePath: "/test/song2.mp3", Title: "七里香", Album: "七里香", Artist: "周杰伦"},
		{FilePath: "/test/song3.mp3", Title: "江南", Album: "江南", Artist: "林俊杰"},
		{FilePath: "/test/song4.mp3", Title: "编号89757", Album: "编号89757", Artist: "林俊杰"},
		{FilePath: "/test/song5.mp3", Title: "一千年以后", Album: "江南", Artist: "林俊杰"},
	}
	for _, s := range songs {
		db.Create(s)
	}

	// 测试获取专辑列表（降序）
	albums, err := repo.GetAllAlbumsWithSongCount(false)
	if err != nil {
		t.Fatalf("Failed to get albums: %v", err)
	}

	if len(albums) != 4 {
		t.Errorf("Expected 4 albums, got %d", len(albums))
	}

	// 验证专辑及其歌曲数量
	albumCounts := make(map[string]int)
	albumArtists := make(map[string]string)
	for _, a := range albums {
		albumCounts[a.Name] = a.SongCount
		albumArtists[a.Name] = a.Artist
	}

	if albumCounts["叶惠美"] != 1 {
		t.Errorf("Expected 叶惠美 to have 1 song, got %d", albumCounts["叶惠美"])
	}
	if albumCounts["七里香"] != 1 {
		t.Errorf("Expected 七里香 to have 1 song, got %d", albumCounts["七里香"])
	}
	if albumCounts["江南"] != 2 {
		t.Errorf("Expected 江南 to have 2 songs, got %d", albumCounts["江南"])
	}

	// 验证专辑对应艺术家
	if albumArtists["江南"] != "林俊杰" {
		t.Errorf("Expected 江南 artist to be 林俊杰, got %s", albumArtists["江南"])
	}

	// 测试升序和降序返回的顺序不同（验证排序确实发生了）
	albumsAsc, err := repo.GetAllAlbumsWithSongCount(true)
	if err != nil {
		t.Fatalf("Failed to get albums (asc): %v", err)
	}

	albumsDesc, err := repo.GetAllAlbumsWithSongCount(false)
	if err != nil {
		t.Fatalf("Failed to get albums (desc): %v", err)
	}

	// 升序和降序的第一个专辑应该不同（验证排序方向不同）
	if albumsAsc[0].Name == albumsDesc[0].Name {
		t.Error("Expected different first album between asc and desc order")
	}
}

func TestAlbumRepository_GetSongsByAlbum(t *testing.T) {
	db := createAlbumTestDB(t)
	repo := NewAlbumRepository(db)

	// 创建测试数据
	songs := []*model.Song{
		{FilePath: "/test/song1.mp3", Title: "晴天", Album: "叶惠美", Artist: "周杰伦"},
		{FilePath: "/test/song2.mp3", Title: "稻香", Album: "叶惠美", Artist: "周杰伦"},
		{FilePath: "/test/song3.mp3", Title: "七里香", Album: "七里香", Artist: "周杰伦"},
		{FilePath: "/test/song4.mp3", Title: "江南", Album: "江南", Artist: "林俊杰"},
	}
	for _, s := range songs {
		db.Create(s)
	}

	// 测试获取特定专辑的歌曲
	found, err := repo.GetSongsByAlbum("叶惠美", "周杰伦", "title", "asc")
	if err != nil {
		t.Fatalf("Failed to get songs by album: %v", err)
	}

	if len(found) != 2 {
		t.Errorf("Expected 2 songs for 叶惠美 by 周杰伦, got %d", len(found))
	}

	// 验证返回的歌曲标题
	titles := make(map[string]bool)
	for _, s := range found {
		titles[s.Title] = true
	}

	if !titles["晴天"] || !titles["稻香"] {
		t.Error("Expected all 叶惠美 songs to be returned")
	}
}

func TestAlbumRepository_GetSongsByAlbum_NotFound(t *testing.T) {
	db := createAlbumTestDB(t)
	repo := NewAlbumRepository(db)

	found, err := repo.GetSongsByAlbum("不存在的专辑", "不存在的艺术家", "title", "asc")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(found) != 0 {
		t.Errorf("Expected 0 songs for non-existent album, got %d", len(found))
	}
}

func TestAlbumRepository_GetAllAlbumsWithSongCount_EmptyAlbum(t *testing.T) {
	db := createAlbumTestDB(t)
	repo := NewAlbumRepository(db)

	// 创建测试数据，包含空专辑
	songs := []*model.Song{
		{FilePath: "/test/song1.mp3", Title: "Song 1", Album: "叶惠美", Artist: "周杰伦"},
		{FilePath: "/test/song2.mp3", Title: "Song 2", Album: ""},
		{FilePath: "/test/song3.mp3", Title: "Song 3", Album: "   "},
	}
	for _, s := range songs {
		db.Create(s)
	}

	// 空专辑不应被返回
	albums, err := repo.GetAllAlbumsWithSongCount(false)
	if err != nil {
		t.Fatalf("Failed to get albums: %v", err)
	}

	if len(albums) != 1 {
		t.Errorf("Expected 1 album (空专辑应被过滤), got %d", len(albums))
	}

	if albums[0].Name != "叶惠美" {
		t.Errorf("Expected album to be 叶惠美, got %s", albums[0].Name)
	}
}

func TestAlbumRepository_GetAllAlbumsWithSongCount_SameAlbumDifferentArtist(t *testing.T) {
	db := createAlbumTestDB(t)
	repo := NewAlbumRepository(db)

	// 创建测试数据 - 同专辑名不同艺术家（如合辑）
	songs := []*model.Song{
		{FilePath: "/test/song1.mp3", Title: "Song 1", Album: "精选集", Artist: "周杰伦"},
		{FilePath: "/test/song2.mp3", Title: "Song 2", Album: "精选集", Artist: "林俊杰"},
		{FilePath: "/test/song3.mp3", Title: "Song 3", Album: "精选集", Artist: "王力宏"},
	}
	for _, s := range songs {
		db.Create(s)
	}

	// 同专辑名不同艺术家应该被视为不同专辑
	albums, err := repo.GetAllAlbumsWithSongCount(false)
	if err != nil {
		t.Fatalf("Failed to get albums: %v", err)
	}

	if len(albums) != 3 {
		t.Errorf("Expected 3 albums (same album name with different artists), got %d", len(albums))
	}

	// 每个艺术家都应该有自己的"精选集"
	artistCounts := make(map[string]int)
	for _, a := range albums {
		artistCounts[a.Artist]++
	}

	if artistCounts["周杰伦"] != 1 || artistCounts["林俊杰"] != 1 || artistCounts["王力宏"] != 1 {
		t.Error("Expected each artist to have their own album entry")
	}
}

func TestAlbumRepository_GetAlbumNameAndArtistByID_ZeroID(t *testing.T) {
	db := createAlbumTestDB(t)
	repo := NewAlbumRepository(db)

	// 创建测试数据
	songs := []*model.Song{
		{FilePath: "/test/song1.mp3", Title: "Song 1", Album: "专辑A", Artist: "艺术家A"},
	}
	for _, s := range songs {
		db.Create(s)
	}

	// ID=0 应该返回错误
	_, _, err := repo.GetAlbumNameAndArtistByID(0)
	if err == nil {
		t.Error("Expected error for ID=0, got nil")
	}
}

func TestAlbumRepository_GetAlbumNameAndArtistByID_OutOfRange(t *testing.T) {
	db := createAlbumTestDB(t)
	repo := NewAlbumRepository(db)

	// 创建测试数据
	songs := []*model.Song{
		{FilePath: "/test/song1.mp3", Title: "Song 1", Album: "专辑A", Artist: "艺术家A"},
	}
	for _, s := range songs {
		db.Create(s)
	}

	// 超出范围的 ID 应该返回错误
	_, _, err := repo.GetAlbumNameAndArtistByID(999)
	if err == nil {
		t.Error("Expected error for out-of-range ID, got nil")
	}
}

func TestAlbumRepository_GetSongsByAlbum_InvalidSortBy(t *testing.T) {
	db := createAlbumTestDB(t)
	repo := NewAlbumRepository(db)

	// 创建测试数据
	songs := []*model.Song{
		{FilePath: "/test/song1.mp3", Title: "晴天", Album: "叶惠美", Artist: "周杰伦"},
	}
	for _, s := range songs {
		db.Create(s)
	}

	// 使用无效的 sort_by 参数应该返回错误
	_, err := repo.GetSongsByAlbum("叶惠美", "周杰伦", "invalid_field", "asc")
	if err == nil {
		t.Error("Expected error for invalid sort_by parameter, got nil")
	}
}

func TestAlbumRepository_GetSongsByAlbum_InvalidOrder(t *testing.T) {
	db := createAlbumTestDB(t)
	repo := NewAlbumRepository(db)

	// 创建测试数据
	songs := []*model.Song{
		{FilePath: "/test/song1.mp3", Title: "晴天", Album: "叶惠美", Artist: "周杰伦"},
	}
	for _, s := range songs {
		db.Create(s)
	}

	// 使用无效的 order 参数应该返回错误
	_, err := repo.GetSongsByAlbum("叶惠美", "周杰伦", "title", "invalid_order")
	if err == nil {
		t.Error("Expected error for invalid order parameter, got nil")
	}
}

func TestAlbumRepository_GetSongsByAlbum_ValidSortParameters(t *testing.T) {
	db := createAlbumTestDB(t)
	repo := NewAlbumRepository(db)

	// 创建测试数据
	songs := []*model.Song{
		{FilePath: "/test/song1.mp3", Title: "晴天", Album: "叶惠美", Artist: "周杰伦", Duration: 180},
		{FilePath: "/test/song2.mp3", Title: "稻香", Album: "叶惠美", Artist: "周杰伦", Duration: 240},
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
		found, err := repo.GetSongsByAlbum("叶惠美", "周杰伦", tc.sortBy, tc.order)
		if err != nil {
			t.Errorf("Unexpected error for sortBy=%s, order=%s: %v", tc.sortBy, tc.order, err)
		}
		if len(found) != 2 {
			t.Errorf("Expected 2 songs for sortBy=%s, order=%s, got %d", tc.sortBy, tc.order, len(found))
		}
	}
}
