package goconf

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"context"
	"net"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/ssh"
)

//MysqlConf 配置
type MysqlConf struct {
	DbURL    string
	LifeTime int
	IdleConn int
	OpenConn int
}

//NewMysqlConf mysql配置
func NewMysqlConf(key string) *MysqlConf {
	sec := nowConfig.Section("mysql")
	dbURL := sec.Key(key + "_mysql_url").MustString("")
	dbUser := sec.Key(key + "_mysql_user").MustString("")
	dbPw := sec.Key(key + "_mysql_pwd").MustString("")
	dbName := sec.Key(key + "_mysql_database").MustString("")
	dbSubURL := dbUser + ":" + dbPw + "@tcp(" + dbURL + ")/" + dbName + "?charset=utf8&parseTime=True&loc=Local"
	LifeTime := sec.Key(key + "_mysql_conn_life_time").MustInt(300)
	IdleConn := sec.Key(key + "_mysql_idle_conn").MustInt(1)
	OpenConn := sec.Key(key + "_mysql_open_conn").MustInt(10)
	return &MysqlConf{
		DbURL:    dbSubURL,
		LifeTime: LifeTime,
		IdleConn: IdleConn, OpenConn: OpenConn,
	}
}

type sshConfig struct {
	Host     string
	User     string
	Password string
	Use      int
}

//GetSSHConfig 获取ssh 配置
func GetSSHConfig() *sshConfig {
	sec := nowConfig.Section("ssh")
	Host := sec.Key("host").MustString("")
	User := sec.Key("user").MustString("")
	Password := sec.Key("password").MustString("")
	Use := sec.Key("use").MustInt(0)
	return &sshConfig{Host: Host, User: User, Password: Password, Use: Use}
}

//InitMysql 初始化mysql
func InitMysql(key string) *sql.DB {
	conf := NewMysqlConf(key)
	dbSubURL := conf.DbURL
	sqlDb, err := sql.Open("mysql", dbSubURL)
	if err != nil {
		msg := fmt.Sprintf("error:%v", err)
		fmt.Println(msg)
		panic(err)
	}
	sqlDb.SetMaxIdleConns(conf.IdleConn)
	sqlDb.SetMaxOpenConns(conf.OpenConn)
	sqlDb.SetConnMaxLifetime(time.Second * time.Duration(conf.LifeTime))
	return sqlDb
}

//InitGorm 初始化
func InitGorm(key string) *gorm.DB {
	conf := NewMysqlConf(key)
	dbSubURL := conf.DbURL
	db, err := gorm.Open("mysql", dbSubURL)
	if err != nil {
		panic(err)
	}
	db.DB().SetMaxIdleConns(conf.IdleConn)
	db.DB().SetMaxOpenConns(conf.OpenConn)
	db.DB().SetConnMaxLifetime(time.Second * time.Duration(conf.LifeTime))
	return db
}

//GetInsertSQL 获取sql
func GetInsertSQL(table string, sets map[string]interface{}) string {
	keys := []string{}
	values := []string{}
	for idx := range sets {
		value := sets[idx]
		switch value.(type) {
		case string:
			keys = append(keys, idx)
			values = append(values, "'"+value.(string)+"'")
		case int, int16, int32, int64, int8, uint, uint16, uint32, uint64, float32, float64:
			keys = append(keys, idx)
			var x string
			switch value.(type) {
			case int:
				x = strconv.FormatInt(int64(value.(int)), 10)
			case int16:
				x = strconv.FormatInt(int64(value.(int16)), 10)
			case int32:
				x = strconv.FormatInt(int64(value.(int32)), 10)
			case int64:
				x = strconv.FormatInt(int64(value.(int64)), 10)
			case int8:
				x = strconv.FormatInt(int64(value.(int8)), 10)
			case uint:
				x = strconv.FormatUint(uint64(value.(uint)), 10)
			case uint16:
				x = strconv.FormatUint(uint64(value.(uint16)), 10)
			case uint32:
				x = strconv.FormatUint(uint64(value.(uint32)), 10)
			case uint64:
				x = strconv.FormatUint(uint64(value.(uint64)), 10)
			case float32:
				x = strconv.FormatFloat(float64(value.(float32)), 'f', -6, 32)
			case float64:
				x = strconv.FormatFloat(float64(value.(float64)), 'f', -6, 64)
			}
			values = append(values, x)
		default:
			panic("GetInsertSQL type error")
		}
	}
	fmt.Println(keys)
	fmt.Println(values)
	sql := fmt.Sprintf("insert into %v (%v) values (%v)", table, strings.Join(keys, ","), strings.Join(values, ","))
	return sql
}

