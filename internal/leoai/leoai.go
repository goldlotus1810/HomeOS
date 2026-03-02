// internal/leoai/leoai.go
// LeoAI — Chief of Intelligence
// Không ra quyết định thời gian thực — đó là việc của AAM.
// LeoAI là bộ não học tập: quản lý cây tri thức, tạo ISL, clustering, ký số bất biến.

package leoai

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"log"
	"sync"
	"time"

	"github.com/goldlotus1810/HomeOS/internal/isl"
	"github.com/goldlotus1810/HomeOS/internal/tree"
)

// ─────────────────────────────────────────────────────────────────
// SHORT-TERM MEMORY — học trước khi commit vào cây
// ─────────────────────────────────────────────────────────────────

// ShortTermMemory là bộ nhớ tạm thời để LeoAI học
// Không ảnh hưởng đến Agent nào đang chạy cho đến khi commit
type ShortTermMemory struct {
	mu       sync.Mutex
	entries  []*MemoryEntry
	maxBytes int64
	usedBytes int64
}

// MemoryEntry là một đơn vị học tập
type MemoryEntry struct {
	ISLAddr    isl.Address
	Payload    []byte
	Immutable  bool
	Confidence float32
	Source     string    // "user_input", "observation", "feedback"
	CreatedAt  time.Time
	Confirmed  bool      // Đã được AAM xác nhận chưa
}

func NewShortTermMemory(maxMB int64) *ShortTermMemory {
	return &ShortTermMemory{
		maxBytes: maxMB * 1024 * 1024,
	}
}

// Add thêm entry vào bộ nhớ ngắn hạn
func (m *ShortTermMemory) Add(entry *MemoryEntry) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	entrySize := int64(len(entry.Payload) + 32)
	if m.usedBytes+entrySize > m.maxBytes {
		return false // Đầy — cần xin phép AAM
	}

	m.entries = append(m.entries, entry)
	m.usedBytes += entrySize
	return true
}

// IsFull kiểm tra bộ nhớ có sắp đầy không (>= 85%)
func (m *ShortTermMemory) IsFull() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return float64(m.usedBytes)/float64(m.maxBytes) >= 0.85
}

// UsagePercent trả về phần trăm sử dụng
func (m *ShortTermMemory) UsagePercent() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	return float64(m.usedBytes) / float64(m.maxBytes) * 100
}

// Flush xóa các entry đã được commit
func (m *ShortTermMemory) Flush(committed []*MemoryEntry) {
	m.mu.Lock()
	defer m.mu.Unlock()
	committedMap := make(map[*MemoryEntry]bool)
	for _, e := range committed { committedMap[e] = true }

	newEntries := m.entries[:0]
	for _, e := range m.entries {
		if !committedMap[e] {
			newEntries = append(newEntries, e)
		} else {
			m.usedBytes -= int64(len(e.Payload) + 32)
		}
	}
	m.entries = newEntries
}

// ─────────────────────────────────────────────────────────────────
// CLUSTERING — tự phân nhóm không cần người dùng
// ─────────────────────────────────────────────────────────────────

// Cluster đại diện cho một nhóm đang hình thành
type Cluster struct {
	ISLAddr  isl.Address
	Members  []*MemoryEntry
	Centroid [3]float32 // Vector trung tâm
	Score    float32    // Độ gắn kết
}

// cluster phân nhóm các entry trong bộ nhớ ngắn hạn
// Đây là unsupervised clustering — không cần người dùng đánh nhãn
func (l *LeoAI) cluster(entries []*MemoryEntry) []*Cluster {
	if len(entries) == 0 {
		return nil
	}

	// Nhóm theo ISL Layer + Group (cùng ngữ nghĩa domain)
	groups := make(map[string][]*MemoryEntry)
	for _, e := range entries {
		key := string([]byte{e.ISLAddr.Layer, e.ISLAddr.Group})
		groups[key] = append(groups[key], e)
	}

	var clusters []*Cluster
	for _, members := range groups {
		if len(members) == 0 { continue }

		// Tính centroid (trung tâm ngữ nghĩa)
		var cx, cy, cz float32
		for _, m := range members {
			// Đơn giản hoá: dùng ISL byte values làm coordinates
			cx += float32(m.ISLAddr.Layer)
			cy += float32(m.ISLAddr.Group)
			cz += float32(m.ISLAddr.Type)
		}
		n := float32(len(members))
		cx, cy, cz = cx/n, cy/n, cz/n

		clusters = append(clusters, &Cluster{
			ISLAddr:  members[0].ISLAddr,
			Members:  members,
			Centroid: [3]float32{cx, cy, cz},
			Score:    n / float32(tree.BranchThreshold),
		})
	}

	return clusters
}

