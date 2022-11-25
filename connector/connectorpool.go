package connector

import "strconv"

type ConnectorPool struct {
	Host       string
	Port       int
	Username   string
	Password   string
	DbPrefix   string
	ThreadNum  int
	ThreadPool chan *Connector
}

// NewConnectorPool: create thread pool with size ThreadNum, fill with *Connector,
// the database name of each connector is config.DbPrefix + thread id
func NewConnectorPool(host string, port int, username string, password string,
	dbPrefix string, threadNum int) (*ConnectorPool, error) {
	connectorPool := &ConnectorPool{
		Host:      host,
		Port:      port,
		Username:  username,
		Password:  password,
		DbPrefix:  dbPrefix,
		ThreadNum: threadNum,
	}
	threadPool := make(chan *Connector, threadNum)
	for i := 0; i < threadNum; i++ {
		conn, err := NewConnector(host, port, username, password, dbPrefix+strconv.Itoa(i))
		if err != nil {
			return nil, err
		}
		threadPool <- conn
	}
	connectorPool.ThreadPool = threadPool
	return connectorPool, nil
}

// ConnectorPool.WaitForFree wait for a free connector
func (connPool *ConnectorPool) WaitForFree() *Connector {
	conn := <-connPool.ThreadPool
	return conn
}

// ConnectorPool.BackToPool: give a connector back to pool
func (connPool *ConnectorPool) BackToPool(conn *Connector) {
	connPool.ThreadPool <- conn
}