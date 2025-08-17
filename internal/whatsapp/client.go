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

	// Para mensajes complejos (imágenes, videos, etc.), usar el procesador avanzado
	messageType := types.GetMessageType(msg)
	if messageType != types.MessageTypeText {
		c.processComplexMessage(msg)
		return
	}

	// Para texto básico, manejar aquí
	text := msg.Message.GetConversation()
	if text == "" {
		if msg.Message.ExtendedTextMessage != nil {
			text = msg.Message.ExtendedTextMessage.GetText()
		}
	}

	if text != "" {
		if msg.Info.IsGroup {
			groupInfo, err := c.whatsAppClient.GetGroupInfo(msg.Info.Chat)
			groupName := "Grupo desconocido"
			if err == nil {
				groupName = groupInfo.Name
			}

			log.Printf("WA [GRUPO:%s] %s: %s",
				groupName,
				msg.Info.PushName,
				text,
			)
		} else {
			log.Printf("WA [PRIVADO] %s: %s",
				msg.Info.PushName,
				text,
			)
		}
	}

	// Llamar al handler personalizado si existe
	if c.messageHandler != nil {
		c.messageHandler(msg)
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
				log.Printf("WA: Codigo QR: %s", evt.Code)
				// TODO: Aquí podríamos generar la imagen del QR
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
