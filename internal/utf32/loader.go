// internal/utf32/loader.go
// Tải toàn bộ bảng Unicode từ unicode.org, lưu local
// Đây là Thân cây bất biến — nền tảng của toàn bộ ISL

package utf32

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	// Nguồn chính thức từ Unicode Consortium
	UnicodeDataURL = "https://www.unicode.org/Public/UCD/latest/ucd/UnicodeData.txt"

	// Lưu local để dùng offline
	LocalCacheFile = "./data/utf32/UnicodeData.txt"
)

// UnicodeEntry là một dòng trong UnicodeData.txt
type UnicodeEntry struct {
	CodePoint  rune   // Mã Unicode (VD: U+1F34C = 🍌)
	Name       string // Tên chính thức (VD: BANANA)
	Category   string // Loại (Lu=Uppercase, Ll=Lowercase, So=Symbol...)
	Block      string // Block (Basic Latin, CJK, Emoji...)
	ISLAddr    string // Địa chỉ ISL được gán (do LeoAI tạo)
}

// DB là cơ sở dữ liệu UTF-32 toàn bộ
type DB struct {
	entries  map[rune]*UnicodeEntry   // codepoint → entry
	byName   map[string]rune          // tên → codepoint
	byBlock  map[string][]rune        // block → danh sách codepoint
	loaded   bool
	loadedAt time.Time
	total    int
}

// NewDB tạo DB mới
func NewDB() *DB {
	return &DB{
		entries: make(map[rune]*UnicodeEntry),
		byName:  make(map[string]rune),
		byBlock: make(map[string][]rune),
	}
}

// Load tải dữ liệu Unicode — từ cache local nếu có, không thì tải từ internet
func (db *DB) Load() error {
	// Kiểm tra cache local
	if _, err := os.Stat(LocalCacheFile); err == nil {
		log.Println("UTF-32: Loading from local cache...")
		return db.loadFromFile(LocalCacheFile)
	}

	// Chưa có cache → tải từ internet
	log.Println("UTF-32: Downloading from unicode.org (one-time setup)...")
	if err := db.download(); err != nil {
		return fmt.Errorf("utf32: download failed: %w", err)
	}

	return db.loadFromFile(LocalCacheFile)
}

// download tải UnicodeData.txt từ unicode.org và lưu local
func (db *DB) download() error {
	// Tạo thư mục nếu chưa có
	dir := filepath.Dir(LocalCacheFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("utf32: create dir: %w", err)
	}

	// HTTP GET với timeout
	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Get(UnicodeDataURL)
	if err != nil {
		return fmt.Errorf("utf32: http get: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("utf32: http status %d", resp.StatusCode)
	}

	// Lưu vào file
	f, err := os.Create(LocalCacheFile)
	if err != nil {
		return fmt.Errorf("utf32: create file: %w", err)
	}
	defer f.Close()

	written, err := io.Copy(f, resp.Body)
	if err != nil {
		return fmt.Errorf("utf32: write file: %w", err)
	}

	log.Printf("UTF-32: Downloaded %.1f MB → %s", float64(written)/1024/1024, LocalCacheFile)
	return nil
}

// loadFromFile parse UnicodeData.txt vào DB
func (db *DB) loadFromFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("utf32: open file: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	count := 0

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		entry, err := parseLine(line)
		if err != nil {
			continue // Bỏ qua dòng lỗi
		}

		db.entries[entry.CodePoint] = entry
		db.byName[strings.ToUpper(entry.Name)] = entry.CodePoint
		db.byBlock[entry.Block] = append(db.byBlock[entry.Block], entry.CodePoint)
		count++
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("utf32: scan: %w", err)
	}

	db.loaded = true
	db.loadedAt = time.Now()
	db.total = count
	log.Printf("UTF-32: Loaded %d characters — Thân cây bất biến sẵn sàng", count)
	return nil
}

// parseLine parse một dòng UnicodeData.txt
// Format: CodePoint;Name;Category;...
// VD: 1F34C;BANANA;So;0;ON;;;;;N;;;;;
func parseLine(line string) (*UnicodeEntry, error) {
	fields := strings.Split(line, ";")
	if len(fields) < 3 {
		return nil, fmt.Errorf("invalid line")
	}

	// Parse codepoint (hex)
	cp, err := strconv.ParseInt(fields[0], 16, 32)
	if err != nil {
		return nil, err
	}

	name := fields[1]
	category := fields[2]

	// Xác định block từ codepoint
	block := codePointToBlock(rune(cp))

	return &UnicodeEntry{
		CodePoint: rune(cp),
		Name:      name,
		Category:  category,
		Block:     block,
	}, nil
}

// Lookup tìm entry theo codepoint
func (db *DB) Lookup(cp rune) (*UnicodeEntry, bool) {
	e, ok := db.entries[cp]
	return e, ok
}

// LookupByName tìm codepoint theo tên
func (db *DB) LookupByName(name string) (rune, bool) {
	cp, ok := db.byName[strings.ToUpper(name)]
	return cp, ok
}

// GetBlock lấy tất cả ký tự trong một block
func (db *DB) GetBlock(block string) []rune {
	return db.byBlock[block]
}

// Stats trả về thống kê
func (db *DB) Stats() map[string]int {
	stats := map[string]int{
		"total":  db.total,
		"blocks": len(db.byBlock),
	}
	for block, cps := range db.byBlock {
		stats["block_"+block] = len(cps)
	}
	return stats
}

// IsLoaded kiểm tra đã load chưa
func (db *DB) IsLoaded() bool {
	return db.loaded
}

// codePointToBlock xác định Unicode block từ codepoint
func codePointToBlock(cp rune) string {
	switch {
	case cp <= 0x007F:
		return "Basic Latin"
	case cp <= 0x00FF:
		return "Latin-1 Supplement"
	case cp <= 0x024F:
		return "Latin Extended"
	case cp <= 0x036F:
		return "Diacritical Marks"
	case cp <= 0x03FF:
		return "Greek"
	case cp <= 0x04FF:
		return "Cyrillic"
	case cp <= 0x05FF:
		return "Hebrew"
	case cp <= 0x06FF:
		return "Arabic"
	case cp <= 0x0FFF:
		return "Indic Scripts"
	case cp >= 0x1EA0 && cp <= 0x1EF9:
		return "Vietnamese"
	case cp >= 0x4E00 && cp <= 0x9FFF:
		return "CJK Unified"
	case cp >= 0x3040 && cp <= 0x30FF:
		return "Japanese"
	case cp >= 0xAC00 && cp <= 0xD7AF:
		return "Korean"
	case cp >= 0x1F300 && cp <= 0x1F9FF:
		return "Emoji"
	case cp >= 0x1F000 && cp <= 0x1F02F:
		return "Mahjong"
	default:
		return "Other"
	}
}
