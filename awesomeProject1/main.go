package main

import (
	"awesomeProject1/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

var jwtSecret = []byte("supersecretkey")

var db *gorm.DB

type User struct {
	ID                  uint   `gorm:"primarykey"`
	Name                string `gorm:"size 100;not null"`
	Email               string `gorm:"size 100;unique;not null"`
	Password            string `gorm:"not null"`
	Lectures_introduced int    `gorm:"size 100"`
	Email_verified      bool   `gorm:"default:false"`
	Verification_token  string `gorm:"size 100;unique"`
	Table_student       []Table_student
	Table_lecture       []Table_lecture
	Table_telegram_bot  []Table_telegram_bot
}
type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
} //для токeнов

type PostSettings struct {
	Meaning string `json:"meaning"` //ФИО
	Marking string `json:"marking"`
}
type Paymentstudent struct {
	ID      uint `json:"id"`
	Payment int  `json:"payment"`
}
type Table_student struct {
	ID             int    `gorm:"primarykey"`
	User_id        uint   `gorm:"index"`
	Name_Student   string `gorm:"size 100;not null"`
	Payment        int    `gorm:"size 100"`
	Theory         int    `gorm:"size 1000"`
	Practice       int    `gorm:"size 1000"`
	Tasks          int    `gorm:"size 1000"`
	Namber_lecture int    `gorm:"size 1000"`
	Alert_payment  bool   `json:"Alertpayment" gorm:"size:1000"`
	Alert_moduls   bool   `json:"Alertmodules" gorm:"size:1000"`
}

type Table_lecture struct {
	ID                int    `gorm:"primarykey"`
	User_id           uint   `gorm:"index"`
	Lecture           string `gorm:"size 100;not null"`
	Lecture_Person_id int    `gorm:"size 100;not null"`
}

type Table_telegram_bot struct {
	User_id     uint   `gorm:"index"`
	Hash        string `gorm:"size 100;unique;not null"`
	First_name  string `gorm:"size 100"`
	Telegram_id int64  `gorm:"size 100"`
	Vhod        bool   `gorm:"size 100;not null"`
}
type printLecture struct {
	Lecture_Person_id int    `json:"Lecture_Person_id"`
	Lecture           string `json:"Lecture"`
}
type PostModuls struct {
	Student_id   uint   `json:"Student_id"`
	Module       string `json:"Module"`
	Lock_lecture bool   `json:"lock_lecture"`
}

type PostLecture struct {
	Number  int    `json:"number"`
	Lecture string `json:"lecture"`
}

type PostStudent struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Lecture int    `json:"lecture"`
}

type resultstruct struct {
	ID           int
	Name_Student string
	Payment      int
	Lecture      string
	Theory       int
	Practice     int
	Tasks        int
}

type PostTelbot struct {
	ModuleAllToggle bool            `json:"moduleAllToggle"`
	Students        []Table_student `json:"students"`
}

func main() {
	dsn := "host=localhost user=postgres password=root dbname=mydatabase port=5432 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Ошибка подключения к БД:", err)
	}

	go RunBot()
	db.AutoMigrate(&User{}, &Table_student{}, &Table_lecture{}, &Table_telegram_bot{})

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	r.GET("/", firstPage)
	r.GET("/registration", regPage)
	r.GET("/authentication", autPage)
	r.GET("/verify-email", verifyEmailPage)
	r.GET("/verify", verifyPage)

	r.POST("/registration", regHandler)
	r.POST("/authentication", autHandler)
	//для авторизованных
	autRoute := r.Group("/")
	autRoute.Use(authMiddleware())
	autRoute.GET("/paymentstudent", paymentstudentPage)
	autRoute.GET("/kabinet", getProfile)
	autRoute.GET("/lecture", lecturePage)
	autRoute.GET("/firstsetting", firstSettinPage)
	autRoute.GET("/notelesson", notelessonPage)
	autRoute.GET("/result", resultPage)
	autRoute.GET("/student", studentPage)
	autRoute.GET("/telbot", telbotPage)
	autRoute.POST("/firstsetting", firstSettingHandler)
	autRoute.POST("/logout", logoutHandler)
	autRoute.POST("/paymentstudent", paymentstudentHandler)
	autRoute.POST("/notelesson", notelessonHandler)
	autRoute.POST("/lecture", lectureHandler)
	autRoute.POST("/student", studentHandler)
	autRoute.POST("/telbot", telbotHandler)

	r.Run(":8080")
}