// ─────────────────────────────────────────────────────────────────
// LEOAI — Chief of Intelligence
// ─────────────────────────────────────────────────────────────────

// LeoAI quản lý tri thức và học tập
type LeoAI struct {
	mu  sync.RWMutex

	// Signing key cho dữ liệu bất biến
	pubKey ed25519.PublicKey
	privKey ed25519.PrivateKey

	// Cây tri thức dài hạn
	knowledgeTree *tree.Tree

	// Bộ nhớ ngắn hạn
	shortTerm *ShortTermMemory

	// ISL dictionary
	dict *isl.Dictionary

	// Giao tiếp với AAM
	confirmCh chan *ConfirmRequest // Xin phép AAM
	learnCh   chan *LearningEvent // Nhận kết quả từ AAM

	stop chan struct{}
}

// ConfirmRequest — LeoAI xin phép AAM trước khi làm gì đó
type ConfirmRequest struct {
	Type        string
	Description string
	Data        interface{}
	Response    chan bool
}

// LearningEvent — kết quả từ AAM để LeoAI học
type LearningEvent struct {
	Input    *isl.ISLMessage
	Output   *isl.ISLMessage
	Feedback string // "good", "bad", "neutral", ""(silent)
}

// AgentLesson — bài học gửi xuống cho Agent
type AgentLesson struct {
	DictVersion uint32
	Updates     []isl.Address
	NewEntries  map[string]isl.Address
}

// New tạo LeoAI mới
func New(treePath string, memoryMB int64) (*LeoAI, error) {
	// Tạo signing key pair cho bất biến
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	t := tree.NewTree(priv)

	l := &LeoAI{
		pubKey:        pub,
		privKey:       priv,
		knowledgeTree: t,
		shortTerm:     NewShortTermMemory(memoryMB),
		dict:          isl.NewDictionary(),
		confirmCh:     make(chan *ConfirmRequest, 50),
		learnCh:       make(chan *LearningEvent, 1000),
		stop:          make(chan struct{}),
	}

	return l, nil
}

// ─────────────────────────────────────────────────────────────────
// LEARNING LOOP — vòng lặp học tập chính
// ─────────────────────────────────────────────────────────────────

// Run khởi động LeoAI learning loop
func (l *LeoAI) Run(ctx context.Context) {
	log.Println("LeoAI: Starting learning loop...")

	// Goroutine học từ events
	go l.learningLoop(ctx)

	// Goroutine clustering định kỳ
	go l.clusteringLoop(ctx)

	// Goroutine kiểm tra bộ nhớ
	go l.memoryWatchdog(ctx)

	<-ctx.Done()
	close(l.stop)
}

// learningLoop liên tục học từ phản hồi AAM
func (l *LeoAI) learningLoop(ctx context.Context) {
	for {
		select {
		case <-l.stop:
			return
		case event := <-l.learnCh:
			l.processLearningEvent(event)
		}
	}
}

// processLearningEvent xử lý một sự kiện học tập
func (l *LeoAI) processLearningEvent(event *LearningEvent) {
	entry := &MemoryEntry{
		ISLAddr:   event.Input.PrimaryAddr,
		Confidence: float32(event.Output.Confidence) / 100.0,
		Source:    "aam_feedback",
		CreatedAt: time.Now(),
	}

	// Điều chỉnh confidence dựa trên feedback
	switch event.Feedback {
	case "good":
		entry.Confidence = min(entry.Confidence*1.1, 1.0)
	case "bad":
		entry.Confidence = entry.Confidence * 0.8
	default:
		// Im lặng → giữ nguyên, tích lũy thêm
	}

	if !l.shortTerm.Add(entry) {
		// Bộ nhớ đầy — cần xin phép AAM
		l.requestMemoryExpansion()
	}
}

// clusteringLoop định kỳ phân cụm dữ liệu trong short-term memory
func (l *LeoAI) clusteringLoop(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-l.stop:
			return
		case <-ticker.C:
			l.runClustering(ctx)
		}
	}
}

