// Song - 歌曲类型定义
export interface Song {
  id: number;
  filePath: string;
  title: string;
  artist: string;
  album: string;
  year: number;
  genre: string;
  trackNum: number;
  duration: number;
  coverPath: string;
  lyrics: string;
  fileHash: string;
  fileSize: number;
  createdAt: string;
  updatedAt: string;
}

// ArtistWithCount - 艺术家及其歌曲数量
export interface ArtistWithCount {
  id: number;
  name: string;
  songCount: number;
}

// AlbumWithCount - 专辑及其歌曲数量
export interface AlbumWithCount {
  id: number;
  name: string;
  artist: string;
  songCount: number;
}
