package agents

import (
	"context"
	"encoding/binary"
	"log"
	"math"
	"sync"
	"time"

	"github.com/goldlotus1810/HomeOS/internal/isl"
	"github.com/goldlotus1810/HomeOS/internal/leoai"
)

// BaseAgent interface mà mọi Agent phải implement
type BaseAgent interface {
	ID() uint16
	Name() string
	OnActivate(ctx context.Context, msg *isl.ISLMessage) (*isl.ISLMessage, error)
	OnLearn(ctx context.Context, lesson *leoai.AgentLesson) error
	OnShutdown()
}

// SilentAgent — im lặng mặc định
// Embed vào mọi Agent cụ thể để có hành vi im lặng mặc định
type SilentAgent struct {
	id          uint16
	name        string
	codec       *isl.ISLCodec
	channel     <-chan []byte
	chiefOut    chan<- []byte
	active      bool
	mu          sync.Mutex
	dictVersion uint32
}

// Run là vòng lặp chính — chỉ lắng nghe, không bao giờ tự phát tín hiệu
func (a *SilentAgent) Run(ctx context.Context, impl BaseAgent) {
	log.Printf("Agent[%s]: Online — silent mode", a.name)

	for {
		select {
		case <-ctx.Done():
			impl.OnShutdown()
			return

		case data, ok := <-a.channel:
			if !ok {
				impl.OnShutdown()
				return
			}

			msg, err := a.codec.Decode(data)
			if err != nil {
				continue
			}

			if msg.TargetID != a.id && msg.TargetID != 0 {
				continue
			}

			a.mu.Lock()
			a.active = true
			a.mu.Unlock()

			var result *isl.ISLMessage

			switch msg.MsgType {
			case isl.MsgActivate:
				result, err = impl.OnActivate(ctx, msg)

			case isl.MsgLearn:
				lesson := decodeLessonFromMsg(msg)
				if lesson != nil {
					_ = impl.OnLearn(ctx, lesson)
					a.dictVersion = lesson.DictVersion
				}

			case isl.MsgDeactivate:
				impl.OnShutdown()
				return

			case isl.MsgEmergency:
				log.Printf("Agent[%s]: Emergency stop", a.name)
				impl.OnShutdown()
				return
			}

			a.mu.Lock()
			a.active = false
			a.mu.Unlock()

			if err == nil && result != nil {
				if encoded, encErr := a.codec.Encode(result); encErr == nil {
					select {
					case a.chiefOut <- encoded:
					default:
					}
				}
			}
		}
	}
}

func decodeLessonFromMsg(msg *isl.ISLMessage) *leoai.AgentLesson {
	if len(msg.Payload) == 0 {
		return nil
	}
	return &leoai.AgentLesson{
		DictVersion: msg.Confidence,
	}
}

// ─────────────────────────────────────────────────────────────────
// VISION AGENT
// ─────────────────────────────────────────────────────────────────

type VisionAgent struct {
	SilentAgent
	recognitionThreshold float32
	thresholdHistory     []float32
}

func NewVisionAgent(id uint16, codec *isl.ISLCodec) *VisionAgent {
	return &VisionAgent{
		SilentAgent: SilentAgent{
			id:    id,
			name:  "VisionAgent",
			codec: codec,
		},
		recognitionThreshold: 0.6,
	}
}

func (a *VisionAgent) ID() uint16   { return a.id }
func (a *VisionAgent) Name() string { return a.name }

func (a *VisionAgent) OnActivate(ctx context.Context, msg *isl.ISLMessage) (*isl.ISLMessage, error) {
	if len(msg.Payload) == 0 {
		return nil, nil
	}

	shape := a.analyzeImage(msg.Payload)
	confidence := shape.matchScore
	var islAddr isl.Address

	if confidence >= a.recognitionThreshold {
		islAddr = shape.toISLAddress()
	}

	a.thresholdHistory = append(a.thresholdHistory, confidence)
	if len(a.thresholdHistory) > 100 {
		a.adaptThreshold()
	}

	return &isl.ISLMessage{
		MsgType:     isl.MsgResponse,
		SenderID:    a.id,
		TargetID:    msg.SenderID,
		PrimaryAddr: islAddr,
		Confidence:  uint32(confidence * 100),
	}, nil
}

