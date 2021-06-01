
package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
)



type server_config struct {
  Port string //`mapstructure:"port"`
}

type db_config struct {
	H_name string //`mapstructure:"hostname"`
	Usr string //`mapstructure: "user"`
	Password string // `mapstructure: "password"`
	Database_name string // `mapstructure:"database_name"`
	Port string //`mapstructure:"port"`
}




var(
  servrcnfig server_config
  dbcnfig db_config
  )


var serverCmd = &cobra.Command{
	Use:   "myapi",
	Short: "This command will turn on myAPI",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		connect()

		r := chi.NewRouter()

		r.Post("/student", Add_New_Student)

		r.Get("/student", Get_All_Students)

		r.Get("/student/{id}", Get_Student_With_ID)

		r.Put("/student/{id}", Update_Student_With_ID)

		r.Delete("/student/{id}", Delete_Student_with_ID)

		err := http.ListenAndServe(":"+servrcnfig.Port, r)

		//changin variable to config
		if err != nil {
			panic(err)
		}
	},
}

var (
	db *gorm.DB
)

type STUDENT_INFO struct {
	S_ID     string `json:"S_ID"`
	Name     string `json:"Name"`
	Village  string `json:"Village"`
	Thana    string `json:"Thana"`
	District string `json:"District"`
}

func connect() {

	//dsn := "host=database user=admin password=secret dbname=database_name port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	dsn:=fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",dbcnfig.H_name,dbcnfig.Usr,dbcnfig.Password,dbcnfig.Database_name,dbcnfig.Port)
	fmt.Println(dsn)
	d, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("err")
	} else {
		fmt.Println("A B C conected")
		fmt.Printf("check")
	}

	db = d

	db.AutoMigrate(&STUDENT_INFO{})
}

func Delete_Student_with_ID(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("content-Type", "application/json")

	var g_student, tm_student []STUDENT_INFO

	ID := chi.URLParam(r, "id")
	db.Where("S_ID = ?", ID).Delete(&g_student)

	err := db.Find(&tm_student).Error

	if err != nil {
		panic(err)
	}

	json.NewEncoder(w).Encode(tm_student)

}

func Add_New_Student(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("content-Type", "application/json")

	var s_student []STUDENT_INFO

	var temp1student STUDENT_INFO

	err := json.NewDecoder(r.Body).Decode(&temp1student)

	db.Create(&temp1student)

	if err != nil {
		panic(err)
	}

	db.Find(&s_student)

	json.NewEncoder(w).Encode(s_student)
}

func Get_All_Students(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var g_student []STUDENT_INFO

	db.Find(&g_student)

	json.NewEncoder(w).Encode(g_student)
}

func Update_Student_With_ID(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("content-Type", "application/json")

	var g_student []STUDENT_INFO
	var p STUDENT_INFO

	ID := chi.URLParam(r, "id")

	var temp STUDENT_INFO
	err := json.NewDecoder(r.Body).Decode(&temp)
	if err != nil {
		fmt.Println(err)
	}

	db.Model(&p).Select("*").Where("S_ID=?", ID).Updates(STUDENT_INFO{S_ID: temp.S_ID, Name: temp.Name, Village: temp.Village, Thana: temp.Thana, District: temp.District})

	db.Find(&g_student)

	json.NewEncoder(w).Encode(g_student)
}

func Get_Student_With_ID(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("content-Type", "application/json")

	var g_student []STUDENT_INFO
	fmt.Println("bug000")
	db.Find(&g_student)

	ID := chi.URLParam(r, "id")
	fmt.Println("Bug1")
	for _, item := range g_student {
		if item.S_ID == ID {
			fmt.Println("Bug2")

			var temp STUDENT_INFO
			temp.S_ID = item.S_ID
			temp.Name = item.Name
			temp.Village = item.Village
			temp.Thana = item.Thana
			temp.District = item.District

			json.NewEncoder(w).Encode(temp)
			return
		}
	}
}




func init() {
	rootCmd.AddCommand(serverCmd)
	cobra.OnInitialize(initconfig1)
	cobra.OnInitialize(initconfig2)
}


func  initconfig1()  {

    viper.SetConfigName("conf1")
    viper.AddConfigPath(".")
    viper.AutomaticEnv()
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	viper.Unmarshal(&dbcnfig)

	fmt.Println("init1start")

	fmt.Println(dbcnfig.H_name)
	fmt.Println(dbcnfig.Password)
	fmt.Println(dbcnfig.Port)
	fmt.Println(dbcnfig.Usr)


	fmt.Println("init1end")


}

func initconfig2()  {
	fmt.Println("init2start")
	viper.SetConfigName("conf")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}
	//servrcnfig.Port=viper.GetString("port")

	viper.Unmarshal(&servrcnfig)
	fmt.Println("init1")
	fmt.Println(servrcnfig.Port)

	fmt.Println("init2end")

}
