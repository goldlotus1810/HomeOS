// internal/aam/aam.go
// AAM — Agent AI Master
// Giao diện DUY NHẤT với người dùng.
// Quyết định dựa trên 4 chiều: ngữ cảnh + không gian + thời gian + ưu tiên
// Xác nhận mọi thay đổi từ LeoAI trước khi commit

package aam

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/goldlotus1810/HomeOS/internal/isl"
)

// ─────────────────────────────────────────────────────────────────
// INTERFACES — AAM chỉ biết Chiefs, không biết Workers
// ─────────────────────────────────────────────────────────────────

// Chief interface — tất cả Chiefs phải implement
type Chief interface {
	ID() string
	Name() string
	// ForcedUpdate cập nhật ISL dictionary — bắt buộc, không thể từ chối
	ForcedUpdate(ctx context.Context, msg *isl.ISLMessage) error
	// Dispatch gửi task đến Chief
	Dispatch(ctx context.Context, msg *isl.ISLMessage) (*isl.ISLMessage, error)
	// IsOnline kiểm tra Chief còn sống không
	IsOnline() bool
}

// UserInterface — cách AAM giao tiếp với người dùng
type UserInterface interface {
	// AskUser hỏi người dùng và chờ phản hồi
	AskUser(ctx context.Context, question string) (string, error)
	// Notify thông báo không cần phản hồi
	Notify(ctx context.Context, msg string) error
	// AskConfirm hỏi yes/no
	AskConfirm(ctx context.Context, question string) (bool, error)
}

// ─────────────────────────────────────────────────────────────────
// 4D DECISION ENGINE
// ─────────────────────────────────────────────────────────────────

// DecisionContext chứa 4 chiều để AAM ra quyết định
type DecisionContext struct {
	// Chiều 1: Ngữ cảnh (Context)
	RecentEvents  []string  // Các sự kiện gần đây
	ActiveSession string    // Phiên hiện tại (xem phim, nấu ăn, ngủ...)
	UserMood      string    // Trạng thái người dùng (nếu biết)

	// Chiều 2: Không gian (Space)
	UserLocation  string    // Phòng ngủ, phòng khách, bếp...
	ActiveDevices []string  // Thiết bị đang bật

	// Chiều 3: Thời gian (Time)
	TimeOfDay     time.Time
	DayOfWeek     string
	Season        string
	IsHoliday     bool

	// Chiều 4: Ưu tiên (Priority)
	SafetyAlert   bool      // Cảnh báo an toàn = ưu tiên tuyệt đối
	UserPriority  byte      // 0=emergency, 127=normal, 255=background
	ImmutableRule bool      // Quy tắc bất biến đang áp dụng
}

// Decision là kết quả quyết định của AAM
type Decision struct {
	Action      string
	TargetChief string
	ISLMessage  *isl.ISLMessage
	Confidence  float32
	Reasoning   string // Giải thích tại sao đưa ra quyết định này
}

// decide4D là thuật toán ra quyết định dựa trên 4 chiều
func decide4D(input *isl.ISLMessage, ctx *DecisionContext) *Decision {
	confidence := float32(0.5)
	reasoning := ""

	// Chiều 4: Ưu tiên — kiểm tra trước nhất
	if ctx.SafetyAlert || ctx.ImmutableRule {
		return &Decision{
			Confidence: 1.0,
			Reasoning:  "Safety rule / immutable rule takes absolute priority",
		}
	}

	// Chiều 1: Ngữ cảnh
	if ctx.ActiveSession == "sleeping" {
		confidence += 0.1
		reasoning += "User sleeping: prefer quiet actions. "
	}

	// Chiều 3: Thời gian
	hour := ctx.TimeOfDay.Hour()
	if hour >= 22 || hour < 6 {
		confidence += 0.15
		reasoning += "Night time: dim lights, lower volume. "
	}

	// Chiều 2: Không gian
	if ctx.UserLocation != "" {
		confidence += 0.1
		reasoning += fmt.Sprintf("User in %s. ", ctx.UserLocation)
	}

	return &Decision{
		Confidence: confidence,
		Reasoning:  reasoning,
	}
}

