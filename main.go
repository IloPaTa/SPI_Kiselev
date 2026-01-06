package main

import (
	"fmt"
	"sync"
)

// AppConfig - структура конфигурации
type AppConfig struct {
	AppName        string
	Version        string
	DatabaseURL    string
	LogLevel       string
	MaxConnections int
	mu             sync.Mutex // Мьютекс для безопасного изменения полей
}

var (
	instance *AppConfig
	once     sync.Once // Примитив для однократного выполнения
)

// GetInstance - потокобезопасный метод получения экземпляра
func GetInstance() *AppConfig {
	// once.Do гарантирует, что функция внутри выполнится ровно один раз
	// даже при вызове из 1000 горутин одновременно
	once.Do(func() {
		fmt.Println("[SINGLETON] Creating AppConfig instance (ONCE)...")
		instance = &AppConfig{
			AppName:        "MySuperApp",
			Version:        "1.0.0",
			DatabaseURL:    "postgres://localhost:5432/db",
			LogLevel:       "INFO",
			MaxConnections: 10,
		}
	})
	return instance
}

// SetLogLevel - потокобезопасный сеттер
func (c *AppConfig) SetLogLevel(level string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.LogLevel = level
	fmt.Printf("[CONFIG] Set log level to %s\n", level)
}

// GetDatabaseURL - геттер
func (c *AppConfig) GetDatabaseURL() string {
	return c.DatabaseURL
}

// Print - вывод текущего состояния
func (c *AppConfig) Print() {
	c.mu.Lock()
	defer c.mu.Unlock()
	fmt.Printf("Config State: {App: %s, Level: %s, Ptr: %p}\n",
		c.AppName, c.LogLevel, c)
}

func main() {
	fmt.Println("=== Singleton Pattern Demo (Thread-Safe) ===")

	// Тест 1: Конкурентный доступ
	fmt.Println("\n--- Concurrent Access Test ---")

	var wg sync.WaitGroup

	// Запускаем 10 горутин одновременно
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			cfg := GetInstance()
			// Просто для демо, чтобы убедиться что адрес один
			if id == 0 || id == 9 {
				fmt.Printf("Goroutine %d -> Config Ptr: %p\n", id, cfg)
			}
		}(i)
	}

	wg.Wait() // Ждем всех корректно
	fmt.Println("All goroutines finished.")

	// Тест 2: Проверка единственности и изменения данных
	fmt.Println("\n--- Data Consistency Test ---")

	cfg1 := GetInstance()
	cfg2 := GetInstance()

	fmt.Printf("cfg1 LogLevel before: %s\n", cfg1.LogLevel)

	// Меняем через первую ссылку
	cfg1.SetLogLevel("DEBUG")

	// Проверяем через вторую ссылку
	fmt.Printf("cfg2 LogLevel after:  %s\n", cfg2.LogLevel)

	if cfg1.LogLevel == cfg2.LogLevel {
		fmt.Println("SUCCESS: Changes reflected globally")
	} else {
		fmt.Println("FAIL: Instances are desynchronized")
	}

	cfg2.Print()
}
