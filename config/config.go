package config

import (
	"context"
	"path"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type ConfigureReader func(ctx context.Context) (func(), error)

var (
	mutex     sync.Mutex
	readerMap map[uintptr]ConfigureReader
)

func LoadConfigWithFile(filename string, paths ...string) error {
	index := strings.LastIndex(filename, ".")
	viper.SetConfigName(filename[:index])
	viper.SetConfigType(filename[index+1:])
	for _, p := range paths {
		viper.AddConfigPath(p)
	}
	return viper.ReadInConfig()
}

func LoadConfigWithFlag(flagName string) error {
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		return err
	}
	filepath := viper.GetString(flagName)
	dir, file := path.Split(filepath)
	return LoadConfigWithFile(file, dir)
}

func configureReaderUintptr(fn ConfigureReader) uintptr {
	return reflect.ValueOf(fn).Pointer()
}

func RegisterConfigureReaders(fns ...ConfigureReader) {
	mutex.Lock()
	defer mutex.Unlock()
	if readerMap == nil {
		readerMap = make(map[uintptr]ConfigureReader)
	}
	for i := 0; i < len(fns); i++ {
		fn := fns[i]
		fnUintptr := configureReaderUintptr(fn)
		if _, ok := readerMap[fnUintptr]; ok {
			continue
		}
		readerMap[fnUintptr] = fn
	}
}

func DoReaderConfiguration(ctx context.Context) (func(), error) {
	mutex.Lock()
	defer mutex.Unlock()
	if readerMap == nil || len(readerMap) == 0 {
		return func() {}, nil
	}
	cleans := make([]func(), 0)
	for _, reader := range readerMap {
		clean, err := reader(ctx)
		if err != nil {
			if clean != nil {
				clean()
			}
			for _, fn := range cleans {
				fn()
			}
			return nil, err
		}
		if clean != nil {
			cleans = append(cleans, clean)
		}
	}
	return func() {
		for _, fn := range cleans {
			fn()
		}
	}, nil
}

func ViperGet[T any](key string, defaults ...T) T {
	var zero T
	// 判断key是否存在
	if !viper.IsSet(key) {
		return defaultOrZero(defaults)
	}
	var result any
	switch any(zero).(type) {
	case uint:
		result = viper.GetUint(key)
	case uint64:
		result = viper.GetUint64(key)
	case int:
		result = viper.GetInt(key)
	case int64:
		result = viper.GetInt64(key)
	case float64:
		result = viper.GetFloat64(key)
	case bool:
		result = viper.GetBool(key)
	case string:
		result = viper.GetString(key)
	case time.Duration:
		result = viper.GetDuration(key)
	case []string:
		result = viper.GetStringSlice(key)
	case []int:
		result = viper.GetIntSlice(key)
	default:
		result = viper.Get(key)
	}
	// 尝试类型转换
	val, ok := result.(T)
	if !ok {
		// 尝试使用反射进行转换
		rv := reflect.ValueOf(result)
		if rv.Type().ConvertibleTo(reflect.TypeOf(zero)) {
			return rv.Convert(reflect.TypeOf(zero)).Interface().(T)
		}
		return defaultOrZero(defaults)
	}
	return val
}

func defaultOrZero[T any](defaults []T) T {
	var zero T
	if len(defaults) != 0 {
		return defaults[0]
	}
	return zero
}