//ExecuteSQL sql转换
func ExecuteSQL(table string, sets map[string]interface{}, db *sql.DB) (int64, error) {
	keys := []string{}
	values := []string{}
	objs := []interface{}{}
	for idx := range sets {
		value := sets[idx]
		switch value.(type) {
		case string:
			keys = append(keys, idx)
			values = append(values, "?")
		case int, int16, int32, int64, int8, uint, uint16, uint32, uint64, float32, float64:
			keys = append(keys, idx)
			values = append(values, "?")
		default:
			panic("GetInsertSQL type error")
		}
		objs = append(objs, value)
	}
	sqlx := fmt.Sprintf("insert into %v (%v) values (%v)", table, strings.Join(keys, ","), strings.Join(values, ","))
	return Execute(table, sqlx, objs, db)
}

//Execute sql 执行
func Execute(table string, sql string, values []interface{}, db *sql.DB) (int64, error) {
	stmt, err := db.Prepare(sql)
	if err != nil {
		fmt.Println(err, "Prepare")
		return 0, err
	}
	result, err := stmt.Exec(values...)
	if err != nil {
		fmt.Println(err, "Exec")
		return 0, err
	}
	return result.LastInsertId()
}

//FormatUpdateSQL 整理uppdate的sql语句
func FormatUpdateSQL(sets map[string]interface{}, where map[string]interface{}) {

}

type ViaSSHDialer struct {
	client *ssh.Client
	_      *context.Context
}

func NewViaSSHDialer(client *ssh.Client) *ViaSSHDialer {
	return &ViaSSHDialer{client, nil}
}

func (self *ViaSSHDialer) Dial(context context.Context, addr string) (net.Conn, error) {
	return self.client.Dial("tcp", addr)
}

type remoteScriptType byte
type remoteShellType byte

const (
	cmdLine remoteScriptType = iota
	rawScript
	scriptFile

	interactiveShell remoteShellType = iota
	nonInteractiveShell
)

type Client struct {
	client *ssh.Client
}

func RegisterDialContext() {
	getSSHConfig := GetSSHConfig()
	client, err := DialWithPasswd(getSSHConfig.Host, getSSHConfig.User, getSSHConfig.Password)
	if err != nil {
		panic(err)
	}
	out, err := client.Cmd("ls -l").Output()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(out))
	// Now we register the ViaSSHDialer with the ssh connection as a parameter
	mysql.RegisterDialContext("mysql+tcp", (NewViaSSHDialer(client.client).Dial))
}

// DialWithPasswd starts a client connection to the given SSH server with passwd authmethod.
func DialWithPasswd(addr, user, passwd string) (*Client, error) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(passwd),
		},
		HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }),
	}

	return Dial("tcp", addr, config)
}

// DialWithKey starts a client connection to the given SSH server with key authmethod.
func DialWithKey(addr, user, keyfile string) (*Client, error) {
	key, err := ioutil.ReadFile(keyfile)
	if err != nil {
		return nil, err
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }),
	}

	return Dial("tcp", addr, config)
}

// DialWithKeyWithPassphrase same as DialWithKey but with a passphrase to decrypt the private key
func DialWithKeyWithPassphrase(addr, user, keyfile string, passphrase string) (*Client, error) {
	key, err := ioutil.ReadFile(keyfile)
	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKeyWithPassphrase(key, []byte(passphrase))
	if err != nil {
		return nil, err
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }),
	}

	return Dial("tcp", addr, config)
}

// Dial starts a client connection to the given SSH server.
// This is wrap the ssh.Dial
func Dial(network, addr string, config *ssh.ClientConfig) (*Client, error) {
	client, err := ssh.Dial(network, addr, config)
	if err != nil {
		return nil, err
	}
	return &Client{
		client: client,
	}, nil
}

func (c *Client) Close() error {
	return c.client.Close()
}

// Cmd create a command on client
func (c *Client) Cmd(cmd string) *remoteScript {
	return &remoteScript{
		_type:  cmdLine,
		client: c.client,
		script: bytes.NewBufferString(cmd + "\n"),
	}
}

// Script
func (c *Client) Script(script string) *remoteScript {
	return &remoteScript{
		_type:  rawScript,
		client: c.client,
		script: bytes.NewBufferString(script + "\n"),
	}
}

// ScriptFile
func (c *Client) ScriptFile(fname string) *remoteScript {
	return &remoteScript{
		_type:      scriptFile,
		client:     c.client,
		scriptFile: fname,
	}
}

type remoteScript struct {
	client     *ssh.Client
	_type      remoteScriptType
	script     *bytes.Buffer
	scriptFile string
	err        error

	stdout io.Writer
	stderr io.Writer
}

