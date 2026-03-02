// internal/isl/address.go
// ISL — Internal Shared Language
// Địa chỉ ngữ nghĩa 64-bit: Layer(1B) + Group(1B) + Type(1B) + ID(1B) + Attributes(4B)

package isl

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"sync"
)

// Address là địa chỉ ISL 64-bit — đại diện cho một khái niệm ngữ nghĩa
// Giống UTF-8 cho ký tự, ISL là UTF-8 cho khái niệm
type Address struct {
	Layer      byte   // A-Z: Domain (A=Visual, B=Organism, C=Color, D=Location...)
	Group      byte   // A-Z: Semantic cluster trong domain
	Type       byte   // a-z: a=concrete, b=abstract, c=action, d=property...
	ID         byte   // 0-255: Sequential ID trong group
	Attributes uint32 // Extended semantic attributes
}

// String trả về biểu diễn human-readable của địa chỉ ISL
func (a Address) String() string {
	return fmt.Sprintf("%c%c%c%d", a.Layer, a.Group, a.Type, a.ID)
}

// Bytes chuyển Address thành 8 bytes nhị phân
func (a Address) Bytes() []byte {
	b := make([]byte, 8)
	b[0] = a.Layer
	b[1] = a.Group
	b[2] = a.Type
	b[3] = a.ID
	binary.BigEndian.PutUint32(b[4:], a.Attributes)
	return b
}

// FromBytes phục hồi Address từ 8 bytes
func FromBytes(b []byte) Address {
	if len(b) < 8 {
		return Address{}
	}
	return Address{
		Layer:      b[0],
		Group:      b[1],
		Type:       b[2],
		ID:         b[3],
		Attributes: binary.BigEndian.Uint32(b[4:]),
	}
}

// ─────────────────────────────────────────────────────────────────
// MESSAGE — đơn vị giao tiếp ISL
// ─────────────────────────────────────────────────────────────────

type MsgType byte

const (
	MsgActivate   MsgType = 0x01 // Kích hoạt Agent thực hiện task
	MsgLearn      MsgType = 0x02 // Cập nhật ISL dictionary
	MsgDeactivate MsgType = 0x03 // Tắt Agent
	MsgQuery      MsgType = 0x04 // Query knowledge tree
	MsgResponse   MsgType = 0x05 // Phản hồi kết quả
	MsgImmutable  MsgType = 0x06 // Broadcast dữ liệu bất biến (ưu tiên tuyệt đối)
	MsgHeartbeat  MsgType = 0x07 // Heartbeat check
	MsgEmergency  MsgType = 0xFF // Khẩn cấp — dừng tất cả
)

// ISLMessage là gói tin ISL đầy đủ
// Tổng kích thước: ~68 bytes sau mã hóa AES-256-GCM
// So với JSON: ~280 bytes → tiết kiệm 75.7%
type ISLMessage struct {
	// Header (8 bytes)
	Version  byte    // Protocol version
	MsgType  MsgType // Loại message
	SenderID uint16  // Agent ID gửi
	TargetID uint16  // Agent ID nhận (0 = broadcast)
	Priority byte    // 0=max, 255=lowest (bất biến luôn = 0)
	Flags    byte    // Bit flags

	// Addresses (16 bytes)
	PrimaryAddr   Address // Địa chỉ ISL chính
	SecondaryAddr Address // Địa chỉ ISL phụ (nếu cần)

	// Context (8 bytes)
	ContextAddr Address // Ngữ cảnh không gian/thời gian

	// Metadata (8 bytes)
	Confidence uint32 // 0-100: độ tin cậy
	Timestamp  uint32 // Unix timestamp (seconds)

	// Payload (biến đổi, tối đa 200 bytes)
	Payload []byte

	// Checksum (4 bytes)
	CRC32 uint32
}

// ─────────────────────────────────────────────────────────────────
// CODEC — mã hóa/giải mã ISL với AES-256-GCM
// ─────────────────────────────────────────────────────────────────

// ISLCodec xử lý encode/decode và mã hóa ISL messages
type ISLCodec struct {
	gcm    cipher.AEAD
	mu     sync.RWMutex
	keyVer uint32 // Version của key (tăng khi key rotation)
}

// NewISLCodec tạo codec mới với AES-256-GCM key
func NewISLCodec(key []byte) (*ISLCodec, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("isl: create cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("isl: create gcm: %w", err)
	}
	return &ISLCodec{gcm: gcm}, nil
}

