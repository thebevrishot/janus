package cli

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	"github.com/qtumproject/janus/pkg/qtum"
	"github.com/qtumproject/janus/pkg/server"
	"github.com/qtumproject/janus/pkg/transformer"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app = kingpin.New("janus", "Qtum adapter to Ethereum JSON RPC")

	accountsFile = app.Flag("accounts", "account addresses (in hex) returned by eth_accounts").File()

	qtumRPC     = app.Flag("qtum-rpc", "URL of qtum RPC service").Envar("QTUM_RPC").Default("").String()
	qtumNetwork = app.Flag("qtum-network", "").Envar("QTUM_NETWORK").Default("regtest").String()
	bind        = app.Flag("bind", "network interface to bind to (e.g. 0.0.0.0) ").Default("localhost").String()
	port        = app.Flag("port", "port to serve proxy").Default("23889").Int()

	devMode = app.Flag("dev", "[Insecure] Developer mode").Default("false").Bool()
)

func loadAccounts(r io.Reader) []string {
	var accounts []string

	if accountsFile != nil {
		s := bufio.NewScanner(*accountsFile)
		for s.Scan() {
			addr := s.Text()
			if addr == "" {
				continue
			}

			buf, err := hex.DecodeString(addr)
			if err != nil {
				continue
			}

			if len(buf) != 20 {
				continue
			}

			accounts = append(accounts, s.Text())
		}
	}

	return accounts
}

func action(pc *kingpin.ParseContext) error {
	addr := fmt.Sprintf("%s:%d", *bind, *port)
	logger := log.NewLogfmtLogger(os.Stdout)

	if !*devMode {
		logger = level.NewFilter(logger, level.AllowWarn())
	}

	var accounts []string
	if *accountsFile != nil {
		accounts = loadAccounts(*accountsFile)
		(*accountsFile).Close()
	}

	fmt.Println("ETH accounts", accounts)

	qtumJSONRPC, err := qtum.NewClient(*qtumRPC, qtum.SetDebug(*devMode), qtum.SetLogger(logger), qtum.SetETHAccounts(accounts))
	if err != nil {
		return errors.Wrap(err, "jsonrpc#New")
	}

	qtumClient, err := qtum.New(qtumJSONRPC, *qtumNetwork)
	if err != nil {
		return errors.Wrap(err, "qtum#New")
	}

	t, err := transformer.New(qtumClient, transformer.DefaultProxies(qtumClient), transformer.SetDebug(*devMode), transformer.SetLogger(logger))
	if err != nil {
		return errors.Wrap(err, "transformer#New")
	}

	s, err := server.New(qtumClient, t, addr, server.SetLogger(logger), server.SetDebug(*devMode))
	if err != nil {
		return errors.Wrap(err, "server#New")
	}

	return s.Start()
}

func Run() {
	kingpin.MustParse(app.Parse(os.Args[1:]))
}

func init() {
	app.Action(action)
}
