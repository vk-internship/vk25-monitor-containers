package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/joho/godotenv"
	proBing "github.com/prometheus-community/pro-bing"
	"github.com/robfig/cron/v3"
)

type Config struct {
	BackendAPIURL string
	Interval      time.Duration
}

type PingResult struct {
	IPAddress       string     `json:"ip_address"`
	PingTime        time.Time  `json:"ping_time"`
	IsSuccess       bool       `json:"is_success"`
	LastSuccessTime *time.Time `json:"last_success_time,omitempty"`
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Ошибка загрузки .env:", err)
	}

	config := Config{
		BackendAPIURL: os.Getenv("BACKEND_API_URL"),
		Interval:      time.Minute,
	}

	if intervalStr := os.Getenv("INTERVAL"); intervalStr != "" {
		if interval, err := time.ParseDuration(intervalStr); err == nil {
			config.Interval = interval
		} else {
			log.Println("Некорректный интервал, используем значение по умолчанию:", err)
		}
	}

	c := cron.New()
	_, err := c.AddFunc(fmt.Sprintf("@every %s", config.Interval), func() {
		log.Println("Задача Cron запущена")

		if err := processContainers(config); err != nil {
			log.Println("Ошибка:", err)
		}
	})

	if err != nil {
		log.Fatal("Ошибка добавления задачи в Cron:", err)
	}

	c.Start()
	log.Printf("Pinger запущен. Интервал: %s. Ожидание задач Cron...", config.Interval)

	select {}
}

func processContainers(config Config) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return fmt.Errorf("ошибка подключения к Docker: %v", err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return fmt.Errorf("ошибка получения контейнеров: %v", err)
	}

	for _, container := range containers {
		var ip string
		for _, network := range container.NetworkSettings.Networks {
			ip = network.IPAddress
			break
		}

		if ip == "" {
			log.Printf("Контейнер %s не имеет IP-адреса", container.ID)
			continue
		}

		isAlive, err := pingIP(ip)
		if err != nil {
			log.Printf("Ошибка пинга %s: %v", ip, err)
			continue
		}

		result := PingResult{
			IPAddress: ip,
			PingTime:  time.Now(),
			IsSuccess: isAlive,
		}

		if isAlive {
			now := time.Now()
			result.LastSuccessTime = &now
		}

		if err := sendToBackend(config.BackendAPIURL, result); err != nil {
			log.Printf("Ошибка отправки данных: %v", err)
		}
	}

	return nil
}

func pingIP(ip string) (bool, error) {
	pinger, err := proBing.NewPinger(ip)
	if err != nil {
		return false, err
	}

	pinger.Count = 1
	pinger.Timeout = 5 * time.Second
	pinger.SetPrivileged(true)

	if err := pinger.Run(); err != nil {
		log.Printf("Ошибка выполнения пинга для %s: %v", ip, err)
		return false, err
	}

	stats := pinger.Statistics()
	return stats.PacketsRecv > 0, nil
}

func sendToBackend(url string, data PingResult) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	resp, err := http.Post(url+"/pings", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("неверный статус код: %d", resp.StatusCode)
	}

	return nil
}