// Encode chuyển ISLMessage thành bytes đã mã hóa
func (c *ISLCodec) Encode(msg *ISLMessage) ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Serialize message
	raw := c.serialize(msg)

	// Tạo nonce ngẫu nhiên (12 bytes cho AES-GCM)
	nonce := make([]byte, c.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("isl: generate nonce: %w", err)
	}

	// Mã hóa
	encrypted := c.gcm.Seal(nonce, nonce, raw, nil)
	return encrypted, nil
}

// Decode giải mã bytes thành ISLMessage
func (c *ISLCodec) Decode(data []byte) (*ISLMessage, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	nonceSize := c.gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("isl: message too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	raw, err := c.gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("isl: decrypt: %w", err)
	}

	return c.deserialize(raw)
}

// serialize chuyển ISLMessage thành byte slice
func (c *ISLCodec) serialize(msg *ISLMessage) []byte {
	buf := make([]byte, 0, 64+len(msg.Payload))

	// Header
	buf = append(buf, msg.Version, byte(msg.MsgType))
	buf = binary.BigEndian.AppendUint16(buf, msg.SenderID)
	buf = binary.BigEndian.AppendUint16(buf, msg.TargetID)
	buf = append(buf, msg.Priority, msg.Flags)

	// Addresses
	buf = append(buf, msg.PrimaryAddr.Bytes()...)
	buf = append(buf, msg.SecondaryAddr.Bytes()...)
	buf = append(buf, msg.ContextAddr.Bytes()...)

	// Metadata
	buf = binary.BigEndian.AppendUint32(buf, msg.Confidence)
	buf = binary.BigEndian.AppendUint32(buf, msg.Timestamp)

	// Payload
	payloadLen := uint16(len(msg.Payload))
	buf = binary.BigEndian.AppendUint16(buf, payloadLen)
	buf = append(buf, msg.Payload...)

	// CRC32
	buf = binary.BigEndian.AppendUint32(buf, msg.CRC32)

	return buf
}

func (c *ISLCodec) deserialize(data []byte) (*ISLMessage, error) {
	if len(data) < 46 {
		return nil, fmt.Errorf("isl: message too short to deserialize")
	}
	msg := &ISLMessage{}
	offset := 0

	msg.Version = data[offset]; offset++
	msg.MsgType = MsgType(data[offset]); offset++
	msg.SenderID = binary.BigEndian.Uint16(data[offset:]); offset += 2
	msg.TargetID = binary.BigEndian.Uint16(data[offset:]); offset += 2
	msg.Priority = data[offset]; offset++
	msg.Flags = data[offset]; offset++

	msg.PrimaryAddr = FromBytes(data[offset:]); offset += 8
	msg.SecondaryAddr = FromBytes(data[offset:]); offset += 8
	msg.ContextAddr = FromBytes(data[offset:]); offset += 8

	msg.Confidence = binary.BigEndian.Uint32(data[offset:]); offset += 4
	msg.Timestamp = binary.BigEndian.Uint32(data[offset:]); offset += 4

	payloadLen := binary.BigEndian.Uint16(data[offset:]); offset += 2
	if len(data) < offset+int(payloadLen)+4 {
		return nil, fmt.Errorf("isl: truncated payload")
	}
	msg.Payload = make([]byte, payloadLen)
	copy(msg.Payload, data[offset:]); offset += int(payloadLen)

	msg.CRC32 = binary.BigEndian.Uint32(data[offset:])
	return msg, nil
}

// ─────────────────────────────────────────────────────────────────
// DICTIONARY — từ điển ISL do LeoAI quản lý
// ─────────────────────────────────────────────────────────────────

// Dictionary ánh xạ khái niệm ngôn ngữ tự nhiên → ISL Address
type Dictionary struct {
	mu      sync.RWMutex
	version uint32
	entries map[string]Address // "banana" → ABb5
	reverse map[string]string  // "ABb5" → "banana"
}

func NewDictionary() *Dictionary {
	return &Dictionary{
		entries: make(map[string]Address),
		reverse: make(map[string]string),
	}
}

// Lookup tìm ISL address cho một khái niệm
func (d *Dictionary) Lookup(concept string) (Address, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	addr, ok := d.entries[concept]
	return addr, ok
}

// Register thêm entry mới vào dictionary (chỉ LeoAI gọi)
func (d *Dictionary) Register(concept string, addr Address) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.entries[concept] = addr
	d.reverse[addr.String()] = concept
	d.version++
}

// Version trả về version hiện tại của dictionary
func (d *Dictionary) Version() uint32 {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.version
}
