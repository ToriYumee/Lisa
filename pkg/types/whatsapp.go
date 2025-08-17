package types

import (
	"time"

	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

// MessageInfo contiene información básica de un mensaje de WhatsApp
type MessageInfo struct {
	ID        string    `json:"id"`
	From      string    `json:"from"`
	To        string    `json:"to"`
	PushName  string    `json:"push_name"`
	Text      string    `json:"text"`
	Timestamp time.Time `json:"timestamp"`
	IsGroup   bool      `json:"is_group"`
	GroupName string    `json:"group_name,omitempty"`
	IsFromMe  bool      `json:"is_from_me"`
}

// MediaMessage representa un archivo multimedia de WhatsApp
type MediaMessage struct {
	Info      MessageInfo `json:"info"`
	Type      MessageType `json:"type"`
	MimeType  string      `json:"mime_type"`
	Filename  string      `json:"filename,omitempty"`
	Caption   string      `json:"caption,omitempty"`
	Data      []byte      `json:"-"` // No serializar los datos binarios
	Size      int64       `json:"size"`
	Duration  int         `json:"duration,omitempty"` // Para audio/video en segundos
	GroupName string      `json:"group_name,omitempty"`

	// Campos para procesamiento
	ProcessedAt     time.Time `json:"processed_at"`
	NeedsProcessing bool      `json:"needs_processing"`
	TextContent     string    `json:"text_content,omitempty"`  // Para documentos de texto
	Transcription   string    `json:"transcription,omitempty"` // Para audio transcrito
}

// GroupInfo información de un grupo de WhatsApp
type GroupInfo struct {
	JID              types.JID `json:"jid"`
	Name             string    `json:"name"`
	Topic            string    `json:"topic"`
	Owner            types.JID `json:"owner"`
	CreatedAt        time.Time `json:"created_at"`
	ParticipantCount int       `json:"participant_count"`
}

// ContactInfo información de un contacto
type ContactInfo struct {
	JID      types.JID `json:"jid"`
	PushName string    `json:"push_name"`
	Name     string    `json:"name"`
	Phone    string    `json:"phone"`
}

// MessageType tipos de mensaje que podemos procesar
type MessageType int

const (
	MessageTypeText MessageType = iota
	MessageTypeImage
	MessageTypeDocument
	MessageTypeAudio
	MessageTypeVideo
	MessageTypeSticker
	MessageTypeContact
	MessageTypeLocation
	MessageTypeUnknown
)

func (mt MessageType) String() string {
	switch mt {
	case MessageTypeText:
		return "text"
	case MessageTypeImage:
		return "image"
	case MessageTypeDocument:
		return "document"
	case MessageTypeAudio:
		return "audio"
	case MessageTypeVideo:
		return "video"
	case MessageTypeSticker:
		return "sticker"
	case MessageTypeContact:
		return "contact"
	case MessageTypeLocation:
		return "location"
	default:
		return "unknown"
	}
}

// MessagePriority prioridades de mensaje para procesamiento
type MessagePriority int

const (
	PriorityLow MessagePriority = iota
	PriorityNormal
	PriorityHigh
	PriorityCritical
)

func (mp MessagePriority) String() string {
	switch mp {
	case PriorityLow:
		return "low"
	case PriorityNormal:
		return "normal"
	case PriorityHigh:
		return "high"
	case PriorityCritical:
		return "critical"
	default:
		return "normal"
	}
}

// Funciones helper para convertir desde whatsmeow events

func NewMessageInfoFromEvent(msg *events.Message) MessageInfo {
	text := msg.Message.GetConversation()
	if text == "" && msg.Message.ExtendedTextMessage != nil {
		text = msg.Message.ExtendedTextMessage.GetText()
	}

	return MessageInfo{
		ID:        msg.Info.ID,
		From:      msg.Info.Chat.String(),
		PushName:  msg.Info.PushName,
		Text:      text,
		Timestamp: msg.Info.Timestamp,
		IsGroup:   msg.Info.IsGroup,
		IsFromMe:  msg.Info.IsFromMe,
	}
}

func NewWhatsAppMessageFromEvent(msg *events.Message) MediaMessage {
	return MediaMessage{
		Info:            NewMessageInfoFromEvent(msg),
		ProcessedAt:     time.Now(),
		NeedsProcessing: true,
	}
}

// GetMessageType determina el tipo de mensaje
func GetMessageType(msg *events.Message) MessageType {
	switch {
	case msg.Message.Conversation != nil:
		return MessageTypeText
	case msg.Message.ExtendedTextMessage != nil:
		return MessageTypeText
	case msg.Message.ImageMessage != nil:
		return MessageTypeImage
	case msg.Message.DocumentMessage != nil:
		return MessageTypeDocument
	case msg.Message.AudioMessage != nil:
		return MessageTypeAudio
	case msg.Message.VideoMessage != nil:
		return MessageTypeVideo
	case msg.Message.StickerMessage != nil:
		return MessageTypeSticker
	case msg.Message.ContactMessage != nil:
		return MessageTypeContact
	case msg.Message.LocationMessage != nil:
		return MessageTypeLocation
	default:
		return MessageTypeUnknown
	}
}
