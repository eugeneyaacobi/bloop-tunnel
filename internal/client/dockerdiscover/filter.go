package dockerdiscover

import "strings"

func looksHTTPish(port int) bool {
	switch port {
	case 80, 81, 3000, 4000, 4173, 5000, 5173, 5432, 5601, 5678, 6006, 8000, 8080, 8081, 8088, 8888, 9000, 9080, 10000:
		return port != 5432
	default:
		return port >= 1024 && port <= 10000
	}
}


func excluded(name, image string) bool {
	target := strings.ToLower(name + " " + image)
	for _, term := range []string{"postgres", "mysql", "mariadb", "redis", "rabbitmq", "kafka", "zookeeper", "mongo", "minio", "nats", "broker", "queue", "db"} {
		if strings.Contains(target, term) {
			return true
		}
	}
	return false
}
