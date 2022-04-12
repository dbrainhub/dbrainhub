package filebeat

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/dbrainhub/dbrainhub/configs"
	"github.com/dbrainhub/dbrainhub/errors"
	"github.com/dbrainhub/dbrainhub/model"
	"github.com/dbrainhub/dbrainhub/utils"
	"github.com/dbrainhub/dbrainhub/utils/logger"
)

const (
	DBTypeMysql = "mysql"
	LocalHost   = "127.0.0.1"
)

type (
	FilebeatService interface {
		StartGatherSlowlog(ctx context.Context) error
	}
)

func NewFilebeatService(agentConf *configs.AgentConfig) (FilebeatService, error) {
	localip, err := utils.GetLocalIP()
	if err != nil {
		logger.Errorf("get local ip error, err: %v", err)
		return nil, err
	}
	return &filebeatService{
		dbInfo: &DBInfo{
			IP:     LocalHost,
			Port:   agentConf.DB.Port,
			User:   agentConf.DB.User,
			Passwd: agentConf.DB.Passwd,
		},
		localip:     localip,
		dbtype:      agentConf.DB.DBType,
		serverAddrs: []string{agentConf.Server.Addr},

		filebeatConfTemplateFile: agentConf.Filebeat.FilebeatConfTemplate,
		moduleConfTemplateFile:   agentConf.Filebeat.ModuleConfTemplate,

		moduleConfDir:     fmt.Sprintf("%s/modules.d", agentConf.Filebeat.HomePath),
		executionFilePath: fmt.Sprintf("%s/filebeat", agentConf.Filebeat.HomePath),
		homePath:          agentConf.Filebeat.HomePath,
		confFilePath:      fmt.Sprintf("%s/filebeat.yml", agentConf.Filebeat.HomePath),

		aliveListenerInterval:          time.Duration(agentConf.Filebeat.AliveListenerInterval) * time.Millisecond,
		aliveListenerHttpRetry:         agentConf.Filebeat.AliveListenerHttpRetry,
		aliveListenerHttpRetryInterval: time.Duration(agentConf.Filebeat.AliveListenerRetryInterval) * time.Millisecond,
		slowlogListenerInterval:        time.Duration(agentConf.Filebeat.SlowlogListenerInterval) * time.Millisecond,
		filebeatStartupTimeout:         time.Duration(agentConf.Filebeat.StartupTimeout) * time.Millisecond,
	}, nil
}

type filebeatService struct {
	filebeatConf *model.FilebeatConf

	localip string
	dbInfo  *DBInfo
	dbtype  string

	// conf template
	filebeatConfTemplateFile string
	moduleConfTemplateFile   string

	moduleConfDir string

	// dbrainhub server addrs
	serverAddrs []string

	// for filebeat startup
	executionFilePath string
	homePath          string
	confFilePath      string
	logPath           string
	dataPath          string

	// for filebeat alive listener
	aliveListenerInterval          time.Duration
	aliveListenerHttpRetry         int
	aliveListenerHttpRetryInterval time.Duration
	filebeatStartupTimeout         time.Duration

	// for slowlog listener
	slowlogListenerInterval time.Duration
}

// There are serveral steps below:
// 1. Generate filebeat conf.
// 2. Begin monitoring slow log. First success will trigger the generation of module conf.
// 3. When module conf generated, start filebeat and filebeat listener.
// 4. if filebeat's startup timeout, panic!
func (f *filebeatService) StartGatherSlowlog(ctx context.Context) error {
	if err := f.generateFilebeatConf(); err != nil {
		return err
	}

	slowLogUpdateFinishChan, err := f.startSlowlogListener(ctx)
	if err != nil {
		return err
	}

	f.startFilebeat(ctx, slowLogUpdateFinishChan)
	f.startAliveListener(ctx)
	return nil

}

func (f *filebeatService) startAliveListener(ctx context.Context) {
	firstSucc := false
	firstSuccChan := make(chan bool)
	errCallback := func(ctx context.Context, err error) {
		if !firstSucc {
			return
		}
		logger.Errorf("filebeat alive listener failed , err: %v", err)
		// TODO: callback server
	}
	succCallback := func(ctx context.Context, succInfo string) {
		firstSucc = true
		select {
		case firstSuccChan <- true:
		default:
		}
	}
	NewAliveListener(fmt.Sprintf("%s:%d", LocalHost, f.filebeatConf.HttpInfo.Port),
		utils.NewHttpClient(f.aliveListenerInterval, f.aliveListenerHttpRetry, f.aliveListenerHttpRetryInterval),
		f.aliveListenerInterval,
		&AliveListenerCallback{
			ErrorCallback: errCallback,
			SuccCallback:  succCallback,
		}).Listen(ctx)

	select {
	case <-firstSuccChan:
	case <-time.After(f.filebeatStartupTimeout):
		logger.Errorf("filebeat startup timeout")
		panic("filebeat startup timeout ...")
	}

}

