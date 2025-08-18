package whatsapp

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mdp/qrterminal/v3"
)

// DisplayQRInTerminal muestra el código QR en la terminal
func (c *Client) DisplayQRInTerminal(code string) {
	log.Println("WA: Generando codigo QR...")

	// Primero intentar con configuración básica para Windows
	config := qrterminal.Config{
		Level:     qrterminal.M,
		Writer:    os.Stdout,
		BlackChar: "██", // Usar bloques sólidos
		WhiteChar: "  ", // Usar espacios dobles
		QuietZone: 1,
	}

	fmt.Println("\n" + strings.Repeat("=", 40))
	fmt.Println("CODIGO QR PARA WHATSAPP")
	fmt.Println(strings.Repeat("=", 40))

	// Intentar generar QR
	qrterminal.GenerateWithConfig(code, config)

	fmt.Println(strings.Repeat("=", 40))
	fmt.Println("Si no ves el QR, usa este codigo:")
	fmt.Printf("URL: %s\n", code)
	fmt.Println("Ve a WhatsApp > Dispositivos vinculados")
	fmt.Println(strings.Repeat("=", 40))
}