// runClustering chạy thuật toán clustering và commit kết quả
func (l *LeoAI) runClustering(ctx context.Context) {
	l.mu.RLock()
	entries := l.shortTerm.entries
	l.mu.RUnlock()

	if len(entries) == 0 {
		return
	}

	clusters := l.cluster(entries)

	for _, c := range clusters {
		if len(c.Members) >= tree.BranchThreshold {
			// Đủ ngưỡng → xin phép AAM để commit vào cây
			l.requestCommit(ctx, c)
		}
	}
}

// requestCommit xin phép AAM để commit một cluster vào cây dài hạn
func (l *LeoAI) requestCommit(ctx context.Context, c *Cluster) {
	req := &ConfirmRequest{
		Type:        "commit",
		Description: l.buildCommitDescription(c),
		Data:        c,
		Response:    make(chan bool, 1),
	}

	l.confirmCh <- req

	// Chờ phản hồi từ AAM (= người dùng)
	select {
	case confirmed := <-req.Response:
		if confirmed {
			l.commitCluster(c)
		}
	case <-time.After(24 * time.Hour):
		// Người dùng không phản hồi trong 24h → giữ lại, thử lại sau
		log.Printf("LeoAI: Commit request timed out, will retry later")
	}
}

// buildCommitDescription tạo mô tả cho request xác nhận
func (l *LeoAI) buildCommitDescription(c *Cluster) string {
	return "LeoAI đã tích lũy đủ dữ liệu để tạo nhóm tri thức mới. Xác nhận?"
}

// commitCluster commit cluster vào cây tri thức dài hạn
func (l *LeoAI) commitCluster(c *Cluster) {
	for _, member := range c.Members {
		dp := tree.DataPoint{
			ISLAddr:   member.ISLAddr,
			Payload:   member.Payload,
			Immutable: member.Immutable,
		}
		if err := l.knowledgeTree.Ingest(dp); err != nil {
			log.Printf("LeoAI: commit error: %v", err)
			continue
		}
	}

	// Flush committed entries từ short-term
	l.shortTerm.Flush(c.Members)
	log.Printf("LeoAI: Committed cluster with %d members to knowledge tree", len(c.Members))
}

// ─────────────────────────────────────────────────────────────────
// MEMORY WATCHDOG
// ─────────────────────────────────────────────────────────────────

// memoryWatchdog theo dõi bộ nhớ ngắn hạn
func (l *LeoAI) memoryWatchdog(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-l.stop:
			return
		case <-ticker.C:
			if l.shortTerm.IsFull() {
				l.requestMemoryExpansion()
			}
		}
	}
}

// requestMemoryExpansion xin phép AAM cấp thêm bộ nhớ ngắn hạn
// LeoAI KHÔNG tự quyết định mở rộng bộ nhớ
func (l *LeoAI) requestMemoryExpansion() {
	usage := l.shortTerm.UsagePercent()
	req := &ConfirmRequest{
		Type:        "memory_expand",
		Description: "Bộ nhớ học tập sắp đầy. Cho phép LeoAI mở rộng thêm?",
		Data:        map[string]interface{}{"usage_percent": usage},
		Response:    make(chan bool, 1),
	}
	l.confirmCh <- req
	// Không chờ — tiếp tục chạy
	// Khi được xác nhận thì expand, không được thì tiếp tục với bộ nhớ hiện tại
}

// ─────────────────────────────────────────────────────────────────
// IMMUTABLE DATA
// ─────────────────────────────────────────────────────────────────

// IngestImmutable nạp dữ liệu bất biến từ người dùng (qua AAM)
// Dữ liệu bất biến được ưu tiên tuyệt đối, broadcast ngay lập tức
func (l *LeoAI) IngestImmutable(addr isl.Address, payload []byte) error {
	dp := tree.DataPoint{
		ISLAddr:   addr,
		Payload:   payload,
		Immutable: true,
	}
	return l.knowledgeTree.Ingest(dp)
}

// ─────────────────────────────────────────────────────────────────
// QUERY
// ─────────────────────────────────────────────────────────────────

// Query tìm kiếm trong cây tri thức
func (l *LeoAI) Query(addr isl.Address, depth int) []tree.QueryResult {
	return l.knowledgeTree.Query(addr, depth, 50)
}

// GetTreeFor3D xuất dữ liệu cây cho web visualization
func (l *LeoAI) GetTreeFor3D() []byte {
	return l.knowledgeTree.ExportFor3D()
}

// ─────────────────────────────────────────────────────────────────
// HELPERS
// ─────────────────────────────────────────────────────────────────

func min(a, b float32) float32 {
	if a < b { return a }
	return b
}
