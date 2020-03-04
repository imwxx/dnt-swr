package lib
import (
    "gopkg.in/yaml.v2"
    "io/ioutil"
    //"fmt"
    //"reflect"
)

var AppConf CONFIG

type CONFIG struct {
    VERFKEY string `yaml:"VerfKey"`
    DBINFO string `yaml:"DbInfo"`
    REDISHOST string `yaml:"RedisHost"`
    REDISPORT string `yaml:"RedisPort"`
    HWDOMAIN string `yaml:"HWDomain"`
    BUILDER  []string `yaml:"Builder,flow"`
    KSNSMAP  []string `yaml:"KSNSMap,flow"`
    NAMESPACEMAP []SNSMAP `yaml:"NamespaceMap"`
}

type SNSMAP struct {
    HARBORREPO string `yaml:"HarborRepo"`
    SWRREPO string `yaml:"SwrRepo"`
}

func LoadConfig(confFile string) (CONFIG, error){
    var dntconf CONFIG
    config, err := ioutil.ReadFile(confFile)
    if err != nil {
        return dntconf, err
    }
    yaml.Unmarshal([]byte(config), &dntconf)
    return dntconf, nil
}
