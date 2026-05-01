package play_test

import (
	"encoding/json"
	"sync"
	"testing"

	"github.com/vpoliakov01/2v2ChessAI/engine/game"
	"github.com/vpoliakov01/2v2ChessAI/engine/play"
)

// Player constants. The game package uses raw ints 0..3 (see engine/game/game.go
// and ui/src/utils/PlayerColors). Naming them here keeps tests readable and
// matches the order the FE uses.
const (
	playerRed    game.Player = 0
	playerBlue   game.Player = 1
	playerYellow game.Player = 2
	playerGreen  game.Player = 3
)

// Message represents a single Message captured by testConnection.
type Message struct {
	MessageType int
	Raw         []byte
	Parsed      play.Message
}

// testConnection wraps a play.Connection together with a thread-safe in-memory
// MessageWriter, so a single object exposes both the dispatch surface
// (ProcessMessage) and the assertion surface (MessagesOfType).
type testConnection struct {
	*play.Connection

	NewMessageCh chan Message

	t        *testing.T
	mu       sync.Mutex
	messages []Message
}

// defaultTestConfig returns a config with all four seats marked as human, so
// the engine never plays moves on its own during a test.
func defaultTestConfig() *play.Config {
	return &play.Config{
		Depth:        1,
		CaptureDepth: 1,
		HumanPlayers: []game.Player{playerRed, playerBlue, playerYellow, playerGreen},
		EvalLimit:    1,
	}
}

// NewConnection builds a testConnection backed by a mock MessageWriter so
// tests can invoke ProcessMessage and inspect captured outbound messages.
func NewConnection(t *testing.T, cfg *play.Config) *testConnection {
	t.Helper()
	if cfg == nil {
		cfg = defaultTestConfig()
	}

	conn := &testConnection{
		t:            t,
		NewMessageCh: make(chan Message, 100),
		messages:     make([]Message, 0),
		mu:           sync.Mutex{},
	}
	conn.Connection = play.NewConnection(conn, cfg)
	return conn
}

// WriteMessage satisfies play.MessageWriter and stores the raw + parsed payload.
func (tc *testConnection) WriteMessage(messageType int, data []byte) error {
	parsed := play.Message{}
	_ = json.Unmarshal(data, &parsed)

	tc.mu.Lock()
	defer tc.mu.Unlock()

	msg := Message{
		MessageType: messageType,
		Raw:         append([]byte(nil), data...),
		Parsed:      parsed,
	}

	tc.messages = append(tc.messages, msg)
	tc.NewMessageCh <- msg
	return nil
}

// Messages returns a snapshot of all captured messages.
func (tc *testConnection) Messages() []Message {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	out := make([]Message, len(tc.messages))
	copy(out, tc.messages)
	return out
}

// MessagesOfType returns captured messages filtered by play.MessageType.
func (tc *testConnection) MessagesOfType(t play.MessageType) []Message {
	all := tc.Messages()
	out := make([]Message, 0, len(all))
	for _, msg := range all {
		if msg.Parsed.Type == t {
			out = append(out, msg)
		}
	}
	return out
}

// WaitForMessagesOfType waits for the specified number of messages of the given type and returns them.
func (tc *testConnection) WaitForMessagesOfType(msgType play.MessageType, count int) []Message {
	tc.t.Helper()
	msgs := make([]Message, 0, count)
	for len(msgs) < count {
		msg := <-tc.NewMessageCh
		if msg.Parsed.Type == msgType {
			msgs = append(msgs, msg)
		}
	}
	return msgs
}

// WaitForMessages waits for the specified number of messages and returns them.
func (tc *testConnection) WaitForMessages(count int) []Message {
	tc.t.Helper()
	msgs := make([]Message, 0, count)
	for len(msgs) < count {
		msg := <-tc.NewMessageCh
		msgs = append(msgs, msg)
	}
	return msgs
}

func (tc *testConnection) ProcessMessage(msgType play.MessageType, data interface{}) {
	tc.t.Helper()
	tc.Connection.ProcessMessage(&play.Message{Type: msgType, Data: data})
}

func availableMovesFromMessage(t *testing.T, msg Message) []play.PGNMove {
	t.Helper()
	return dataFromMessage[[]play.PGNMove](t, msg)
}

func dataFromMessage[T any](t *testing.T, msg Message) T {
	t.Helper()
	var out T
	raw, err := json.Marshal(msg.Parsed.Data)
	if err != nil {
		t.Fatalf("failed to remarshal message data: %v", err)
	}
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("failed to decode message payload into %T: %v", out, err)
	}
	return out
}

// requireSingleMessage asserts that exactly one message of the given type was
// captured and returns it.
func requireSingleMessage(t *testing.T, conn *testConnection, msgType play.MessageType) Message {
	t.Helper()
	msgs := conn.MessagesOfType(msgType)
	if len(msgs) != 1 {
		t.Fatalf("expected exactly 1 %q message, got %d", msgType, len(msgs))
	}
	return msgs[0]
}

// requireOnlyMessage asserts that no messages besides a message of the given type were captured.
func requireOnlyMessage(t *testing.T, conn *testConnection, msgType play.MessageType) Message {
	t.Helper()
	msgs := conn.Messages()
	if len(msgs) != 1 {
		t.Fatalf("expected exactly 1 %q message, got %d", msgType, len(msgs))
	}
	if msgs[0].Parsed.Type != msgType {
		t.Fatalf("expected only %q message, got %q", msgType, msgs[0].Parsed.Type)
	}
	return msgs[0]
}
