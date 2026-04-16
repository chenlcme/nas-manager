package id3

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// MusicMetadata - 音乐元数据
type MusicMetadata struct {
	Title    string
	Artist   string
	Album    string
	Year     int
	TrackNum int
	Genre    string
	Duration int // 秒
	Cover    []byte // 封面图片数据
	Lyrics   string
	FileHash string
	FileSize int64
}

// Parser - ID3 解析器
type Parser struct{}

// NewParser - 创建解析器
func NewParser() *Parser {
	return &Parser{}
}

// ParseFile - 解析音乐文件
func (p *Parser) ParseFile(filePath string) (*MusicMetadata, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// 获取文件大小
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}
	fileSize := fileInfo.Size()

	// 计算文件哈希
	fileHash, err := p.CalculateHash(file)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate hash: %w", err)
	}

	// 重置文件指针
	file.Seek(0, io.SeekStart)

	// 获取文件扩展名
	ext := strings.ToLower(filepath.Ext(filePath))

	// 解析 ID3 标签
	metadata := &MusicMetadata{
		FileHash: fileHash,
		FileSize: fileSize,
	}

	// 根据文件类型解析
	switch ext {
	case ".mp3":
		p.parseMP3(file, metadata)
	case ".flac":
		p.parseFLAC(file, metadata)
	case ".ogg":
		p.parseOGG(file, metadata)
	case ".m4a":
		p.parseM4A(file, metadata)
	case ".wav":
		p.parseWAV(file, metadata)
	default:
		// 对于不支持的格式，只提取基本信息
		metadata.Duration = p.estimateDuration(fileSize, ext)
	}

	return metadata, nil
}

