// internal/utf32/mapper.go
// Ánh xạ khái niệm ngôn ngữ tự nhiên → ISL address
// Dựa trên Unicode category và block để xác định domain ISL

package utf32

import (
	"fmt"
	"strings"
	"sync"

	"github.com/goldlotus1810/HomeOS/internal/isl"
)

// Mapper ánh xạ Unicode codepoint → ISL address
type Mapper struct {
	mu      sync.RWMutex
	db      *DB
	cache   map[rune]isl.Address   // codepoint → ISL address (cache)
	reverse map[string]rune        // ISL string → codepoint
}

// NewMapper tạo mapper mới
func NewMapper(db *DB) *Mapper {
	m := &Mapper{
		db:      db,
		cache:   make(map[rune]isl.Address),
		reverse: make(map[string]rune),
	}
	return m
}

// Map ánh xạ một codepoint → ISL address
// Đây là hàm cốt lõi: Unicode → ISL domain
func (m *Mapper) Map(cp rune) (isl.Address, error) {
	m.mu.RLock()
	if addr, ok := m.cache[cp]; ok {
		m.mu.RUnlock()
		return addr, nil
	}
	m.mu.RUnlock()

	entry, ok := m.db.Lookup(cp)
	if !ok {
		return isl.Address{}, fmt.Errorf("mapper: codepoint U+%04X not found", cp)
	}

	addr := m.entryToISL(entry)

	m.mu.Lock()
	m.cache[cp] = addr
	m.reverse[addr.String()] = cp
	m.mu.Unlock()

	return addr, nil
}

// MapConcept ánh xạ tên khái niệm → ISL address
// Ví dụ: "BANANA" → ISL address của 🍌
func (m *Mapper) MapConcept(concept string) (isl.Address, error) {
	// Tìm theo tên Unicode chính thức
	cp, ok := m.db.LookupByName(strings.ToUpper(concept))
	if !ok {
		return isl.Address{}, fmt.Errorf("mapper: concept '%s' not found in UTF-32", concept)
	}
	return m.Map(cp)
}

// Reverse tìm codepoint từ ISL address
func (m *Mapper) Reverse(addr isl.Address) (rune, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	cp, ok := m.reverse[addr.String()]
	return cp, ok
}

// entryToISL chuyển UnicodeEntry → ISL address
// Logic:
//   Layer = Unicode domain (A=Visual/Emoji, B=Organism, C=Color...)
//   Group = Unicode block group (a=Latin, b=CJK, c=Emoji...)
//   Type  = Unicode category (a=letter, b=symbol, c=number...)
//   ID    = sequential ID trong nhóm
func (m *Mapper) entryToISL(entry *UnicodeEntry) isl.Address {
	layer := blockToISLLayer(entry.Block)
	group := blockToISLGroup(entry.Block)
	typ := categoryToISLType(entry.Category)
	id := byte(entry.CodePoint & 0xFF) // Simplified: dùng byte thấp

	return isl.Address{
		Layer: layer,
		Group: group,
		Type:  typ,
		ID:    id,
	}
}

// blockToISLLayer ánh xạ Unicode block → ISL Layer (A-Z)
// Layer = domain ngữ nghĩa cao nhất
func blockToISLLayer(block string) byte {
	switch block {
	case "Basic Latin", "Latin-1 Supplement", "Latin Extended":
		return 'L' // Language/Text
	case "Vietnamese":
		return 'L' // Language
	case "CJK Unified":
		return 'L' // Language
	case "Japanese":
		return 'L'
	case "Korean":
		return 'L'
	case "Arabic":
		return 'L'
	case "Hebrew":
		return 'L'
	case "Greek":
		return 'L'
	case "Cyrillic":
		return 'L'
	case "Emoji":
		return 'E' // Emoji/Expression
	case "Mahjong":
		return 'G' // Game/Symbol
	case "Indic Scripts":
		return 'L'
	default:
		return 'Z' // Unknown
	}
}

// blockToISLGroup ánh xạ block → ISL Group (A-Z)
func blockToISLGroup(block string) byte {
	switch block {
	case "Basic Latin":
		return 'A' // ASCII base
	case "Latin-1 Supplement":
		return 'B'
	case "Latin Extended":
		return 'C'
	case "Vietnamese":
		return 'V'
	case "CJK Unified":
		return 'Z' // Z for CJK (large set)
	case "Japanese":
		return 'J'
	case "Korean":
		return 'K'
	case "Arabic":
		return 'R' // aRabic
	case "Hebrew":
		return 'H'
	case "Greek":
		return 'G'
	case "Cyrillic":
		return 'Y'
	case "Emoji":
		return 'E'
	default:
		return 'O' // Other
	}
}

// categoryToISLType ánh xạ Unicode category → ISL Type (a-z)
// Unicode categories: Lu=Uppercase Letter, Ll=Lowercase, So=Symbol...
func categoryToISLType(cat string) byte {
	if len(cat) < 1 {
		return 'z'
	}
	switch cat[0] {
	case 'L': // Letter
		return 'a'
	case 'N': // Number
		return 'b'
	case 'P': // Punctuation
		return 'c'
	case 'S': // Symbol
		return 'd'
	case 'Z': // Separator
		return 'e'
	case 'C': // Control
		return 'f'
	case 'M': // Mark
		return 'g'
	default:
		return 'z'
	}
}

// PreloadAll tạo ISL address cho toàn bộ Unicode (chạy khi khởi động)
// Chạy async để không block
func (m *Mapper) PreloadAll() int {
	count := 0
	for cp, entry := range m.db.entries {
		addr := m.entryToISL(entry)
		m.mu.Lock()
		m.cache[cp] = addr
		m.reverse[addr.String()] = cp
		m.mu.Unlock()
		count++
	}
	return count
}
