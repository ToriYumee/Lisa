package whatsapp

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"Lisa/pkg/types"

	"go.mau.fi/whatsmeow/types/events"
)

// MediaHandler maneja diferentes tipos de medios
type MediaHandler func(media *types.MediaMessage)

type MediaHandlers struct {
	Image    MediaHandler
	Audio    MediaHandler
	Video    MediaHandler
	Document MediaHandler
}

func (c *Client) SetMediaHandlers(handlers MediaHandlers) {
	// Método simplificado por ahora
	log.Println("WA: Media handlers configurados")
}

func (c *Client) processComplexMessage(msg *events.Message) {
	// Determinar tipo de mensaje
	messageType := types.GetMessageType(msg)

	// Obtener información básica del mensaje
	groupName := ""
	if msg.Info.IsGroup {
		groupInfo, err := c.whatsAppClient.GetGroupInfo(msg.Info.Chat)
		if err == nil {
			groupName = groupInfo.Name
		} else {
			groupName = "Grupo desconocido"
		}
	}

	// Procesar según el tipo de mensaje
	switch messageType {
	case types.MessageTypeText:
		c.handleTextMessage(msg, groupName)
	case types.MessageTypeImage:
		c.handleImageMessage(msg, groupName)
	case types.MessageTypeAudio:
		c.handleAudioMessage(msg, groupName)
	case types.MessageTypeVideo:
		c.handleVideoMessage(msg, groupName)
	case types.MessageTypeDocument:
		c.handleDocumentMessage(msg, groupName)
	default:
		c.logMessage(msg, groupName, fmt.Sprintf("[%s] No soportado", messageType.String()))
	}
}

func (c *Client) handleTextMessage(msg *events.Message, groupName string) {
	text := msg.Message.GetConversation()
	if text == "" && msg.Message.ExtendedTextMessage != nil {
		text = msg.Message.ExtendedTextMessage.GetText()
	}

	if text != "" {
		c.logMessage(msg, groupName, text)
	}
}

func (c *Client) handleImageMessage(msg *events.Message, groupName string) {
	imageMsg := msg.Message.GetImageMessage()
	if imageMsg == nil {
		return
	}

	caption := imageMsg.GetCaption()
	mimetype := imageMsg.GetMimetype()

	// Log del mensaje
	logText := fmt.Sprintf("[IMAGEN] %s", mimetype)
	if caption != "" {
		logText += fmt.Sprintf(" - Caption: %s", caption)
	}
	c.logMessage(msg, groupName, logText)

	// TODO: Implementar descarga de imagen cuando sea necesario
}

func (c *Client) handleAudioMessage(msg *events.Message, groupName string) {
	audioMsg := msg.Message.GetAudioMessage()
	if audioMsg == nil {
		return
	}

	mimetype := audioMsg.GetMimetype()
	duration := audioMsg.GetSeconds()

	logText := fmt.Sprintf("[AUDIO] %s - Duracion: %ds", mimetype, duration)
	c.logMessage(msg, groupName, logText)

	// TODO: Implementar descarga de audio cuando sea necesario
}

func (c *Client) handleVideoMessage(msg *events.Message, groupName string) {
	videoMsg := msg.Message.GetVideoMessage()
	if videoMsg == nil {
		return
	}

	caption := videoMsg.GetCaption()
	mimetype := videoMsg.GetMimetype()
	duration := videoMsg.GetSeconds()

	logText := fmt.Sprintf("[VIDEO] %s - Duracion: %ds", mimetype, duration)
	if caption != "" {
		logText += fmt.Sprintf(" - Caption: %s", caption)
	}
	c.logMessage(msg, groupName, logText)

	// TODO: Implementar descarga de video cuando sea necesario
}

func (c *Client) handleDocumentMessage(msg *events.Message, groupName string) {
	docMsg := msg.Message.GetDocumentMessage()
	if docMsg == nil {
		return
	}

	filename := docMsg.GetFileName()
	mimetype := docMsg.GetMimetype()
	fileSize := docMsg.GetFileLength()

	logText := fmt.Sprintf("[DOCUMENTO] %s (%s) - Size: %d bytes", filename, mimetype, fileSize)
	c.logMessage(msg, groupName, logText)

	// TODO: Implementar descarga de documento cuando sea necesario
}

func (c *Client) logMessage(msg *events.Message, groupName, content string) {
	if msg.Info.IsGroup {
		log.Printf("WA [GRUPO:%s] %s: %s", groupName, msg.Info.PushName, content)
	} else {
		log.Printf("WA [PRIVADO] %s: %s", msg.Info.PushName, content)
	}
}

// Helper para determinar si un archivo es de texto y puede ser procesado
func (c *Client) isProcessableTextFile(filename, mimetype string) bool {
	textMimeTypes := []string{
		"text/plain",
		"text/csv",
		"application/json",
		"text/xml",
		"application/xml",
	}

	textExtensions := []string{
		".txt", ".log", ".csv", ".json", ".xml", ".md",
	}

	// Verificar mimetype
	for _, mt := range textMimeTypes {
		if strings.Contains(mimetype, mt) {
			return true
		}
	}

	// Verificar extensión
	ext := strings.ToLower(filepath.Ext(filename))
	for _, textExt := range textExtensions {
		if ext == textExt {
			return true
		}
	}

	return false
}
