package internal

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	g "gvb/internal/global"
	"gvb/internal/model"
	"os"
)

var Lg *zap.Logger

func InitLogger(conf *g.Config) {
	writeSyncerDebug := getLogWriter(conf.Log.DebugFileName, conf.Log.MaxSize, conf.Log.MaxBackups, conf.Log.MaxAge)
	writeSyncerInfo := getLogWriter(conf.Log.InfoFileName, conf.Log.MaxSize, conf.Log.MaxBackups, conf.Log.MaxAge)
	writeSyncerWarn := getLogWriter(conf.Log.WarnFileName, conf.Log.MaxSize, conf.Log.MaxBackups, conf.Log.MaxAge)
	encoder := getEncoder()
	//文件输出
	debugCore := zapcore.NewCore(encoder, writeSyncerDebug, zapcore.DebugLevel)
	infoCore := zapcore.NewCore(encoder, writeSyncerInfo, zapcore.InfoLevel)
	warnCore := zapcore.NewCore(encoder, writeSyncerWarn, zapcore.WarnLevel)
	//标准输出
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	std := zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel)
	core := zapcore.NewTee(debugCore, infoCore, warnCore, std)
	Lg = zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(Lg) // 替换zap包中全局的logger实例，后续在其他包中只需使用zap.L()调用即可
	return
}
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
	}
	return zapcore.AddSync(lumberJackLogger)
}

// 数据库初始化
func InitDatabase(conf *g.Config) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.Mysql.Username,
		conf.Mysql.Password,
		conf.Mysql.Host,
		conf.Mysql.Port,
		conf.Mysql.Dbname,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		zap.L().Debug("数据库连接失败...")
	}
	zap.L().Info("数据库 连接成功...")
	// 数据库迁移
	if conf.Server.DbAutoMigrate {
		if err = model.MakeMigrate(db); err != nil {
			zap.L().Debug("数据库迁移失败")
		}
	}
	return db
}
func InitRedis(conf *g.Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Addr,
		Password: conf.Redis.Password,
		DB:       conf.Redis.DB,
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		zap.L().Debug("redis 连接失败！！！")
	}
	zap.L().Info("redis 连接成功...")
	return rdb
}
