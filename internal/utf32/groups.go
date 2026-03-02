// internal/utf32/groups.go
// Phân nhóm ký tự theo ngôn ngữ, loại bỏ trùng lặp
// En group = ASCII cơ bản
// Vn group = chỉ chứa ký tự THÊM VÀO so với ASCII
// ZH group = ký tự tượng hình CJK
// → Khi cần tiếng Việt: gọi En + Vn → dung lượng giảm

package utf32

import "sort"

// GroupID định danh nhóm ngôn ngữ
type GroupID byte

const (
	GroupEn  GroupID = 'E' // English / ASCII
	GroupVn  GroupID = 'V' // Vietnamese (chỉ phần thêm vào)
	GroupZH  GroupID = 'C' // Chinese / CJK
	GroupJA  GroupID = 'J' // Japanese
	GroupKO  GroupID = 'K' // Korean
	GroupAR  GroupID = 'A' // Arabic
	GroupHE  GroupID = 'H' // Hebrew
	GroupGR  GroupID = 'G' // Greek
	GroupCY  GroupID = 'Y' // Cyrillic
	GroupEM  GroupID = 'M' // Emoji / Symbols
	GroupOTH GroupID = 'O' // Other
)

// Group là một nhóm ký tự ngôn ngữ
type Group struct {
	ID         GroupID
	Name       string
	CodePoints []rune            // Danh sách codepoint UNIQUE (không trùng với En)
	Shared     []rune            // Codepoint dùng chung với En (có mã En+GroupID)
}

// Groups chứa tất cả nhóm ngôn ngữ
type Groups struct {
	db     *DB
	groups map[GroupID]*Group
}

// NewGroups tạo và phân nhóm từ DB
func NewGroups(db *DB) *Groups {
	g := &Groups{
		db:     db,
		groups: make(map[GroupID]*Group),
	}
	g.build()
	return g
}

// build phân nhóm toàn bộ ký tự Unicode
func (g *Groups) build() {
	// Khởi tạo các nhóm
	groupDefs := []struct {
		id   GroupID
		name string
	}{
		{GroupEn, "English/ASCII"},
		{GroupVn, "Vietnamese"},
		{GroupZH, "Chinese/CJK"},
		{GroupJA, "Japanese"},
		{GroupKO, "Korean"},
		{GroupAR, "Arabic"},
		{GroupHE, "Hebrew"},
		{GroupGR, "Greek"},
		{GroupCY, "Cyrillic"},
		{GroupEM, "Emoji/Symbols"},
		{GroupOTH, "Other"},
	}

	for _, def := range groupDefs {
		g.groups[def.id] = &Group{
			ID:   def.id,
			Name: def.name,
		}
	}

	// ASCII set — nền tảng của GroupEn
	asciiSet := make(map[rune]bool)
	for cp := rune(0x0020); cp <= rune(0x007E); cp++ {
		asciiSet[cp] = true
		g.groups[GroupEn].CodePoints = append(g.groups[GroupEn].CodePoints, cp)
	}

	// Phân loại toàn bộ ký tự
	for cp, entry := range g.db.entries {
		gid := blockToGroup(entry.Block)
		grp := g.groups[gid]
		if grp == nil {
			grp = g.groups[GroupOTH]
		}

		if gid == GroupEn {
			continue // Đã xử lý ASCII ở trên
		}

		// Kiểm tra có trùng với ASCII không
		if asciiSet[cp] {
			// Trùng → đánh dấu shared, không thêm vào unique
			grp.Shared = append(grp.Shared, cp)
		} else {
			// Không trùng → thêm vào unique của nhóm này
			grp.CodePoints = append(grp.CodePoints, cp)
		}
	}

	// Sắp xếp để nhất quán
	for _, grp := range g.groups {
		sort.Slice(grp.CodePoints, func(i, j int) bool {
			return grp.CodePoints[i] < grp.CodePoints[j]
		})
	}
}

// blockToGroup ánh xạ Unicode block → GroupID
func blockToGroup(block string) GroupID {
	switch block {
	case "Basic Latin":
		return GroupEn
	case "Latin-1 Supplement", "Latin Extended":
		return GroupEn // Vẫn là Latin
	case "Vietnamese":
		return GroupVn
	case "CJK Unified":
		return GroupZH
	case "Japanese":
		return GroupJA
	case "Korean":
		return GroupKO
	case "Arabic":
		return GroupAR
	case "Hebrew":
		return GroupHE
	case "Greek":
		return GroupGR
	case "Cyrillic":
		return GroupCY
	case "Emoji", "Mahjong":
		return GroupEM
	default:
		return GroupOTH
	}
}

// Get lấy nhóm theo ID
func (g *Groups) Get(id GroupID) *Group {
	return g.groups[id]
}

// EncodeLang trả về danh sách codepoint cần thiết cho một ngôn ngữ
// Ví dụ: EncodeLang("vi") = En + Vn (không trùng lặp)
// → Tiết kiệm dung lượng so với lưu toàn bộ riêng biệt
func (g *Groups) EncodeLang(lang string) []rune {
	switch lang {
	case "en":
		return g.groups[GroupEn].CodePoints

	case "vi":
		// En + phần unique của Vn (không có ASCII trùng)
		result := make([]rune, len(g.groups[GroupEn].CodePoints))
		copy(result, g.groups[GroupEn].CodePoints)
		result = append(result, g.groups[GroupVn].CodePoints...)
		return result

	case "zh":
		result := make([]rune, len(g.groups[GroupEn].CodePoints))
		copy(result, g.groups[GroupEn].CodePoints)
		result = append(result, g.groups[GroupZH].CodePoints...)
		return result

	case "ja":
		result := make([]rune, len(g.groups[GroupEn].CodePoints))
		copy(result, g.groups[GroupEn].CodePoints)
		result = append(result, g.groups[GroupJA].CodePoints...)
		result = append(result, g.groups[GroupZH].CodePoints...) // Kanji dùng CJK
		return result

	default:
		return g.groups[GroupEn].CodePoints
	}
}

// Stats thống kê số ký tự mỗi nhóm
func (g *Groups) Stats() map[string]int {
	stats := make(map[string]int)
	for id, grp := range g.groups {
		stats[string([]byte{byte(id)})] = len(grp.CodePoints)
	}
	return stats
}
