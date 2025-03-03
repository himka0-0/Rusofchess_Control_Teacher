package router

import (
	"awesomeProject1/controllers"
	"awesomeProject1/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.GET("/", controllers.FirstPage)
	r.GET("/registration", controllers.RegPage)
	r.GET("/authentication", controllers.AutPage)
	r.GET("/verify-email", controllers.VerifyEmailPage)
	r.GET("/recovery-password", controllers.RecMailPage)
	r.GET("/verify", controllers.VerifyPage)
	r.GET("/recoveryPassword", controllers.RecoveryPasPage)
	r.POST("/registration", controllers.RegHandler)
	r.POST("/authentication", controllers.AutHandler)
	r.POST("/recoveryPassword", controllers.RecoveryPasHandler)
	r.POST("/recovery-password", controllers.RecMailHandler)
	autRoute := r.Group("/")
	autRoute.Use(middlewares.AuthMiddleware())
	{
		autRoute.GET("/paymentstudent", controllers.PaymentstudentPage)
		autRoute.GET("/kabinet", controllers.GetProfile)
		autRoute.GET("/lecture", controllers.LecturePage)
		autRoute.GET("/firstsetting", controllers.FirstSettinPage)
		autRoute.GET("/notelesson", controllers.NotelessonPage)
		autRoute.GET("/result", controllers.ResultPage)
		autRoute.GET("/student", controllers.StudentPage)
		autRoute.GET("/telbot", controllers.TelbotPage)
		autRoute.GET("/we", controllers.WePage)
		autRoute.GET("/instuction", controllers.InstructionPage)
		autRoute.POST("/firstsetting", controllers.FirstSettingHandler)
		autRoute.POST("/logout", controllers.LogoutHandler)
		autRoute.POST("/paymentstudent", controllers.PaymentstudentHandler)
		autRoute.POST("/notelesson", controllers.NotelessonHandler)
		autRoute.POST("/lecture", controllers.LectureHandler)
		autRoute.POST("/student", controllers.StudentHandler)
		autRoute.POST("/telbot", controllers.TelbotHandler)
	}
}
