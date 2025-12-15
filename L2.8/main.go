package main

import (
	"fmt"
	"os"

	"github.com/beevik/ntp"
)

// L2.8 Получение точного времени (NTP)
func main() {
	// Получаем время через NTP
	time, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		// Выводим ошибку в STDERR и завершаем с ненулевым кодом
		fmt.Fprintf(os.Stderr, "error getting time: %v\n", err)
		os.Exit(1)
	}
	// Выводим время
	fmt.Println(time)
}
