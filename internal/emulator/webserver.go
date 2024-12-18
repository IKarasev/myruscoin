package emulator

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	HTTP_ADDR            string = "127.0.0.1"
	HTTP_PORT            string = "8080"
	RSS_READ_UPDATE_TIME        = time.Millisecond * 100
	OP_PAUSE_MILISEC            = time.Millisecond * 500
	WITH_LOG                    = false
)

type EmulatorWeb struct {
	RcMngr            *RuscoinMngr
	E                 *echo.Echo
	rssChan           RssChan
	RssReadUpdateTime time.Duration
	ctx               context.Context
}

func (e *EmulatorWeb) TestRoutine() {

}

// Loads web server settings from Enviroment variables
//
// RUSCOIN_HTTP_ADDR  - ip address for web server to listen to
//
// RUSCOIN_HTTP_PORT  - port for web server to listen
//
// RUSCOIN_RSS_UPDATE - send update period in Milliseconds for RSS messages
func LoadSettingsFromEnv() error {
	errStr := ""
	if v := os.Getenv("RUSCOIN_HTTP_ADDR"); v != "" {
		HTTP_ADDR = v
	}
	if v := os.Getenv("RUSCOIN_HTTP_PORT"); v != "" {
		if _, err := strconv.Atoi(v); err != nil {
			errStr += "Failed to paser HTTP_PORT env variable\n"
		}
		HTTP_PORT = v
	}
	if v := os.Getenv("RUSCOIN_HTTP_ADDR"); v != "" {
		HTTP_ADDR = v
	}
	if v := os.Getenv("RUSCOIN_RSS_UPDATE"); v != "" {
		if c, err := strconv.Atoi(v); err == nil {
			RSS_READ_UPDATE_TIME = time.Millisecond * time.Duration(c)
		} else {
			errStr += "Failed to pase RUSCOIN_RSS_UPDATE env variable\n"
		}
	}
	if v := os.Getenv("OP_PAUSE_MILISEC"); v != "" {
		if c, err := strconv.Atoi(v); err == nil {
			OP_PAUSE_MILISEC = time.Millisecond * time.Duration(c)
		} else {
			errStr += "Failed to pase OP_PAUSE_MILISEC env variable\n"
		}
	}
	if v := os.Getenv("WITH_LOG"); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			WITH_LOG = b
		} else {
			errStr += "Failed to pase WITH_LOG env variable\n"
		}
	}
	if errStr != "" {
		return fmt.Errorf(errStr)
	}
	return nil

}

func NewEmulatorWeb() *EmulatorWeb {
	return &EmulatorWeb{
		RcMngr:            NewRuscoinMngr(),
		E:                 echo.New(),
		rssChan:           make(RssChan),
		RssReadUpdateTime: RSS_READ_UPDATE_TIME,
		ctx:               context.Background(),
	}
}

func (wb *EmulatorWeb) DefaultRcManager() *EmulatorWeb {
	wb.RcMngr = DefaultRuscoinMngr()
	return wb
}

func (wb *EmulatorWeb) StartWithLogger() {
	wb.E.Use(middleware.Logger())
	wb.E.Logger.Fatal(wb.Start())
}

func (wb *EmulatorWeb) Start() error {
	fmt.Println(HTTP_ADDR + ":" + HTTP_PORT)
	ctx, ctxDone := context.WithCancel(wb.ctx)
	wb.ctx = ctx
	defer ctxDone()
	defer close(wb.rssChan)

	wb.initRoutes()
	err := wb.E.Start(HTTP_ADDR + ":" + HTTP_PORT)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (wb *EmulatorWeb) initRoutes() {
	wd, _ := os.Getwd()
	wd = wd + "/assets"

	wb.E.Static("/static", "assets")

	wb.E.GET("/", wb.HandleIndex)

	wb.E.GET("/test", wb.HandleTest)

	wb.E.GET("/sse", wb.HandleSse)

	wb.E.GET("/nodelist", wb.HandleNodeList)

	wb.E.GET("/selectminer", wb.HandleMinerSelect)

	wb.E.GET("/tick", wb.HandleTick)

	wb.E.GET("/settings", wb.HandleEimulationSettings)

	gNode := wb.E.Group("/node")
	gNode.GET("/slist", wb.HandleNodeSelectList)
	gNode.POST("/info", wb.HandleNodeInfo)
	gNode.POST("/block", wb.HandleBlockDetails)
	gNode.POST("/block/tr", wb.HandleBlockTransactions)

	gWallet := wb.E.Group("/wallet")
	gWallet.POST("/slist", wb.HandleWalletList)
	gWallet.POST("/utxotable", wb.HandleWalletUtxoTable)
	gWallet.POST("/addtr", wb.HandleAddTransaction)
	gWallet.POST("/blocktr", wb.HandleWalletBlockTr)

	gEvil := wb.E.Group("/evil")
	gEvil.GET("/load", wb.HandleEvilLoad)
	gEvil.GET("/steal", wb.HandleEvilSteal)
	gEvil.GET("/mine", wb.HandleEvilMine)
	gEvil.GET("/inject", wb.HandleEvilInject)
	gEvil.GET("/send", wb.HandleEvilSend)

	gEvilSet := gEvil.Group("/set")
	gEvilSet.POST("/height", wb.HandleEvilSetHeihgt)
	gEvilSet.POST("/time", wb.HandleEvilSetTime)
	gEvilSet.POST("/hash", wb.HandleEvilSetHash)
	gEvilSet.POST("/nonce", wb.HandleEvilSetInt)
	gEvilSet.POST("/coinbase", wb.HandleEvilSetInt)
	gEvilSet.POST("/tr/sign", wb.HandleEvilSetTrHashValue)
	gEvilSet.POST("/tr/pk", wb.HandleEvilSetTrHashValue)
	gEvilSet.POST("/tr/utxo", wb.HandleEvilSetTrUtxo)

	gEvilDel := gEvil.Group("/del")
	gEvilDel.POST("/utxo", wb.HandleEvilDelUtxo)
	gEvilDel.POST("/tr", wb.HandleEvilDelTr)

	gEvilAdd := gEvil.Group("/add")
	gEvilAdd.GET("/tr", wb.HandleEvilAddTr)
	gEvilAdd.POST("/utxo", wb.HandleEvilAddUtxo)
}
