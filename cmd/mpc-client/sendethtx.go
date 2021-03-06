package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/anyswap/mpc-client/cmd/utils"
	"github.com/anyswap/mpc-client/log"
	"github.com/anyswap/mpc-client/mpcrpc"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/urfave/cli/v2"
)

var (
	sendEthTxCommand = &cli.Command{
		Action:      sendEthTx,
		Name:        "sendethtx",
		Usage:       "send eth-like transaction",
		ArgsUsage:   "",
		Description: ``,
		Flags: []cli.Flag{
			pubkeyFlag,
			gidFlag,
			thresholdFlag,
			signModeFlag,
			signMemoFlag,
			mpcServerFlag,
			mpcKeystoreFlag,
			mpcPasswordFlag,
			signTypeFlag,
			apiPrefixFlag,
			rpcTimeoutFlag,
			signTimeoutFlag,
			gatewaysFlag,
			chainIDFlag,
			createContractFlag,
			fromAddrFlag,
			toAddrFlag,
			nonceFlag,
			valueFlag,
			gasLimitFlag,
			gasPriceFlag,
			inputFlag,
			dryrunFlag,
		},
	}
)

type sendEthTxArgs struct {
	gateways []string
	from     common.Address
	to       common.Address
	gasLimit uint64
	gasPrice *big.Int
	chainID  *big.Int
	accNonce *big.Int
	value    *big.Int
	input    []byte
	dryrun   bool

	createContract bool
}

type ethClientAndURL struct {
	cli *ethclient.Client
	url string
}

var (
	txArgs sendEthTxArgs

	ethClients []*ethClientAndURL

	bgCtx = context.Background()
)

func checkSendEthTxArguments(ctx *cli.Context) (err error) {
	txArgs.gateways = ctx.StringSlice(gatewaysFlag.Name)
	txArgs.gasLimit = ctx.Uint64(gasLimitFlag.Name)
	txArgs.dryrun = ctx.Bool(dryrunFlag.Name)

	fromAddrStr := ctx.String(fromAddrFlag.Name)
	if !common.IsHexAddress(fromAddrStr) {
		return fmt.Errorf("wrong from address %v", fromAddrStr)
	}
	txArgs.from = common.HexToAddress(fromAddrStr)

	txArgs.createContract = ctx.Bool(createContractFlag.Name)
	toAddrStr := ctx.String(toAddrFlag.Name)
	if !txArgs.createContract {
		if !common.IsHexAddress(toAddrStr) {
			return fmt.Errorf("wrong to address %v", toAddrStr)
		}
		txArgs.to = common.HexToAddress(toAddrStr)
	} else if toAddrStr != "" {
		return errors.New("create contract tx forbid specify 'to' address")
	}

	var ok bool
	gasPriceStr := ctx.String(gasPriceFlag.Name)
	txArgs.gasPrice, ok = new(big.Int).SetString(gasPriceStr, 0)
	if !ok {
		return fmt.Errorf("wrong gas price %v", gasPriceStr)
	}

	nodeChainIDStr := ctx.String(chainIDFlag.Name)
	txArgs.chainID, ok = new(big.Int).SetString(nodeChainIDStr, 0)
	if !ok {
		return fmt.Errorf("wrong chain Id %v", nodeChainIDStr)
	}

	accNonceStr := ctx.String(nonceFlag.Name)
	if accNonceStr != "" {
		txArgs.accNonce, ok = new(big.Int).SetString(accNonceStr, 0)
		if !ok {
			return fmt.Errorf("wrong account nonce %v", accNonceStr)
		}
	}

	valueStr := ctx.String(valueFlag.Name)
	if valueStr != "" {
		txArgs.value, ok = new(big.Int).SetString(valueStr, 0)
		if !ok {
			return fmt.Errorf("wrong value %v", valueStr)
		}
	}

	inputData := ctx.String(inputFlag.Name)
	if inputData != "" {
		txArgs.input, err = hexutil.Decode(inputData)
		if err != nil {
			return fmt.Errorf("wrong input data %v, err=%v", inputData, err)
		}
	}

	log.Info("check arguments pass")
	return nil
}

