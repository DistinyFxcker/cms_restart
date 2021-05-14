package runtime

import (
	"errors"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"log"
	"textprojcet/domain"
	"textprojcet/schema/mysql"
	"textprojcet/schema/sqlite"
	"upper.io/db.v3/lib/sqlbuilder"
	mysqlDbV2 "upper.io/db.v3/mysql"
	sqliteDbV2 "upper.io/db.v3/sqlite"
)
type Runtime interface {
	//Init 初始化运行
	Init() error

	//Start 开启运行时
	Start() error

	//Stop 停止运行时
	Stop()
}
func New(database sqlbuilder.Database, dbType string , fileDb *badger.DB , prot int , apiPort int , p2pPort int) Runtime{

}

type runtime struct {
	database sqlbuilder.Database
	dbType string
	fileDb *badger.DB
	port int
	apiport int
	p2pport int

	services []domain.Service
}

func (rt *runtime) Init() error {
	rt.ctx = domain.NewContext(rt.database,rt.dbType)

	if err := func() error(){
		dbDriver := rt.database.ConnectionURL()
		if _, ok :=dbDriver.(sqliteDbV2.ConnectionURL);ok{//链接数据库，
			if err := sqlite.Migrate(rt.ctx);err != nil {//执行数据库迁移
				return err
			}
		}
		if _, ok := dbDriver.(mysqlDbV2.ConnectionURL);ok{
			if err := mysql.Migrate(rt.ctx);err != nil{
				return err
			}
		}
		return errors.New("不支持的数据库类型，仅支持Sqlite和mysql")
	}(); err !=nil{
		return err
	}
}

func (rt *runtime) register(srv domain.Service) error {
	log.Println("注册", srv.Name())

	for i := range rt.services{
		if srv.Name() == rt.services[i].Name(){
			return fmt.Errorf("服务:%s已经存在",srv.Name())
		}
	}
	rt.services = append(rt.services , srv)
	return nil
}