func (a *VisionAgent) OnLearn(_ context.Context, lesson *leoai.AgentLesson) error {
	a.dictVersion = lesson.DictVersion
	return nil
}

func (a *VisionAgent) OnShutdown() {
	log.Printf("Agent[VisionAgent]: Shutdown")
}

// ─────────────────────────────────────────────────────────────────
// QUADTREE DECOMPOSITION
// ─────────────────────────────────────────────────────────────────

type QuadNode struct {
	X, Y, W, H int
	HasEdge    bool
	Children   [4]*QuadNode
	Midpoints  [2]Point
}

type Point struct{ X, Y float32 }

type ShapeSignature struct {
	skeleton    []Point
	matchScore  float32
	matchedChar rune
}

func (a *VisionAgent) analyzeImage(imageData []byte) *ShapeSignature {
	width, height := 64, 64
	pixels := a.decodePixels(imageData, width, height)
	edges := a.detectEdges(pixels, width, height)
	root := a.buildQuadTree(edges, 0, 0, width, height, 0)
	skeleton := a.extractSkeleton(root)
	sig := &ShapeSignature{skeleton: skeleton}
	sig.matchScore, sig.matchedChar = a.matchWithUTF32(skeleton)
	return sig
}

func (a *VisionAgent) buildQuadTree(edges [][]bool, x, y, w, h, depth int) *QuadNode {
	node := &QuadNode{X: x, Y: y, W: w, H: h}
	cx, cy := x+w/2, y+h/2
	node.Midpoints[0] = Point{float32(cx), float32(cy)}
	node.HasEdge = a.hasEdgeInRegion(edges, x, y, w, h)

	if w <= 4 || h <= 4 || depth > 6 {
		return node
	}

	if !a.midpointOnEdge(edges, cx, cy) {
		hw, hh := w/2, h/2
		node.Children[0] = a.buildQuadTree(edges, x, y, hw, hh, depth+1)
		node.Children[1] = a.buildQuadTree(edges, x+hw, y, hw, hh, depth+1)
		node.Children[2] = a.buildQuadTree(edges, x, y+hh, hw, hh, depth+1)
		node.Children[3] = a.buildQuadTree(edges, x+hw, y+hh, hw, hh, depth+1)
	}
	return node
}

func (a *VisionAgent) extractSkeleton(root *QuadNode) []Point {
	var points []Point
	a.collectEdgePoints(root, &points)
	return a.connectPoints(points)
}

func (a *VisionAgent) collectEdgePoints(node *QuadNode, points *[]Point) {
	if node == nil {
		return
	}
	if node.HasEdge {
		*points = append(*points, Point{
			X: float32(node.X + node.W/2),
			Y: float32(node.Y + node.H/2),
		})
	}
	for _, child := range node.Children {
		a.collectEdgePoints(child, points)
	}
}

func (a *VisionAgent) connectPoints(points []Point) []Point {
	if len(points) < 2 {
		return points
	}
	ordered := []Point{points[0]}
	remaining := make([]Point, len(points)-1)
	copy(remaining, points[1:])

	for len(remaining) > 0 {
		last := ordered[len(ordered)-1]
		nearestIdx := 0
		nearestDist := float32(1e9)
		for i, p := range remaining {
			d := ptDist(last, p)
			if d < nearestDist {
				nearestDist = d
				nearestIdx = i
			}
		}
		ordered = append(ordered, remaining[nearestIdx])
		remaining = append(remaining[:nearestIdx], remaining[nearestIdx+1:]...)
	}
	return ordered
}

