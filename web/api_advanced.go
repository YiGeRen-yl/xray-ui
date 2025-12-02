package web

import (
    "encoding/json"
    "net/http"

    "github.com/gin-gonic/gin"

    "github.com/YiGeRen-yl/xray-ui/config"
)

type advancedConfigDTO struct {
    Outbounds json.RawMessage `json:"outbounds"`
    Routing   json.RawMessage `json:"routing"`
    DNS       json.RawMessage `json:"dns"`
}

func registerAdvancedConfigRoutes(r *gin.Engine) {
    api := r.Group("/api/advanced")

    api.GET("/config", func(c *gin.Context) {
        cfg, err := config.GetAdvancedConfig()
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"success": false, "msg": err.Error()})
            return
        }
        c.JSON(http.StatusOK, gin.H{
            "success": true,
            "data": gin.H{
                "outbounds": cfg.Outbounds,
                "routing":   cfg.Routing,
                "dns":       cfg.DNS,
            },
        })
    })

    api.POST("/config", func(c *gin.Context) {
        var body advancedConfigDTO
        if err := c.ShouldBindJSON(&body); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"success": false, "msg": err.Error()})
            return
        }

        cfg := &config.AdvancedConfig{
            Outbounds: body.Outbounds,
            Routing:   body.Routing,
            DNS:       body.DNS,
        }

        if err := config.SaveAdvancedConfig(cfg); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"success": false, "msg": err.Error()})
            return
        }

        c.JSON(http.StatusOK, gin.H{"success": true})
    })
}