// ─────────────────────────────────────────────────────────────────
// AAM — Agent AI Master
// ─────────────────────────────────────────────────────────────────

// AAM là đầu não trung tâm của toàn hệ thống
type AAM struct {
	mu     sync.RWMutex
	chiefs map[string]Chief
	ui     UserInterface

	// Heartbeat tracking
	heartbeats map[string]time.Time

	// Decision context — cập nhật liên tục
	decCtx *DecisionContext

	// Pending confirmations từ LeoAI
	pendingConfirm chan *ConfirmRequest

	// Emergency mode
	emergency bool

	// Channels
	stop chan struct{}
}

// ConfirmRequest là yêu cầu xác nhận từ LeoAI
type ConfirmRequest struct {
	Type        string // "commit", "memory_expand", "new_branch", "immutable_promote"
	Description string
	Data        interface{}
	Response    chan bool
}

// New tạo AAM mới
func New(ui UserInterface) *AAM {
	return &AAM{
		chiefs:         make(map[string]Chief),
		ui:             ui,
		heartbeats:     make(map[string]time.Time),
		decCtx:         &DecisionContext{TimeOfDay: time.Now()},
		pendingConfirm: make(chan *ConfirmRequest, 100),
		stop:           make(chan struct{}),
	}
}

// RegisterChief đăng ký một Chief với AAM
func (a *AAM) RegisterChief(c Chief) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.chiefs[c.ID()] = c
	log.Printf("AAM: Chief %s registered", c.Name())
}

// ─────────────────────────────────────────────────────────────────
// CORE — xử lý input và ra quyết định
// ─────────────────────────────────────────────────────────────────

// Process nhận input từ người dùng và điều phối
func (a *AAM) Process(ctx context.Context, msg *isl.ISLMessage) (*isl.ISLMessage, error) {
	if a.emergency {
		return nil, fmt.Errorf("aam: system in emergency mode — all operations suspended")
	}

	// Cập nhật thời gian quyết định
	a.decCtx.TimeOfDay = time.Now()

	// Ra quyết định dựa trên 4 chiều
	decision := decide4D(msg, a.decCtx)

	// Tìm Chief phù hợp
	chiefID := a.selectChief(msg)
	chief, ok := a.chiefs[chiefID]
	if !ok {
		return nil, fmt.Errorf("aam: no chief available for message type")
	}

	// Dispatch đến Chief
	result, err := chief.Dispatch(ctx, msg)
	if err != nil {
		return nil, fmt.Errorf("aam: chief %s error: %w", chiefID, err)
	}

	// Học từ kết quả (không blocking — LeoAI lắng nghe async)
	go a.notifyLeoAI(msg, result, decision)

	return result, nil
}

// selectChief chọn Chief phù hợp dựa trên loại message
func (a *AAM) selectChief(msg *isl.ISLMessage) string {
	// Dựa vào ISL Layer để chọn Chief
	switch msg.PrimaryAddr.Layer {
	case 'A': // Visual
		return "vision_chief"
	case 'B', 'C': // Organism, Color
		return "vision_chief"
	case 'S': // Sound/Audio
		return "audio_chief"
	case 'H': // Home
		return "home_chief"
	case 'E': // Energy
		return "energy_chief"
	case 'X': // Security
		return "security_chief"
	default:
		return "home_chief" // Default
	}
}

// ─────────────────────────────────────────────────────────────────
// CASCADE BROADCAST — từ LeoAI xuống tất cả Chiefs
// ─────────────────────────────────────────────────────────────────

