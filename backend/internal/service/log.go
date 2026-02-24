package service

import (
	"fmt"
	stdlog "log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/skye-z/harbor/internal/data"
	"github.com/skye-z/harbor/internal/util/response"
	"xorm.io/xorm"
)

// 日志服务
type LogService struct {
	engine *xorm.Engine
}

var logChannel chan *data.SystemLog
var logBufferSize = 100
var logEngine *xorm.Engine

func init() {
	logChannel = make(chan *data.SystemLog, logBufferSize)
}

func SetLogEngine(engine *xorm.Engine) {
	logEngine = engine
	go func() {
		for logEntry := range logChannel {
			retries := 0
			maxRetries := 5

			for retries < maxRetries {
				_, err := logEngine.Insert(logEntry)
				if err == nil {
					stdlog.Printf("[Log] Inserted log: Type=%s, Action=%s, Target=%s", logEntry.Type, logEntry.Action, logEntry.Target)
					break
				}

				retries++
				if retries < maxRetries {
					stdlog.Printf("[Log] Database busy, retry %d/%d: %v", retries, maxRetries, err)
					time.Sleep(time.Duration(retries) * 100 * time.Millisecond)
				} else {
					stdlog.Printf("[Log] Failed to insert log after %d attempts: %v", maxRetries, err)
				}
			}
		}
	}()
}

// 创建日志服务实例
func NewLogService(engine *xorm.Engine) *LogService {
	if logEngine == nil {
		SetLogEngine(engine)
	}
	return &LogService{engine: engine}
}

// 日志类型常量
const (
	LogTypeSystem    = "system"
	LogTypeContainer = "container"
	LogTypeImage     = "image"
	LogTypeNetwork   = "network"
	LogTypeVolume    = "volume"
)

// 日志级别常量
const (
	LogLevelInfo    = "info"
	LogLevelWarning = "warning"
	LogLevelError   = "error"
)

// 记录系统启动日志
func (s *LogService) LogSystemStartup() error {
	log := &data.SystemLog{
		Type:      LogTypeSystem,
		Level:     LogLevelInfo,
		Action:    "startup",
		Target:    "system",
		Message:   "Harbor Docker 管理面板已启动",
		CreatedAt: time.Now(),
	}
	select {
	case logChannel <- log:
	default:
	}
	return nil
}

// 记录系统关闭日志
func (s *LogService) LogSystemShutdown() error {
	log := &data.SystemLog{
		Type:      LogTypeSystem,
		Level:     LogLevelInfo,
		Action:    "shutdown",
		Target:    "system",
		Message:   "Harbor Docker 管理面板已关闭",
		CreatedAt: time.Now(),
	}
	select {
	case logChannel <- log:
	default:
	}
	return nil
}

// 记录容器操作日志
func (s *LogService) LogContainerAction(action string, containerName string, containerID string, username string, userID int, level string) error {
	message := fmt.Sprintf("容器操作: %s - %s (%s)", action, containerName, containerID)
	log := &data.SystemLog{
		Type:      LogTypeContainer,
		Level:     level,
		Action:    action,
		Target:    containerName,
		TargetID:  containerID,
		Message:   message,
		Username:  username,
		UserID:    userID,
		CreatedAt: time.Now(),
	}
	select {
	case logChannel <- log:
	default:
	}
	return nil
}

// 记录镜像操作日志
func (s *LogService) LogImageAction(action string, imageName string, username string, userID int, level string) error {
	message := fmt.Sprintf("镜像操作: %s - %s", action, imageName)
	log := &data.SystemLog{
		Type:      LogTypeImage,
		Level:     level,
		Action:    action,
		Target:    imageName,
		Message:   message,
		Username:  username,
		UserID:    userID,
		CreatedAt: time.Now(),
	}
	select {
	case logChannel <- log:
	default:
	}
	return nil
}

// 记录网络操作日志
func (s *LogService) LogNetworkAction(action string, networkName string, networkID string, username string, userID int, level string) error {
	message := fmt.Sprintf("网络操作: %s - %s (%s)", action, networkName, networkID)
	log := &data.SystemLog{
		Type:      LogTypeNetwork,
		Level:     level,
		Action:    action,
		Target:    networkName,
		TargetID:  networkID,
		Message:   message,
		Username:  username,
		UserID:    userID,
		CreatedAt: time.Now(),
	}
	select {
	case logChannel <- log:
	default:
	}
	return nil
}

// 记录卷操作日志
func (s *LogService) LogVolumeAction(action string, volumeName string, username string, userID int, level string) error {
	message := fmt.Sprintf("卷操作: %s - %s", action, volumeName)
	log := &data.SystemLog{
		Type:      LogTypeVolume,
		Level:     level,
		Action:    action,
		Target:    volumeName,
		Message:   message,
		Username:  username,
		UserID:    userID,
		CreatedAt: time.Now(),
	}
	select {
	case logChannel <- log:
	default:
	}
	return nil
}

// 记录通用日志
func (s *LogService) Log(logType string, level string, action string, target string, targetID string, message string, username string, userID int) error {
	log := &data.SystemLog{
		Type:      logType,
		Level:     level,
		Action:    action,
		Target:    target,
		TargetID:  targetID,
		Message:   message,
		Username:  username,
		UserID:    userID,
		CreatedAt: time.Now(),
	}
	select {
	case logChannel <- log:
	default:
	}
	return nil
}

func (s *LogService) LogLogin(username string, userID int, ipAddress string) error {
	if logEngine == nil {
		return fmt.Errorf("log engine not initialized")
	}

	message := fmt.Sprintf("用户登录: %s (ID: %d)", username, userID)
	logEntry := &data.SystemLog{
		Type:      LogTypeSystem,
		Level:     LogLevelInfo,
		Action:    "login",
		Target:    "user",
		TargetID:  fmt.Sprintf("%d", userID),
		Message:   message,
		Username:  username,
		UserID:    userID,
		IPAddress: ipAddress,
		CreatedAt: time.Now(),
	}
	select {
	case logChannel <- logEntry:
		return nil
	default:
		return fmt.Errorf("log channel is full")
	}
}

func (s *LogService) LogContainerEvent(eventType string, containerName string, containerID string, status string) error {
	if logEngine == nil {
		return fmt.Errorf("log engine not initialized")
	}

	message := fmt.Sprintf("容器事件: %s - %s (%s)", eventType, containerName, status)
	logEntry := &data.SystemLog{
		Type:      LogTypeContainer,
		Level:     LogLevelInfo,
		Action:    eventType,
		Target:    containerName,
		TargetID:  containerID,
		Message:   message,
		CreatedAt: time.Now(),
	}
	select {
	case logChannel <- logEntry:
		return nil
	default:
		return fmt.Errorf("log channel is full")
	}
}

// 获取日志列表
func (s *LogService) GetList(logType string, page int, pageSize int) ([]*data.SystemLog, int64, error) {
	var logs []*data.SystemLog
	var total int64

	var query *xorm.Session

	if logType != "" {
		query = s.engine.Where("type = ?", logType)
	} else {
		query = s.engine.Table(&data.SystemLog{})
	}

	// 统计总数
	count, err := query.Count()
	if err != nil {
		return nil, 0, err
	}
	total = count

	// 分页查询
	offset := (page - 1) * pageSize
	err = query.OrderBy("created_at DESC").Limit(pageSize, offset).Find(&logs)
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// GetRecent HTTP handler for recent logs
func (s *LogService) GetRecent(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	logType := c.Query("type")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	logs, _, err := s.GetList(logType, 1, limit)
	if err != nil {
		response.Error(c, err.Error())
		return
	}

	response.Success(c, logs)
}
