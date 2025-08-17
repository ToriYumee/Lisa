package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"Lisa/internal/config"
	"Lisa/internal/whatsapp"
)

func main() {
	// Banner de inicio
	printBanner()

	log.Println("INICIANDO: Lisa Bot...")

	// Cargar configuración
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("ERROR: Error cargando configuracion: %v", err)
	}

	log.Printf("OK: Configuracion cargada correctamente")
	log.Printf("MODO: %s", cfg.Server.Environment)
	log.Printf("PUERTO: %s", cfg.Server.Port)
	log.Printf("LOG LEVEL: %s", cfg.Server.LogLevel)

	// Mostrar estado de las configuraciones
	printConfigStatus(cfg)

	// Crear contexto con cancelación
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Inicializar servicios
	log.Println("INIT: Inicializando servicios...")

	// 1. WhatsApp Client
	log.Println("WA: Inicializando cliente WhatsApp...")
	waClient, err := whatsapp.NewClient(cfg)
	if err != nil {
		log.Fatalf("ERROR: No se pudo crear cliente WhatsApp: %v", err)
	}

	// Conectar WhatsApp
	log.Println("WA: Conectando a WhatsApp...")
	if err := waClient.Connect(); err != nil {
		log.Fatalf("ERROR: No se pudo conectar a WhatsApp: %v", err)
	}

	// Configurar cleanup al salir
	defer func() {
		log.Println("WA: Desconectando WhatsApp...")
		waClient.Disconnect()
	}()

	// TODO: Inicializar otros servicios
	/*
		// 2. Discord Bot
		log.Println("DC: Inicializando bot de Discord...")

		// 3. Gemini AI
		log.Println("AI: Inicializando Gemini AI...")

		// 4. Jira Client
		log.Println("JIRA: Inicializando cliente Jira...")
	*/

	// Goroutine para mostrar status
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				log.Println("CONTEXTO: Contexto cancelado, deteniendo servicios...")
				return
			case <-ticker.C:
				waStatus := "DESCONECTADO"
				if waClient.IsConnected() {
					waStatus = "CONECTADO"
				}
				log.Printf("STATUS: Lisa Bot funcionando - WhatsApp: %s", waStatus)
			}
		}
	}()

	// TODO: Aquí inicializaremos los servicios uno por uno
	/*
		// Inicializar servicios
		log.Println("INIT: Inicializando servicios...")

		// 1. Base de datos
		log.Println("DB: Conectando a PostgreSQL...")

		// 2. WhatsApp Client
		log.Println("WA: Inicializando cliente WhatsApp...")

		// 3. Discord Bot
		log.Println("DC: Inicializando bot de Discord...")

		// 4. Gemini AI
		log.Println("AI: Inicializando Gemini AI...")

		// 5. Jira Client
		log.Println("JIRA: Inicializando cliente Jira...")
	*/

	log.Println("OK: Lisa Bot iniciado correctamente")
	log.Println("INFO: Presiona Ctrl+C para detener el bot")

	// Canal para manejar señales del sistema
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Esperar señal de interrupción
	<-sigChan
	log.Println("STOP: Senal de interrupcion recibida, deteniendo Lisa Bot...")

	// Dar tiempo para cleanup graceful
	cancel()
	time.Sleep(2 * time.Second)

	log.Println("DONE: Lisa Bot detenido correctamente")
}

func printBanner() {
	banner := `
╔══════════════════════════════════════════════════════════════════╗
║                                                                  ║
║    ██╗     ██╗███████╗ █████╗     ██████╗  ██████╗ ████████╗    ║
║    ██║     ██║██╔════╝██╔══██╗    ██╔══██╗██╔═══██╗╚══██╔══╝    ║
║    ██║     ██║███████╗███████║    ██████╔╝██║   ██║   ██║       ║
║    ██║     ██║╚════██║██╔══██║    ██╔══██╗██║   ██║   ██║       ║
║    ███████╗██║███████║██║  ██║    ██████╔╝╚██████╔╝   ██║       ║
║    ╚══════╝╚═╝╚══════╝╚═╝  ╚═╝    ╚═════╝  ╚═════╝    ╚═╝       ║
║                                                                  ║
║    Discord + WhatsApp + Jira + AI Assistant Bot                 ║
║    Version: 1.0.0                                               ║
║                                                                  ║
╚══════════════════════════════════════════════════════════════════╝
`
	fmt.Println(banner)
}

func printConfigStatus(cfg *config.Config) {
	log.Println("CONFIG: Estado de configuraciones:")

	// Discord
	if cfg.Discord.Token != "" {
		log.Println("  OK: Discord: Configurado")
	} else {
		log.Println("  ERROR: Discord: Token faltante")
	}

	// WhatsApp (PostgreSQL)
	if cfg.WhatsApp.Password != "" {
		log.Printf("  OK: WhatsApp: PostgreSQL configurado (%s:%d/%s)",
			cfg.WhatsApp.Host, cfg.WhatsApp.Port, cfg.WhatsApp.Database)
	} else {
		log.Println("  ERROR: WhatsApp: Configuracion de PostgreSQL incompleta")
	}

	// Jira
	if cfg.Jira.URL != "" && cfg.Jira.Token != "" {
		log.Printf("  OK: Jira: Configurado (%s)", cfg.Jira.URL)
	} else {
		log.Println("  WARN: Jira: Configuracion incompleta (opcional)")
	}

	// Gemini AI
	if cfg.Gemini.APIKey != "" {
		log.Printf("  OK: Gemini AI: Configurado (modelo: %s)", cfg.Gemini.Model)
	} else {
		log.Println("  ERROR: Gemini AI: API Key faltante")
	}

	log.Println("")
}
