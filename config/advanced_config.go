package config

import (
    "encoding/json"
    "errors"
    "io/ioutil"
    "os"
    "path/filepath"
    "sync"
)

type AdvancedConfig struct {
    Outbounds json.RawMessage `json:"outbounds"` // JSON 数组
    Routing   json.RawMessage `json:"routing"`   // JSON 对象
    DNS       json.RawMessage `json:"dns"`       // JSON 对象
}

var (
    advancedCfg     *AdvancedConfig
    advancedCfgOnce sync.Once
    advancedCfgLock sync.Mutex
)

func getAdvancedConfigPath() string {
    // 如果你项目有统一的配置目录，可以改成对应目录
    return filepath.Join("/etc/xray-ui", "advanced.json")
}

func loadAdvancedConfigFromFile() (*AdvancedConfig, error) {
    path := getAdvancedConfigPath()
    _, err := os.Stat(path)
    if errors.Is(err, os.ErrNotExist) {
        // 没有就返回空配置
        return &AdvancedConfig{}, nil
    }
    if err != nil {
        return nil, err
    }

    data, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, err
    }

    cfg := &AdvancedConfig{}
    if len(data) == 0 {
        return cfg, nil
    }

    if err := json.Unmarshal(data, cfg); err != nil {
        return nil, err
    }
    return cfg, nil
}

func GetAdvancedConfig() (*AdvancedConfig, error) {
    var err error
    advancedCfgOnce.Do(func() {
        advancedCfg, err = loadAdvancedConfigFromFile()
    })
    return advancedCfg, err
}

func SaveAdvancedConfig(cfg *AdvancedConfig) error {
    if cfg == nil {
        cfg = &AdvancedConfig{}
    }

    advancedCfgLock.Lock()
    defer advancedCfgLock.Unlock()

    path := getAdvancedConfigPath()
    if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
        return err
    }

    data, err := json.MarshalIndent(cfg, "", "  ")
    if err != nil {
        return err
    }

    if err := ioutil.WriteFile(path, data, 0o600); err != nil {
        return err
    }

    advancedCfg = cfg
    return nil
}

