package config

//Conf 配置结构体
type Conf struct {
	Listen       string
	DataDir      string
	DatabaseName string
	TablePrefix  string
	AppName      string
	Secret       string
}

//V current config value
var V *Conf

func Setup(listen, dataDir, databaseName, tablePrefix string) error {
	V = &Conf{
		Listen:       listen,
		DataDir:      dataDir,
		DatabaseName: databaseName,
		TablePrefix:  tablePrefix,
	}
	return nil
}