// CalculateHash - 计算文件 SHA256 哈希
func (p *Parser) CalculateHash(file *os.File) (string, error) {
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// parseMP3 - 解析 MP3 文件的 ID3 标签
func (p *Parser) parseMP3(file *os.File, metadata *MusicMetadata) {
	// 读取文件开头
	header := make([]byte, 10)
	if _, err := file.Read(header); err != nil {
		return
	}

	// 检查 ID3v2 标签
	if string(header[0:3]) == "ID3" {
		p.parseID3v2(file, metadata, header)
	}

	// 估算时长（基于文件大小）
	metadata.Duration = p.estimateDuration(metadata.FileSize, ".mp3")
}

// parseID3v2 - 解析 ID3v2 标签
func (p *Parser) parseID3v2(file *os.File, metadata *MusicMetadata, header []byte) {
	// 获取标签大小
	size := int(header[6])<<21 | int(header[7])<<14 | int(header[8])<<7 | int(header[9])

	// 跳过到帧数据
	file.Seek(10, io.SeekStart)

	frames := make(map[string][]byte)
	for pos := int64(0); pos < int64(size)-10; {
		frameHeader := make([]byte, 10)
		if _, err := file.Read(frameHeader); err != nil {
			break
		}
		pos += 10

		// 检查帧 ID
		if frameHeader[0] == 0 {
			break // 填充字节
		}

		frameID := string(frameHeader[0:4])
		frameSize := int(frameHeader[4])<<24 | int(frameHeader[5])<<16 | int(frameHeader[6])<<8 | int(frameHeader[7])

		if frameSize <= 0 || frameSize > 1024*1024 {
			break
		}

		frameData := make([]byte, frameSize)
		if _, err := file.Read(frameData); err != nil {
			break
		}
		pos += int64(frameSize)

		frames[frameID] = frameData
	}

	// 提取常用标签
	if title, ok := frames["TIT2"]; ok {
		metadata.Title = p.decodeText(title)
	}
	if artist, ok := frames["TPE1"]; ok {
		metadata.Artist = p.decodeText(artist)
	}
	if album, ok := frames["TALB"]; ok {
		metadata.Album = p.decodeText(album)
	}
	if year, ok := frames["TYER"]; ok {
		metadata.Year, _ = strconv.Atoi(p.decodeText(year))
	}
	if track, ok := frames["TRCK"]; ok {
		trackStr := p.decodeText(track)
		if parts := strings.Split(trackStr, "/"); len(parts) > 0 {
			metadata.TrackNum, _ = strconv.Atoi(parts[0])
		}
	}
	if genre, ok := frames["TCON"]; ok {
		metadata.Genre = p.decodeText(genre)
	}
	if lyrics, ok := frames["USLT"]; ok {
		metadata.Lyrics = p.decodeText(lyrics)
	}

	// 提取封面
	if apic, ok := frames["APIC"]; ok {
		metadata.Cover = apic
	}
}

// decodeText - 解码文本（简化版，实际需要处理多种编码）
func (p *Parser) decodeText(data []byte) string {
	if len(data) == 0 {
		return ""
	}

	// 跳过编码字节
	offset := 1
	if data[0] == 0x00 {
		// ISO-8859-1
		return string(data[offset:])
	} else if data[0] == 0x01 {
		// UTF-16
		return p.decodeUTF16(data[offset:])
	} else if data[0] == 0x02 {
		// UTF-16BE
		return p.decodeUTF16BE(data[offset:])
	} else if data[0] == 0x03 {
		// UTF-8
		return string(data[offset:])
	}

	return string(data)
}

// decodeUTF16 - 解码 UTF-16 (LE)
func (p *Parser) decodeUTF16(data []byte) string {
	if len(data) < 2 {
		return ""
	}
	// 奇数长度时丢弃最后一个字节
	if len(data)%2 != 0 {
		data = data[:len(data)-1]
	}
	runes := make([]rune, 0, len(data)/2)
	for i := 0; i < len(data)-1; i += 2 {
		runes = append(runes, rune(data[i])|rune(data[i+1])<<8)
	}
	return string(runes)
}

// decodeUTF16BE - 解码 UTF-16BE
func (p *Parser) decodeUTF16BE(data []byte) string {
	if len(data) < 2 {
		return ""
	}
	// 奇数长度时丢弃最后一个字节
	if len(data)%2 != 0 {
		data = data[:len(data)-1]
	}
	runes := make([]rune, 0, len(data)/2)
	for i := 0; i < len(data)-1; i += 2 {
		runes = append(runes, rune(data[i])<<8|rune(data[i+1]))
	}
	return string(runes)
}

// parseFLAC - 解析 FLAC 文件
func (p *Parser) parseFLAC(file *os.File, metadata *MusicMetadata) {
	// 读取 FLAC 头部
	header := make([]byte, 4)
	if _, err := file.Read(header); err != nil {
		return
	}

	if string(header[0:4]) != "fLaC" {
		return
	}

	// 读取元数据块
	for {
		blockHeader := make([]byte, 4)
		if _, err := file.Read(blockHeader); err != nil {
			break
		}

		isLast := (blockHeader[0] & 0x80) != 0
		blockType := blockHeader[0] & 0x7F
		blockSize := int(blockHeader[1])<<16 | int(blockHeader[2])<<8 | int(blockHeader[3])

		if blockType == 4 {
			// Vorbis 评论块
			p.parseVorbisComment(file, metadata, blockSize)
		} else if blockType == 6 {
			// 图片块
			pictureData := make([]byte, blockSize)
			if _, err := file.Read(pictureData); err != nil {
				break
			}
			metadata.Cover = pictureData
		} else {
			file.Seek(int64(blockSize), io.SeekCurrent)
		}

		if isLast {
			break
		}
	}

	// 估算时长
	metadata.Duration = p.estimateDuration(metadata.FileSize, ".flac")
}

// parseVorbisComment - 解析 Vorbis 评论
func (p *Parser) parseVorbisComment(file *os.File, metadata *MusicMetadata, blockSize int) {
	// 跳过 vendor 长度和数据
	vendorLen := make([]byte, 4)
	if _, err := file.Read(vendorLen); err != nil {
		return
	}
	vLen := int(vendorLen[0]) | int(vendorLen[1])<<8 | int(vendorLen[2])<<16 | int(vendorLen[3])<<24
	file.Seek(int64(vLen), io.SeekCurrent)

	// 读取评论数量
	commentCountData := make([]byte, 4)
	if _, err := file.Read(commentCountData); err != nil {
		return
	}
	count := int(commentCountData[0]) | int(commentCountData[1])<<8 | int(commentCountData[2])<<16 | int(commentCountData[3])<<24

	for i := 0; i < count; i++ {
		lenData := make([]byte, 4)
		if _, err := file.Read(lenData); err != nil {
			break
		}
		commentLen := int(lenData[0]) | int(lenData[1])<<8 | int(lenData[2])<<16 | int(lenData[3])<<24

		comment := make([]byte, commentLen)
		if _, err := file.Read(comment); err != nil {
			break
		}

		// 解析 key=value 格式
		if idx := bytesIndex(comment, '='); idx != -1 {
			key := string(comment[:idx])
			value := string(comment[idx+1:])

			switch strings.ToUpper(key) {
			case "TITLE":
				metadata.Title = value
			case "ARTIST":
				metadata.Artist = value
			case "ALBUM":
				metadata.Album = value
			case "DATE", "YEAR":
				metadata.Year, _ = strconv.Atoi(value)
			case "TRACKNUMBER":
				metadata.TrackNum, _ = strconv.Atoi(value)
			case "GENRE":
				metadata.Genre = value
			case "LYRICS":
				metadata.Lyrics = value
			}
		}
	}
}

// bytesIndex - 查找字节切片中的字符位置
func bytesIndex(data []byte, target byte) int {
	for i, b := range data {
		if b == target {
			return i
		}
	}
	return -1
}

// parseOGG - 解析 OGG 文件（简化）
func (p *Parser) parseOGG(file *os.File, metadata *MusicMetadata) {
	metadata.Duration = p.estimateDuration(metadata.FileSize, ".ogg")
}

// parseM4A - 解析 M4A 文件（简化）
func (p *Parser) parseM4A(file *os.File, metadata *MusicMetadata) {
	metadata.Duration = p.estimateDuration(metadata.FileSize, ".m4a")
}

// parseWAV - 解析 WAV 文件（简化）
func (p *Parser) parseWAV(file *os.File, metadata *MusicMetadata) {
	// 读取 WAV 头部
	header := make([]byte, 44)
	if _, err := file.Read(header); err != nil {
		return
	}

	// 解析采样率、位深、通道数计算时长
	if string(header[0:4]) == "RIFF" && string(header[8:12]) == "WAVE" {
		sampleRate := binary.LittleEndian.Uint32(header[24:28])
		bitsPerSample := binary.LittleEndian.Uint16(header[34:36])
		numChannels := binary.LittleEndian.Uint16(header[22:24])
		dataSize := binary.LittleEndian.Uint32(header[40:44])

		if sampleRate > 0 && bitsPerSample > 0 && numChannels > 0 {
			bytesPerSample := bitsPerSample / 8
			totalSamples := dataSize / (uint32(bytesPerSample) * uint32(numChannels))
			metadata.Duration = int(totalSamples / sampleRate)
		}
	}
}

// estimateDuration - 根据文件大小估算时长
func (p *Parser) estimateDuration(fileSize int64, ext string) int {
	// 根据文件扩展名估算平均比特率
	var bitrate int64
	switch ext {
	case ".mp3":
		bitrate = 128 * 1024 / 8 // 128kbps
	case ".flac":
		bitrate = 800 * 1024 / 8 // 800kbps
	case ".ogg":
		bitrate = 128 * 1024 / 8
	case ".m4a":
		bitrate = 128 * 1024 / 8
	case ".wav":
		bitrate = 1411 * 1024 / 8 // CD quality
	default:
		bitrate = 128 * 1024 / 8
	}

	if bitrate == 0 {
		return 0
	}
	return int(fileSize / bitrate)
}

// GetDuration - 获取音频文件时长（秒）
func (p *Parser) GetDuration(filePath string) (int, error) {
	metadata, err := p.ParseFile(filePath)
	if err != nil {
		return 0, err
	}
	return metadata.Duration, nil
}

// CalculateFileHash - 计算文件哈希
func (p *Parser) CalculateFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	return p.CalculateHash(file)
}