func (a *VisionAgent) matchWithUTF32(skeleton []Point) (float32, rune) {
	if len(skeleton) < 3 {
		return 0, 0
	}

	curvature := a.calculateCurvature(skeleton)
	aspectRatio := a.calculateAspectRatio(skeleton)
	symmetry := a.calculateSymmetry(skeleton)

	patterns := map[rune][3]float32{
		'🍌': {0.7, 0.3, 0.4},
		'🏠': {0.1, 0.8, 0.9},
		'🐱': {0.4, 0.6, 0.7},
		'🌳': {0.3, 0.4, 0.8},
	}

	bestScore := float32(0)
	var bestChar rune

	for char, pattern := range patterns {
		score := 1.0 - (ptAbs(curvature-pattern[0])+ptAbs(aspectRatio-pattern[1])+ptAbs(symmetry-pattern[2]))/3.0
		if score > bestScore {
			bestScore = score
			bestChar = char
		}
	}
	return bestScore, bestChar
}

func (s *ShapeSignature) toISLAddress() isl.Address {
	switch s.matchedChar {
	case '🍌':
		return isl.Address{Layer: 'B', Group: 'b', Type: 'b', ID: 1}
	case '🏠':
		return isl.Address{Layer: 'H', Group: 'a', Type: 'a', ID: 1}
	case '🐱':
		return isl.Address{Layer: 'B', Group: 'a', Type: 'a', ID: 5}
	default:
		return isl.Address{Layer: 'Z', Group: 'z', Type: 'z', ID: 0}
	}
}

func (a *VisionAgent) adaptThreshold() {
	if len(a.thresholdHistory) < 10 {
		return
	}
	var sum float32
	for _, v := range a.thresholdHistory {
		sum += v
	}
	mean := sum / float32(len(a.thresholdHistory))

	var variance float32
	for _, v := range a.thresholdHistory {
		diff := v - mean
		variance += diff * diff
	}
	stddev := float32(math.Sqrt(float64(variance / float32(len(a.thresholdHistory)))))

	a.recognitionThreshold = mean - stddev*0.5
	if a.recognitionThreshold < 0.3 {
		a.recognitionThreshold = 0.3
	}
	if a.recognitionThreshold > 0.9 {
		a.recognitionThreshold = 0.9
	}
	a.thresholdHistory = a.thresholdHistory[:0]
}

// ─────────────────────────────────────────────────────────────────
// SENSOR AGENT
// ─────────────────────────────────────────────────────────────────

type SensorAgent struct {
	SilentAgent
	lastValue float64
	threshold float64
	history   []float64
	islAddr   isl.Address
}

func (s *SensorAgent) ID() uint16   { return s.id }
func (s *SensorAgent) Name() string { return s.name }

func (s *SensorAgent) Poll(raw float64) *isl.ISLMessage {
	delta := math.Abs(raw - s.lastValue)
	if delta < s.threshold {
		return nil
	}

	s.lastValue = raw
	s.history = append(s.history, delta)

	if len(s.history) >= 50 {
		s.adaptSensorThreshold()
	}

	return &isl.ISLMessage{
		MsgType:     isl.MsgActivate,
		SenderID:    s.id,
		PrimaryAddr: s.islAddr,
		Payload:     encodeFloat64(raw),
		Timestamp:   uint32(time.Now().Unix()),
	}
}

func (s *SensorAgent) adaptSensorThreshold() {
	var sum float64
	for _, v := range s.history {
		sum += v
	}
	mean := sum / float64(len(s.history))

	var variance float64
	for _, v := range s.history {
		diff := v - mean
		variance += diff * diff
	}
	stddev := math.Sqrt(variance / float64(len(s.history)))
	s.threshold = stddev * 0.5
	if s.threshold < 0.01 {
		s.threshold = 0.01
	}
	s.history = s.history[:0]
}