// get запросы
func firstPage(c *gin.Context) {
	c.HTML(http.StatusOK, "nachalo.html", nil)
} //первая страница
func regPage(c *gin.Context) {
	c.HTML(http.StatusOK, "registration.html", nil)
} //стр регистрации
func autPage(c *gin.Context) {
	c.HTML(http.StatusOK, "authentication.html", nil)
} //стр аутентификации
func verifyEmailPage(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Токен отсутствует"})
		return
	}
	var user User
	if err := db.Where("verification_token=?", token).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Недействительный токен"})
		return
	}

	user.Email_verified = true
	user.Verification_token = ""
	db.Save(&user)
	c.HTML(http.StatusOK, "verifyEmail.html", gin.H{"message": "Email успешно подтвержден! Теперь вы можете войти."})
} //подтверждение почты
func verifyPage(c *gin.Context) {
	c.HTML(http.StatusOK, "verify.html", nil)
}
func firstSettinPage(c *gin.Context) {
	userData, exists := c.Get("User") // Берем пользователя из контекста
	if !exists || userData == nil {
		log.Println("Пользователь не найден в middleware")
		c.HTML(http.StatusUnauthorized, "error.html", gin.H{"error": "Пользователь не авторизован"})
		return
	}
	user := userData.(User)
	c.HTML(http.StatusOK, "firstSetting.html", gin.H{
		"User":          user,
		"signification": signification,
	})
} //стр первой настройки
func lecturePage(c *gin.Context) {
	userData, exists := c.Get("User")
	if !exists || userData == nil {
		log.Println("Пользователь не найден в middleware")
		c.HTML(http.StatusUnauthorized, "error.html", gin.H{"error": "Пользователь не авторизован"})
		return
	}
	user, ok := userData.(User)
	if !ok {
		log.Println("ошибка в переводе из формата")
		return
	}
	UserID := user.ID
	var data []printLecture
	err := db.Model(&Table_lecture{}).Select("Lecture_Person_id,Lecture").Where("User_id=?", UserID).Find(&data).Error
	if err != nil {
		log.Println("Не записывается инфа для передачи")
	}
	c.HTML(http.StatusOK, "lecture.html", gin.H{
		"data": data,
		"User": user,
	})
} //стр управления лекциями
func studentPage(c *gin.Context) {
	userData, exists := c.Get("User")
	if !exists || userData == nil {
		log.Println("проблема с вытаскиванием инфы про юзера", exists)
	}
	user := userData.(User)
	var lecture []Table_lecture
	err := db.Model(&Table_lecture{}).Select("Lecture_Person_id,Lecture").Where("User_id=?", user.ID).Find(&lecture).Error
	if err != nil {
		log.Println("ошибка при вытаскивании лекций", err)
	}
	var students []Table_student
	err = db.Model(&Table_student{}).Select("ID,Name_Student,Namber_lecture").Where("User_id=?", user.ID).Find(&students).Error
	if err != nil {
		log.Println("Ошибка при вытаскивании студента", err)
	}

	c.HTML(http.StatusOK, "student.html", gin.H{
		"students": students,
		"User":     user,
		"lecture":  lecture,
	})
} //стр управления учениками
func telbotPage(c *gin.Context) {
	userData, exists := c.Get("User")
	if !exists || userData == nil {
		log.Println("проблема с вытаскиванием инфы про юзера", exists)
	}
	user := userData.(User)
	var students []Table_student
	err := db.Model(&Table_student{}).Select("ID,Name_Student,Alert_payment,Alert_moduls").Where("User_id=?", user.ID).Find(&students).Error
	if err != nil {
		log.Println("стр телеграмбота ошибка с вытаскиванием студентов")
	}
	var vhod bool
	err = db.Model(&Table_telegram_bot{}).Select("Vhod").Where("User_id=?", user.ID).Scan(&vhod).Error
	c.HTML(http.StatusOK, "telgrambot.html", gin.H{
		"students": students,
		"User":     user,
		"vhod":     vhod,
	})
}
func notelessonPage(c *gin.Context) {
	userData, exists := c.Get("User") // Берем пользователя из контекста
	if !exists || userData == nil {
		log.Println("Пользователь не найден в middleware")
		c.HTML(http.StatusUnauthorized, "error.html", gin.H{"error": "Пользователь не авторизован"})
		return
	}
	user := userData.(User)
	tokenstr, _ := c.Cookie("token")
	claims := jwt.MapClaims{}
	_, _ = jwt.ParseWithClaims(tokenstr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	emailFromToken, _ := claims["email"].(string)
	var userID uint
	err := db.Model(&User{}).Select("id").Where("email = ?", emailFromToken).Scan(&userID).Error
	if err != nil {
		log.Println("не могу наути пользователя")
		return
	}
	var students []Table_student
	err = db.Model(&Table_student{}).Select("ID,Name_Student").Where("User_id = ?", userID).Scan(&students).Error
	if err != nil {
		log.Println(err)
	}
	c.HTML(http.StatusOK, "notelesson.html", gin.H{
		"students": students,
		"User":     user,
	})

} //стр отметки урока
func paymentstudentPage(c *gin.Context) {
	userData, exists := c.Get("User") // Берем пользователя из контекста
	if !exists || userData == nil {
		log.Println("Пользователь не найден в middleware")
		c.HTML(http.StatusUnauthorized, "error.html", gin.H{"error": "Пользователь не авторизован"})
		return
	}
	user := userData.(User)
	tokenstr, _ := c.Cookie("token")
	claims := jwt.MapClaims{}
	_, _ = jwt.ParseWithClaims(tokenstr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	emailFromToken, _ := claims["email"].(string)
	var userID uint
	err := db.Model(&User{}).Select("id").Where("email = ?", emailFromToken).Scan(&userID).Error
	if err != nil {
		log.Println("не могу наути пользователя")
		return
	}
	var students []Table_student
	err = db.Model(&Table_student{}).Select("ID,Name_Student").Where("User_id = ?", userID).Scan(&students).Error
	if err != nil {
		log.Println(err)
	}
	c.HTML(http.StatusOK, "paymentstudent.html", gin.H{
		"students": students,
		"User":     user,
	})
} //стр записи об оплате учеников
func resultPage(c *gin.Context) {
	userData, exists := c.Get("User") // Берем пользователя из контекста
	if !exists || userData == nil {
		log.Println("Пользователь не найден в middleware")
		c.HTML(http.StatusUnauthorized, "error.html", gin.H{"error": "Пользователь не авторизован"})
		return
	}
	user := userData.(User)
	emailData, exists := c.Get("email") // Берем пользователя из контекста
	if !exists || emailData == nil {
		log.Println("Пользователь не найден в middleware")
		c.HTML(http.StatusUnauthorized, "error.html", gin.H{"error": "Пользователь не авторизован"})
		return
	}
	email := emailData.(string)
	var userID uint
	err := db.Model(&User{}).Select("ID").Where("email=?", email).Scan(&userID).Error
	if err != nil {
		log.Println("в итоге не смог вытащить id", err)
	}
	var output []Table_student
	err = db.Model(&Table_student{}).Select("ID,Name_Student,Payment,Theory,Practice,Tasks,Namber_lecture").Where("User_id = ?", userID).Find(&output).Error
	if err != nil {
		log.Println("не смог создать таблицу", err)
	}
	var lecture []Table_lecture
	err = db.Model(&Table_lecture{}).Select("Lecture_Person_id,Lecture").Where("User_id=?", userID).Find(&lecture).Error
	lectureMap := make(map[int]string)
	for _, lec := range lecture {
		lectureMap[lec.Lecture_Person_id] = lec.Lecture
	}
	var outputtofront []resultstruct
	for _, el := range output {
		if el.Namber_lecture == 0 {
			outputtofront = append(outputtofront, resultstruct{ID: el.ID, Name_Student: el.Name_Student, Payment: el.Payment, Lecture: "Лекция не выбрана", Theory: el.Theory, Practice: el.Practice, Tasks: el.Tasks})
		} else {
			value, _ := lectureMap[el.Namber_lecture]
			outputtofront = append(outputtofront, resultstruct{ID: el.ID, Name_Student: el.Name_Student, Payment: el.Payment, Lecture: value, Theory: el.Theory, Practice: el.Practice, Tasks: el.Tasks})
		}
	}

	c.HTML(http.StatusOK, "result.html", gin.H{
		"outputtofront": outputtofront,
		"User":          user,
	})
} //вывод таблицы со всеми значениями

// post запросы
func regHandler(c *gin.Context) {
	var user User
	if er := c.ShouldBindJSON(&user); er != nil {
		fmt.Println("ошибка парсинга при реге", er)
	}

	hashPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("проблемы хеширования", err)
	}
	user.Password = string(hashPass)

	verifyToken := utils.GenerationToken()
	user.Verification_token = verifyToken

	err = db.Create(&user).Error
	if err != nil {
		log.Println("Ошибка создания пользователя в бд", err)
	}
	go func() {
		hash := hashIDAndEmail(user.ID, user.Email)
		err = db.Create(&Table_telegram_bot{User_id: user.ID, Hash: hash, Vhod: false}).Error
		if err != nil {
			log.Println("Ошибка создания пользователя в таблице телеграм бота", err)
		}
	}()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Регистрация успешна!Нобходимо подтвердить почту",
	})
	go utils.SendVerificationEmail(user.Email, verifyToken)
}
func autHandler(c *gin.Context) {
	var input User
	if err := c.ShouldBindJSON(&input); err != nil {
		fmt.Println("Ошибка парсинга", err)
	}
	var user User
	if er := db.Where("email=?", input.Email).First(&user).Error; er != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный email или пароль"})
		return
	}

	if !user.Email_verified {
		c.JSON(http.StatusForbidden, gin.H{"error": "Подтвердите email перед входом!"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный email или пароль"})
		return
	}
	token := GenerateJwt(user.Email)
	c.SetCookie("token", token, 3600, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"redirect": "/kabinet",
	})
}
func firstSettingHandler(c *gin.Context) {
	var input PostSettings
	tokenstr, _ := c.Cookie("token")
	claims := jwt.MapClaims{}
	_, _ = jwt.ParseWithClaims(tokenstr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	emailFromToken, _ := claims["email"].(string)
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Println("ошибка парсинга начальной настройки")
	}
	//запись учеников
	if input.Marking == "1" {
		var userID uint
		err := db.Model(&User{}).Select("id").Where("email = ?", emailFromToken).Scan(&userID).Error
		if err != nil {
			fmt.Println("Ошибка:", err)
		}
		db.Create(&Table_student{User_id: userID, Name_Student: input.Meaning, Alert_payment: true, Alert_moduls: true})
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "Ученик сохранен"})
	}
	//запись лекций
	if input.Marking == "0" {
		var lection User
		var num_lectures_int int
		err := db.Where("email =?", emailFromToken).First(&lection).Error
		if err != nil {
			log.Println("ошибка при определении количества лекций")
		}
		lection.Lectures_introduced += 1
		num_lectures_int = lection.Lectures_introduced
		if err = db.Save(&lection).Error; err != nil {
			log.Println("Ошибка при обновлении данных ученика:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Ошибка при сохранении"})
			return
		}
		var UserID uint
		err = db.Model(&User{}).Select("id").Where("email =?", emailFromToken).Scan(&UserID).Error
		if err != nil {
			log.Println("ошибка при определениии UserID", err)
		}
		db.Create(&Table_lecture{User_id: UserID, Lecture: input.Meaning, Lecture_Person_id: num_lectures_int})
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "Лекция сохранена"})
	}
}
func notelessonHandler(c *gin.Context) {
	var input PostModuls
	err := c.ShouldBindJSON(&input)
	if err != nil {
		log.Println(err)
	}
	var student Table_student
	if err = db.Where("id=?", input.Student_id).First(&student).Error; err != nil {
		log.Println("Ученик не найден:", err)
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Ученик не найден"})
		return
	}
	switch input.Module {
	case "Теория":
		student.Theory += 1
		student.Practice = 0
		student.Tasks = 0
		if (student.Namber_lecture != 0) && (input.Lock_lecture == false) {
			student.Namber_lecture += 1
		}
		if student.Theory > 2 && student.Alert_moduls == true {
			message := "Вы провели 3 лекции подряд,не забывайте, что практика и задачи тоже важны!Ученик:"
			messageBot(message, student.Name_Student, student.User_id)
		}
	case "Практика":
		student.Theory = 0
		student.Practice += 1
		student.Tasks = 0
		if student.Practice > 2 && student.Alert_moduls == true {
			message := "Вы провели 3 практики подряд,не забывайте, что теория и задачи тоже важны!Ученик:"
			messageBot(message, student.Name_Student, student.User_id)
		}
	case "Задачи":
		student.Theory = 0
		student.Practice = 0
		student.Tasks += 1
		if student.Tasks > 2 && student.Alert_moduls == true {
			message := "Ученик решает задача 3 урока подряд,не забывайте, что теория и практика тоже важны!Ученик:"
			messageBot(message, student.Name_Student, student.User_id)
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Неверный модуль"})
		return
	}
	student.Payment -= 1
	if student.Payment < 1 && student.Alert_payment == true {
		message := "На балансе ученика не достаточно средств, напомни об оплате за уроки!Ученик:"
		messageBot(message, student.Name_Student, student.User_id)
	}
	if err = db.Save(&student).Error; err != nil {
		log.Println("Ошибка при обновлении данных ученика:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Ошибка при сохранении"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "massage": "Сохранено"})
}
func paymentstudentHandler(c *gin.Context) {
	var input Paymentstudent
	err := c.ShouldBindJSON(&input)
	if err != nil {
		log.Println(err)
	}
	var student Table_student
	if err = db.Where("id = ?", input.ID).First(&student).Error; err != nil {
		log.Println("Ученик не найден:", err)
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "Ученик не найден"})
		return
	}
	student.Payment += input.Payment
	if err = db.Save(&student).Error; err != nil {
		log.Println("Ошибка при обновлении данных ученика:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Ошибка при сохранении"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Оплата сохранена"})
}
func lectureHandler(c *gin.Context) {
	var input []PostLecture
	err := c.ShouldBindJSON(&input)
	if err != nil {
		log.Println("ошибка парсинга json")
	}
	userData, exists := c.Get("User")
	if !exists || userData == nil {
		log.Println("Пользователь не найден в middleware")
		c.HTML(http.StatusUnauthorized, "error.html", gin.H{"error": "Пользователь не авторизован"})
		return
	}
	user := userData.(User)
	var IDShnik []uint
	err = db.Model(&Table_lecture{}).Select("id").Where("User_id=?", user.ID).Find(&IDShnik).Error
	if err != nil {
		log.Println("не смог найти лекций пользователя")
	}
	//только изменеие название или порядка
	if len(IDShnik) == len(input) {
		for idx, el := range IDShnik {
			Element_input := input[idx]
			Lecture_Element := Element_input.Lecture
			err = db.Model(&Table_lecture{}).Where("id=?", el).Update("Lecture", Lecture_Element).Error
			if err != nil {
				log.Println("Проблема с обновлением в бд")
			}
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "massage": "изменения сохранены"})
	}
	//изменение порядка и(или) удаление
	if len(IDShnik) > len(input) {
		for idx, el := range IDShnik {
			if idx < len(input) {
				Element_input := input[idx]
				Lecture_Element := Element_input.Lecture
				err = db.Model(&Table_lecture{}).Where("id=?", el).Updates(map[string]interface{}{
					"Lecture":           Lecture_Element,
					"Lecture_Person_id": idx + 1,
				}).Error
				if err != nil {
					log.Println("Проблема с обновлением в бд при удалении")
				}
			} else {
				err = db.Where("id=?", el).Delete(&Table_lecture{}).Error
				if err != nil {
					log.Println("Не удалось удалить запись", err)
				}
			}
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "измененение сохранено"})
	}
	//изменение порядка и(или) добавление новой лекции
	if len(IDShnik) < len(input) {
		for idx, el := range input {
			if idx < len(IDShnik) {
				err = db.Model(&Table_lecture{}).Where("id=?", IDShnik[idx]).Updates(map[string]interface{}{
					"Lecture":           el.Lecture,
					"Lecture_Person_id": idx + 1,
				}).Error
				if err != nil {
					log.Println("Проблема с обновлением в бд при удалении")
				}
			} else {
				err = db.Create(&Table_lecture{Lecture: el.Lecture, User_id: user.ID, Lecture_Person_id: idx + 1}).Error
				if err != nil {
					log.Println("Не получилось создать лекцию")
				}
			}
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "сохранено"})
	}
}
func studentHandler(c *gin.Context) {
	userData, exists := c.Get("User")
	if !exists || userData == nil {
		log.Println("проблема с вытаскиванием инфы про юзера", exists)
	}
	user, ok := userData.(User)
	if !ok {
		log.Println("Ошибка приведения userData к User")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка авторизации"})
		return
	}
	var students []Table_student
	err := db.Model(&Table_student{}).Select("ID,Name_Student,Namber_lecture").Where("User_id=?", user.ID).Find(&students).Error
	if err != nil {
		log.Println("Ошибка при вытаскивании студента", err)
	}
	var input []PostStudent
	if err = c.ShouldBindJSON(&input); err != nil {
		log.Println("Ошибка при парсинге входных данных", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные входные данные"})
		return
	}
	studentMap := make(map[int]Table_student)
	for _, student := range students {
		studentMap[student.ID] = student
	}
	var updatedInput []PostStudent
	var newStudentIDs []int

	for _, el := range input {
		if el.ID == 0 {
			newStudent := Table_student{User_id: user.ID, Name_Student: el.Name, Namber_lecture: el.Lecture, Alert_payment: true, Alert_moduls: true}
			err = db.Create(&newStudent).Error
			if err != nil {
				log.Println("ошибка при сохранении ученика", err)
			} else {
				updatedInput = append(updatedInput, PostStudent{ID: newStudent.ID, Name: newStudent.Name_Student, Lecture: newStudent.Namber_lecture})
				newStudentIDs = append(newStudentIDs, newStudent.ID)
			}
		} else {
			updatedInput = append(updatedInput, el)
		}
	}

	var existingStudentIDs []int
	for _, student := range students {
		existingStudentIDs = append(existingStudentIDs, student.ID)
	}

	var inputStudentIDs []int
	for _, el := range updatedInput {
		inputStudentIDs = append(inputStudentIDs, el.ID)
	}

	var studentsToDelete []int
	for _, id := range existingStudentIDs {
		if !contains(inputStudentIDs, id) {
			studentsToDelete = append(studentsToDelete, id)
		}
	}

	if len(studentsToDelete) > 0 {
		db.Where("user_id = ? AND id IN (?)", user.ID, studentsToDelete).Delete(&Table_student{})
	}

	for _, el := range updatedInput {
		if existing, found := studentMap[el.ID]; found {
			if el.Lecture != existing.Namber_lecture {
				if err = db.Model(&Table_student{}).Where("ID = ?", el.ID).UpdateColumn("Namber_lecture", el.Lecture).Error; err != nil {
					log.Println("Ошибка при обновлении лекций", err)
				}
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Всё сохранено"})
}
func telbotHandler(c *gin.Context) {
	userData, exists := c.Get("User")
	if !exists || userData == nil {
		log.Println("проблема с вытаскиванием инфы про юзера", exists)
	}
	user := userData.(User)
	var input PostTelbot
	err := c.ShouldBindJSON(&input)
	if err != nil {
		log.Println("Ошибка при парсинге поступления со стр настройки телеграмм бота", err)
	}
	err = db.Model(&Table_telegram_bot{}).Where("User_id=?", user.ID).Update("Vhod", input.ModuleAllToggle).Error
	if err != nil {
		log.Println("ошибка в обнавлении разрешения на отправку уведомлений", err)
	}
	fmt.Println(input.ModuleAllToggle)
	var outStudent_Alertpay []Table_student
	var outStudent_Alertmod []Table_student
	err = db.Model(&Table_student{}).Select("ID,Alert_payment").Where("User_id=?", user.ID).Find(&outStudent_Alertpay).Error
	if err != nil {
		log.Println("стр телеграмбота ошибка с вытаскиванием студентов")
	}
	err = db.Model(&Table_student{}).Select("ID,Alert_moduls").Where("User_id=?", user.ID).Find(&outStudent_Alertmod).Error
	if err != nil {
		log.Println("стр телеграмбота ошибка с вытаскиванием студентов")
	}
	ModulMap := make(map[int]bool)
	for _, modul := range outStudent_Alertmod {
		ModulMap[modul.ID] = modul.Alert_moduls
	}
	PaymentMap := make(map[int]bool)
	for _, Pay := range outStudent_Alertpay {
		PaymentMap[Pay.ID] = Pay.Alert_payment
	}
	for _, el := range input.Students {
		if ModulMap[el.ID] != el.Alert_moduls {
			err = db.Model(&Table_student{}).Where("ID=?", el.ID).Update("Alert_moduls", el.Alert_moduls).Error
			if err != nil {
				log.Println("ошибка в записи уведомлений модулей", err)
			}
		}
	}
	for _, el := range input.Students {
		if PaymentMap[el.ID] != el.Alert_payment {
			err = db.Model(&Table_student{}).Where("ID=?", el.ID).Update("Alert_payment", el.Alert_payment).Error
			if err != nil {
				log.Println("ошибка в записи уведомлений оплаты", err)
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Всё сохранено"})
}
func GenerateJwt(email string) string {
	expirationTime := time.Now().Add(24 * time.Hour) // Токен на 24 часа
	claims := &Claims{
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(jwtSecret)
	return tokenString
}
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("token")
		if err != nil {
			c.Redirect(http.StatusSeeOther, "/authentication")
			c.Abort()
			return
		}

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.Redirect(http.StatusSeeOther, "/authentication")
			c.Abort()
			return
		}
		emailFromToken, ok := claims["email"].(string)
		if !ok {
			c.Redirect(http.StatusSeeOther, "/authentication")
			c.Abort()
			return
		}
		var user User
		err = db.Where("email = ?", emailFromToken).First(&user).Error
		if err != nil {
			log.Println("Ошибка получения пользователя:", err)
			c.Redirect(http.StatusSeeOther, "/authentication")
			c.Abort()
			return
		}
		c.Set("email", emailFromToken)
		c.Set("User", user)
		c.Next()
	}
}
func getProfile(c *gin.Context) {
	email, exists := c.Get("email")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Ошибка авторизации"})
		return
	}

	var user User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Пользователь не найден"})
		return
	}

	c.HTML(http.StatusOK, "kabinet.html", gin.H{"User": user})
}
func logoutHandler(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Вы успешно вышли из системы", "redirect": "/"})
}
func contains(slice []int, item int) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