// Run
func (rs *remoteScript) Run() error {
	if rs.err != nil {
		fmt.Println(rs.err)
		return rs.err
	}

	if rs._type == cmdLine {
		return rs.runCmds()
	} else if rs._type == rawScript {
		return rs.runScript()
	} else if rs._type == scriptFile {
		return rs.runScriptFile()
	} else {
		return errors.New("Not supported remoteScript type")
	}
}

func (rs *remoteScript) Output() ([]byte, error) {
	if rs.stdout != nil {
		return nil, errors.New("Stdout already set")
	}
	var out bytes.Buffer
	rs.stdout = &out
	err := rs.Run()
	return out.Bytes(), err
}

func (rs *remoteScript) SmartOutput() ([]byte, error) {
	if rs.stdout != nil {
		return nil, errors.New("Stdout already set")
	}
	if rs.stderr != nil {
		return nil, errors.New("Stderr already set")
	}

	var (
		stdout bytes.Buffer
		stderr bytes.Buffer
	)
	rs.stdout = &stdout
	rs.stderr = &stderr
	err := rs.Run()
	if err != nil {
		return stderr.Bytes(), err
	}
	return stdout.Bytes(), err
}

func (rs *remoteScript) Cmd(cmd string) *remoteScript {
	_, err := rs.script.WriteString(cmd + "\n")
	if err != nil {
		rs.err = err
	}
	return rs
}

func (rs *remoteScript) SetStdio(stdout, stderr io.Writer) *remoteScript {
	rs.stdout = stdout
	rs.stderr = stderr
	return rs
}

func (rs *remoteScript) runCmd(cmd string) error {
	session, err := rs.client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	session.Stdout = rs.stdout
	session.Stderr = rs.stderr

	if err := session.Run(cmd); err != nil {
		return err
	}
	return nil
}

func (rs *remoteScript) runCmds() error {
	for {
		statment, err := rs.script.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if err := rs.runCmd(statment); err != nil {
			return err
		}
	}

	return nil
}
func (rs *remoteScript) runScript() error {
	session, err := rs.client.NewSession()
	if err != nil {
		return err
	}

	session.Stdin = rs.script
	session.Stdout = rs.stdout
	session.Stderr = rs.stderr

	if err := session.Shell(); err != nil {
		return err
	}
	if err := session.Wait(); err != nil {
		return err
	}

	return nil
}

func (rs *remoteScript) runScriptFile() error {
	var buffer bytes.Buffer
	file, err := os.Open(rs.scriptFile)
	if err != nil {
		return err
	}
	_, err = io.Copy(&buffer, file)
	if err != nil {
		return err
	}

	rs.script = &buffer
	return rs.runScript()
}

type remoteShell struct {
	client         *ssh.Client
	requestPty     bool
	terminalConfig *TerminalConfig

	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

type TerminalConfig struct {
	Term   string
	Height int
	Weight int
	Modes  ssh.TerminalModes
}

// Terminal create a interactive shell on client.
func (c *Client) Terminal(config *TerminalConfig) *remoteShell {
	return &remoteShell{
		client:         c.client,
		terminalConfig: config,
		requestPty:     true,
	}
}

// Shell create a noninteractive shell on client.
func (c *Client) Shell() *remoteShell {
	return &remoteShell{
		client:     c.client,
		requestPty: false,
	}
}

func (rs *remoteShell) SetStdio(stdin io.Reader, stdout, stderr io.Writer) *remoteShell {
	rs.stdin = stdin
	rs.stdout = stdout
	rs.stderr = stderr
	return rs
}

// Start start a remote shell on client
func (rs *remoteShell) Start() error {
	session, err := rs.client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	if rs.stdin == nil {
		session.Stdin = os.Stdin
	} else {
		session.Stdin = rs.stdin
	}
	if rs.stdout == nil {
		session.Stdout = os.Stdout
	} else {
		session.Stdout = rs.stdout
	}
	if rs.stderr == nil {
		session.Stderr = os.Stderr
	} else {
		session.Stderr = rs.stderr
	}
	if rs.requestPty {
		tc := rs.terminalConfig
		if tc == nil {
			tc = &TerminalConfig{
				Term:   "xterm",
				Height: 40,
				Weight: 80,
			}
		}
		if err := session.RequestPty(tc.Term, tc.Height, tc.Weight, tc.Modes); err != nil {
			return err
		}
	}

	if err := session.Shell(); err != nil {
		return err
	}

	if err := session.Wait(); err != nil {
		return err
	}

	return nil
}