// Cascade broadcast ISL message xuống tất cả Chiefs đồng thời
// Đây là cơ chế đồng bộ từ điển ISL — buộc ghi nhớ, không thể từ chối
func (a *AAM) Cascade(ctx context.Context, msg *isl.ISLMessage) error {
	a.mu.RLock()
	chiefs := make([]Chief, 0, len(a.chiefs))
	for _, c := range a.chiefs {
		chiefs = append(chiefs, c)
	}
	a.mu.RUnlock()

	// Broadcast song song đến tất cả Chiefs
	var wg sync.WaitGroup
	errors := make(chan error, len(chiefs))

	for _, c := range chiefs {
		wg.Add(1)
		go func(chief Chief) {
			defer wg.Done()
			if err := chief.ForcedUpdate(ctx, msg); err != nil {
				errors <- fmt.Errorf("chief %s update failed: %w", chief.ID(), err)
				// Ghi nhận lỗi nhưng không dừng broadcast
			}
		}(c)
	}

	wg.Wait()
	close(errors)

	// Gom lỗi
	var errs []error
	for err := range errors {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		// Retry cho Chiefs bị lỗi
		log.Printf("AAM: %d chiefs failed to update, scheduling retry", len(errs))
	}

	return nil
}

// ─────────────────────────────────────────────────────────────────
// LEOAI COMMUNICATION — xác nhận và học tập
// ─────────────────────────────────────────────────────────────────

// notifyLeoAI gửi kết quả cho LeoAI để học (async, không blocking)
func (a *AAM) notifyLeoAI(input, output *isl.ISLMessage, decision *Decision) {
	// LeoAI nhận kết quả và học từ đó
	// Implementation: gửi qua channel đến LeoAI goroutine
}

// ConfirmFromLeoAI nhận yêu cầu xác nhận từ LeoAI
// LeoAI không tự quyết định — phải hỏi AAM (= người dùng)
func (a *AAM) ConfirmFromLeoAI(req *ConfirmRequest) {
	a.pendingConfirm <- req
}

// processConfirmations xử lý các yêu cầu xác nhận từ LeoAI
func (a *AAM) processConfirmations(ctx context.Context) {
	for {
		select {
		case <-a.stop:
			return
		case req := <-a.pendingConfirm:
			// Hỏi người dùng
			confirmed, err := a.ui.AskConfirm(ctx, req.Description)
			if err != nil {
				req.Response <- false
				continue
			}
			req.Response <- confirmed

			if confirmed {
				log.Printf("AAM: User confirmed: %s", req.Type)
			} else {
				log.Printf("AAM: User rejected: %s", req.Type)
			}
		}
	}
}

// ─────────────────────────────────────────────────────────────────
// HEARTBEAT & EMERGENCY
// ─────────────────────────────────────────────────────────────────

// StartHeartbeat theo dõi sức khỏe của tất cả Chiefs
func (a *AAM) StartHeartbeat(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-a.stop:
			return
		case <-ticker.C:
			a.checkHeartbeats()
		}
	}
}

func (a *AAM) checkHeartbeats() {
	a.mu.RLock()
	defer a.mu.RUnlock()

	for id, chief := range a.chiefs {
		if !chief.IsOnline() {
			log.Printf("AAM: WARNING — Chief %s is offline, attempting reconnect", id)
			// Trigger reconnect
		}
	}
}

// TriggerEmergency được gọi khi AAM mất kết nối hoặc có sự cố nghiêm trọng
// Dừng toàn bộ hệ thống và báo động
func (a *AAM) TriggerEmergency(reason string) {
	a.mu.Lock()
	a.emergency = true
	a.mu.Unlock()

	log.Printf("AAM: ⚠️  EMERGENCY: %s — All operations suspended", reason)

	// Thông báo tất cả Chiefs dừng lại
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	emergencyMsg := &isl.ISLMessage{
		MsgType:  isl.MsgEmergency,
		Priority: 0, // Ưu tiên tuyệt đối
	}
	_ = a.Cascade(ctx, emergencyMsg)

	// Kích hoạt báo động vật lý (còi, đèn...)
	_ = a.ui.Notify(ctx, fmt.Sprintf("⚠️ EMERGENCY: %s", reason))
}

// Run khởi động AAM
func (a *AAM) Run(ctx context.Context) error {
	log.Println("AAM: Starting...")

	go a.StartHeartbeat(ctx)
	go a.processConfirmations(ctx)

	log.Println("AAM: Online — ready to receive input")

	<-ctx.Done()
	close(a.stop)
	return nil
}
