package main

import (
	"github.com/urfave/cli/v2"
)

var (
	pubkeyFlag = &cli.StringFlag{
		Name:  "pubkey",
		Usage: "mpc public key",
	}
	msgHashFlag = &cli.StringFlag{
		Name:  "msghash",
		Usage: "mpc sign message hash",
	}
	signMessageFlag = &cli.StringFlag{
		Name:  "signmsg",
		Usage: "mpc sign message of hex string",
	}
	msgContextFlag = &cli.StringFlag{
		Name:  "msgcontext",
		Usage: "mpc sign message context",
	}
	keyIDFlag = &cli.StringFlag{
		Name:  "key",
		Usage: "mpc sign key ID",
	}
	nonInteractiveFlag = &cli.BoolFlag{
		Name:  "non-interactive",
		Usage: "open non interactive mode",
	}
	agreeSignFlag = &cli.BoolFlag{
		Name:  "agree",
		Usage: "agree sgin non-interactively",
	}
	disagreeSignFlag = &cli.BoolFlag{
		Name:  "disagree",
		Usage: "disagree sgin non-interactively",
	}
	apiPrefixFlag = &cli.StringFlag{
		Name:  "apiPrefix",
		Usage: "mpc rpc apiPrefix",
		Value: "dcrm_",
	}
	rpcTimeoutFlag = &cli.Uint64Flag{
		Name:  "rpcTimeout",
		Usage: "mpc rpc timeout of seconds",
		Value: 20,
	}
	signTimeoutFlag = &cli.Uint64Flag{
		Name:  "signTimeout",
		Usage: "mpc sign timeout of seconds",
		Value: 120,
	}
	signTypeFlag = &cli.StringFlag{
		Name:  "keytype",
		Usage: "mpc sign algorithm type",
		Value: "ECDSA",
	}
	gidFlag = &cli.StringFlag{
		Name:  "gid",
		Usage: "mpc sign group ID",
	}
	groupIDFlag = &cli.StringFlag{
		Name:  "gid",
		Usage: "mpc group ID",
	}
	thresholdFlag = &cli.StringFlag{
		Name:  "ts",
		Usage: "mpc sign threshold",
		Value: "3/5",
	}
	signModeFlag = &cli.Uint64Flag{
		Name:  "mode",
		Usage: "mpc sign mode (private=1/managed=0)",
		Value: 0,
	}
	signMemoFlag = &cli.StringFlag{
		Name:  "memo",
		Usage: "mpc sign memo text",
	}
	mpcServerFlag = &cli.StringFlag{
		Name:  "url",
		Usage: "mpc server URL",
	}
	mpcUserFlag = &cli.StringFlag{
		Name:  "user",
		Usage: "mpc user address",
	}
	mpcKeystoreFlag = &cli.StringFlag{
		Name:  "keystore",
		Usage: "mpc user keystore file",
	}
	mpcPasswordFlag = &cli.StringFlag{
		Name:  "passwd",
		Usage: "mpc user password file",
	}
	mpcDKGFlag = &cli.BoolFlag{
		Name:  "dkg",
		Usage: "is mpc public key generation",
	}
	enodeSigsFlag = &cli.StringSliceFlag{
		Name:  "sig",
		Usage: "group member enodes sigs (multiple)",
	}
	showEnodeSigFlag = &cli.BoolFlag{
		Name:  "sig",
		Usage: "show enode sig",
	}
	expiredIntervalFlag = &cli.Int64Flag{
		Name:  "expiredInterval",
		Usage: "expired interval of seconds",
	}

	gatewaysFlag = &cli.StringSliceFlag{
		Name:  "gateway",
		Usage: "gateway URLs of full nodes (multiple)",
	}
	chainIDFlag = &cli.StringFlag{
		Name:  "chainID",
		Usage: "blockchain ID",
	}
	createContractFlag = &cli.BoolFlag{
		Name:  "createContract",
		Usage: "tx is to create contract",
	}
	fromAddrFlag = &cli.StringFlag{
		Name:  "from",
		Usage: "tx sender address",
	}
	toAddrFlag = &cli.StringFlag{
		Name:  "to",
		Usage: "tx receiver address",
	}
	inputFlag = &cli.StringFlag{
		Name:  "input",
		Usage: "tx input data",
	}
	gasLimitFlag = &cli.Uint64Flag{
		Name:  "gas",
		Usage: "tx gas limit",
		Value: 90000,
	}
	gasPriceFlag = &cli.StringFlag{
		Name:  "gasPrice",
		Usage: "tx gas price in Wei",
	}
	nonceFlag = &cli.StringFlag{
		Name:  "nonce",
		Usage: "tx nonce",
	}
	valueFlag = &cli.StringFlag{
		Name:  "value",
		Usage: "tx value of native coins",
	}
	dryrunFlag = &cli.BoolFlag{
		Name:  "dryrun",
		Usage: "dry run",
	}
)
