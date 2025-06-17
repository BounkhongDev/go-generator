package generators

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"sort"
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
	CreateConfigTimezonse(projectName)
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
	CreateMiddleware(projectName)
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
		fmt.Fprintf(destination, "	\"%s/migrations\"\n", projectName)
		fmt.Fprintf(destination, ")\n\n")

		fmt.Fprintf(destination, "func main() {\n")
		fmt.Fprintf(destination, "	// Connect to database\n")
		fmt.Fprintf(destination, "	db, err := database.PostgresConnection()\n")
		fmt.Fprintf(destination, "	if err != nil {\n")
		fmt.Fprintf(destination, "		log.Fatal(err)\n")
		fmt.Fprintf(destination, "	}\n\n")

		fmt.Fprintf(destination, "	// Run migrations\n")
		fmt.Fprintf(destination, "	if err := migrations.MigrateAll(db); err != nil {\n")
		fmt.Fprintf(destination, "		log.Fatalf(\"Migration failed: %%v\", err)\n")
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
		fmt.Println("⚠️  File already exists:", file)
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
	if _, err := os.Stat(file); os.IsNotExist(err) {
		destination, err := os.Create(file)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer destination.Close()

		fmt.Fprintf(destination, "package responses\n\n")
		fmt.Fprintf(destination, "import (\n")
		fmt.Fprintf(destination, "	\"net/http\"\n")
		fmt.Fprintf(destination, "	\"%s/errs\"\n", projectName)
		fmt.Fprintf(destination, "	\"github.com/gofiber/fiber/v2\"\n")
		fmt.Fprintf(destination, ")\n\n")

		fmt.Fprintf(destination, "// ValidationError represents a single validation error.\n")
		fmt.Fprintf(destination, "type ValidationError struct {\n")
		fmt.Fprintf(destination, "	Field   string `json:\"field\"`\n")
		fmt.Fprintf(destination, "	Message string `json:\"message\"`\n")
		fmt.Fprintf(destination, "}\n\n")

		fmt.Fprintf(destination, "// APIResponse defines the standard response structure.\n")
		fmt.Fprintf(destination, "type APIResponse struct {\n")
		fmt.Fprintf(destination, "	Success bool        `json:\"success\"`\n")
		fmt.Fprintf(destination, "	Message string      `json:\"message\"`\n")
		fmt.Fprintf(destination, "	Data    interface{} `json:\"data\"`\n")
		fmt.Fprintf(destination, "	Errors  interface{} `json:\"errors\"`\n")
		fmt.Fprintf(destination, "}\n\n")

		fmt.Fprintf(destination, "// NewErrorResponse handles application and system errors.\n")
		fmt.Fprintf(destination, "func NewErrorResponse(ctx *fiber.Ctx, err error) error {\n")
		fmt.Fprintf(destination, "	var code int\n")
		fmt.Fprintf(destination, "	var message string\n")
		fmt.Fprintf(destination, "	switch e := err.(type) {\n")
		fmt.Fprintf(destination, "	case errs.AppError:\n")
		fmt.Fprintf(destination, "		code = e.Status\n")
		fmt.Fprintf(destination, "		message = e.Message\n")
		fmt.Fprintf(destination, "	default:\n")
		fmt.Fprintf(destination, "		code = http.StatusUnprocessableEntity\n")
		fmt.Fprintf(destination, "		message = err.Error()\n")
		fmt.Fprintf(destination, "	}\n")
		fmt.Fprintf(destination, "	return ctx.Status(code).JSON(APIResponse{\n")
		fmt.Fprintf(destination, "		Success: false,\n")
		fmt.Fprintf(destination, "		Message: message,\n")
		fmt.Fprintf(destination, "		Data:    nil,\n")
		fmt.Fprintf(destination, "		Errors:  nil,\n")
		fmt.Fprintf(destination, "	})\n")
		fmt.Fprintf(destination, "}\n\n")

		fmt.Fprintf(destination, "// NewValidationError sends validation error response.\n")
		fmt.Fprintf(destination, "func NewValidationError(ctx *fiber.Ctx, errors []ValidationError) error {\n")
		fmt.Fprintf(destination, "	return ctx.Status(http.StatusUnprocessableEntity).JSON(APIResponse{\n")
		fmt.Fprintf(destination, "		Success: false,\n")
		fmt.Fprintf(destination, "		Message: \"Validation failed\",\n")
		fmt.Fprintf(destination, "		Data:    nil,\n")
		fmt.Fprintf(destination, "		Errors:  errors,\n")
		fmt.Fprintf(destination, "	})\n")
		fmt.Fprintf(destination, "}\n\n")

		fmt.Fprintf(destination, "// NewSuccessResponse sends a 200 OK response with data.\n")
		fmt.Fprintf(destination, "func NewSuccessResponse(ctx *fiber.Ctx, message string, data interface{}) error {\n")
		fmt.Fprintf(destination, "	return ctx.Status(http.StatusOK).JSON(APIResponse{\n")
		fmt.Fprintf(destination, "		Success: true,\n")
		fmt.Fprintf(destination, "		Message: message,\n")
		fmt.Fprintf(destination, "		Data:    data,\n")
		fmt.Fprintf(destination, "		Errors:  nil,\n")
		fmt.Fprintf(destination, "	})\n")
		fmt.Fprintf(destination, "}\n\n")

		fmt.Fprintf(destination, "// NewCreatedResponse sends a 201 Created response with data.\n")
		fmt.Fprintf(destination, "func NewCreatedResponse(ctx *fiber.Ctx, message string, data interface{}) error {\n")
		fmt.Fprintf(destination, "	return ctx.Status(http.StatusCreated).JSON(APIResponse{\n")
		fmt.Fprintf(destination, "		Success: true,\n")
		fmt.Fprintf(destination, "		Message: message,\n")
		fmt.Fprintf(destination, "		Data:    data,\n")
		fmt.Fprintf(destination, "		Errors:  nil,\n")
		fmt.Fprintf(destination, "	})\n")
		fmt.Fprintf(destination, "}\n")

		fmt.Println("Created responses package successfully:", file)
	} else {
		fmt.Println("⚠️ File already exists:", file)
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

func CreateConfigTimezonse(projectName string) {
	pathFolder := "config"
	if _, err := os.Stat(pathFolder); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(pathFolder, os.ModePerm)
		if err != nil {
			fmt.Println("Failed to create folder:", err)
			return
		}
	}

	file := pathFolder + "/timezone.go"
	if _, err := os.Stat(file); os.IsNotExist(err) {
		destination, err := os.Create(file)
		if err != nil {
			fmt.Println("Failed to create file:", err)
			return
		}
		defer destination.Close()

		content := `package config

import (
	"log"
	"time"
	_ "time/tzdata"
)

func init() {
	location, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		log.Fatal("ERROR_LOADING_TIMEZONE", err)
	}
	time.Local = location
}
`
		_, err = destination.WriteString(content)
		if err != nil {
			fmt.Println("Failed to write content:", err)
			return
		}

		fmt.Println("Created Config Timezone successfully:", file)
	} else {
		fmt.Println("File already exists:", file)
	}
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
	if _, err := os.Stat(file); os.IsNotExist(err) {
		destination, err := os.Create(file)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer destination.Close()

		fmt.Fprintf(destination, "package paginates\n\n")
		fmt.Fprintf(destination, "import (\n")
		fmt.Fprintf(destination, "	\"gorm.io/gorm\"\n")
		fmt.Fprintf(destination, "	\"gorm.io/gorm/clause\"\n")
		fmt.Fprintf(destination, ")\n\n")

		fmt.Fprintf(destination, "// PaginateRequest holds pagination and filter parameters\n")
		fmt.Fprintf(destination, "type PaginateRequest struct {\n")
		fmt.Fprintf(destination, "	Limit     int    `json:\"limit\"`\n")
		fmt.Fprintf(destination, "	Page      int    `json:\"page\"`\n")
		fmt.Fprintf(destination, "	Status    string `json:\"status\"`\n")
		fmt.Fprintf(destination, "	Search    string `json:\"search\"`\n")
		fmt.Fprintf(destination, "	OrderBy   string `json:\"order_by\"`\n")
		fmt.Fprintf(destination, "	SortBy    string `json:\"sort_by\"`\n")
		fmt.Fprintf(destination, "	StartDate string `json:\"start_date\"`\n")
		fmt.Fprintf(destination, "	EndDate   string `json:\"end_date\"`\n")
		fmt.Fprintf(destination, "	UserID    string `json:\"user_id\"`\n")
		fmt.Fprintf(destination, "}\n\n")

		fmt.Fprintf(destination, "// PaginatedResponse is the inner data field of a paginated API response\n")
		fmt.Fprintf(destination, "type PaginatedResponse struct {\n")
		fmt.Fprintf(destination, "	TotalItems    int         `json:\"total_items\"`\n")
		fmt.Fprintf(destination, "	ItemsPerPage  int         `json:\"items_per_page\"`\n")
		fmt.Fprintf(destination, "	CurrentPage   int         `json:\"current_page\"`\n")
		fmt.Fprintf(destination, "	TotalPages    int         `json:\"total_pages\"`\n")
		fmt.Fprintf(destination, "	NextPage      int         `json:\"next_page\"`\n")
		fmt.Fprintf(destination, "	PreviousPage  *int        `json:\"previous_page\"`\n")
		fmt.Fprintf(destination, "	Rows          interface{} `json:\"rows\"`\n")
		fmt.Fprintf(destination, "}\n\n")

		fmt.Fprintf(destination, "// Paginate applies limit, offset, and preload to the DB query\n")
		fmt.Fprintf(destination, "func Paginate(db *gorm.DB, paginate PaginateRequest, resultModel interface{}) (*PaginatedResponse, error) {\n")
		fmt.Fprintf(destination, "	if paginate.Limit <= 0 {\n")
		fmt.Fprintf(destination, "		paginate.Limit = 10\n")
		fmt.Fprintf(destination, "	}\n")
		fmt.Fprintf(destination, "	if paginate.Page <= 0 {\n")
		fmt.Fprintf(destination, "		paginate.Page = 1\n")
		fmt.Fprintf(destination, "	}\n\n")

		fmt.Fprintf(destination, "	var total int64\n")
		fmt.Fprintf(destination, "	db.Count(&total)\n")
		fmt.Fprintf(destination, "	totalPages := (int(total) + paginate.Limit - 1) / paginate.Limit\n")
		fmt.Fprintf(destination, "	offset := (paginate.Page - 1) * paginate.Limit\n\n")

		fmt.Fprintf(destination, "	result := db.Limit(paginate.Limit).\n")
		fmt.Fprintf(destination, "		Offset(offset).Preload(clause.Associations).Find(resultModel)\n")
		fmt.Fprintf(destination, "	if result.Error != nil {\n")
		fmt.Fprintf(destination, "		return nil, result.Error\n")
		fmt.Fprintf(destination, "	}\n\n")

		fmt.Fprintf(destination, "	nextPage := paginate.Page + 1\n")
		fmt.Fprintf(destination, "	if nextPage > totalPages {\n")
		fmt.Fprintf(destination, "		nextPage = 0\n")
		fmt.Fprintf(destination, "	}\n\n")

		fmt.Fprintf(destination, "	var previousPage *int\n")
		fmt.Fprintf(destination, "	if paginate.Page > 1 {\n")
		fmt.Fprintf(destination, "		prev := paginate.Page - 1\n")
		fmt.Fprintf(destination, "		previousPage = &prev\n")
		fmt.Fprintf(destination, "	}\n\n")

		fmt.Fprintf(destination, "	pagination := &PaginatedResponse{\n")
		fmt.Fprintf(destination, "		TotalItems:   int(total),\n")
		fmt.Fprintf(destination, "		ItemsPerPage: paginate.Limit,\n")
		fmt.Fprintf(destination, "		CurrentPage:  paginate.Page,\n")
		fmt.Fprintf(destination, "		TotalPages:   totalPages,\n")
		fmt.Fprintf(destination, "		NextPage:     nextPage,\n")
		fmt.Fprintf(destination, "		PreviousPage: previousPage,\n")
		fmt.Fprintf(destination, "		Rows:         resultModel,\n")
		fmt.Fprintf(destination, "	}\n")
		fmt.Fprintf(destination, "	return pagination, nil\n")
		fmt.Fprintf(destination, "}\n")

		fmt.Println("Created pagination helper successfully:", file)
	} else {
		fmt.Println("⚠️ File already exists:", file)
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
	CreateTestsStructure(filename, projectName)
	CreateMigrations(filename, projectName)
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
		// lower := strings.ToLower(string(upper[0])) + upper[1:]

		fmt.Fprintf(destination, "package %s\n\n", filename)
		fmt.Fprintf(destination, "import (\n")
		fmt.Fprintf(destination, "\t\"github.com/gofiber/fiber/v2\"\n")
		fmt.Fprintf(destination, ")\n\n")

		// fmt.Fprintf(destination, "type %sController interface {\n", upper)
		// fmt.Fprintf(destination, "\tPingController(ctx *fiber.Ctx) error\n")
		// fmt.Fprintf(destination, "}\n\n")

		fmt.Fprintf(destination, "type %sController struct {\n", upper)
		fmt.Fprintf(destination, "\tservice %sService\n", upper) // use the local service type without import
		fmt.Fprintf(destination, "}\n\n")

		fmt.Fprintf(destination, "func New%sController(service %sService) %sController {\n", upper, upper, upper)
		fmt.Fprintf(destination, "\treturn %sController{service: service}\n", upper)
		fmt.Fprintf(destination, "}\n\n")

		fmt.Fprintf(destination, "func (c *%sController) PingController(ctx *fiber.Ctx) error {\n", upper)
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

func CreateTestsStructure(filename string, projectName string) {
	testFolder := "tests/" + filename

	if err := os.MkdirAll(testFolder, os.ModePerm); err != nil {
		fmt.Println("Error creating test directory:", err)
		return
	}

	upper := strings.Replace(
		cases.Title(language.Und, cases.NoLower).String(strings.ReplaceAll(filename, "_", " ")),
		" ", "", -1,
	)
	lower := strings.ToLower(upper) // Fixed the undefined error by initializing lower

	createFile(testFolder+"/"+filename+"_service_test.go", fmt.Sprintf(`package %s

	import (
		"testing"
		"github.com/stretchr/testify/assert"
	)

	func Test%sService(t *testing.T) {
		// TODO: Write tests for %s service
		assert.True(t, true)
	}`, filename, upper, lower))

	createFile(testFolder+"/"+filename+"_repository_mock.go", fmt.Sprintf(`package %s

	import "github.com/stretchr/testify/mock"

	type %sRepositoryMock struct {
		mock.Mock
	}

	// TODO: Add mock implementations
	`, filename, upper))

	fmt.Println("Test structure created successfully at", testFolder)
}

func createFile(path, content string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
		defer file.Close()
		file.WriteString(content)
		fmt.Println("Created file:", path)
	} else {
		fmt.Println("File already exists:", path)
	}
}

func CreateMiddleware(projectName string) {
	pathFolder := "middleware"
	if _, err := os.Stat(pathFolder); os.IsNotExist(err) {
		err := os.Mkdir(pathFolder, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	file := pathFolder + "/logging.go"
	if _, err := os.Stat(file); os.IsNotExist(err) {
		destination, err := os.Create(file)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer destination.Close()

		fmt.Fprintf(destination, "package middleware\n\n")

		fmt.Fprintf(destination, "import (\n")
		fmt.Fprintf(destination, "\t\"%s/logs\"\n", projectName)
		fmt.Fprintf(destination, "\t\"time\"\n")
		fmt.Fprintf(destination, "\n")
		fmt.Fprintf(destination, "\t\"github.com/gofiber/fiber/v2\"\n")
		fmt.Fprintf(destination, "\t\"go.uber.org/zap\"\n")
		fmt.Fprintf(destination, "\t\"go.uber.org/zap/zapcore\"\n")
		fmt.Fprintf(destination, ")\n\n")

		fmt.Fprintf(destination, "func LogInfo(ctx *fiber.Ctx) error {\n")
		fmt.Fprintf(destination, "\tbeginTime := time.Now()\n")
		fmt.Fprintf(destination, "\treqHeader := ctx.GetReqHeaders()\n")
		fmt.Fprintf(destination, "\treqPath := ctx.OriginalURL()\n\n")

		fmt.Fprintf(destination, "\tif err := ctx.Next(); err != nil {\n")
		fmt.Fprintf(destination, "\t\treturn err\n")
		fmt.Fprintf(destination, "\t}\n\n")

		fmt.Fprintf(destination, "\tlatency := time.Since(beginTime)\n")
		fmt.Fprintf(destination, "\tresHeader := ctx.GetRespHeaders()\n\n")

		fmt.Fprintf(destination, "\tlogs.Info(\n")
		fmt.Fprintf(destination, "\t\t\"\",\n")
		fmt.Fprintf(destination, "\t\tctx.Context(),\n")
		fmt.Fprintf(destination, "\t\tzapcore.Field{\n")
		fmt.Fprintf(destination, "\t\t\tKey:    \"latency\",\n")
		fmt.Fprintf(destination, "\t\t\tType:   zapcore.DurationType,\n")
		fmt.Fprintf(destination, "\t\t\tString: latency.String(),\n")
		fmt.Fprintf(destination, "\t\t},\n")
		fmt.Fprintf(destination, "\t\tzap.Any(\"req_header\", reqHeader),\n")
		fmt.Fprintf(destination, "\t\tzap.Any(\"req_path\", reqPath),\n")
		fmt.Fprintf(destination, "\t\tzap.Any(\"res_header\", resHeader),\n")
		fmt.Fprintf(destination, "\t\tzap.Any(\"requester_ip\", ctx.IP()),\n")
		fmt.Fprintf(destination, "\t)\n\n")

		fmt.Fprintf(destination, "\treturn nil\n")
		fmt.Fprintf(destination, "}\n")

		fmt.Println("Created logging.go successfully:", file)
	} else {
		fmt.Println("File already exists!", file)
	}
}

func CreateMigrations(filename string, projectName string) {
	migrationDir := "migrations/"
	filePath := migrationDir + "migrations.go"

	// Ensure migrations directory exists
	if _, err := os.Stat(migrationDir); os.IsNotExist(err) {
		err := os.MkdirAll(migrationDir, os.ModePerm)
		if err != nil {
			fmt.Println("Failed to create migrations directory:", err)
			return
		}
	}

	// Format model name: "user_account" -> "UserAccount"
	modelName := strings.Replace(
		cases.Title(language.Und, cases.NoLower).String(strings.ReplaceAll(filename, "_", " ")),
		" ", "", -1,
	)

	// Build import and model reference
	importPath := fmt.Sprintf("%s/internal/%s", projectName, filename)
	modelLine := fmt.Sprintf("&%s.%s{}", filename, modelName)

	// Maps for uniqueness
	imports := map[string]struct{}{"gorm.io/gorm": {}} // pre-add gorm
	models := map[string]struct{}{}

	// Parse existing file if present
	if content, err := os.ReadFile(filePath); err == nil {
		lines := strings.Split(string(content), "\n")
		inImport := false
		inMigrate := false

		for _, line := range lines {
			trim := strings.TrimSpace(line)

			if trim == "import (" {
				inImport = true
				continue
			}
			if inImport {
				if trim == ")" {
					inImport = false
					continue
				}
				trim = strings.Trim(trim, `"`)
				if trim != "" {
					imports[trim] = struct{}{}
				}
				continue
			}

			if strings.HasPrefix(trim, "func MigrateAll") {
				inMigrate = true
				continue
			}
			if inMigrate && strings.HasPrefix(trim, "&") {
				models[strings.TrimRight(trim, ",")] = struct{}{}
			}
		}
	}

	// Add new entries
	imports[importPath] = struct{}{}
	models[modelLine] = struct{}{}

	// Sort imports and models
	importList := make([]string, 0, len(imports))
	for imp := range imports {
		importList = append(importList, imp)
	}
	sort.Strings(importList)

	modelList := make([]string, 0, len(models))
	for model := range models {
		modelList = append(modelList, model)
	}
	sort.Strings(modelList)

	// Build output
	var builder strings.Builder
	fmt.Fprintf(&builder, "package migrations\n\n")
	fmt.Fprintf(&builder, "import (\n")
	for _, imp := range importList {
		fmt.Fprintf(&builder, "\t\"%s\"\n", imp)
	}
	fmt.Fprintf(&builder, ")\n\n")

	fmt.Fprintf(&builder, "func MigrateAll(db *gorm.DB) error {\n")
	fmt.Fprintf(&builder, "\treturn db.AutoMigrate(\n")
	for _, model := range modelList {
		fmt.Fprintf(&builder, "\t\t%s,\n", model)
	}
	fmt.Fprintf(&builder, "\t)\n")
	fmt.Fprintf(&builder, "}\n")

	// Write to file
	if err := os.WriteFile(filePath, []byte(builder.String()), 0644); err != nil {
		fmt.Println("Failed to write migrations.go:", err)
		return
	}

	fmt.Println("migrations/migrations.go updated with model:", modelLine)
}
