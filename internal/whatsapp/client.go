package whatsapp

import (
	"context"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"

	"Lisa/internal/config"
	"Lisa/pkg/types"
)

type Client struct {
	whatsAppClient *whatsmeow.Client
	container      *sqlstore.Container
	messageHandler MessageHandler
	logger         waLog.Logger
	ctx            context.Context
	cancel         context.CancelFunc
}

type MessageHandler func(*events.Message)

type Config struct {
	DatabaseURI string
	LogLevel    string
}

func NewClient(cfg *config.Config) (*Client, error) {
	ctx, cancel := context.WithCancel(context.Background())

	// Crear logger
	logger := waLog.Stdout("WhatsApp", cfg.WhatsApp.LogLevel, true)

	// Conectar al almacenamiento PostgreSQL
	container, err := sqlstore.New(ctx, "postgres", cfg.WhatsApp.DatabaseURI, logger)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("fallo al conectar a PostgreSQL: %v", err)
	}

	// Obtener o crear dispositivo
	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("fallo al obtener el dispositivo: %v", err)
	}

	// Crear cliente WhatsApp
	whatsAppClient := whatsmeow.NewClient(deviceStore, nil)

	client := &Client{
		whatsAppClient: whatsAppClient,
		container:      container,
		logger:         logger,
		ctx:            ctx,
		cancel:         cancel,
	}

	// Configurar manejadores de eventos
	client.setupEventHandlers()

	return client, nil
}

func (c *Client) setupEventHandlers() {
	c.whatsAppClient.AddEventHandler(func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Message:
			c.handleMessage(v)
		case *events.Receipt:
			c.logger.Infof("Mensaje entregado: %s", v.MessageIDs[0])
		case *events.Connected:
			c.logger.Infof("Cliente WhatsApp conectado")
		case *events.Disconnected:
			c.logger.Warnf("Cliente WhatsApp desconectado")
		case *events.LoggedOut:
			c.logger.Warnf("Cliente WhatsApp sesion cerrada")
		}
	})
}

func (c *Client) handleMessage(msg *events.Message) {
	// Ignorar mensajes enviados por nosotros
	if msg.Info.IsFromMe {
		return
	}

	// Obtener información del grupo si es necesario
	groupName := ""
	if msg.Info.IsGroup {
		groupInfo, err := c.whatsAppClient.GetGroupInfo(msg.Info.Chat)
		if err == nil {
			groupName = groupInfo.Name
		} else {
			groupName = "Grupo desconocido"
		}
	}

	// Intentar obtener texto del mensaje
	text := msg.Message.GetConversation()
	if text == "" && msg.Message.GetExtendedTextMessage() != nil {
		text = msg.Message.GetExtendedTextMessage().GetText()
	}

	// Si hay texto, es un mensaje de texto
	if text != "" {
		if msg.Info.IsGroup {
			log.Printf("WA [GRUPO:%s] %s: %s", groupName, msg.Info.PushName, text)
		} else {
			log.Printf("WA [PRIVADO] %s: %s", msg.Info.PushName, text)
		}
	} else {
		// Para otros tipos de mensaje, determinar el tipo
		messageType := types.GetMessageType(msg)

		// Log básico para tipos no-texto
		typeStr := messageType.String()
		if msg.Info.IsGroup {
			log.Printf("WA [GRUPO:%s] %s: [%s]", groupName, msg.Info.PushName, typeStr)
		} else {
			log.Printf("WA [PRIVADO] %s: [%s]", msg.Info.PushName, typeStr)
		}

		// Procesar tipos específicos si es necesario
		if messageType != types.MessageTypeUnknown {
			c.processComplexMessage(msg)
		}
	}

	// Llamar al handler personalizado si existe
	if c.messageHandler != nil {
		c.messageHandler(msg)
	}
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
	case types.MessageTypeImage:
		c.handleImageMessage(msg, groupName)
	case types.MessageTypeAudio:
		c.handleAudioMessage(msg, groupName)
	case types.MessageTypeVideo:
		c.handleVideoMessage(msg, groupName)
	case types.MessageTypeDocument:
		c.handleDocumentMessage(msg, groupName)
	case types.MessageTypeSticker:
		c.logMessage(msg, groupName, "[STICKER]")
	case types.MessageTypeContact:
		c.logMessage(msg, groupName, "[CONTACTO]")
	case types.MessageTypeLocation:
		c.logMessage(msg, groupName, "[UBICACIÓN]")
	default:
		c.logMessage(msg, groupName, fmt.Sprintf("[%s] No soportado", messageType.String()))
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
	logText := "[IMAGEN]"
	if mimetype != "" {
		logText += fmt.Sprintf(" (%s)", mimetype)
	}
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

	logText := "[AUDIO]"
	if mimetype != "" {
		logText += fmt.Sprintf(" (%s)", mimetype)
	}
	if duration > 0 {
		logText += fmt.Sprintf(" - %ds", duration)
	}
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

	logText := "[VIDEO]"
	if mimetype != "" {
		logText += fmt.Sprintf(" (%s)", mimetype)
	}
	if duration > 0 {
		logText += fmt.Sprintf(" - %ds", duration)
	}
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

	logText := "[DOCUMENTO]"
	if filename != "" {
		logText += fmt.Sprintf(" %s", filename)
	}
	if mimetype != "" {
		logText += fmt.Sprintf(" (%s)", mimetype)
	}
	if fileSize > 0 {
		logText += fmt.Sprintf(" - %d bytes", fileSize)
	}
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

func (c *Client) SetMessageHandler(handler MessageHandler) {
	c.messageHandler = handler
}

func (c *Client) Connect() error {
	// Verificar si ya hay una sesión
	if c.whatsAppClient.Store.ID == nil {
		// Primera vez - necesita QR
		qrChan, _ := c.whatsAppClient.GetQRChannel(c.ctx)
		err := c.whatsAppClient.Connect()
		if err != nil {
			return fmt.Errorf("no se pudo conectar: %v", err)
		}

		log.Println("WA: Escanea el codigo QR para iniciar sesion...")

		for evt := range qrChan {
			if evt.Event == "code" {
				// Usar la nueva función para mostrar el QR en terminal
				c.DisplayQRInTerminal(evt.Code)
			} else if evt.Event == "success" {
				log.Println("WA: QR escaneado exitosamente")
				break
			}
		}
	} else {
		// Sesión existente - conectar directamente
		err := c.whatsAppClient.Connect()
		if err != nil {
			return fmt.Errorf("no se pudo conectar: %v", err)
		}
		log.Println("WA: Sesion restaurada exitosamente")
	}

	return nil
}

func (c *Client) Disconnect() {
	if c.whatsAppClient != nil {
		c.whatsAppClient.Disconnect()
	}
	if c.cancel != nil {
		c.cancel()
	}
	log.Println("WA: Cliente desconectado")
}

func (c *Client) IsConnected() bool {
	return c.whatsAppClient != nil && c.whatsAppClient.IsConnected()
}

func (c *Client) GetJID() string {
	if c.whatsAppClient.Store.ID == nil {
		return ""
	}
	return c.whatsAppClient.Store.ID.String()
}

// Métodos para enviar mensajes (para uso futuro)
func (c *Client) SendTextMessage(jid, text string) error {
	// TODO: Implementar envío de mensajes si es necesario
	return fmt.Errorf("envio de mensajes no implementado aun")
}
