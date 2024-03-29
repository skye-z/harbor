/*
宿主机监控

BetaX Harbor
Copyright © 2024 SkyeZhang <skai-zhang@hotmail.com>
*/

package monitor

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/skye-z/harbor/model"
	"github.com/skye-z/harbor/util"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"xorm.io/xorm"
)

type DeviceInfo struct {
	HostID      string `json:"id"`
	Type        string `json:"type"`
	Platform    string `json:"platform"`
	Family      string `json:"family"`
	Version     string `json:"version"`
	Arch        string `json:"arch"`
	BootTime    uint64 `json:"bootTime"`
	UpTime      uint64 `json:"upTime"`
	CpuName     string `json:"cpuName"`
	CpuPhysical int    `json:"cpuPhysical"`
	CpuLogical  int    `json:"cpuLogical"`
	Swap        uint64 `json:"swap"`
	Memory      uint64 `json:"memory"`
	Disk        uint64 `json:"disk"`
}

func GetDeviceInfo() DeviceInfo {
	info, _ := host.Info()

	physicalCnt, _ := cpu.Counts(false)
	logicalCnt, _ := cpu.Counts(true)
	cpus, _ := cpu.Info()
	cpuName := strings.Replace(cpus[0].ModelName, " CPU", "", -1)
	cpuName = strings.Replace(cpuName, " Processor", "", -1)
	cpuName = strings.Replace(cpuName, " @", "", -1)
	cpuName = strings.Replace(cpuName, "(R)", "", -1)
	cpuName = strings.Replace(cpuName, "(TM)", "", -1)
	arch := info.KernelArch
	arch = strings.Replace(arch, "_", " ", -1)

	v, _ := mem.VirtualMemory()
	s, _ := mem.SwapMemory()
	virtualTotal := v.Total / 1024 / 1024
	swapTotal := s.Total / 1024 / 1024

	var diskTotal uint64 = 0
	// 获取根目录使用情况
	diskUsage, err := disk.Usage("/")
	if err != nil {
		// 获取磁盘的分区信息
		diskInfo, _ := disk.Partitions(true)
		if diskInfo[0].Mountpoint == "/" {
			usage, _ := disk.Usage("/")
			diskTotal += usage.Total
		} else {
			for i := 0; i < len(diskInfo); i++ {
				usage, _ := disk.Usage(diskInfo[i].Mountpoint)
				diskTotal += usage.Total
			}
		}
		diskTotal = diskTotal / 1024 / 1024

	} else {
		diskTotal = diskUsage.Total / 1024 / 1024
	}

	return DeviceInfo{
		HostID:      info.HostID,
		Type:        info.OS,
		Platform:    info.Platform,
		Family:      info.PlatformFamily,
		Version:     info.KernelVersion,
		Arch:        arch,
		BootTime:    info.BootTime,
		UpTime:      info.Uptime,
		CpuName:     cpuName,
		CpuPhysical: physicalCnt,
		CpuLogical:  logicalCnt,
		Swap:        swapTotal,
		Memory:      virtualTotal,
		Disk:        diskTotal,
	}
}

type SystemUse struct {
	CPU    []float64              `json:"cpu"`
	Avg    *load.AvgStat          `json:"avg"`
	Memory *mem.VirtualMemoryStat `json:"memory"`
	Disk   *disk.UsageStat        `json:"disk"`
}

func GetUse() SystemUse {
	cpuPercent, _ := cpu.Percent(0, false)
	avg, _ := load.Avg()
	memInfo, _ := mem.VirtualMemory()
	diskUsage, _ := disk.Usage("/")
	return SystemUse{
		CPU:    cpuPercent,
		Avg:    avg,
		Memory: memInfo,
		Disk:   diskUsage,
	}
}

func ListenHostOverhead(engine *xorm.Engine) {
	logger := &model.LogModel{
		DB: engine,
	}
	noticeNumber := 0
	init := true
	for {
		if init {
			log.Println("[Monitor] start listening events for host")
			init = false
		} else {
			interval := util.GetInt("alarm.interval") * 60
			loadThreshold := util.GetFloat64("alarm.loadThreshold")
			memoryThreshold := util.GetFloat64("alarm.memoryThreshold")
			diskThreshold := util.GetFloat64("alarm.diskThreshold")

			use := GetUse()
			// 通知计数为0才检查
			if noticeNumber == 0 {
				if use.Avg.Load1 > loadThreshold || use.Avg.Load5 > loadThreshold || use.Avg.Load15 > loadThreshold {
					log.Println("[Monitor] system load high")
					logger.AddLog("host", "highLoad", fmt.Sprintf("负载: %.2f %.2f %.2f", use.Avg.Load1, use.Avg.Load5, use.Avg.Load15))
					if NoticeHighLoad() {
						noticeNumber += 1
					}
				}
				if use.Memory.UsedPercent > memoryThreshold {
					log.Println("[Monitor] memory run out")
					num := strconv.FormatFloat(use.Memory.UsedPercent, 'f', 2, 64)
					logger.AddLog("host", "runOut", "内存已使用"+num+"%")
					if NoticeRunOut("内存", num) {
						noticeNumber += 1
					}
				}
				if use.Disk.UsedPercent > diskThreshold {
					log.Println("[Monitor] disk run out")
					num := strconv.FormatFloat(use.Disk.UsedPercent, 'f', 2, 64)
					logger.AddLog("host", "runOut", "磁盘已使用"+num+"%")
					if NoticeRunOut("磁盘", num) {
						noticeNumber += 1
					}
				}
				// 每30分钟告警1次
				if noticeNumber > 0 {
					noticeNumber = interval
				}
			} else if noticeNumber > 0 {
				noticeNumber--
			} else if noticeNumber < 0 {
				noticeNumber = 0
			}
			// 记录日志

		}
		time.Sleep(time.Minute)
	}
}