func (s *SensorAgent) OnActivate(ctx context.Context, msg *isl.ISLMessage) (*isl.ISLMessage, error) {
	return nil, nil
}

func (s *SensorAgent) OnLearn(_ context.Context, lesson *leoai.AgentLesson) error {
	s.dictVersion = lesson.DictVersion
	return nil
}

func (s *SensorAgent) OnShutdown() {
	log.Printf("Agent[SensorAgent]: Shutdown")
}

// ─────────────────────────────────────────────────────────────────
// HELPERS
// ─────────────────────────────────────────────────────────────────

func (a *VisionAgent) decodePixels(data []byte, w, h int) [][]uint8 {
	pixels := make([][]uint8, h)
	for i := range pixels {
		pixels[i] = make([]uint8, w)
	}
	for i, b := range data {
		if i >= w*h {
			break
		}
		pixels[i/w][i%w] = b
	}
	return pixels
}

func (a *VisionAgent) detectEdges(pixels [][]uint8, w, h int) [][]bool {
	edges := make([][]bool, h)
	for i := range edges {
		edges[i] = make([]bool, w)
	}
	for y := 1; y < h-1; y++ {
		for x := 1; x < w-1; x++ {
			gx := int(pixels[y][x+1]) - int(pixels[y][x-1])
			gy := int(pixels[y+1][x]) - int(pixels[y-1][x])
			magnitude := math.Sqrt(float64(gx*gx + gy*gy))
			edges[y][x] = magnitude > 30
		}
	}
	return edges
}

func (a *VisionAgent) hasEdgeInRegion(edges [][]bool, x, y, w, h int) bool {
	for dy := 0; dy < h && y+dy < len(edges); dy++ {
		for dx := 0; dx < w && x+dx < len(edges[0]); dx++ {
			if edges[y+dy][x+dx] {
				return true
			}
		}
	}
	return false
}

func (a *VisionAgent) midpointOnEdge(edges [][]bool, cx, cy int) bool {
	if cy < len(edges) && cx < len(edges[0]) {
		return edges[cy][cx]
	}
	return false
}

func (a *VisionAgent) calculateCurvature(pts []Point) float32 {
	if len(pts) < 3 {
		return 0
	}
	var total float32
	for i := 1; i < len(pts)-1; i++ {
		dx1 := pts[i].X - pts[i-1].X
		dy1 := pts[i].Y - pts[i-1].Y
		dx2 := pts[i+1].X - pts[i].X
		dy2 := pts[i+1].Y - pts[i].Y
		cross := dx1*dy2 - dy1*dx2
		total += float32(math.Abs(float64(cross)))
	}
	return total / float32(len(pts))
}

func (a *VisionAgent) calculateAspectRatio(pts []Point) float32 {
	if len(pts) == 0 {
		return 0
	}
	minX, maxX, minY, maxY := pts[0].X, pts[0].X, pts[0].Y, pts[0].Y
	for _, p := range pts {
		if p.X < minX { minX = p.X }
		if p.X > maxX { maxX = p.X }
		if p.Y < minY { minY = p.Y }
		if p.Y > maxY { maxY = p.Y }
	}
	if maxX-minX == 0 {
		return 0
	}
	return (maxY - minY) / (maxX - minX)
}

func (a *VisionAgent) calculateSymmetry(pts []Point) float32 {
	if len(pts) < 2 {
		return 0
	}
	var cx float32
	for _, p := range pts {
		cx += p.X
	}
	cx /= float32(len(pts))
	var sym float32
	for _, p := range pts {
		sym += float32(math.Abs(float64(p.X - cx)))
	}
	return 1.0 - sym/float32(len(pts)*100)
}

func ptDist(a, b Point) float32 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	return float32(math.Sqrt(float64(dx*dx + dy*dy)))
}

func ptAbs(x float32) float32 {
	if x < 0 {
		return -x
	}
	return x
}

func encodeFloat64(v float64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, math.Float64bits(v))
	return b
}