func sendEthTx(ctx *cli.Context) (err error) {
	utils.SetLogger(ctx)
	mpcCfg.NeedKeyStore = true
	err = checkAndInitMpcConfig(ctx, true)
	if err != nil {
		return err
	}
	err = checkSendEthTxArguments(ctx)
	if err != nil {
		return err
	}

	err = dailGateways(txArgs.gateways)
	if err != nil {
		return err
	}

	var nonce uint64
	if txArgs.accNonce != nil {
		nonce = txArgs.accNonce.Uint64()
	} else {
		nonce, err = getPendingNonce(txArgs.from)
		if err != nil {
			log.Error("get account nonce failed", "account", txArgs.from.String(), "err", err)
			return err
		}
		log.Info("get account nonce success", "account", txArgs.from.String(), "nonce", nonce)
	}

	var rawTx *types.Transaction
	if txArgs.createContract {
		rawTx = types.NewContractCreation(nonce, txArgs.value, txArgs.gasLimit, txArgs.gasPrice, txArgs.input)
	} else {
		rawTx = types.NewTransaction(nonce, txArgs.to, txArgs.value, txArgs.gasLimit, txArgs.gasPrice, txArgs.input)
	}
	log.Info("create raw tx success")
	_ = printTx(rawTx, true)

	chainSigner := types.NewEIP155Signer(txArgs.chainID)
	msgHash := chainSigner.Hash(rawTx)
	txJSON, err := json.Marshal(rawTx)
	if err != nil {
		log.Error("json marshal tx failed")
		return err
	}
	msgContext := []string{"ethtx", string(txJSON), txArgs.chainID.String()}
	if signMemoArg != "" {
		msgContext = append(msgContext, signMemoArg)
	}

	keyID, rsvs, err := mpcrpc.DoSign(mpcPublicKey, []string{msgHash.String()}, msgContext)
	if err != nil {
		log.Error("mpc sign failed", "err", err)
		return err
	}
	log.Info("mpc sign success", "keyID", keyID)

	if len(rsvs) != 1 {
		log.Error("mpc sign result rsv count is wrong", "have", len(rsvs), "want", 1)
		return errors.New("mpc sign result rsv count is wrong")
	}
	rsv := rsvs[0]

	signature := common.FromHex(rsv)
	if len(signature) != crypto.SignatureLength {
		log.Error("mpc sign result rsv length is wrong", "rsv", rsv)
		return errors.New("mpc sign result rsv length is wrong")
	}

	signedTx, err := rawTx.WithSignature(chainSigner, signature)
	if err != nil {
		log.Error("sign tx failed", "err", err)
		return err
	}

	sender, err := types.Sender(chainSigner, signedTx)
	if err != nil {
		log.Error("get sender from signed tx failed", "err", err)
		return err
	}

	if sender != txArgs.from {
		log.Error("sender mismatch", "signer", sender.String(), "sender", txArgs.from.String())
		return errors.New("sender mismatch")
	}

	txHash := signedTx.Hash().String()

	log.Info("mpc sign tx success", "txHash", txHash, "sender", sender.String())
	_ = printTx(signedTx, false)

	if !txArgs.dryrun {
		err = sendSignedTransaction(signedTx)
		if err != nil {
			log.Error("send tx failed", "err", err)
			return err
		}
		log.Info("send tx success", "txHash", txHash)
	}
	return nil
}

func dailGateways(gateways []string) (err error) {
	ethClients = make([]*ethClientAndURL, 0, len(gateways))
	cliURLs := make([]string, 0, len(gateways))
	var ethClient *ethclient.Client
	for _, gateway := range gateways {
		ethClient, err = ethclient.Dial(gateway)
		if err != nil {
			log.Warn("dail gateway failed", "url", gateway, "err", err)
			continue
		}
		ethClients = append(ethClients, &ethClientAndURL{cli: ethClient, url: gateway})
		cliURLs = append(cliURLs, gateway)
	}
	if len(ethClients) > 0 {
		log.Info("dail gateways success", "clients", cliURLs)
		return nil
	}
	return err
}

func getPendingNonce(account common.Address) (maxNonce uint64, err error) {
	var success bool
	var nonce uint64
	for _, ethClient := range ethClients {
		nonce, err = ethClient.cli.PendingNonceAt(bgCtx, account)
		if err != nil {
			log.Warn("get pending nonce failed", "account", account.String(), "url", ethClient.url, "err", err)
			continue
		}
		success = true
		if nonce > maxNonce {
			maxNonce = nonce
		}
	}
	if success {
		return maxNonce, nil
	}
	return 0, err
}

func sendSignedTransaction(signedTx *types.Transaction) (err error) {
	var success bool
	for _, ethClient := range ethClients {
		err = ethClient.cli.SendTransaction(bgCtx, signedTx)
		if err != nil {
			log.Warn("send tx failed", "hash", signedTx.Hash().String(), "url", ethClient.url, "err", err)
			continue
		}
		success = true
	}
	if success {
		return nil
	}
	return err
}

func printTx(tx *types.Transaction, jsonFmt bool) error {
	if jsonFmt {
		bs, err := json.MarshalIndent(tx, "", "  ")
		if err != nil {
			return fmt.Errorf("json marshal err %v", err)
		}
		fmt.Println(string(bs))
		_, r, _ := tx.RawSignatureValues()
		if r == nil || r.Sign() == 0 {
			fmt.Printf("tx value is %v, nonce is %v, gasPrice is %v, gasLimit is %v\n", tx.Value(), tx.Nonce(), tx.GasPrice(), tx.Gas())
		} else {
			fmt.Printf("tx chainID is %v, value is %v, nonce is %v, gasPrice is %v, gasLimit is %v\n", tx.ChainId(), tx.Value(), tx.Nonce(), tx.GasPrice(), tx.Gas())
		}
	} else {
		bs, err := tx.MarshalBinary()
		if err != nil {
			return fmt.Errorf("marshal tx err %v", err)
		}
		fmt.Println(hexutil.Bytes(bs))
	}
	return nil
}
