package generators

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	WORKDIR = "internal/"
)

func GenerateInitialStructure() {
	projectName, err := getProjectName()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	CreateConfigEnv(projectName)
	CreateDatabaseConnection(projectName)
	CreateLoggers(projectName)
	CreatePagination(projectName)
	CreateAppErrs()
	CreateRoutes()
	CreateFiberRoutes(projectName)
	CreateHandleResponse(projectName)
	CreateValidation()
	CreateMainGo(projectName)
	CreateSrcDir()
	CreateExampleConfig()
	GenerateModules("example")
}

func CreateMainGo(projectName string) {
	file := "main.go"
	if _, err := os.Stat(file); os.IsNotExist(err) {
		destination, err := os.Create(file)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer destination.Close()

		fmt.Fprintf(destination, "package main\n\n")

		fmt.Fprintf(destination, "import (\n")
		fmt.Fprintf(destination, "	\"encoding/json\"\n")
		fmt.Fprintf(destination, "	\"fmt\"\n")
		fmt.Fprintf(destination, "	\"log\"\n")
		fmt.Fprintf(destination, "	\"github.com/gofiber/fiber/v2\"\n")
		fmt.Fprintf(destination, "	\"github.com/gofiber/fiber/v2/middleware/cors\"\n")
		fmt.Fprintf(destination, "	\"github.com/gofiber/fiber/v2/middleware/logger\"\n")
		fmt.Fprintf(destination, "	\"%s/config\"\n", projectName)
		fmt.Fprintf(destination, "	\"%s/database\"\n", projectName)
		fmt.Fprintf(destination, "	\"%s/routes\"\n", projectName)
		fmt.Fprintf(destination, "	\"%s/internal/example\"\n", projectName)
		fmt.Fprintf(destination, ")\n\n")

		fmt.Fprintf(destination, "func main() {\n")
		fmt.Fprintf(destination, "	// Connect to database\n")
		fmt.Fprintf(destination, "	db, err := database.PostgresConnection()\n")
		fmt.Fprintf(destination, "	if err != nil {\n")
		fmt.Fprintf(destination, "		log.Fatal(err)\n")
		fmt.Fprintf(destination, "	}\n\n")

		fmt.Fprintf(destination, "	// Initialize example module\n")
		fmt.Fprintf(destination, "	exampleRepo := example.NewExampleRepository(db)\n")
		fmt.Fprintf(destination, "	exampleService := example.NewExampleService(exampleRepo)\n")
		fmt.Fprintf(destination, "	exampleController := example.NewExampleController(exampleService)\n\n")

		fmt.Fprintf(destination, "	// TODO: Add more modules here...\n\n")

		fmt.Fprintf(destination, "	app := fiber.New(fiber.Config{\n")
		fmt.Fprintf(destination, "		JSONEncoder: json.Marshal,\n")
		fmt.Fprintf(destination, "		JSONDecoder: json.Unmarshal,\n")
		fmt.Fprintf(destination, "	})\n")
		fmt.Fprintf(destination, "	app.Use(logger.New())\n")
		fmt.Fprintf(destination, "	app.Use(cors.New())\n\n")

		fmt.Fprintf(destination, "	// Register routes\n")
		fmt.Fprintf(destination, "	routes.NewFiberRoutes(exampleController).Install(app)\n\n")

		fmt.Fprintf(destination, "	log.Fatal(app.Listen(fmt.Sprintf(\":%%s\", config.Env(\"app.port\"))))\n")
		fmt.Fprintf(destination, "}\n")

		fmt.Println("Created main.go successfully:", file)
	} else {
		fmt.Println("File already exists!", file)
	}
}