func (f *filebeatService) startFilebeat(ctx context.Context, slowLogUpdateFinishChan <-chan bool) {
	// wait slow log updated
	<-slowLogUpdateFinishChan

	NewFilebeatOperation(&localExecutor{}, f.executionFilePath, f.confFilePath, f.homePath, f.logPath, f.dataPath).Startup(ctx)
}

func (f *filebeatService) startSlowlogListener(ctx context.Context) (<-chan bool, error) {
	db, err := f.dbInfo.GetDB(f.dbtype)
	if err != nil {
		return nil, err
	}
	slowLogQuery := model.NewSlowLogInfoQuerier(db)

	errCallback := func(ctx context.Context, err error) {
		logger.Errorf("slowlog listener err: %v", err)
		// TODO: callback the server

	}
	firstFinishChan := make(chan bool)
	changedCallback := func(ctx context.Context, oldPath, newPath string) {
		if err := f.generateModuleConf(newPath); err != nil {
			errCallback(ctx, err)
			return
		}
		select {
		case firstFinishChan <- true:
		default:
		}
	}
	NewSlowLogPathListener(slowLogQuery, f.slowlogListenerInterval, &SlowLogPathCallback{
		ChangedCallback: changedCallback,
		ErrorCallback:   errCallback,
	}).Listen(ctx)
	return firstFinishChan, nil

}

// template file -> filebeat.yml
func (f *filebeatService) generateFilebeatConf() error {
	filebeatConfTemplate, err := utils.ReadFile(f.filebeatConfTemplateFile)
	if err != nil {
		logger.Errorf("read filebeat conf file error, file: %s, err: %v", f.filebeatConfTemplateFile, err)
		return errors.AgentConfigError("read filebeat_conf error")
	}

	if err := NewConfValidator().ValidateFilebeatConf(filebeatConfTemplate); err != nil {
		logger.Errorf("filebeat conf validate err: %v", err)
		return err
	}

	filebeatConfStr := NewFilebeatConfGenerator(f.localip, f.dbInfo.Port, f.serverAddrs).Generate(filebeatConfTemplate)
	if err := utils.OverwriteToFile(f.confFilePath, filebeatConfStr); err != nil {
		logger.Errorf("write to filebeat conf file error, file: %s, err: %v", f.confFilePath, err)
		return err
	}
	f.filebeatConf, err = model.NewFileBeatConfFactory().NewFilebeatConf(filebeatConfStr)
	return err
}

// template file -> module conf file
func (f *filebeatService) generateModuleConf(inputPaths string) error {
	moduleConfTemplate, err := utils.ReadFile(f.moduleConfTemplateFile)
	if err != nil {
		logger.Errorf("read module conf file error, file: %s, err: %v", f.moduleConfTemplateFile, err)
		return errors.AgentConfigError("read module_conf error")
	}

	if err := NewConfValidator().ValidateModuleConf(moduleConfTemplate); err != nil {
		logger.Errorf("module conf validate err: %v", err)
		return err
	}

	moduleConfStr := NewModuleConfGenerator(inputPaths).Generate(moduleConfTemplate)
	if err := utils.OverwriteToFile(f.getModuleFilePath(), moduleConfStr); err != nil {
		logger.Errorf("write to module conf file error, file: %s, err: %v", f.getModuleFilePath(), err)
		return err
	}
	return nil
}

func (f *filebeatService) getModuleFilePath() string {
	return fmt.Sprintf("%s/%s.yml", f.moduleConfDir, f.dbtype)
}

type DBInfo struct {
	IP     string
	Port   int
	User   string
	Passwd string
}

func (d *DBInfo) GetDB(dbtype string) (*sql.DB, error) {
	switch dbtype {
	case DBTypeMysql:
		return sql.Open(dbtype,
			fmt.Sprintf("%s:%s@tcp(%s:%d)/", d.User, d.Passwd, d.IP, d.Port))
	}
	return nil, errors.AgentConfigError("invalid dbtype: %s", dbtype)
}
