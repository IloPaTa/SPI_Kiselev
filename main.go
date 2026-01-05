package main

import (
	"fmt"
	"time"
)

// AppConfig - структура конфигурации
type AppConfig struct {
	AppName        string
	Version        string
	DatabaseURL    string
	LogLevel       string
	MaxConnections int
}

// instance - переменная для хранения единственного экземпляра
var instance *AppConfig

// GetInstance - метод получения экземпляра (НЕ ПОТОКОБЕЗОПАСНЫЙ!)
func GetInstance() *AppConfig {
	if instance == nil {
		// Имитация тяжелой инициализации
		fmt.Println("[SINGLETON] Initializing AppConfig instance...")
		time.Sleep(100 * time.Millisecond)

		instance = &AppConfig{
			AppName:        "MySuperApp",
			Version:        "1.0.0",
			DatabaseURL:    "postgres://localhost:5432/db",
			LogLevel:       "INFO",
			MaxConnections: 10,
		}
	} else {
		fmt.Println("[SINGLETON] Returning existing instance")
	}
	return instance
}

// Print - метод вывода конфигурации
func (c *AppConfig) Print() {
	fmt.Printf("Config: {App: %s, Ver: %s, DB: %s, Log: %s}\n",
		c.AppName, c.Version, c.DatabaseURL, c.LogLevel)
}

func main() {
	fmt.Println("=== Singleton Pattern Demo (Unsafe) ===")

	// Тест 1: Последовательный вызов (работает нормально)
	fmt.Println("\n--- Sequential Access ---")
	cfg1 := GetInstance()
	cfg2 := GetInstance()

	if cfg1 == cfg2 {
		fmt.Println("SUCCESS: Both variables point to the same instance")
	} else {
		fmt.Println("FAIL: Variables point to different instances")
	}

	// Тест 2: Конкурентный вызов (покажет проблему)
	// В этой версии инициализация может произойти ДВАЖДЫ
	fmt.Println("\n--- Concurrent Access ---")

	// Сбрасываем инстанс для теста
	instance = nil

	for i := 0; i < 5; i++ {
		go func(id int) {
			cfg := GetInstance()
			fmt.Printf("Goroutine %d got config: %p\n", id, cfg)
		}(i)
	}

	// Ждем завершения горутин
	time.Sleep(1 * time.Second)
}
