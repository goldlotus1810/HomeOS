// scripts/seed_utf32.go
// Chạy một lần để tải toàn bộ Unicode về máy
// go run scripts/seed_utf32.go

package main

import (
	"log"

	"github.com/goldlotus1810/HomeOS/internal/utf32"
)

func main() {
	log.Println("🌱 Seeding UTF-32 database...")

	db := utf32.NewDB()
	if err := db.Load(); err != nil {
		log.Fatalf("❌ Failed: %v", err)
	}

	stats := db.Stats()
	log.Printf("✅ Total characters: %d", stats["total"])
	log.Printf("✅ Total blocks: %d", stats["blocks"])

	// Phân nhóm ngôn ngữ
	groups := utf32.NewGroups(db)
	gstats := groups.Stats()
	log.Printf("📊 Group stats:")
	for g, count := range gstats {
		log.Printf("   %s: %d chars", g, count)
	}

	// Test mapper
	mapper := utf32.NewMapper(db)
	log.Println("🔄 Pre-loading ISL address map...")
	count := mapper.PreloadAll()
	log.Printf("✅ Mapped %d codepoints → ISL addresses", count)

	// Test lookup
	testConcepts := []string{"BANANA", "HOUSE", "CAT", "SUN"}
	for _, c := range testConcepts {
		addr, err := mapper.MapConcept(c)
		if err != nil {
			log.Printf("   %-10s → not found", c)
			continue
		}
		log.Printf("   %-10s → ISL: %s", c, addr.String())
	}

	log.Println("✅ UTF-32 seed complete — Thân cây bất biến sẵn sàng!")
}
