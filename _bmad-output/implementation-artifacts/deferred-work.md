# Deferred Work

## Deferred from: code review of story 2-3-folder-view-browse (2026-04-16)

- **全表扫描加载所有歌曲到内存** [folder.go:48] — Known trade-off to avoid SQLite REVERSE function limitation. Current approach acceptable for local NAS with moderate song counts.
- **无分页支持** [folder.go:26,47] — Local NAS scope limits data volume. Consider pagination if scale increases.
- **动态 ID 排序变化时不稳定** [folder.go:86] — Consistent with existing ArtistRepository and AlbumRepository pattern. IDs are display-only, not persistent references.
- **. 文件夹名处理** [folder.go:33] — Edge case: songs with path "./song.mp3" rare. Not worth special handling.
- **非 ASCII Unicode 大小写排序** [folder.go:69,73] — Go's strings.ToLower has locale-specific behavior. Current implementation acceptable for most CJK users.
