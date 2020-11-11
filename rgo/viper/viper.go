package viper

import(
	"log"
	"os"
	"github.com/spf13/viper"
)

//实例化config,这次是指定目录了
func NewConfig(cnfName,cnfType string)*viper.Viper{
	viper := viper.New()
	viper.SetConfigName(cnfName)
	viper.SetConfigType(cnfType)

	env := "dev"
	runmod := os.Getenv("RGO_RUNMOD")
	if runmod == "master"{
		env = "master"
	}
	if runmod == "test"{
		env = "test"
	}
	cnfDirPath := "./config/"+env
	
	viper.AddConfigPath(cnfDirPath)
	if err := viper.ReadInConfig();err != nil {
		log.Println(err)
	}
	return viper
}