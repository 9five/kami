package router

import (
	"kami/config"
	_encryptionUsecase "kami/encryption/usecase"
	_kamiOrderHandler "kami/kamiOrder/delivery/handler"
	_kamiOrderRepo "kami/kamiOrder/repository/postgresql"
	_kamiOrderUsecase "kami/kamiOrder/usecase"
	_kamiUserHandler "kami/kamiUser/delivery/handler"
	_kamiUserRepo "kami/kamiUser/repository/postgresql"
	_kamiUserUsecase "kami/kamiUser/usecase"
	_lotteryHandler "kami/lottery/delivery/handler"
	_lotteryRepo "kami/lottery/repository/postgresql"
	_lotteryUsecase "kami/lottery/usecase"
	_middlewareUsecase "kami/middleware/usecase"
	_twilioServiceRepo "kami/twilioService/repository"
	_twilioServiceUsecase "kami/twilioService/usecase"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "kami/docs"
)

var middlewareUsecase = _middlewareUsecase.NewMiddlewareUsecase(config.JwtKey)

func NewRouter() *gin.Engine {
	r := gin.Default()
	setupRouter(r)
	return r
}

func setupRouter(router *gin.Engine) {
	setupCORS(router)
	setupSwagger(router)
	setupLoginSys(router)
	setupUser(router)
	setupOrderProcess(router)
	setuplottery(router)
}

func setupCORS(router *gin.Engine) {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = config.AllowOrigin
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Authorization")
	router.Use(cors.New(corsConfig))
}

func setupSwagger(router *gin.Engine) {
	if gin.Mode() == gin.DebugMode {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
}

func setupLoginSys(router *gin.Engine) {
	twilioServiceRepo := _twilioServiceRepo.NewTwilioServiceRepository()
	twilioServiceUsecase := _twilioServiceUsecase.NewTwilioServiceUsecase(twilioServiceRepo)
	kamiUserRepo := _kamiUserRepo.NewPostgresqlKamiUserRepository(config.DB)
	kamiUserUsecase := _kamiUserUsecase.NewKamiUserUsecase(config.JwtKey, kamiUserRepo)
	encryptionUsecase := _encryptionUsecase.NewEncryptionUsecase(config.DesKey)
	_kamiUserHandler.NewLoginHandler(router.Group("/api/login"), kamiUserUsecase, twilioServiceUsecase, encryptionUsecase)
}

func setupUser(router *gin.Engine) {
	kamiUserRepo := _kamiUserRepo.NewPostgresqlKamiUserRepository(config.DB)
	kamiUserUsecase := _kamiUserUsecase.NewKamiUserUsecase(config.JwtKey, kamiUserRepo)
	_kamiUserHandler.NewKamiUserHandler(router.Group("/api/user", middlewareUsecase.VerifyToken), kamiUserUsecase)
}

func setupOrderProcess(router *gin.Engine) {
	kamiUserRepo := _kamiUserRepo.NewPostgresqlKamiUserRepository(config.DB)
	kamiUserUsecase := _kamiUserUsecase.NewKamiUserUsecase(config.JwtKey, kamiUserRepo)
	kamiOrderRepo := _kamiOrderRepo.NewPostgresqlKamiOrderRepository(config.DB)
	kamiOrderUsecase := _kamiOrderUsecase.NewKamiOrderUsecase(kamiOrderRepo)
	_kamiOrderHandler.NewKamiOrderHandler(router.Group("/api/order", middlewareUsecase.VerifyToken), kamiOrderUsecase, kamiUserUsecase)
}

func setuplottery(router *gin.Engine) {
	prizePoolRepo := _lotteryRepo.NewPostgresqlPrizePoolRepository(config.DB)
	prizePoolUsecase := _lotteryUsecase.NewLotteryPrizePoolUsecase(config.AwsBucket, prizePoolRepo)
	prizeCardRepo := _lotteryRepo.NewPostgresqlPrizeCardRepository(config.DB)
	prizeCardUsecase := _lotteryUsecase.NewLotteryPrizeCardUsecase(config.AwsBucket, prizeCardRepo)
	kamiUserRepo := _kamiUserRepo.NewPostgresqlKamiUserRepository(config.DB)
	kamiUserUsecase := _kamiUserUsecase.NewKamiUserUsecase(config.JwtKey, kamiUserRepo)
	_lotteryHandler.NewLotteryHandler(router.Group("/api/lottery", middlewareUsecase.VerifyToken), prizePoolUsecase, prizeCardUsecase, kamiUserUsecase)
}
