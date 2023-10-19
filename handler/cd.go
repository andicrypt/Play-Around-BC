package handler

import (
	"crypto/ecdsa"
	"strings"

	seaportPkg "github.com/andicrypt/play-around-bc/generated_contracts"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)
		 

type Handler interface {
	GetHandlers() map[string]func(c *gin.Context, sender, selector, callData string)
}
type DCHandler struct {

	// Seaport struct go
	seaport *seaportPkg.Seaport
	seaportABI abi.ABI

	handlers map[string]func(method abi.Method, name string, args ...interface{}) ([]byte, uint64, error)

	privateKey *ecdsa.PrivateKey
	db *gorm.DB
}

func NewDCHandler(config *Config, db *gorm.DB) (*DCHandler, error) {
	client, err := ethclient.Dial(config.RPC)
	if err != nil {
		return nil, err
	}
	seaport, err := seaportPkg.NewSeaport(common.HexToAddress(config.Registry), client)
	if err != nil {
		return nil, err
	}
	seaportABI, err := abi.JSON(strings.NewReader(seaportPkg.SeaportMetaData.ABI))

	if err != nil {
		return nil, err
	}
	privateKey, err := crypto.HexToECDSA(config.PrivateKey)
	if err != nil {
		return nil, err
	}

	h := &DCHandler{
		seaport:	seaport,
		seaportABI: seaportABI,
		privateKey: privateKey,
		db: db,
	}

	h.handlers = map[string]func(abi.Method, string, ...interface{}) ([]byte, uint64, error) {
	}

	return nil, nil

}

func (h* DCHandler) GetHandlers() map[string]func(c *gin.Context, sender, selector, callData string) {
	return map[string] func(c *gin.Context, sender, selector, callData string) {

	}
} 