func CreateSrcDir() {
	pathFolder := "internal"
	if _, err := os.Stat(pathFolder); os.IsNotExist(err) {
		err := os.Mkdir(pathFolder, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func CreateValidation() {
	pathFolder := "validation"
	if _, err := os.Stat(pathFolder); os.IsNotExist(err) {
		err := os.Mkdir(pathFolder, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	path := pathFolder + "/"
	file := path + "fiber.go"
	var _, err = os.Stat(file)

	if os.IsNotExist(err) {
		destination, err := os.Create(file)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer destination.Close()

		fmt.Fprintf(destination, "package validation\n\n")
		fmt.Fprintf(destination, "import \"github.com/go-playground/validator/v10\"\n\n")
		fmt.Fprintf(destination, "type ErrorResponse struct {\n")
		fmt.Fprintf(destination, "	FailedField string `json:\"failed_field\"`\n")
		fmt.Fprintf(destination, "	Tag         string `json:\"tag\"`\n")
		fmt.Fprintf(destination, "	Value       string `json:\"value\"`\n")
		fmt.Fprintf(destination, "}\n\n")
		fmt.Fprintf(destination, "func ValidateStruct(myStruct interface{}) (string, error) {\n")
		fmt.Fprintf(destination, "	var errorX []*ErrorResponse\n")
		fmt.Fprintf(destination, "	validate := validator.New()\n")
		fmt.Fprintf(destination, "	err := validate.Struct(myStruct)\n")
		fmt.Fprintf(destination, "	if err != nil {\n")
		fmt.Fprintf(destination, "		for _, err := range err.(validator.ValidationErrors) {\n")
		fmt.Fprintf(destination, "			var element ErrorResponse\n")
		fmt.Fprintf(destination, "			element.FailedField = err.Field() + \" \" + err.Tag() + \" \" + err.Param()\n")
		fmt.Fprintf(destination, "			element.Tag = err.Tag()\n")
		fmt.Fprintf(destination, "			element.Value = err.Param()\n")
		fmt.Fprintf(destination, "			errorX = append(errorX, &element)\n")
		fmt.Fprintf(destination, "		}\n")
		fmt.Fprintf(destination, "	}\n")
		fmt.Fprintf(destination, "	if errorX != nil {\n")
		fmt.Fprintf(destination, "		return errorX[0].FailedField, err\n")
		fmt.Fprintf(destination, "	}\n")
		fmt.Fprintf(destination, "	return \"\", nil\n")
		fmt.Fprintf(destination, "}\n")

		fmt.Println("Created Validation successfully:", file)
	} else {
		fmt.Println("File already exists!", file)
	}
}

func CreateDatabaseConnection(projectName string) {
	pathFolder := "database"
	if _, err := os.Stat(pathFolder); os.IsNotExist(err) {
		err := os.Mkdir(pathFolder, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	path := pathFolder + "/"
	file := path + "postgres.go"
	var _, err = os.Stat(file)

	if os.IsNotExist(err) {
		destination, err := os.Create(file)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer destination.Close()

		fmt.Fprintf(destination, "package database\n\n")

		fmt.Fprintf(destination, "import (\n")
		fmt.Fprintf(destination, "\t\"fmt\"\n")
		fmt.Fprintf(destination, "\t\"log\"\n")
		fmt.Fprintf(destination, "\t\"%s/config\"\n", projectName)
		fmt.Fprintf(destination, "\t\"time\"\n\n")
		fmt.Fprintf(destination, "\t\"gorm.io/driver/postgres\"\n")
		fmt.Fprintf(destination, "\t\"gorm.io/gorm\"\n")
		fmt.Fprintf(destination, "\t\"gorm.io/gorm/logger\"\n")
		fmt.Fprintf(destination, ")\n\n")

		fmt.Fprintf(destination, "type SqlLogger struct {\n")
		fmt.Fprintf(destination, "\tlogger.Interface\n")
		fmt.Fprintf(destination, "}\n\n")

		fmt.Fprintf(destination, "var openConnectionDB *gorm.DB\n")
		fmt.Fprintf(destination, "var err error\n\n")

		fmt.Fprintf(destination, "func PostgresConnection() (*gorm.DB, error) {\n")
		fmt.Fprintf(destination, "\tmyDSN := fmt.Sprintf(\"host=%%v user=%%v password=%%v dbname=%%v port=%%v sslmode=disable TimeZone=Asia/Bangkok\",\n")
		fmt.Fprintf(destination, "\t\tconfig.Env(\"postgres.host\"),\n")
		fmt.Fprintf(destination, "\t\tconfig.Env(\"postgres.user\"),\n")
		fmt.Fprintf(destination, "\t\tconfig.Env(\"postgres.password\"),\n")
		fmt.Fprintf(destination, "\t\tconfig.Env(\"postgres.database\"),\n")
		fmt.Fprintf(destination, "\t\tconfig.Env(\"postgres.port\"),\n")
		fmt.Fprintf(destination, "\t)\n\n")

		fmt.Fprintf(destination, "\tfmt.Println(\"CONNECTING_TO_POSTGRES_DB\")\n")
		fmt.Fprintf(destination, "\topenConnectionDB, err = gorm.Open(postgres.Open(myDSN), &gorm.Config{\n")
		fmt.Fprintf(destination, "\t\tNowFunc: func() time.Time {\n")
		fmt.Fprintf(destination, "\t\t\tti, _ := time.LoadLocation(\"Asia/Bangkok\")\n")
		fmt.Fprintf(destination, "\t\t\treturn time.Now().In(ti)\n")
		fmt.Fprintf(destination, "\t\t},\n")
		fmt.Fprintf(destination, "\t})\n")
		fmt.Fprintf(destination, "\tif err != nil {\n")
		fmt.Fprintf(destination, "\t\tlog.Fatal(\"ERROR_PING_POSTGRES\", err)\n")
		fmt.Fprintf(destination, "\t\treturn nil, err\n")
		fmt.Fprintf(destination, "\t}\n")
		fmt.Fprintf(destination, "\tfmt.Println(\"POSTGRES_CONNECTED\")\n")
		fmt.Fprintf(destination, "\treturn openConnectionDB, nil\n")
		fmt.Fprintf(destination, "}\n")

		fmt.Println("Created Database Connection successfully:", file)
	} else {
		fmt.Println("File already exists!", file)
	}
}

func CreateHandleResponse(projectName string) {
	pathFolder := "responses"
	if _, err := os.Stat(pathFolder); os.IsNotExist(err) {
		err := os.Mkdir(pathFolder, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	path := pathFolder + "/"
	file := path + "handle_responses.go"
	var _, err = os.Stat(file)

	if os.IsNotExist(err) {
		destination, err := os.Create(file)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer destination.Close()

		fmt.Fprintf(destination, "package responses\n\n")
		fmt.Fprintf(destination, "import (\n")
		fmt.Fprintf(destination, "	\"github.com/gofiber/fiber/v2\"\n")
		fmt.Fprintf(destination, "	\"%s/errs\"\n", projectName)
		fmt.Fprintf(destination, "	\"net/http\"\n")
		fmt.Fprintf(destination, ")\n\n")
		fmt.Fprintf(destination, "var (\n")
		fmt.Fprintf(destination, "	code    int\n")
		fmt.Fprintf(destination, "	message string\n")
		fmt.Fprintf(destination, ")\n\n")
		fmt.Fprintf(destination, "type ErrorResponse struct {\n")
		fmt.Fprintf(destination, "	Status bool   `json:\"status\"`\n")
		fmt.Fprintf(destination, "	Error  string `json:\"error\"`\n")
		fmt.Fprintf(destination, "}\n\n")
		fmt.Fprintf(destination, "func NewErrorResponses(ctx *fiber.Ctx, err error) error {\n")
		fmt.Fprintf(destination, "	switch e := err.(type) {\n")
		fmt.Fprintf(destination, "	case errs.AppError:\n")
		fmt.Fprintf(destination, "		code = e.Status\n")
		fmt.Fprintf(destination, "		message = e.Message\n")
		fmt.Fprintf(destination, "	case error:\n")
		fmt.Fprintf(destination, "		code = http.StatusUnprocessableEntity\n")
		fmt.Fprintf(destination, "		message = err.Error()\n")
		fmt.Fprintf(destination, "	}\n")
		fmt.Fprintf(destination, "	errorResponse := ErrorResponse{\n")
		fmt.Fprintf(destination, "		Status: false,\n")
		fmt.Fprintf(destination, "		Error:  message,\n")
		fmt.Fprintf(destination, "	}\n")
		fmt.Fprintf(destination, "	return ctx.Status(code).JSON(errorResponse)\n")
		fmt.Fprintf(destination, "}\n\n")
		fmt.Fprintf(destination, "func NewSuccessResponse(ctx *fiber.Ctx, data interface{}) error {\n")
		fmt.Fprintf(destination, "	return ctx.Status(http.StatusOK).JSON(fiber.Map{\n")
		fmt.Fprintf(destination, "		\"status\": true,\n")
		fmt.Fprintf(destination, "		\"data\":   data,\n")
		fmt.Fprintf(destination, "	})\n")
		fmt.Fprintf(destination, "}\n\n")
		fmt.Fprintf(destination, "func NewSuccessMsg(ctx *fiber.Ctx, msg interface{}) error {\n")
		fmt.Fprintf(destination, "	return ctx.Status(http.StatusOK).JSON(fiber.Map{\n")
		fmt.Fprintf(destination, "		\"status\": true,\n")
		fmt.Fprintf(destination, "		\"msg\":    msg,\n")
		fmt.Fprintf(destination, "	})\n")
		fmt.Fprintf(destination, "}\n\n")
		fmt.Fprintf(destination, "func NewCreateSuccessResponse(ctx *fiber.Ctx, data interface{}) error {\n")
		fmt.Fprintf(destination, "	return ctx.Status(http.StatusCreated).JSON(fiber.Map{\n")
		fmt.Fprintf(destination, "		\"status\": true,\n")
		fmt.Fprintf(destination, "		\"data\":   data,\n")
		fmt.Fprintf(destination, "	})\n")
		fmt.Fprintf(destination, "}\n\n")
		fmt.Fprintf(destination, "func NewSuccessMessage(ctx *fiber.Ctx, data interface{}) error {\n")
		fmt.Fprintf(destination, "	return ctx.Status(http.StatusOK).JSON(fiber.Map{\n")
		fmt.Fprintf(destination, "		\"status\":  true,\n")
		fmt.Fprintf(destination, "		\"message\": data,\n")
		fmt.Fprintf(destination, "	})\n")
		fmt.Fprintf(destination, "}\n\n")
		fmt.Fprintf(destination, "func NewErrorValidate(ctx *fiber.Ctx, data interface{}) error {\n")
		fmt.Fprintf(destination, "	validateError := fiber.Map{\n")
		fmt.Fprintf(destination, "		\"error\":  data,\n")
		fmt.Fprintf(destination, "		\"status\": false,\n")
		fmt.Fprintf(destination, "	}\n")
		fmt.Fprintf(destination, "	return ctx.Status(http.StatusUnprocessableEntity).JSON(validateError)\n")
		fmt.Fprintf(destination, "}\n")

		fmt.Println("Created responses package successfully:", file)
	} else {
		fmt.Println("File already exists!", file)
	}
}

func CreateConfigEnv(projectName string) {
	pathFolder := "config"
	if _, err := os.Stat(pathFolder); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(pathFolder, os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
	}

	path := "config/"
	file := path + "env.go"
	var _, err = os.Stat(file)

	if os.IsNotExist(err) {
		destination, err := os.Create(file)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer destination.Close()
		fmt.Fprintf(destination, "package config\n\n")

		fmt.Fprintf(destination, "import (\n")
		fmt.Fprintf(destination, "\t\"fmt\"\n")
		fmt.Fprintf(destination, "\t\"os\"\n")
		fmt.Fprintf(destination, "\t\"strings\"\n\n")
		fmt.Fprintf(destination, "\t\"github.com/spf13/viper\"\n")
		fmt.Fprintf(destination, ")\n\n")

		fmt.Fprintf(destination, "func init() {\n")
		fmt.Fprintf(destination, "\t// Set up Viper for environment variable handling\n")
		fmt.Fprintf(destination, "\tviper.AutomaticEnv()\n")
		fmt.Fprintf(destination, "\tviper.SetEnvKeyReplacer(strings.NewReplacer(\".\", \"_\"))\n\n")
		fmt.Fprintf(destination, "\t// Try loading config.yaml first\n")
		fmt.Fprintf(destination, "\tviper.SetConfigName(\"config\")\n")
		fmt.Fprintf(destination, "\tviper.SetConfigType(\"yaml\")\n")
		fmt.Fprintf(destination, "\tviper.AddConfigPath(\"./\")\n\n")
		fmt.Fprintf(destination, "\terr := viper.ReadInConfig()\n")
		fmt.Fprintf(destination, "\tif err != nil {\n")
		fmt.Fprintf(destination, "\t\t// If config.yaml not found, try loading .env instead\n")
		fmt.Fprintf(destination, "\t\tfmt.Println(\"config.yaml not found, trying .env file\")\n")
		fmt.Fprintf(destination, "\t\tviper.SetConfigName(\".env\")\n")
		fmt.Fprintf(destination, "\t\tviper.SetConfigType(\"env\")\n")
		fmt.Fprintf(destination, "\t\terr = viper.ReadInConfig()\n")
		fmt.Fprintf(destination, "\t\tif err != nil {\n")
		fmt.Fprintf(destination, "\t\t\tfmt.Println(\"ERROR_READING_CONFIG_FILE\", err)\n")
		fmt.Fprintf(destination, "\t\t\treturn\n")
		fmt.Fprintf(destination, "\t\t}\n")
		fmt.Fprintf(destination, "\t}\n")
		fmt.Fprintf(destination, "\tfmt.Println(\"SUCCESS_READING_CONFIG_FILE\")\n")
		fmt.Fprintf(destination, "}\n\n")

		fmt.Fprintf(destination, "func GetEnv(key, defaultValue string) string {\n")
		fmt.Fprintf(destination, "\t// Prioritize environment variables (for Kubernetes Secrets)\n")
		fmt.Fprintf(destination, "\tif val, found := os.LookupEnv(key); found {\n")
		fmt.Fprintf(destination, "\t\treturn val\n")
		fmt.Fprintf(destination, "\t}\n\n")
		fmt.Fprintf(destination, "\t// Fallback to Viper if the environment variable is not set\n")
		fmt.Fprintf(destination, "\treadValue := viper.GetString(key)\n")
		fmt.Fprintf(destination, "\tif readValue == \"\" {\n")
		fmt.Fprintf(destination, "\t\treturn defaultValue\n")
		fmt.Fprintf(destination, "\t}\n")
		fmt.Fprintf(destination, "\treturn readValue\n")
		fmt.Fprintf(destination, "}\n\n")

		fmt.Fprintf(destination, "func Env(key string) string {\n")
		fmt.Fprintf(destination, "\t// Prioritize environment variables (for Kubernetes Secrets)\n")
		fmt.Fprintf(destination, "\tif val, found := os.LookupEnv(key); found {\n")
		fmt.Fprintf(destination, "\t\treturn val\n")
		fmt.Fprintf(destination, "\t}\n")
		fmt.Fprintf(destination, "\treturn viper.GetString(key)\n")
		fmt.Fprintf(destination, "}\n")
	} else {
		fmt.Println("File already exists!", file)
		return
	}

	fmt.Println("Created Config successfully", file)
}

func CreateAppErrs() {
	pathFolder := "errs"
	if _, err := os.Stat(pathFolder); os.IsNotExist(err) {
		err := os.Mkdir(pathFolder, os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
	}

	path := "errs/"
	file := path + "errors.go"
	if _, err := os.Stat(file); os.IsNotExist(err) {
		destination, err := os.Create(file)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer destination.Close()

		fmt.Fprintf(destination, "package errs\n\n")
		fmt.Fprintf(destination, "import \"net/http\"\n\n")
		fmt.Fprintf(destination, "type AppError struct {\n")
		fmt.Fprintf(destination, "	Status  int\n")
		fmt.Fprintf(destination, "	Message string\n")
		fmt.Fprintf(destination, "}\n\n")
		fmt.Fprintf(destination, "func (a AppError) Error() string {\n")
		fmt.Fprintf(destination, "	return a.Message\n")
		fmt.Fprintf(destination, "}\n\n")
		fmt.Fprintf(destination, "func NewError(code int, errMsg string) error {\n")
		fmt.Fprintf(destination, "	return AppError{\n")
		fmt.Fprintf(destination, "		Status:  code,\n")
		fmt.Fprintf(destination, "		Message: errMsg,\n")
		fmt.Fprintf(destination, "	}\n")
		fmt.Fprintf(destination, "}\n\n")
		fmt.Fprintf(destination, "func ErrorBadRequest(errorMessage string) error {\n")
		fmt.Fprintf(destination, "	return AppError{\n")
		fmt.Fprintf(destination, "		Status:  http.StatusBadRequest,\n")
		fmt.Fprintf(destination, "		Message: errorMessage,\n")
		fmt.Fprintf(destination, "	}\n")
		fmt.Fprintf(destination, "}\n\n")
		fmt.Fprintf(destination, "func ErrorUnprocessableEntity(errorMessage string) error {\n")
		fmt.Fprintf(destination, "	return AppError{\n")
		fmt.Fprintf(destination, "		Status:  http.StatusUnprocessableEntity,\n")
		fmt.Fprintf(destination, "		Message: errorMessage,\n")
		fmt.Fprintf(destination, "	}\n")
		fmt.Fprintf(destination, "}\n\n")
		fmt.Fprintf(destination, "func ErrorInternalServerError(errorMessage string) error {\n")
		fmt.Fprintf(destination, "	return AppError{\n")
		fmt.Fprintf(destination, "		Status:  http.StatusInternalServerError,\n")
		fmt.Fprintf(destination, "		Message: errorMessage,\n")
		fmt.Fprintf(destination, "	}\n")
		fmt.Fprintf(destination, "}\n")

		fmt.Println("Created AppErrs successfully", file)
	} else {
		fmt.Println("File already exists!", file)
	}
}

func CreateLoggers(projectName string) {
	pathFolder := "logs"
	if _, err := os.Stat(pathFolder); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(pathFolder, os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
	}

	path := "logs/"
	file := path + "loggers.go"
	var _, err = os.Stat(file)

	if os.IsNotExist(err) {
		destination, err := os.Create(file)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer destination.Close()

		// Write the logging package code to the file
		fmt.Fprintf(destination, "package logs\n\n")

		fmt.Fprintf(destination, "import (\n")
		fmt.Fprintf(destination, "\t\"context\"\n")
		fmt.Fprintf(destination, "\t\"fmt\"\n")
		fmt.Fprintf(destination, "\n")
		fmt.Fprintf(destination, "\t\"go.uber.org/zap\"\n")
		fmt.Fprintf(destination, "\t\"go.uber.org/zap/zapcore\"\n")
		fmt.Fprintf(destination, ")\n\n")

		fmt.Fprintf(destination, "var log *zap.Logger\n")
		fmt.Fprintf(destination, "var err error\n\n")

		fmt.Fprintf(destination, "func init() {\n")
		fmt.Fprintf(destination, "\tconfig := zap.NewProductionConfig()\n")
		fmt.Fprintf(destination, "\tconfig.EncoderConfig.TimeKey = \"timestamp\"\n")
		fmt.Fprintf(destination, "\tconfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder\n")
		fmt.Fprintf(destination, "\tconfig.EncoderConfig.StacktraceKey = \"\"\n")
		fmt.Fprintf(destination, "\tlog, err = config.Build(zap.AddCallerSkip(1))\n\n")
		fmt.Fprintf(destination, "\tconfig.OutputPaths = []string{\"stdout\"} // Change to \"stderr\" if needed\n\n")
		fmt.Fprintf(destination, "\tif err != nil {\n")
		fmt.Fprintf(destination, "\t\tfmt.Println(err)\n")
		fmt.Fprintf(destination, "\t\treturn\n")
		fmt.Fprintf(destination, "\t}\n")
		fmt.Fprintf(destination, "}\n\n")

		fmt.Fprintf(destination, "func Info(message string, ctx context.Context, fields ...zap.Field) {\n")
		fmt.Fprintf(destination, "\trequestId := ctx.Value(\"requestid\")\n")
		fmt.Fprintf(destination, "\tif requestId != nil {\n")
		fmt.Fprintf(destination, "\t\tfields = append(fields, zap.String(\"request_id\", requestId.(string)))\n")
		fmt.Fprintf(destination, "\t}\n")
		fmt.Fprintf(destination, "\tlog.Info(message, fields...)\n")
		fmt.Fprintf(destination, "}\n\n")

		fmt.Fprintf(destination, "func Error(message interface{}, ctx context.Context, fields ...zap.Field) {\n")
		fmt.Fprintf(destination, "\trequestId := ctx.Value(\"requestid\")\n")
		fmt.Fprintf(destination, "\tif requestId != nil {\n")
		fmt.Fprintf(destination, "\t\tfields = append(fields, zap.String(\"request_id\", requestId.(string)))\n")
		fmt.Fprintf(destination, "\t}\n\n")
		fmt.Fprintf(destination, "\tswitch v := message.(type) {\n")
		fmt.Fprintf(destination, "\tcase error:\n")
		fmt.Fprintf(destination, "\t\tlog.Error(v.Error(), fields...)\n")
		fmt.Fprintf(destination, "\tcase string:\n")
		fmt.Fprintf(destination, "\t\tlog.Error(v, fields...)\n")
		fmt.Fprintf(destination, "\t}\n")
		fmt.Fprintf(destination, "}\n\n")

		fmt.Fprintf(destination, "func Debug(message string, fields ...zap.Field) {\n")
		fmt.Fprintf(destination, "\tlog.Debug(message, fields...)\n")
		fmt.Fprintf(destination, "}\n")

		fmt.Println("Created Loggers successfully", file)
	} else {
		fmt.Println("File already exists!", file)
		return
	}
}

func CreatePagination(projectName string) {
	pathFolder := "paginates"
	if _, err := os.Stat(pathFolder); os.IsNotExist(err) {
		err := os.Mkdir(pathFolder, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	path := pathFolder + "/"
	file := path + "pagination.go"
	var _, err = os.Stat(file)

	if os.IsNotExist(err) {
		destination, err := os.Create(file)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer destination.Close()

		fmt.Fprintf(destination, "package paginates\n\n")
		fmt.Fprintf(destination, "import (\n")
		fmt.Fprintf(destination, "\t\"gorm.io/gorm\"\n")
		fmt.Fprintf(destination, "\t\"gorm.io/gorm/clause\"\n")
		fmt.Fprintf(destination, ")\n\n")

		fmt.Fprintf(destination, "type PaginateRequest struct {\n")
		fmt.Fprintf(destination, "\tLimit     int    `json:\"limit\"`\n")
		fmt.Fprintf(destination, "\tPage      int    `json:\"skip\"`\n")
		fmt.Fprintf(destination, "\tStatus    string `json:\"status\"`\n")
		fmt.Fprintf(destination, "\tSearch    string `json:\"search\"`\n")
		fmt.Fprintf(destination, "\tOrderBy   string `json:\"order_by\"`\n")
		fmt.Fprintf(destination, "\tSortBy    string `json:\"sort_by\"`\n")
		fmt.Fprintf(destination, "\tStartDate string `json:\"start_date\"`\n")
		fmt.Fprintf(destination, "\tEndDate   string `json:\"end_date\"`\n")
		fmt.Fprintf(destination, "\tUserID    string `json:\"user_id\"`\n")
		fmt.Fprintf(destination, "}\n\n")

		fmt.Fprintf(destination, "type PaginatedResponse struct {\n")
		fmt.Fprintf(destination, "\tRows         interface{} `json:\"rows\"`\n")
		fmt.Fprintf(destination, "\tCount        int         `json:\"count\"`\n")
		fmt.Fprintf(destination, "\tCountPage    int         `json:\"countPage\"`\n")
		fmt.Fprintf(destination, "\tCurrentPage  int         `json:\"currentPage\"`\n")
		fmt.Fprintf(destination, "\tNextPage     int         `json:\"nextPage\"`\n")
		fmt.Fprintf(destination, "\tPreviousPage int         `json:\"previousPage\"`\n")
		fmt.Fprintf(destination, "}\n\n")

		fmt.Fprintf(destination, "func Paginate(db *gorm.DB, paginate PaginateRequest, resultModel interface{}) (*PaginatedResponse, error) {\n")
		fmt.Fprintf(destination, "\tif paginate.Limit <= 0 {\n")
		fmt.Fprintf(destination, "\t\tpaginate.Limit = 10\n")
		fmt.Fprintf(destination, "\t}\n")
		fmt.Fprintf(destination, "\tvar total int64\n\n")
		fmt.Fprintf(destination, "\tdb.Count(&total)\n")
		fmt.Fprintf(destination, "\tcountPage := (int(total) + paginate.Limit - 1) / paginate.Limit\n")
		fmt.Fprintf(destination, "\toffset := (paginate.Page - 1) * paginate.Limit\n\n")
		fmt.Fprintf(destination, "\tresult := db.Limit(paginate.Limit).\n")
		fmt.Fprintf(destination, "\t\tOffset(offset).Preload(clause.Associations).Find(resultModel)\n")
		fmt.Fprintf(destination, "\tif result.Error != nil {\n")
		fmt.Fprintf(destination, "\t\treturn nil, result.Error\n")
		fmt.Fprintf(destination, "\t}\n\n")
		fmt.Fprintf(destination, "\tnextPage := paginate.Page + 1\n")
		fmt.Fprintf(destination, "\tif nextPage > countPage {\n")
		fmt.Fprintf(destination, "\t\tnextPage = 0\n")
		fmt.Fprintf(destination, "\t}\n\n")
		fmt.Fprintf(destination, "\tpreviousPage := paginate.Page - 1\n")
		fmt.Fprintf(destination, "\tif previousPage < 1 {\n")
		fmt.Fprintf(destination, "\t\tpreviousPage = 0\n")
		fmt.Fprintf(destination, "\t}\n\n")
		fmt.Fprintf(destination, "\tpagination := &PaginatedResponse{\n")
		fmt.Fprintf(destination, "\t\tCount:        int(total),\n")
		fmt.Fprintf(destination, "\t\tCountPage:    countPage,\n")
		fmt.Fprintf(destination, "\t\tCurrentPage:  paginate.Page,\n")
		fmt.Fprintf(destination, "\t\tNextPage:     nextPage,\n")
		fmt.Fprintf(destination, "\t\tPreviousPage: previousPage,\n")
		fmt.Fprintf(destination, "\t\tRows:         resultModel,\n")
		fmt.Fprintf(destination, "\t}\n")
		fmt.Fprintf(destination, "\treturn pagination, nil\n")
		fmt.Fprintf(destination, "}\n")

		fmt.Println("Created Pagination successfully:", file)
	} else {
		fmt.Println("File already exists!", file)
	}
}

func CreateRoutes() {
	pathFolder := "routes"
	if _, err := os.Stat(pathFolder); os.IsNotExist(err) {
		err := os.Mkdir(pathFolder, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	path := pathFolder + "/"
	file := path + "routes.go"
	var _, err = os.Stat(file)

	if os.IsNotExist(err) {
		destination, err := os.Create(file)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer destination.Close()

		fmt.Fprintf(destination, "package routes\n\n")
		fmt.Fprintf(destination, "import \"github.com/gofiber/fiber/v2\"\n\n")
		fmt.Fprintf(destination, "type Routes interface {\n")
		fmt.Fprintf(destination, "	Install(app *fiber.App)\n")
		fmt.Fprintf(destination, "}\n")

		fmt.Println("Created Routes successfully:", file)
	} else {
		fmt.Println("File already exists!", file)
	}
}

func CreateFiberRoutes(projectName string) {
	pathFolder := "routes"
	if _, err := os.Stat(pathFolder); os.IsNotExist(err) {
		err := os.Mkdir(pathFolder, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	file := pathFolder + "/fiber_routes.go"
	if _, err := os.Stat(file); os.IsNotExist(err) {
		destination, err := os.Create(file)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer destination.Close()

		fmt.Fprintf(destination, "package routes\n\n")

		fmt.Fprintf(destination, "import (\n")
		fmt.Fprintf(destination, "\t\"github.com/gofiber/fiber/v2\"\n")
		fmt.Fprintf(destination, "\t\"%s/internal/example\"\n", projectName)
		fmt.Fprintf(destination, ")\n\n")

		fmt.Fprintf(destination, "type fiberRoutes struct {\n")
		fmt.Fprintf(destination, "\texampleController example.ExampleController\n")
		fmt.Fprintf(destination, "}\n\n")

		fmt.Fprintf(destination, "func (r fiberRoutes) Install(app *fiber.App) {\n")
		fmt.Fprintf(destination, "\troute := app.Group(\"/api/\", func(ctx *fiber.Ctx) error {\n")
		fmt.Fprintf(destination, "\t\treturn ctx.Next()\n")
		fmt.Fprintf(destination, "\t})\n")
		fmt.Fprintf(destination, "\troute.Get(\"ping\", r.exampleController.PingController)\n")
		fmt.Fprintf(destination, "}\n\n")

		fmt.Fprintf(destination, "func NewFiberRoutes(\n")
		fmt.Fprintf(destination, "\texampleController example.ExampleController,\n")
		fmt.Fprintf(destination, ") Routes {\n")
		fmt.Fprintf(destination, "\treturn &fiberRoutes{\n")
		fmt.Fprintf(destination, "\t\texampleController: exampleController,\n")
		fmt.Fprintf(destination, "\t}\n")
		fmt.Fprintf(destination, "}\n")

		fmt.Println("Created fiber_routes.go successfully:", file)
	} else {
		fmt.Println("File already exists!", file)
	}
}

func GenerateModules(filename string) {
	filename = strings.ToLower(filename)

	projectName, err := getProjectName()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	CreateRequests(filename)
	CreateResponses(filename)
	CreateModels(filename)
	CreateRepositories(filename, projectName)
	CreateServices(filename, projectName)
	CreateControllers(filename, projectName)
}

func CreateRequests(filename string) {
	pathFolder := WORKDIR + filename
	if _, err := os.Stat(pathFolder); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(pathFolder, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	file := pathFolder + "/" + filename + "_request.go"
	if _, err := os.Stat(file); os.IsNotExist(err) {
		destination, err := os.Create(file)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer destination.Close()

		fmt.Fprintf(destination, "package %s\n", filename)
		fmt.Println("Created Request successfully", file)
	} else {
		fmt.Println("File already exists!", file)
	}
}

func CreateResponses(filename string) {
	pathFolder := WORKDIR + filename
	if _, err := os.Stat(pathFolder); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(pathFolder, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	file := pathFolder + "/" + filename + "_response.go"
	if _, err := os.Stat(file); os.IsNotExist(err) {
		destination, err := os.Create(file)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer destination.Close()

		fmt.Fprintf(destination, "package %s\n", filename)
		fmt.Println("Created Response successfully", file)
	} else {
		fmt.Println("File already exists!", file)
	}
}

func CreateModels(filename string) {
	pathFolder := WORKDIR + filename
	if _, err := os.Stat(pathFolder); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(pathFolder, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	file := pathFolder + "/" + filename + "_model.go"
	if _, err := os.Stat(file); os.IsNotExist(err) {
		destination, err := os.Create(file)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer destination.Close()

		upperString := strings.Replace(
			cases.Title(language.Und, cases.NoLower).String(strings.ReplaceAll(filename, "_", " ")),
			" ", "", -1,
		)

		fmt.Fprintf(destination, "package %s\n\n", filename)
		fmt.Fprintf(destination, "import \"gorm.io/gorm\"\n\n")
		fmt.Fprintf(destination, "type %s struct {\n", upperString)
		fmt.Fprintf(destination, "\tgorm.Model\n")
		fmt.Fprintf(destination, "}\n")
		fmt.Println("Created Model successfully", file)
	} else {
		fmt.Println("File already exists!", file)
	}
}

func CreateRepositories(filename string, projectName string) {
	pathFolder := WORKDIR + filename
	if _, err := os.Stat(pathFolder); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(pathFolder, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	file := pathFolder + "/" + filename + "_repository.go"
	if _, err := os.Stat(file); os.IsNotExist(err) {
		destination, err := os.Create(file)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer destination.Close()

		upper := strings.Replace(
			cases.Title(language.Und, cases.NoLower).String(strings.ReplaceAll(filename, "_", " ")),
			" ", "", -1,
		)
		lower := strings.ToLower(string(upper[0])) + upper[1:]

		fmt.Fprintf(destination, "package %s\n\n", filename)
		fmt.Fprintf(destination, "import (\n")
		fmt.Fprintf(destination, "\t\"gorm.io/gorm\"\n")
		fmt.Fprintf(destination, ")\n\n")

		fmt.Fprintf(destination, "type %sRepository interface {\n", upper)
		fmt.Fprintf(destination, "\t// Insert your function interface\n")
		fmt.Fprintf(destination, "}\n\n")

		fmt.Fprintf(destination, "type %sRepository struct {\n", lower)
		fmt.Fprintf(destination, "\tdb *gorm.DB\n")
		fmt.Fprintf(destination, "}\n\n")

		fmt.Fprintf(destination, "func New%sRepository(db *gorm.DB) %sRepository {\n", upper, upper)
		fmt.Fprintf(destination, "\tdb.AutoMigrate(%s{})\n", upper) // referencing the model directly
		fmt.Fprintf(destination, "\treturn &%sRepository{db: db}\n", lower)
		fmt.Fprintf(destination, "}\n")

		fmt.Println("Created Repository successfully", file)
	} else {
		fmt.Println("File already exists!", file)
	}
}

func CreateServices(filename string, projectName string) {
	pathFolder := WORKDIR + filename
	if _, err := os.Stat(pathFolder); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(pathFolder, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	file := pathFolder + "/" + filename + "_service.go"
	if _, err := os.Stat(file); os.IsNotExist(err) {
		destination, err := os.Create(file)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer destination.Close()

		upper := strings.Replace(
			cases.Title(language.Und, cases.NoLower).String(strings.ReplaceAll(filename, "_", " ")),
			" ", "", -1,
		)
		lower := strings.ToLower(string(upper[0])) + upper[1:]

		fmt.Fprintf(destination, "package %s\n\n", filename)

		fmt.Fprintf(destination, "type %sService interface {\n", upper)
		fmt.Fprintf(destination, "\t// Insert your function interface\n")
		fmt.Fprintf(destination, "}\n\n")

		fmt.Fprintf(destination, "type %sService struct {\n", lower)
		fmt.Fprintf(destination, "\trepo %sRepository\n", upper)
		fmt.Fprintf(destination, "}\n\n")

		fmt.Fprintf(destination, "func New%sService(repo %sRepository) %sService {\n", upper, upper, upper)
		fmt.Fprintf(destination, "\treturn &%sService{repo: repo}\n", lower)
		fmt.Fprintf(destination, "}\n")

		fmt.Println("Created Service successfully", file)
	} else {
		fmt.Println("File already exists!", file)
	}
}

func CreateControllers(filename string, projectName string) {
	pathFolder := WORKDIR + filename
	if _, err := os.Stat(pathFolder); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(pathFolder, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	file := pathFolder + "/" + filename + "_controller.go"
	if _, err := os.Stat(file); os.IsNotExist(err) {
		destination, err := os.Create(file)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer destination.Close()

		upper := strings.Replace(
			cases.Title(language.Und, cases.NoLower).String(strings.ReplaceAll(filename, "_", " ")),
			" ", "", -1,
		)
		lower := strings.ToLower(string(upper[0])) + upper[1:]

		fmt.Fprintf(destination, "package %s\n\n", filename)
		fmt.Fprintf(destination, "import (\n")
		fmt.Fprintf(destination, "\t\"github.com/gofiber/fiber/v2\"\n")
		fmt.Fprintf(destination, ")\n\n")

		fmt.Fprintf(destination, "type %sController interface {\n", upper)
		fmt.Fprintf(destination, "\tPingController(ctx *fiber.Ctx) error\n")
		fmt.Fprintf(destination, "}\n\n")

		fmt.Fprintf(destination, "type %sController struct {\n", lower)
		fmt.Fprintf(destination, "\tservice %sService\n", upper) // use the local service type without import
		fmt.Fprintf(destination, "}\n\n")

		fmt.Fprintf(destination, "func New%sController(service %sService) %sController {\n", upper, upper, upper)
		fmt.Fprintf(destination, "\treturn &%sController{service: service}\n", lower)
		fmt.Fprintf(destination, "}\n\n")

		fmt.Fprintf(destination, "func (c *%sController) PingController(ctx *fiber.Ctx) error {\n", lower)
		fmt.Fprintf(destination, "\treturn ctx.JSON(fiber.Map{\n")
		fmt.Fprintf(destination, "\t\t\"message\": \"pong\",\n")
		fmt.Fprintf(destination, "\t})\n")
		fmt.Fprintf(destination, "}\n")

		fmt.Println("Created Controller successfully", file)
	} else {
		fmt.Println("File already exists!", file)
	}
}

// getProjectName reads the project name from the go.mod file
func getProjectName() (string, error) {
	file, err := os.Open("go.mod")
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "module") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				return parts[1], nil
			}
		}
	}

	return "", errors.New("could not determine module name")
}

func CreateExampleConfig() {
	pathFolder := "./"
	if _, err := os.Stat(pathFolder); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(pathFolder, os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
	}

	path := pathFolder
	file := path + "example.config.yaml"
	var _, err = os.Stat(file)

	if os.IsNotExist(err) {
		destination, err := os.Create(file)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer destination.Close()

		// Write the logging package code to the file
		fmt.Fprintf(destination, "app:\n")
		fmt.Fprintf(destination, "  port: 8080\n")
		fmt.Fprintf(destination, "\n")
		fmt.Fprintf(destination, "secrete:\n")
		fmt.Fprintf(destination, "  jwt: \"secrete\"\n")
		fmt.Fprintf(destination, "\n")
		fmt.Fprintf(destination, "postgres:\n")
		fmt.Fprintf(destination, "  host: localhost\n")
		fmt.Fprintf(destination, "  port: 5432\n")
		fmt.Fprintf(destination, "  user: postgres\n")
		fmt.Fprintf(destination, "  password: postgres\n")
		fmt.Fprintf(destination, "  database: postgresdb\n")
		fmt.Fprintf(destination, "\n")
		fmt.Fprintf(destination, "redis:\n")
		fmt.Fprintf(destination, "  host: localhost\n")
		fmt.Fprintf(destination, "  port: 6479\n")

		fmt.Println("Created Example Config successfully", file)
	} else {
		fmt.Println("File already exists!", file)
		return
	}
}
